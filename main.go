package main

import (
	"flag"
	"fmt"
	"gosnmp/service"
	"os"
)

func main() {
	cfgFile := flag.String("config", "", "config file")
	flag.Parse()

	if *cfgFile == "" {
		fmt.Fprintf(os.Stderr, "please specified config file\n")
		os.Exit(1)
	}

	svc, err := service.New(*cfgFile)
	if err != nil {
		panic(err)
	}

	if err := svc.Run(); err != nil {
		panic(err)
	}
}
