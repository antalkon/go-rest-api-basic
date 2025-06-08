package config

type App struct {
	Name    string `env:"NAME,required"`
	Version string `env:"VERSION,required"`
}

type HTTP struct {
	Port           string `env:"PORT,required"`
	UsePreforkMode bool   `env:"USE_PREFORK_MODE" envDefault:"false"`
}

type Log struct {
	Level string `env:"LEVEL" envDefault:"info"`
}

type PG struct {
	PoolMax int    `env:"POOL_MAX" envDefault:"10"`
	URL     string `env:"URL,required"`
}

type GRPC struct {
	Port string `env:"PORT,required"`
}

type RMQ struct {
	ServerExchange string `env:"RPC_SERVER,required"`
	ClientExchange string `env:"RPC_CLIENT,required"`
	URL            string `env:"URL,required"`
}

type Metrics struct {
	Enabled bool `env:"ENABLED" envDefault:"true"`
}

type Swagger struct {
	Enabled bool `env:"ENABLED" envDefault:"false"`
}

type Redis struct {
	Addr     string `env:"ADDR,required"`
	Password string `env:"PASSWORD"`
	DB       int    `env:"DB" envDefault:"0"`
}

type S3 struct {
	Endpoint  string `env:"ENDPOINT,required"`
	AccessKey string `env:"ACCESS_KEY,required"`
	SecretKey string `env:"SECRET_KEY,required"`
	Bucket    string `env:"BUCKET,required"`
	UseSSL    bool   `env:"USE_SSL" envDefault:"false"`
}
