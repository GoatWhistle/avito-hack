package infra

import (
	"bufio"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type fakeRedis struct {
	mu       sync.Mutex
	replies  map[string]string
	fallback string
}

func respArray(items ...string) string {
	var b strings.Builder
	b.WriteString("*" + strconv.Itoa(len(items)) + "\r\n")
	for _, item := range items {
		b.WriteString(item)
	}

	return b.String()
}

func respBulk(value string) string {
	return "$" + strconv.Itoa(len(value)) + "\r\n" + value + "\r\n"
}

func respNil() string { return "$-1\r\n" }

func respInt(value int64) string { return ":" + strconv.FormatInt(value, 10) + "\r\n" }

func startFakeRedis(t *testing.T, replies map[string]string) *redis.Client {
	t.Helper()

	full := map[string]string{
		"HELLO": "-ERR unknown command 'HELLO'\r\n",
		"AUTH":  "+OK\r\n",
		"PING":  "+PONG\r\n",
	}
	for name, reply := range replies {
		full[name] = reply
	}
	server := &fakeRedis{replies: full, fallback: "+OK\r\n"}

	var config net.ListenConfig
	listener, err := config.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}
			go server.serve(conn)
		}
	}()

	client := redis.NewClient(&redis.Options{
		Addr:         listener.Addr().String(),
		Protocol:     2,
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		MaxRetries:   -1,
	})

	t.Cleanup(func() {
		_ = client.Close()
		_ = listener.Close()
		<-done
	})

	return client
}

func (f *fakeRedis) serve(conn net.Conn) {
	defer func() { _ = conn.Close() }()

	reader := bufio.NewReader(conn)
	for {
		command, err := readCommand(reader)
		if err != nil {
			return
		}
		if _, err := conn.Write([]byte(f.reply(command))); err != nil {
			return
		}
	}
}

func (f *fakeRedis) reply(command []string) string {
	f.mu.Lock()
	defer f.mu.Unlock()

	name := strings.ToUpper(command[0])
	if reply, ok := f.replies[name]; ok {
		return reply
	}

	return f.fallback
}

func arrayCount(header string) (count int, ok bool) {
	parsed, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(header, "*")))
	if err != nil || parsed <= 0 {
		return 0, false
	}

	return parsed, true
}

func readCommand(reader *bufio.Reader) ([]string, error) {
	header, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	count, ok := arrayCount(header)
	if !ok {
		return []string{"UNKNOWN"}, nil
	}

	parts := make([]string, 0, count)
	for range count {
		if _, err := reader.ReadString('\n'); err != nil {
			return nil, err
		}
		value, readErr := reader.ReadString('\n')
		if readErr != nil {
			return nil, readErr
		}
		parts = append(parts, strings.TrimRight(value, "\r\n"))
	}

	return parts, nil
}
