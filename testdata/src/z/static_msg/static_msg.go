package static_msg

import (
	"fmt"

	"go.uber.org/zap"
)

const constMsg = "constant message"

var varMsg = "variable message"

func tests(logger *zap.Logger) {
	logger.Info(constMsg)                      // OK
	logger.Info("static message")              // OK
	logger.Info(varMsg)                        // want `message should be a string literal or a constant`
	logger.Info(fmt.Sprintf("dynamic: %d", 1)) // want `message should be a string literal or a constant`
}
