package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigOptional_PrunesObsoleteSaveCooldownStatusKey(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	// Legacy key should be removed on normal load/backfill so config stays relevant.
	if err := os.WriteFile(configPath, []byte("port: 8317\nsave-cooldown-status: true\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if _, err := LoadConfigOptional(configPath, false); err != nil {
		t.Fatalf("LoadConfigOptional: %v", err)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if strings.Contains(string(data), "save-cooldown-status") {
		t.Fatalf("expected obsolete save-cooldown-status key pruned, got:\n%s", data)
	}
}

func TestParseConfigBytes_IgnoresUnknownSaveCooldownStatus(t *testing.T) {
	// Field removed from Config; unknown keys are ignored by yaml.Unmarshal.
	cfg, err := ParseConfigBytes([]byte("port: 8317\nsave-cooldown-status: false\n"))
	if err != nil {
		t.Fatalf("ParseConfigBytes: %v", err)
	}
	if cfg.Port != 8317 {
		t.Fatalf("port = %d, want 8317", cfg.Port)
	}
}
