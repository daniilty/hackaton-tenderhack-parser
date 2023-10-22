package main

import "flag"

type flags struct {
	configPath string
}

func initFlags() *flags {
	f := &flags{}

	flag.StringVar(&f.configPath, "config", "./config.json", "путь до конфига")
	flag.Parse()

	return f
}
