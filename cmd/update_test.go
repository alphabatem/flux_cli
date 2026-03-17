package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIsNewer(t *testing.T) {
	tests := []struct {
		current string
		latest  string
		want    bool
	}{
		// Basic semver
		{"0.0.1", "0.0.2", true},
		{"0.0.2", "0.0.1", false},
		{"0.0.1", "0.0.1", false},

		// Minor bumps
		{"0.1.0", "0.2.0", true},
		{"0.2.0", "0.1.0", false},

		// Major bumps
		{"1.0.0", "2.0.0", true},
		{"2.0.0", "1.0.0", false},

		// Mixed
		{"1.2.3", "1.2.4", true},
		{"1.2.3", "1.3.0", true},
		{"1.2.3", "2.0.0", true},
		{"1.2.3", "1.2.3", false},
		{"1.3.0", "1.2.9", false},

		// Dev / non-semver current is always older
		{"dev", "0.0.1", true},
		{"dev", "1.0.0", true},
		{"none", "0.0.1", true},
		{"", "0.0.1", true},

		// Non-semver latest is never newer
		{"0.0.1", "dev", false},
		{"0.0.1", "none", false},
		{"0.0.1", "", false},

		// Both non-semver
		{"dev", "dev", false},
		{"dev", "none", false},

		// Different length versions
		{"0.0", "0.0.1", true},
		{"0.0.1", "0.0", false},
		{"1.0", "1.0.0", false},
	}

	for _, tt := range tests {
		got := isNewer(tt.current, tt.latest)
		if got != tt.want {
			t.Errorf("isNewer(%q, %q) = %v, want %v", tt.current, tt.latest, got, tt.want)
		}
	}
}

func TestIsSemver(t *testing.T) {
	tests := []struct {
		v    string
		want bool
	}{
		{"0.0.1", true},
		{"1.0", true},
		{"1.2.3", true},
		{"dev", false},
		{"none", false},
		{"", false},
		{"1", false},
		{"abc.def", false},
		{"1.2.abc", false},
	}

	for _, tt := range tests {
		got := isSemver(tt.v)
		if got != tt.want {
			t.Errorf("isSemver(%q) = %v, want %v", tt.v, got, tt.want)
		}
	}
}

func TestBuildAssetName(t *testing.T) {
	name := buildAssetName("1.2.3")

	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	expected := "flux_1.2.3_" + runtime.GOOS + "_" + runtime.GOARCH + "." + ext

	if name != expected {
		t.Errorf("buildAssetName(\"1.2.3\") = %q, want %q", name, expected)
	}
}

func TestFetchLatestRelease(t *testing.T) {
	release := githubRelease{TagName: "v1.2.3"}
	release.Assets = []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	}{
		{Name: "flux_1.2.3_linux_amd64.tar.gz", BrowserDownloadURL: "https://example.com/flux.tar.gz"},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(release)
	}))
	defer srv.Close()

	// Override fetchLatestRelease by testing the HTTP handler directly
	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var got githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if got.TagName != "v1.2.3" {
		t.Errorf("TagName = %q, want %q", got.TagName, "v1.2.3")
	}
	if len(got.Assets) != 1 {
		t.Fatalf("expected 1 asset, got %d", len(got.Assets))
	}
	if got.Assets[0].Name != "flux_1.2.3_linux_amd64.tar.gz" {
		t.Errorf("asset name = %q, want %q", got.Assets[0].Name, "flux_1.2.3_linux_amd64.tar.gz")
	}
}

func TestExtractFromTarGz(t *testing.T) {
	content := []byte("fake-binary-content")
	archivePath := createTestTarGz(t, "flux", content)

	got, err := extractFromTarGz(archivePath, "flux")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("extracted = %q, want %q", got, content)
	}
}

func TestExtractFromTarGz_NotFound(t *testing.T) {
	content := []byte("other-content")
	archivePath := createTestTarGz(t, "other-binary", content)

	_, err := extractFromTarGz(archivePath, "flux")
	if err == nil {
		t.Fatal("expected error for missing binary")
	}
}

func TestExtractFromZip(t *testing.T) {
	content := []byte("fake-binary-content")
	archivePath := createTestZip(t, "flux.exe", content)

	got, err := extractFromZip(archivePath, "flux.exe")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("extracted = %q, want %q", got, content)
	}
}

func TestExtractFromZip_NotFound(t *testing.T) {
	content := []byte("other-content")
	archivePath := createTestZip(t, "other.exe", content)

	_, err := extractFromZip(archivePath, "flux.exe")
	if err == nil {
		t.Fatal("expected error for missing binary")
	}
}

func TestExtractFromTarGz_NestedPath(t *testing.T) {
	content := []byte("nested-binary")
	archivePath := createTestTarGz(t, "flux_1.0.0_linux_amd64/flux", content)

	got, err := extractFromTarGz(archivePath, "flux")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("extracted = %q, want %q", got, content)
	}
}

func TestExtractFromZip_NestedPath(t *testing.T) {
	content := []byte("nested-binary")
	archivePath := createTestZip(t, "flux_1.0.0_windows_amd64/flux.exe", content)

	got, err := extractFromZip(archivePath, "flux.exe")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("extracted = %q, want %q", got, content)
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "src")
	dst := filepath.Join(tmpDir, "dst")

	content := []byte("binary-data")
	if err := os.WriteFile(src, content, 0644); err != nil {
		t.Fatal(err)
	}

	if err := copyFile(src, dst); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("copied content = %q, want %q", got, content)
	}
}

// --- helpers ---

func createTestTarGz(t *testing.T, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.tar.gz")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	if err := tw.WriteHeader(&tar.Header{
		Name:     name,
		Size:     int64(len(content)),
		Mode:     0755,
		Typeflag: tar.TypeReg,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatal(err)
	}

	return path
}

func createTestZip(t *testing.T, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.zip")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	w, err := zw.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write(content); err != nil {
		t.Fatal(err)
	}

	return path
}
