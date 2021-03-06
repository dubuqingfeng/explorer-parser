package main

import (
	"flag"
	"fmt"
	"github.com/dubuqingfeng/explorer-parser/consumer/config"
	"github.com/dubuqingfeng/explorer-parser/consumer/filters"
	"github.com/dubuqingfeng/explorer-parser/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
	"os"
)

func init() {
	app := &cli.App{
		Name:  "Consumer",
		Usage: "",
		Action: func(c *cli.Context) error {
			set := flag.NewFlagSet("contrive", 0)
			nc := cli.NewContext(c.App, set, c)
			fmt.Printf("Load config from %#v \n", nc.String("config"))
			config.InitConfig(nc.String("config"))
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "./config/consumer.json",
				Usage:   "Load configuration from `FILE`",
			},
		},
	}
	app.Version = "0.0.1"
	app.Name = "Explorer Parser Consumer"
	app.Run(os.Args)

	utils.InitLog(config.Config.Log.Level, config.Config.Log.Path, config.Config.Log.Filename)
}

func main() {
	log.Info("consumer start")
	fmt.Printf("config: %#v\n", config.Config.APPName)

	multiCoin := make([]filters.Filter, 0)

	for _, value := range config.Config.EnableCoin {
		newFilter := newFilter(value)
		if newFilter != nil {
			multiCoin = append(multiCoin, newFilter)
		}
	}

	for _, value := range multiCoin {
		// go func
		go value.Filter("test")
	}
}

func newFilter(coin string) filters.Filter {
	switch coin {
	case "btc":
		return filters.NewBTCFilter()
	case "xmr":
		return filters.NewBTCFilter()
	}
	return nil
}
