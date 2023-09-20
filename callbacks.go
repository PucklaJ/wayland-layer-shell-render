package wayland_layer_shell_render

import (
	"fmt"

	zwlr "github.com/PucklaJ/wayland-layer-shell-render/wayland/unstable/wlr-layer-shell-v1"
	"github.com/rajveermalviya/go-wayland/wayland/client"
	zxdg "github.com/rajveermalviya/go-wayland/wayland/unstable/xdg-output-v1"
)

func RegistryGlobal(ctx *Context, e client.RegistryGlobalEvent) {
	fmt.Printf("RegistryGlobal: interface=\"%s\" name=%d version=%d\n", e.Interface, e.Name, e.Version)

	switch e.Interface {
	case "wl_shm":
		ctx.Shm = client.NewShm(ctx.Display.Context())
		ctx.Registry.Bind(e.Name, e.Interface, 1, ctx.Shm)
	case "wl_seat":
		ctx.Seat = client.NewSeat(ctx.Display.Context())
		ctx.Registry.Bind(e.Name, e.Interface, 1, ctx.Seat)
	case "wl_compositor":
		ctx.Compositor = client.NewCompositor(ctx.Display.Context())
		ctx.Registry.Bind(e.Name, e.Interface, 1, ctx.Compositor)
	case "zwlr_layer_shell_v1":
		ctx.LayerShell = zwlr.NewLayerShell(ctx.Display.Context())
		ctx.Registry.Bind(e.Name, e.Interface, 1, ctx.LayerShell)
	case "wl_output":
		fmt.Println("TODO: Add outputs")
	case "zxdg_output_manager_v1":
		ctx.OutputManager = zxdg.NewOutputManager(ctx.Display.Context())
		ctx.Registry.Bind(e.Name, e.Interface, 2, ctx.OutputManager)
	}
}

func RegistryGlobalRemove(ctx *Context, e client.RegistryGlobalRemoveEvent) {
	fmt.Printf("RegistryGlobalRemove: name=%d\n", e.Name)
}
