package allow_args_on_same_line

import "go.uber.org/zap"

func tests(logger *zap.Logger, sugar *zap.SugaredLogger) {
	logger.Info("msg", zap.String("k1", "v1"))                   // OK
	logger.Info("msg", zap.String("k1", "v1"), zap.Int("k2", 2)) // OK
	sugar.Infow("msg", "k1", "v1", "k2", 2)                      // OK

	// This is also OK
	logger.Info("msg",
		zap.String("k1", "v1"),
		zap.Int("k2", 2),
	)

	// This is also OK
	sugar.Infow("msg",
		"k1", "v1",
		"k2", 2,
	)
}
