package captcha

import (
	"image/color"
	"image/draw"
	"math"
	"math/rand"
)

func drawDistort(img draw.Image, w, h int) {
	for y := range h {
		for x := range w {
			newX := x + int(10*math.Sin(float64(y)/50))
			if newX >= 0 && newX < w {
				img.Set(x, y, img.At(newX, y))
			} else {
				img.Set(x, y, color.Black) // 在边缘之外用黑色填补
			}
		}
	}
}

// drawPoints 绘制干扰点
func drawPoints(img draw.Image, w, h, points int) {
	for range points {
		img.Set(rand.Intn(w)+1, rand.Intn(h)+1, randColor())
	}
}

func drawNoiseLines(img draw.Image, w, h, lines int) {
	lx := w / 5
	rx := w - lx
	for i := range lines {
		x := rand.Intn(lx)
		y := rand.Intn(h)
		if i%3 == 0 {
			drawLine(img, x, y, rand.Intn(lx)+rx, rand.Intn(h)) // 直线
		} else {
			drawArcLine(img, x, y, w, h) // 弧线
		}
	}
}

// drawLine 画直线 x0,y0 起点 x1,y1终点
// Bresenham算法(https://zh.wikipedia.org/zh-cn/布雷森漢姆直線演算法#最佳化)
func drawLine(img draw.Image, x0, y0, x1, y1 int) {
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
}

func drawArcLine(img draw.Image, x, y, width, height int) {
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
			drawLine(img, lx, ly, endx, endy)
		}
		lx, ly = endx, endy
	}
}
