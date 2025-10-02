package linters

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/v1nvn/zaplint"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("zaplint", New)
}

type ZapLintPlugin struct {
	settings *zaplint.Options
}

func New(settings any) (register.LinterPlugin, error) {
	var opts *zaplint.Options

	// If settings are provided, decode them using DecodeSettings
	// The zero values of Options struct serve as defaults (all checks disabled by default)
	if settings != nil {
		s, err := register.DecodeSettings[zaplint.Options](settings)
		if err != nil {
			return nil, err
		}
		opts = &s
	}

	// Pass opts to zaplint.New() which will use zero values as defaults if opts is nil
	return &ZapLintPlugin{settings: opts}, nil
}

func (z *ZapLintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		zaplint.New(z.settings),
	}, nil
}

func (z *ZapLintPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
