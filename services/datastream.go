package services

import (
	"fmt"

	"github.com/alphabatem/flux_cli/internal/httpclient"
	"github.com/alphabatem/flux_cli/pkg/context"
)

const DATASTREAM_SVC = "datastream_svc"

type DataStreamService struct {
	context.DefaultService

	Client *httpclient.Client
}

func (s *DataStreamService) Id() string {
	return DATASTREAM_SVC
}

func (s *DataStreamService) Configure(ctx *context.Context) error {
	cfg := ctx.Service(CONFIG_SVC).(*ConfigService).Config()

	s.Client = httpclient.New(
		cfg.DataStream.BaseURL,
		cfg.DataStream.APIKey,
		"Authorization",
	)

	return s.DefaultService.Configure(ctx)
}

// GetPrices returns bulk token prices for the given comma-separated mint IDs.
func (s *DataStreamService) GetPrices(ids string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/prices?ids=%s", ids), &result)
	return result, err
}

// GetTokens returns a paginated list of tokens.
func (s *DataStreamService) GetTokens(cursor string, limit int) (interface{}, error) {
	path := fmt.Sprintf("/tokens?limit=%d", limit)
	if cursor != "" {
		path += "&cursor=" + cursor
	}
	var result interface{}
	err := s.Client.Get(path, &result)
	return result, err
}

// GetToken returns details for a specific token.
func (s *DataStreamService) GetToken(id string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/tokens/%s", id), &result)
	return result, err
}

// GetTokenDetails returns extended token info with price and market cap.
func (s *DataStreamService) GetTokenDetails(id string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/tokens/%s/details", id), &result)
	return result, err
}

// GetTokenCandles returns OHLCV candle data for a token.
func (s *DataStreamService) GetTokenCandles(id string, interval string, from, to int64, count int) (interface{}, error) {
	path := fmt.Sprintf("/tokens/%s/candles?interval=%s&count=%d", id, interval, count)
	if from > 0 {
		path += fmt.Sprintf("&from=%d", from)
	}
	if to > 0 {
		path += fmt.Sprintf("&to=%d", to)
	}
	var result interface{}
	err := s.Client.Get(path, &result)
	return result, err
}

// GetTokenHolders returns holders for a token.
func (s *DataStreamService) GetTokenHolders(id string, limit, page int) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/tokens/%s/holders?limit=%d&page=%d", id, limit, page), &result)
	return result, err
}

// GetTokenHoldersTop returns top holders for a token.
func (s *DataStreamService) GetTokenHoldersTop(id string, limit int) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/tokens/%s/holders/top?limit=%d", id, limit), &result)
	return result, err
}

// GetTokenHoldersDistribution returns holder distribution for a token.
func (s *DataStreamService) GetTokenHoldersDistribution(id string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/tokens/%s/holders/distribution", id), &result)
	return result, err
}

// GetTokenHoldersCount returns the holder count for a token.
func (s *DataStreamService) GetTokenHoldersCount(id string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/tokens/%s/holders/count", id), &result)
	return result, err
}

// GetTokenTraders returns trader stats for a token.
func (s *DataStreamService) GetTokenTraders(id string, limit, page int, all bool) (interface{}, error) {
	path := fmt.Sprintf("/tokens/%s/traders?limit=%d&page=%d", id, limit, page)
	if all {
		path += "&all=true"
	}
	var result interface{}
	err := s.Client.Get(path, &result)
	return result, err
}

// GetTokenTrades returns trade history for a token.
func (s *DataStreamService) GetTokenTrades(id string, limit int, from, to int64) (interface{}, error) {
	path := fmt.Sprintf("/tokens/%s/trades?limit=%d", id, limit)
	if from > 0 {
		path += fmt.Sprintf("&from=%d", from)
	}
	if to > 0 {
		path += fmt.Sprintf("&to=%d", to)
	}
	var result interface{}
	err := s.Client.Get(path, &result)
	return result, err
}

// GetTokenPrice returns the current price for a token.
func (s *DataStreamService) GetTokenPrice(id string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/tokens/%s/price", id), &result)
	return result, err
}

// GetTokenMcap returns the market cap for a token.
func (s *DataStreamService) GetTokenMcap(id string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/tokens/%s/mcap", id), &result)
	return result, err
}

// GetTokenStats returns volume change stats for a token.
func (s *DataStreamService) GetTokenStats(id string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/tokens/%s/stats/change", id), &result)
	return result, err
}

// GetStatsNew returns newly detected tokens.
func (s *DataStreamService) GetStatsNew() (interface{}, error) {
	var result interface{}
	err := s.Client.Get("/stats/tokens/new", &result)
	return result, err
}

// GetStatsTop returns top tokens by volume.
func (s *DataStreamService) GetStatsTop(limit int, interval string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/stats/tokens/top?limit=%d&interval=%s", limit, interval), &result)
	return result, err
}

// GetStatsTrending returns trending tokens.
func (s *DataStreamService) GetStatsTrending(limit int, interval string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/stats/tokens/trending?limit=%d&interval=%s", limit, interval), &result)
	return result, err
}

// GetStatsVolume returns tokens by volume.
func (s *DataStreamService) GetStatsVolume(limit int) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/stats/tokens/volume?limit=%d", limit), &result)
	return result, err
}

// GetTradersTop returns top traders.
func (s *DataStreamService) GetTradersTop(limit int, from, to int64) (interface{}, error) {
	path := fmt.Sprintf("/stats/traders/top?limit=%d", limit)
	if from > 0 {
		path += fmt.Sprintf("&from=%d", from)
	}
	if to > 0 {
		path += fmt.Sprintf("&to=%d", to)
	}
	var result interface{}
	err := s.Client.Get(path, &result)
	return result, err
}

// GetTraderDetail returns details for a specific trader.
func (s *DataStreamService) GetTraderDetail(id string) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/traders/%s/detail", id), &result)
	return result, err
}

// GetTraderTokens returns token holdings for a trader.
func (s *DataStreamService) GetTraderTokens(id string, limit int) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/traders/%s/tokens?limit=%d", id, limit), &result)
	return result, err
}

// GetTraderPnL returns PnL for a trader.
func (s *DataStreamService) GetTraderPnL(id string, limit int) (interface{}, error) {
	var result interface{}
	err := s.Client.Get(fmt.Sprintf("/traders/%s/tokens/pnl?limit=%d", id, limit), &result)
	return result, err
}

// GetTraderTrades returns trade history for a trader.
func (s *DataStreamService) GetTraderTrades(id string, limit int, from, to int64) (interface{}, error) {
	path := fmt.Sprintf("/traders/%s/trades?limit=%d", id, limit)
	if from > 0 {
		path += fmt.Sprintf("&from=%d", from)
	}
	if to > 0 {
		path += fmt.Sprintf("&to=%d", to)
	}
	var result interface{}
	err := s.Client.Get(path, &result)
	return result, err
}

// GetBulkCandles returns bulk candle data for multiple tokens.
func (s *DataStreamService) GetBulkCandles(ids string, interval string, from, to int64, count int) (interface{}, error) {
	path := fmt.Sprintf("/bulk/tokens/candles?ids=%s&interval=%s&count=%d", ids, interval, count)
	if from > 0 {
		path += fmt.Sprintf("&from=%d", from)
	}
	if to > 0 {
		path += fmt.Sprintf("&to=%d", to)
	}
	var result interface{}
	err := s.Client.Get(path, &result)
	return result, err
}
