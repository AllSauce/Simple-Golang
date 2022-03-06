//Stefan Nilsson 2013-02-27

//This program creates pictures of Julia sets (en.wikipedia.org/wiki/Julia_set).
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type ComplexFunc func(complex128) complex128

var Funcs []ComplexFunc = []ComplexFunc{
	func(z complex128) complex128 { return z*z - 0.61803398875 },
	func(z complex128) complex128 { return z*z + complex(0, 1) },
	func(z complex128) complex128 { return z*z + complex(-0.835, -0.2321) },
	func(z complex128) complex128 { return z*z + complex(0.45, 0.1428) },
	func(z complex128) complex128 { return z*z*z + 0.400 },
	func(z complex128) complex128 { return cmplx.Exp(z*z*z) - 0.621 },
	func(z complex128) complex128 { return (z*z+z)/cmplx.Log(z) + complex(0.268, 0.060) },
	func(z complex128) complex128 { return cmplx.Sqrt(cmplx.Sinh(z*z)) + complex(0.065, 0.122) },
}

func main() {
	runtime.GOMAXPROCS(10)
	before := time.Now()
	var wg sync.WaitGroup
	for n, fn := range Funcs {
		wg.Add(1)
		//before := time.Now()

		go CreatePng("picture-"+strconv.Itoa(n)+".png", fn, 1024, &wg)
		//after := time.Now()
		//println("time: ", after.Sub(before)/1000000, "ms")

	}

	wg.Wait()
	after := time.Now()
	println("time: ", after.Sub(before)/1000000, "ms")

}

// CreatePng creates a PNG picture file with a Julia image of size n x n.
func CreatePng(filename string, f ComplexFunc, n int, wg *sync.WaitGroup) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	err = png.Encode(file, Julia(f, n))
	if err != nil {
		log.Fatal(err)
	}
	wg.Done()

	return
}

// Julia returns an image of size n x n of the Julia set for f.
func Julia(f ComplexFunc, n int) image.Image {
	bounds := image.Rect(-n/2, -n/2, n/2, n/2)
	img := image.NewRGBA(bounds)
	s := float64(n / 4)

	theCounter := 0
	for i := bounds.Min.X; i < bounds.Max.X; i++ {

		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {

			switch {
			case theCounter <= 1000:

				theCounter++
				go func(i, j int) {
					n := Iterate(f, complex(float64(i)/s, float64(j)/s), 256)
					//Print green color for concurrent goroutines
					r := uint8(10)
					g := uint8(30)
					b := uint8(n % 32 * 8)
					img.Set(i, j, color.RGBA{r, g, b, 255})

					theCounter--
				}(i, j)

			default:
				n := Iterate(f, complex(float64(i)/s, float64(j)/s), 256)
				//Print red for non concurrent goroutines
				r := uint8(30)
				g := uint8(10)
				b := uint8(n % 32 * 8)
				img.Set(i, j, color.RGBA{r, g, b, 255})

			}
		}
	}
	return img
}

// Iterate sets z_0 = z, and repeatedly computes z_n = f(z_{n-1}), n â‰¥ 1,
// until |z_n| > 2  or n = max and returns this n.
func Iterate(f ComplexFunc, z complex128, max int) (n int) {
	for ; n < max; n++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			break
		}
		z = f(z)
	}
	return
}

/****Questions and Answers:*****/

//Q: How many CPUs does you program use?
//A: 6 cores 

//Q: How much faster is your parallell version?
//A: Total time improves by approximately 5000ms. Takes only about half the time: from 10000ms to 5000ms
