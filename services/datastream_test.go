package services

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alphabatem/flux_cli/internal/httpclient"
)

// newTestDataStreamService creates a DataStreamService that captures request paths.
func newTestDataStreamService(t *testing.T, capturedPath *string) (*DataStreamService, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if capturedPath != nil {
			*capturedPath = r.URL.RequestURI()
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"data":"ok"}`))
	}))

	svc := &DataStreamService{
		Client: httpclient.New(srv.URL, "test-key", "Authorization"),
	}
	return svc, srv
}

func TestDataStream_GetPrices(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetPrices("mint1,mint2,mint3")
	if path != "/prices?ids=mint1,mint2,mint3" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetTokens(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokens("", 50)
	if path != "/tokens?limit=50" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetTokens_WithCursor(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokens("cursor123", 100)
	if !strings.Contains(path, "cursor=cursor123") {
		t.Errorf("expected cursor in path, got: %s", path)
	}
	if !strings.Contains(path, "limit=100") {
		t.Errorf("expected limit in path, got: %s", path)
	}
}

func TestDataStream_GetToken(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetToken("So11111111111111111111111111111111111111112")
	if path != "/tokens/So11111111111111111111111111111111111111112" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetTokenDetails(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokenDetails("mint1")
	if path != "/tokens/mint1/details" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetTokenCandles(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokenCandles("mint1", "15", 1000, 2000, 50)
	if !strings.Contains(path, "/tokens/mint1/candles") {
		t.Errorf("expected candles path, got: %s", path)
	}
	if !strings.Contains(path, "interval=15") {
		t.Errorf("expected interval, got: %s", path)
	}
	if !strings.Contains(path, "count=50") {
		t.Errorf("expected count, got: %s", path)
	}
	if !strings.Contains(path, "from=1000") {
		t.Errorf("expected from, got: %s", path)
	}
	if !strings.Contains(path, "to=2000") {
		t.Errorf("expected to, got: %s", path)
	}
}

func TestDataStream_GetTokenCandles_NoTimestamps(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokenCandles("mint1", "5", 0, 0, 100)
	if strings.Contains(path, "from=") {
		t.Errorf("from should not be in path when 0, got: %s", path)
	}
	if strings.Contains(path, "to=") {
		t.Errorf("to should not be in path when 0, got: %s", path)
	}
}

func TestDataStream_GetTokenHolders(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokenHolders("mint1", 50, 2)
	if path != "/tokens/mint1/holders?limit=50&page=2" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetTokenHoldersTop(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokenHoldersTop("mint1", 200)
	if path != "/tokens/mint1/holders/top?limit=200" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetTokenHoldersDistribution(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokenHoldersDistribution("mint1")
	if path != "/tokens/mint1/holders/distribution" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetTokenHoldersCount(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokenHoldersCount("mint1")
	if path != "/tokens/mint1/holders/count" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetTokenTrades(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokenTrades("mint1", 25, 100, 200)
	if !strings.Contains(path, "limit=25") {
		t.Errorf("expected limit, got: %s", path)
	}
	if !strings.Contains(path, "from=100") {
		t.Errorf("expected from, got: %s", path)
	}
}

func TestDataStream_GetTokenTrades_NoTimestamps(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokenTrades("mint1", 50, 0, 0)
	if strings.Contains(path, "from=") || strings.Contains(path, "to=") {
		t.Errorf("timestamps should not be in path when 0, got: %s", path)
	}
}

func TestDataStream_GetTokenPrice(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokenPrice("mint1")
	if path != "/tokens/mint1/price" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetTokenMcap(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTokenMcap("mint1")
	if path != "/tokens/mint1/mcap" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetStatsNew(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetStatsNew()
	if path != "/stats/tokens/new" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetStatsTop(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetStatsTop(20, "60")
	if path != "/stats/tokens/top?limit=20&interval=60" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetStatsTrending(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetStatsTrending(50, "30")
	if path != "/stats/tokens/trending?limit=50&interval=30" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetTradersTop(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTradersTop(10, 1000, 2000)
	if !strings.Contains(path, "/stats/traders/top") {
		t.Errorf("expected traders/top, got: %s", path)
	}
	if !strings.Contains(path, "from=1000") {
		t.Errorf("expected from, got: %s", path)
	}
}

func TestDataStream_GetTraderDetail(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTraderDetail("wallet1")
	if path != "/traders/wallet1/detail" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetTraderPnL(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetTraderPnL("wallet1", 50)
	if path != "/traders/wallet1/tokens/pnl?limit=50" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestDataStream_GetBulkCandles(t *testing.T) {
	var path string
	svc, srv := newTestDataStreamService(t, &path)
	defer srv.Close()

	svc.GetBulkCandles("m1,m2", "15", 1000, 2000, 50)
	if !strings.Contains(path, "/bulk/tokens/candles") {
		t.Errorf("expected bulk candles path, got: %s", path)
	}
	if !strings.Contains(path, "ids=m1,m2") {
		t.Errorf("expected ids, got: %s", path)
	}
}
