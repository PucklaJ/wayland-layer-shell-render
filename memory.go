package wayland_layer_shell_render

import (
	"errors"
	"os"

	"github.com/rajveermalviya/go-wayland/wayland/client"
	"golang.org/x/sys/unix"
)

type SharedMemory struct {
	File   *os.File
	Data   []byte
	Buffer *client.Buffer
}

func CreateSharedMemory(shm *client.Shm, width, height int32, format uint32) (mem *SharedMemory, err error) {
	stride := width * 4
	size := stride * height

	dir := os.Getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		err = errors.New("XDG_RUNTIME_DIR is not defined")
		return
	}

	mem = new(SharedMemory)

	mem.File, err = os.CreateTemp(dir, "wlsr_*")
	if err != nil {
		return
	}

	err = mem.File.Truncate(int64(size))
	if err != nil {
		return
	}
	os.Remove(mem.File.Name())

	mem.Data, err = unix.Mmap(int(mem.File.Fd()), 0, int(size), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		return
	}

	pool, err := shm.CreatePool(int(mem.File.Fd()), size)
	if err != nil {
		mem.File.Close()
		unix.Munmap(mem.Data)
		return
	}

	mem.Buffer, err = pool.CreateBuffer(0, width, height, stride, format)
	if err != nil {
		mem.File.Close()
		unix.Munmap(mem.Data)
		return
	}
	pool.Destroy()

	return
}

func (mem *SharedMemory) Destroy() {
	mem.File.Close()
	unix.Munmap(mem.Data)
	mem.Buffer.Destroy()
}
