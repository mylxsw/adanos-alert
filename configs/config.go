package configs

import (
	"encoding/json"
	"time"
)

type Config struct {
	MongoURI          string
	MongoDB           string
	APIToken          string
	UseLocalDashboard bool
	ConsoleColor      bool

	AggregationPeriod   time.Duration
	ActionTriggerPeriod time.Duration
}

func (conf *Config) Serialize() string {
	rs, _ := json.Marshal(conf)
	return string(rs)
}
