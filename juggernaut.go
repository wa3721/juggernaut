package main

import (
	"github.com/sirupsen/logrus"
	"judgement/app"
)

func main() {
	err := app.NewApp("juggernaut").Run(":9090")
	if err != nil {
		logrus.Fatalf("start http server failed! error: %v", err)
	}
}
