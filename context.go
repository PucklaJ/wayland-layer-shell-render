package wayland_layer_shell_render

import (
	"errors"

	zwlr "github.com/PucklaJ/wayland-layer-shell-render/wayland/unstable/wlr-layer-shell-v1"
	zxdg "github.com/rajveermalviya/go-wayland/wayland/unstable/xdg-output-v1"

	"github.com/rajveermalviya/go-wayland/wayland/client"
)

type Context struct {
	Display       *client.Display
	Registry      *client.Registry
	Shm           *client.Shm
	Compositor    *client.Compositor
	Seat          *client.Seat
	Pointer       *client.Pointer
	LayerShell    *zwlr.LayerShell
	OutputManager *zxdg.OutputManager

	Running bool
}

func NewContext() (ctx Context, err error) {
	ctx.Display, err = client.Connect("")
	if err != nil {
		return
	}

	ctx.Registry, err = ctx.Display.GetRegistry()
	if err != nil {
		return
	}

	ctx.Registry.SetGlobalHandler(func(e client.RegistryGlobalEvent) {
		RegistryGlobal(&ctx, e)
	})
	ctx.Registry.SetGlobalRemoveHandler(func(e client.RegistryGlobalRemoveEvent) {
		RegistryGlobalRemove(&ctx, e)
	})
	displayRoundTrip(ctx.Display)

	// Check for sufficient support
	if ctx.Shm == nil {
		err = errors.New("No shm support")
		return
	}
	if ctx.Compositor == nil {
		err = errors.New("No compositor support")
		return
	}
	if ctx.Seat == nil {
		err = errors.New("No seat support")
		return
	}
	if ctx.LayerShell == nil {
		err = errors.New("No layer shell support")
		return
	}
	if ctx.OutputManager == nil {
		err = errors.New("No xdg output manager support")
		return
	}

	ctx.Seat.SetCapabilitiesHandler(func(e client.SeatCapabilitiesEvent) {
		SeatCapabilities(&ctx, e)
	})
	displayRoundTrip(ctx.Display)

	if ctx.Pointer == nil {
		err = errors.New("No pointer support")
		return
	}

	ctx.Pointer.SetEnterHandler(func(e client.PointerEnterEvent) {
		PointerEnter(&ctx, e)
	})
	ctx.Pointer.SetLeaveHandler(func(e client.PointerLeaveEvent) {
		PointerLeave(&ctx, e)
	})
	ctx.Pointer.SetMotionHandler(func(e client.PointerMotionEvent) {
		PointerMotion(&ctx, e)
	})
	ctx.Pointer.SetButtonHandler(func(e client.PointerButtonEvent) {
		PointerButton(&ctx, e)
	})
	ctx.Pointer.SetAxisHandler(func(e client.PointerAxisEvent) {
		PointerAxis(&ctx, e)
	})

	return
}

func displayRoundTrip(d *client.Display) {
	callback, err := d.Sync()
	if err != nil {
		return
	}
	defer callback.Destroy()

	done := false
	callback.SetDoneHandler(func(_ client.CallbackDoneEvent) {
		done = true
	})

	for !done {
		d.Context().Dispatch()
	}
}
