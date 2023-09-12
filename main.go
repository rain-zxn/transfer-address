package main

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/rain-zxn/transfer-address/manager"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "transfer-address",
		Usage: "transfer-address",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "config.json",
				Usage: "configuration file",
			},
		},
		Before: Init,
		Commands: []*cli.Command{
			&cli.Command{
				Name:   manager.BATCH_CREATE_ACCOUNT,
				Usage:  "Batch Create some new eth keystore accounts",
				Action: command(manager.BATCH_CREATE_ACCOUNT),
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:     "count",
						Usage:    "accounts count",
						Required: true,
					},
				},
			},
			&cli.Command{
				Name:   manager.BATCH_TRANSFER_TOKEN,
				Usage:  "Batch transfer token to child accounts",
				Action: command(manager.BATCH_TRANSFER_TOKEN),
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:     "starttosubaccount",
						Usage:    "starttosubaccount",
						Required: true,
					},
					&cli.Int64Flag{
						Name:     "endtosubaccount",
						Usage:    "endtosubaccount",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "token",
						Usage:    "usdt or matic",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "tokenamount",
						Usage:    "tokenamount",
						Required: true,
					},
					&cli.BoolFlag{
						Name:     "estimate",
						Usage:    "need estimate",
						Required: false,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error("Start error", "err", err)
	}
}

func command(method string) func(*cli.Context) error {
	return func(c *cli.Context) error {
		err := manager.HandleCommand(method, c)
		if err != nil {
			log.Error("Failure", "command", method, "err", err)
		} else {
			log.Info("Command was executed successful!")
		}
		return nil
	}
}

func Init(ctx *cli.Context) (err error) {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	return nil
}
