package config

type Auth struct {
	Jwt Jwt
}

type Jwt struct {
	SigningKey string
	Expired    int
}
