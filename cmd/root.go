package main

import (
	"flag"
	"log"

	"github.com/staple-org/staple/pkg"
)

func init() {
	flag.BoolVar(&pkg.Opts.AutoTLS, "auto-tls", false, "--auto-tls")
	flag.StringVar(&pkg.Opts.CacheDir, "cache-dir", "", "--cache-dir /home/user/.server/.cache")
	flag.StringVar(&pkg.Opts.ServerKeyPath, "server-key-path", "", "--server-key-file /home/user/.server/server.key")
	flag.StringVar(&pkg.Opts.ServerCrtPath, "server-crt-path", "", "--server-crt-file /home/user/.server/server.crt")
	flag.StringVar(&pkg.Opts.Port, "port", "9998", "--port 443")
	flag.StringVar(&pkg.Opts.Hostname, "hostname", "localhost", "--hostname staple-clipper.org")
	flag.StringVar(&pkg.Opts.GlobalTokenKey, "token-key", "", "--token-key <random-data>")
	flag.Parse()
}

func main() {
	if err := pkg.Serve(); err != nil {
		log.Fatal("Failure starting Stapler: ", err)
	}
}
