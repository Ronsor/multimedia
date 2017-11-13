// Package media allows access to Woven media files (vjpeg-15,rawpcm8u8khz)
package multimedia
//import "time"
//import "errors"
//import "encoding/binary"
//import "bufio"
import "bytes"
import "image/jpeg"
import "image"
import "io"
type AudioStream struct{
	Reader io.Reader
	FrameSize int
}
// NewAudioStream returns a new audio stream
func NewAudioStream(r io.Reader) (*AudioStream, error) {
	return &AudioStream{Reader:r,FrameSize:533}, nil
}
// GetFrame returns a frame of audio (default 522 bytes for 15 fps video)
func (a *AudioStream) GetFrame() ([]byte, error) {
	ret := make([]byte, a.FrameSize)
	_, err := a.Reader.Read(ret)
	if err != nil { return nil, err }
	return ret, nil
}
type VideoStream struct{
	Reader io.Reader
	FPS int
	CurFrame int
}
// NewVideoStream returns a video stream
func NewVideoStream(r io.Reader) (*VideoStream, error) {
//	sig := make([]byte, 4)
//	_, e := r.Read(sig)
//	if e != nil { return nil, e }
//	if string(sig) != "RSAV" { return nil, errors.New("Invalid signature") }
//	fps := make([]byte, 1)
//	_, e2 := r.Read(fps)
//	if e2 != nil { return nil, e2 }
	return &VideoStream{Reader:r,CurFrame:0,FPS:15}, nil
}
// GetRawFrame returns one raw jpeg frame
func (v *VideoStream) GetRawFrame() ([]byte, error) {
	abc := make([]byte, 0, 4096)
	tmp := make([]byte, 1)
	lst := byte(0)
	for {
		_, err := v.Reader.Read(tmp)
		if err != nil { return nil, err }
		abc = append(abc, tmp[0])
		if lst == 0xFF && tmp[0] == 0xD9 {
			return abc, nil
		}
		lst = tmp[0]
	}

//	return v.Reader.ReadBytes(0xD9)
}
// GetFrame returns a frame of video as an image
func (v *VideoStream) GetFrame() (image.Image, error) {
//	defer func() { junk := []byte("0"); v.Reader.Read(junk); println(junk[0]) } ()
	v.CurFrame++
	rawf, err := v.GetRawFrame()
	if err != nil { return nil, err }
	return jpeg.Decode(bytes.NewReader(rawf))
}
type Stream struct{
	Video *VideoStream
	Audio *AudioStream
}
// NewStream creates a Woven media stream from a reader
func NewStream(cr io.Reader) (*Stream, error) {
	r := cr
	d := &Stream{}
	v, err := NewVideoStream(r)
	d.Video = v
	d.Audio, err = NewAudioStream(r)
	if err != nil { return nil, err }
	return d, nil
}
// GetFrame returns one frame of audio and video
func (s *Stream) GetFrame() ([]byte, image.Image, error) {
	i, e := s.Video.GetFrame()
	if e != nil { return nil, nil, e }
	a, e := s.Audio.GetFrame()
	if e != nil { return nil, nil, e }
	return a, i, nil
}
// GetRawFrame is the same, but returns a raw video frame instead
func (s *Stream) GetRawFrame() ([]byte, []byte, error) {
	i, e := s.Video.GetRawFrame()
	if e != nil { return nil, nil, e }
	a, e := s.Audio.GetFrame()
	if e != nil { return nil, nil, e }
	return a, i, nil
}
