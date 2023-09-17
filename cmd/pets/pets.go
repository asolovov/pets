package main

import (
	"pets/internal"
)

// main is a main app endpoint
func main() {
	app := internal.NewApp()

	app.Run()
	app.Stop()
}
