package logx

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
)

type Logger struct {
	*zap.Logger
}

type Config struct {
	//格式: json 或 console
	Encoding string `json:"encoding"`
	// debug info warn error fatal
	Level string `json:"level"`
	// 日志文件路径
	Filename string `json:"filename"`
	// 每个日志文件保存的最大尺寸 单位：M
	MaxSize int `json:"maxsize"`
	// 日志文件最多保存多少个备份
	MaxBackups int `json:"maxbackups"`
	// 文件最多保存多少天
	MaxAge int `json:"maxage"`
	//是否打印屏幕
	Console bool `json:"console"`
	//默认字段
	InitialFields map[string]interface{} `json:"initialFields"`
}

func ParseConfig(data string) *Config {
	return ParseConfigByte([]byte(data))
}

func ParseConfigByte(data []byte) *Config {
	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		log.Fatalf("logx parse config: %v, error:%v", string(data), err.Error())
	}
	return cfg
}

func NewDefaultConfig() *Config {
	defaultConf := &Config{
		Encoding:      "json",
		Level:         "info",
		Filename:      "logx.log",
		MaxSize:       128,
		MaxBackups:    30,
		MaxAge:        30,
		Console:       true,
		InitialFields: make(map[string]interface{}),
	}
	return defaultConf
}

func NewZapLogger(cfg *Config) *Logger {
	if cfg == nil {
		cfg = NewDefaultConfig()
	}

	hook := lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   true, // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch cfg.Level {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
	case "fatal":
		atomicLevel.SetLevel(zap.FatalLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	var syncs zapcore.WriteSyncer
	if cfg.Console {
		syncs = zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(&hook))
	} else {
		syncs = zapcore.AddSync(&hook)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		syncs,
		atomicLevel,
	)
	var options []zap.Option = []zap.Option{
		// 开启开发模式，堆栈跟踪
		zap.AddCaller(),
		// 开启文件及行号
		zap.Development(),
	}

	//默认字段
	for k, v := range cfg.InitialFields {
		options = append(options, zap.Fields(zap.Any(k, v)))
	}

	return &Logger{
		zap.New(core, options...)}
}
