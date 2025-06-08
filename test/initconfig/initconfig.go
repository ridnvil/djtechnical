package initconfig

import (
	"context"
	"github.com/caarlos0/env/v11"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type EnvConfig struct {
	Host      string `env:"DATABASE_DSN" envDefault:"host=192.168.100.73 user=ridwan password=M1r34cl3 dbname=payslip port=5432 sslmode=disable TimeZone=Asia/Jakarta"`
	SecretKEY string `env:"SECRET_KEY" envDefault:"M1r34cl3"`
	RedisAddr string `env:"REDIS_ADDR" envDefault:"192.168.100.73:4004"`
}

func InitDatabase() *gorm.DB {
	var envCfg EnvConfig
	if err := env.Parse(&envCfg); err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(envCfg.Host), &gorm.Config{
		//Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	return db
}

func InitRedis(ctx context.Context) *redis.Client {
	var envCfg EnvConfig
	if err := env.Parse(&envCfg); err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     envCfg.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	return client
}
