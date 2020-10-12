package main

import (
	"github.com/lroman242/redirective/config"
	"github.com/lroman242/redirective/registry"
)

func main() {
	conf := config.ParseConsole()

	r := registry.NewRegistry(conf)
}
