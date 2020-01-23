package main

import (
	"flag"
	"log"

	"github.com/staple-org/staple/pkg/config"

	"github.com/staple-org/staple/pkg"
)

func init() {
	flag.BoolVar(&config.Opts.AutoTLS, "auto-tls", false, "--auto-tls")
	flag.StringVar(&config.Opts.CacheDir, "cache-dir", "", "--cache-dir /home/user/.server/.cache")
	flag.StringVar(&config.Opts.ServerKeyPath, "server-key-path", "", "--server-key-path /home/user/.server/server.key")
	flag.StringVar(&config.Opts.ServerCrtPath, "server-crt-path", "", "--server-crt-path /home/user/.server/server.crt")
	flag.StringVar(&config.Opts.Port, "port", "9998", "--port 443")
	flag.StringVar(&config.Opts.Hostname, "hostname", "localhost", "--hostname staple-clipper.org")
	flag.StringVar(&config.Opts.GlobalTokenKey, "token-key", "", "--token-key <random-data>")
	flag.StringVar(&config.Opts.Database.ConnectionURL,
		"database-connection-url",
		"postgresql://localhost/staples?user=staple&password=password123",
		"--database-connection-url postgresql://localhost/staples?user=staple&password=password123")
	flag.StringVar(&config.Opts.Mailer.Domain, "mg-domain", "", "--mg-domain <MG_DOMAIN>")
	flag.StringVar(&config.Opts.Mailer.ApiKey, "mg-api-key", "", "--mg-api-key <MG_API_KEY>")
	flag.Parse()
}

func main() {
	if err := pkg.Serve(); err != nil {
		log.Fatal("Failure starting Stapler: ", err)
	}
}
