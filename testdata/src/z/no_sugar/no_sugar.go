package no_sugar

import "go.uber.org/zap"

func tests(logger *zap.Logger) {
	sugar := logger.Sugar()
	sugar.Info("hello")               // want `sugared logger should not be used`
	sugar.Infow("hello", "key", 1)    // want `sugared logger should not be used`
	sugar.With("key", 1).Info("test") // want `sugared logger should not be used`
}
