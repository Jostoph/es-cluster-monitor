package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	Log *zap.Logger

	syncInit sync.Once
)

func Init() {

	syncInit.Do(func() {

		// config
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder

		// encoder
		encoder := zapcore.NewJSONEncoder(cfg)

		// redirect level >= Error to stdout
		highPrio := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		// redirect level < Error to stderr
		lowPrio := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		})

		core := zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), highPrio),
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), lowPrio),
		)

		Log = zap.New(core)
		// zap.RedirectStdLog(Log)
	})
}
