package nanovgo4

import (
	"reflect"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// #include <stdlib.h>
import "C"

/* THIS FILE IS LARGELY A COPY FROM github.com/goxjs */

// Enum is equivalent to GLenum, and is normally used with one of the
// constants defined in this package.
type Enum uint32

// Attrib identifies the location of a specific attribute variable.
type Attrib struct {
	Value uint
}

// Program identifies a compiled shader program.
type Program struct {
	Value uint32
}

// Shader identifies a GLSL shader.
type Shader struct {
	Value uint32
}

// Uniform identifies the location of a specific uniform variable.
type Uniform struct {
	Value int32
}

// A Texture identifies a GL texture unit.
type Texture struct {
	Value uint32
}

// Buffer identifies a GL buffer object.
type Buffer struct {
	Value uint32
}

func ShaderSource(s Shader, src string) {
	glsource, free := Strs(src + "\x00")
	gl.ShaderSource(s.Value, 1, glsource, nil)
	free()
}

// Strs takes a list of Go strings (with or without null-termination) and
// returns their C counterpart.
//
// The returned free function must be called once you are done using the strings
// in order to free the memory.
//
// If no strings are provided as a parameter this function will panic.
func Strs(strs ...string) (cstrs **uint8, free func()) {
	if len(strs) == 0 {
		panic("Strs: expected at least 1 string")
	}

	// Allocate a contiguous array large enough to hold all the strings' contents.
	n := 0
	for i := range strs {
		n += len(strs[i])
	}
	data := C.malloc(C.size_t(n))

	// Copy all the strings into data.
	dataSlice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(data),
		Len:  n,
		Cap:  n,
	}))
	css := make([]*uint8, len(strs)) // Populated with pointers to each string.
	offset := 0
	for i := range strs {
		copy(dataSlice[offset:offset+len(strs[i])], strs[i][:]) // Copy strs[i] into proper data location.
		css[i] = (*uint8)(unsafe.Pointer(&dataSlice[offset]))   // Set a pointer to it.
		offset += len(strs[i])
	}

	return (**uint8)(&css[0]), func() { C.free(data) }
}

func GetShaderi(s Shader, pname Enum) int {
	var result int32
	gl.GetShaderiv(s.Value, uint32(pname), &result)
	return int(result)
}

func GetProgrami(p uint32, pname Enum) int {
	var result int32
	gl.GetProgramiv(p, uint32(pname), &result)
	return int(result)
}

func (v Attrib) Valid() bool  { return v.Value != 0 }
func (v Program) Valid() bool { return v.Value != 0 }
func (v Shader) Valid() bool  { return v.Value != 0 }
func (v Buffer) Valid() bool  { return v.Value != 0 }
func (v Texture) Valid() bool { return v.Value != 0 }
func (v Uniform) Valid() bool { return v.Value != 0 }

func GetUniformLocation(p Program, name string) Uniform {
	return Uniform{Value: gl.GetUniformLocation(p.Value, gl.Str(name+"\x00"))}
}

func DeleteTexture(v Texture) {
	gl.DeleteTextures(1, &v.Value)
}
func BindTexture(target Enum, t Texture) {
	gl.BindTexture(uint32(target), t.Value)
}

func StencilFunc(fn Enum, ref int, mask uint32) {
	gl.StencilFunc(uint32(fn), int32(ref), mask)
}

func Uniform4fv(dst Uniform, src []float32) {
	gl.Uniform4fv(dst.Value, int32(len(src)/4), &src[0])
}

func CreateBuffer() Buffer {
	var b Buffer
	gl.GenBuffers(1, &b.Value)
	return b
}

func CreateVao() uint32 {
	var b uint32
	gl.GenVertexArrays(1, &b)
	return b
}

func CreateTexture() Texture {
	var t Texture
	gl.GenTextures(1, &t.Value)
	return t
}

// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexImage2D.xhtml
func TexImage2D(target Enum, level int, width, height int, format Enum, ty Enum, data []byte) {
	p := unsafe.Pointer(nil)
	if len(data) > 0 {
		p = gl.Ptr(&data[0])
	} else {
		data = make([]byte, width*height*1*100)
		p = gl.Ptr(&data[0])
	}

	gl.TexImage2D(uint32(target), int32(level), int32(format), int32(width), int32(height), 0, uint32(format), uint32(ty), p)

}

func TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	gl.TexSubImage2D(uint32(target), int32(level), int32(x), int32(y), int32(width), int32(height), uint32(format), uint32(ty), gl.Ptr(&data[0]))
}
func BufferData(target Enum, src []byte, usage Enum) {
	gl.BufferData(uint32(target), int(len(src)), gl.Ptr(&src[0]), uint32(usage))
}

func Uniform2fv(dst Uniform, src []float32) {
	gl.Uniform2fv(dst.Value, int32(len(src)/2), &src[0])
}

func DeleteBuffer(v Buffer) {
	gl.DeleteBuffers(1, &v.Value)
}

func GetShaderInfoLog(s Shader) string {
	var logLength int32
	gl.GetShaderiv(s.Value, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetShaderInfoLog(s.Value, logLength, nil, &logBuffer[0])
	return gl.GoStr(&logBuffer[0])
}
func GetProgramInfoLog(p Program) string {
	var logLength int32
	gl.GetProgramiv(p.Value, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetProgramInfoLog(p.Value, logLength, nil, &logBuffer[0])
	return gl.GoStr(&logBuffer[0])
}
