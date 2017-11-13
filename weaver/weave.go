package main
import "fmt"
import "os"
import "github.com/ronsor/multimedia"
func stderr(v interface{}) {
	os.Stderr.Write([]byte(v.(string)))
}
func stderrln(v interface{}) {
	os.Stderr.Write([]byte(v.(string)))
	os.Stderr.Write([]byte("\n"))
}
func abort(v interface{}) {
	stderrln(v)
	os.Exit(1)
}
func main() {
	silence := make([]byte, 533)
	if len(os.Args) != 3 {
		stderrln("Usage: weave [video.vjpg] [audio.pcm] > [out.wov]")
		os.Exit(1)
	}
	vf, err := os.Open(os.Args[1])
	if err != nil { abort(err) }
	af, err := os.Open(os.Args[2])
	if err != nil { abort(err) }
	v, _ := multimedia.NewVideoStream(vf)
	a, _ := multimedia.NewAudioStream(af)
	frames := 0
	for {
		fr, err := v.GetRawFrame()
		if err != nil { break }
		os.Stdout.Write(fr)
		afr, err := a.GetFrame()
		stderr(fmt.Sprintf("\r%d frames (%d:%02d)", frames, frames/15/60, frames/15%60))
		if err != nil {
			os.Stdout.Write(silence)
		} else {
			os.Stdout.Write(afr)
		}
		frames++
	}
}
