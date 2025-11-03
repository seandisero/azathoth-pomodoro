package azathoth

import "fmt"

type AzColor string

const ESC AzColor = "\033["

const NC AzColor = ESC + "0m"
const BG_NC AzColor = ESC + "49m"

const (
	T_RED   AzColor = ESC + "0;31m"
	T_GREEN AzColor = ESC + "0;32m"
)

const (
	BG_RED   AzColor = ESC + "41m"
	BG_GREEN AzColor = ESC + "42m"
)

func PrintWithColor(text string, color AzColor, bg AzColor) {
	fmt.Print(color)
	fmt.Print(bg)
	fmt.Print(text)
	fmt.Print(NC)
}

func DisableCursor() {
	fmt.Print("\033[?25l")
}

func EnableCursor() {
	fmt.Print("\033[?25h")
}
