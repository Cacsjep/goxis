package axevent

import (
	"fmt"
)

// KeyValueEntrie is a key-value pair for an AXEventKeyValueSet.
type KeyValueEntrie struct {
	Key       string
	Namespace *string
	Value     interface{}
	ValueType AXEventValueType
}

// Mark a key in the AXEventKeyValueSet as a source. A source key is an identifier used to distinguish between multiple instances of the same event declaration.
// E.g. if a device has multiple I/O ports then event declarations that represent the state of each port will have the same keys but different values.
// The key that represents which port the event represents should be marked as source and the key which represents the state should be marked as data.
// Please note that although it is possible to mark more than one key as a source, only events with zero or one source keys can be used to trigger actions.
type AxEventKeyValueSetSourceMark struct {
	Key       string
	Namespace *string
}

// Mark a key in the AXEventKeyValueSet as data. A data key is a key that represents the state of what the event represents.
// E.g. an event declaration that represents an I/O port should have a key marked as data which represents the state, high or low, of the port.
// Please note that although it is possible to mark more than one key as data, only events with one and only one data key can be used to trigger actions.
type AxEventKeyValueSetDataMark struct {
	Key       string
	Namespace *string
}

// Mark a key in AXEventKeyValueSet with an user defined tag.
type AxEventKeyValueSetUserDefineMark struct {
	Key       string
	Namespace *string
	Tag       *string
}

// Set the nice names of a key/value pair in the AXEventKeyValueSet.
// Nice names can be used to display human-readable information about the key/value pair.
type AxEventKeyValueSetNiceNames struct {
	Key           string
	Namespace     *string
	KeyNiceName   *string
	ValueNiceName *string
}

// NewTns1AxisEvent creates a new AXEventKeyValueSet with the given topics and key-value pairs.
func NewTns1AxisEvent(topic0 string, topic1 string, topic2 *string, topic3 *string, keyvalues []*KeyValueEntrie) (*AXEventKeyValueSet, error) {
	kvs := NewAXEventKeyValueSet()
	if err := kvs.AddKeyValue("topic0", &OnfivNameSpaceTns1, topic0, AXValueTypeString); err != nil {
		return nil, fmt.Errorf("failed to add key-value for topic0: %w", err)
	}
	if err := kvs.AddKeyValue("topic1", &OnfivNameSpaceTnsAxis, topic1, AXValueTypeString); err != nil {
		return nil, fmt.Errorf("failed to add key-value for topic1: %w", err)
	}
	if topic2 != nil {
		if err := kvs.AddKeyValue("topic2", &OnfivNameSpaceTnsAxis, *topic2, AXValueTypeString); err != nil {
			return nil, fmt.Errorf("failed to add key-value for topic2: %w", err)
		}
	}
	if topic3 != nil {
		if err := kvs.AddKeyValue("topic3", &OnfivNameSpaceTnsAxis, *topic3, AXValueTypeString); err != nil {
			return nil, fmt.Errorf("failed to add key-value for topic2: %w", err)
		}
	}
	if keyvalues != nil {
		for _, kv := range keyvalues {
			if err := kvs.AddKeyValue(kv.Key, kv.Namespace, kv.Value, kv.ValueType); err != nil {
				return nil, fmt.Errorf("failed to add key-value for %s: %w", kv.Key, err)
			}
		}

	}
	return kvs, nil
}

// TnsAxisEvent creates a new AXEventKeyValueSet with the given topics and key-value pairs.
func NewTnsAxisEvent(topic0 string, topic1 string, topic2 *string, topic3 *string, keyvalues []*KeyValueEntrie) (*AXEventKeyValueSet, error) {
	kvs := NewAXEventKeyValueSet()
	if err := kvs.AddKeyValue("topic0", &OnfivNameSpaceTnsAxis, topic0, AXValueTypeString); err != nil {
		return nil, fmt.Errorf("failed to add key-value for topic0: %w", err)
	}
	if err := kvs.AddKeyValue("topic1", &OnfivNameSpaceTnsAxis, topic1, AXValueTypeString); err != nil {
		return nil, fmt.Errorf("failed to add key-value for topic1: %w", err)
	}
	if topic2 != nil {
		if err := kvs.AddKeyValue("topic2", &OnfivNameSpaceTnsAxis, *topic2, AXValueTypeString); err != nil {
			return nil, fmt.Errorf("failed to add key-value for topic2: %w", err)
		}
	}
	if topic3 != nil {
		if err := kvs.AddKeyValue("topic3", &OnfivNameSpaceTnsAxis, *topic3, AXValueTypeString); err != nil {
			return nil, fmt.Errorf("failed to add key-value for topic3: %w", err)
		}
	}
	if keyvalues != nil {
		for _, kv := range keyvalues {
			if err := kvs.AddKeyValue(kv.Key, kv.Namespace, kv.Value, kv.ValueType); err != nil {
				return nil, fmt.Errorf("failed to add key-value for %s: %w", kv.Key, err)
			}
		}

	}
	return kvs, nil
}
