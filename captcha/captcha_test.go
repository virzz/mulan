package captcha

import (
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"testing"
)

func TestDraw(t *testing.T) {
	c, err := New(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	for i := range 3 {
		img, data := c.Draw()
		fmt.Printf("%+v\n\n", data)
		f, _ := os.OpenFile(fmt.Sprintf("./test%d.png", i), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		png.Encode(f, img)
		f.Close()
	}
	exec.Command("sh", "-c", "open ./test*.png").Run()
}
