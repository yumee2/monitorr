package main

import (
	"fmt"
	"monitorr/internal/config"
)

const cfgPath = "./config.yaml"

func main() {
	cfg, err := config.ReadConfig(cfgPath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", cfg)
}
