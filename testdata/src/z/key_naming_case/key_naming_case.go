package key_naming_case

import "go.uber.org/zap"

func tests(logger *zap.Logger, sugar *zap.SugaredLogger) {
	logger.Info("msg", zap.String("user_id", "123"))  // OK
	logger.Info("msg", zap.String("userID", "123"))   // want `keys should be written in snake_case`

	sugar.Infow("msg", "request_id", "abc") // OK
	sugar.Infow("msg", "requestID", "abc")  // want `keys should be written in snake_case`
}
