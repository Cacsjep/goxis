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
		ErrorChan:   make(chan *MDBProviderError, 10),
		MessageChan: make(chan T, 100),
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
	con, err := MDBConnectionCreate(func(err error) {
		if err != nil {
			mdb.ErrorChan <- &MDBProviderError{
				Err:     err,
				ErrType: MDBProviderErrorTypeConnection,
			}
		}
	})
	if err != nil {
		mdb.ErrorChan <- &MDBProviderError{
			Err:     err,
			ErrType: MDBProviderErrorTypeConnection,
		}
		return
	}
	mdb.con = con

	subConfig, config_create_err := MDBSubscriberConfigCreate(mdb.Topic, mdb.Source, func(msg *Message) {
		if msg.Payload == "" {
			mdb.ErrorChan <- &MDBProviderError{
				Err:     fmt.Errorf("empty payload received"),
				ErrType: MDBProviderErrorTypeEmptyPayload,
			}
			return
		}
		// Parse JSON into the appropriate type
		var t T
		parsed, parseErr := t.TransformMessage(msg.Payload)
		if parseErr != nil {
			mdb.ErrorChan <- &MDBProviderError{
				Err:     parseErr,
				ErrType: MDBProviderErrorTypeInvalidMessage,
			}
			return
		}

		// Type assertion safety
		typedMessage, ok := parsed.(T)
		if !ok {
			mdb.ErrorChan <- &MDBProviderError{
				Err:     errors.New("failed to cast parsed message to the expected type"),
				ErrType: MDBProviderErrorTypeParseMessage,
			}
			return
		}

		// Send the parsed message to the channel
		mdb.MessageChan <- typedMessage
	})
	if config_create_err != nil {
		mdb.ErrorChan <- &MDBProviderError{
			Err:     err,
			ErrType: MDBProviderErrorTypeSubscriberConfigCreate,
		}
		mdb.cleanup()
		return
	}
	mdb.subConfig = subConfig

	subscriber, create_async_err := MDBSubscriberCreateAsync(con, subConfig, func(onDone error) {
		if onDone != nil {
			mdb.ErrorChan <- &MDBProviderError{
				Err:     onDone,
				ErrType: MDBProviderErrorSubscribeDone,
			}
		}
	})
	if create_async_err != nil {
		mdb.ErrorChan <- &MDBProviderError{
			Err:     err,
			ErrType: MDBProviderErrorSubscribe,
		}
		mdb.cleanup()
		return
	}
	mdb.subscriber = subscriber
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

// cleanup handles errors and cleans up resources.
func (mdb *MDBProvider[T]) cleanup() {
	mdb.Disconnect()
}
