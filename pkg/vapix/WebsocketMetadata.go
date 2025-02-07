package vapix

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Cacsjep/goxis/pkg/dbus"
	"github.com/Cacsjep/goxis/pkg/utils"
	"github.com/gorilla/websocket"
)

// https://help.axis.com/en-us/axis-os-knowledge-base#metadata-via-websocket

// VapixWsMetadataStreamRequest represents a request to configure the VAPIX WebSocket metadata stream.
type VapixWsMetadataStreamRequest struct {
	APIVersion string                                  `json:"apiVersion"` // API version of the request, e.g., "1.0".
	Method     string                                  `json:"method"`     // Method to execute, e.g., "events:configure".
	Params     VapixWsMetadataStreamRequestEventParams `json:"params"`     // Request parameters containing event filters.
}

// VapixWsMetadataStreamRequestEventParams defines the parameters for a metadata stream request.
type VapixWsMetadataStreamRequestEventParams struct {
	EventFilterList []VapixWsMetadataStreamRequestEventFilter `json:"eventFilterList"` // List of filters to apply to the events stream.
	ChannelFilter   []string                                  `json:"channelFilter"`
}

// VapixWsMetadataStreamRequestEventFilter represents an event filter for the metadata stream.
type VapixWsMetadataStreamRequestEventFilter struct {
	TopicFilter   string `json:"topicFilter"`   // Specifies the topic to filter, e.g., "tns1:Device/IO/VirtualInput".
	ContentFilter string `json:"contentFilter"` // Specifies additional filtering using XPath syntax.
}

// VapixWsMetadataStreamResponse represents the response structure for the VAPIX metadata stream.
type VapixWsMetadataStreamResponse struct {
	APIVersion string                              `json:"apiVersion"` // API version of the response.
	Method     string                              `json:"method"`     // Method executed, e.g., "events:configure".
	Error      *VapixError                         `json:"error"`      // Error information, if any (assumes VapixError is defined elsewhere).
	Params     VapixWsMetadataStreamResponseParams `json:"params"`     // Response parameters containing notifications.
}

// VapixWsMetadataStreamResponseParams defines the response parameters containing event notifications.
type VapixWsMetadataStreamResponseParams struct {
	Notification VapixWsMetadataStreamResponseNotification `json:"notification"` // Notification details for a metadata event.
}

// VapixWsMetadataStreamResponseNotification represents a single notification from the metadata stream.
type VapixWsMetadataStreamResponseNotification struct {
	Topic     string                 `json:"topic"`     // The topic of the notification, e.g., "tns1:Device/IO/VirtualInput".
	Timestamp int64                  `json:"timestamp"` // Timestamp of the notification in milliseconds since epoch.
	Message   map[string]interface{} `json:"message"`   // Generic message content as key-value pairs.
}

// VapixWsMetadataConsumer represents a consumer for VAPIX WebSocket metadata streams.
type VapixWsMetadataConsumer struct {
	Username      string
	Password      string
	EventFilters  []VapixWsMetadataStreamRequestEventFilter
	RequestConfig *VapixWsMetadataStreamRequest
}

// NewVapixWsMetadataStreamRequest creates a new VAPIX WebSocket metadata stream request with the provided filters.
func NewVapixWsMetadataStreamReqeust(eventsFilter []VapixWsMetadataStreamRequestEventFilter, channelFilters []string) *VapixWsMetadataStreamRequest {
	return &VapixWsMetadataStreamRequest{
		APIVersion: "1.0",
		Method:     "events:configure",
		Params: VapixWsMetadataStreamRequestEventParams{
			EventFilterList: eventsFilter,
			ChannelFilter:   channelFilters,
		},
	}
}

// NewVapixWsMetadataConsumer creates a new VAPIX WebSocket metadata consumer with optional event filters.
func NewVapixWsMetadataConsumer(eventFilters *[]VapixWsMetadataStreamRequestEventFilter, channelFilters []string) *VapixWsMetadataConsumer {
	vwmc := &VapixWsMetadataConsumer{
		EventFilters: []VapixWsMetadataStreamRequestEventFilter{},
	}
	// Use default filter if none provided.
	if eventFilters == nil {
		vwmc.EventFilters = append(vwmc.EventFilters, VapixWsMetadataStreamRequestEventFilter{TopicFilter: ""})
	} else {
		vwmc.EventFilters = *eventFilters
	}
	// Create the metadata stream request configuration.
	vwmc.RequestConfig = NewVapixWsMetadataStreamReqeust(vwmc.EventFilters, channelFilters)
	return vwmc
}

// Connect establishes a WebSocket connection to the VAPIX metadata stream and sends the configuration request.
func (vwmc *VapixWsMetadataConsumer) Connect(sources string) (*websocket.Conn, error) {

	// Retrieve credentials (assumes dbus.RetrieveVapixCredentials is implemented).
	username, password, err := dbus.RetrieveVapixCredentials("root")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve VAPIX credentials: %s", err.Error())
	}
	vwmc.Username = username
	vwmc.Password = password

	// Set up WebSocket headers with authentication.
	headers := http.Header{}
	headers.Set("Authorization", utils.BasicAuthHeader(username, password))

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(INTERNAL_VAPIX_WS_METADATA_STREAM_ENDPOINT+sources, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %v", err)
	}

	// Send the configuration request.
	err = conn.WriteJSON(vwmc.RequestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to send configuration request: %s", err.Error())
	}
	return conn, nil
}

// DumpRequest returns the JSON representation of the metadata stream request for debugging purposes.
func (vwmc *VapixWsMetadataConsumer) DumpRequest() (string, error) {
	if vwmc.RequestConfig == nil {
		return "", fmt.Errorf("RequestConfig is nil")
	}
	jsonData, err := json.MarshalIndent(vwmc.RequestConfig, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
