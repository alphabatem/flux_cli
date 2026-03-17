package services

import (
	"fmt"

	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/pkg/context"
)

const RUGCHECK_SVC = "rugcheck_svc"

type RugCheckService struct {
	context.DefaultService

	Client *httpclient.Client
}

func (s *RugCheckService) Id() string {
	return RUGCHECK_SVC
}

func (s *RugCheckService) Configure(ctx *context.Context) error {
	cfg := ctx.Service(CONFIG_SVC).(*ConfigService).Config()

	s.Client = httpclient.New(
		cfg.RugCheck.BaseURL,
		cfg.RugCheck.APIKey,
		"X-API-KEY",
	)

	return s.DefaultService.Configure(ctx)
}

// GetReport returns a detailed security report for a token.
func (s *RugCheckService) GetReport(mint string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/tokens/%s/report", mint), &result)
	return result, err
}

// GetSummary returns a summary security report for a token.
func (s *RugCheckService) GetSummary(mint string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/tokens/%s/report/summary", mint), &result)
	return result, err
}

// Scan performs a real-time token security scan.
func (s *RugCheckService) Scan(chain, contractAddress string, includeDexScreener bool) (interface{}, error) {
	path := fmt.Sprintf("/tokens/scan/%s/%s", chain, contractAddress)
	if includeDexScreener {
		path += "?includeDexScreenerData=true"
	}
	var result interface{}
	err := s.Client.Get(path, &result)
	return result, err
}

// Search searches tokens by name or symbol.
func (s *RugCheckService) Search(query, chain string, pageSize int) (interface{}, error) {
	path := fmt.Sprintf("/tokens/search?query=%s", query)
	if chain != "" {
		path += "&chain=" + chain
	}
	if pageSize > 0 {
		path += fmt.Sprintf("&pageSize=%d", pageSize)
	}
	var result interface{}
	err := s.Client.Get(path, &result)
	return result, err
}

// GetStatsNew returns recently detected tokens.
func (s *RugCheckService) GetStatsNew() (interface{}, error) {
	var result interface{}
	err := s.Client.Get("/stats/new_tokens", &result)
	return result, err
}

// GetStatsRecent returns most viewed tokens.
func (s *RugCheckService) GetStatsRecent() (interface{}, error) {
	var result interface{}
	err := s.Client.Get("/stats/recent", &result)
	return result, err
}

// GetStatsTrending returns trending tokens.
func (s *RugCheckService) GetStatsTrending() (interface{}, error) {
	var result interface{}
	err := s.Client.Get("/stats/trending", &result)
	return result, err
}

// GetStatsVerified returns verified tokens.
func (s *RugCheckService) GetStatsVerified() (interface{}, error) {
	var result interface{}
	err := s.Client.Get("/stats/verified", &result)
	return result, err
}

// GetWalletRisk returns a wallet risk assessment.
func (s *RugCheckService) GetWalletRisk(chain, walletAddress string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/wallets/risk-rating/%s/%s", chain, walletAddress), &result)
	return result, err
}

// GetChains returns all supported chains.
func (s *RugCheckService) GetChains() (interface{}, error) {
	var result interface{}
	err := s.Client.Get("/utils/chains", &result)
	return result, err
}
