package captcha

import (
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"
)

func TestDraw(t *testing.T) {
	c, err := New(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	for i := range 3 {
		img, data := c.Draw()
		t.Logf("%+v\n\n", data)
		f, _ := os.OpenFile(fmt.Sprintf("./test%d.png", i), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		png.Encode(f, img)
		f.Close()
	}
	defer func() {
		time.Sleep(time.Second * 10)
		os.Remove("./test0.png")
		os.Remove("./test1.png")
		os.Remove("./test2.png")
	}()
	exec.Command("sh", "-c", "open ./test*.png").Run()
}

func TestConcurrentDraw(t *testing.T) {
	c, err := New(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	const goroutines = 10
	const iterations = 5

	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				c.Draw()
			}
		}()
	}

	wg.Wait()
	t.Logf("成功完成 %d 个goroutine，每个执行 %d 次Draw()调用，耗时 %s", goroutines, iterations, time.Since(start))
}

func BenchmarkDraw(b *testing.B) {
	c, err := New(defaultConfig)
	if err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		c.Draw()
	}
}
