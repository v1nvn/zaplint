package allow_dynamic_msg

import (
	"fmt"

	"go.uber.org/zap"
)

const constMsg = "constant message"

var varMsg = "variable message"

func tests(logger *zap.Logger) {
	logger.Info(constMsg)                      // OK
	logger.Info("static message")              // OK
	logger.Info(varMsg)                        // OK
	logger.Info(fmt.Sprintf("dynamic: %d", 1)) // OK
}
