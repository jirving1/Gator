package main

import (
	"blogaggregator/internal/config"

	"blogaggregator/internal/database"
)

type State struct {
	cfgPtr *config.Config
	db     *database.Queries
}
