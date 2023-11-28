package manipulator

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

// CheckPasswordComplexity ...
func CheckPasswordComplexity(ps string) error {
	if len(ps) < 8 {
		return fmt.Errorf("password should contains 8 characters ")
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	//symbol := `[!@#~$%^&*()+|_]{1}`

	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		return fmt.Errorf("password need number :%v", err)
	}
	if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
		return fmt.Errorf("password need lower case :%v", err)
	}
	if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
		return fmt.Errorf("password need Upper case :%v", err)
	}
	// if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
	// 	return fmt.Errorf("password need symbol :%v", err)
	// }

	return nil
}

// FormatPhonenumber ...
func FormatPhonenumber(phonenumber string) string {
	if strings.HasPrefix(phonenumber, `+`) {
		phonenumber = phonenumber[1:]
	}
	if strings.HasPrefix(phonenumber, `0`) {
		phonenumber = `62` + phonenumber[1:]
	}
	if strings.HasPrefix(phonenumber, `8`) {
		phonenumber = `62` + phonenumber
	}
	return phonenumber
}

func PointerStringToStringValue(valueNil *string) string {
	if valueNil == nil {
		return ``
	}
	return *valueNil
}

// ParseTemplate
func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	template, err := template.ParseFiles(wd + "/templates/" + templateFileName)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	err = template.Execute(buffer, data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}