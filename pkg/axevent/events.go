package axevent

import (
	"fmt"
	"reflect"
	"strings"
)

// KeyValueEntrie is a key-value pair for an AXEventKeyValueSet.
type KeyValueEntrie struct {
	Key       string
	Namespace *string
	Value     interface{}
	ValueType AXEventValueType
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

// UnmarshalEvent unmarshals the given event into the provided struct.
func UnmarshalEvent(e *Event, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("value must be a pointer to a struct")
	}

	for i := 0; i < val.Elem().NumField(); i++ {
		field := val.Elem().Field(i)
		if !field.CanSet() {
			continue
		}
		fieldType := val.Elem().Type().Field(i)

		key := fieldType.Tag.Get("eventKey")
		if key == "" {
			key = strings.ToLower(fieldType.Name)
		}

		switch field.Kind() {
		case reflect.Int:
			if intValue, err := e.Kvs.GetInteger(key, nil); err == nil {
				field.SetInt(int64(intValue))
			} else {
				return fmt.Errorf("error getting integer for key %s: %v", key, err)
			}
		case reflect.Float64:
			if fValue, err := e.Kvs.GetDouble(key, nil); err == nil {
				field.SetFloat(fValue)
			} else {
				return fmt.Errorf("error getting double for key %s: %v", key, err)
			}
		case reflect.String:
			if sValue, err := e.Kvs.GetString(key, nil); err == nil {
				field.SetString(sValue)
			} else {
				return fmt.Errorf("error getting string for key %s: %v", key, err)
			}
		case reflect.Bool:
			if boolValue, err := e.Kvs.GetBoolean(key, nil); err == nil {
				field.SetBool(boolValue)
			} else {
				return fmt.Errorf("error getting boolean for key %s: %v", key, err)
			}
		}
	}

	return nil
}
