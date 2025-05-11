package main

import (
	"flag"
	"9fans.net/go/draw"
	"9fans.net/go/draw/memdraw"
	"fmt"
	"image"
	"log"
	"time"
)

const (
	width = 17
	height = 13
)

func main() {
	usemem := flag.Bool("usemem", false, "use memdraw")
	screxp := flag.Bool("screxp", true, "export 'unload' pix after drawing directly to screen")
	flag.Parse()

	errch := make(chan error, 1)
	rectStr := fmt.Sprintf("%dx%d", width, height)
	log.Printf("rectStr=%s", rectStr)
	display, err := draw.Init(errch, "", "hello", "800x600")
	if err != nil { panic(err.Error()) }

	r := image.Rectangle{Max: image.Point{X: width, Y: height}}
	buf := make([]byte, 1024*1024*32)

	if *screxp {
		log.Printf("export from actual screen")
		img, err := display.AllocImage(r, draw.ARGB32, false, draw.Blue)
		if err != nil { panic(err.Error()) }
		display.ScreenImage.Draw(r, img, nil, image.ZP)
		display.Flush()

		nbuf, err := display.ScreenImage.Unload(r, buf)
		if err != nil { panic(err.Error()) }
		buf = buf[:nbuf]

		log.Printf("17x13x4=%v", 17*13*4)
		log.Printf("buf len=%v", len(buf))
		log.Printf("buf=%+x", buf)

		_, err = display.ScreenImage.Load(r.Add(image.Point{100, 100}), buf)
		if err != nil { panic(err.Error()) }
		display.Flush()

		<-time.After(time.Hour)
		return
	}

	if *usemem {
		log.Printf("using memdraw")
		img, err := memdraw.AllocImage(r, display.ScreenImage.Pix)
		if err != nil {
			panic(err.Error())
		}
		memdraw.FillColor(img, draw.Black)
		nbuf, err := memdraw.Unload(img, r, buf)
		if err != nil { panic(err.Error()) }
		buf = buf[:nbuf]
	} else {
		log.Printf("using screen directly")
		img, err := display.AllocImage(r, display.ScreenImage.Pix, false, draw.Black)
		if err != nil { panic(err.Error()) }
		nbuf, err := img.Unload(r, buf)
		if err != nil { panic(err.Error()) }
		buf = buf[:nbuf]
	}

	log.Printf("17x13x4=%v", 17*13*4)
	log.Printf("buf len=%v", len(buf))
	log.Printf("buf=%+x", buf)

	_, err = display.ScreenImage.Load(r, buf)
	if err != nil { panic(err.Error()) }
	display.Flush()

	<-time.After(time.Hour)
}
