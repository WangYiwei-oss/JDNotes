package main

import (
	"fmt"
	"github.com/WangYiwei-oss/cli"
	"os"
)

func main() {
	var a string
	app := &cli.App{
		Name:  "app名字",
		Usage: "帮助信息",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "a",
				Usage:       "控制我说什么话",
				Destination: &a,
				EnvVars:     []string{"LANG"},
			},
			&cli.BoolFlag{
				Name:        "b",
				Aliases:     []string{"c", "d"},
				Value:       true,
				Usage:       "控制我高兴不高兴",
				DefaultText: "高兴",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("app启动就执行了这个")
			if c.Bool("b") == true {
				fmt.Println("我高兴地说：", a)
			} else {
				fmt.Println("我不高兴地说：", a)
			}
			return nil
		},
	}
	app.Run(os.Args)
}
