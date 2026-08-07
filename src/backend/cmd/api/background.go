package main

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/avito-hack/backend/internal/module/pet"
	"github.com/avito-hack/backend/internal/shared/kafka"
)

const defaultFlushInterval = 15 * time.Second

type background struct {
	consumer *kafka.Consumer
	producer *kafka.Producer
	petModul *pet.Module
	log      *slog.Logger
}

func (b *background) run(ctx context.Context) func() {
	var wg sync.WaitGroup

	if b.consumer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := b.consumer.Run(ctx); err != nil && ctx.Err() == nil {
				b.log.Error("kafka consumer stopped", slog.Any("error", err))
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		b.flushLoop(ctx)
	}()

	return func() {
		if b.consumer != nil {
			b.consumer.Close()
		}
		if b.producer != nil {
			b.producer.Close()
		}
		wg.Wait()
	}
}

func (b *background) flushLoop(ctx context.Context) {
	interval := b.petModul.FlushInterval
	if interval <= 0 {
		interval = defaultFlushInterval
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := b.petModul.Flusher.Handle(ctx); err != nil && ctx.Err() == nil {
				b.log.Warn("flush pet hot state", slog.Any("error", err))
			}
		}
	}
}
