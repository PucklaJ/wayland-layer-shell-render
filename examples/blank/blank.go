package main

import (
	"fmt"
	"os"

	wlsr "github.com/PucklaJ/wayland-layer-shell-render"
)

func main() {
	_, err := wlsr.NewContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create context: %s\n", err)
		os.Exit(1)
	}
}
