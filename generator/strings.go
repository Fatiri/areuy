package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mathRand "math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	lowerCharSet   = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
	numberCharSet  = "0123456789"
)

func RandomString(
	passwordLength int,
	noUpperChar, noSpecialChar, noNumberChar, allowRepeat bool) (string, error) {
	var generatedRandomString string
	letters := lowerCharSet

	if !noUpperChar {
		letters += upperCharSet
	}

	if !noSpecialChar {
		letters += specialCharSet
	}

	if !noNumberChar {
		letters += numberCharSet
	}

	for i := 0; i < passwordLength; i++ {
		char, err := randomElement(letters)
		if err != nil {
			return "", err
		}

		if !allowRepeat && strings.Contains(generatedRandomString, char) {
			i--
			continue
		}

		generatedRandomString, err = randomInsert(generatedRandomString, char)
		if err != nil {
			return "", err
		}
	}

	return generatedRandomString, nil
}

func GenerateStringByDate() string {
	year := strconv.Itoa(time.Now().Year())
	month := strconv.Itoa(int(time.Now().Month()))
	day := strconv.Itoa(time.Now().Day())

	return year + month + day
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand *mathRand.Rand = mathRand.New(
	mathRand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateRandomStringUpperCase(prefix string, length int) string {
	return fmt.Sprintf("%s%s", prefix, StringWithCharset(length, charset))
}

func GenerateUniqCode(prefix string) string {
	prefixTimeStamp := time.Now().UnixNano() / int64(time.Millisecond)
	code := fmt.Sprintf(`%v`, prefixTimeStamp)
	code = strings.ReplaceAll(code, `-`, ``)
	code = prefix + code
	return code
}

func GenerateRandomCode(walletID int, code string) string {
	//generate increment wallet id 6 digits
	walletIdString := strconv.Itoa(int(walletID))
	for i := 0; i < 10-len(walletIdString); i++ {
		walletIdString = "0" + walletIdString
	}

	//generate 2 digits month string
	year, month, _ := time.Now().Date()
	monthString := strconv.Itoa(int(month))

	//generate 2 digits year string
	yearSting := strconv.Itoa(year)[2:]
	if month < 10 {
		monthString = "0" + monthString
	}
	return code + "-" + monthString + yearSting + "-" + walletIdString
}

func randomInsert(s, val string) (string, error) {
	if s == "" {
		return val, nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(s)+1)))
	if err != nil {
		return "", err
	}

	i := n.Int64()
	return s[0:i] + val + s[i:], nil
}

func randomElement(s string) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(s))))
	if err != nil {
		return "", err
	}
	return string(s[n.Int64()]), nil
}
