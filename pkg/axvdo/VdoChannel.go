package axvdo

/*
#cgo pkg-config: glib-2.0 gio-2.0 gio-unix-2.0 vdostream
#include <vdo-stream.h>
#include <vdo-channel.h>
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/Cacsjep/goxis/pkg/clib"
	"github.com/Cacsjep/goxis/pkg/glib"
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-channel_8h.html
type VdoChannel struct {
	Ptr *C.VdoChannel
}

// Get an existing channel.
// Create a new VdoChannel object representing an existing channel session with the specified channel number.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-channel_8h.html#ae5223c6a11ae8f3583457c4917a0b820
func VdoChannelGet(channel_nbr uint) (*VdoChannel, error) {
	cError := clib.NewError()
	ch := C.vdo_channel_get(C.uint(channel_nbr), (**C.GError)(unsafe.Pointer(cError.Ptr)))
	if err := cError.IsError(); err != nil {
		return nil, err
	}
	return &VdoChannel{Ptr: ch}, nil
}

// Gets all existing channels.
// Create new VdoChannel objects representing each existing channel and return those in a GList.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-channel_8h.html#ab1032730452f1617fd004f9aa46ded70
func VdoChannelGetAll() ([]*VdoStream, error) {
	var streams []*VdoStream
	cError := clib.NewError()
	vdoStreamsPtr := uintptr(
		unsafe.Pointer(
			C.vdo_channel_get_all(
				(**C.GError)(unsafe.Pointer(cError.Ptr)),
			),
		),
	)

	if err := cError.IsError(); err != nil {
		return streams, err
	}

	vdoStreamsList := glib.WrapList(vdoStreamsPtr)
	vdoStreamsList.DataWrapper(wrapVdoStream)
	vdoStreamsList.Foreach(func(item interface{}) {
		vdoStream, ok := item.(*VdoStream)
		if !ok {
			panic("VdoChannelGetAll: item is not of type *VdoStream")
		}
		streams = append(streams, vdoStream)
	})
	return streams, nil
}

// Gets all existing channels matching a filter.
// Create new VdoChannel objects representing existing channels matching a filter and return those in a VdoStream List.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-channel_8h.html#a045ce99305ee3e07ce2a243a6fd28861
func VdoChannelGetFilterd(filter *VdoMap) ([]*VdoStream, error) {
	var streams []*VdoStream
	cError := clib.NewError()
	vdoStreamsPtr := uintptr(
		unsafe.Pointer(
			C.vdo_channel_get_filtered(
				filter.Ptr,
				(**C.GError)(unsafe.Pointer(cError.Ptr)),
			),
		),
	)

	if err := cError.IsError(); err != nil {
		return streams, err
	}

	vdoStreamsList := glib.WrapList(vdoStreamsPtr)
	vdoStreamsList.DataWrapper(wrapVdoStream)
	vdoStreamsList.Foreach(func(item interface{}) {
		vdoStream, ok := item.(*VdoStream)
		if !ok {
			panic("VdoChannelGetAll: item is not of type *VdoStream")
		}
		streams = append(streams, vdoStream)
	})
	return streams, nil
}

// Get the ID for this channel.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-channel_8h.html#a2ba1def293a637c15586eb871a586596
func (c *VdoChannel) GetId() (channel_id uint) {
	cid := C.vdo_channel_get_id(c.Ptr)
	return uint(cid)
}

// Get the info for this channel.
// This function is called in order to get a pointer to the info map for this channel.
// The returned pointer is a pointer to a newly allocated VdoMap owned by the caller of this function
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-channel_8h.html#a19ca10165dba07f8295cd0933bfcaa49
func (c *VdoChannel) GetInfo() (*VdoMap, error) {
	cError := clib.NewError()
	infoMap := C.vdo_channel_get_info(c.Ptr, (**C.GError)(unsafe.Pointer(cError.Ptr)))
	if err := cError.IsError(); err != nil {
		return nil, err
	}
	return NewVdoMapFromC(infoMap), nil
}

// Fetch all valid resolutions for a channel with specified stream format.
// Get a VdoResolutionSet of valid resolutions.
// Specifying filter or NULL for default.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-channel_8h.html#ab95177576e046dd6a42c9f87013089ec
func (c *VdoChannel) GetResolutions(filter *VdoMap) ([]VdoResolution, error) {
	cError := clib.NewError()
	var filterMap *C.VdoMap
	if filter != nil {
		filterMap = filter.Ptr
	} else {
		filterMap = nil
	}
	resoSet := C.vdo_channel_get_resolutions(
		c.Ptr,
		filterMap,
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsError(); err != nil {
		return nil, err
	}

	if resoSet == nil {
		return nil, errors.New("VdoResolutionSet is nil ")
	}

	defer C.g_free(C.gpointer(resoSet))

	count := int(resoSet.count)
	resolutions := make([]VdoResolution, count)
	firstResolutionPtr := uintptr(unsafe.Pointer(resoSet)) + unsafe.Sizeof(resoSet.count)
	for i := 0; i < count; i++ {
		resolutionPtr := (*C.VdoResolution)(unsafe.Pointer(firstResolutionPtr + uintptr(i)*unsafe.Sizeof(C.VdoResolution{})))
		resolutions = append(resolutions, VdoResolution{Width: int(resolutionPtr.width), Height: int(resolutionPtr.height)})
	}
	return resolutions, nil
}

// Get the settings for this channel.
// This function is called in order to get a pointer to the settings map for this channel.
// The returned pointer is a pointer to a newly allocated VdoMap owned by the caller of this function. For example:
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-channel_8h.html#a13a6f7939f63317f8b8798e9f3a75ea0
func (c *VdoChannel) GetSettings() (*VdoMap, error) {
	cError := clib.NewError()
	settingsMap := C.vdo_channel_get_settings(c.Ptr, (**C.GError)(unsafe.Pointer(cError.Ptr)))
	if err := cError.IsError(); err != nil {
		return nil, err
	}
	return NewVdoMapFromC(settingsMap), nil
}

// Update the framerate for the specified channel.
// This function is invoked in order to update the framerate of a channel.
// All streams that are connected to this channel may be changed.
// Only streams with a higher fps set than the channel framerate will be affected.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-channel_8h.html#a44e1103d8690783c53103d326e9da5e0
func (c *VdoChannel) SetFramerate(framerate float32) error {
	cError := clib.NewError()
	success := C.vdo_channel_set_framerate(
		c.Ptr,
		C.gdouble(framerate),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)
	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to set framerate"); err != nil {
		return err
	}
	return nil
}

// Update the settings for this channel.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-channel_8h.html#a3e17d1d5abf72b3c826a70284eb1ae99
func (c *VdoChannel) SetSettings(settings *VdoMap) error {
	cError := clib.NewError()
	success := C.vdo_channel_set_settings(
		c.Ptr,
		settings.Ptr,
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)
	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to set settings"); err != nil {
		return err
	}
	return nil
}

func wrapVdoStream(ptr unsafe.Pointer) interface{} {
	return &VdoStream{
		Ptr: (*C.VdoStream)(ptr),
	}
}
