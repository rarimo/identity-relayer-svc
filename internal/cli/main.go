package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gitlab.com/rarimo/relayer-svc/internal/services"

	"gitlab.com/rarimo/relayer-svc/internal/services/relayer"

	"gitlab.com/rarimo/relayer-svc/internal/config"

	"github.com/alecthomas/kingpin"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v2/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

func Run(args []string) {
	defer func() {
		if rvr := recover(); rvr != nil {
			logan.New().WithRecover(rvr).Fatal("app panicked")

		}
		os.Stdout.Sync()
		os.Stderr.Sync()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := config.New(kv.MustFromEnv())
	log := cfg.Log()

	var wg sync.WaitGroup
	run := func(f func(config.Config, context.Context)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				os.Stdout.Sync()
				os.Stderr.Sync()

				if rvr := recover(); rvr != nil {
					err := errors.FromPanic(rvr)
					logan.New().WithError(err).Fatal("one of the services panicked")
				}
			}()
			f(cfg, ctx)
		}()
	}

	app := kingpin.New("relayer-svc", "")

	runCmd := app.Command("run", "run command")
	runAllCmd := runCmd.Command("all", "run all services")
	runRelayerCmd := runCmd.Command("relayer", "run relayer routines only")
	runSchedulerCmd := runCmd.Command("scheduler", "run scheduler only")

	cmd, err := app.Parse(args[1:])
	if err != nil {
		log.WithError(err).Fatal("failed to parse arguments")
	}

	switch cmd {
	case runAllCmd.FullCommand():
		log.Info("starting all services")
		run(services.RunScheduler)
		run(relayer.Run)
		run(services.RunQueueCleaner)
	case runSchedulerCmd.FullCommand():
		log.Info("starting scheduler")
		run(services.RunScheduler)
	case runRelayerCmd.FullCommand():
		log.Info("starting all services")
		run(relayer.Run)
		run(services.RunQueueCleaner)
	default:
		log.Fatal("unknown command %s", cmd)
	}

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	// making WaitGroup usable in select
	wgch := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgch)
	}()

	select {
	// listening for runners stop
	case <-wgch:
		cfg.Log().Warn("all services stopped")
	// listening for OS signals
	case <-gracefulStop:
		cfg.Log().Info("received signal to stop")
		cancel()
		<-wgch
	}
}
