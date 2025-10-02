package msg_style_capitalized

import "go.uber.org/zap"

func tests(logger *zap.Logger) {
	logger.Info("Capitalized message") // OK
	logger.Info("message")             // want `message should be capitalized`
}
