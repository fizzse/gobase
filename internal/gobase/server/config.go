package server

import (
	"log"
	"os"

	"github.com/fizzse/gobase/internal/gobase/server/consumer"
	"github.com/fizzse/gobase/internal/gobase/server/rest"
	"github.com/fizzse/gobase/internal/gobase/server/rpc"

	"github.com/fizzse/gobase/pkg/cache/redis"
	"github.com/fizzse/gobase/pkg/db"
	"github.com/fizzse/gobase/pkg/logger"
	"github.com/fizzse/gobase/pkg/trace"
	"github.com/spf13/viper"
)

const (
	envKey  = "ENV_CLUSTER"
	prodEnv = "prod"
	devEnv  = "dev"
	testEnv = "test"
)

func init() {
	// 读取环境变量 根据不同的环境变量 读不同的配置文件
	env := os.Getenv(envKey)
	configPath := "config/config.yaml"

	switch env {
	case prodEnv:
		configPath = "config/prod.yaml"
	case testEnv:
		configPath = "config/test.yaml"
	default:
		configPath = "config/dev.yaml"
	}

	log.Printf("server get env: ENV_CLUSTER value: %s use config path: %s\n", env, configPath)

	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("viper read failed: ", err)
	}
}

func loadTraceConfig() *trace.Config {
	config := &trace.Config{
		Agent:       "127.0.0.1:6831",
		Sampling:    "http://127.0.0.1:5778/sampling",
		ServiceName: "gobase",
		LogSpan:     false,
		Type:        "const",
		Param:       1,
	}

	configType := "jaeger"
	if err := viper.UnmarshalKey(configType, config); err != nil {
		log.Printf("viper get: %s config failed %v : use defalut config\n", configType, err)
	}

	log.Printf("%s config info: %+v\n", configType, config)
	return config
}

func loadLoggerConfig() *logger.Config {
	config := &logger.Config{Drive: logger.ZapStdDrive, Level: 0}

	configType := "logger"
	if err := viper.UnmarshalKey(configType, config); err != nil {
		log.Printf("viper get: %s config failed %v : use defalut config\n", configType, err)
	}

	return config
}

func loadRestConfig() *rest.Config {
	config := &rest.Config{
		Host:       "0.0.0.0",
		Port:       8080,
		DebugModel: true,
	}

	configType := "rest"
	if err := viper.UnmarshalKey(configType, config); err != nil {
		log.Printf("viper get: %s config failed %v : use defalut config\n", configType, err)
	}

	log.Printf("%s config info: %+v\n", configType, config)
	return config
}

func loadGrpcConfig() *rpc.Config {
	config := &rpc.Config{
		Host:       "0.0.0.0",
		Port:       8081,
		DebugModel: true,
	}

	configType := "grpc"
	if err := viper.UnmarshalKey(configType, config); err != nil {
		log.Printf("viper get: %s config failed %v : use defalut config\n", configType, err)
	}

	log.Printf("%s config info: %+v\n", configType, config)
	return config
}

func loadDbConfig() *db.Config {
	config := &db.Config{
		Drive:    "mysql",
		Address:  "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "s",
		DbName:   "mysql",
		Charset:  "utf8",
	}

	configType := "mysql"
	if err := viper.UnmarshalKey(configType, &config); err != nil {
		log.Printf("viper get: %s config failed %v : use defalut config\n", configType, err)
	}

	log.Printf("%s config info: %+v\n", configType, config)
	return config
}

func loadRedisConfig() *redis.Config {
	config := &redis.Config{}
	config.Mode = redis.ModeSingle
	config.Single.Addr = "127.0.0.1:6379"
	config.Password = "s"

	configType := "redis"
	if err := viper.UnmarshalKey(configType, &config); err != nil {
		log.Printf("viper get: %s config failed %v : use defalut config\n", configType, err)
	}

	log.Printf("%s config info: %+v\n", configType, config)
	return config
}

func loadConsumerConfig() *consumer.WorkerConfig {
	config := &consumer.WorkerConfig{
		Broker:      []string{"127.0.0.1:9092"},
		WorkerCount: 3,
		Topic:       "hello",
	}
	if err := viper.UnmarshalKey("consumer", config); err != nil {
		log.Printf("viper get: %s config failed %v : use defalut config\n", "consumer", err)
	}

	return config
}
