package config

type Apiserver struct {
	Http Http
}

type Http struct {
	Scheme  string
	Server  string
	BaseUri string
}
