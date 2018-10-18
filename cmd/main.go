package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
)

var (
	appName string = "popmeet-api"
	version string = "0.0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = version
	app.Copyright = "(c) 2018 - Ricardo Pinto"
	app.Usage = "Popmeet API"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "listen, l",
			Value:  "0.0.0.0:8000",
			Usage:  "Address and port on which the API will accept HTTP requests",
			EnvVar: "LISTEN",
		},
		cli.StringFlag{
			Name:   "log-folder, lf",
			Value:  "",
			Usage:  `Log folder path for access and application logging. Default "stdout"`,
			EnvVar: "LOG_FOLDER",
		},
		cli.StringFlag{
			Name:   "database-file, d",
			Value:  "",
			Usage:  "Database configuration used by the API to connect to database",
			EnvVar: "DATABASE_FILE",
		},
		cli.StringFlag{
			Name:   "security-file, sf",
			Value:  "",
			Usage:  "Security configuration",
			EnvVar: "SECURITY_FILE",
		},
		cli.StringFlag{
			Name:   "ssl-cert",
			Value:  "",
			Usage:  "Define SSL certificate to accept HTTPS requests",
			EnvVar: "SSL_CERT",
		},
		cli.StringFlag{
			Name:   "ssl-key",
			Value:  "",
			Usage:  "Define SSL key to accept HTTPS requests",
			EnvVar: "SSL_KEY",
		},
	}

	app.Action = Handler
	app.Run(os.Args)
}
