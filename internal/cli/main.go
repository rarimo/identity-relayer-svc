package cli

import (
	"context"

	"github.com/rarimo/identity-relayer-svc/internal/services"
	"github.com/rarimo/identity-relayer-svc/internal/services/ingester"

	"github.com/rarimo/identity-relayer-svc/internal/config"

	"github.com/alecthomas/kingpin"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
)

func Run(args []string) bool {
	defer func() {
		if rvr := recover(); rvr != nil {
			logan.New().WithRecover(rvr).Fatal("app panicked")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.New(kv.MustFromEnv())
	log := cfg.Log()

	log.Info("Running service")

	app := kingpin.New("relayer-svc", "")

	runCmd := app.Command("run", "run command")
	runAllCmd := runCmd.Command("all", "run all services")

	// Running migrations
	migrateCmd := app.Command("migrate", "migrate command")
	migrateUpCmd := migrateCmd.Command("up", "migrate db up")
	migrateDownCmd := migrateCmd.Command("down", "migrate db down")

	cmd, err := app.Parse(args[1:])
	if err != nil {
		log.WithError(err).Fatal("failed to parse arguments")
	}

	switch cmd {
	case runAllCmd.FullCommand():
		go ingester.NewService(cfg, ingester.NewStateIngester(cfg)).Run(ctx)
		go ingester.NewService(cfg, ingester.NewGistIngester(cfg)).Run(ctx)
		err = services.NewServer(cfg).Run()
	case migrateUpCmd.FullCommand():
		err = MigrateUp(cfg)
	case migrateDownCmd.FullCommand():
		err = MigrateDown(cfg)
	default:
		log.Fatal("unknown command %s", cmd)
	}

	if err != nil {
		log.WithError(err).Error("failed to exec cmd")
		return false
	}

	return true
}
