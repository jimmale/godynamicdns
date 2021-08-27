package config

import "time"

type Configuration struct {
	Debug   bool
	Domains []Domain `toml:"domain"`
}

type Domain struct {
	Username  string
	Password  string
	Hostname  string
	Frequency duration
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
