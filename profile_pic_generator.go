package api

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
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
	img := RandomImage()
	np, nm, ne, neb := rand.Intn(10)+1, rand.Intn(10)+1, rand.Intn(10)+1, rand.Intn(10)+1
	pencil := LoadImage("./res/pencil_" + fmt.Sprint(np) + ".png")
	mouth := LoadImage("./res/mouth_" + fmt.Sprint(nm) + ".png")
	eyes := LoadImage("./res/eyes_" + fmt.Sprint(ne) + ".png")
	eyebrows := LoadImage("./res/eyebrows_" + fmt.Sprint(neb) + ".png")
	draw.Draw(img, pencil.Bounds(), pencil, image.Point{0, 0}, draw.Over)
	draw.Draw(img, mouth.Bounds(), mouth, image.Point{0, 0}, draw.Over)
	draw.Draw(img, eyes.Bounds(), eyes, image.Point{0, 0}, draw.Over)
	draw.Draw(img, eyebrows.Bounds(), eyebrows, image.Point{0, 0}, draw.Over)
	err := png.Encode(c.Writer, img)
	if err != nil {
		fmt.Println(err)
	}
}

func Uint8n(n int) uint8 {
	return uint8(rand.Intn(n))
}

func LoadImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	image, err := png.Decode(f)
	if err != nil {
		fmt.Println(err)
	}
	return image
}
