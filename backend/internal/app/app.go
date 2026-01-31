// Package app - Application entry point
package app

import (
	"fmt"

	"github.com/asliddinberdiev/eirsystem/config"
	httpDelivery "github.com/asliddinberdiev/eirsystem/internal/delivery/http"
	"github.com/asliddinberdiev/eirsystem/internal/repository"
	"github.com/asliddinberdiev/eirsystem/internal/server"
	"github.com/asliddinberdiev/eirsystem/internal/service"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/minio"
	"github.com/asliddinberdiev/eirsystem/pkg/postgres"
	"github.com/asliddinberdiev/eirsystem/pkg/redis"
	"github.com/asliddinberdiev/eirsystem/pkg/seed"
	"github.com/asliddinberdiev/eirsystem/pkg/telegram"
)

func New() {
	cfg, err := config.Load("./config")
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.Logger, cfg.App.IsDev())
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	telegram.Init(log.Named("Telegram"), cfg.App.TelegramBotToken, cfg.App.TelegramChatID)
	appLog := log.Named("APP")

	failOnError := func(title string, err error) {
		msg := fmt.Sprintf("🚨 <b>CRITICAL ERROR</b>\n\n🆔 <b>%s</b>\n❌ <pre>%v</pre>", title, err)
		telegram.Send(msg)
		telegram.Close()
		appLog.Fatal(title, logger.Error(err))
	}

	defer func() {
		if r := recover(); r != nil {
			telegram.Send(fmt.Sprintf("🚨 <b>PANIC RECOVERED</b>\n\n<pre>%v</pre>", r))
			telegram.Close()
			appLog.Fatal("Application panicked", logger.Any("panic", r))
		}
	}()

	gormPsql, err := postgres.New(&cfg.Postgres, cfg.App.IsDev(), log.Named("GORM"))
	if err != nil {
		failOnError("Postgres connection failed", err)
	}

	sqlDB, err := gormPsql.DB()
	if err != nil {
		failOnError("Postgres sql.DB retrieval failed", err)
	}
	defer sqlDB.Close()
	appLog.Info("Connected to postgres")

	if err := postgres.RunMigrations(sqlDB, "./migrations"); err != nil {
		failOnError("Database migration failed", err)
	}
	appLog.Info("Database migrations applied")

	if err := seed.SeedSuperAdmin(log.Named("SEED"), gormPsql, cfg.SeedSuperAdmin); err != nil {
		failOnError("Super Admin seeding failed", err)
	}

	redisClient, err := redis.New(&cfg.Redis)
	if err != nil {
		failOnError("Redis connection failed", err)
	}
	defer redisClient.Close()
	appLog.Info("Connected to redis")

	minioClient, err := minio.New(&cfg.Minio, cfg.App.IsDev(), log.Named("MINIO"))
	if err != nil {
		failOnError("Minio connection failed", err)
	}
	appLog.Info("Connected to minio")

	repository := repository.New(cfg, log.Named("REPOSITORY"), gormPsql, redisClient)
	service := service.New(cfg, log.Named("SERVICE"), minioClient, repository)

	h := httpDelivery.New(cfg, log.Named("HTTP"), redisClient.Client, service)
	srv := server.New(&cfg.App, log.Named("SERVER"), h.InitRouter())

	if err := srv.Run(); err != nil {
		failOnError("Server forced to shutdown", err)
	}
}
