package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/axlicense"
	"github.com/Cacsjep/goxis/pkg/axmdb"
	"github.com/Cacsjep/goxis/pkg/axparameter"
	"github.com/Cacsjep/goxis/pkg/axvdo"
	"github.com/Cacsjep/goxis/pkg/glib"
	"github.com/stretchr/testify/assert"
)

func main() {
	testing.Main(
		nil,
		[]testing.InternalTest{
			//{"EventTests", EventTests},
			//{"TestVdoMapOperations", TestVdoMapOperations},
			//{"VdoMapTest", VdoMapTest},
			//{"VdoChannelTest", VdoChannelTest},
			//{"TestVdoStream", TestVdoStream},
			//{"LicenseTest", LicenseTest},
			//{"ParamTests", ParamTests},
			//{"EventHandlerTests", EventHandlerTests},
			{"MdbTests", MdbTests},
		},
		nil, nil,
	)
}

func TestVdoFrameOperations(t *testing.T, existingBuffer *axvdo.VdoBuffer) {

	data, err := existingBuffer.GetBytes()
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Greater(t, len(data), 0)

	frame, err := existingBuffer.GetFrame()
	assert.NoError(t, err)
	assert.NotNil(t, frame.Ptr, "Frame should not be nil")

	frameType := frame.GetFrameType()
	assert.NotNil(t, frameType, "Frame type should not be nil")

	// Test GetSequenceNbr
	seqNum := frame.GetSequenceNbr()
	assert.Equal(t, seqNum, uint(0), "Sequence number should be positive")

	// Test GetTimestamp
	timestamp := frame.GetTimestamp()
	assert.Greater(t, timestamp, uint64(0), "Timestamp should be positive")

	// Test GetCustomTimestamp
	customTimestamp := frame.GetCustomTimestamp()
	assert.NotZero(t, customTimestamp, "Custom timestamp should not be zero")

	// Test GetSize
	size := frame.GetSize()
	assert.Greater(t, size, uint(0), "Size should be positive")

	// Test GetHeaderSize
	headerSize := frame.GetHeaderSize()
	assert.GreaterOrEqual(t, headerSize, 0, "Header size should be non-negative")

	// Test GetFd
	fd := frame.GetFd()
	assert.Greater(t, fd, -1, "File descriptor should be valid")

	// Test GetExtraInfo
	extraInfo := frame.GetExtraInfo()
	assert.NotNil(t, extraInfo, "Extra info should not be nil")

	// Test GetOpaque
	opaque := frame.GetOpaque()
	assert.Nil(t, opaque, "Opaque should not be nil")

	// Test GetIsLastBuffer
	isLastBuffer := frame.GetIsLastBuffer()
	assert.NotNil(t, isLastBuffer, "IsLastBuffer should not be nil")

	// Assuming SetSize and other setters affect the frame; otherwise, they're not testable without effect
	frame.SetSize(1024)
	newSize := frame.GetSize()
	assert.Equal(t, newSize, uint(1024), "Frame size should be updated to 1024")
}

func TestVdoBufferOperations(t *testing.T, existingBuffer *axvdo.VdoBuffer) {
	// Test GetID
	id := existingBuffer.GetId()
	assert.NotZero(t, id, "Expected non-zero ID for the buffer")

	// Test GetFd
	fd, err := existingBuffer.GetFd()
	assert.NoError(t, err)
	assert.Greater(t, fd, -1, "Expected valid file descriptor")

	// Test GetOffset
	offset := existingBuffer.GetOffset()
	assert.GreaterOrEqual(t, offset, int64(0), "Expected non-negative offset")

	// Test GetCapacity
	capacity := existingBuffer.GetCapacity()
	assert.NotZero(t, capacity, "Expected non-zero capacity")

	// Test IsComplete
	isComplete := existingBuffer.IsComplete()
	// This assert depends on your expectations of the buffer's completeness
	assert.True(t, isComplete, "Expected valid completeness state")

	// Test GetOpaque - Use carefully
	opaque := existingBuffer.GetOpaque()
	assert.Nil(t, opaque, "Expected nil opaque pointer")

	// Test GetData - Use carefully
	data, err := existingBuffer.GetData()
	assert.NoError(t, err)
	assert.NotNil(t, data, "Expected non-nil data pointer")

	TestVdoFrameOperations(t, existingBuffer)
}

func TestVdoStream(t *testing.T) {
	settings := axvdo.NewVdoMap()
	settings.SetUint32("channel", 1)
	settings.SetUint32("format", uint32(axvdo.VdoFormatJPEG))
	//settings.SetUint32("width", 1920)
	//settings.SetUint32("height", 1080)

	assert.NotNil(t, settings.Ptr)

	s, err := axvdo.NewStream(settings)
	assert.NoError(t, err)
	assert.NotNil(t, s.Ptr)
	defer s.Unref()

	id := s.GetId()
	assert.NotEqual(t, -1, id, "Stream ID should be valid")

	err = s.Start()
	assert.Nil(t, err, "Starting stream should not produce an error")

	streamSettings, err := s.GetSettings()
	assert.Nil(t, err, "Failed to get stream settings")
	defer streamSettings.Unref() // Clean up resources

	err = s.SetSettings(streamSettings)
	assert.Nil(t, err, "Failed to update stream settings")

	err = s.ForceKeyFrame()
	assert.Nil(t, err, "Failed to force key frame")

	err = s.SetFramerate(30.0) // Assuming 30.0 is a valid framerate
	assert.Nil(t, err, "Setting framerate should not fail")

	buffer, err := s.BufferAlloc()
	assert.Error(t, err)
	assert.Equal(t, err.(*axvdo.VdoError).Code, axvdo.VdoErrorCodeInvalidArgument)
	defer func() {
		if buffer != nil {
			err := s.BufferUnref(buffer)
			assert.Nil(t, err, "Failed to unref buffer")
		}
	}()

	intentMap := axvdo.NewVdoMap() // Setup intent map as needed
	intentMap.SetUint32("intent", uint32(axvdo.VdoIntentEventFD))
	err = s.Attach(intentMap)
	assert.Nil(t, err, "Failed to attach with intent")
	defer intentMap.Unref()

	// Get the stream file descriptor
	fd, err := s.GetFd()
	assert.Nil(t, err, "Getting stream file descriptor should not fail")
	assert.Greater(t, fd, -1, "File descriptor should be valid")

	// TODO: Fix this !! returns -1
	//eventFd, err := s.GetEventFd()
	//assert.Nil(t, err, "Getting event file descriptor should not fail")
	//assert.Greater(t, eventFd, -1, "Event file descriptor should be valid")

	snapshotBuffer, err := axvdo.Snapshot(settings)
	assert.Nil(t, err, "Taking a snapshot should not fail")
	assert.NotNil(t, snapshotBuffer.Ptr, "snapshotBuffer.Ptr is nil")

	streams, err := axvdo.StreamGetAll()
	assert.NoError(t, err)
	assert.NotNil(t, streams)
	assert.Greater(t, len(streams), 0)

	TestVdoBufferOperations(t, snapshotBuffer)
	snapshotBuffer.Unref()
	s.Stop()
}

func TestVdoMapOperations(t *testing.T) {
	// Assuming NewVdoMap is a constructor that initializes VdoMap correctly.
	vdoMap := axvdo.NewVdoMap()
	anotherMap := axvdo.NewVdoMap()

	// Initially, both maps should be empty
	assert.True(t, vdoMap.Empty(), "vdoMap should be empty initially")
	assert.Equal(t, 0, vdoMap.Size(), "vdoMap should have size 0 initially")

	// Add some entries to vdoMap
	vdoMap.SetString("key1", "value1")
	vdoMap.SetInt32("key2", 1234)

	// Now, vdoMap should not be empty and should contain the keys
	assert.False(t, vdoMap.Empty(), "vdoMap should not be empty after adding entries")
	assert.Equal(t, 2, vdoMap.Size(), "vdoMap should have size 2 after adding entries")
	assert.True(t, vdoMap.Contains("key1"), "vdoMap should contain key1")
	assert.True(t, vdoMap.Contains("key2"), "vdoMap should contain key2")

	// Testing Swap operation
	vdoMap.Swap(anotherMap)
	assert.True(t, vdoMap.Empty(), "vdoMap should be empty after swap")
	assert.False(t, anotherMap.Empty(), "anotherMap should not be empty after swap")

	// Reset maps to original state for further testing
	vdoMap.Swap(anotherMap)

	// Testing Merge operation
	anotherMap.SetString("key3", "value3")
	vdoMap.Merge(anotherMap)
	assert.Equal(t, 3, vdoMap.Size(), "vdoMap should have size 3 after merge")
	assert.True(t, vdoMap.Contains("key3"), "vdoMap should contain key3 after merge")

	// Testing Remove operation
	vdoMap.Remove("key3")
	assert.False(t, vdoMap.Contains("key3"), "vdoMap should not contain key3 after removal")
	assert.Equal(t, 2, vdoMap.Size(), "vdoMap should have size 2 after removing key3")

	// Testing Equal operation
	assert.False(t, vdoMap.Equal(anotherMap), "vdoMap should not be equal to anotherMap after modifications")

	vdoMap.Unref()
	anotherMap.Unref()
}

func VdoMapTest(t *testing.T) {
	vdoMap := axvdo.NewVdoMap()
	byteValue := byte(255)
	boolValue := true
	int16Value := int16(-32768)
	uint16Value := uint16(65535)
	int32Value := int32(-2147483648)
	uint32Value := uint32(4294967295)
	int64Value := int64(-9223372036854775808)
	uint64Value := uint64(18446744073709551615)
	doubleValue := 3.14159
	stringValue := "Test String"

	// Set values using setters
	vdoMap.SetByte("byteKey", byteValue)
	vdoMap.SetBoolean("boolKey", boolValue)
	vdoMap.SetInt16("int16Key", int16Value)
	vdoMap.SetUint16("uint16Key", uint16Value)
	vdoMap.SetInt32("int32Key", int32Value)
	vdoMap.SetUint32("uint32Key", uint32Value)
	vdoMap.SetInt64("int64Key", int64Value)
	vdoMap.SetUint64("uint64Key", uint64Value)
	vdoMap.SetDouble("doubleKey", doubleValue)
	vdoMap.SetString("stringKey", stringValue)

	// Get values using getters and assert they match what was set
	assert.Equal(t, byteValue, vdoMap.GetByte("byteKey", 0), "Byte value did not match")
	assert.Equal(t, boolValue, vdoMap.GetBoolean("boolKey", false), "Boolean value did not match")
	assert.Equal(t, int16Value, vdoMap.GetInt16("int16Key", 0), "Int16 value did not match")
	assert.Equal(t, uint16Value, vdoMap.GetUint16("uint16Key", 0), "Uint16 value did not match")
	assert.Equal(t, int32Value, vdoMap.GetInt32("int32Key", 0), "Int32 value did not match")
	assert.Equal(t, uint32Value, vdoMap.GetUint32("uint32Key", 0), "Uint32 value did not match")
	assert.Equal(t, int64Value, vdoMap.GetInt64("int64Key", 0), "Int64 value did not match")
	assert.Equal(t, uint64Value, vdoMap.GetUint64("uint64Key", 0), "Uint64 value did not match")
	assert.Equal(t, doubleValue, vdoMap.GetDouble("doubleKey", 0), "Double value did not match")
	assert.Equal(t, stringValue, vdoMap.GetString("stringKey", "foo"), "String value did not match")

	vdoMap.Unref()
}

func VdoChannelTest(t *testing.T) {
	VDO_CHANNEL := uint(1)
	s, err := axvdo.VdoChannelGet(VDO_CHANNEL)
	assert.NoError(t, err)
	assert.Equal(t, VDO_CHANNEL, s.GetId())
	resos, err := s.GetResolutions(nil)
	assert.NoError(t, err)
	assert.Greater(t, len(resos), 0)

	info, err := s.GetInfo()
	assert.NoError(t, err)
	assert.NotNil(t, info)
	info.Unref()

	m2, err := s.GetSettings()
	assert.NoError(t, err)
	assert.NotNil(t, m2)

	err = s.SetSettings(m2)
	assert.NoError(t, err)

	err = s.SetFramerate(25)
	assert.NoError(t, err)
	m2.Unref()

	streams, err := axvdo.VdoChannelGetAll()
	assert.NoError(t, err)
	assert.NotNil(t, streams)
	assert.Greater(t, len(streams), 0)

	filtermap := axvdo.NewVdoMap()
	filtermap.SetUint32("format", uint32(axvdo.VdoFormatH264))
	streams2, err := axvdo.VdoChannelGetFilterd(filtermap)
	assert.NoError(t, err)
	assert.NotNil(t, streams2)
	assert.Greater(t, len(streams2), 0)
	filtermap.Unref()
}

func LicenseTest(t *testing.T) {
	s := axlicense.LicensekeyVerify("test", 414614, 1, 0)
	assert.False(t, s)
	// TODO: Bring this to work
	//t1, err := acap.LicensekeyGetExpDate("test")
	//fmt.Println(t1, err)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func ParamTests(t *testing.T) {
	appname := "test"
	p, err := axparameter.AXParameterNew(appname)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	param := randSeq(10)
	// Add Test
	err = p.Add(param, "yes", "bool:no,yes")
	assert.NoError(t, err)

	// Set Test 1
	err = p.Set(param, "no", true)
	assert.NoError(t, err)

	// Set Test 2
	err = p.Set("unkown", "no", true)
	assert.Error(t, err)

	// Get Test 2
	pvalue2, err := p.Get(param)
	assert.NoError(t, err)
	assert.Equal(t, "no", pvalue2)

	// Get Test 3
	pvalue3, err := p.Get("Icantfound")
	assert.Error(t, err)
	assert.Equal(t, "", pvalue3)

	// List test
	params, err := p.List()
	assert.NoError(t, err)
	assert.Greater(t, len(params), 0)

	err = p.RegisterCallback(param, func(name, value string, userdata any) {
		fmt.Println("Param Callback Invoked", name, value, userdata)
	}, "mydata")

	assert.NoError(t, err)

	m := glib.NewMainLoop()
	assert.NotNil(t, m)

	l, err := p.List()
	fmt.Println(l, err)

	fmt.Println("Starting gorutine for callback check (10Sec)")
	go time.AfterFunc(time.Second*10, func() {
		m.Quit()
		p.UnregisterCallback(param)
		// Remove Test 1
		err = p.Remove(param)
		assert.NoError(t, err)

		// Remove Test 2
		err = p.Remove("idontknow")
		assert.Error(t, err)

		p.Free()
	})
	m.Run()
}

func MdbTests(t *testing.T) {
	con, err := axmdb.MDBConnectionCreate(func(cerr error) {
		fmt.Println("Error callback invoked", cerr)
		assert.Error(t, cerr)
	})
	assert.NoError(t, err)
	assert.NotNil(t, con)

	subc, err := axmdb.MDBSubscriberConfigCreate("com.axis.analytics_scene_description.v0.beta", "1", func(msg *axmdb.Message) {
		t.Log("Message callback invoked", msg)
	})
	assert.NoError(t, err)

	sub, err := axmdb.MDBSubscriberCreateAsync(con, subc, func(donerr error) {
		fmt.Println("Subscriber created", donerr)
	})

	sub.Destroy()
	subc.Destroy()
	con.Destroy()

}

func EventHandlerTests(t *testing.T) {
	set := axevent.NewAXEventKeyValueSet()
	namespacet1 := "tns1"
	err := set.AddKeyValue("topic0", &namespacet1, "Device", axevent.AXValueTypeString)
	assert.NoError(t, err)
	handler := axevent.NewEventHandler()
	subscription, err2 := handler.Subscribe(set, func(subscription int, event *axevent.AXEvent, userdata any) {
		fmt.Println("EVT Callback invoked", subscription, event.GetTimestamp(), userdata.(string))
	}, "myuserdata")
	assert.NotNil(t, subscription)
	assert.Equal(t, 1, subscription)
	assert.NoError(t, err2)

	declaration, err3 := handler.Declare(set, true, func(declaration int, userdata any) {
		fmt.Println("Event declared successfully", declaration, userdata.(string))
	}, "foobar")
	assert.NotNil(t, declaration)
	assert.NoError(t, err3)

	set2 := axevent.NewAXEventKeyValueSet()
	assert.NotNil(t, set2)
	err4 := set2.AddKeyValue("feature", nil, "myfeature", axevent.AXValueTypeString)
	assert.NoError(t, err4)
	err5 := set2.AddKeyValue("enabled", nil, true, axevent.AXValueTypeBool)
	assert.NoError(t, err5)

	/* declaration2, err6 := handler.DeclareFromTemplate(set2, "com.vendor.PropertyState.Example", func(declaration int, userdata any) {
		fmt.Println("Event with templ declared successfully", declaration, userdata.(string))
	}, "bazfoo")
	assert.NotNil(t, declaration2)
	assert.NoError(t, err6) */

	m := glib.NewMainLoop()
	assert.NotNil(t, m)
	go time.AfterFunc(time.Second*5, func() {
		m.Quit()
		err := handler.Unsubscribe(1)
		assert.NoError(t, err)
		err2 := handler.Undeclare(1)
		assert.NoError(t, err2)
		err3 := handler.Undeclare(2)
		assert.NoError(t, err3)
		handler.Free()
		set.Free()
		set2.Free()
	})
	m.Run()

}

func EventTests(t *testing.T) {
	set := axevent.NewAXEventKeyValueSet()
	assert.NotNil(t, set)

	namespace := "tnaxis"
	err := set.AddKeyValue("topic0", &namespace, "port", axevent.AXValueTypeString)
	assert.NoError(t, err)

	vtype, err := set.GetValueType("topic0", &namespace)
	assert.NoError(t, err)
	assert.Equal(t, vtype, axevent.AXValueTypeString)

	str, err2 := set.GetString("topic0", &namespace)
	assert.NoError(t, err2)
	assert.Equal(t, str, "port")

	err = set.AddKeyValue("topic1", nil, "foobar", axevent.AXValueTypeString)
	assert.NoError(t, err)

	str2, err3 := set.GetString("topic1", nil)
	assert.NoError(t, err3)
	assert.Equal(t, str2, "foobar")

	vtype, err = set.GetValueType("topic1", nil)
	assert.NoError(t, err)
	assert.Equal(t, vtype, axevent.AXValueTypeString)

	err = set.AddKeyValue("topic2", nil, nil, axevent.AXValueTypeString)
	assert.NoError(t, err)

	vtype, err = set.GetValueType("topic2", nil)
	assert.NoError(t, err)
	assert.Equal(t, vtype, axevent.AXValueTypeString)

	err = set.RemoveKey("topic2", nil)
	assert.NoError(t, err)

	err = set.AddKeyValue("topic2", nil, true, axevent.AXValueTypeBool)
	assert.NoError(t, err)

	b1, err3 := set.GetBoolean("topic2", nil)
	assert.NoError(t, err3)
	assert.Equal(t, true, b1)

	err = set.AddKeyValue("topic3", nil, 1, axevent.AXValueTypeInt)
	assert.NoError(t, err)

	i1, err4 := set.GetInteger("topic3", nil)
	assert.NoError(t, err4)
	assert.Equal(t, 1, i1)

	err = set.AddKeyValue("topic4", nil, 1.2, axevent.AXValueTypeDouble)
	assert.NoError(t, err)

	f1, err5 := set.GetDouble("topic4", nil)
	assert.NoError(t, err5)
	assert.Equal(t, 1.2, f1)

	err6 := set.MarkAsSource("topic4", nil)
	assert.NoError(t, err6)

	err7 := set.MarkAsData("topic3", nil)
	assert.NoError(t, err7)

	tag := "mytag"
	err8 := set.MarkAsUserDefined("topic2", nil, &tag)
	assert.NoError(t, err8)

	now := time.Now()
	evt := axevent.NewAxEvent(set, &now)
	assert.NotNil(t, evt)
	assert.NotNil(t, evt.GetTimestamp())
	evt.Free()

	set2 := axevent.NewAXEventKeyValueSet()
	evt2 := axevent.NewAxEvent(set2, nil)
	assert.NotNil(t, evt2)

	evt2.Free()

	set.Free()
	set2.Free()
}
