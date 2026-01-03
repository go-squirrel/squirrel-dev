package config

type Auth struct {
	DefaultUn string
	DefaultPw string
	Jwt       Jwt
}

type Jwt struct {
	SigningKey string
	Expired    int
}
