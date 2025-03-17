package captcha

import (
	"bytes"
	"math"
	"math/rand"
	"strconv"
)

const (
	equationAdd = iota // x + y = z
	equationSub        // x - y = z
	equationMul        // x * y = z
	equationDiv        // x / y = z
)

func randomEquation(bit int) (content, result string) {
	var resultNum int64
	buf := new(bytes.Buffer)
	minNum := int64(math.Pow10(bit - 1))
	switch rand.Intn(4) {
	case equationAdd: // left + right = result
		var left, right int64
		resultNum = nBitNum(bit)
		for resultNum < minNum*2 {
			resultNum = nBitNum(bit)
		}
		for left < minNum || right < minNum {
			left = rand.Int63n(resultNum-minNum) + minNum
			right = resultNum - left
		}
		buf.WriteString(strconv.FormatInt(left, 10))
		buf.WriteByte('+')
		buf.WriteString(strconv.FormatInt(resultNum-left, 10))
	case equationSub: // left - right = result
		var left, right int64
		for left < minNum*2 {
			left = nBitNum(bit)
		}
		for resultNum < minNum || right < minNum {
			resultNum = rand.Int63n(left-minNum) + minNum
			right = left - resultNum
		}
		buf.WriteString(strconv.FormatInt(left, 10))
		buf.WriteByte('-')
		buf.WriteString(strconv.FormatInt(right, 10))
	case equationMul: // left * right = result
		left, right := nBitNum(bit), nBitNum(bit)
		resultNum = left * right
		buf.WriteString(strconv.FormatInt(left, 10))
		buf.WriteByte('*')
		buf.WriteString(strconv.FormatInt(right, 10))
	case equationDiv: // left / right = result
		left, right := nBitNum(bit), nBitNum(bit)
		resultNum = left * right
		resultNum, left = left, resultNum
		buf.WriteString(strconv.FormatInt(left, 10))
		buf.WriteByte('/')
		buf.WriteString(strconv.FormatInt(right, 10))
	}
	return buf.String(), strconv.FormatInt(resultNum, 10)
}
