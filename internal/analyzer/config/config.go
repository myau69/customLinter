package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	Rules     RulesConfig     `json:"rules"`
	AutoFix   AutoFixConfig   `json:"autofix"`
	Sensitive SensitiveConfig `json:"sensitive"`
}

type RulesConfig struct {
	Lowercase bool `json:"lowercase"`
	English   bool `json:"english"`
	Specials  bool `json:"specials"`
	Sensitive bool `json:"sensitive"`
}

type AutoFixConfig struct {
	Enabled   bool `json:"enabled"`
	Lowercase bool `json:"lowercase"`
	Specials  bool `json:"specials"`
}

type SensitiveConfig struct {
	Patterns []string `json:"patterns"`
}

func Default() Config {
	return Config{
		Rules: RulesConfig{
			Lowercase: true,
			English:   true,
			Specials:  true,
			Sensitive: true,
		},
		AutoFix: AutoFixConfig{
			Enabled:   false,
			Lowercase: true,
			Specials:  true,
		},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if strings.TrimSpace(path) == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	cfg.Sensitive.Patterns = cleanPatterns(cfg.Sensitive.Patterns)
	return cfg, nil
}

func cleanPatterns(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))

	for _, raw := range in {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}

		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		out = append(out, trimmed)
	}

	return out
}
