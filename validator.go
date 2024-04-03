package validator

import (
	"encoding/base64"
	"reflect"
	"regexp"
	"strings"
)

type VALIDATE_FUNC func(value reflect.Value) (string, bool)

// We only compare the magic number of the given data
var FILE_SIGNATURES = map[string][]byte{
	"png": {0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
	"jpg": {0xFF, 0xD8, 0xFF, 0xE0},
}

/* The first index of the returned item is always
 * the tag name of the `json` tag
 */
func Validate(typedStruct interface{}) ([][]string, bool) {
	errMessages := [][]string{}
	valueOf := reflect.ValueOf(typedStruct)
	typeOf := reflect.TypeOf(typedStruct)
	numFields := typeOf.NumField()

	for i := 0; i < numFields; i++ {
		field := typeOf.Field(i)
		value := reflect.Indirect(valueOf).FieldByIndex(field.Index)
		tagName, hasJson := field.Tag.Lookup("json")
		errs := []string{}

		if !hasJson {
			// Field is skipped due to the abscence of the json tag
			continue
		}

		if tag, ok := field.Tag.Lookup("validate"); ok {
			// iterate over the rules separated by "|"
			for _, rule := range strings.Split(tag, "|") {
				if rule == "required" {
					validate(validateRequired, value, &errs)
				} else if rule == "email" {
					validate(validateEmail, value, &errs)
				} else if strings.HasPrefix(rule, "ext:") {
                    // No need to further check if the data isn't a valid base64
                    decodedData, err := base64.StdEncoding.DecodeString(value.Interface().(string))

                    if err != nil {
                        errs = append(errs, "invalid base64 data")
                        break
                    }

					extensionString := strings.ToLower(string(rule[len("ext:"):]))
					extensions := strings.Split(extensionString, ",")

					tempErr := []string{}
					isValid := false

					for i := 0; i < len(extensions); i++ {
						/**
						 * It only takes one success validation to returns true
						 * and ignore the rest of the errors
						 */
						extension := string(extensions[i])
						if _, ok := FILE_SIGNATURES[extension]; !ok {
							// Unsupported file extension
							continue
						}

						errMsg, ok := validateFileFormat(reflect.ValueOf(decodedData), extension)

						if ok {
							isValid = true
							break
						}

						tempErr = append(tempErr, errMsg)
					}

                    // in case there's no extension/valid extension supplied?
					if !isValid && len(tempErr) > 0 {
						errs = append(errs, tempErr...)
					}
				}
			}
		}

		if len(errs) > 0 {
			errs = append([]string{tagName}, errs...)
			errMessages = append(errMessages, errs)
		}
	}

	return errMessages, len(errMessages) == 0
}

func validate(fn VALIDATE_FUNC, rule reflect.Value, errs *[]string) {
	if msg, ok := fn(rule); !ok {
		*errs = append(*errs, msg)
	}
}

func validateRequired(value reflect.Value) (string, bool) {
	kind := value.Kind()

	if (kind == reflect.Pointer && value.IsNil()) ||
		(kind == reflect.String && value.Interface().(string) == "") {
		return "field is empty", false
	}

	return "", true
}

func validateEmail(value reflect.Value) (string, bool) {
	if value.Kind() != reflect.String {
		return "invalid email format", false
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !re.MatchString(value.Interface().(string)) {
		return "invalid email format", false
	}

	return "", true
}

func validateFileFormat(value reflect.Value, format string) (string, bool) {
	if value.Kind() != reflect.Slice {
		// Should be a slice of bytes of the actual content
		return "invalid file data", false
	}

    data := value.Interface().([]byte)

	signature := FILE_SIGNATURES[format]

	if len(data) < len(signature) {
		return "invalid " + format + " file", false
	}

	for i := range signature {
		if data[i] != signature[i] {
			// should we reduce more the information given?
			return "invalid " + format + " file", false
		}
	}

	return "", true
}
