package no_args_on_sep_lines

import "go.uber.org/zap"

func tests(logger *zap.Logger, sugar *zap.SugaredLogger) {
	logger.Info("msg", zap.String("k1", "v1"))                   // OK
	logger.Info("msg", zap.String("k1", "v1"), zap.Int("k2", 2)) // want `arguments should be put on separate lines`
	sugar.Infow("msg", "k1", "v1", "k2", 2)                      // want `arguments should be put on separate lines`

	// This is OK
	logger.Info("msg",
		zap.String("k1", "v1"),
		zap.Int("k2", 2),
	)

	// This is OK
	sugar.Infow("msg",
		"k1", "v1",
		"k2", 2,
	)
}
