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
	Noise struct {
		Lines    int
		Points   int
		Rotates  int
		Distorts int
	}
	Config struct {
		Width, Height int         // width and height of the captcha image
		Bg, Ft        color.Color // background and font color
		Size          float64     // font size
		DPI           float64     // DPI
		Bit           int         // Length of captcha text
		Noise
	}
	Data struct {
		Content string
		Result  string
		Expire  int64
	}
	Captcha struct {
		*Config
		font *sfnt.Font
		face font.Face
	}
)

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
	return nil
}

func New(c *Config) (*Captcha, error) {
	err := c.Validate()
	if err != nil {
		return nil, err
	}
	font, err := sfnt.Parse(fontData)
	if err != nil {
		return nil, err
	}
	if c.DPI == 0 {
		c.DPI = 72
	}
	if c.Size == 0 {
		c.Size = 24
	}
	face, err := parseFont(font, c.Size, c.DPI)
	if err != nil {
		return nil, err
	}
	return &Captcha{Config: c, font: font, face: face}, nil
}

func parseFont(f *opentype.Font, size, dpi float64) (font.Face, error) {
	return opentype.NewFace(f, &opentype.FaceOptions{
		Size: size, DPI: dpi, Hinting: font.HintingNone,
	})
}

func (c *Captcha) Draw() (image.Image, *Data) {
	// 创建本地副本，避免并发修改共享状态
	width, height := c.Width, c.Height
	// 为每次Draw调用创建新的font face，彻底消除并发问题
	fontFace, err := parseFont(c.font, c.Size, c.DPI)
	if err != nil {
		fontFace = c.face
	}
	// 生成随机字符串
	text, result := randomEquation(c.Bit)
	// 创建完全独立的drawer对象
	drawer := &font.Drawer{
		Src:  image.NewUniform(c.Ft),
		Face: fontFace,
		Dot:  fixed.P(0, fontFace.Metrics().Ascent.Floor()),
	}
	// 字符串空间计算
	drawBound, advance := drawer.BoundString(text)
	textWidth := advance.Ceil() + 30
	textHeight := (drawBound.Max.Y - drawBound.Min.Y).Ceil() + 20
	// 使用局部变量，而不是修改共享状态
	if textWidth > width {
		width = textWidth
	}
	if textHeight > height {
		height = textHeight
	}
	// 初始化底图
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.NewUniform(c.Bg), image.Point{}, draw.Over)
	// 使用独立的drawer绘制文本
	drawTextConcurrent(img, text, drawer, fontFace)
	// 绘制噪点
	drawNoise(img, width, height, c.Noise)
	return img, &Data{Content: text, Result: result}
}

func drawTextConcurrent(img draw.Image, text string, drawer *font.Drawer, face font.Face) {
	drawer.Dst = img
	// 字符串居中
	b, _ := drawer.BoundString(text)
	width := b.Max.X / fixed.I(len(text))
	// 水平居中
	x := (fixed.I(img.Bounds().Max.X) - b.Max.X) / 2
	dot := drawer.Dot
	drawer.Dot.X = x
	// 计算垂直方向的中心位置
	baseY := (fixed.I(img.Bounds().Max.Y) + face.Metrics().Height) / 2
	for _, t := range text {
		// 设置字符的Y坐标，基于中心位置随机浮动
		drawer.Dot.Y = baseY + fixed.I(rand.Intn(10)-10)
		drawer.DrawBytes([]byte{byte(t)})
		drawer.Dot.X += width
	}
	drawer.Dot = dot
}

func drawNoise(img draw.Image, w, h int, n Noise) {
	// draw distorts 扭曲 - 简单的扭曲变形
	drawDistort(img, w, h)
	// draw points 点
	drawPoints(img, w, h, n.Points)
	// draw lines 线
	drawNoiseLines(img, w, h, n.Lines)
}
