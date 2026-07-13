package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigOptional_SaveCooldownStatusMissingDefaultsTrueAndBackfills(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("port: 8317\nhost: \"\"\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadConfigOptional(configPath, false)
	if err != nil {
		t.Fatalf("LoadConfigOptional() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("expected config")
	}
	if !cfg.SaveCooldownStatus {
		t.Fatal("expected SaveCooldownStatus=true when key is missing")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(data), "save-cooldown-status:") {
		t.Fatalf("expected config file to backfill save-cooldown-status, got:\n%s", data)
	}
}

func TestLoadConfigOptional_SaveCooldownStatusExplicitFalsePreserved(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("port: 8317\nsave-cooldown-status: false\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadConfigOptional(configPath, false)
	if err != nil {
		t.Fatalf("LoadConfigOptional() error = %v", err)
	}
	if cfg.SaveCooldownStatus {
		t.Fatal("expected SaveCooldownStatus=false when explicitly set")
	}
}

func TestParseConfigBytes_SaveCooldownStatusMissingDefaultsTrue(t *testing.T) {
	cfg, err := ParseConfigBytes([]byte("port: 8317\n"))
	if err != nil {
		t.Fatalf("ParseConfigBytes() error = %v", err)
	}
	if !cfg.SaveCooldownStatus {
		t.Fatal("expected SaveCooldownStatus=true when key is missing")
	}
}

func TestParseConfigBytes_SaveCooldownStatusExplicitFalse(t *testing.T) {
	cfg, err := ParseConfigBytes([]byte("port: 8317\nsave-cooldown-status: false\n"))
	if err != nil {
		t.Fatalf("ParseConfigBytes() error = %v", err)
	}
	if cfg.SaveCooldownStatus {
		t.Fatal("expected SaveCooldownStatus=false when explicitly set")
	}
}
