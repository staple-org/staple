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
	flag.Parse()
}

func main() {
	if err := pkg.Serve(); err != nil {
		log.Fatal("Failure starting Stapler: ", err)
	}
}
