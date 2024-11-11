package env

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createLogger(this_service string) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// level := strings.ToLower(os.Getenv("DEBUG_LEVEL"))
	// WARNING: THIS MUST RUN AFTER THE ENV IS PARSED/SET UP
	s, _ := Env.GetUserService(this_service)
	level := s.GetParamString("debug_level")

	zap_level := zap.InfoLevel
	switch level {
	case "debug":
		zap_level = zap.DebugLevel
	case "info":
		zap_level = zap.InfoLevel
	case "warn":
		zap_level = zap.WarnLevel
	case "error":
		zap_level = zap.ErrorLevel
	default:
		zap_level = zap.InfoLevel
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap_level),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(config.Build())
}

func SetupLogging(this_service string) {
	zap.ReplaceGlobals(zap.Must(createLogger(this_service), nil))
}
