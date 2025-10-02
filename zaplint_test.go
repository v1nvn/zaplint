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
		"no global":                   {opts: Options{AllowGlobal: false, AllowSugar: true, AllowRawKeys: true, AllowArgsOnSameLine: true}, dir: "no_global"},
		"allow global":                {opts: Options{AllowGlobal: true, AllowSugar: true, AllowRawKeys: true, AllowArgsOnSameLine: true}, dir: "allow_global"},
		"no sugar":                    {opts: Options{AllowSugar: false, AllowGlobal: true, AllowRawKeys: true, AllowArgsOnSameLine: true}, dir: "no_sugar"},
		"allow sugar":                 {opts: Options{AllowSugar: true, AllowGlobal: true, AllowRawKeys: true, AllowArgsOnSameLine: true}, dir: "allow_sugar"},
		"static message":              {opts: Options{AllowDynamicMsg: false, AllowGlobal: true, AllowSugar: true, AllowRawKeys: true, AllowArgsOnSameLine: true}, dir: "no_dynamic_msg"},
		"allow dynamic message":       {opts: Options{AllowDynamicMsg: true, AllowGlobal: true, AllowSugar: true, AllowRawKeys: true, AllowArgsOnSameLine: true}, dir: "allow_dynamic_msg"},
		"message style (lowercased)":  {opts: Options{MsgStyle: "lowercased", AllowGlobal: true, AllowSugar: true, AllowRawKeys: true, AllowArgsOnSameLine: true}, dir: "msg_style_lowercased"},
		"message style (capitalized)": {opts: Options{MsgStyle: "capitalized", AllowGlobal: true, AllowSugar: true, AllowRawKeys: true, AllowArgsOnSameLine: true}, dir: "msg_style_capitalized"},
		"no raw keys":                 {opts: Options{AllowRawKeys: false, AllowGlobal: true, AllowSugar: true, AllowArgsOnSameLine: true}, dir: "no_raw_keys"},
		"allow raw keys":              {opts: Options{AllowRawKeys: true, AllowGlobal: true, AllowSugar: true, AllowArgsOnSameLine: true}, dir: "allow_raw_keys"},
		"key naming case":             {opts: Options{KeyNamingCase: "snake", AllowGlobal: true, AllowSugar: true, AllowRawKeys: true, AllowArgsOnSameLine: true}, dir: "key_naming_case"},
		"forbidden keys":              {opts: Options{ForbiddenKeys: []string{"time", "level", "msg"}, AllowGlobal: true, AllowSugar: true, AllowRawKeys: true, AllowArgsOnSameLine: true}, dir: "forbidden_keys"},
		"arguments on separate lines": {opts: Options{AllowArgsOnSameLine: false, AllowGlobal: true, AllowSugar: true, AllowRawKeys: true}, dir: "no_args_on_sep_lines"},
		"allow args on same line":     {opts: Options{AllowArgsOnSameLine: true, AllowGlobal: true, AllowSugar: true, AllowRawKeys: true}, dir: "allow_args_on_same_line"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			analyzer := New(&tt.opts)
			testdata := analysistest.TestData()
			analysistest.Run(t, testdata, analyzer, "z/"+tt.dir)
		})
	}
}
