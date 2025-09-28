package friendlyerrors

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var errorMessages = map[string]string{
	// Comparisons between fields
	"eqcsfield":  "the field '%s' must be equal to the field '%s'",
	"eqfield":    "the field '%s' must be equal to the field '%s'",
	"gtcsfield":  "the field '%s' must be greater than the field '%s'",
	"gtecsfield": "the field '%s' must be greater than or equal to the field '%s'",
	"gtefield":   "the field '%s' must be greater than or equal to the field '%s'",
	"gtfield":    "the field '%s' must be greater than the field '%s'",
	"ltcsfield":  "the field '%s' must be less than the field '%s'",
	"ltecsfield": "the field '%s' must be less than or equal to the field '%s'",
	"ltefield":   "the field '%s' must be less than or equal to the field '%s'",
	"ltfield":    "the field '%s' must be less than the field '%s'",
	"necsfield":  "the field '%s' must not be equal to the field '%s'",
	"nefield":    "the field '%s' must not be equal to the field '%s'",

	// Network validations
	"cidr":             "the field '%s' must be a valid CIDR notation",
	"cidrv4":           "the field '%s' must be a valid IPv4 CIDR notation",
	"cidrv6":           "the field '%s' must be a valid IPv6 CIDR notation",
	"datauri":          "the field '%s' must be a valid data URI",
	"fqdn":             "the field '%s' must be a valid fully qualified domain name",
	"hostname":         "the field '%s' must be a valid hostname according to RFC 952",
	"hostname_port":    "the field '%s' must be a valid host:port format",
	"hostname_rfc1123": "the field '%s' must be a valid hostname according to RFC 1123",
	"ip":               "the field '%s' must be a valid IP address",
	"ip4_addr":         "the field '%s' must be a valid IPv4 address",
	"ip6_addr":         "the field '%s' must be a valid IPv6 address",
	"ip_addr":          "the field '%s' must be a valid IP address",
	"ipv4":             "the field '%s' must be a valid IPv4 address",
	"ipv6":             "the field '%s' must be a valid IPv6 address",
	"mac":              "the field '%s' must be a valid MAC address",
	"tcp4_addr":        "the field '%s' must be a valid TCP IPv4 address",
	"tcp6_addr":        "the field '%s' must be a valid TCP IPv6 address",
	"tcp_addr":         "the field '%s' must be a valid TCP address",
	"udp4_addr":        "the field '%s' must be a valid UDP IPv4 address",
	"udp6_addr":        "the field '%s' must be a valid UDP IPv6 address",
	"udp_addr":         "the field '%s' must be a valid UDP address",
	"unix_addr":        "the field '%s' must be a valid Unix domain socket address",
	"uri":              "the field '%s' must be a valid URI",
	"url":              "the field '%s' must be a valid URL",
	"http_url":         "the field '%s' must be a valid HTTP URL",
	"url_encoded":      "the field '%s' must be a valid URL encoded string",
	"urn_rfc2141":      "the field '%s' must be a valid URN according to RFC 2141",

	// String validations
	"alpha":           "the field '%s' must contain only alphabetic characters",
	"alphanum":        "the field '%s' must contain only alphanumeric characters",
	"alphanumunicode": "the field '%s' must contain only alphanumeric Unicode characters",
	"alphaunicode":    "the field '%s' must contain only alphabetic Unicode characters",
	"ascii":           "the field '%s' must contain only ASCII characters",
	"boolean":         "the field '%s' must be a valid boolean value",
	"contains":        "the field '%s' must contain the text '%s'",
	"containsany":     "the field '%s' must contain any of the characters '%s'",
	"containsrune":    "the field '%s' must contain the character '%s'",
	"endsnotwith":     "the field '%s' must not end with '%s'",
	"endswith":        "the field '%s' must end with '%s'",
	"excludes":        "the field '%s' must not contain the text '%s'",
	"excludesall":     "the field '%s' must not contain any of the characters '%s'",
	"excludesrune":    "the field '%s' must not contain the character '%s'",
	"lowercase":       "the field '%s' must be in lowercase",
	"multibyte":       "the field '%s' must contain multi-byte characters",
	"number":          "the field '%s' must be a valid number",
	"numeric":         "the field '%s' must be a valid numeric value",
	"printascii":      "the field '%s' must contain only printable ASCII characters",
	"startsnotwith":   "the field '%s' must not start with '%s'",
	"startswith":      "the field '%s' must start with '%s'",
	"uppercase":       "the field '%s' must be in uppercase",

	// Format validations
	"base64":                        "the field '%s' must be a valid Base64 string",
	"base64url":                     "the field '%s' must be a valid Base64URL string",
	"base64rawurl":                  "the field '%s' must be a valid Base64RawURL string",
	"bic":                           "the field '%s' must be a valid Business Identifier Code (BIC)",
	"bcp47_language_tag":            "the field '%s' must be a valid BCP 47 language tag",
	"btc_addr":                      "the field '%s' must be a valid Bitcoin address",
	"btc_addr_bech32":               "the field '%s' must be a valid Bitcoin Bech32 address",
	"credit_card":                   "the field '%s' must be a valid credit card number",
	"mongodb":                       "the field '%s' must be a valid MongoDB ObjectID",
	"mongodb_connection_string":     "the field '%s' must be a valid MongoDB connection string",
	"cron":                          "the field '%s' must be a valid cron expression",
	"spicedb":                       "the field '%s' must be a valid SpiceDb identifier",
	"datetime":                      "the field '%s' must be a valid datetime format '%s'",
	"e164":                          "the field '%s' must be a valid E.164 formatted phone number",
	"ein":                           "the field '%s' must be a valid U.S. Employer Identification Number",
	"email":                         "the field '%s' must be a valid email address",
	"eth_addr":                      "the field '%s' must be a valid Ethereum address",
	"hexadecimal":                   "the field '%s' must be a valid hexadecimal string",
	"hexcolor":                      "the field '%s' must be a valid hex color",
	"hsl":                           "the field '%s' must be a valid HSL color",
	"hsla":                          "the field '%s' must be a valid HSLA color",
	"html":                          "the field '%s' must contain valid HTML tags",
	"html_encoded":                  "the field '%s' must be HTML encoded",
	"isbn":                          "the field '%s' must be a valid ISBN",
	"isbn10":                        "the field '%s' must be a valid ISBN-10",
	"isbn13":                        "the field '%s' must be a valid ISBN-13",
	"issn":                          "the field '%s' must be a valid ISSN",
	"iso3166_1_alpha2":              "the field '%s' must be a valid ISO 3166-1 alpha-2 country code",
	"iso3166_1_alpha3":              "the field '%s' must be a valid ISO 3166-1 alpha-3 country code",
	"iso3166_1_alpha_numeric":       "the field '%s' must be a valid ISO 3166-1 numeric country code",
	"iso3166_2":                     "the field '%s' must be a valid ISO 3166-2 country subdivision code",
	"iso4217":                       "the field '%s' must be a valid ISO 4217 currency code",
	"json":                          "the field '%s' must be valid JSON",
	"jwt":                           "the field '%s' must be a valid JSON Web Token (JWT)",
	"latitude":                      "the field '%s' must be a valid latitude",
	"longitude":                     "the field '%s' must be a valid longitude",
	"luhn_checksum":                 "the field '%s' must pass the Luhn algorithm checksum",
	"postcode_iso3166_alpha2":       "the field '%s' must be a valid postcode for the country '%s'",
	"postcode_iso3166_alpha2_field": "the field '%s' must be a valid postcode for the country in field '%s'",
	"rgb":                           "the field '%s' must be a valid RGB color",
	"rgba":                          "the field '%s' must be a valid RGBA color",
	"ssn":                           "the field '%s' must be a valid Social Security Number",
	"timezone":                      "the field '%s' must be a valid timezone",
	"uuid":                          "the field '%s' must be a valid UUID",
	"uuid3":                         "the field '%s' must be a valid UUID v3",
	"uuid3_rfc4122":                 "the field '%s' must be a valid UUID v3 RFC4122",
	"uuid4":                         "the field '%s' must be a valid UUID v4",
	"uuid4_rfc4122":                 "the field '%s' must be a valid UUID v4 RFC4122",
	"uuid5":                         "the field '%s' must be a valid UUID v5",
	"uuid5_rfc4122":                 "the field '%s' must be a valid UUID v5 RFC4122",
	"uuid_rfc4122":                  "the field '%s' must be a valid UUID RFC4122",
	"md4":                           "the field '%s' must be a valid MD4 hash",
	"md5":                           "the field '%s' must be a valid MD5 hash",
	"sha256":                        "the field '%s' must be a valid SHA256 hash",
	"sha384":                        "the field '%s' must be a valid SHA384 hash",
	"sha512":                        "the field '%s' must be a valid SHA512 hash",
	"ripemd128":                     "the field '%s' must be a valid RIPEMD-128 hash",
	"ripemd160":                     "the field '%s' must be a valid RIPEMD-160 hash",
	"tiger128":                      "the field '%s' must be a valid TIGER128 hash",
	"tiger160":                      "the field '%s' must be a valid TIGER160 hash",
	"tiger192":                      "the field '%s' must be a valid TIGER192 hash",
	"semver":                        "the field '%s' must be a valid semantic version",
	"ulid":                          "the field '%s' must be a valid ULID",
	"cve":                           "the field '%s' must be a valid CVE identifier",

	// Numeric/string comparisons
	"eq":             "the field '%s' must be equal to '%s'",
	"eq_ignore_case": "the field '%s' must be equal to '%s' ignoring case",
	"gt":             "the field '%s' must be greater than %s",
	"gte":            "the field '%s' must be greater than or equal to %s",
	"lt":             "the field '%s' must be less than %s",
	"lte":            "the field '%s' must be less than or equal to %s",
	"ne":             "the field '%s' must not be equal to '%s'",
	"ne_ignore_case": "the field '%s' must not be equal to '%s' ignoring case",

	// Length/size validations for strings
	"len": "the field '%s' must have a length of %s",
	"max": "the field '%s' must be no more than %s",
	"min": "the field '%s' must have a length of %s",

	// Other validations
	"oneof":                "the field '%s' must be one of the following values: %s",
	"required":             "the field '%s' is required",
	"required_if":          "the field '%s' is required when %s",
	"required_unless":      "the field '%s' is required unless %s",
	"required_with":        "the field '%s' is required when %s is present",
	"required_with_all":    "the field '%s' is required when all of %s are present",
	"required_without":     "the field '%s' is required when %s is not present",
	"required_without_all": "the field '%s' is required when none of %s are present",
	"excluded_if":          "the field '%s' is excluded when %s",
	"excluded_unless":      "the field '%s' is excluded unless %s",
	"excluded_with":        "the field '%s' is excluded when %s is present",
	"excluded_with_all":    "the field '%s' is excluded when all of %s are present",
	"excluded_without":     "the field '%s' is excluded when %s is not present",
	"excluded_without_all": "the field '%s' is excluded when none of %s are present",
	"unique":               "the field '%s' must contain unique values",
	"validateFn":           "the field '%s' failed custom validation",

	// Aliases
	"iscolor":      "the field '%s' must be a valid color (hex, rgb, rgba, hsl, or hsla)",
	"country_code": "the field '%s' must be a valid country code (ISO 3166-1 alpha-2, alpha-3, or numeric)",
}

// getJSONFieldName returns the JSON field name for a struct field
func getJSONFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return field.Name
	}
	parts := strings.Split(jsonTag, ",")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return field.Name
}

func formatErrorMessage(err validator.FieldError, jsonFieldName string) string {
	tag := err.Tag()
	param := err.Param()

	if messageFormat, exists := errorMessages[tag]; exists {
		var message string
		switch tag {
		case "required", "dir", "file", "image", "isdefault", "unique", "validateFn",
			"lowercase", "uppercase", "multibyte", "ascii", "printascii",
			"boolean", "alpha", "alphanum", "alphanumunicode", "alphaunicode",
			"ip", "ipv4", "ipv6", "email", "url", "uri", "http_url", "uuid",
			"json", "credit_card", "isbn", "isbn10", "isbn13", "issn", "ssn",
			"hexadecimal", "hexcolor", "rgb", "rgba", "hsl", "hsla", "iscolor",
			"country_code", "timezone", "latitude", "longitude",
			"cidr", "hostname", "fqdn", "mac", "base64", "jwt", "html", "html_encoded",
			"number", "numeric", "semver", "ulid", "cve", "ein", "btc_addr",
			"btc_addr_bech32", "eth_addr", "bic", "bcp47_language_tag", "mongodb",
			"mongodb_connection_string", "cron", "spicedb", "luhn_checksum", "e164":
			message = fmt.Sprintf(messageFormat, jsonFieldName)

		case "eq", "ne", "gt", "gte", "lt", "lte", "eq_ignore_case", "ne_ignore_case":
			message = fmt.Sprintf(messageFormat, jsonFieldName, param)

		case "len":
			message = fmt.Sprintf(messageFormat, jsonFieldName, param)

		case "max", "min":
			if err.Kind() == reflect.String {
				message = fmt.Sprintf(errorMessages["len"], jsonFieldName, param)
			} else {
				if tag == "min" {
					message = fmt.Sprintf("the field '%s' must be at least %s", jsonFieldName, param)
				} else {
					message = fmt.Sprintf("the field '%s' must be no more than %s", jsonFieldName, param)
				}
			}

		case "oneof":
			message = fmt.Sprintf(messageFormat, jsonFieldName, param)

		case "eqfield", "nefield", "gtfield", "gtefield", "ltfield", "ltefield",
			"eqcsfield", "necsfield", "gtcsfield", "gtecsfield", "ltcsfield", "ltecsfield":
			message = fmt.Sprintf(messageFormat, jsonFieldName, param)

		case "contains", "containsany", "containsrune", "excludes", "excludesall",
			"excludesrune", "endswith", "endsnotwith", "startswith", "startsnotwith":
			message = fmt.Sprintf(messageFormat, jsonFieldName, param)

		case "required_if", "required_unless", "required_with", "required_with_all",
			"required_without", "required_without_all", "excluded_if", "excluded_unless",
			"excluded_with", "excluded_with_all", "excluded_without", "excluded_without_all":
			message = fmt.Sprintf(messageFormat, jsonFieldName, param)

		case "postcode_iso3166_alpha2", "postcode_iso3166_alpha2_field":
			message = fmt.Sprintf(messageFormat, jsonFieldName, param)

		case "datetime":
			message = fmt.Sprintf(messageFormat, jsonFieldName, param)

		default:
			message = fmt.Sprintf(messageFormat, jsonFieldName)
		}
		return message
	}
	return err.Error()
}

func deepValidator(s interface{}, validate *validator.Validate) map[string]string {
	errors := make(map[string]string)
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	// Validate the struct itself
	var errs validator.ValidationErrors
	err := validate.Struct(s)
	if err != nil {
		errs = err.(validator.ValidationErrors)
		for _, e := range errs {
			fieldName := e.StructField()
			if f, ok := typ.FieldByName(fieldName); ok {
				errors[getJSONFieldName(f)] = formatErrorMessage(e, getJSONFieldName(f))
			}
		}
	}

	// Process struct fields recursively
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		jsonFieldName := getJSONFieldName(fieldType)

		switch field.Kind() {
		case reflect.Struct:
			nestedErrors := deepValidator(field.Interface(), validate)
			for k, v := range nestedErrors {
				errors[fmt.Sprintf("%s.%s", jsonFieldName, k)] = v
			}
		case reflect.Ptr:
			if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
				nestedErrors := deepValidator(field.Interface(), validate)
				for k, v := range nestedErrors {
					errors[fmt.Sprintf("%s.%s", jsonFieldName, k)] = v
				}
			}
		case reflect.Slice, reflect.Array:
			// Only process if the slice/array contains structs
			if field.Type().Elem().Kind() == reflect.Struct {
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j)
					nestedErrors := deepValidator(elem.Interface(), validate)
					for k, v := range nestedErrors {
						errors[fmt.Sprintf("%s[%d].%s", jsonFieldName, j, k)] = v
					}
				}
			}
		}
	}

	return errors
}

// ValidateStruct validates a struct and returns JSON field names in error messages
func ValidateStruct(s interface{}) map[string]string {
	if s == nil {
		return nil
	}
	validate := validator.New()
	return deepValidator(s, validate)
}

// ValidateStructDeep validates a struct with deep nested validation using JSON field names
func ValidateStructDeep(s interface{}) map[string]string {
	return ValidateStruct(s)
}

// ValidateSliceDeep validates a slice/array with deep nested validation using JSON field names
func ValidateSliceDeep(slice interface{}) map[string]string {
	if slice == nil {
		return nil
	}
	validate := validator.New()
	errors := make(map[string]string)
	val := reflect.ValueOf(slice)

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return map[string]string{"root": "provided value is neither a struct nor a slice/array"}
	}

	// Only process slices/arrays that contain structs
	elemType := val.Type().Elem()
	if elemType.Kind() != reflect.Struct && 
	   (elemType.Kind() == reflect.Ptr && elemType.Elem().Kind() != reflect.Struct) {
		return nil // Don't validate slices of non-struct types
	}

	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		if elem.Kind() == reflect.Struct || 
		   (elem.Kind() == reflect.Ptr && !elem.IsNil() && elem.Elem().Kind() == reflect.Struct) {
			nestedErrors := deepValidator(elem.Interface(), validate)
			for k, v := range nestedErrors {
				errors[fmt.Sprintf("[%d].%s", i, k)] = v
			}
		}
	}
	return errors
}

// ValidateAny validates any input (struct, slice, or array) with deep nesting using JSON field names
func ValidateAny(input interface{}) map[string]string {
	if input == nil {
		return nil
	}
	v := reflect.ValueOf(input)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		return ValidateStructDeep(input)
	case reflect.Slice, reflect.Array:
		return ValidateSliceDeep(input)
	default:
		return map[string]string{"root": "provided value is neither a struct nor a slice/array"}
	}
}