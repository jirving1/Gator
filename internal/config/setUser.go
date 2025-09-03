package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func (cfg *Config) SetUser(name string) error {

	if len(name) == 0 {
		return fmt.Errorf("name too short")
	}
	err := write(cfg, name)
	if err != nil {
		return err
	}
	fmt.Println("user set to", name)
	return nil
}

func write(cfg *Config, name string) error {
	cfg.CurrentUsername = name
	newJson, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	fi, err := os.Lstat(path)
	if err != nil {
		return err
	}
	filePerms := fi.Mode().Perm()
	err = os.WriteFile(path, newJson, filePerms)
	if err != nil {
		return err
	}
	return nil

}
