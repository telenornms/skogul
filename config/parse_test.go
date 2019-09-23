package config

import (
	"testing"
	//"fmt"
)

func TestFile(t *testing.T) {
	c, err := File("test.json")
	if err != nil {
		t.Errorf("File() failed: %v", err)
	}
	if c == nil {
		t.Errorf("File() returned nil config")
	}
}

func TestHelpSender(t *testing.T) {
	doc, err := HelpSender("mysql")
	if err != nil {
		t.Errorf("HelpSender(\"mysql\") didn't work: %v", err)
	}
	doc.Print()
}
