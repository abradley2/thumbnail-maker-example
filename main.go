package main

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"syscall/js"

	"golang.org/x/image/draw"
)

type OutputWriter struct {
	B   []byte
	Len *int
}

func (o OutputWriter) Write(p []byte) (n int, err error) {
	len_copied := copy(o.B[*o.Len:], p)
	*o.Len += len_copied
	return len(p), nil
}

func resizeMe(src_bytes_value js.Value, output_bytes_value js.Value) int {
	src_bytes := make([]byte, src_bytes_value.Length())
	js.CopyBytesToGo(src_bytes, src_bytes_value)

	output_bytes := make([]byte, output_bytes_value.Length())

	src_reader := bytes.NewReader(src_bytes)
	src_image, _ := png.Decode(src_reader)

	dst := image.NewRGBA(image.Rect(0, 0, 64, 64))

	draw.NearestNeighbor.Scale(dst, dst.Rect, src_image, src_image.Bounds(), draw.Over, nil)

	var output io.Writer
	var out_len int

	output = OutputWriter{
		B:   output_bytes,
		Len: &out_len,
	}

	png.Encode(output, dst)

	js.CopyBytesToJS(output_bytes_value, output_bytes)

	return out_len
}

func registerHandler() {
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return resizeMe(args[0], args[1])
	})
	js.Global().Set("makeThumbnail", cb)
}

func main() {
	registerHandler()
	select {}
}
