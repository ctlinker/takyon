package ui

import (
	"fmt"
	"os"
)

const (
	green  = "\033[32m"
	red    = "\033[31m"
	blue   = "\033[34m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)

func Info(msg string, args ...any) {
	fmt.Printf(blue+"[*] "+msg+reset+"\n", args...)
}

func Success(msg string, args ...any) {
	fmt.Printf(green+"[+] "+msg+reset+"\n", args...)
}

func Error(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, red+"[!] "+msg+reset+"\n", args...)
}

func Fatal(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, red+"[!] "+msg+reset+"\n", args...)
	os.Exit(1)
}

func Step(msg string, args ...any) {
	fmt.Printf("[→] "+msg+"\n", args...)
}

// your warning function
func Warn(msg string, args ...any) {
	fmt.Printf(yellow+"[?] "+msg+reset+"\n", args...)
}

var AbortErr = fmt.Errorf("Aborting operation")
