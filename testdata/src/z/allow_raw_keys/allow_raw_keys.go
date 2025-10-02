package allow_raw_keys

import "go.uber.org/zap"

const UserIDKey = "user_id"

func tests(logger *zap.Logger, sugar *zap.SugaredLogger) {
	logger.Info("msg", zap.String(UserIDKey, "123")) // OK
	logger.Info("msg", zap.String("user_id", "123")) // OK

	sugar.Infow("msg", UserIDKey, "123") // OK
	sugar.Infow("msg", "user_id", "123") // OK
}
