package logx_test

import (
	"github.com/GamerCode/logx"
	"testing"
)

func TestLoggerT(t *testing.T) {
	conf := `
		{
			"encoding": "json",
			"level": "info",
			"filename": "logx.log",
			"maxsize": 128,
			"maxbackups": 30,
			"maxage": 30,
			"console": true,
			"initialFields": {
				"app": "LogX",
				"version": "1.0.0"
			}
		}
	`
	
	cfg := logx.ParseConfig(conf)

	//
	// cfg := logx.NewDefaultConfig()

	l := logx.NewZapLogger(cfg)
	l.With(logx.String("host", "10.8.8.8")).Info("info message")
}
