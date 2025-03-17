package captcha

import (
	_ "embed"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"math/rand"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

//go:embed monaco_ascii.ttf
var fontData []byte

type (
	Config struct {
		Width, Height                    int         // width and height of the captcha image
		Bg, Ft                           color.Color // background and font color
		Face                             font.Face   // font face
		Size                             float64     // font size
		DPI                              float64     // DPI
		Lines, Points, Rotates, Distorts int         // Number of
		Bit                              int         // Length of captcha text
	}

	Data struct {
		Content string
		Result  string
		Expire  int64
	}

	Captcha struct {
		*Config
		drawer *font.Drawer
	}
)

func (c *Config) ParseFont(buf []byte) error {
	obj, err := sfnt.Parse(buf)
	if err != nil {
		return err
	}
	if c.DPI == 0 {
		c.DPI = 72
	}
	if c.Size == 0 {
		c.Size = 24
	}
	c.Face, err = opentype.NewFace(obj, &opentype.FaceOptions{
		Size: c.Size, DPI: c.DPI, Hinting: font.HintingNone,
	})
	return err
}

func (c *Config) Validate() error {
	if c.Width <= 0 {
		return errors.New("width must be greater than 0")
	}
	if c.Height <= 0 {
		return errors.New("height must be greater than 0")
	}
	if c.Bg == nil {
		return errors.New("background color is nil")
	}
	if c.Ft == nil {
		return errors.New("font color is nil")
	}
	if c.Bit <= 0 {
		return errors.New("bit must be greater than 0")
	}
	if c.Face == nil {
		return errors.New("font face is nil")
	}
	return nil
}

func init() {
	defaultConfig.ParseFont(fontData)
}

func New(c *Config) (*Captcha, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if c.Face == nil {
		if err := c.ParseFont(fontData); err != nil {
			return nil, err
		}
	}
	// 初始化字体 & 字体画笔
	drawer := &font.Drawer{
		Src:  image.NewUniform(c.Ft),
		Face: c.Face,
		Dot:  fixed.P(0, c.Face.Metrics().Ascent.Floor()),
	}
	return &Captcha{drawer: drawer, Config: c}, nil
}

func (c *Captcha) Draw() (image.Image, *Data) {
	// 生成随机字符串
	text, result := randomEquation(c.Bit)
	// 字符串空间计算
	drawBound, advance := c.drawer.BoundString(text)
	textWidth := advance.Ceil() + 30
	textHeight := (drawBound.Max.Y - drawBound.Min.Y).Ceil() + 20
	if textWidth > c.Width {
		c.Width = textWidth
	}
	if textHeight > c.Height {
		c.Height = textHeight
	}
	// 初始化底图
	img := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
	draw.Draw(img, img.Bounds(), image.NewUniform(c.Bg), image.Point{}, draw.Over)
	// 绘制字符串
	c.drawText(img, text)
	// 绘制噪点
	c.drawNoise(img)
	return img, &Data{Content: text, Result: result}
}

func (c *Captcha) drawText(img draw.Image, text string) *Captcha {
	c.drawer.Dst = img
	// 字符串居中
	b, _ := c.drawer.BoundString(text)
	width := b.Max.X / fixed.I(len(text))
	// 水平居中
	x := (fixed.I(img.Bounds().Max.X) - b.Max.X) / 2
	dot := c.drawer.Dot
	c.drawer.Dot.X = x
	// 计算垂直方向的中心位置
	baseY := (fixed.I(img.Bounds().Max.Y) + c.Face.Metrics().Height) / 2
	for _, t := range text {
		// 设置字符的Y坐标，基于中心位置随机浮动
		c.drawer.Dot.Y = baseY + fixed.I(rand.Intn(10)-10)
		c.drawer.DrawBytes([]byte{byte(t)})
		c.drawer.Dot.X += width
	}
	c.drawer.Dot = dot
	return c
}

func (c *Captcha) drawNoise(img draw.Image) *Captcha {
	// draw rotates 旋转
	// c.drawRotate(img)
	// draw distorts 扭曲 - 简单的扭曲变形
	c.drawDistort(img)
	// draw points 点
	c.drawPoints(img)
	// draw lines 线
	c.drawNoiseLines(img)
	return c
}
