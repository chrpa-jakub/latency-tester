package main

import (
	"os"
)

func main() {
	webPages := ParseArgs(os.Args)
	StartMeasuring(webPages)
}
