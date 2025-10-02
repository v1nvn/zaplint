# zaplint

A Go linter that ensures consistent code style when using `go.uber.org/zap`.

## 📌 About

The `go.uber.org/zap` library provides both structured logging through `zap.Logger` and a more flexible sugared API through `zap.SugaredLogger`. While teams may have different preferences about which API to use and how to structure their logs, most agree on one thing: it should be consistent. With `zaplint` you can enforce various rules for `go.uber.org/zap` based on your preferred code style.

## 🚀 Features

* Disallow using global loggers (enabled by default)
* Disallow using the sugared logger (enabled by default)
* Disallow dynamic log messages (enabled by default)
* Enforce message style - lowercased (enabled by default)
* Disallow using raw string keys (enabled by default)
* Enforce key naming convention - snake (enabled by default)
* Disallow specific keys (optional)
* Disallow putting arguments on the same line (enabled by default)

## 📦 Install

### golangci-lint integration

`zaplint` can be used as a custom linter in [`golangci-lint`][1] v2.4.0+. See [GOLANGCI_PLUGIN.md](GOLANGCI_PLUGIN.md) for detailed setup instructions.

Quick setup:

```bash
go get github.com/v1nvn/zaplint@latest
```

Add to your `.golangci.yml`:

```yaml
linters-settings:
  custom:
    zaplint:
      type: "module"
      description: Ensure consistent code style when using go.uber.org/zap
      # All settings are optional - defaults shown below
      # settings:
      #   allow-global: false       # Disallow global loggers (default)
      #   allow-sugar: false        # Disallow sugared logger (default)
      #   allow-dynamic-msg: false  # Disallow dynamic messages (default)
      #   msg-style: "lowercased"   # Enforce lowercased messages (default)
      #   allow-raw-keys: false     # Disallow raw keys (default)
      #   key-naming-case: "snake"  # Enforce snake_case keys (default)
      #   allow-args-on-same-line: false  # Disallow args on same line (default)

linters:
  enable:
    - zaplint
```

### Standalone usage

You can also use `zaplint` standalone:

```bash
go install github.com/v1nvn/zaplint/cmd/zaplint@latest
```

## 📋 Usage

Run `zaplint` on your Go packages with opinionated defaults enabled:

```bash
zaplint ./...
```

The default configuration is opinionated and enforces strict logging practices. All checks are enabled by default. You can customize behavior by passing flags:

```bash
# Allow global loggers
zaplint -allow-global ./...

# Allow sugared logger
zaplint -allow-sugar ./...

# Allow dynamic messages
zaplint -allow-dynamic-msg ./...

# Allow raw keys
zaplint -allow-raw-keys ./...

# Disable message style check
zaplint -msg-style= ./...

# Add forbidden keys
zaplint -forbidden-keys=password,secret ./...
```

### No global

Some projects prefer to pass loggers as explicit dependencies.
The `no-global` option causes `zaplint` to report the use of global loggers:

```go
zap.L().Info("user logged in") // zaplint: global logger should not be used
zap.S().Info("user logged in") // zaplint: global logger should not be used
```

### No sugar

Some teams prefer to use the structured `zap.Logger` exclusively for better performance and type safety.
The `no-sugar` option causes `zaplint` to report any use of `zap.SugaredLogger`:

```go
logger := zap.NewExample().Sugar()
logger.Info("user logged in") // zaplint: sugared logger should not be used
```

### Static messages

To get the most out of structured logging, you may want to require log messages to be static.
The `static-msg` option causes `zaplint` to report non-static messages:

```go
logger.Info(fmt.Sprintf("user with id %d logged in", 42)) // zaplint: message should be a string literal or a constant
```

The report can be fixed by moving dynamic values to structured fields:

```go
logger.Info("user logged in", zap.Int("user_id", 42))
```

### Message style

The `msg-style` option causes `zaplint` to check log messages for a particular style.

Possible values are `lowercased` (report messages that begin with an uppercase letter)...

```go
logger.Info("User logged in") // zaplint: message should be lowercased
```

...and `capitalized` (report messages that begin with a lowercase letter):

```go
logger.Info("user logged in") // zaplint: message should be capitalized
```

Special cases such as acronyms (e.g. `HTTP`, `U.S.`) are ignored.

### No raw keys

To prevent typos, you may want to forbid the use of raw keys altogether.
The `no-raw-keys` option causes `zaplint` to report the use of strings as keys
(including `zap.Field` calls, e.g. `zap.Int("user_id", 42)`):

```go
logger.Info("user logged in", zap.Int("user_id", 42)) // zaplint: raw keys should not be used
```

This report can be fixed by using either constants...

```go
const UserID = "user_id"

logger.Info("user logged in", zap.Int(UserID, 42))
```

...or custom `zap.Field` constructors:

```go
func UserID(value int) zap.Field { return zap.Int("user_id", value) }

logger.Info("user logged in", UserID(42))
```

### Key naming convention

To ensure consistency in logs, you may want to enforce a single key naming convention.
The `key-naming-case` option causes `zaplint` to report keys written in a case other than the given one:

```go
logger.Info("user logged in", zap.Int("user-id", 42)) // zaplint: keys should be written in snake_case
```

Possible values are `snake`, `kebab`, `camel`, or `pascal`.

### Forbidden keys

To prevent accidental use of reserved log keys, you may want to forbid specific keys altogether.
The `forbidden-keys` option causes `zaplint` to report the use of forbidden keys:

```go
logger.Info("user logged in", zap.String("reserved", "value")) // zaplint: "reserved" key is forbidden and should not be used
```

For example, when using custom log processors or exporters, you may want to forbid keys that conflict with your logging infrastructure's reserved fields.

### Arguments on separate lines

To improve code readability, you may want to put arguments on separate lines, especially when using the structured logger.
The `args-on-sep-lines` option causes `zaplint` to report 2+ arguments on the same line:

```go
logger.Info("user logged in", zap.Int("user_id", 42), zap.String("ip", "192.0.2.0")) // zaplint: arguments should be put on separate lines
```

This report can be fixed by reformatting the code:

```go
logger.Info("user logged in",
    zap.Int("user_id", 42),
    zap.String("ip", "192.0.2.0"),
)
```

For `SugaredLogger` methods with the `w` suffix (e.g., `Infow`, `Errorw`), key-value pairs are allowed on the same line, but different pairs should be on separate lines:

```go
// OK: each key-value pair on its own line
sugar.Infow("user logged in",
    "user_id", 42,
    "ip", "192.0.2.0",
)

// Not OK: multiple pairs on the same line
sugar.Infow("user logged in", "user_id", 42, "ip", "192.0.2.0") // zaplint: arguments should be put on separate lines
```

[1]: https://golangci-lint.run
[2]: https://github.com/v1nvn/zaplint/releases
