package api

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomImage() *image.RGBA {
	r, g, b := Uint8n(256), Uint8n(256), Uint8n(256)
	c := color.NRGBA{r, g, b, 255}
	i := image.NewRGBA(image.Rect(0, 0, 500, 500))
	draw.Draw(i, i.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
	return i
}

func GenerateProfilePicture(c *gin.Context) {
	err := png.Encode(c.Writer, RandomImage())
	if err != nil {
		panic(err)
	}
}

func Uint8n(n int) uint8 {
	return uint8(rand.Intn(n))
}
