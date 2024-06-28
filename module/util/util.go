package util

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

var src = rand.NewSource(time.Now().UnixNano())

func IsRoot() bool {
	return os.Geteuid() == 0
}

func GetRandomString(n int) string {
	allBytes := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	letterBytes := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	result := new(bytes.Buffer)
	result.WriteByte(allBytes[rand.Intn(len(letterBytes))])
	for i := 0; i < n-1; i++ {
		result.WriteByte(allBytes[rand.Intn(len(allBytes))])
	}
	return result.String()
}

func GetRandomLowerString(n int) string {
	allBytes := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	result := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		result = append(result, allBytes[rand.Intn(len(allBytes))])
	}
	return string(result)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func Init() {
	timelocal, _ := time.LoadLocation("Asia/Chongqing")
	time.Local = timelocal
}

func MD5(str string) string {
	c := md5.New()
	c.Write([]byte(str))
	return hex.EncodeToString(c.Sum(nil))
}

// BytesToString 提供快速的 bytes 转 string 的方法，
// 不需要拷贝内存，因此比直接转换速度快，
// Do NOT use this function unless you know what you are doing
func BytesToString(bytes []byte) (s string) {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	stringHeader.Data = sliceHeader.Data
	stringHeader.Len = sliceHeader.Len
	return
}

// StringToBytes 提供快速的 string 转 bytes 的方法，不需要拷贝内存
// 转换后的 bytes 不可修改，因为 golang 中 string 是 immutable 的
// Do NOT use this function unless you know what you are doing
func StringToBytes(s string) (b []byte) {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sliceHeader.Data = stringHeader.Data
	sliceHeader.Len = stringHeader.Len
	sliceHeader.Cap = stringHeader.Len
	return
}

// Base64EncodeString 计算字符串的base64值
func Base64EncodeString(in string) string {
	return base64.StdEncoding.EncodeToString(StringToBytes(in))
}

// Base64EncodeBytes 计算Bytes的base64值
func Base64EncodeBytes(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

// BaseDecodeToBytes base64解码为bytes
// 忽略了报错
func Base64DecodeToBytes(in string) []byte {
	out, _ := base64.StdEncoding.DecodeString(in)
	return out
}

// Base64DecodeToString base64解码为string
func Base64DecodeToString(in string) string {
	out, _ := base64.StdEncoding.DecodeString(in)
	return string(out)
}

// Base64DecodeBytesToString base64解码bytes为string
// 为了兼容xray格式增加的，没什么实际用处
func Base64DecodeBytesToString(in []byte) string {
	out, _ := base64.StdEncoding.DecodeString(BytesToString(in))
	return string(out)
}

func RandomInt(from, to int) int {
	if from > to {
		to, from = from, to
	}
	mod := to - from
	mod += 1
	return rand.New(src).Int()%mod + from
}

func RandomInt64(from, to int64) int64 {
	if from > to {
		to, from = from, to
	}
	mod := to - from
	mod += 1
	return rand.New(src).Int63()%mod + from
}

func GetCurrentTimeStamp() int64 {
	return time.Now().UnixNano() / 1e6
}

func ByteToAscii(buff []byte) string {
	result := new(bytes.Buffer)
	for _, c := range buff {
		if c <= 31 || c >= 127 {
			result.WriteString(fmt.Sprintf("\\x%02x", c))
			continue
		}
		result.WriteByte(c)
	}
	return result.String()
}

func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
