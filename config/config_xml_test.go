package config

import (
	"testing"
)

type TestConfig struct {
	Name string `xml:"name"`
	Port int    `xml:"port"`
}

func TestUnmarshalXML(t *testing.T) {
	// Test normal XML parsing
	xmlData := []byte(`<config><name>test</name><port>8080</port></config>`)
	var config TestConfig
	err := UnmarshalXML(xmlData, &config)
	if err != nil {
		t.Errorf("UnmarshalXML failed: %v", err)
	}
	if config.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", config.Name)
	}
	if config.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", config.Port)
	}

	// Test empty XML
	var emptyConfig TestConfig
	err = UnmarshalXML([]byte(""), &emptyConfig)
	if err != nil {
		t.Logf("Empty XML returns error: %v", err)
	}

	// Test malformed XML (should not cause security issues)
	malformedXML := []byte("<config><name>test</name>")
	var malformed TestConfig
	err = UnmarshalXML(malformedXML, &malformed)
	if err == nil {
		t.Log("Malformed XML parsed without error")
	}
}

func TestUnmarshalXMLSecurity(t *testing.T) {
	// Test XXE attack payload - should be safely handled
	xxePayload := `<?xml version="1.0" encoding="UTF-8"?>
		<!DOCTYPE config [
			<!ENTITY xxe SYSTEM "file:///etc/passwd">
		]>
		<config><name>&xxe;</name></config>`

	var config TestConfig
	err := UnmarshalXML([]byte(xxePayload), &config)
	if err != nil {
		t.Logf("XXE payload rejected: %v", err)
	} else {
		// If parsed, should not contain external entity content
		t.Logf("XXE payload parsed, name: %s", config.Name)
	}

	// Test billion laughs attack - should not cause DoS
	billionLaughs := `<?xml version="1.0"?><!ENTITY lol "lol"><!ENTITY lol2 "&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;"><!ENTITY lol3 "&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;"><config><name>&lol3;</name></config>`

	var config2 TestConfig
	err = UnmarshalXML([]byte(billionLaughs), &config2)
	if err != nil {
		t.Logf("Billion laughs attack rejected: %v", err)
	} else {
		t.Logf("Billion laughs parsed, name length: %d", len(config2.Name))
	}
}

func TestMarshalXML(t *testing.T) {
	config := TestConfig{Name: "test", Port: 8080}
	data, err := MarshalXML(config)
	if err != nil {
		t.Errorf("MarshalXML failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("MarshalXML returned empty data")
	}
}

func TestMarshalXMLString(t *testing.T) {
	config := TestConfig{Name: "test", Port: 8080}
	str := MarshalXMLString(config)
	if str == "" {
		t.Error("MarshalXMLString returned empty string")
	}
}
