package main

import (
	"flag"
	"log"

	"github.com/staple-org/staple/pkg"
	"github.com/staple-org/staple/pkg/config"
)

func init() {
	flag.BoolVar(&config.Opts.DevMode, "dev", false, "--dev")
	flag.BoolVar(&config.Opts.AutoTLS, "auto-tls", false, "--auto-tls")
	flag.StringVar(&config.Opts.CacheDir, "cache-dir", "", "--cache-dir /home/user/.server/.cache")
	flag.StringVar(&config.Opts.ServerKeyPath, "server-key-path", "", "--server-key-path /home/user/.server/server.key")
	flag.StringVar(&config.Opts.ServerCrtPath, "server-crt-path", "", "--server-crt-path /home/user/.server/server.crt")
	flag.StringVar(&config.Opts.Port, "port", "9998", "--port 443")
	flag.StringVar(&config.Opts.Hostname, "hostname", "localhost", "--hostname staple-clipper.org")
	flag.StringVar(&config.Opts.GlobalTokenKey, "token-key", "", "--token-key <random-data>")
	flag.StringVar(&config.Opts.Database.Hostname, "staple-db-hostname", "localhost", "--staple-db-hostname localhost")
	flag.StringVar(&config.Opts.Database.Database, "staple-db-database", "staples", "--staple-db-database staples")
	flag.StringVar(&config.Opts.Database.Username, "staple-db-username", "staple", "--staple-db-username staple")
	flag.StringVar(&config.Opts.Database.Password, "staple-db-password", "password123", "--staple-db-password password123")
	flag.StringVar(&config.Opts.Mailer.Domain, "mg-domain", "", "--mg-domain <MG_DOMAIN>")
	flag.StringVar(&config.Opts.Mailer.APIKey, "mg-api-key", "", "--mg-api-key <MG_API_KEY>")
	flag.BoolVar(&config.Opts.Debug, "debug", false, "--debug")
	flag.Parse()
}

func main() {
	if err := pkg.Serve(); err != nil {
		log.Fatal("Failure starting Stapler: ", err)
	}
}
