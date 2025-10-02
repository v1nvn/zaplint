package allow_sugar

import "go.uber.org/zap"

func tests(logger *zap.Logger) {
	sugar := logger.Sugar()
	sugar.Info("hello")               // OK
	sugar.Infow("hello", "key", 1)    // OK
	sugar.With("key", 1).Info("test") // OK
}
