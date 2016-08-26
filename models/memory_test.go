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

var slugifyTests = []struct {
	before string
	after  string
}{
	{"hello world", "hello-world"},
	{"-hello world-", "hello-world"},
	{"HelLo WoRlD", "hello-world"},
	{"  hello world  ", "hello-world"},
	{"!@#$%^&*()_+=`~,./?><;:'\"[]{}\\|hello world", "hello-world"},
}

func TestSlugify(t *testing.T) {
	for _, pair := range slugifyTests {
		result := slugify(pair.before)
		if result != pair.after {
			t.Errorf("expected %s, got %s\n", pair.after, result)
		}
	}
}
