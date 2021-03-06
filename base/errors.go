package base

// #cgo CFLAGS: -I../ex/inc
// #include "fmod_errors.h"
import "C"

import (
  "fmt"
  "errors"
  "reflect"
)

// Takes an interface{} rather than C.FMOD_RESULT because ex and event both
// will have their own definition of C.FMOD_RESULT as far as go is concerned,
// and we don't want to have to typecase all of those when we can just do a
// little work here instead.
func ResultToError(fmod_result interface{}) error {
  value := reflect.ValueOf(fmod_result)
  converter := converters[value.Kind()]
  if converter == nil {
    return errors.New(fmt.Sprintf("Unexpected type: %v", value.Type()))
  }
  err, ok := error_map[C.FMOD_RESULT(converter(value))]
  if !ok {
    return errors.New(fmt.Sprintf("Unknown FMOD_RESULT: %d", fmod_result))
  }
  return err
}

var error_map map[C.FMOD_RESULT]error
var converters map[reflect.Kind]func(reflect.Value) uint64

func init() {
  error_map = map[C.FMOD_RESULT]error{
    C.FMOD_OK:                         nil,
    C.FMOD_ERR_ALREADYLOCKED:          errors.New("Tried to call lock a second time before unlock was called."),
    C.FMOD_ERR_BADCOMMAND:             errors.New("Tried to call a function on a data type that does not allow this type of functionality (ie calling Sound::lock on a streaming sound)."),
    C.FMOD_ERR_CDDA_DRIVERS:           errors.New("Neither NTSCSI nor ASPI could be initialised."),
    C.FMOD_ERR_CDDA_INIT:              errors.New("An error occurred while initialising the CDDA subsystem."),
    C.FMOD_ERR_CDDA_INVALID_DEVICE:    errors.New("Couldn't find the specified device."),
    C.FMOD_ERR_CDDA_NOAUDIO:           errors.New("No audio tracks on the specified disc."),
    C.FMOD_ERR_CDDA_NODEVICES:         errors.New("No CD/DVD devices were found. "),
    C.FMOD_ERR_CDDA_NODISC:            errors.New("No disc present in the specified drive."),
    C.FMOD_ERR_CDDA_READ:              errors.New("A CDDA read error occurred."),
    C.FMOD_ERR_CHANNEL_ALLOC:          errors.New("Error trying to allocate a channel."),
    C.FMOD_ERR_CHANNEL_STOLEN:         errors.New("The specified channel has been reused to play another sound."),
    C.FMOD_ERR_COM:                    errors.New("A Win32 COM related error occured. COM failed to initialize or a QueryInterface failed meaning a Windows codec or driver was not installed properly."),
    C.FMOD_ERR_DMA:                    errors.New("DMA Failure.  See debug output for more information."),
    C.FMOD_ERR_DSP_CONNECTION:         errors.New("DSP connection error.  Connection possibly caused a cyclic dependancy.  Or tried to connect a tree too many units deep (more than 128)."),
    C.FMOD_ERR_DSP_FORMAT:             errors.New("DSP Format error.  A DSP unit may have attempted to connect to this network with the wrong format."),
    C.FMOD_ERR_DSP_NOTFOUND:           errors.New("DSP connection error.  Couldn't find the DSP unit specified."),
    C.FMOD_ERR_DSP_RUNNING:            errors.New("DSP error.  Cannot perform this operation while the network is in the middle of running.  This will most likely happen if a connection or disconnection is attempted in a DSP callback."),
    C.FMOD_ERR_DSP_TOOMANYCONNECTIONS: errors.New("DSP connection error.  The unit being connected to or disconnected should only have 1 input or output."),
    C.FMOD_ERR_FILE_BAD:               errors.New("Error loading file."),
    C.FMOD_ERR_FILE_COULDNOTSEEK:      errors.New("Couldn't perform seek operation.  This is a limitation of the medium (ie netstreams) or the file format."),
    C.FMOD_ERR_FILE_DISKEJECTED:       errors.New("Media was ejected while reading."),
    C.FMOD_ERR_FILE_EOF:               errors.New("End of file unexpectedly reached while trying to read essential data (truncated data?)."),
    C.FMOD_ERR_FILE_NOTFOUND:          errors.New("File not found."),
    C.FMOD_ERR_FILE_UNWANTED:          errors.New("Unwanted file access occured."),
    C.FMOD_ERR_FORMAT:                 errors.New("Unsupported file or audio format."),
    C.FMOD_ERR_HTTP:                   errors.New("A HTTP error occurred. This is a catch-all for HTTP errors not listed elsewhere."),
    C.FMOD_ERR_HTTP_ACCESS:            errors.New("The specified resource requires authentication or is forbidden."),
    C.FMOD_ERR_HTTP_PROXY_AUTH:        errors.New("Proxy authentication is required to access the specified resource."),
    C.FMOD_ERR_HTTP_SERVER_ERROR:      errors.New("A HTTP server error occurred."),
    C.FMOD_ERR_HTTP_TIMEOUT:           errors.New("The HTTP request timed out."),
    C.FMOD_ERR_INITIALIZATION:         errors.New("FMOD was not initialized correctly to support this function."),
    C.FMOD_ERR_INITIALIZED:            errors.New("Cannot call this command after System::init."),
    C.FMOD_ERR_INTERNAL:               errors.New("An error occured that wasn't supposed to.  Contact support."),
    C.FMOD_ERR_INVALID_ADDRESS:        errors.New("On Xbox 360, this memory address passed to FMOD must be physical, (ie allocated with XPhysicalAlloc.)"),
    C.FMOD_ERR_INVALID_FLOAT:          errors.New("Value passed in was a NaN, Inf or denormalized float."),
    C.FMOD_ERR_INVALID_HANDLE:         errors.New("An invalid object handle was used."),
    C.FMOD_ERR_INVALID_PARAM:          errors.New("An invalid parameter was passed to this function."),
    C.FMOD_ERR_INVALID_POSITION:       errors.New("An invalid seek position was passed to this function."),
    C.FMOD_ERR_INVALID_SPEAKER:        errors.New("An invalid speaker was passed to this function based on the current speaker mode."),
    C.FMOD_ERR_INVALID_SYNCPOINT:      errors.New("The syncpoint did not come from this sound handle."),
    C.FMOD_ERR_INVALID_VECTOR:         errors.New("The vectors passed in are not unit length, or perpendicular."),
    C.FMOD_ERR_MAXAUDIBLE:             errors.New("Reached maximum audible playback count for this sound's soundgroup."),
    C.FMOD_ERR_MEMORY:                 errors.New("Not enough memory or resources."),
    C.FMOD_ERR_MEMORY_CANTPOINT:       errors.New("Can't use FMOD_OPENMEMORY_POINT on non PCM source data, or non mp3/xma/adpcm data if FMOD_CREATECOMPRESSEDSAMPLE was used."),
    C.FMOD_ERR_MEMORY_SRAM:            errors.New("Not enough memory or resources on console sound ram."),
    C.FMOD_ERR_NEEDS2D:                errors.New("Tried to call a command on a 3d sound when the command was meant for 2d sound."),
    C.FMOD_ERR_NEEDS3D:                errors.New("Tried to call a command on a 2d sound when the command was meant for 3d sound."),
    C.FMOD_ERR_NEEDSHARDWARE:          errors.New("Tried to use a feature that requires hardware support.  (ie trying to play a GCADPCM compressed sound in software on Wii)."),
    C.FMOD_ERR_NEEDSSOFTWARE:          errors.New("Tried to use a feature that requires the software engine.  Software engine has either been turned off, or command was executed on a hardware channel which does not support this feature."),
    C.FMOD_ERR_NET_CONNECT:            errors.New("Couldn't connect to the specified host."),
    C.FMOD_ERR_NET_SOCKET_ERROR:       errors.New("A socket error occurred.  This is a catch-all for socket-related errors not listed elsewhere."),
    C.FMOD_ERR_NET_URL:                errors.New("The specified URL couldn't be resolved."),
    C.FMOD_ERR_NET_WOULD_BLOCK:        errors.New("Operation on a non-blocking socket could not complete immediately."),
    C.FMOD_ERR_NOTREADY:               errors.New("Operation could not be performed because specified sound/DSP connection is not ready."),
    C.FMOD_ERR_OUTPUT_ALLOCATED:       errors.New("Error initializing output device, but more specifically, the output device is already in use and cannot be reused."),
    C.FMOD_ERR_OUTPUT_CREATEBUFFER:    errors.New("Error creating hardware sound buffer."),
    C.FMOD_ERR_OUTPUT_DRIVERCALL:      errors.New("A call to a standard soundcard driver failed, which could possibly mean a bug in the driver or resources were missing or exhausted."),
    C.FMOD_ERR_OUTPUT_ENUMERATION:     errors.New("Error enumerating the available driver list. List may be inconsistent due to a recent device addition or removal."),
    C.FMOD_ERR_OUTPUT_FORMAT:          errors.New("Soundcard does not support the minimum features needed for this soundsystem (16bit stereo output)."),
    C.FMOD_ERR_OUTPUT_INIT:            errors.New("Error initializing output device."),
    C.FMOD_ERR_OUTPUT_NOHARDWARE:      errors.New("FMOD_HARDWARE was specified but the sound card does not have the resources necessary to play it."),
    C.FMOD_ERR_OUTPUT_NOSOFTWARE:      errors.New("Attempted to create a software sound but no software channels were specified in System::init."),
    C.FMOD_ERR_PAN:                    errors.New("Panning only works with mono or stereo sound sources."),
    C.FMOD_ERR_PLUGIN:                 errors.New("An unspecified error has been returned from a 3rd party plugin."),
    C.FMOD_ERR_PLUGIN_INSTANCES:       errors.New("The number of allowed instances of a plugin has been exceeded."),
    C.FMOD_ERR_PLUGIN_MISSING:         errors.New("A requested output, dsp unit type or codec was not available."),
    C.FMOD_ERR_PLUGIN_RESOURCE:        errors.New("A resource that the plugin requires cannot be found. (ie the DLS file for MIDI playback)"),
    C.FMOD_ERR_PRELOADED:              errors.New("The specified sound is still in use by the event system, call EventSystem::unloadFSB before trying to release it."),
    C.FMOD_ERR_PROGRAMMERSOUND:        errors.New("The specified sound is still in use by the event system, wait for the event which is using it finish with it."),
    C.FMOD_ERR_RECORD:                 errors.New("An error occured trying to initialize the recording device."),
    C.FMOD_ERR_REVERB_INSTANCE:        errors.New("Specified instance in FMOD_REVERB_PROPERTIES couldn't be set. Most likely because it is an invalid instance number or the reverb doesnt exist."),
    C.FMOD_ERR_SUBSOUND_ALLOCATED:     errors.New("This subsound is already being used by another sound, you cannot have more than one parent to a sound.  Null out the other parent's entry first."),
    C.FMOD_ERR_SUBSOUND_CANTMOVE:      errors.New("Shared subsounds cannot be replaced or moved from their parent stream, such as when the parent stream is an FSB file."),
    C.FMOD_ERR_SUBSOUND_MODE:          errors.New("The subsound's mode bits do not match with the parent sound's mode bits.  See documentation for function that it was called with."),
    C.FMOD_ERR_SUBSOUNDS:              errors.New("The error occured because the sound referenced contains subsounds when it shouldn't have, or it doesn't contain subsounds when it should have.  The operation may also not be able to be performed on a parent sound, or a parent sound was played without setting up a sentence first."),
    C.FMOD_ERR_TAGNOTFOUND:            errors.New("The specified tag could not be found or there are no tags."),
    C.FMOD_ERR_TOOMANYCHANNELS:        errors.New("The sound created exceeds the allowable input channel count.  This can be increased using the maxinputchannels parameter in System::setSoftwareFormat."),
    C.FMOD_ERR_UNIMPLEMENTED:          errors.New("Something in FMOD hasn't been implemented when it should be! contact support!"),
    C.FMOD_ERR_UNINITIALIZED:          errors.New("This command failed because System::init or System::setDriver was not called."),
    C.FMOD_ERR_UNSUPPORTED:            errors.New("A command issued was not supported by this object.  Possibly a plugin without certain callbacks specified."),
    C.FMOD_ERR_UPDATE:                 errors.New("An error caused by System::update occured."),
    C.FMOD_ERR_VERSION:                errors.New("The version number of this file format is not supported."),
    C.FMOD_ERR_EVENT_FAILED:           errors.New("An Event failed to be retrieved, most likely due to 'just fail' being specified as the max playbacks behavior."),
    C.FMOD_ERR_EVENT_INFOONLY:         errors.New("Can't execute this command on an EVENT_INFOONLY event."),
    C.FMOD_ERR_EVENT_INTERNAL:         errors.New("An error occured that wasn't supposed to.  See debug log for reason."),
    C.FMOD_ERR_EVENT_MAXSTREAMS:       errors.New("Event failed because 'Max streams' was hit when FMOD_EVENT_INIT_FAIL_ON_MAXSTREAMS was specified."),
    C.FMOD_ERR_EVENT_MISMATCH:         errors.New("FSB mismatches the FEV it was compiled with, the stream/sample mode it was meant to be created with was different, or the FEV was built for a different platform."),
    C.FMOD_ERR_EVENT_NAMECONFLICT:     errors.New("A category with the same name already exists."),
    C.FMOD_ERR_EVENT_NOTFOUND:         errors.New("The requested event, event group, event category or event property could not be found."),
    C.FMOD_ERR_EVENT_NEEDSSIMPLE:      errors.New("Tried to call a function on a complex event that's only supported by simple events."),
    C.FMOD_ERR_EVENT_GUIDCONFLICT:     errors.New("An event with the same GUID already exists."),
    C.FMOD_ERR_EVENT_ALREADY_LOADED:   errors.New("The specified project or bank has already been loaded. Having multiple copies of the same project loaded simultaneously is forbidden."),
    C.FMOD_ERR_MUSIC_UNINITIALIZED:    errors.New("Music system is not initialized probably because no music data is loaded."),
    C.FMOD_ERR_MUSIC_NOTFOUND:         errors.New("The requested music entity could not be found."),
    C.FMOD_ERR_MUSIC_NOCALLBACK:       errors.New("The music callback is required, but it has not been set."),
  }
  from_int := func(v reflect.Value) uint64 {
    return uint64(v.Int())
  }
  from_uint := func(v reflect.Value) uint64 {
    return uint64(v.Uint())
  }
  converters = map[reflect.Kind]func(reflect.Value) uint64{
    reflect.Int:     from_int,
    reflect.Int8:    from_int,
    reflect.Int16:   from_int,
    reflect.Int32:   from_int,
    reflect.Int64:   from_int,
    reflect.Uint:    from_uint,
    reflect.Uint8:   from_uint,
    reflect.Uint16:  from_uint,
    reflect.Uint32:  from_uint,
    reflect.Uint64:  from_uint,
    reflect.Uintptr: from_uint,
  }
}
