package pkg

// Config contains configuration options for the inner server of Beemo.
type Config struct {
	AutoTLS       bool
	CacheDir      string
	ServerKeyPath string
	ServerCrtPath string
	Port          string
	Hostname      string
}

// Opts represents server side configuration for Beemo.
var Opts = Config{}
