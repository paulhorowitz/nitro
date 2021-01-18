package validate

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Validator interface {
	Validate(input string) error
}

// HostnameValidator is used to validate a provided hostname
type HostnameValidator struct{}

func (v *HostnameValidator) Validate(input string) error {
	// check length
	if len(input) < 3 {
		return fmt.Errorf("hostname must be more than 3 characters")
	}

	// check for spaces
	if strings.Contains(input, " ") {
		return fmt.Errorf("hostname must not include spaces")
	}

	// check for special characters
	if strings.ContainsAny(input, "!@#$%^&*()") {
		return fmt.Errorf("hostname must not include any special characters")
	}

	return nil
}

type PHPVersionValidator struct{}

func (v *PHPVersionValidator) Validate(input string) error {
	switch input {
	case "8.0", "7.4", "7.3", "7.2", "7.1", "7.0":
		return nil
	}

	return fmt.Errorf("the PHP inputrsion %q is not valid", input)
}

func Hostname(v string) error {
	msg := "you must provide a valid domain, without a TLD and only lowercase"

	if strings.Contains(v, " ") {
		return errors.New(msg)
	}

	lower := strings.ToLower(v)
	if lower != v {
		return errors.New(msg)
	}

	return nil
}

// path will check is a fail
func Path(p string) error {
	f, err := os.Stat(p)
	if err != nil {

		return err
	}

	if f.IsDir() {
		return nil
	}

	return errors.New("you must provide a path to a valid directory")
}

// PHPVersion takes a string that represents a PHP version to install and returns and error if that PHP version
// is not yet supported.
func PHPVersion(v string) error {
	switch v {
	case "7.4":
		return nil
	case "7.3":
		return nil
	case "7.2":
		return nil
	}

	return fmt.Errorf("the PHP version %q is not valid", v)
}

func Memory(v string) error {
	if !strings.HasSuffix(v, "G") {
		return errors.New("memory must end with a G")
	}

	return nil
}

func DiskSize(v string) error {
	if !strings.HasSuffix(v, "G") {
		return errors.New("disk must end with a G")
	}

	return nil
}

func MachineName(v string) error {
	if v == "" {
		return errors.New("machine name cannot be empty")
	}
	if strings.Contains(v, " ") {
		return errors.New("machine name cannot contain spaces")
	}

	return nil
}

func MaxExecutionTime(v string) error {
	_, err := strconv.Atoi(v)
	if err != nil {
		return errors.New("max_execution_time must be a valid integer")
	}

	return nil
}

func MaxInputVars(v string) error {
	num, err := strconv.Atoi(v)
	if err != nil {
		return errors.New("max_input_vars must be a valid integer")
	}

	if num >= 10000 {
		return errors.New("max_input_vars must be less than 10000")
	}

	return nil
}

func IsMegabytes(v string) error {
	if len(v) == 1 {
		return errors.New("memory must be larger than 1 character (e.g. 256M)")
	}

	if !strings.HasSuffix(v, "M") {
		return errors.New("memory must end with a M")
	}

	return nil
}

func PhpMaxFileUploads(v string) error {
	num, err := strconv.Atoi(v)
	if err != nil {
		return errors.New("max_input_vars must be a valid integer")
	}

	if num >= 500 {
		return errors.New("max_file_uploads must be less than 500")
	}

	return nil
}
