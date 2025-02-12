package config

import "time"

type App struct {
	Port               int // start app port
	IsDebug            bool
	UpdateLiveInterval time.Duration
	MaxOddValue        float64
}
