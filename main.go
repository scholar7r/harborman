package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"code.0x7r.com/scholar7r/harborman/internal/cfg"
	"code.0x7r.com/scholar7r/harborman/internal/handler"
	"code.0x7r.com/scholar7r/harborman/internal/notifier"
)

const (
	readTimeout  = 30 * time.Second
	writeTimeout = 30 * time.Second
	idleTimeout  = 60 * time.Second
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
		Level: slog.LevelInfo,
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

	slog.Info(
		"harborman service listening",
		slog.String("addr", c.Listen),
	)

	notifiers := notifier.FromCfg(c.Notifiers)
	handler := handler.NewNotifyHandler(notifiers)

	server := &http.Server{
		Addr:         c.Listen,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	if err = server.ListenAndServe(); err != nil {
		slog.Error(
			"failed to listen",
			slog.String("addr", c.Listen),
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}
}
