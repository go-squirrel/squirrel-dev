package config

type Agent struct {
	Http Http
}

type Http struct {
	Scheme  string
	BaseUrl string
}
