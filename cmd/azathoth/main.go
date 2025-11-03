package main

import "github.com/seandisero/azathoth-pomodoro"

func main() {
	azathoth := azathoth.NewAzathoth(azathoth.WithDefaultConfig())
	azathoth.Start()
}
