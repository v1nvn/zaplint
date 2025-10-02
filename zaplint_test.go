package zaplint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tests := map[string]struct {
		opts Options
		dir  string
	}{
		"no global":                   {opts: Options{NoGlobal: true}, dir: "no_global"},
		"no sugar":                    {opts: Options{NoSugar: true}, dir: "no_sugar"},
		"static message":              {opts: Options{StaticMsg: true}, dir: "static_msg"},
		"message style (lowercased)":  {opts: Options{MsgStyle: "lowercased"}, dir: "msg_style_lowercased"},
		"message style (capitalized)": {opts: Options{MsgStyle: "capitalized"}, dir: "msg_style_capitalized"},
		"no raw keys":                 {opts: Options{NoRawKeys: true}, dir: "no_raw_keys"},
		"key naming case":             {opts: Options{KeyNamingCase: "snake"}, dir: "key_naming_case"},
		"forbidden keys":              {opts: Options{ForbiddenKeys: []string{"time", "level", "msg"}}, dir: "forbidden_keys"},
		"arguments on separate lines": {opts: Options{ArgsOnSepLines: true}, dir: "args_on_sep_lines"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			analyzer := New(&tt.opts)
			testdata := analysistest.TestData()
			analysistest.Run(t, testdata, analyzer, "z/"+tt.dir)
		})
	}
}
