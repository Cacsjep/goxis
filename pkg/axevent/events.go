package axevent

import "fmt"

type keyvalues struct {
	key        string
	namespace  *string
	value      interface{}
	value_type AXEventValueType
}

func tnsAxisEvent(topic0 string, topic1 string, topic2 *string, topic3 *string, keyvalues []*keyvalues) (*AXEventKeyValueSet, error) {
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
			if err := kvs.AddKeyValue(kv.key, kv.namespace, kv.value, kv.value_type); err != nil {
				return nil, fmt.Errorf("failed to add key-value for %s: %w", kv.key, err)
			}
		}

	}
	return kvs, nil
}

func NewStringPointer(value string) *string {
	return &value
}

func NewIntPointer(value int) *int {
	return &value
}
