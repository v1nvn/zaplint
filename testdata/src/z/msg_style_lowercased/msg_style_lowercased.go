package msg_style_lowercased

import "go.uber.org/zap"

func tests(logger *zap.Logger) {
	logger.Info("lowercase message") // OK
	logger.Info("Message")           // want `message should be lowercased`
}
