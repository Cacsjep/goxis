package axmdb

import (
	"errors"
	"fmt"
	"sync"
)

// MDBProvider is a generic provider for messages of type T.
type MDBProvider[T MessageType] struct {
	Topic       string
	Source      string
	ErrorChan   chan *MDBProviderError
	MessageChan chan T
	con         *MDBConnection
	subConfig   *MDBSubscriberConfig
	subscriber  *MDBSubscriber
	once        sync.Once
}

type MDBProviderErrorType int

const (
	MDBProviderErrorTypeConnection MDBProviderErrorType = iota
	MDBProviderErrorTypeSubscriberConfigCreate
	MDBProviderErrorTypeSubscriberCreate
	MDBProviderErrorTypeInvalidMessage
	MDBProviderErrorTypeParseMessage
	MDBProviderErrorTypeEmptyPayload
	MDBProviderErrorSubscribeDone
	MDBProviderErrorSubscribe
)

type MDBProviderError struct {
	Err     error
	ErrType MDBProviderErrorType
}

// NewMDBProvider creates a new MDBProvider.
func NewMDBProvider[T MessageType](source string) (*MDBProvider[T], error) {
	topic, err := resolveTopic[T]()
	if err != nil {
		return nil, err
	}
	return &MDBProvider[T]{
		Topic:       topic,
		Source:      source,
		ErrorChan:   make(chan *MDBProviderError, 5),
		MessageChan: make(chan T, 10),
	}, nil
}

// resolveTopic determines the topic based on the type T.
func resolveTopic[T MessageType]() (string, error) {
	switch any(new(T)).(type) {
	case *SceneDescription:
		return "com.axis.analytics_scene_description.v0.beta", nil
	case *ConsolidatedTrack:
		return "com.axis.consolidated_track.v1.beta", nil
	default:
		return "", errors.New("unsupported message type")
	}
}

// Connect initializes the MDBProvider and processes JSON strings dynamically.
func (mdb *MDBProvider[T]) Connect() {
	con, err := MDBConnectionCreate(func(conerr error) {
		if !mdb.safeToProcced(conerr, MDBProviderErrorTypeConnection) {
			return
		}
	})
	if !mdb.safeToProcced(err, MDBProviderErrorTypeConnection) {
		return
	}
	mdb.con = con

	subConfig, config_create_err := MDBSubscriberConfigCreate(mdb.Topic, mdb.Source, func(msg *Message) {
		if msg.Payload == "" {
			if !mdb.safeToProcced(fmt.Errorf("empty payload received"), MDBProviderErrorTypeEmptyPayload) {
				return
			}
			return
		}
		// Parse JSON into the appropriate type
		var t T
		parsed, parseErr := t.TransformMessage(msg.Payload)
		if !mdb.safeToProcced(parseErr, MDBProviderErrorTypeInvalidMessage) {
			return
		}

		// Type assertion safety
		typedMessage, ok := parsed.(T)
		if !ok {
			if !mdb.safeToProcced(errors.New("failed to cast parsed message to the expected type"), MDBProviderErrorTypeParseMessage) {
				return
			}
			return
		}

		// Send the parsed message to the channel
		mdb.MessageChan <- typedMessage
	})
	if !mdb.safeToProcced(config_create_err, MDBProviderErrorSubscribe) {
		return
	}
	mdb.subConfig = subConfig

	subscriber, create_async_err := MDBSubscriberCreateAsync(con, subConfig, func(onDone error) {
		if !mdb.safeToProcced(onDone, MDBProviderErrorSubscribe) {
			return
		}
	})
	if !mdb.safeToProcced(create_async_err, MDBProviderErrorSubscribe) {
		return
	}
	mdb.subscriber = subscriber
}

func (mdb *MDBProvider[T]) safeToProcced(err error, errType MDBProviderErrorType) bool {
	if err != nil {
		mdb.ErrorChan <- &MDBProviderError{
			Err:     err,
			ErrType: errType,
		}
		mdb.Disconnect()
		return false
	}
	return true
}

// Disconnect cleans up resources.
func (mdb *MDBProvider[T]) Disconnect() {
	mdb.once.Do(func() {
		if mdb.subscriber != nil {
			mdb.subscriber.Destroy()
		}
		if mdb.subConfig != nil {
			mdb.subConfig.Destroy()
		}
		if mdb.con != nil {
			mdb.con.Destroy()
		}
		close(mdb.MessageChan)
		close(mdb.ErrorChan)
	})
}
