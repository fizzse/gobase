package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	ZapDrive    = "zap"
	ZapStdDrive = "zapStd"
)

type Config struct {
	Drive       string        `yaml:"drive"`       // 驱动
	Path        string        `yaml:"path"`        // 路径
	File        string        `yaml:"file"`        // 文件名
	MaxAge      time.Duration `yaml:"maxAge"`      // 最大保存时间
	MaxBackups  int           `yaml:"maxBackups"`  // 最大文件个数
	MaxFileSize int           `yaml:"maxFileSize"` // 文件大小
	Level       int8          `yaml:"level"`       // 级别
	Compress    bool          `yaml:"compress"`    // 是否压缩
	LocalTime   bool          `yaml:"localTime"`   // 是否使用当地时间
	//RotationTime time.Duration `yaml:"rotationTime"` //
}

func New(cfg *Config) (*zap.SugaredLogger, error) {
	switch cfg.Drive {
	case ZapDrive:
		return NewZapLogger(cfg)

	default:
		return NewZapSTDLogger(cfg)
	}
}

func NewZapLogger(cfg *Config) (*zap.SugaredLogger, error) {
	enConfig := zap.NewProductionEncoderConfig() //生成配置

	// 时间格式
	enConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	level := zapcore.Level(cfg.Level)

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", cfg.Path, cfg.File), //filePath
		MaxSize:    cfg.MaxFileSize,                          // megabytes
		MaxBackups: 10000,
		MaxAge:     int(cfg.MaxAge / 24), //days
		Compress:   cfg.Compress,         // disabled by default
		LocalTime:  cfg.LocalTime,
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(enConfig), //编码器配置
		w,                                //打印到控制台和文件
		level,                            //日志等级
	)

	logger := zap.New(core, zap.AddCaller())
	zapLogger := logger.Sugar()
	return zapLogger, nil
}

func NewZapSTDLogger(cfg *Config) (*zap.SugaredLogger, error) {
	enConfig := zap.NewProductionEncoderConfig() //生成配置
	enConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	level := zapcore.Level(cfg.Level)
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(level), // 日志级别
		Development:       true,                        // 开发模式，堆栈跟踪
		DisableStacktrace: true,                        // 关闭堆栈追踪
		Encoding:          "json",                      // 输出格式 console 或 json
		EncoderConfig:     enConfig,                    // 编码器配置
		// InitialFields:    map[string]interface{}{"serviceName": "spikeProxy"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout"}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	zapLogger := logger.Sugar()
	return zapLogger, nil
}
