package no_global

import "go.uber.org/zap"

func tests() {
	zap.L().Info("using global logger")         // want `global logger should not be used`
	zap.S().Info("using global sugared logger") // want `global logger should not be used`
}
