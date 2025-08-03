package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
  f, err := os.Open("JohnieTre-white.png")
  if err != nil {
    panic(err)
  }
  defer f.Close()
  img, err := png.Decode(f)
  if err != nil {
    panic(err)
  }
  bounds := img.Bounds()
  newImg := image.NewRGBA(bounds)
  for x := 0; x < bounds.Dx(); x++ {
    for y := 0; y < bounds.Dy(); y++ {
      r, g, b, a := img.At(x, y).RGBA()
      c := color.RGBA64{
        R: uint16(r),
        G: uint16(g),
        B: uint16(b),
        A: uint16(a),
      }
      if r > 0xFFF && g > 0xFFF && b > 0xFFF {
        c.A = 0
        newImg.SetRGBA64(x, y, c)
      } else {
        newImg.SetRGBA64(x, y, c)
      }
    }
  }
  out, err := os.Create("JohnieTre-trans.png")
  if err != nil {
    panic(err)
  }
  if err := png.Encode(out, newImg); err != nil {
    panic(err)
  }
  fmt.Println("done")
}
