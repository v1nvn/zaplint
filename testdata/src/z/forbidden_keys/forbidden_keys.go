package forbidden_keys

import "go.uber.org/zap"

func tests(logger *zap.Logger, sugar *zap.SugaredLogger) {
	logger.Info("msg", zap.String("user_id", "123")) // OK
	logger.Info("msg", zap.String("msg", "hello"))   // want `\"msg\" key is forbidden and should not be used`

	sugar.Infow("another message", "level", "info") // want `\"level\" key is forbidden and should not be used`
}
