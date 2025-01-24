package axevent

// <tns1:Device>
// 	<tnsaxis:IO>
// 		<SupervisedPort wstop:topic="true">
// 			<tt:MessageDescription IsProperty="true">
// 				<tt:Source>
// 					<tt:SimpleItemDescription Name="port" Type="xsd:int"></tt:SimpleItemDescription>
// 				</tt:Source>
// 				<tt:Data>
// 					<tt:SimpleItemDescription Name="state" Type="xsd:string"></tt:SimpleItemDescription>
// 					<tt:SimpleItemDescription Name="tampered" Type="xsd:boolean"></tt:SimpleItemDescription>
// 				</tt:Data>
// 			</tt:MessageDescription>
// 		</SupervisedPort>
// 	</tnsaxis:IO>
// </tns1:Device>
func DeviceIoSupervisedPortEventKvs(port *int, tampered *bool, state *string) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "IO"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "SupervisedPort"),
		NewIntKeyValueEntrie("port", port),
		NewBoolKeyValueEntrie("tampered", tampered),
		NewStringKeyValueEntrie("state", state),
	})
}

type DeviceIoSupervisedPortEvent struct {
	Port     int    `eventKey:"port"`
	Tampered bool   `eventKey:"tampered"`
	State    string `eventKey:"state"`
}

// <tns1:Device>
// 	<tnsaxis:IO>
//		<VirtualPort wstop:topic="true">
//			<tt:MessageDescription IsProperty="true">
//				<tt:Source>
//					<tt:SimpleItemDescription Name="port" Type="xsd:int"></tt:SimpleItemDescription>
//				</tt:Source>
//				<tt:Data>
//					<tt:SimpleItemDescription Name="state" Type="xsd:boolean"></tt:SimpleItemDescription>
//				</tt:Data>
//			</tt:MessageDescription>
//		</VirtualPort>
// 	</tnsaxis:IO>
// </tns1:Device>
func DeviceIoVirtualPortEventKvs(port *int, state *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "IO"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "VirtualPort"),
		NewIntKeyValueEntrie("port", port),
		NewBoolKeyValueEntrie("state", state),
	})
}

type DeviceIoVirtualPortEvent struct {
	Port  int  `eventKey:"port"`
	State bool `eventKey:"state"`
}

// <tns1:Device>
// 	<tnsaxis:IO>
//		<OutputPort wstop:topic="true">
//			<tt:MessageDescription IsProperty="true">
//				<tt:Source>
//					<tt:SimpleItemDescription Name="port" Type="xsd:int"></tt:SimpleItemDescription>
//				</tt:Source>
//				<tt:Data>
//					<tt:SimpleItemDescription Name="state" Type="xsd:boolean"></tt:SimpleItemDescription>
//				</tt:Data>
//			</tt:MessageDescription>
//		</OutputPort>
// 	</tnsaxis:IO>
// </tns1:Device>
func DeviceIoOutputPortEventKvs(port *int, state *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "IO"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "OutputPort"),
		NewIntKeyValueEntrie("port", port),
		NewBoolKeyValueEntrie("state", state),
	})
}

type DeviceIoOutputPortEvent struct {
	Port  int  `eventKey:"port"`
	State bool `eventKey:"state"`
}

// <tns1:Device>
// 	<tnsaxis:IO>
//		<VirtualInput wstop:topic="true">
//			<tt:MessageDescription IsProperty="true">
//				<tt:Source>
//					<tt:SimpleItemDescription Name="port" Type="xsd:int"></tt:SimpleItemDescription>
//				</tt:Source>
//				<tt:Data>
//					<tt:SimpleItemDescription Name="active" Type="xsd:boolean"></tt:SimpleItemDescription>
//				</tt:Data>
//			</tt:MessageDescription>
//		</VirtualInput>
// 	</tnsaxis:IO>
// </tns1:Device>
func DeviceIoVirtualInputEventKvs(port *int, active *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "IO"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "VirtualInput"),
		NewIntKeyValueEntrie("port", port),
		NewBoolKeyValueEntrie("active", active),
	})
}

type DeviceIoVirtualInputEvent struct {
	Port   int  `eventKey:"port"`
	Active bool `eventKey:"active"`
}

// <tns1:Device>
// 	<tnsaxis:IO>
//		<Port wstop:topic="true">
//			<tt:MessageDescription IsProperty="true">
//				<tt:Source>
//					<tt:SimpleItemDescription Name="port" Type="xsd:int"></tt:SimpleItemDescription>
//				</tt:Source>
//				<tt:Data>
//					<tt:SimpleItemDescription Name="state" Type="xsd:boolean"></tt:SimpleItemDescription>
//				</tt:Data>
//			</tt:MessageDescription>
//		</Port>
// 	</tnsaxis:IO>
// </tns1:Device>
func DeviceIoPortEventKvs(port *int, state *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "IO"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "Port"),
		NewIntKeyValueEntrie("port", port),
		NewBoolKeyValueEntrie("state", state),
	})
}

type DeviceIoPortEvent struct {
	Port  int  `eventKey:"port"`
	State bool `eventKey:"state"`
}

// <tns1:Device>
//	<tnsaxis:Sensor>
//		<PIR wstop:topic="true">
//			<tt:MessageDescription IsProperty="true">
//				<tt:Source>
//					<tt:SimpleItemDescription Name="sensor" Type="xsd:int"></tt:SimpleItemDescription>
//				</tt:Source>
//				<tt:Data>
//					<tt:SimpleItemDescription Name="state" Type="xsd:boolean"></tt:SimpleItemDescription>
//				</tt:Data>
//			</tt:MessageDescription>
//		</PIR>
//	</tnsaxis:Sensor>
// </tns1:Device>
func DeviceSensorPIREventKvs(sensor *int, state *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "Sensor"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "PIR"),
		NewIntKeyValueEntrie("sensor", sensor),
		NewBoolKeyValueEntrie("state", state),
	})
}

type DeviceSensorPIREvent struct {
	Sensor int  `eventKey:"sensor"`
	State  bool `eventKey:"state"`
}

// <tns1:Device>
//	<tnsaxis:Light>
//		<Status wstop:topic="true">
//			<tt:MessageDescription IsProperty="true">
//				<tt:Source>
//					<tt:SimpleItemDescription Name="id" Type="xsd:int"></tt:SimpleItemDescription>
//				</tt:Source>
//				<tt:Data>
//					<tt:SimpleItemDescription Name="state" Type="xsd:string"></tt:SimpleItemDescription>
//				</tt:Data>
//			</tt:MessageDescription>
//		</Status>
//	</tnsaxis:Light>
// </tns1:Device>
func DeviceLightStatusEventKvs(id *int, state *string) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "Light"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "Status"),
		NewIntKeyValueEntrie("id", id),
		NewStringKeyValueEntrie("state", state),
	})
}

type DeviceLightStatusEvent struct {
	Id    int    `eventKey:"id"`
	State string `eventKey:"state"`
}

// <tns1:Device>
//	<tnsaxis:Status>
//		<SystemReady wstop:topic="true">
//			<tt:MessageDescription IsProperty="true">
//				<tt:Data>
//					<tt:SimpleItemDescription Name="ready" Type="xsd:boolean"></tt:SimpleItemDescription>
//				</tt:Data>
//			</tt:MessageDescription>
//		</SystemReady>
//	</tnsaxis:Status>
// </tns1:Device>
func DeviceStatusSystemReadyEventKvs(ready *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "Status"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "SystemReady"),
		NewBoolKeyValueEntrie("ready", ready),
	})
}

type DeviceStatusSystemReadyEvent struct {
	Ready bool `eventKey:"ready"`
}

// <tns1:Device>
//	<tnsaxis:Status>
//		<Temperature>
//			<Inside wstop:topic="true">
//				<tt:MessageDescription IsProperty="true">
//					<tt:Data>
//						<tt:SimpleItemDescription Name="sensor_level" Type="xsd:boolean"></tt:SimpleItemDescription>
//					</tt:Data>
//				</tt:MessageDescription>
//			</Inside>
//		</Temperature>
//	</tnsaxis:Status>
// </tns1:Device>
func DeviceStatusTemperatureInsideEventKvs(sensor_level *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "Status"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "Temperature"),
		NewTopicKeyValueEntrie("topic3", &OnfivNameSpaceTnsAxis, "Inside"),
		NewBoolKeyValueEntrie("sensor_level", sensor_level),
	})
}

type DeviceStatusTemperatureInsideEvent struct {
	SensorLevel bool `eventKey:"sensor_level"`
}

// <tns1:Device>
//	<tnsaxis:Status>
//		<Temperature>
//			<Above wstop:topic="true">
//				<tt:MessageDescription IsProperty="true">
//					<tt:Data>
//						<tt:SimpleItemDescription Name="sensor_level" Type="xsd:boolean"></tt:SimpleItemDescription>
//					</tt:Data>
//				</tt:MessageDescription>
//			</Above>
//		</Temperature>
//	</tnsaxis:Status>
// </tns1:Device>
func DeviceStatusTemperatureAboveEventKvs(sensor_level *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "Status"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "Temperature"),
		NewTopicKeyValueEntrie("topic3", &OnfivNameSpaceTnsAxis, "Above"),
		NewBoolKeyValueEntrie("sensor_level", sensor_level),
	})
}

type DeviceStatusTemperatureAboveEvent struct {
	SensorLevel bool `eventKey:"sensor_level"`
}

// <tns1:Device>
//	<tnsaxis:Status>
//		<Temperature>
//			<Above_or_below wstop:topic="true">
//				<tt:MessageDescription IsProperty="true">
//				<tt:Data>
//					<tt:SimpleItemDescription Name="sensor_level" Type="xsd:boolean"></tt:SimpleItemDescription>
//				</tt:Data>
//				</tt:MessageDescription>
//			</Above_or_below>
//		</Temperature>
//	</tnsaxis:Status>
// </tns1:Device>
func DeviceStatusTemperatureAboveOrBelowEventKvs(sensor_level *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "Status"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "Temperature"),
		NewTopicKeyValueEntrie("topic3", &OnfivNameSpaceTnsAxis, "Above_or_below"),
		NewBoolKeyValueEntrie("sensor_level", sensor_level),
	})
}

type DeviceStatusTemperatureAboveOrBelowEvent struct {
	SensorLevel bool `eventKey:"sensor_level"`
}

// <tns1:Device>
//	<tnsaxis:Status>
//		<Temperature>
//			<Below wstop:topic="true">
//				<tt:MessageDescription IsProperty="true">
//				<tt:Data>
//					<tt:SimpleItemDescription Name="sensor_level" Type="xsd:boolean"></tt:SimpleItemDescription>
//				</tt:Data>
//				</tt:MessageDescription>
//			</Above_or_below>
//		</Temperature>
//	</tnsaxis:Status>
// </tns1:Device>
func DeviceStatusTemperatureBelowEventKvs(sensor_level *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "Status"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "Temperature"),
		NewTopicKeyValueEntrie("topic3", &OnfivNameSpaceTnsAxis, "Below"),
		NewBoolKeyValueEntrie("sensor_level", sensor_level),
	})
}

type DeviceStatusTemperatureBelowEvent struct {
	SensorLevel bool `eventKey:"sensor_level"`
}

// <tns1:Device>
//	<HardwareFailure>
//		<PowerSupplyFailure>
//			<tnsaxis:PTZPowerFailure wstop:topic="true">
//				<tt:MessageDescription IsProperty="true">
//					<tt:Source>
//						<tt:SimpleItemDescription Name="Token" Type="tt:ReferenceToken"></tt:SimpleItemDescription>
//					</tt:Source>
//					<tt:Data>
//						<tt:SimpleItemDescription Name="Failed" Type="xsd:boolean"></tt:SimpleItemDescription>
//					</tt:Data>
//				</tt:MessageDescription>
//			</tnsaxis:PTZPowerFailure>
//		</PowerSupplyFailure>
//	</HardwareFailure>
// </tns1:Device>
func DeviceHardwareFailurePowerSupplyFailurePTZPowerFailureEventKvs(token *int, failed *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTns1, "HardwareFailure"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTns1, "PowerSupplyFailure"),
		NewTopicKeyValueEntrie("topic3", &OnfivNameSpaceTnsAxis, "PTZPowerFailure"),
		NewIntKeyValueEntrie("Token", token),
		NewBoolKeyValueEntrie("Failed", failed),
	})
}

type DeviceHardwareFailurePowerSupplyFailurePTZPowerFailureEvent struct {
	Token  int  `eventKey:"Token"`
	Failed bool `eventKey:"Failed"`
}

// <tns1:Device>
//	<Trigger>
//		<DigitalInput wstop:topic="true">
//			<tt:MessageDescription IsProperty="true">
//				<tt:Source>
//					<tt:SimpleItemDescription Name="InputToken" Type="tt:ReferenceToken"></tt:SimpleItemDescription>
//				</tt:Source>
//				<tt:Data>
//					<tt:SimpleItemDescription Name="LogicalState" Type="xsd:boolean"></tt:SimpleItemDescription>
//				</tt:Data>
//			</tt:MessageDescription>
//		</DigitalInput>
//	</Trigger>
// </tns1:Device>
func DeviceTriggerDigitalInputEventKvs(inputToken *int, logicalState *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTns1, "Trigger"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTns1, "DigitalInput"),
		NewIntKeyValueEntrie("InputToken", inputToken),
		NewBoolKeyValueEntrie("LogicalState", logicalState),
	})
}

type DeviceTriggerDigitalInputEvent struct {
	InputToken   int  `eventKey:"InputToken"`
	LogicalState bool `eventKey:"LogicalState"`
}

// <tns1:Device>
//	<Trigger>
//		<Relay wstop:topic="true">
//			<tt:MessageDescription IsProperty="true">
//			<tt:Source>
//				<tt:SimpleItemDescription Name="RelayToken" Type="tt:ReferenceToken"></tt:SimpleItemDescription>
//			</tt:Source>
//			<tt:Data>
//				<tt:SimpleItemDescription Name="LogicalState" Type="tt:RelayLogicalState"></tt:SimpleItemDescription>
//			</tt:Data>
//			</tt:MessageDescription>
//		</Relay>
//	</Trigger>
// </tns1:Device>
func DeviceTriggerRelayEventKvs(relayToken *int, logicalState *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTns1, "Trigger"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTns1, "Relay"),
		NewIntKeyValueEntrie("RelayToken", relayToken),
		NewBoolKeyValueEntrie("LogicalState", logicalState),
	})
}

type DeviceTriggerRelayEvent struct {
	RelayToken   int  `eventKey:"RelayToken"`
	LogicalState bool `eventKey:"LogicalState"`
}

// <tns1:Device>
//	<tnsaxis:RingPowerLimitExceeded wstop:topic="true">
//		<tt:MessageDescription IsProperty="true">
//			<tt:Source>
//				<tt:SimpleItemDescription Name="input" Type="xsd:int"></tt:SimpleItemDescription>
//			</tt:Source>
//			<tt:Data>
//				<tt:SimpleItemDescription Name="limit_exceeded" Type="xsd:boolean"></tt:SimpleItemDescription>
//			</tt:Data>
//		</tt:MessageDescription>
//	</tnsaxis:RingPowerLimitExceeded>
// </tns1:Device>
func DeviceRingPowerLimitExceededEventKvs(input *int, limitExceeded *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "RingPowerLimitExceeded"),
		NewIntKeyValueEntrie("input", input),
		NewBoolKeyValueEntrie("limit_exceeded", limitExceeded),
	})
}

type RingPowerLimitExceededEvent struct {
	Input         int  `eventKey:"input"`
	LimitExceeded bool `eventKey:"limit_exceeded"`
}

// <tns1:LightControl>
// 	<tnsaxis:LightStatusChanged>
// 		<Status wstop:topic="true">
// 			<tt:MessageDescription IsProperty="true">
// 				<tt:Data>
// 					<tt:SimpleItemDescription Name="state" Type="xsd:string"></tt:SimpleItemDescription>
// 				</tt:Data>
// 			</tt:MessageDescription>
// 		</Status>
// 	</tnsaxis:LightStatusChanged>
// </tns1:LightControl>
func LightControlLightStatusChangedEventKvs(state *string) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "LightControl"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "LightStatusChanged"),
		NewStringKeyValueEntrie("state", state),
	})
}

type LightControlLightStatusChangedEvent struct {
	State string `eventKey:"state"`
}

// <tns1:VideoSource>
// 	<tnsaxis:LiveStreamAccessed wstop:topic="true">
// 		<tt:MessageDescription IsProperty="true">
// 			<tt:Data>
// 				<tt:SimpleItemDescription Name="accessed" Type="xsd:boolean"></tt:SimpleItemDescription>
// 			</tt:Data>
// 		</tt:MessageDescription>
// 	</tnsaxis:LiveStreamAccessed>
// </tns1:VideoSource>
func VideoSourceLiveStreamAccessedEventKvs(accessed *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "VideoSource"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "LiveStreamAccessed"),
		NewBoolKeyValueEntrie("accessed", accessed),
	})
}

type VideoSourceLiveStreamAccessedEvent struct {
	Accessed bool `eventKey:"accessed"`
}

// <tns1:VideoSource>
// 	<tnsaxis:DayNightVision wstop:topic="true">
// 		<tt:MessageDescription IsProperty="true">
// 			<tt:Source>
// 				<tt:SimpleItemDescription Name="VideoSourceConfigurationToken" Type="xsd:int"></tt:SimpleItemDescription>
// 			</tt:Source>
// 			<tt:Data>
// 				<tt:SimpleItemDescription Name="day" Type="xsd:boolean"></tt:SimpleItemDescription>
// 			</tt:Data>
// 		</tt:MessageDescription>
// 	</tnsaxis:DayNightVision>
// </tns1:VideoSource>
func VideoSourceDayNightVisionEventKvs(videoSourceConfigurationToken *int, day *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "VideoSource"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "DayNightVision"),
		NewIntKeyValueEntrie("VideoSourceConfigurationToken", videoSourceConfigurationToken),
		NewBoolKeyValueEntrie("day", day),
	})
}

type VideoSourceDayNightVisionEvent struct {
	VideoSourceConfigurationToken int  `eventKey:"VideoSourceConfigurationToken"`
	Day                           bool `eventKey:"day"`
}

// <tns1:VideoSource>
// 	<tnsaxis:Tampering wstop:topic="true">
// 		<tt:MessageDescription IsProperty="false">
// 			<tt:Source>
// 				<tt:SimpleItemDescription Name="channel" Type="xsd:int"></tt:SimpleItemDescription>
// 			</tt:Source>
// 			<tt:Data>
// 				<tt:SimpleItemDescription Name="tampering" Type="xsd:int"></tt:SimpleItemDescription>
// 			</tt:Data>
// 		</tt:MessageDescription>
// 	</tnsaxis:Tampering>
// </tns1:VideoSource>
func VideoSourceTamperingEventKvs(channel *int, tampering *int) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "VideoSource"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "Tampering"),
		NewIntKeyValueEntrie("channel", channel),
		NewIntKeyValueEntrie("tampering", tampering),
	})
}

type VideoSourceTamperingEvent struct {
	Channel   int `eventKey:"channel"`
	Tampering int `eventKey:"tampering"`
}

// <tns1:VideoSource>
// 	<tnsaxis:ABR wstop:topic="true">
// 		<tt:MessageDescription IsProperty="true">
// 			<tt:Source>
// 				<tt:SimpleItemDescription Name="VideoSourceConfigurationToken" Type="xsd:int"></tt:SimpleItemDescription>
// 			</tt:Source>
// 			<tt:Data>
// 				<tt:SimpleItemDescription Name="abr_error" Type="xsd:boolean"></tt:SimpleItemDescription>
// 			</tt:Data>
// 		</tt:MessageDescription>
// 	</tnsaxis:ABR>
// </tns1:VideoSource>
func VideoSourceABREventKvs(videoSourceConfigurationToken *int, abrError *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "VideoSource"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "ABR"),
		NewIntKeyValueEntrie("VideoSourceConfigurationToken", videoSourceConfigurationToken),
		NewBoolKeyValueEntrie("abr_error", abrError),
	})
}

type VideoSourceABREvent struct {
	VideoSourceConfigurationToken int  `eventKey:"VideoSourceConfigurationToken"`
	AbrError                      bool `eventKey:"abr_error"`
}

// <tns1:VideoSource>
// 	<GlobalSceneChange>
// 		<ImagingService wstop:topic="true">
// 			<tt:MessageDescription IsProperty="true">
// 				<tt:Source>
// 					<tt:SimpleItemDescription Name="Source" Type="tt:ReferenceToken"></tt:SimpleItemDescription>
// 				</tt:Source>
// 				<tt:Data>
// 					<tt:SimpleItemDescription Name="State" Type="xsd:boolean"></tt:SimpleItemDescription>
// 				</tt:Data>
// 			</tt:MessageDescription>
// 		</ImagingService>
// 	</GlobalSceneChange>
// </tns1:VideoSource>
func VideoSourceGlobalSceneChangeEventKvs(source *int, state *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "VideoSource"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTns1, "GlobalSceneChange"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTns1, "ImagingService"),
		NewIntKeyValueEntrie("Source", source),
		NewBoolKeyValueEntrie("State", state),
	})
}

type VideoSourceGlobalSceneChangeEvent struct {
	Source int  `eventKey:"Source"`
	State  bool `eventKey:"State"`
}

// <tns1:VideoSource>
// 	<MotionAlarm wstop:topic="true">
// 		<tt:MessageDescription IsProperty="true">
// 			<tt:Source>
// 				<tt:SimpleItemDescription Name="Source" Type="tt:ReferenceToken"></tt:SimpleItemDescription>
// 			</tt:Source>
// 			<tt:Data>
// 				<tt:SimpleItemDescription Name="State" Type="xsd:boolean"></tt:SimpleItemDescription>
// 			</tt:Data>
// 		</tt:MessageDescription>
// 	</MotionAlarm>
// </tns1:VideoSource>
func VideoSourceMotionAlarmEventKvs(source *int, state *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "VideoSource"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTns1, "MotionAlarm"),
		NewIntKeyValueEntrie("Source", source),
		NewBoolKeyValueEntrie("State", state),
	})
}

type VideoSourceMotionAlarmEvent struct {
	Source int  `eventKey:"Source"`
	State  bool `eventKey:"State"`
}

// <tns1:PTZController>
// 	<tnsaxis:PTZError wstop:topic="true">
// 		<tt:MessageDescription IsProperty="false">
// 			<tt:Source>
// 				<tt:SimpleItemDescription Name="channel" Type="xsd:int"></tt:SimpleItemDescription>
// 			</tt:Source>
// 			<tt:Data>
// 				<tt:SimpleItemDescription Name="ptz_error" Type="xsd:string"></tt:SimpleItemDescription>
// 			</tt:Data>
// 		</tt:MessageDescription>
// 	</tnsaxis:PTZError>
// </tns1:PTZController>
func PTZControllerPTZErrorEventKvs(channel *int, ptzError *string) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "PTZController"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "PTZError"),
		NewIntKeyValueEntrie("channel", channel),
		NewStringKeyValueEntrie("ptz_error", ptzError),
	})
}

type PTZControllerPTZErrorEvent struct {
	Channel  int    `eventKey:"channel"`
	PTZError string `eventKey:"ptz_error"`
}

// <tns1:PTZController>
// 	<tnsaxis:PTZReady wstop:topic="true">
// 		<tt:MessageDescription IsProperty="true">
// 			<tt:Source>
// 				<tt:SimpleItemDescription Name="channel" Type="xsd:int"></tt:SimpleItemDescription>
// 			</tt:Source>
// 			<tt:Data>
// 				<tt:SimpleItemDescription Name="ready" Type="xsd:boolean"></tt:SimpleItemDescription>
// 			</tt:Data>
// 		</tt:MessageDescription>
// 	</tnsaxis:PTZReady>
// </tns1:PTZController>
func PTZControllerPTZReadyEventKvs(channel *int, ready *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "PTZController"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "PTZReady"),
		NewIntKeyValueEntrie("channel", channel),
		NewBoolKeyValueEntrie("ready", ready),
	})
}

type PTZControllerPTZReadyEvent struct {
	Channel int  `eventKey:"channel"`
	Ready   bool `eventKey:"ready"`
}

// <tns1:Media>
// 	<ConfigurationChanged wstop:topic="true">
// 		<tt:MessageDescription IsProperty="false">
// 			<tt:Source>
// 				<tt:SimpleItemDescription Name="Type" Type="xsd:string"></tt:SimpleItemDescription>
// 				<tt:SimpleItemDescription Name="Token" Type="tt:ReferenceToken"></tt:SimpleItemDescription>
// 			</tt:Source>
// 		</tt:MessageDescription>
// 	</ConfigurationChanged>
// </tns1:Media>
func MediaConfigurationChangedEventKvs(eventType *string, token *string) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Media"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTns1, "ConfigurationChanged"),
		NewStringKeyValueEntrie("Type", eventType),
		NewStringKeyValueEntrie("Token", token),
	})
}

type MediaConfigurationChangedEvent struct {
	Type  string `eventKey:"Type"`
	Token string `eventKey:"Token"`
}

// <tns1:Media>
// 	<ProfileChanged wstop:topic="true">
// 		<tt:MessageDescription IsProperty="false">
// 			<tt:Source>
// 				<tt:SimpleItemDescription Name="Token" Type="tt:ReferenceToken"></tt:SimpleItemDescription>
// 			</tt:Source>
// 		</tt:MessageDescription>
// 	</ProfileChanged>
// </tns1:Media>
func MediaProfileChangedEventKvs(token *string) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Media"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTns1, "ProfileChanged"),
		NewStringKeyValueEntrie("Token", token),
	})
}

type MediaProfileChangedEvent struct {
	Token string `eventKey:"Token"`
}

// <tnsaxis:CameraApplicationPlatform>
// 	<ObjectAnalytics>
// 		<Device1Scenario1 wstop:topic="true">
// 			<tt:MessageDescription IsProperty="true">
// 				<tt:Data>
// 					<tt:SimpleItemDescription Name="active" Type="xsd:boolean"></tt:SimpleItemDescription>
// 				</tt:Data>
// 			</tt:MessageDescription>
// 		</Device1Scenario1>
// 	</ObjectAnalytics>
// </tnsaxis:CameraApplicationPlatform>
func CameraApplicationPlatformDevice1Scenario1EventKvs(active *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTnsAxis, "CameraApplicationPlatform"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "ObjectAnalytics"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "Device1Scenario1"),
		NewBoolKeyValueEntrie("active", active),
	})
}

type CameraApplicationPlatformDevice1Scenario1Event struct {
	Active bool `eventKey:"active"`
}

// <tnsaxis:CameraApplicationPlatform>
// 	<ObjectAnalytics>
// 		<Device1ScenarioANY wstop:topic="true">
// 			<tt:MessageDescription IsProperty="true">
// 				<tt:Data>
// 					<tt:SimpleItemDescription Name="active" Type="xsd:boolean"></tt:SimpleItemDescription>
// 				</tt:Data>
// 			</tt:MessageDescription>
// 		</Device1ScenarioANY>
// 	</ObjectAnalytics>
// </tnsaxis:CameraApplicationPlatform>
func CameraApplicationPlatformDevice1ScenarioANYEventKvs(active *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTnsAxis, "CameraApplicationPlatform"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "ObjectAnalytics"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "Device1ScenarioANY"),
		NewBoolKeyValueEntrie("active", active),
	})
}

type CameraApplicationPlatformDevice1ScenarioANYEvent struct {
	Active bool `eventKey:"active"`
}

// <tnsaxis:CameraApplicationPlatform>
// 	<ObjectAnalytics>
// 		<xinternal_data wstop:topic="true">
// 			<tt:MessageDescription IsProperty="false">
// 				<tt:Data>
// 					<tt:SimpleItemDescription Name="svgframe" Type="xsd:string"></tt:SimpleItemDescription>
// 				</tt:Data>
// 			</tt:MessageDescription>
// 		</xinternal_data>
// 	</ObjectAnalytics>
// </tnsaxis:CameraApplicationPlatform>
func CameraApplicationPlatformXInternalDataEventKvs(svgFrame *string) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTnsAxis, "CameraApplicationPlatform"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "ObjectAnalytics"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "xinternal_data"),
		NewStringKeyValueEntrie("svgframe", svgFrame),
	})
}

type CameraApplicationPlatformXInternalDataEvent struct {
	SvgFrame string `eventKey:"svgframe"`
}

// <tnsaxis:Storage>
// 	<Alert wstop:topic="true">
// 		<tt:MessageDescription IsProperty="true">
// 			<tt:Source>
// 				<tt:SimpleItemDescription Name="disk_id" Type="xsd:string"></tt:SimpleItemDescription>
// 			</tt:Source>
// 			<tt:Data>
// 				<tt:SimpleItemDescription Name="alert" Type="xsd:boolean"></tt:SimpleItemDescription>
// 				<tt:SimpleItemDescription Name="overall_health" Type="xsd:int"></tt:SimpleItemDescription>
// 				<tt:SimpleItemDescription Name="temperature" Type="xsd:int"></tt:SimpleItemDescription>
// 				<tt:SimpleItemDescription Name="wear" Type="xsd:int"></tt:SimpleItemDescription>
// 			</tt:Data>
// 		</tt:MessageDescription>
// 	</Alert>
// </tnsaxis:Storage>
func StorageAlertEventKvs(diskID *string, alert *bool, overallHealth *int, temperature *int, wear *int) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTnsAxis, "Storage"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "Alert"),
		NewStringKeyValueEntrie("disk_id", diskID),
		NewBoolKeyValueEntrie("alert", alert),
		NewIntKeyValueEntrie("overall_health", overallHealth),
		NewIntKeyValueEntrie("temperature", temperature),
		NewIntKeyValueEntrie("wear", wear),
	})
}

type StorageAlertEvent struct {
	DiskID        string `eventKey:"disk_id"`
	Alert         bool   `eventKey:"alert"`
	OverallHealth int    `eventKey:"overall_health"`
	Temperature   int    `eventKey:"temperature"`
	Wear          int    `eventKey:"wear"`
}

// <tnsaxis:Storage>
// 	<Disruption wstop:topic="true">
// 		<tt:MessageDescription IsProperty="true">
// 			<tt:Source>
// 				<tt:SimpleItemDescription Name="disk_id" Type="xsd:string"></tt:SimpleItemDescription>
// 			</tt:Source>
// 			<tt:Data>
// 				<tt:SimpleItemDescription Name="disruption" Type="xsd:boolean"></tt:SimpleItemDescription>
// 			</tt:Data>
// 		</tt:MessageDescription>
// 	</Disruption>
// </tnsaxis:Storage>
func StorageDisruptionEventKvs(diskID *string, disruption *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTnsAxis, "Storage"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "Disruption"),
		NewStringKeyValueEntrie("disk_id", diskID),
		NewBoolKeyValueEntrie("disruption", disruption),
	})
}

type StorageDisruptionEvent struct {
	DiskID     string `eventKey:"disk_id"`
	Disruption bool   `eventKey:"disruption"`
}

// <tnsaxis:Storage>
// 	<Recording wstop:topic="true">
// 		<tt:MessageDescription IsProperty="true">
// 			<tt:Data>
// 				<tt:SimpleItemDescription Name="recording" Type="xsd:boolean"></tt:SimpleItemDescription>
// 			</tt:Data>
// 		</tt:MessageDescription>
// 	</Recording>
// </tnsaxis:Storage>
func StorageRecordingEventKvs(recording *bool) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTnsAxis, "Storage"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "Recording"),
		NewBoolKeyValueEntrie("recording", recording),
	})
}

type StorageRecordingEvent struct {
	Recording bool `eventKey:"recording"`
}
