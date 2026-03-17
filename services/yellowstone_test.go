package services

import (
	"testing"

	"github.com/alphabatem/flux_cli/dto"
)

func TestResolveYellowstoneURL(t *testing.T) {
	tests := []struct {
		name string
		cfg  dto.FluxRPCConfig
		want string
	}{
		{
			name: "us region",
			cfg:  dto.FluxRPCConfig{Region: "us"},
			want: "https://yellowstone.us.fluxrpc.com",
		},
		{
			name: "eu region",
			cfg:  dto.FluxRPCConfig{Region: "eu"},
			want: "https://yellowstone.eu.fluxrpc.com",
		},
		{
			name: "default region",
			cfg:  dto.FluxRPCConfig{},
			want: "https://yellowstone.us.fluxrpc.com",
		},
	}

	for _, tt := range tests {
		got := ResolveYellowstoneURL(&tt.cfg)
		if got != tt.want {
			t.Fatalf("%s: got %q want %q", tt.name, got, tt.want)
		}
	}
}
