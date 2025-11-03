package azathoth

import "time"

type AzathothConfig struct {
	ShowMilliseconds bool
	WorkIntervalTime time.Time
	RestIntervalTime time.Time
}

func NewAzathothConfig(options ...AzConfigOption) AzathothConfig {
	config := AzathothConfig{}
	for _, opt := range options {
		opt(&config)
	}
	return config
}

func GetDefaultConfig() AzathothConfig {
	config := NewAzathothConfig(WithMilliseconds())
	config.WorkIntervalTime = time.Time{}.Add(25 * time.Minute)
	config.RestIntervalTime = time.Time{}.Add(5 * time.Minute)
	config.ShowMilliseconds = true
	return config
}
