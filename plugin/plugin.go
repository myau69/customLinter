package plugin

import (
	"fmt"

	"github.com/myau69/customLinter/internal/analyzer"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

const LinterName = "customlinter"

func init() {
	register.Plugin(LinterName, New)
}

type Plugin struct {
	configPath string
}

func New(settings any) (register.LinterPlugin, error) {
	p := &Plugin{}

	if m, ok := settings.(map[string]any); ok {
		if v, ok := m["config"].(string); ok {
			p.configPath = v
		}
	}

	return p, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	a := analyzer.Analyzer
	if p.configPath != "" {
		if err := a.Flags.Set("config", p.configPath); err != nil {
			return nil, fmt.Errorf("set config flag: %w", err)
		}
	}
	return []*analysis.Analyzer{a}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
