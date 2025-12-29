package config

type Config struct {
	AppName string `mapstructure:"appName"`
	BaseUrl string `mapstructure:"baseUrl"`

	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Database struct {
		Host string `mapstructure:"host"`
		User string `mapstructure:"user"`
		Password string `mapstrucure:"password"`
		Port string `mapstructure:"port"`
		DatabaseName string `mapstructure:"databaseName"`
		SslMode string `mapstructure:"sslMode"`
		} `mapstructure:"database"`

	SMTP struct {
		Host string `mapstructure:"host"`
		Port int `mapstructure:"port"`
		SenderName string `mapsctructure:"senderName"`
	} `mapstructure:"smtp"`

}