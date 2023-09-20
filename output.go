package wayland_layer_shell_render

import (
	zwlr "github.com/PucklaJ/wayland-layer-shell-render/wayland/unstable/wlr-layer-shell-v1"
	"github.com/rajveermalviya/go-wayland/wayland/client"
	zxdg "github.com/rajveermalviya/go-wayland/wayland/unstable/xdg-output-v1"
)

type ConfigureEvent struct {
	Width  uint32
	Height uint32
}

type RenderOutput struct {
	Output       *client.Output
	XdgOutput    *zxdg.Output
	Name         string
	Surface      *client.Surface
	LayerSurface *zwlr.LayerSurface
	Memory       *SharedMemory

	LastSize   [2]int32 // w, h
	OutputDims [4]int32 // x, y, w, h

	ConfigureEvent *ConfigureEvent
}
