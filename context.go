package wayland_layer_shell_render

import (
	"errors"
	"fmt"
	"os"
	"time"

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

	Outputs []RenderOutput

	Running bool
}

func NewContext() (ctx *Context, err error) {
	ctx = new(Context)

	ctx.Display, err = client.Connect("")
	if err != nil {
		return
	}

	ctx.Registry, err = ctx.Display.GetRegistry()
	if err != nil {
		return
	}

	ctx.Registry.SetGlobalHandler(func(e client.RegistryGlobalEvent) {
		RegistryGlobal(ctx, e)
	})
	ctx.Registry.SetGlobalRemoveHandler(func(e client.RegistryGlobalRemoveEvent) {
		RegistryGlobalRemove(ctx, e)
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
	if len(ctx.Outputs) == 0 {
		err = errors.New("No outputs")
		return
	}

	ctx.Seat.SetCapabilitiesHandler(func(e client.SeatCapabilitiesEvent) {
		SeatCapabilities(ctx, e)
	})
	displayRoundTrip(ctx.Display)

	if ctx.Pointer == nil {
		err = errors.New("No pointer support")
		return
	}

	ctx.Pointer.SetEnterHandler(func(e client.PointerEnterEvent) {
		PointerEnter(ctx, e)
	})
	ctx.Pointer.SetLeaveHandler(func(e client.PointerLeaveEvent) {
		PointerLeave(ctx, e)
	})
	ctx.Pointer.SetMotionHandler(func(e client.PointerMotionEvent) {
		PointerMotion(ctx, e)
	})
	ctx.Pointer.SetButtonHandler(func(e client.PointerButtonEvent) {
		PointerButton(ctx, e)
	})
	ctx.Pointer.SetAxisHandler(func(e client.PointerAxisEvent) {
		PointerAxis(ctx, e)
	})

	for i := range ctx.Outputs {
		ro := &ctx.Outputs[i]
		var oErr error
		ro.XdgOutput, oErr = ctx.OutputManager.GetXdgOutput(ro.Output)
		if oErr != nil {
			err = fmt.Errorf("Failed to get XDG output: %s", oErr)
			return
		}

		ro.XdgOutput.SetLogicalPositionHandler(func(e zxdg.OutputLogicalPositionEvent) {
			XdgOutputLogicalPosition(ro, e)
		})
		ro.XdgOutput.SetLogicalSizeHandler(func(e zxdg.OutputLogicalSizeEvent) {
			XdgOutputLogicalSize(ro, e)
		})
		ro.XdgOutput.SetDoneHandler(func(e zxdg.OutputDoneEvent) {
			XdgOutputDone(ro, e)
		})
		ro.XdgOutput.SetNameHandler(func(e zxdg.OutputNameEvent) {
			XdgOutputName(ro, e)
		})
		ro.XdgOutput.SetDescriptionHandler(func(e zxdg.OutputDescriptionEvent) {
			XdgOutputDescription(ro, e)
		})

		ro.Surface, oErr = ctx.Compositor.CreateSurface()
		if oErr != nil {
			err = fmt.Errorf("Failed to create surface: %s", oErr)
			return
		}

		ro.LayerSurface, oErr = ctx.LayerShell.GetLayerSurface(ro.Surface, ro.Output, uint32(zwlr.LayerShellLayerOverlay), fmt.Sprint("wayland-layer-shell-render-", i))
		if oErr != nil {
			err = fmt.Errorf("Failed to create layer surface: %s", oErr)
			return
		}

		ro.LayerSurface.SetConfigureHandler(func(e zwlr.LayerSurfaceConfigureEvent) {
			LayerSurfaceConfigure(ro, e)
		})
		ro.LayerSurface.SetClosedHandler(func(e zwlr.LayerSurfaceClosedEvent) {
			LayerSurfaceClosed(ro, e)
		})

		ro.LayerSurface.SetAnchor(uint32(zwlr.LayerSurfaceAnchorTop | zwlr.LayerSurfaceAnchorLeft | zwlr.LayerSurfaceAnchorRight | zwlr.LayerSurfaceAnchorBottom))
		ro.LayerSurface.SetKeyboardInteractivity(0)
		ro.LayerSurface.SetExclusiveZone(-1)
		ro.Surface.Commit()
	}
	displayRoundTrip(ctx.Display)

	return
}

func (ctx *Context) Run() {
	ctx.Running = true

	for ctx.Running {
		displayRoundTrip(ctx.Display)

		for i := range ctx.Outputs {
			ro := &ctx.Outputs[i]

			if ro.ConfigureEvent != nil {
				if ro.Memory != nil {
					ro.Memory.Destroy()
				}

				var err error
				ro.Memory, err = CreateSharedMemory(ctx.Shm, int32(ro.ConfigureEvent.Width), int32(ro.ConfigureEvent.Height), uint32(client.ShmFormatArgb8888))
				if err != nil {
					ro.Memory = nil
					fmt.Fprintf(os.Stderr, "Failed to create memory for \"%s\": %s\n", ro.Name, err)
				}

				ro.ConfigureEvent = nil
			}
		}

		for i := range ctx.Outputs {
			ro := &ctx.Outputs[i]

			if ro.Memory != nil {
				for j := 0; j < len(ro.Memory.Data); j += 4 {
					ro.Memory.Data[j+0] = 100
					ro.Memory.Data[j+1] = 20
					ro.Memory.Data[j+2] = 60
					ro.Memory.Data[j+3] = 10
				}

				ro.Surface.Attach(ro.Memory.Buffer, 0, 0)
				ro.Surface.Damage(0, 0, ro.OutputDims[2], ro.OutputDims[3])
				ro.Surface.Commit()
			}
		}

		time.Sleep(16 * time.Millisecond)
	}
}

func (ctx *Context) Destroy() {
	for i := range ctx.Outputs {
		ro := &ctx.Outputs[i]

		ro.LayerSurface.Destroy()
		ro.Surface.Destroy()
		if ro.Memory != nil {
			ro.Memory.Destroy()
		}
		ro.XdgOutput.Destroy()
	}

	ctx.LayerShell.Destroy()
	ctx.OutputManager.Destroy()
	ctx.Shm.Destroy()
	ctx.Registry.Destroy()
	ctx.Display.Destroy()
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
