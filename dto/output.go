package dto

// CLIResponse is the standard envelope for all CLI output.
// AI agents can rely on this consistent structure.
type CLIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   *CLIError   `json:"error"`
}

type CLIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type CLIMeta struct {
	Service    string `json:"service,omitempty"`
	Endpoint   string `json:"endpoint,omitempty"`
	DurationMs int64  `json:"duration_ms,omitempty"`
}

// Exit codes for programmatic error handling
const (
	ExitSuccess      = 0
	ExitGeneralError = 1
	ExitUsageError   = 2
	ExitAPIError     = 64
	ExitAuthError    = 65
	ExitServiceDown  = 69
	ExitConfigError  = 78
)
