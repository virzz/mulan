package captcha

import (
	"image/color"
	"image/draw"
	"math"
	"math/rand"
)

// // drawRotate 旋转
// func (c *Captcha) drawRotate(img draw.Image) *Captcha {
// 	// 设置旋转角度（逆时针, 单位为度）
// 	radian := float64(c.Rotates) * math.Pi / 180.0
// 	// 计算新图像的尺寸
// 	cx, cy := float64(c.Width)/2, float64(c.Height)/2
// 	// 计算旋转之后新的图像尺寸
// 	newW := int(math.Abs(float64(c.Width)*math.Cos(radian)) + math.Abs(float64(c.Height)*math.Sin(radian)))
// 	newH := int(math.Abs(float64(c.Height)*math.Cos(radian)) + math.Abs(float64(c.Width)*math.Sin(radian)))
// 	// 创建新图像
// 	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
// 	draw.Draw(dst, dst.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
// 	// 旋转图像
// 	for y := 0; y < c.Height; y++ {
// 		for x := 0; x < c.Width; x++ {
// 			xf := float64(x) - cx
// 			yf := float64(y) - cy
// 			newX := int(xf*math.Cos(radian) - yf*math.Sin(radian) + float64(newW)/2)
// 			newY := int(yf*math.Cos(radian) + xf*math.Sin(radian) + float64(newH)/2)
// 			if newX >= 0 && newX < newW && newY >= 0 && newY < newH {
// 				img.Set(newX, newY, img.At(x, y))
// 			}
// 		}
// 	}
// 	// img = dst
// 	return c
// }

// drawDistort 扭曲
func (c *Captcha) drawDistort(img draw.Image) *Captcha {
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			newX := x + int(10*math.Sin(float64(y)/50))
			if newX >= 0 && newX < c.Width {
				img.Set(x, y, img.At(newX, y))
			} else {
				img.Set(x, y, color.Black) // 在边缘之外用黑色填补
			}
		}
	}
	return c
}

// drawPoints 绘制干扰点
func (c *Captcha) drawPoints(img draw.Image) *Captcha {
	for i := 0; i <= c.Points; i++ {
		img.Set(rand.Intn(c.Width)+1, rand.Intn(c.Height)+1, randColor())
	}
	return c
}

func (c *Captcha) drawNoiseLines(img draw.Image) *Captcha {
	lx := c.Width / 5
	rx := c.Width - lx
	for i := 0; i <= c.Lines; i++ {
		x := rand.Intn(lx)
		y := rand.Intn(c.Height)
		if i%3 == 0 {
			c.drawLine(img, x, y, rand.Intn(lx)+rx, rand.Intn(c.Height)) // 直线
		} else {
			c.drawArcLine(img, x, y, c.Width, c.Height) // 弧线
		}
	}
	return c
}

// drawLine 画直线 x0,y0 起点 x1,y1终点
// Bresenham算法(https://zh.wikipedia.org/zh-cn/布雷森漢姆直線演算法#最佳化)
func (c *Captcha) drawLine(img draw.Image, x0, y0, x1, y1 int) *Captcha {
	// 判断斜率是否大于1
	steep := abs(y1-y0) > abs(x1-x0)
	if steep {
		// 如果斜率大于1，交换x和y坐标
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}
	// 确保从左到右绘制
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	// 计算增量
	dx := x1 - x0
	dy := abs(y1 - y0)
	err := dx / 2
	y := y0
	// 确定y方向步进值
	ystep := 1
	if y0 >= y1 {
		ystep = -1
	}
	// 绘制线条
	for x := x0; x <= x1; x++ {
		if steep {
			img.Set(y, x, randColor(true))
		} else {
			img.Set(x, y, randColor(true))
		}
		err -= dy
		if err < 0 {
			y += ystep
			err += dx
		}
	}
	return c
}

// drawArcLine 绘制弧线
func (c *Captcha) drawArcLine(img draw.Image, x, y, width, height int) *Captcha {
	width = rand.Intn(width) + 50
	height = rand.Intn(height)
	startAngle, endAngle := rand.Intn(360), rand.Intn(360)
	var lx, ly, endx, endy int
	if (startAngle % 360) == (endAngle % 360) {
		startAngle, endAngle = 360, 0
	} else {
		if startAngle > 360 {
			startAngle = startAngle % 360
		}
		if endAngle > 360 {
			endAngle = endAngle % 360
		}
		for startAngle < 0 {
			startAngle += 360
		}
		for endAngle < startAngle {
			endAngle += 360
		}
		if startAngle == endAngle {
			startAngle = 0
			endAngle = 360
		}
	}
	for i := startAngle; i <= endAngle; i++ {
		endx = gdCosT[i%360]*width/2048 + x
		endy = gdSinT[i%360]*height/2048 + y
		if i != startAngle {
			c.drawLine(img, lx, ly, endx, endy)
		}
		lx, ly = endx, endy
	}
	return c
}
