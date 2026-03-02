package server

import (
	"strings"
	"testing"
)

// TestGenerateASCII tests the core ASCII generation logic
func TestGenerateASCII(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		banner     string
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "valid request",
			text:       "Hello",
			banner:     "standard",
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name:       "default banner",
			text:       "Hello",
			banner:     "",
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name:       "empty text",
			text:       "",
			banner:     "standard",
			wantStatus: 400,
			wantErr:    true,
		},
		{
			name:       "text too long",
			text:       strings.Repeat("a", 1001),
			banner:     "standard",
			wantStatus: 400,
			wantErr:    true,
		},
		{
			name:       "invalid banner",
			text:       "Hello",
			banner:     "invalid",
			wantStatus: 404,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, status, err := GenerateASCII(tt.text, tt.banner)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateASCII() error = %v, wantErr %v", err, tt.wantErr)
			}

			if status != tt.wantStatus {
				t.Errorf("GenerateASCII() status = %v, want %v", status, tt.wantStatus)
			}

			if !tt.wantErr && result == "" {
				t.Errorf("GenerateASCII() returned empty result for valid input")
			}
		})
	}
}
