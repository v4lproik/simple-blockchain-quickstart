package utils

type ApiConf struct {
	Env    string `env:"SBQ_ENV,required"`
	Server struct {
		Address string `env:"SBQ_SERVER_ADDRESS,required"`
		Port    int    `env:"SBQ_SERVER_PORT,required"`
		Options struct {
			IsSsl    bool   `env:"SBQ_SERVER_IS_SSL" envDefault:"false"`
			CertFile string `env:"SBQ_SERVER_CERT_FILE"`
			KeyFile  string `env:"SBQ_SERVER_KEY_FILE"`
		}
		HttpCors struct {
			AllowedOrigins []string `env:"SBQ_SERVER_HTTP_CORS_ALLOWED_ORIGINS,required" envSeparator:","`
			AllowedMethods []string `env:"SBQ_SERVER_HTTP_CORS_ALLOWED_METHODS,required" envSeparator:","`
			AllowedHeaders []string `env:"SBQ_SERVER_HTTP_CORS_ALLOWED_HEADERS,required" envSeparator:","`
		}
	}
}
