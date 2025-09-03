package validation

import (
	"strings"
	"testing"
)

// Test structures for validation
type TestUser struct {
	Name     string `validate:"required,min=2,max=50"`
	Email    string `validate:"required,email"`
	Age      int    `validate:"min=18,max=120"`
	Website  string `validate:"url"`
	Username string `validate:"required,pattern=^[a-zA-Z0-9_]+$"`
}

type TestProduct struct {
	Title       string   `validate:"required,min=3"`
	Description string   `validate:"max=500"`
	Tags        []string `validate:"min=1,max=10"`
	Price       int      `validate:"min=0"`
}

func TestFieldValidator_Interface(t *testing.T) {
	validator := NewFieldValidator()
	if validator == nil {
		t.Fatal("NewFieldValidator() returned nil")
	}
}

func TestFieldValidator_ValidateValidStruct(t *testing.T) {
	validator := NewFieldValidator()

	user := TestUser{
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      25,
		Website:  "https://example.com",
		Username: "john_doe123",
	}

	err := validator.Validate(user)
	if err != nil {
		t.Errorf("Validate() valid struct error = %v", err)
	}
}

func TestFieldValidator_ValidatePointer(t *testing.T) {
	validator := NewFieldValidator()

	user := &TestUser{
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      25,
		Website:  "https://example.com",
		Username: "john_doe123",
	}

	err := validator.Validate(user)
	if err != nil {
		t.Errorf("Validate() valid pointer error = %v", err)
	}
}

func TestFieldValidator_ValidateNilPointer(t *testing.T) {
	validator := NewFieldValidator()

	var user *TestUser = nil

	err := validator.Validate(user)
	if err == nil {
		t.Error("Validate() nil pointer should return error")
	}

	if !strings.Contains(err.Error(), "cannot be nil") {
		t.Errorf("Validate() nil pointer error = %v, should contain 'cannot be nil'", err)
	}
}

func TestFieldValidator_ValidateNonStruct(t *testing.T) {
	validator := NewFieldValidator()

	err := validator.Validate("not a struct")
	if err == nil {
		t.Error("Validate() non-struct should return error")
	}

	if !strings.Contains(err.Error(), "must be a struct") {
		t.Errorf("Validate() non-struct error = %v, should contain 'must be a struct'", err)
	}
}

func TestFieldValidator_RequiredValidation(t *testing.T) {
	validator := NewFieldValidator()

	// Missing required name
	user := TestUser{
		Email:    "john@example.com",
		Age:      25,
		Website:  "https://example.com",
		Username: "john_doe123",
	}

	err := validator.Validate(user)
	if err == nil {
		t.Error("Validate() missing required field should return error")
	}

	if !strings.Contains(err.Error(), "Name") || !strings.Contains(err.Error(), "required") {
		t.Errorf("Validate() required error = %v, should mention Name and required", err)
	}
}

func TestFieldValidator_MinValidation(t *testing.T) {
	validator := NewFieldValidator()

	// Name too short
	user := TestUser{
		Name:     "J",
		Email:    "john@example.com",
		Age:      25,
		Website:  "https://example.com",
		Username: "john_doe123",
	}

	err := validator.Validate(user)
	if err == nil {
		t.Error("Validate() field below minimum should return error")
	}

	if !strings.Contains(err.Error(), "Name") || !strings.Contains(err.Error(), "at least") {
		t.Errorf("Validate() min error = %v, should mention Name and 'at least'", err)
	}
}

func TestFieldValidator_MaxValidation(t *testing.T) {
	validator := NewFieldValidator()

	// Name too long
	user := TestUser{
		Name:     strings.Repeat("a", 51),
		Email:    "john@example.com",
		Age:      25,
		Website:  "https://example.com",
		Username: "john_doe123",
	}

	err := validator.Validate(user)
	if err == nil {
		t.Error("Validate() field above maximum should return error")
	}

	if !strings.Contains(err.Error(), "Name") || !strings.Contains(err.Error(), "at most") {
		t.Errorf("Validate() max error = %v, should mention Name and 'at most'", err)
	}
}

func TestFieldValidator_EmailValidation(t *testing.T) {
	validator := NewFieldValidator()

	testCases := []struct {
		email     string
		shouldErr bool
	}{
		{"john@example.com", false},
		{"user.name@domain.co.uk", false},
		{"test+tag@example.org", false},
		{"invalid-email", true},
		{"@example.com", true},
		{"john@", true},
		{"john@.com", true},
		{"", true}, // Required field
	}

	for _, tc := range testCases {
		user := TestUser{
			Name:     "John Doe",
			Email:    tc.email,
			Age:      25,
			Website:  "https://example.com",
			Username: "john_doe123",
		}

		err := validator.Validate(user)
		if tc.shouldErr && err == nil {
			t.Errorf("Validate() email '%s' should return error", tc.email)
		}
		if !tc.shouldErr && err != nil {
			t.Errorf("Validate() email '%s' should not return error, got %v", tc.email, err)
		}
	}
}

func TestFieldValidator_URLValidation(t *testing.T) {
	validator := NewFieldValidator()

	testCases := []struct {
		url       string
		shouldErr bool
	}{
		{"https://example.com", false},
		{"http://test.org", false},
		{"https://subdomain.example.com/path", false},
		{"invalid-url", true},
		{"ftp://example.com", true}, // Only http/https allowed
		{"", false},                 // Optional field
	}

	for _, tc := range testCases {
		user := TestUser{
			Name:     "John Doe",
			Email:    "john@example.com",
			Age:      25,
			Website:  tc.url,
			Username: "john_doe123",
		}

		err := validator.Validate(user)
		if tc.shouldErr && err == nil {
			t.Errorf("Validate() url '%s' should return error", tc.url)
		}
		if !tc.shouldErr && err != nil {
			t.Errorf("Validate() url '%s' should not return error, got %v", tc.url, err)
		}
	}
}

func TestFieldValidator_PatternValidation(t *testing.T) {
	validator := NewFieldValidator()

	testCases := []struct {
		username  string
		shouldErr bool
	}{
		{"john_doe123", false},
		{"user123", false},
		{"test_user", false},
		{"john-doe", true}, // Dash not allowed
		{"john doe", true}, // Space not allowed
		{"john@doe", true}, // Special chars not allowed
		{"", true},         // Required field
	}

	for _, tc := range testCases {
		user := TestUser{
			Name:     "John Doe",
			Email:    "john@example.com",
			Age:      25,
			Website:  "https://example.com",
			Username: tc.username,
		}

		err := validator.Validate(user)
		if tc.shouldErr && err == nil {
			t.Errorf("Validate() username '%s' should return error", tc.username)
		}
		if !tc.shouldErr && err != nil {
			t.Errorf("Validate() username '%s' should not return error, got %v", tc.username, err)
		}
	}
}

func TestFieldValidator_SliceValidation(t *testing.T) {
	validator := NewFieldValidator()

	// Valid product
	product := TestProduct{
		Title:       "Test Product",
		Description: "A great product",
		Tags:        []string{"tag1", "tag2"},
		Price:       100,
	}

	err := validator.Validate(product)
	if err != nil {
		t.Errorf("Validate() valid product error = %v", err)
	}

	// Too few tags
	product.Tags = []string{}
	err = validator.Validate(product)
	if err == nil {
		t.Error("Validate() too few tags should return error")
	}

	// Too many tags
	product.Tags = make([]string, 11)
	for i := range product.Tags {
		product.Tags[i] = "tag"
	}
	err = validator.Validate(product)
	if err == nil {
		t.Error("Validate() too many tags should return error")
	}
}

func TestFieldValidator_IntValidation(t *testing.T) {
	validator := NewFieldValidator()

	// Age too low
	user := TestUser{
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      17,
		Website:  "https://example.com",
		Username: "john_doe123",
	}

	err := validator.Validate(user)
	if err == nil {
		t.Error("Validate() age below minimum should return error")
	}

	// Age too high
	user.Age = 121
	err = validator.Validate(user)
	if err == nil {
		t.Error("Validate() age above maximum should return error")
	}
}

func TestFieldValidator_MultipleErrors(t *testing.T) {
	validator := NewFieldValidator()

	// Multiple validation errors
	user := TestUser{
		Name:     "",              // Required, min
		Email:    "invalid-email", // Email format
		Age:      17,              // Min age
		Website:  "invalid-url",   // URL format
		Username: "john-doe",      // Pattern
	}

	err := validator.Validate(user)
	if err == nil {
		t.Error("Validate() multiple errors should return error")
	}

	errorMsg := err.Error()

	// Check that multiple errors are included
	if !strings.Contains(errorMsg, "Name") {
		t.Error("Validate() should include Name error")
	}
	if !strings.Contains(errorMsg, "Email") {
		t.Error("Validate() should include Email error")
	}
	if !strings.Contains(errorMsg, "Age") {
		t.Error("Validate() should include Age error")
	}
}

func TestFieldValidator_UnknownRule(t *testing.T) {
	validator := NewFieldValidator()

	type TestStruct struct {
		Field string `validate:"unknown_rule"`
	}

	test := TestStruct{Field: "value"}

	err := validator.Validate(test)
	if err == nil {
		t.Error("Validate() unknown rule should return error")
	}

	if !strings.Contains(err.Error(), "unknown validation rule") {
		t.Errorf("Validate() unknown rule error = %v, should mention 'unknown validation rule'", err)
	}
}
