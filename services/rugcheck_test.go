package services

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alphabatem/flux_cli/internal/httpclient"
)

func newTestRugCheckService(t *testing.T, capturedPath *string) (*RugCheckService, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if capturedPath != nil {
			*capturedPath = r.URL.RequestURI()
		}
		// Verify RugCheck uses X-API-KEY header
		if r.Header.Get("X-API-KEY") == "" && r.Header.Get("Authorization") != "" {
			t.Error("RugCheck should use X-API-KEY, not Authorization")
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"data":"ok"}`))
	}))

	svc := &RugCheckService{
		Client: httpclient.New(srv.URL, "rc-key", "X-API-KEY"),
	}
	return svc, srv
}

func TestRugCheck_GetReport(t *testing.T) {
	var path string
	svc, srv := newTestRugCheckService(t, &path)
	defer srv.Close()

	svc.GetReport("mint1")
	if path != "/tokens/mint1/report" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestRugCheck_GetSummary(t *testing.T) {
	var path string
	svc, srv := newTestRugCheckService(t, &path)
	defer srv.Close()

	svc.GetSummary("mint1")
	if path != "/tokens/mint1/report/summary" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestRugCheck_Scan(t *testing.T) {
	var path string
	svc, srv := newTestRugCheckService(t, &path)
	defer srv.Close()

	svc.Scan("solana", "addr1", false)
	if path != "/tokens/scan/solana/addr1" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestRugCheck_Scan_WithDexScreener(t *testing.T) {
	var path string
	svc, srv := newTestRugCheckService(t, &path)
	defer srv.Close()

	svc.Scan("solana", "addr1", true)
	if !strings.Contains(path, "includeDexScreenerData=true") {
		t.Errorf("expected dexscreener param, got: %s", path)
	}
}

func TestRugCheck_Search(t *testing.T) {
	var path string
	svc, srv := newTestRugCheckService(t, &path)
	defer srv.Close()

	svc.Search("BONK", "", 0)
	if path != "/tokens/search?query=BONK" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestRugCheck_Search_WithChainAndPageSize(t *testing.T) {
	var path string
	svc, srv := newTestRugCheckService(t, &path)
	defer srv.Close()

	svc.Search("SOL", "solana", 20)
	if !strings.Contains(path, "query=SOL") {
		t.Errorf("expected query, got: %s", path)
	}
	if !strings.Contains(path, "chain=solana") {
		t.Errorf("expected chain, got: %s", path)
	}
	if !strings.Contains(path, "pageSize=20") {
		t.Errorf("expected pageSize, got: %s", path)
	}
}

func TestRugCheck_GetStatsNew(t *testing.T) {
	var path string
	svc, srv := newTestRugCheckService(t, &path)
	defer srv.Close()

	svc.GetStatsNew()
	if path != "/stats/new_tokens" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestRugCheck_GetStatsRecent(t *testing.T) {
	var path string
	svc, srv := newTestRugCheckService(t, &path)
	defer srv.Close()

	svc.GetStatsRecent()
	if path != "/stats/recent" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestRugCheck_GetStatsTrending(t *testing.T) {
	var path string
	svc, srv := newTestRugCheckService(t, &path)
	defer srv.Close()

	svc.GetStatsTrending()
	if path != "/stats/trending" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestRugCheck_GetStatsVerified(t *testing.T) {
	var path string
	svc, srv := newTestRugCheckService(t, &path)
	defer srv.Close()

	svc.GetStatsVerified()
	if path != "/stats/verified" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestRugCheck_GetWalletRisk(t *testing.T) {
	var path string
	svc, srv := newTestRugCheckService(t, &path)
	defer srv.Close()

	svc.GetWalletRisk("solana", "wallet1")
	if path != "/wallets/risk-rating/solana/wallet1" {
		t.Errorf("unexpected path: %s", path)
	}
}

func TestRugCheck_GetChains(t *testing.T) {
	var path string
	svc, srv := newTestRugCheckService(t, &path)
	defer srv.Close()

	svc.GetChains()
	if path != "/utils/chains" {
		t.Errorf("unexpected path: %s", path)
	}
}
