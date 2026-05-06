package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/cronwatch/internal/alerter"
	"github.com/example/cronwatch/internal/config"
	"github.com/example/cronwatch/internal/notifier"
	"github.com/example/cronwatch/internal/runner"
	"github.com/example/cronwatch/internal/tracker"
	"github.com/example/cronwatch/internal/watcher"
)

func main() {
	cfgPath := flag.String("config", "cronwatch.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	sender, err := alerter.FromConfig(cfg)
	if err != nil {
		log.Fatalf("failed to build alerter: %v", err)
	}

	trk := tracker.New()
	ntf := notifier.New(trk, sender)
	disp := runner.NewDispatcher(trk, ntf)
	wch := watcher.New(cfg, trk, ntf)

	log.Printf("cronwatch starting, monitoring %d job(s)", len(cfg.Jobs))

	disp.Start(cfg)
	wch.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down cronwatch")
	wch.Stop()
	disp.Stop()
}
