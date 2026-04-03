package dto

type Config struct {
	DataStream DataStreamConfig `json:"datastream"`
	FluxRPC    FluxRPCConfig    `json:"fluxrpc"`
	RugCheck   RugCheckConfig   `json:"rugcheck"`
	Output     OutputConfig     `json:"output"`
}

type DataStreamConfig struct {
	APIKey  string `json:"api_key,omitempty"`
	BaseURL string `json:"base_url,omitempty"`
}

type FluxRPCConfig struct {
	APIKey  string `json:"api_key,omitempty"`
	BaseURL string `json:"base_url,omitempty"`
	Region  string `json:"region,omitempty"` // "eu" or "us" — overrides base_url
}

type RugCheckConfig struct {
	APIKey  string `json:"api_key,omitempty"`
	BaseURL string `json:"base_url,omitempty"`
}

type OutputConfig struct {
	Format string `json:"format,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		DataStream: DataStreamConfig{
			BaseURL: "https://data.fluxbeam.xyz",
		},
		FluxRPC: FluxRPCConfig{
			Region: "us",
		},
		RugCheck: RugCheckConfig{
			BaseURL: "https://api.rugcheck.xyz/v1",
		},
		Output: OutputConfig{
			Format: "json",
		},
	}
}
