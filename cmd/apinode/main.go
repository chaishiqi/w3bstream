package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	"github.com/iotexproject/w3bstream/service/apinode"
	"github.com/iotexproject/w3bstream/service/apinode/config"
	"github.com/iotexproject/w3bstream/service/apinode/persistence"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to get config"))
	}
	cfg.Print()
	slog.Info("apinode config loaded")

	p, err := persistence.NewPersistence(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	apinode := apinode.NewAPINode(cfg, p)

	if err := apinode.Start(); err != nil {
		log.Fatal(err)
	}
	defer apinode.Stop()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
