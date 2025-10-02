// Package zaplint implements the zaplint analyzer.
package zaplint

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"iter"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/ettle/strcase"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

// Options are options for the zaplint analyzer.
type Options struct {
	NoGlobal       bool     // Enforce not using global loggers.
	NoSugar        bool     // Enforce not using the sugared logger.
	StaticMsg      bool     // Enforce using static messages.
	MsgStyle       string   // Enforce message style ("lowercased" or "capitalized").
	NoRawKeys      bool     // Enforce using constants instead of raw keys.
	KeyNamingCase  string   // Enforce key naming convention ("snake", "kebab", "camel", or "pascal").
	ForbiddenKeys  []string // Enforce not using specific keys.
	ArgsOnSepLines bool     // Enforce putting arguments on separate lines.
}

// New creates a new zaplint analyzer.
func New(opts *Options) *analysis.Analyzer {
	if opts == nil {
		opts = &Options{
			NoGlobal:       true,
			NoSugar:        true,
			StaticMsg:      true,
			MsgStyle:       styleLowercased,
			KeyNamingCase:  snakeCase,
			ArgsOnSepLines: true,
		}
	}
	return &analysis.Analyzer{
		Name:     "zaplint",
		Doc:      "ensure consistent code style when using go.uber.org/zap",
		Flags:    *flags(opts),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			if err := validateOptions(opts); err != nil {
				return nil, err
			}
			run(pass, opts)
			return nil, nil
		},
	}
}

type logFuncInfo struct {
	IsSugar   bool
	IsW       bool
	MsgPos    int
	ArgsStart int
	HasMsg    bool
}

var zapFuncs = map[string]logFuncInfo{
	"go.uber.org/zap.L":                        {},
	"go.uber.org/zap.S":                        {},
	"(*go.uber.org/zap.Logger).Debug":          {IsSugar: false, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.Logger).Info":           {IsSugar: false, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.Logger).Warn":           {IsSugar: false, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.Logger).Error":          {IsSugar: false, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.Logger).DPanic":         {IsSugar: false, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.Logger).Panic":          {IsSugar: false, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.Logger).Fatal":          {IsSugar: false, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.Logger).With":           {IsSugar: false, ArgsStart: 0, HasMsg: false},
	"(*go.uber.org/zap.Logger).Sugar":          {IsSugar: false, ArgsStart: 0, HasMsg: false},
	"(*go.uber.org/zap.SugaredLogger).Debug":   {IsSugar: true, ArgsStart: 0, HasMsg: false},
	"(*go.uber.org/zap.SugaredLogger).Info":    {IsSugar: true, ArgsStart: 0, HasMsg: false},
	"(*go.uber.org/zap.SugaredLogger).Warn":    {IsSugar: true, ArgsStart: 0, HasMsg: false},
	"(*go.uber.org/zap.SugaredLogger).Error":   {IsSugar: true, ArgsStart: 0, HasMsg: false},
	"(*go.uber.org/zap.SugaredLogger).DPanic":  {IsSugar: true, ArgsStart: 0, HasMsg: false},
	"(*go.uber.org/zap.SugaredLogger).Panic":   {IsSugar: true, ArgsStart: 0, HasMsg: false},
	"(*go.uber.org/zap.SugaredLogger).Fatal":   {IsSugar: true, ArgsStart: 0, HasMsg: false},
	"(*go.uber.org/zap.SugaredLogger).Debugf":  {IsSugar: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).Infof":   {IsSugar: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).Warnf":   {IsSugar: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).Errorf":  {IsSugar: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).DPanicf": {IsSugar: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).Panicf":  {IsSugar: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).Fatalf":  {IsSugar: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).Debugw":  {IsSugar: true, IsW: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).Infow":   {IsSugar: true, IsW: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).Warnw":   {IsSugar: true, IsW: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).Errorw":  {IsSugar: true, IsW: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).DPanicw": {IsSugar: true, IsW: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).Panicw":  {IsSugar: true, IsW: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).Fatalw":  {IsSugar: true, IsW: true, MsgPos: 0, ArgsStart: 1, HasMsg: true},
	"(*go.uber.org/zap.SugaredLogger).With":    {IsSugar: true, ArgsStart: 0, HasMsg: false},
}

func run(pass *analysis.Pass, opts *Options) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}
	inspector.Preorder(nodeFilter, func(node ast.Node) {
		visit(pass, node.(*ast.CallExpr), opts)
	})
}

// cleanVendorPath removes vendor prefixes from package paths.
func cleanVendorPath(path string) string {
	const vendor = "vendor/"
	i := strings.LastIndex(path, vendor)
	if i == -1 {
		return path
	}
	slashBefore := i - 1
	start := 0
	if slashBefore >= 0 {
		j := slashBefore - 1
		for j >= 0 {
			c := path[j]
			if strings.ContainsRune("/*[] ,\t)", rune(c)) {
				break
			}
			j--
		}
		start = j + 1
	}
	return path[:start] + path[i+len(vendor):]
}

func visit(pass *analysis.Pass, call *ast.CallExpr, opts *Options) {
	fn := typeutil.StaticCallee(pass.TypesInfo, call)
	if fn == nil {
		return
	}
	originalFullName := fn.FullName()
	cleanedFullName := cleanVendorPath(originalFullName)

	info, ok := zapFuncs[cleanedFullName]
	if !ok {
		return
	}

	// Get the position for reporting - use selector position if available for better error location
	reportPos := call.Pos()
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		reportPos = sel.Sel.Pos()
	}

	if opts.NoGlobal {
		if cleanedFullName == "go.uber.org/zap.L" || cleanedFullName == "go.uber.org/zap.S" {
			pass.Reportf(reportPos, "global logger should not be used")
			return
		}
	}
	if opts.NoSugar && info.IsSugar {
		// For chained calls like sugar.With().Info(), only report on the inner call (.With)
		// to avoid duplicate diagnostics. Skip if the receiver is a sugared logger call.
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if innerCall, ok := sel.X.(*ast.CallExpr); ok {
				if innerFn := typeutil.StaticCallee(pass.TypesInfo, innerCall); innerFn != nil {
					innerFullName := cleanVendorPath(innerFn.FullName())
					if innerInfo, ok := zapFuncs[innerFullName]; ok && innerInfo.IsSugar {
						// This call's receiver is a sugared logger call, skip reporting
						// (the inner call will be reported instead)
						return
					}
				}
			}
		}
		pass.Reportf(reportPos, "sugared logger should not be used")
		return
	}

	logArgs := call.Args[info.ArgsStart:]

	if opts.StaticMsg && info.HasMsg && len(call.Args) > info.MsgPos {
		msgArg := call.Args[info.MsgPos]
		if !isStaticMsg(pass.TypesInfo, msgArg) {
			pass.Reportf(msgArg.Pos(), "message should be a string literal or a constant")
		}
	}

	if opts.MsgStyle != "" && info.HasMsg && len(call.Args) > info.MsgPos {
		checkMsgStyle(pass, call.Args[info.MsgPos], opts.MsgStyle)
	}

	keys := allKeys(pass, fn.Name(), info, logArgs)
	checkAllKeys(pass, opts, keys)

	if opts.ArgsOnSepLines && areArgsOnSameLine(pass.Fset, info.IsW, logArgs) {
		pass.Reportf(call.Pos(), "arguments should be put on separate lines")
	}
}

func allKeys(pass *analysis.Pass, funcName string, fnInfo logFuncInfo, args []ast.Expr) iter.Seq[ast.Expr] {
	return func(yield func(key ast.Expr) bool) {
		if !fnInfo.IsSugar {
			for _, arg := range args {
				if call, ok := arg.(*ast.CallExpr); ok {
					if callee := typeutil.StaticCallee(pass.TypesInfo, call); callee != nil && callee.Pkg() != nil {
						pkgPath := cleanVendorPath(callee.Pkg().Path())
						if pkgPath == "go.uber.org/zap" && len(call.Args) > 0 {
							if !yield(call.Args[0]) {
								return
							}
						}
					}
				}
			}
		} else if fnInfo.IsW || funcName == "With" {
			for i := 0; i < len(args); i += 2 {
				if i < len(args) {
					if !yield(args[i]) {
						return
					}
				}
			}
		}
	}
}

func areArgsOnSameLine(fset *token.FileSet, isW bool, logArgs []ast.Expr) bool {
	if len(logArgs) <= 1 {
		return false
	}

	// For W-style functions (Infow, Errorw, etc.), allow key-value pairs on the same line
	// Check that pairs are on separate lines from other pairs
	if isW {
		if len(logArgs)%2 != 0 {
			return false // odd number of args, can't check pairs properly
		}
		lines := make(map[int]bool)
		for i := 0; i < len(logArgs); i += 2 {
			// Check the line of the key (every even index)
			line := fset.Position(logArgs[i].Pos()).Line
			if lines[line] {
				return true // two pairs on the same line
			}
			lines[line] = true
		}
		return false
	}

	// For non-W functions, check that no two arguments are on the same line
	lines := make(map[int]bool)
	for _, arg := range logArgs {
		line := fset.Position(arg.Pos()).Line
		if lines[line] {
			return true
		}
		lines[line] = true
	}
	return false
}

var errInvalidValue = errors.New("invalid value")

const (
	snakeCase        = "snake"
	kebabCase        = "kebab"
	camelCase        = "camel"
	pascalCase       = "pascal"
	styleLowercased  = "lowercased"
	styleCapitalized = "capitalized"
)

func validateOptions(opts *Options) error {
	switch opts.MsgStyle {
	case "", styleLowercased, styleCapitalized:
	default:
		return fmt.Errorf("zaplint: Options.MsgStyle=%s: %w", opts.MsgStyle, errInvalidValue)
	}
	switch opts.KeyNamingCase {
	case "", snakeCase, kebabCase, camelCase, pascalCase:
	default:
		return fmt.Errorf("zaplint: Options.KeyNamingCase=%s: %w", opts.KeyNamingCase, errInvalidValue)
	}
	return nil
}

func flags(opts *Options) *flag.FlagSet {
	fset := flag.NewFlagSet("zaplint", flag.ContinueOnError)
	fset.BoolVar(&opts.NoGlobal, "no-global", opts.NoGlobal, "enforce not using global loggers (zap.L() and zap.S())")
	fset.BoolVar(&opts.NoSugar, "no-sugar", opts.NoSugar, "enforce using zap.Logger over zap.SugaredLogger")
	fset.BoolVar(&opts.StaticMsg, "static-msg", opts.StaticMsg, "enforce using static messages")
	fset.StringVar(&opts.MsgStyle, "msg-style", opts.MsgStyle, "enforce message style (lowercased|capitalized)")
	fset.BoolVar(&opts.NoRawKeys, "no-raw-keys", opts.NoRawKeys, "enforce using constants instead of raw keys")
	fset.StringVar(&opts.KeyNamingCase, "key-naming-case", opts.KeyNamingCase, "enforce key naming convention (snake|kebab|camel|pascal)")
	fset.BoolVar(&opts.ArgsOnSepLines, "args-on-sep-lines", opts.ArgsOnSepLines, "enforce putting arguments on separate lines")
	fset.Func("forbidden-keys", "comma-separated list of forbidden keys", func(s string) error {
		if s != "" {
			opts.ForbiddenKeys = append(opts.ForbiddenKeys, strings.Split(s, ",")...)
		}
		return nil
	})
	return fset
}

func isStaticMsg(info *types.Info, msg ast.Expr) bool {
	switch msg := msg.(type) {
	case *ast.BasicLit:
		return msg.Kind == token.STRING
	case *ast.Ident:
		if obj := info.ObjectOf(msg); obj != nil {
			_, isConst := obj.(*types.Const)
			return isConst
		}
		return false
	case *ast.BinaryExpr:
		if msg.Op != token.ADD {
			return false
		}
		return isStaticMsg(info, msg.X) && isStaticMsg(info, msg.Y)
	default:
		return false
	}
}

func checkMsgStyle(pass *analysis.Pass, msg ast.Expr, style string) {
	lit, ok := msg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return
	}
	value, err := strconv.Unquote(lit.Value)
	if err != nil || value == "" {
		return
	}
	firstRune := []rune(value)[0]
	isValid := false
	switch style {
	case styleLowercased:
		if unicode.IsLower(firstRune) || !unicode.IsLetter(firstRune) {
			isValid = true
		}
	case styleCapitalized:
		if unicode.IsUpper(firstRune) || !unicode.IsLetter(firstRune) {
			isValid = true
		}
	}
	if !isValid {
		pass.Reportf(msg.Pos(), "message should be %s", style)
	}
}

func checkAllKeys(pass *analysis.Pass, opts *Options, keys iter.Seq[ast.Expr]) {
	caseFn, caseName := getCaseConverter(opts.KeyNamingCase)
	for keyExpr := range keys {
		// keyRender := render(pass.Fset, keyExpr)
		if opts.NoRawKeys {
			if _, ok := keyExpr.(*ast.BasicLit); ok {
				pass.Reportf(keyExpr.Pos(), "raw keys should not be used")
			}
		}
		keyName, ok := getKeyName(keyExpr)
		if !ok {
			continue
		}
		if len(opts.ForbiddenKeys) > 0 && slices.Contains(opts.ForbiddenKeys, keyName) {
			pass.Reportf(keyExpr.Pos(), "%q key is forbidden and should not be used", keyName)
		}
		if caseFn != nil && keyName != caseFn(keyName) {
			pass.Report(analysis.Diagnostic{
				Pos:     keyExpr.Pos(),
				Message: fmt.Sprintf("keys should be written in %s", caseName),
				SuggestedFixes: []analysis.SuggestedFix{{
					Message: fmt.Sprintf("Change to %q", caseFn(keyName)),
					TextEdits: []analysis.TextEdit{{
						Pos: keyExpr.Pos(), End: keyExpr.End(),
						NewText: []byte(strconv.Quote(caseFn(keyName))),
					}},
				}},
			})
		}
	}
}

func getKeyName(key ast.Expr) (string, bool) {
	if lit, ok := key.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		value, err := strconv.Unquote(lit.Value)
		if err == nil {
			return value, true
		}
	}
	return "", false
}

func getCaseConverter(style string) (func(string) string, string) {
	switch style {
	case snakeCase:
		return strcase.ToSnake, "snake_case"
	case kebabCase:
		return strcase.ToKebab, "kebab-case"
	case camelCase:
		return strcase.ToCamel, "camelCase"
	case pascalCase:
		return strcase.ToPascal, "PascalCase"
	}
	return nil, ""
}
