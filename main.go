package main

import (
	"fmt"
	"os"
)

func main() {
	if err := ghreposCmd.Execute(); err != nil {
		abort(err)
	}
}

func abort(err error) {
	fmt.Println(err)
	os.Exit(-1)
}
