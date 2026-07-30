package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"code.0x7r.com/scholar7r/harborman/internal/cfg"
	"code.0x7r.com/scholar7r/harborman/internal/handler"
	"code.0x7r.com/scholar7r/harborman/internal/notifier"
)

var cfgPath string

func FromArg() {
	flag.StringVar(
		&cfgPath,
		"c",
		"./harborman.yaml",
		"specify configuration file path",
	)

	flag.Parse()
}

func main() {
	FromArg()

	if cfgPath == "" {
		flag.Usage()
	}

	c, err := cfg.FromFile(cfgPath)
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

	if err = http.ListenAndServe(c.Listen, handler); err != nil {
		slog.Error(
			"failed to listen",
			slog.String("addr", c.Listen),
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}
}
