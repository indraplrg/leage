package configs

type Config struct {
	AppName string `mapstructure:"appName"`
	BaseUrl string `mapstructure:"baseUrl"`

	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Database struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
		SslMode string `mapstructure:"sslMode"`
		} `mapstructure:"database"`

	SMTP struct {
		Host string `mapstructure:"host"`
		Port int `mapstructure:"port"`
		SenderName string `mapsctructure:"senderName"`
	} `mapstructure:"smtp"`

	REDIS struct {
		Host string `mapstructure:"host"`
		Port int `mapstructure:"port"`
	} `mapstructure:"redis"`
}