package main

import (
	"github.com/clearcodecn/wetalk/cmd/web"
	"github.com/urfave/cli"
	"log"
	"os"
)

var (
	app *cli.App
)

const appName = "wetalk"
const version = "dev"

func init() {
	app = cli.NewApp()
	app.Name = appName
	app.HelpName = appName
	app.Usage = "a modern im app"
	app.Version = version

	app.Commands = []cli.Command{
		web.Web,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
