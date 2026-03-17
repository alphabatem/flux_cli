package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
)

const githubRepo = "alphabatem/flux_cli"

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Flux CLI to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		updated, newVersion, err := runUpdate(true)
		if err != nil {
			return err
		}
		if updated {
			output.PrintSuccess(cmd, map[string]string{
				"status":  "updated",
				"version": newVersion,
			}, nil)
		} else {
			output.PrintSuccess(cmd, map[string]string{
				"status":  "up_to_date",
				"version": Version,
			}, nil)
		}
		return nil
	},
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// checkForUpdate runs a silent update check on startup.
func checkForUpdate() {
	updated, newVersion, _ := runUpdate(false)
	if updated {
		fmt.Fprintf(os.Stderr, "flux: auto-updated to %s\n", newVersion)
	}
}

// runUpdate checks for a new version and updates the binary if needed.
// Returns (updated bool, newVersion string, err error).
func runUpdate(verbose bool) (bool, string, error) {
	release, err := fetchLatestRelease()
	if err != nil {
		if verbose {
			return false, "", fmt.Errorf("failed to check for updates: %w", err)
		}
		return false, "", nil
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(Version, "v")

	if !isNewer(current, latest) {
		return false, "", nil
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Updating v%s → v%s...\n", current, latest)
	}

	// Find the right asset
	assetName := buildAssetName(latest)
	var downloadURL string
	for _, a := range release.Assets {
		if a.Name == assetName {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		if verbose {
			return false, "", fmt.Errorf("no binary found for %s/%s (expected %s)", runtime.GOOS, runtime.GOARCH, assetName)
		}
		return false, "", nil
	}

	if err := downloadAndReplace(downloadURL); err != nil {
		if verbose {
			return false, "", fmt.Errorf("update failed: %w", err)
		}
		return false, "", nil
	}

	return true, latest, nil
}

func fetchLatestRelease() (*githubRelease, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", githubRepo))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}

func buildAssetName(version string) string {
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	return fmt.Sprintf("flux_%s_%s_%s.%s", version, runtime.GOOS, runtime.GOARCH, ext)
}

// isNewer returns true if latest is a higher semver than current.
func isNewer(current, latest string) bool {
	// Non-semver current (e.g. "dev", "none") is always older than any release
	if !isSemver(current) {
		return isSemver(latest)
	}
	if !isSemver(latest) {
		return false
	}

	cp := strings.Split(current, ".")
	lp := strings.Split(latest, ".")

	maxLen := len(cp)
	if len(lp) > maxLen {
		maxLen = len(lp)
	}

	for i := 0; i < maxLen; i++ {
		var c, l int
		if i < len(cp) {
			c, _ = strconv.Atoi(cp[i])
		}
		if i < len(lp) {
			l, _ = strconv.Atoi(lp[i])
		}
		if l > c {
			return true
		}
		if l < c {
			return false
		}
	}
	return false
}

func isSemver(v string) bool {
	parts := strings.Split(v, ".")
	if len(parts) < 2 {
		return false
	}
	for _, p := range parts {
		if _, err := strconv.Atoi(p); err != nil {
			return false
		}
	}
	return true
}

func downloadAndReplace(url string) error {
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download returned HTTP %d", resp.StatusCode)
	}

	// Download archive to temp file
	tmpArchive, err := os.CreateTemp("", "flux-update-archive-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmpArchive.Name())

	if _, err := io.Copy(tmpArchive, resp.Body); err != nil {
		tmpArchive.Close()
		return fmt.Errorf("downloading archive: %w", err)
	}
	tmpArchive.Close()

	// Extract the binary
	binaryData, err := extractBinary(tmpArchive.Name())
	if err != nil {
		return fmt.Errorf("extracting binary: %w", err)
	}

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding current executable: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("resolving executable path: %w", err)
	}

	// Write new binary to temp file next to the executable
	tmpBin, err := os.CreateTemp(filepath.Dir(execPath), "flux-new-*")
	if err != nil {
		return fmt.Errorf("creating temp binary: %w", err)
	}
	tmpBinPath := tmpBin.Name()
	defer os.Remove(tmpBinPath)

	if _, err := tmpBin.Write(binaryData); err != nil {
		tmpBin.Close()
		return fmt.Errorf("writing new binary: %w", err)
	}
	tmpBin.Close()

	if err := os.Chmod(tmpBinPath, 0755); err != nil {
		return fmt.Errorf("setting permissions: %w", err)
	}

	// Atomic replace
	if err := os.Rename(tmpBinPath, execPath); err != nil {
		// On Windows or cross-device, fall back to copy
		return copyFile(tmpBinPath, execPath)
	}

	return nil
}

// extractBinary pulls the "flux" (or "flux.exe") binary out of a tar.gz or zip archive.
func extractBinary(archivePath string) ([]byte, error) {
	binaryName := "flux"
	if runtime.GOOS == "windows" {
		binaryName = "flux.exe"
	}

	if strings.HasSuffix(archivePath, ".zip") || runtime.GOOS == "windows" {
		return extractFromZip(archivePath, binaryName)
	}
	return extractFromTarGz(archivePath, binaryName)
}

func extractFromTarGz(archivePath, binaryName string) ([]byte, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if filepath.Base(hdr.Name) == binaryName && hdr.Typeflag == tar.TypeReg {
			return io.ReadAll(tr)
		}
	}

	return nil, fmt.Errorf("binary %q not found in archive", binaryName)
}

func extractFromZip(archivePath, binaryName string) ([]byte, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if filepath.Base(f.Name) == binaryName {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}

	return nil, fmt.Errorf("binary %q not found in archive", binaryName)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0755)
}
