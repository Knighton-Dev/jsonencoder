package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestEncodeJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "test.json with special chars",
			input:    `{"Server":"192.168.1.1","Database":"cool_db","User_Id":"super_user","Password":"kRp-CK@D2DCc3d9QoZG3WBBg@i2j!g"}`,
			expected: `"{\"Server\": \"192.168.1.1\", \"Database\": \"cool_db\", \"User_Id\": \"super_user\", \"Password\": \"kRp-CK@D2DCc3d9QoZG3WBBg@i2j!g\"}"`,
			wantErr:  false,
		},
		{
			name:     "simple object",
			input:    `{"key": "value"}`,
			expected: `"{\"key\": \"value\"}"`,
			wantErr:  false,
		},
		{
			name:     "object with numbers",
			input:    `{"age": 30, "active": true}`,
			expected: `"{\"age\": 30, \"active\": true}"`,
			wantErr:  false,
		},
		{
			name:     "array",
			input:    `["item1", "item2"]`,
			expected: `"[\"item1\", \"item2\"]"`,
			wantErr:  false,
		},
		{
			name:     "nested object",
			input:    `{"user": {"name": "John", "age": 30}}`,
			expected: `"{\"user\": {\"name\": \"John\", \"age\": 30}}"`,
			wantErr:  false,
		},
		{
			name:    "invalid JSON",
			input:   `{"invalid": json}`,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   ``,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encodeJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Both result and expected are JSON-encoded strings, so unmarshal them for comparison
				var gotStr, wantStr string
				if err := json.Unmarshal([]byte(result), &gotStr); err == nil {
					if err := json.Unmarshal([]byte(tt.expected), &wantStr); err == nil {
						// Try to unmarshal the string contents as JSON for logical comparison
						var gotObj, wantObj interface{}
						if err1 := json.Unmarshal([]byte(gotStr), &gotObj); err1 == nil {
							if err2 := json.Unmarshal([]byte(wantStr), &wantObj); err2 == nil {
								if !equalJSON(gotObj, wantObj) {
									t.Errorf("encodeJSON() = %v, want %v", result, tt.expected)
								}
								return
							}
						}
					}
				}
				// Fallback: compare the encoded string directly
				if result != tt.expected {
					t.Errorf("encodeJSON() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

// equalJSON compares two unmarshaled JSON objects for deep equality
func equalJSON(a, b interface{}) bool {
	switch aVal := a.(type) {
	case map[string]interface{}:
		bVal, ok := b.(map[string]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for k, v := range aVal {
			if !equalJSON(v, bVal[k]) {
				return false
			}
		}
		return true
	case []interface{}:
		bVal, ok := b.([]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for i := range aVal {
			if !equalJSON(aVal[i], bVal[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}

func TestDecodeJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple object",
			input:    `"{\"key\": \"value\"}"`,
			expected: `{"key": "value"}`,
			wantErr:  false,
		},
		{
			name:     "object with numbers",
			input:    `"{\"age\": 30, \"active\": true}"`,
			expected: `{"age": 30, "active": true}`,
			wantErr:  false,
		},
		{
			name:     "array",
			input:    `"[\"item1\", \"item2\"]"`,
			expected: `["item1", "item2"]`,
			wantErr:  false,
		},
		{
			name:    "invalid encoded string",
			input:   `not-a-json-string`,
			wantErr: true,
		},
		{
			name:    "decoded result not valid JSON",
			input:   `"not json"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := decodeJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("decodeJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReadFromFile(t *testing.T) {
	// Create a temporary file for testing
	tempFile := "/tmp/test_json_encoder.json"
	content := `{"test": "content"}`

	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tempFile)

	result, err := readFromFile(tempFile)
	if err != nil {
		t.Errorf("readFromFile() error = %v", err)
		return
	}

	if result != content {
		t.Errorf("readFromFile() = %v, want %v", result, content)
	}
}

func TestReadFromFileWithWhitespace(t *testing.T) {
	// Create a temporary file with whitespace for testing
	tempFile := "/tmp/test_json_encoder_ws.json"
	content := "  \n  {\"test\": \"content\"}  \n  "
	expected := `{"test": "content"}`

	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tempFile)

	result, err := readFromFile(tempFile)
	if err != nil {
		t.Errorf("readFromFile() error = %v", err)
		return
	}

	if result != expected {
		t.Errorf("readFromFile() = %v, want %v", result, expected)
	}
}

func TestRoundTrip(t *testing.T) {
	// Test encoding then decoding returns the original
	original := `{"name": "John Doe", "age": 30, "hobbies": ["reading", "coding"]}`

	encoded, err := encodeJSON(original)
	if err != nil {
		t.Fatalf("Failed to encode JSON: %v", err)
	}

	decoded, err := decodeJSON(encoded)
	if err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	// Compare decoded and original as JSON objects for logical equality
	var gotObj, wantObj interface{}
	if err := json.Unmarshal([]byte(decoded), &gotObj); err != nil {
		t.Fatalf("Decoded output is not valid JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(original), &wantObj); err != nil {
		t.Fatalf("Original input is not valid JSON: %v", err)
	}
	if !equalJSON(gotObj, wantObj) {
		t.Errorf("Round trip failed: got %v, want %v", decoded, original)
	}
}
