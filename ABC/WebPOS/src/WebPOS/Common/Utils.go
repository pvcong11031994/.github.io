package Common

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"strconv"
	"github.com/goframework/gf/ext"
	"os"
	"WebPOS/WebApp"
	"github.com/goframework/gf"
	"runtime"
)

// https://dev-backlog.rsp.honto.jp/backlog/alias/wiki/54805
// 7.2
func GeneratePass(user_id, pass string) string {
	cryptSha512 := sha512.Sum512([]byte(user_id + pass))
	return hex.EncodeToString(cryptSha512[:])
}

// Generate secret pass with md5 encrypt
func GeneratePassMd5(user_id, pass string) string {
	hash := md5.New()
	io.WriteString(hash, pass)

	pwmd5 := fmt.Sprintf("%x", hash.Sum(nil))

	io.WriteString(hash, SHA_PRIVATE_KEY_1)
	io.WriteString(hash, user_id)
	io.WriteString(hash, SHA_PRIVATE_KEY_2)
	io.WriteString(hash, pwmd5)
	encryptPass := fmt.Sprintf("%x", hash.Sum(nil))

	return encryptPass
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

func StringArrayEquals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func PercentDiv(a, b int) int {
	if b != 0 {
		percent := float64(a) / float64(b) * float64(100)
		return int(percent + 0.5) // Round
	} else {
		return 0
	}
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numberRunes = []rune("0123456789")

// Create random string with length
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Create random string number with length
func RandNumber(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numberRunes[rand.Intn(len(numberRunes))]
	}
	return string(b)
}

// Create random string with other string
func GenerateString(input string) string {
	cryptSha512 := sha512.Sum512([]byte(input))
	return hex.EncodeToString(cryptSha512[:])
}
// Write log debug error
func LogErr(err error) {
	if err != nil {
		//log.Println(err)
		log.Println(err.Error())
	}
}

func ToInterfaceArray(arrStr []string) []interface{} {
	b := make([]interface{}, len(arrStr))
	for i := range arrStr {
		b[i] = arrStr[i]
	}
	return b
}

//Create list sql parameter ?
func SQLPara(arrStr []string) string {
	return strings.Trim(strings.Repeat("?,", len(arrStr)), ",")
}

// Convert link to <a> tag html
func RenderLink(content string) template.HTML {
	content = template.HTMLEscapeString(content)
	linkPat := regexp.MustCompile(`(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`)
	content = linkPat.ReplaceAllStringFunc(content, func(link string) string {
		content = content[1 : len(content)-1]
		return `<a href="` + link + `" target="blank">` + link + `</a>`
	})
	return template.HTML(strings.Replace(content, "\n", "<br/>", -1))
}

func IF(flag bool, trueVal, falseVal interface{}) interface{} {
	if flag {
		return trueVal
	}
	return falseVal
}

//Size in bytes to KB/MB/GB/TB
func HumanReadableFileSize(byteSize int64) string {
	if byteSize > 1024*1024*1024*1024 {
		tbSize := byteSize / (1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%dTB", tbSize)
	}
	if byteSize > 1024*1024*1024 {
		gbSize := byteSize / (1024 * 1024 * 1024)
		return fmt.Sprintf("%dGB", gbSize)
	}
	if byteSize > 1024*1024 {
		mbSize := byteSize / (1024 * 1024)
		return fmt.Sprintf("%dMB", mbSize)
	}
	if byteSize > 1024 {
		kbSize := byteSize / (1024)
		return fmt.Sprintf("%dKB", kbSize)
	}

	return fmt.Sprintf("%dB", byteSize)
}

func GenerateJAN(code string) string {

	//前後の空白をトリムする
	code = strings.TrimSpace(code)
	//全半角ハイフンを取り除く
	code = strings.Replace(code, "-","",-1)

	if len(code) == 8 {
		code = "00000" + code //短縮JAN
	} else if len(code) == 9 || len(code) == 10 {
		code = "978" + code //ISBN
	} else if len(code) == 12 {
		//書籍JAN・雑誌JANと異なる場合
		if !strings.HasPrefix(code, "978") || !strings.HasPrefix(code, "491") {
			code = "0" + code //UPC
			return code[:len(code)-1] + GetCheckDigit(code)
		}
	}

	//その以外の場合、そのまま戻る
	if len(code) != 13 {
		return code
	}

	checkDigit := GetCheckDigit(code)

	code = code[:len(code)-1] + checkDigit

	return code
}

func GetCheckDigit(code string) string {

	odd := 0
	even := 0
	for k := 0; k < len(code)-1; k++ {
		num, err := strconv.Atoi(string(code[k]))
		if err != nil {
			return ""
		}

		if k % 2 == 0{
			even += num
		} else {
			odd += num
		}
	}

	return strconv.Itoa((10 - (odd * 3 + even) % 10) % 10)
}

// Ouput log to log file
func LogOutput(logContents string) {

	if runtime.GOOS != _WINDOWS_OS {
		logFilePath := _PATH_LOG + "application.log"

		if !ext.FolderExists(_PATH_LOG) {
			os.MkdirAll(_PATH_LOG, os.ModePerm)
		}

		f, _ := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		log.SetOutput(f)

		log.Println(fmt.Sprintf("%v", logContents))
	}
}

// Ouput log query to log file
func LogSQL(ctx *gf.Context, query string) {

	logLevel := ctx.Config.Int(WebApp.CONFIG_KEY_SERVER_LOG_LEVEL, 0)
	if logLevel == 1 {
		LogOutput(fmt.Sprintf("%v \n %v", ctx.UrlPath, query))
	}
}

// Join concatenates the elements of a to create a single string. The separator string
// sep is placed between elements in the resulting string.
// option add character first and end of item in array
func JoinArray(a []string, first, last, sep string) string {
	tempArr := []string{}
	for _, item := range a {
		tempArr = append(tempArr, first + item + last)
	}
	a = tempArr
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	case 2:
		// Special case for common small values.
		// Remove if golang.org/issue/6714 is fixed
		return a[0] + sep + a[1]
	case 3:
		// Special case for common small values.
		// Remove if golang.org/issue/6714 is fixed
		return a[0] + sep + a[1] + sep + a[2]
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	b := make([]byte, n)
	bp := copy(b, a[0])
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}
	return string(b)
}