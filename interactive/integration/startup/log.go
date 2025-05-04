package startup

import (
	"github.com/rwpp/RzWeLook/pkg/logger"
)

func InitLog() logger.LoggerV1 {
	return logger.NewNoOpLogger()
}
