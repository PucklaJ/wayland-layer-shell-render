package main

import (
	"fmt"
	"os"

	wlsr "github.com/PucklaJ/wayland-layer-shell-render"
)

func main() {
	ctx, err := wlsr.NewContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create context: %s\n", err)
		os.Exit(1)
	}

	ctx.Run()

	ctx.Destroy()
}
