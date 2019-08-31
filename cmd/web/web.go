package web

import (
	"fmt"
	"github.com/clearcodecn/wetalk/api/http"
	"github.com/clearcodecn/wetalk/configs"
	"github.com/clearcodecn/wetalk/pkg/fs"
	"github.com/clearcodecn/wetalk/pkg/mail"
	"github.com/clearcodecn/wetalk/pkg/sms"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	Web cli.Command
)

func init() {
	Web = cli.Command{
		Name:   "web",
		Usage:  "web server management",
		Action: runWeb,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "c,config",
				Usage:    "config file path",
				EnvVar:   "WETALK_WEB_CONFIG",
				FilePath: "/etc/wetalk/web.yaml",
				Required: true,
				Value:    "/etc/wetalk/web.yaml",
			},
		},
	}
}

func runWeb(ctx *cli.Context) error {
	data, err := ioutil.ReadFile(ctx.String("c"))
	if err != nil {
		return err
	}

	config, err := configs.ParseWebConfig(data)
	if err != nil {
		return errors.Wrap(err, "unable to parse config")
	}

	var opts []http.Option

	// load sms
	if config.SmsConfig.Enable {
		if _, ok := sms.Senders[config.SmsConfig.Driver]; !ok {
			return fmt.Errorf("can't find sms driver: %s", config.SmsConfig.Driver)
		}

		data, err := yaml.Marshal(config.SmsConfig.Content)
		if err != nil {
			return err
		}
		sender := sms.Senders[config.SmsConfig.Driver]()
		if err := sender.Init(data); err != nil {
			return fmt.Errorf("init sms sender failed: %s", err)
		}

		opts = append(opts, http.WithSmsSender(sender))
	}

	if config.EmailConfig.Enable {
		if _, ok := mail.Senders[config.EmailConfig.Driver]; !ok {
			return fmt.Errorf("can't find sms driver: %s", config.EmailConfig.Driver)
		}

		data, err := yaml.Marshal(config.EmailConfig.Content)
		if err != nil {
			return err
		}
		sender := mail.Senders[config.EmailConfig.Driver]()
		if err := sender.Init(data); err != nil {
			return fmt.Errorf("init mail sender failed: %s", err)
		}

		opts = append(opts, http.WithEmailSender(sender))
	}

	// upload config
	if d, ok := fs.Uploaders[config.UploadConfig.Driver]; ok {
		data, err := yaml.Marshal(config.UploadConfig.Content)
		if err != nil {
			return err
		}
		uploader := d()
		if err = uploader.Init(data); err != nil {
			return errors.Wrap(err, "init uploader failed")
		}
		opts = append(opts, http.WithUploader(uploader))
	}

	s := http.NewServer(config, opts...)
	return s.Run()
}
