package main

import (
	"github.com/reechou/robot-zbshipin/config"
	"github.com/reechou/robot-zbshipin/controller"
)

func main() {
	controller.NewLogic(config.NewConfig()).Run()
}
