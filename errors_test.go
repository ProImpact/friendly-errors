package friendlyerrors_test

import (
	"testing"

	friendlyerrors "github.com/ProImpact/friendly-errors"
)

// TestStruct defines a test struct with various JSON tags
type TestStruct struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Age         int    `json:"age" validate:"required,min=18"`
	Password    string `json:"password" validate:"required,min=8"`
	ConfirmPass string `json:"confirm_password" validate:"eqfield=Password"`
	Website     string `json:"website" validate:"url"`
	Username    string `json:"username" validate:"required,alphanum"`
}

// TestNestedStruct defines a nested test struct
type TestNestedStruct struct {
	Address string `json:"address" validate:"required"`
	City    string `json:"city" validate:"required"`
	ZipCode string `json:"zip_code" validate:"required,numeric"`
}

// TestComplexStruct defines a complex struct with nested structures
type TestComplexStruct struct {
	User     TestStruct       `json:"user" validate:"required"`
	Contacts []TestStruct     `json:"contacts" validate:"required,min=1"`
	Address  TestNestedStruct `json:"address" validate:"required"`
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected map[string]string
	}{
		{
			name: "valid struct",
			input: TestStruct{
				Name:        "John Doe",
				Email:       "john@example.com",
				Age:         25,
				Password:    "password123",
				ConfirmPass: "password123",
				Website:     "https://example.com",
				Username:    "johndoe",
			},
			expected: nil,
		},
		{
			name: "invalid struct with required fields",
			input: TestStruct{
				Name:        "",
				Email:       "invalid-email",
				Age:         15,
				Password:    "123",
				ConfirmPass: "different",
				Website:     "not-a-url",
				Username:    "john doe", // contains space
			},
			expected: map[string]string{
				"name":             "the field 'name' is required",
				"email":            "the field 'email' must be a valid email address",
				"age":              "the field 'age' must be at least 18",
				"password":         "the field 'password' must have a length of 8",
				"confirm_password": "the field 'confirm_password' must be equal to the field 'Password'",
				"website":          "the field 'website' must be a valid URL",
				"username":         "the field 'username' must contain only alphanumeric characters",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := friendlyerrors.ValidateStruct(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d errors, got %d", len(tt.expected), len(result))
				t.Logf("Expected: %v", tt.expected)
				t.Logf("Got: %v", result)
				return
			}

			for expectedField, expectedMsg := range tt.expected {
				if actualMsg, exists := result[expectedField]; !exists {
					t.Errorf("Expected error for field '%s', but got none", expectedField)
				} else if actualMsg != expectedMsg {
					t.Errorf("For field '%s': expected '%s', got '%s'", expectedField, expectedMsg, actualMsg)
				}
			}
		})
	}
}

func TestValidateStructDeep(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected map[string]string
	}{
		{
			name: "valid complex struct",
			input: TestComplexStruct{
				User: TestStruct{
					Name:        "John Doe",
					Email:       "john@example.com",
					Age:         25,
					Password:    "password123",
					ConfirmPass: "password123",
					Website:     "https://example.com",
					Username:    "johndoe",
				},
				Contacts: []TestStruct{
					{
						Name:        "Jane Doe",
						Email:       "jane@example.com",
						Age:         23,
						Password:    "password456",
						ConfirmPass: "password456",
						Website:     "https://jane.com",
						Username:    "janedoe",
					},
				},
				Address: TestNestedStruct{
					Address: "123 Main St",
					City:    "New York",
					ZipCode: "10001",
				},
			},
			expected: nil,
		},
		{
			name: "invalid complex struct with nested errors",
			input: TestComplexStruct{
				User: TestStruct{
					Name:        "",
					Email:       "invalid-email",
					Age:         15,
					Password:    "123",
					ConfirmPass: "different",
					Website:     "not-a-url",
					Username:    "john doe",
				},
				Contacts: []TestStruct{
					{
						Name:        "",
						Email:       "another-invalid",
						Age:         10,
						Password:    "short",
						ConfirmPass: "mismatch",
						Website:     "bad-url",
						Username:    "bad user",
					},
				},
				Address: TestNestedStruct{
					Address: "",
					City:    "",
					ZipCode: "abc",
				},
			},
			expected: map[string]string{
				"user.name":                    "the field 'user.name' is required",
				"user.email":                   "the field 'user.email' must be a valid email address",
				"user.age":                     "the field 'user.age' must be at least 18",
				"user.password":                "the field 'user.password' must have a length of 8",
				"user.confirm_password":        "the field 'user.confirm_password' must be equal to the field 'Password'",
				"user.website":                 "the field 'user.website' must be a valid URL",
				"user.username":                "the field 'user.username' must contain only alphanumeric characters",
				"contacts[0].name":             "the field 'contacts[0].name' is required",
				"contacts[0].email":            "the field 'contacts[0].email' must be a valid email address",
				"contacts[0].age":              "the field 'contacts[0].age' must be at least 18",
				"contacts[0].password":         "the field 'contacts[0].password' must have a length of 8",
				"contacts[0].confirm_password": "the field 'contacts[0].confirm_password' must be equal to the field 'Password'",
				"contacts[0].website":          "the field 'contacts[0].website' must be a valid URL",
				"contacts[0].username":         "the field 'contacts[0].username' must contain only alphanumeric characters",
				"address.address":              "the field 'address.address' is required",
				"address.city":                 "the field 'address.city' is required",
				"address.zip_code":             "the field 'address.zip_code' must be a valid numeric value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := friendlyerrors.ValidateStructDeep(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d errors, got %d", len(tt.expected), len(result))
				t.Logf("Expected: %v", tt.expected)
				t.Logf("Got: %v", result)
				return
			}

			for expectedField, expectedMsg := range tt.expected {
				if actualMsg, exists := result[expectedField]; !exists {
					t.Errorf("Expected error for field '%s', but got none", expectedField)
				} else if actualMsg != expectedMsg {
					t.Errorf("For field '%s': expected '%s', got '%s'", expectedField, expectedMsg, actualMsg)
				}
			}
		})
	}
}

func TestValidateSliceDeep(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected map[string]string
	}{
		{
			name: "valid slice of structs",
			input: []TestStruct{
				{
					Name:        "John Doe",
					Email:       "john@example.com",
					Age:         25,
					Password:    "password123",
					ConfirmPass: "password123",
					Website:     "https://example.com",
					Username:    "johndoe",
				},
				{
					Name:        "Jane Doe",
					Email:       "jane@example.com",
					Age:         23,
					Password:    "password456",
					ConfirmPass: "password456",
					Website:     "https://jane.com",
					Username:    "janedoe",
				},
			},
			expected: nil,
		},
		{
			name: "invalid slice of structs",
			input: []TestStruct{
				{
					Name:        "",
					Email:       "invalid-email",
					Age:         15,
					Password:    "123",
					ConfirmPass: "different",
					Website:     "not-a-url",
					Username:    "john doe",
				},
				{
					Name:        "",
					Email:       "another-invalid",
					Age:         10,
					Password:    "short",
					ConfirmPass: "mismatch",
					Website:     "bad-url",
					Username:    "bad user",
				},
			},
			expected: map[string]string{
				"[0].name":             "the field '[0].name' is required",
				"[0].email":            "the field '[0].email' must be a valid email address",
				"[0].age":              "the field '[0].age' must be at least 18",
				"[0].password":         "the field '[0].password' must have a length of 8",
				"[0].confirm_password": "the field '[0].confirm_password' must be equal to the field 'Password'",
				"[0].website":          "the field '[0].website' must be a valid URL",
				"[0].username":         "the field '[0].username' must contain only alphanumeric characters",
				"[1].name":             "the field '[1].name' is required",
				"[1].email":            "the field '[1].email' must be a valid email address",
				"[1].age":              "the field '[1].age' must be at least 18",
				"[1].password":         "the field '[1].password' must have a length of 8",
				"[1].confirm_password": "the field '[1].confirm_password' must be equal to the field 'Password'",
				"[1].website":          "the field '[1].website' must be a valid URL",
				"[1].username":         "the field '[1].username' must contain only alphanumeric characters",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := friendlyerrors.ValidateSliceDeep(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d errors, got %d", len(tt.expected), len(result))
				t.Logf("Expected: %v", tt.expected)
				t.Logf("Got: %v", result)
				return
			}

			for expectedField, expectedMsg := range tt.expected {
				if actualMsg, exists := result[expectedField]; !exists {
					t.Errorf("Expected error for field '%s', but got none", expectedField)
				} else if actualMsg != expectedMsg {
					t.Errorf("For field '%s': expected '%s', got '%s'", expectedField, expectedMsg, actualMsg)
				}
			}
		})
	}
}

func TestValidateAny(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected map[string]string
	}{
		{
			name: "valid struct",
			input: TestStruct{
				Name:        "John Doe",
				Email:       "john@example.com",
				Age:         25,
				Password:    "password123",
				ConfirmPass: "password123",
				Website:     "https://example.com",
				Username:    "johndoe",
			},
			expected: nil,
		},
		{
			name: "valid slice",
			input: []TestStruct{
				{
					Name:        "John Doe",
					Email:       "john@example.com",
					Age:         25,
					Password:    "password123",
					ConfirmPass: "password123",
					Website:     "https://example.com",
					Username:    "johndoe",
				},
			},
			expected: nil,
		},
		{
			name: "invalid struct",
			input: TestStruct{
				Name:        "",
				Email:       "invalid-email",
				Age:         15,
				Password:    "123",
				ConfirmPass: "different",
				Website:     "not-a-url",
				Username:    "john doe",
			},
			expected: map[string]string{
				"name":             "the field 'name' is required",
				"email":            "the field 'email' must be a valid email address",
				"age":              "the field 'age' must be at least 18",
				"password":         "the field 'password' must have a length of 8",
				"confirm_password": "the field 'confirm_password' must be equal to the field 'Password'",
				"website":          "the field 'website' must be a valid URL",
				"username":         "the field 'username' must contain only alphanumeric characters",
			},
		},
		{
			name: "invalid slice",
			input: []TestStruct{
				{
					Name:        "",
					Email:       "invalid-email",
					Age:         15,
					Password:    "123",
					ConfirmPass: "different",
					Website:     "not-a-url",
					Username:    "john doe",
				},
			},
			expected: map[string]string{
				"[0].name":             "the field '[0].name' is required",
				"[0].email":            "the field '[0].email' must be a valid email address",
				"[0].age":              "the field '[0].age' must be at least 18",
				"[0].password":         "the field '[0].password' must have a length of 8",
				"[0].confirm_password": "the field '[0].confirm_password' must be equal to the field 'Password'",
				"[0].website":          "the field '[0].website' must be a valid URL",
				"[0].username":         "the field '[0].username' must contain only alphanumeric characters",
			},
		},
		{
			name:     "invalid input type",
			input:    "not a struct or slice",
			expected: map[string]string{"root": "provided value is neither a struct nor a slice/array"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := friendlyerrors.ValidateAny(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d errors, got %d", len(tt.expected), len(result))
				t.Logf("Expected: %v", tt.expected)
				t.Logf("Got: %v", result)
				return
			}

			for expectedField, expectedMsg := range tt.expected {
				if actualMsg, exists := result[expectedField]; !exists {
					t.Errorf("Expected error for field '%s', but got none", expectedField)
				} else if actualMsg != expectedMsg {
					t.Errorf("For field '%s': expected '%s', got '%s'", expectedField, expectedMsg, actualMsg)
				}
			}
		})
	}
}

func TestValidateStructWithPointer(t *testing.T) {
	// Test with pointer to struct
	invalidStruct := &TestStruct{
		Name:        "",
		Email:       "invalid-email",
		Age:         15,
		Password:    "123",
		ConfirmPass: "different",
		Website:     "not-a-url",
		Username:    "john doe",
	}

	expected := map[string]string{
		"name":             "the field 'name' is required",
		"email":            "the field 'email' must be a valid email address",
		"age":              "the field 'age' must be at least 18",
		"password":         "the field 'password' must have a length of 8",
		"confirm_password": "the field 'confirm_password' must be equal to the field 'Password'",
		"website":          "the field 'website' must be a valid URL",
		"username":         "the field 'username' must contain only alphanumeric characters",
	}

	result := friendlyerrors.ValidateStruct(invalidStruct)

	if len(result) != len(expected) {
		t.Errorf("Expected %d errors, got %d", len(expected), len(result))
		return
	}

	for expectedField, expectedMsg := range expected {
		if actualMsg, exists := result[expectedField]; !exists {
			t.Errorf("Expected error for field '%s', but got none", expectedField)
		} else if actualMsg != expectedMsg {
			t.Errorf("For field '%s': expected '%s', got '%s'", expectedField, expectedMsg, actualMsg)
		}
	}
}

func TestValidateStructWithNilPointer(t *testing.T) {
	var nilStruct *TestStruct

	result := friendlyerrors.ValidateStruct(nilStruct)
	if result != nil {
		t.Errorf("Expected nil result for nil pointer, got %v", result)
	}
}

func TestValidateStructDeepWithPointer(t *testing.T) {
	invalidComplexStruct := &TestComplexStruct{
		User: TestStruct{
			Name:  "",
			Email: "invalid-email",
			Age:   15,
		},
		Contacts: []TestStruct{
			{
				Name:  "",
				Email: "another-invalid",
				Age:   10,
			},
		},
		Address: TestNestedStruct{
			Address: "",
			City:    "",
		},
	}

	expected := map[string]string{
		"user.name":         "the field 'user.name' is required",
		"user.email":        "the field 'user.email' must be a valid email address",
		"user.age":          "the field 'user.age' must be at least 18",
		"contacts[0].name":  "the field 'contacts[0].name' is required",
		"contacts[0].email": "the field 'contacts[0].email' must be a valid email address",
		"contacts[0].age":   "the field 'contacts[0].age' must be at least 18",
		"address.address":   "the field 'address.address' is required",
		"address.city":      "the field 'address.city' is required",
	}

	result := friendlyerrors.ValidateStructDeep(invalidComplexStruct)

	if len(result) != len(expected) {
		t.Errorf("Expected %d errors, got %d", len(expected), len(result))
		t.Logf("Expected: %v", expected)
		t.Logf("Got: %v", result)
		return
	}

	for expectedField, expectedMsg := range expected {
		if actualMsg, exists := result[expectedField]; !exists {
			t.Errorf("Expected error for field '%s', but got none", expectedField)
		} else if actualMsg != expectedMsg {
			t.Errorf("For field '%s': expected '%s', got '%s'", expectedField, expectedMsg, actualMsg)
		}
	}
}
