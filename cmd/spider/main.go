package main

import (
	"github.com/lyu302/stock/cmd"
	"github.com/lyu302/stock/cmd/spider/option"
	"github.com/lyu302/stock/cmd/spider/server"
	"github.com/spf13/pflag"
	"log"
	"os"
)

func main()  {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		cmd.DisplayVersion("spider")
	}

	c := option.NewSpider()
	c.AddFlags(pflag.CommandLine)
	pflag.Parse()

	if err := server.Run(c); err != nil {
		log.Printf("Spider Start Error: %s", err)
		os.Exit(1)
	}
}
