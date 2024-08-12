package config_watcher

import (
	"log"
	"log-tracer/internal/config"

	"github.com/fsnotify/fsnotify"
)

type ConfigWatcher struct {
	ConfigPath string
	ReloadFunc func(*config.ProducerConfig)
}

func NewConfigWatcher(configPath string, reloadFunc func(*config.ProducerConfig)) *ConfigWatcher {
	return &ConfigWatcher{
		ConfigPath: configPath,
		ReloadFunc: reloadFunc,
	}
}

func (cw *ConfigWatcher) WatchConfig() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(cw.ConfigPath)
	if err != nil {
		return err
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("Config file modified:", event.Name)
				cw.ReloadConfig()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Println("Error watching config file:", err)
		}
	}
}

func (cw *ConfigWatcher) ReloadConfig() {
	config, err := config.LoadProducerConfig(cw.ConfigPath)
	if err != nil {
		log.Println("Failed to reload config:", err)
		return
	}
	cw.ReloadFunc(config)
}
