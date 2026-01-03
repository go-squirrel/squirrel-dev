package main

import (
	"log"
	
	"squirrel-dev/cmd/app/app"
)

func main() {
	cmd := app.NewServerCommand()
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
