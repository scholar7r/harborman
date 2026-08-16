// Package cfg provides internal configuration implementation
package cfg

import (
	"os"

	"go.yaml.in/yaml/v4"
)

type Cfg struct {
	Debug     bool          `yaml:"debug"`
	Listen    string        `yaml:"listen"`
	Notifiers []NotifierCfg `yaml:"notifiers"`
}

type NotifierCfg struct {
	Type NotifierType `yaml:"type"`
	URL  string       `yaml:"url"`
	// authorization header the notifier requires, unverified when empty
	Authorization string `yaml:"authorization"`
}

type NotifierType string

const (
	NotifierTypeDiscord NotifierType = "discord"
	NotifierTypeLark    NotifierType = "lark"
)

func FromFile(filePath string) (*Cfg, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Cfg
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
