package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code.0x7r.com/scholar7r/harborman/internal/cfg"
	"code.0x7r.com/scholar7r/harborman/internal/handler"
	"code.0x7r.com/scholar7r/harborman/internal/notifier"
)

const (
	shutdownTimeout = 3 * time.Second
	readTimeout     = 30 * time.Second
	writeTimeout    = 30 * time.Second
	idleTimeout     = 60 * time.Second
)

type option struct {
	cfgPath string
}

func fromArg() option {
	var opts option

	flag.StringVar(
		&opts.cfgPath,
		"c",
		"./harborman.yaml",
		"specify configuration file path",
	)

	flag.Parse()

	return opts
}

func main() {
	opts := fromArg()

	if opts.cfgPath == "" {
		flag.Usage()
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	c, err := cfg.FromFile(opts.cfgPath)
	if err != nil {
		slog.Error(
			"failed to load configuration",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	client := &http.Client{Timeout: readTimeout}
	notifiers := notifier.FromCfg(client, c.Notifiers)
	handler := handler.NewNotifyHandler(notifiers)

	server := &http.Server{
		Addr:         c.Listen,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			slog.Error(
				"failed to listen",
				slog.String("addr", c.Listen),
				slog.String("error", err.Error()),
			)

			os.Exit(1)
		}
	}()

	slog.Info(
		"harborman service listening",
		slog.String("addr", c.Listen),
	)

	<-sigChan
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		slog.Error(
			"server forced to shutdown",
			slog.String("error", err.Error()),
		)

		//nolint: gocritic // force shutdown do not cares
		os.Exit(1)
	}

	slog.Info("harborman service stopped")
}
