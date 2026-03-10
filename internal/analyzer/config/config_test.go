package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefault(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.Rules.Lowercase || !cfg.Rules.English || !cfg.Rules.Specials || !cfg.Rules.Sensitive {
		t.Fatalf("default rules should be enabled: %#v", cfg.Rules)
	}
}

func TestLoadFromJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "customlinter.json")

	content := `{
		"rules": {"lowercase": true, "english": false, "specials": true, "sensitive": true},
		"autofix": {"enabled": true, "lowercase": true, "specials": false},
		"sensitive": {"patterns": ["session_id", "Session_ID", "", " auth-token "]}
	}`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Rules.English {
		t.Fatalf("english rule should be disabled: %#v", cfg.Rules)
	}
	if !cfg.AutoFix.Enabled || !cfg.AutoFix.Lowercase || cfg.AutoFix.Specials {
		t.Fatalf("unexpected autofix config: %#v", cfg.AutoFix)
	}
	if got, want := len(cfg.Sensitive.Patterns), 2; got != want {
		t.Fatalf("patterns len = %d, want %d, patterns=%v", got, want, cfg.Sensitive.Patterns)
	}
}
