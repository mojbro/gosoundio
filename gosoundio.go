// Package gosoundio is a Go wrapper for libsoundio, cross-platform library
// for real-time audio input and output.
package gosoundio

/*
#cgo LDFLAGS: -lsoundio
#include "gosoundio.h"
*/
import "C"
import (
	"errors"
	"unsafe"

	pointer "github.com/mattn/go-pointer"
)

// Error represents a libsoundio error.
type Error int

// SoundIo Errors.
const (
	ErrorNoMem               Error = Error(C.SoundIoErrorNoMem)
	ErrorInitAudioBackend    Error = Error(C.SoundIoErrorInitAudioBackend)
	ErrorSystemResources     Error = Error(C.SoundIoErrorSystemResources)
	ErrorOpeningDevice       Error = Error(C.SoundIoErrorOpeningDevice)
	ErrorNoSuchDevice        Error = Error(C.SoundIoErrorNoSuchDevice)
	ErrorInvalid             Error = Error(C.SoundIoErrorInvalid)
	ErrorBackendUnavailable  Error = Error(C.SoundIoErrorBackendUnavailable)
	ErrorStreaming           Error = Error(C.SoundIoErrorStreaming)
	ErrorIncompatibleDevice  Error = Error(C.SoundIoErrorIncompatibleDevice)
	ErrorNoSuchClient        Error = Error(C.SoundIoErrorNoSuchClient)
	ErrorIncompatibleBackend Error = Error(C.SoundIoErrorIncompatibleBackend)
	ErrorBackendDisconnected Error = Error(C.SoundIoErrorBackendDisconnected)
	ErrorInterrupted         Error = Error(C.SoundIoErrorInterrupted)
	ErrorUnderflow           Error = Error(C.SoundIoErrorUnderflow)
	ErrorEncodingString      Error = Error(C.SoundIoErrorEncodingString)
)

func (e Error) Error() string {
	return C.GoString(C.soundio_strerror(C.int(e)))
}

type SoundIO struct {
	context unsafe.Pointer
}

func CreateSoundIO() (*SoundIO, error) {
	context := C.soundio_create()
	if context == nil {
		return nil, ErrorNoMem
	}
	return &SoundIO{context: unsafe.Pointer(C.soundio_create())}, nil
}

func (s SoundIO) cPtr() *C.struct_SoundIo {
	return (*C.struct_SoundIo)(s.context)
}

func (s *SoundIO) Destroy() {
	C.soundio_destroy(s.cPtr())
	s.context = nil
}

func (s SoundIO) Connect() error {
	r := C.gosoundio_connect(s.context)
	if r != 0 {
		return Error(r)
	}
	return nil
}

func (s SoundIO) Disconnect() {
	C.soundio_disconnect(s.cPtr())
}

func (s SoundIO) DefaultOutputDeviceIndex() (index int, err error) {
	i := int(C.soundio_default_output_device_index(s.cPtr()))
	if i < 0 {
		return i, ErrorNoSuchDevice
	}
	return i, nil
}

func (s SoundIO) GetOutputDevice(index int) (*Device, error) {
	ptr := C.soundio_get_output_device(s.cPtr(), C.int(index))
	if ptr == nil {
		return nil, errors.New("error getting output device")
	}
	return &Device{ptr: unsafe.Pointer(ptr)}, nil
}

type Device struct {
	ptr unsafe.Pointer
}

func (d Device) cPtr() *C.struct_SoundIoDevice {
	return (*C.struct_SoundIoDevice)(d.ptr)
}

// Unref releases the device reference. Call this when you are done using the
// device.
func (d Device) Unref() {
	C.soundio_device_unref(d.cPtr())
}

// Id returns the the device ID, a string that uniquely identifies this device.
func (d Device) Id() string {
	return C.GoString((*d.cPtr()).id)
}

// Name returns a user-friendly name that describes this device.
func (d Device) Name() string {
	return C.GoString((*d.cPtr()).name)
}

// IsRaw returns true if you are directly opening the hardware device without
// going through a proxy such as dmix, PulseAudio, or JACK.
//
// When you open a raw device, other applications on the computer are not able
// to simultaneously access the device. Raw devices do not perform automatic
// resampling and thus tend to have fewer formats available.
func (d Device) IsRaw() bool {
	return bool((*d.cPtr()).is_raw)
}

type OutStream struct {
	ptr unsafe.Pointer
}

type WriteCallback struct {
	Func    func(*OutStream, int, int)
	oStream *OutStream
}

//export writeCallbackProxy
func writeCallbackProxy(outstream *C.struct_SoundIoOutStream, frameCountMin C.int, frameCountMax C.int) {
	wcb := pointer.Restore(outstream.userdata).(*WriteCallback)
	wcb.Func(wcb.oStream, int(frameCountMin), int(frameCountMax))
}

func NewOutStream(d Device, f func(*OutStream, int, int)) (*OutStream, error) {
	cOutStream := C.create_outstream(d.cPtr())
	if cOutStream == nil {
		return nil, ErrorNoMem
	}
	outStream := &OutStream{ptr: unsafe.Pointer(cOutStream)}
	cOutStream.userdata = pointer.Save(&WriteCallback{
		Func:    f,
		oStream: outStream,
	})
	return outStream, nil
}

func (o OutStream) cPtr() *C.struct_SoundIoOutStream {
	return (*C.struct_SoundIoOutStream)(o.ptr)
}

func (o *OutStream) Destroy() {
	C.soundio_outstream_destroy(o.cPtr())
}

func (o *OutStream) Open() error {
	r := C.soundio_outstream_open(o.cPtr())
	if r != 0 {
		return Error(r)
	}
	return nil
}

func (o *OutStream) ClearBuffer() error {
	r := C.soundio_outstream_clear_buffer(o.cPtr())
	if r != 0 {
		return Error(r)
	}
	return nil
}

// Temporary function that starts out stream and waits for events forever
func (o *OutStream) Run() {
	C.run(o.cPtr())
}
