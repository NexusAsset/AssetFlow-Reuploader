package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type ghRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

func parseVer(v string) []int {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	if i := strings.IndexAny(v, "-+ "); i >= 0 {
		v = v[:i]
	}
	parts := strings.Split(v, ".")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		n, _ := strconv.Atoi(strings.TrimSpace(p))
		out = append(out, n)
	}
	return out
}

// verLess reports whether version a is strictly older than b.
func verLess(a, b string) bool {
	pa, pb := parseVer(a), parseVer(b)
	n := len(pa)
	if len(pb) > n {
		n = len(pb)
	}
	for i := 0; i < n; i++ {
		var x, y int
		if i < len(pa) {
			x = pa[i]
		}
		if i < len(pb) {
			y = pb[i]
		}
		if x != y {
			return x < y
		}
	}
	return false
}

func (s *Server) fetchLatest() (*ghRelease, error) {
	if s.updateRepo == "" {
		return nil, fmt.Errorf("no update repo configured")
	}
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+s.updateRepo+"/releases/latest", nil)
	req.Header.Set("User-Agent", "AssetFlowReuploader")
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := (&http.Client{Timeout: 12 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		msg := string(b)
		if len(msg) > 120 {
			msg = msg[:120]
		}
		return nil, fmt.Errorf("github %d: %s", resp.StatusCode, msg)
	}
	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}
	return &rel, nil
}

// platformAsset picks the right release asset for the running OS. Windows wants
// the .exe; mac wants the build whose name carries the arch.
func (rel *ghRelease) platformAsset() (string, string) {
	pick := func(match func(string) bool) (string, string) {
		for _, a := range rel.Assets {
			if match(strings.ToLower(a.Name)) {
				return a.Name, a.URL
			}
		}
		return "", ""
	}
	switch runtime.GOOS {
	case "windows":
		return pick(func(n string) bool { return strings.HasSuffix(n, ".exe") })
	case "darwin":
		arch := "amd64"
		if runtime.GOARCH == "arm64" {
			arch = "arm64"
		}
		return pick(func(n string) bool { return strings.Contains(n, "mac") && strings.Contains(n, arch) })
	default:
		return "", ""
	}
}

func (s *Server) handleUpdateCheck(w http.ResponseWriter, r *http.Request) {
	rel, err := s.fetchLatest()
	if err != nil {
		writeJSON(w, map[string]any{"current": s.appVersion, "hasUpdate": false, "error": err.Error()})
		return
	}
	latest := strings.TrimPrefix(rel.TagName, "v")
	_, url := rel.platformAsset()
	writeJSON(w, map[string]any{
		"current":   s.appVersion,
		"latest":    latest,
		"hasUpdate": verLess(s.appVersion, latest) && url != "",
		"notes":     rel.Body,
		"url":       url,
	})
}

func (s *Server) handleUpdateApply(w http.ResponseWriter, r *http.Request) {
	if runtime.GOOS != "windows" {
		writeJSON(w, map[string]any{"ok": false, "error": "in-place update is Windows-only for now; grab the latest build from GitHub"})
		return
	}
	rel, err := s.fetchLatest()
	if err != nil {
		writeJSON(w, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	if !verLess(s.appVersion, strings.TrimPrefix(rel.TagName, "v")) {
		writeJSON(w, map[string]any{"ok": false, "error": "already up to date"})
		return
	}
	_, url := rel.platformAsset()
	if url == "" {
		writeJSON(w, map[string]any{"ok": false, "error": "no Windows build in the latest release"})
		return
	}
	exe, err := os.Executable()
	if err != nil {
		writeJSON(w, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	newPath := exe + ".new"
	if err := downloadTo(url, newPath); err != nil {
		writeJSON(w, map[string]any{"ok": false, "error": "download failed: " + err.Error()})
		return
	}
	if err := launchUpdater(exe, newPath); err != nil {
		_ = os.Remove(newPath)
		writeJSON(w, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	s.push("info", "Update downloaded. Restarting into the new version...")
	writeJSON(w, map[string]any{"ok": true})
	go func() { time.Sleep(700 * time.Millisecond); os.Exit(0) }()
}

func downloadTo(url, path string) error {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "AssetFlowReuploader")
	resp, err := (&http.Client{Timeout: 5 * time.Minute}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// launchUpdater writes a batch that waits for this process to exit, swaps the new
// exe over the old one, relaunches it, and deletes itself, then starts it detached.
func launchUpdater(exe, newPath string) error {
	dir := filepath.Dir(exe)
	bat := filepath.Join(dir, "assetflow-update.bat")
	pid := strconv.Itoa(os.Getpid())
	script := "@echo off\r\n" +
		":wait\r\n" +
		"tasklist /FI \"PID eq " + pid + "\" 2>nul | find \"" + pid + "\" >nul\r\n" +
		"if not errorlevel 1 ( timeout /t 1 /nobreak >nul & goto wait )\r\n" +
		"move /Y \"" + newPath + "\" \"" + exe + "\" >nul\r\n" +
		"start \"\" \"" + exe + "\"\r\n" +
		"del \"%~f0\"\r\n"
	if err := os.WriteFile(bat, []byte(script), 0o644); err != nil {
		return err
	}
	return exec.Command("cmd", "/c", "start", "", "/min", bat).Start()
}
