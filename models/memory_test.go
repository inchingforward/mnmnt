package models

import "testing"

var requiredFieldsTests = []struct {
	memory  *Memory
	message string
}{
	{&Memory{}, "memory with no fields"},
	{&Memory{Title: "A Title"}, "memory with no details"},
	{&Memory{Title: "A Title", Details: "Details"}, "memory with missing lat/long"},
}

func TestMemoryRequiresTitle(t *testing.T) {
	for _, rft := range requiredFieldsTests {
		err := AddMemory(rft.memory)
		if err == nil {
			t.Errorf("expected an error for %v\n", rft.message)
		}
	}
}
