package main

import (
	"image-storage/app"
)

var apphandler app.App

func main() {

	var err error

	apphandler = app.NewApp()
	apphandler.Init()
	if err = apphandler.Migrate(); err != nil {
		panic(err.Error())
	}

	apphandler.Listen()

}
