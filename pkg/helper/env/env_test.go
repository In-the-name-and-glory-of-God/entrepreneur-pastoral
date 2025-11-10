package env

import (
	"os"
	"testing"
	"time"
)

func TestGetString(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		def      string
		expected string
		setEnv   bool
	}{
		{
			name:     "existing environment variable",
			key:      "TEST_STRING_VAR",
			value:    "test-value",
			def:      "default",
			expected: "test-value",
			setEnv:   true,
		},
		{
			name:     "non-existing environment variable",
			key:      "NON_EXISTING_VAR",
			def:      "default-value",
			expected: "default-value",
			setEnv:   false,
		},
		{
			name:     "empty string value",
			key:      "EMPTY_STRING_VAR",
			value:    "",
			def:      "default",
			expected: "",
			setEnv:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before test
			os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			result := GetString(tt.key, tt.def)

			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGetStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		def      []string
		expected []string
		setEnv   bool
	}{
		{
			name:     "space-separated values",
			key:      "TEST_SLICE_VAR",
			value:    "value1 value2 value3",
			def:      []string{"default"},
			expected: []string{"value1", "value2", "value3"},
			setEnv:   true,
		},
		{
			name:     "single value",
			key:      "TEST_SINGLE_VAR",
			value:    "single",
			def:      []string{"default"},
			expected: []string{"single"},
			setEnv:   true,
		},
		{
			name:     "non-existing variable",
			key:      "NON_EXISTING_SLICE",
			def:      []string{"default1", "default2"},
			expected: []string{"default1", "default2"},
			setEnv:   false,
		},
		{
			name:     "empty string",
			key:      "EMPTY_SLICE_VAR",
			value:    "",
			def:      []string{"default"},
			expected: []string{""},
			setEnv:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before test
			os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			result := GetStringSlice(tt.key, tt.def)

			if len(result) != len(tt.expected) {
				t.Fatalf("Expected slice length %d, got %d", len(tt.expected), len(result))
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("At index %d: expected '%s', got '%s'", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		def      int
		expected int
		setEnv   bool
	}{
		{
			name:     "valid positive integer",
			key:      "TEST_INT_VAR",
			value:    "42",
			def:      0,
			expected: 42,
			setEnv:   true,
		},
		{
			name:     "valid negative integer",
			key:      "TEST_NEG_INT_VAR",
			value:    "-10",
			def:      0,
			expected: -10,
			setEnv:   true,
		},
		{
			name:     "zero value",
			key:      "TEST_ZERO_VAR",
			value:    "0",
			def:      99,
			expected: 0,
			setEnv:   true,
		},
		{
			name:     "invalid integer",
			key:      "TEST_INVALID_INT",
			value:    "not-a-number",
			def:      100,
			expected: 100,
			setEnv:   true,
		},
		{
			name:     "non-existing variable",
			key:      "NON_EXISTING_INT",
			def:      50,
			expected: 50,
			setEnv:   false,
		},
		{
			name:     "float value",
			key:      "TEST_FLOAT_VAR",
			value:    "3.14",
			def:      10,
			expected: 10,
			setEnv:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before test
			os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			result := GetInt(tt.key, tt.def)

			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		def      bool
		expected bool
		setEnv   bool
	}{
		{
			name:     "true value",
			key:      "TEST_BOOL_TRUE",
			value:    "true",
			def:      false,
			expected: true,
			setEnv:   true,
		},
		{
			name:     "false value",
			key:      "TEST_BOOL_FALSE",
			value:    "false",
			def:      true,
			expected: false,
			setEnv:   true,
		},
		{
			name:     "1 as true",
			key:      "TEST_BOOL_ONE",
			value:    "1",
			def:      false,
			expected: true,
			setEnv:   true,
		},
		{
			name:     "0 as false",
			key:      "TEST_BOOL_ZERO",
			value:    "0",
			def:      true,
			expected: false,
			setEnv:   true,
		},
		{
			name:     "invalid boolean",
			key:      "TEST_BOOL_INVALID",
			value:    "not-a-bool",
			def:      true,
			expected: true,
			setEnv:   true,
		},
		{
			name:     "non-existing variable",
			key:      "NON_EXISTING_BOOL",
			def:      false,
			expected: false,
			setEnv:   false,
		},
		{
			name:     "uppercase TRUE",
			key:      "TEST_BOOL_UPPER",
			value:    "TRUE",
			def:      false,
			expected: true,
			setEnv:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before test
			os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			result := GetBool(tt.key, tt.def)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetDuration(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		def      time.Duration
		expected time.Duration
		setEnv   bool
	}{
		{
			name:     "seconds duration",
			key:      "TEST_DURATION_SEC",
			value:    "30s",
			def:      time.Minute,
			expected: 30 * time.Second,
			setEnv:   true,
		},
		{
			name:     "minutes duration",
			key:      "TEST_DURATION_MIN",
			value:    "5m",
			def:      time.Second,
			expected: 5 * time.Minute,
			setEnv:   true,
		},
		{
			name:     "hours duration",
			key:      "TEST_DURATION_HOUR",
			value:    "2h",
			def:      time.Minute,
			expected: 2 * time.Hour,
			setEnv:   true,
		},
		{
			name:     "combined duration",
			key:      "TEST_DURATION_COMBINED",
			value:    "1h30m",
			def:      time.Second,
			expected: time.Hour + 30*time.Minute,
			setEnv:   true,
		},
		{
			name:     "invalid duration",
			key:      "TEST_DURATION_INVALID",
			value:    "not-a-duration",
			def:      10 * time.Second,
			expected: 10 * time.Second,
			setEnv:   true,
		},
		{
			name:     "non-existing variable",
			key:      "NON_EXISTING_DURATION",
			def:      5 * time.Minute,
			expected: 5 * time.Minute,
			setEnv:   false,
		},
		{
			name:     "milliseconds duration",
			key:      "TEST_DURATION_MS",
			value:    "500ms",
			def:      time.Second,
			expected: 500 * time.Millisecond,
			setEnv:   true,
		},
		{
			name:     "zero duration",
			key:      "TEST_DURATION_ZERO",
			value:    "0s",
			def:      time.Minute,
			expected: 0,
			setEnv:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before test
			os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			result := GetDuration(tt.key, tt.def)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEnvironmentIsolation(t *testing.T) {
	// This test ensures that environment variables don't interfere with each other

	key1 := "TEST_ISO_VAR1"
	key2 := "TEST_ISO_VAR2"

	os.Setenv(key1, "value1")
	os.Setenv(key2, "value2")
	defer os.Unsetenv(key1)
	defer os.Unsetenv(key2)

	result1 := GetString(key1, "default")
	result2 := GetString(key2, "default")

	if result1 != "value1" {
		t.Errorf("Expected 'value1' for key1, got '%s'", result1)
	}

	if result2 != "value2" {
		t.Errorf("Expected 'value2' for key2, got '%s'", result2)
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Test that concurrent access to environment variables is safe
	key := "TEST_CONCURRENT_VAR"
	os.Setenv(key, "123")
	defer os.Unsetenv(key)

	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			result := GetInt(key, 0)
			if result != 123 {
				t.Errorf("Expected 123, got %d", result)
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
