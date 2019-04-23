package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	//Log is shared acorss the application using zap
	Log      *zap.Logger
	onceInit sync.Once
)

//Init initialize
func Init(lvl int) error {
	var err error

	onceInit.Do(func() {
		// First, define our level-handling logic.
		defaultLvl := zapcore.Level(lvl)

		//Prioritize more higher level log
		high := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		low := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= defaultLvl && lvl < zapcore.ErrorLevel
		})
		infos := zapcore.Lock(os.Stdout)
		errors := zapcore.Lock(os.Stderr)
		//zap log output configuration
		ecfg := zap.NewProductionEncoderConfig()
		ecfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(ecfg)

		core := zapcore.NewTee(
			zapcore.NewCore(encoder, errors, high),
			zapcore.NewCore(encoder, infos, low),
		)

		//construct log
		Log = zap.New(core)
		zap.RedirectStdLog(Log)
	})

	return err
}
