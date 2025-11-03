package main

import (
	"azathoth"
)

func main() {
	azathoth := azathoth.NewAzathoth(azathoth.WithDefaultConfig())
	azathoth.Start()
}
