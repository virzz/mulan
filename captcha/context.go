package captcha

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"sync"
	"time"

	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var (
	captchaMap = cmap.New[*Data]()
	once       sync.Once
	std        *Captcha
)

func Init() {
	std, _ = New(defaultConfig)
	go once.Do(func() {
		for {
			time.Sleep(time.Minute)
			captchaMap.IterCb(func(key string, value *Data) {
				if value.Expire <= time.Now().Unix() {
					captchaMap.Remove(key)
				}
			})
		}
	})
}

// Check 验证验证码是否正确,返回错误类型
func Check(id, code string) (bool, error) {
	if data, ok := captchaMap.Get(id); ok {
		if data.Expire <= time.Now().Unix() {
			captchaMap.Remove(id)
			return false, ErrExpired
		}
		if data.Result == code {
			captchaMap.Remove(id)
			return true, nil
		}
	}
	return false, ErrInvalid
}

// Check 验证验证码是否正确
func CheckOk(id, code string) (ok bool) {
	ok, _ = Check(id, code)
	return
}

func create() (id, result string, img image.Image) {
	if std == nil {
		panic("plz init")
	}
	img, data := std.Draw()
	id = uuid.New().String()
	data.Expire = time.Now().Add(time.Minute * 5).Unix()
	captchaMap.Set(id, data)
	if r := recover(); r != nil {
		return create()
	}
	return id, data.Result, img
}

func createBytes() (id, result string, buf []byte, err error) {
	id, result, img := create()
	_buf := new(bytes.Buffer)
	err = png.Encode(_buf, img)
	if err != nil {
		return "", "", nil, err
	}
	return id, result, _buf.Bytes(), nil
}

func createB64() (id, result, data string, err error) {
	id, result, buf, err := createBytes()
	if err != nil {
		return "", "", "", err
	}
	data = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf)
	return id, result, data, nil
}

// Create 创建验证码，返回图片对象image.Image
func Create() (id, result string, img image.Image) { return create() }

// CreateBytes 创建验证码，返回图片数据[]byte
func CreateBytes() (id, result string, buf []byte, err error) { return createBytes() }

// CreateB64 创建验证码，返回图片base64
func CreateB64() (id, result string, data string, err error) { return createB64() }
