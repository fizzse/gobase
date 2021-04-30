package server

import (
	"log"
	"os"

	"github.com/fizzse/gobase/internal/gobase/server/rpc"

	"github.com/fizzse/gobase/internal/gobase/server/consumer"
	"github.com/fizzse/gobase/internal/gobase/server/rest"
	"github.com/fizzse/gobase/pkg/cache/redis"
	"github.com/fizzse/gobase/pkg/db"
	"github.com/fizzse/gobase/pkg/logger"
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
		configPath = "config/config_prod.yaml"
	case testEnv:
		configPath = "config/config_test.yaml"
	default:
		configPath = "config/config.yaml"
	}

	log.Printf("server get env: ENV_CLUSTER value: %s use config path: %s\n", env, configPath)

	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("viper read failed: ", err)
	}
}

func LoadLoggerConfig() *logger.Config {
	config := &logger.Config{Drive: logger.ZapStdDrive, Level: -1}

	if err := viper.UnmarshalKey("logger", config); err != nil {
		log.Printf("viper get: %s config failed %v : use defalut config\n", "logger", err)
	}

	return config
}

func LoadRestConfig() *rest.Config {
	return &rest.Config{
		Host:       "0.0.0.0",
		Port:       8080,
		DebugModel: true,
	}
}

func LoadGrpcConfig() *rpc.Config {
	return &rpc.Config{
		Host:       "0.0.0.0",
		Port:       8081,
		DebugModel: true,
	}
}

func LoadDbConfig() *db.Config {
	return &db.Config{
		Drive:    "mysql",
		Address:  "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "s",
		DbName:   "mysql",
		Charset:  "utf8",
	}
}

func LoadRedisConfig() *redis.Config {
	return &redis.Config{
		Host:     "127.0.0.1",
		Port:     "6379",
		Password: "s",
	}
}

func LoadConsumerConfig() *consumer.WorkerConfig {
	config := &consumer.WorkerConfig{
		Broker:      []string{"127.0.0.1:9200"},
		WorkerCount: 3,
		Topic:       "hello",
	}
	if err := viper.UnmarshalKey("consumer", config); err != nil {
		log.Printf("viper get: %s config failed %v : use defalut config\n", "consumer", err)
	}

	return config
}
