package validation

import (
	"aviasales-shared-core/domain/interfaces"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// FieldValidator provides field validation functionality
type FieldValidator struct{}

// NewFieldValidator creates a new FieldValidator
func NewFieldValidator() interfaces.Validator {
	return &FieldValidator{}
}

// Validate validates a struct using reflection and tags
func (v *FieldValidator) Validate(data interface{}) error {
	value := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	// Handle pointers
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return fmt.Errorf("validation target cannot be nil")
		}
		value = value.Elem()
		typ = typ.Elem()
	}

	if value.Kind() != reflect.Struct {
		return fmt.Errorf("validation target must be a struct")
	}

	errors := []string{}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		fieldName := fieldType.Name
		fieldValue := field.Interface()

		fieldErrors := v.validateField(fieldName, fieldValue, tag)
		errors = append(errors, fieldErrors...)
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// validateField validates a single field
func (v *FieldValidator) validateField(name string, value interface{}, tag string) []string {
	errors := []string{}
	rules := strings.Split(tag, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		parts := strings.Split(rule, "=")
		ruleName := parts[0]
		var ruleValue string
		if len(parts) > 1 {
			ruleValue = parts[1]
		}

		if err := v.applyRule(name, value, ruleName, ruleValue); err != nil {
			errors = append(errors, err.Error())
		}
	}

	return errors
}

// applyRule applies a validation rule
func (v *FieldValidator) applyRule(name string, value interface{}, ruleName, ruleValue string) error {
	switch ruleName {
	case "required":
		return v.validateRequired(name, value)
	case "min":
		return v.validateMin(name, value, ruleValue)
	case "max":
		return v.validateMax(name, value, ruleValue)
	case "email":
		return v.validateEmail(name, value)
	case "url":
		return v.validateURL(name, value)
	case "pattern":
		return v.validatePattern(name, value, ruleValue)
	default:
		return fmt.Errorf("unknown validation rule: %s", ruleName)
	}
}

// validateRequired checks if field is not empty
func (v *FieldValidator) validateRequired(name string, value interface{}) error {
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.String:
		if val.String() == "" {
			return fmt.Errorf("field '%s' is required", name)
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		if val.Len() == 0 {
			return fmt.Errorf("field '%s' is required", name)
		}
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			return fmt.Errorf("field '%s' is required", name)
		}
	case reflect.Invalid:
		return fmt.Errorf("field '%s' is required", name)
	}

	return nil
}

// validateMin checks minimum value/length
func (v *FieldValidator) validateMin(name string, value interface{}, minStr string) error {
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.String:
		if len(val.String()) < parseInt(minStr) {
			return fmt.Errorf("field '%s' must be at least %s characters", name, minStr)
		}
	case reflect.Slice, reflect.Array:
		if val.Len() < parseInt(minStr) {
			return fmt.Errorf("field '%s' must have at least %s items", name, minStr)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() < int64(parseInt(minStr)) {
			return fmt.Errorf("field '%s' must be at least %s", name, minStr)
		}
	}

	return nil
}

// validateMax checks maximum value/length
func (v *FieldValidator) validateMax(name string, value interface{}, maxStr string) error {
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.String:
		if len(val.String()) > parseInt(maxStr) {
			return fmt.Errorf("field '%s' must be at most %s characters", name, maxStr)
		}
	case reflect.Slice, reflect.Array:
		if val.Len() > parseInt(maxStr) {
			return fmt.Errorf("field '%s' must have at most %s items", name, maxStr)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() > int64(parseInt(maxStr)) {
			return fmt.Errorf("field '%s' must be at most %s", name, maxStr)
		}
	}

	return nil
}

// validateEmail checks email format
func (v *FieldValidator) validateEmail(name string, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("field '%s' must be a string for email validation", name)
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(str) {
		return fmt.Errorf("field '%s' must be a valid email address", name)
	}

	return nil
}

// validateURL checks URL format
func (v *FieldValidator) validateURL(name string, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("field '%s' must be a string for URL validation", name)
	}

	// Empty string is allowed for optional URL fields
	if str == "" {
		return nil
	}

	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(str) {
		return fmt.Errorf("field '%s' must be a valid URL", name)
	}

	return nil
}

// validatePattern checks pattern match
func (v *FieldValidator) validatePattern(name string, value interface{}, pattern string) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("field '%s' must be a string for pattern validation", name)
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern for field '%s': %v", name, err)
	}

	if !regex.MatchString(str) {
		return fmt.Errorf("field '%s' does not match required pattern", name)
	}

	return nil
}

// parseInt parses string to int with default 0
func parseInt(str string) int {
	if str == "" {
		return 0
	}

	var result int
	fmt.Sscanf(str, "%d", &result)
	return result
}
