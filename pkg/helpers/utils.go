package helpers

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

func GeneratePassword(password string) string {
	sha1pwd := sha1.Sum([]byte(strings.ToUpper(password)))
	return strings.ToUpper(hex.EncodeToString(sha1pwd[:]))
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToQrBase64(data string) string {
	qrImgData, err := qrcode.Encode(data, qrcode.High, 256)
	if err != nil {
		return ""
	}
	base64 := base64.StdEncoding.EncodeToString(qrImgData)
	return base64
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandSeq :
func RandSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Mask : mask string
func Mask(str string, ts ...int) string {
	result := ""
	if ts[0] == 1 {
		i := 0
		for _, c := range str {
			if i%2 == 1 {
				result += "*"
			} else {
				result += string(c)
			}
			i++
		}
		return result
	}
	return str
}

func StringInArr(a string, list []string) bool {
	if len(list) == 0 {
		return false
	}

	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GenerateOtp() int {
	k, _ := strconv.Atoi(os.Getenv("OTP_CODE_LENGTH"))
	min := int(math.Pow10(k - 1))
	max := int(math.Pow10(k) - 1)
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func GenerateHashDocument(first_name, last_name, birth_date, gender, country_code, passport_number string) string {
	hash := md5.Sum([]byte(strings.ToLower(first_name + last_name + gender + birth_date + country_code + passport_number)))
	return hex.EncodeToString(hash[:])
}

func UniqueIntSlice(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
