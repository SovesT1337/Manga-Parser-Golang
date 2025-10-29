package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"go_scripts/config"
	"go_scripts/database"
	"go_scripts/internal/bot"
	"go_scripts/internal/fsm"
	"go_scripts/internal/logger"
	"go_scripts/internal/scheduler"
)

func main() {
	logger.BotInfo("starting telegram-bot service")
	_ = godotenv.Load()
	os.Setenv("SERVICE", "telegram-bot")
	c, err := config.Load()
	if err != nil {
		logger.BotError("config: %v", err)
		os.Exit(1)
	}

	botURL := c.TelegramAPI + c.TelegramToken
	if err := database.InitDB(c.DatabaseDSN); err != nil {
		logger.DatabaseError("init db: %v", err)
		os.Exit(1)
	}

	manager := fsm.NewManager(24*time.Hour, 10*time.Minute)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start scheduler
	sched := &scheduler.Runner{BotURL: botURL, ChannelID: c.SchedulerTelegramChannelID, IntervalSec: c.SchedulerIntervalSec, SubscribeURL: c.SubscribeLinkURL}
	go sched.Run(ctx)

	// Start bot updates loop
	b := bot.New(botURL, manager)
	go b.Run(ctx)

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
	manager.Shutdown()
}
