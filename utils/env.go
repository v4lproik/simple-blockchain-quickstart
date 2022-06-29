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
	Auth struct {
		IsAuthenticationActivated bool `env:"SBQ_IS_AUTHENTICATION_ACTIVATED,required"`
		IsJwksEndpointActivated   bool `env:"SBQ_IS_JKMS_ACTIVATED,required"`
		Jwt                       struct {
			Signing struct {
				KeyPath   string `env:"SBQ_JWT_KEY_PATH,required"`
				KeyId     string `env:"SBQ_JWT_KEY_ID,required"`
				ExpiresIn int    `env:"SBQ_JWT_EXPIRES_IN_HOURS,required"`
				Domain    string `env:"SBQ_JWT_DOMAIN,required"`
				Audience  string `env:"SBQ_JWT_AUDIENCE,required"`
				Issuer    string `env:"SBQ_JWT_ISSUER,required"`
				Algo      string `env:"SBQ_JWT_ALGO,required"`
			}
			Verifying struct {
				JkmsUrl                        string `env:"SBQ_JWT_JKMS_URL,required"`
				JkmsRefreshCacheIntervalInMin  int    `env:"SBQ_JWT_JKMS_REFRESH_CACHE_INTERVAL_IN_MIN,required"`
				JkmsRefreshCacheRateLimitInMin int    `env:"SBQ_JWT_JKMS_REFRESH_CACHE_RATE_LIMIT_IN_MIN,required"`
				JkmsRefreshCacheTimeoutInSec   int    `env:"SBQ_JWT_JKMS_REFRESH_CACHE_TIMEOUT_IN_SEC,required"`
			}
		}
	}
}
