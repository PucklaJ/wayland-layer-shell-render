package wayland_layer_shell_render

import (
	"fmt"
	"os"

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
		var ro RenderOutput

		ro.Output = client.NewOutput(ctx.Display.Context())
		ctx.Registry.Bind(e.Name, e.Interface, 3, ro.Output)

		ctx.Outputs = append(ctx.Outputs, ro)
	case "zxdg_output_manager_v1":
		ctx.OutputManager = zxdg.NewOutputManager(ctx.Display.Context())
		ctx.Registry.Bind(e.Name, e.Interface, 2, ctx.OutputManager)
	}
}

func RegistryGlobalRemove(ctx *Context, e client.RegistryGlobalRemoveEvent) {
	fmt.Printf("RegistryGlobalRemove: name=%d\n", e.Name)
}

func SeatCapabilities(ctx *Context, e client.SeatCapabilitiesEvent) {
	fmt.Printf("SeatCapabilities: caps=%d\n", e.Capabilities)

	if e.Capabilities&uint32(client.SeatCapabilityPointer) != 0 {
		var err error
		ctx.Pointer, err = ctx.Seat.GetPointer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get pointer: %s\n", err)
		}
	} else {
		ctx.Pointer = nil
	}
}

func PointerEnter(ctx *Context, e client.PointerEnterEvent) {

}

func PointerLeave(ctx *Context, e client.PointerLeaveEvent) {

}

func PointerMotion(ctx *Context, e client.PointerMotionEvent) {

}

func PointerButton(ctx *Context, e client.PointerButtonEvent) {
	ctx.Running = false
}

func PointerAxis(ctx *Context, e client.PointerAxisEvent) {

}

func XdgOutputLogicalPosition(ro *RenderOutput, e zxdg.OutputLogicalPositionEvent) {
	ro.OutputDims[0] = e.X
	ro.OutputDims[1] = e.Y
}

func XdgOutputLogicalSize(ro *RenderOutput, e zxdg.OutputLogicalSizeEvent) {
	ro.OutputDims[2] = e.Width
	ro.OutputDims[3] = e.Height
}

func XdgOutputDone(ro *RenderOutput, e zxdg.OutputDoneEvent) {

}

func XdgOutputName(ro *RenderOutput, e zxdg.OutputNameEvent) {
	ro.Name = e.Name
}

func XdgOutputDescription(ro *RenderOutput, e zxdg.OutputDescriptionEvent) {

}

func LayerSurfaceConfigure(ro *RenderOutput, e zwlr.LayerSurfaceConfigureEvent) {
	ro.ConfigureEvent = new(ConfigureEvent)
	ro.ConfigureEvent.Width = e.Width
	ro.ConfigureEvent.Height = e.Height
}

func LayerSurfaceClosed(ro *RenderOutput, e zwlr.LayerSurfaceClosedEvent) {

}
