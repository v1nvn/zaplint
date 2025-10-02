package allow_global

import "go.uber.org/zap"

func tests() {
	zap.L().Info("using global logger")         // OK
	zap.S().Info("using global sugared logger") // OK
}
