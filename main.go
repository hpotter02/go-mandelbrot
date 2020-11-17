package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"sync"
)

const (
	imgY = 100
	imgX = 100
	numY = 4
	numX = 4
	size = 1000
)

type cRecantlge struct {
	Min complex128
	Max complex128
}

func main() {
	/*for i := 0; i < 3000; i++ {
		math.Pow(1, 2)
		zoom := float64(math.Pow(1.01, float64(-i)))
		img := getImage(image.Rect(0, 0, 1280, 720), zoomToPoint(cRecantlge{Min: complex(-2.24, -1.26), Max: complex(2.24, 1.26)}, zoom, complex(0.31822, -0.447110006)))
		f, err := os.Create(fmt.Sprintf("img/%04d.png", i))
		if err != nil {
			fmt.Println(err)
		}
		png.Encode(f, img)
		fmt.Printf("Finished: %04d/2000\n", i)
	}*/
	img1 := getImage(image.Rect(0, 0, 10000, 10000), cRecantlge{Min: complex(0.2, 0.2), Max: complex(0.5, 0.5)})

	f1, err := os.Create("tmp1.png")
	if err != nil {
		fmt.Println(err)
	}
	png.Encode(f1, img1)
}

func getImage(size image.Rectangle, frame cRecantlge) image.Image {
	img := image.NewNRGBA(size)
	xlen := float64(size.Max.X - size.Min.X)
	ylen := float64(size.Max.Y - size.Min.Y)
	wg := sync.WaitGroup{}
	for i := 0; i < int(xlen); i++ {
		for ii := 0; ii < int(ylen); ii++ {
			wg.Add(1)
			go func(r float64, i float64) {
				c := complex(real(frame.Min)+(real(frame.Max)-real(frame.Min))*r/xlen, imag(frame.Min)+(imag(frame.Max)-imag(frame.Min))*i/ylen)
				img.Set(int(r), int(i), getPointRGB(c))
				wg.Done()
			}(float64(i), float64(ii))
		}
	}
	wg.Wait()
	return img
}

func zoomToPoint(in cRecantlge, k float64, p complex128) cRecantlge {
	//mid := (in.Max - in.Min) / 2
	finalmin := complex(real(in.Min)*k, imag(in.Min)*k) + p
	finalmax := complex(real(in.Max)*k, imag(in.Max)*k) + p
	//relmin := (in.Min - mid)
	//relmax := (in.Max - mid)

	return cRecantlge{
		Min: finalmin,
		Max: finalmax,
	}
}

func iterate(n complex128, c complex128) complex128 {
	return n*n + c
}

func getPoint(p complex128) (complex128, int) {
	n := iterate(0, p)
	for i := 0; i < 1000; i++ {
		n = iterate(n, p)
		if cmplx.IsNaN(n) {
			return n, i
		}
	}
	return n, 255
}

func getPointRGB(p complex128) color.RGBA {
	v, i := getPoint(p)
	if cmplx.IsNaN(v) {
		return getColor(float64(i))
	}
	return color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	}

}

func getColor(i float64) color.RGBA {

	r := math.Sin(i*math.Pi/180) * 255
	g := math.Sin((i*2)*math.Pi/180) * 255
	b := math.Sin((i*4)*math.Pi/180) * 255

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(255),
	}
}
