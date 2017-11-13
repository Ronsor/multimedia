package main
import "github.com/ronsor/multimedia"
import "os"
import "log"

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil { log.Fatal(err) }
	d, err := multimedia.NewVideoStream(f)
	if err != nil { log.Fatal(err) }
	log.Printf("JPEG Video, FPS=%d", d.FPS)
	frames := 0
	sized := false
	prevsz := int64(0)
	for {
		frame, err := d.GetFrame()
		if err != nil {
			log.Printf("Read %d frames.", frames)
			log.Fatal(err)
		}
		frames++
		if !sized {
			log.Printf("Video dimensions: %dx%d", frame.Bounds().Max.X, frame.Bounds().Max.Y)
			sized = true
			b, _ := f.Seek(0,1)
			log.Printf("%d bytes", b-prevsz)
			prevsz = b
		}
	}
}
