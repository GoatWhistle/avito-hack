package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectPhotoType(t *testing.T) {
	t.Parallel()

	jpeg := append([]byte{0xFF, 0xD8, 0xFF}, make([]byte, 32)...)
	png := append([]byte("\x89PNG\r\n\x1a\n"), make([]byte, 32)...)
	webp := append([]byte("RIFF\x00\x00\x00\x00WEBPVP8 "), make([]byte, 32)...)
	gif := append([]byte("GIF89a"), make([]byte, 32)...)
	text := []byte("this is definitely not an image, just plain ascii text content")

	tests := []struct {
		name    string
		raw     []byte
		want    string
		wantErr bool
	}{
		{name: "jpeg", raw: jpeg, want: "image/jpeg"},
		{name: "png", raw: png, want: "image/png"},
		{name: "webp", raw: webp, want: "image/webp"},
		{name: "gif is rejected", raw: gif, wantErr: true},
		{name: "text renamed to jpg is rejected", raw: text, wantErr: true},
		{name: "empty is rejected", raw: nil, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := detectPhotoType(tt.raw)

			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
