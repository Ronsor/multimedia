package main
import "github.com/ronsor/multimedia"
import "os"
import "log"
import "net/http"
func main() {
	f, err := os.Open(os.Args[1])
	if err != nil { log.Fatal(err) }
	d, err := multimedia.NewVideoStream(f)
	if err != nil { log.Fatal(err) }
	log.Printf("JPEG Video, FPS=%d", d.FPS)
	frames := 0
//	prevsz := int64(0)
	singleframe := make(chan []byte, 4)
	
	go func() { for {
		frame, err := d.GetRawFrame()
		if err != nil {
			log.Printf("Read %d frames.", frames)
			f.Seek(0,0)
		}
		singleframe <- frame
		frames++
	} } ()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte("<img style=width:100%;height:100% id=n src=/v.jpg>"))
		w.Write([]byte("<script>setInterval(function() { document.getElementById('n').src = '/v.jpg?' + Math.random() }, 50);</script>"))
	})
	http.HandleFunc("/v.jpg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "image/jpeg")
		w.Header().Add("Refresh", "0")
		data := <- singleframe
		w.Write(data)
	})
	http.ListenAndServe(os.Args[2], nil)
}
