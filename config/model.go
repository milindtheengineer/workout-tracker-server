package config

type Config struct {
	Debug                bool
	DatabaseURI          string
	DatabaseName         string
	DatabaseReadTimeout  uint
	DatabaseWriteTimeout uint
	DatabaseUserName     string
	DatabasePassword     string
	DatabaseAuthSource   string
	Token                string
}
