package web

import (
	"fmt"
	"github.com/clearcodecn/log"
	"github.com/clearcodecn/wetalk/api/http"
	"github.com/clearcodecn/wetalk/configs"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/zap"
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
			cli.BoolFlag{
				Name:  "reload",
				Usage: "enable config watch reload",
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

	if ctx.Bool("reload") {
		opts = append(opts, http.WithConfigReloader(configReload(ctx.String("c"))))
	}

	// load sms
	if config.SmsConfig.Enable {
		if _, ok := configs.SmsSenders[config.SmsConfig.Driver]; !ok {
			return fmt.Errorf("can't find sms driver: %s", config.SmsConfig.Driver)
		}

		data, err := yaml.Marshal(config.SmsConfig.Content)
		if err != nil {
			return err
		}
		sender := configs.SmsSenders[config.SmsConfig.Driver]()
		if err := sender.Init(data); err != nil {
			return fmt.Errorf("init sms sender failed: %s", err)
		}

		opts = append(opts, http.WithSmsSender(sender))
	}

	if config.EmailConfig.Enable {
		if _, ok := configs.MailSenders[config.EmailConfig.Driver]; !ok {
			return fmt.Errorf("can't find sms driver: %s", config.EmailConfig.Driver)
		}

		data, err := yaml.Marshal(config.EmailConfig.Content)
		if err != nil {
			return err
		}
		sender := configs.MailSenders[config.EmailConfig.Driver]()
		if err := sender.Init(data); err != nil {
			return fmt.Errorf("init mail sender failed: %s", err)
		}

		opts = append(opts, http.WithEmailSender(sender))
	}

	// upload config
	if d, ok := configs.Uploaders[config.UploadConfig.Driver]; ok {
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

	if config.DbConfig.Driver == "" || config.DbConfig.Dsn == "" {
		return fmt.Errorf("dbconfig can not be empty")
	}

	s := http.NewServer(config, opts...)
	return s.Run()
}

func configReload(path string) func() chan configs.WebConfig {
	return func() chan configs.WebConfig {
		ch := make(chan configs.WebConfig)
		w, err := fsnotify.NewWatcher()
		if err != nil {
			panic(err)
		}
		_ = w.Add(path)
		go func() {
			for {
				select {
				case ev := <-w.Events:
					if ev.Op == fsnotify.Write {
						data, err := ioutil.ReadFile(path)
						if err != nil {
							log.Error("failed to read config", zap.Error(err))
							continue
						}
						cfg, err := configs.ParseWebConfig(data)
						if err != nil {
							log.Error("unable to parse config", zap.Error(err))
							continue
						}
						ch <- cfg
					}
				case <-w.Errors:
					w.Close()
					close(ch)
					return
				}
			}
		}()
		return ch
	}
}
