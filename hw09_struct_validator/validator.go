package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	lenRule   = "len"
	inRule    = "in"
	regexRule = "regexp"
	minRule   = "min"
	maxRule   = "max"
)

func Validate(v interface{}) error {
	r := reflect.ValueOf(v)
	if r.Kind() != reflect.Struct {
		return ErrInvalidArg
	}

	var errs ValidationErrors

	for _, f := range reflect.VisibleFields(r.Type()) {
		err := validateField(f, r, &errs)
		if err != nil {
			return err
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateField(f reflect.StructField, r reflect.Value, errs *ValidationErrors) error {
	validateTag := f.Tag.Get("validate")
	if validateTag == "" {
		return nil
	}

	validationRules, err := extractRules(validateTag)
	if err != nil {
		return err
	}

	fieldValue := r.FieldByIndex(f.Index)

	switch fieldKind := f.Type.Kind(); fieldKind { //nolint:exhaustive
	case reflect.String:
		err := validateString(validationRules, fieldValue.String())
		if errors.Is(err, ErrValidateField) {
			*errs = append(*errs, ValidationError{
				Field: f.Name,
				Err:   err,
			})
		} else if err != nil {
			return err
		}
		return nil
	case reflect.Int:
		err := validateInt(validationRules, int(fieldValue.Int()))
		if errors.Is(err, ErrValidateField) {
			*errs = append(*errs, ValidationError{
				Field: f.Name,
				Err:   err,
			})
		} else if err != nil {
			return err
		}
		return nil
	case reflect.Slice:
		err := validateSliceField(f, fieldValue, validationRules, errs)
		if err != nil {
			return err
		}
	default:
		return nil
	}

	return nil
}

func extractRules(validateTag string) (map[string]string, error) {
	rules := make(map[string]string)
	for _, rule := range strings.Split(validateTag, "|") {
		nameVal := strings.Split(rule, ":")
		if len(nameVal) != 2 {
			return nil, ErrInvalidRule
		}
		rules[nameVal[0]] = nameVal[1]
	}
	return rules, nil
}

func validateString(rules map[string]string, value string) error {
	for ruleName, ruleVal := range rules {
		switch ruleName {
		case lenRule:
			requiredLen, err := strconv.Atoi(ruleVal)
			if err != nil {
				return err
			}
			if len(value) != requiredLen {
				return NewErrValidateLen(requiredLen)
			}
		case regexRule:
			regexp, err := regexp.Compile(ruleVal)
			if err != nil {
				return err
			}
			if !regexp.MatchString(value) {
				return NewErrValidateRegex(ruleVal)
			}
		case inRule:
			allowedValues := strings.Split(ruleVal, ",")
			valueInAllowed := func(val string, allowed []string) bool {
				for _, av := range allowed {
					if val == av {
						return true
					}
				}
				return false
			}(value, allowedValues)

			if !valueInAllowed {
				return NewErrValidateIn(ruleVal)
			}
		default:
			return NewErrUnknnownRule(ruleName)
		}
	}
	return nil
}

func validateInt(rules map[string]string, value int) error {
	for ruleName, ruleValue := range rules {
		switch ruleName {
		case minRule:
			minValue, err := strconv.Atoi(ruleValue)
			if err != nil {
				return ErrInvalidRuleIntVal
			}
			if value < minValue {
				return NewErrValidateMin(minValue)
			}
		case maxRule:
			maxValue, err := strconv.Atoi(ruleValue)
			if err != nil {
				return ErrInvalidRuleIntVal
			}
			if value > maxValue {
				return NewErrValidateMax(maxValue)
			}
		case inRule:
			allowedValues := strings.Split(ruleValue, ",")
			stringVal := strconv.Itoa(value)
			valueInAllowed := func(val string, allowed []string) bool {
				for _, av := range allowed {
					if val == av {
						return true
					}
				}
				return false
			}(stringVal, allowedValues)

			if !valueInAllowed {
				return NewErrValidateIn(ruleValue)
			}
		default:
			return NewErrUnknnownRule(ruleName)
		}
	}
	return nil
}

func validateSliceField(
	f reflect.StructField,
	fieldValue reflect.Value,
	validationRules map[string]string,
	errs *ValidationErrors,
) error {
	switch elemKind := f.Type.Elem().Kind(); elemKind { //nolint:exhaustive
	case reflect.String:
		elems := fieldValue.Interface().([]string)
		for i, val := range elems {
			err := validateString(validationRules, val)
			if errors.Is(err, ErrValidateField) {
				*errs = append(*errs, ValidationError{
					Field: fmt.Sprintf("%s[%d]", f.Name, i),
					Err:   err,
				})
			} else if err != nil {
				return err
			}
		}
	case reflect.Int:
		elems := fieldValue.Interface().([]int)
		for i, val := range elems {
			err := validateInt(validationRules, val)
			if errors.Is(err, ErrValidateField) {
				*errs = append(*errs, ValidationError{
					Field: fmt.Sprintf("%s[%d]", f.Name, i),
					Err:   err,
				})
			} else if err != nil {
				return err
			}
		}
	default:
		return nil
	}

	return nil
}
