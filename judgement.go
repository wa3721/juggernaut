package main

import (
	"github.com/sirupsen/logrus"
	"judgement/app"
)

func main() {
	err := app.NewApp("juggernaut").Run(":8080")
	if err != nil {
		logrus.Fatalf("start http server failed! error: %v", err)
	}
}
