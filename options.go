package azathoth

import "time"

type AzOption func(*Azathoth)

func WithDefaultConfig() AzOption {
	config := GetDefaultConfig()
	return func(a *Azathoth) {
		a.showMilliseconds = config.ShowMilliseconds
		a.workInterval = config.WorkIntervalTime
		a.restInterval = config.RestIntervalTime
	}
}

func WithWorkRestPeriod(w, r time.Duration) AzOption {
	return func(a *Azathoth) {
		a.workInterval = time.Time{}.Add(w)
		a.restInterval = time.Time{}.Add(r)
	}
}

type AzConfigOption func(*AzathothConfig)

func WithMilliseconds() AzConfigOption {
	return func(ao *AzathothConfig) {
		ao.ShowMilliseconds = true
	}
}
