package minecraft

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aimjel/minecraft/chat"
	"image/png"
	"math"
	"os"
)

type Status struct {
	enc *json.Encoder

	buf *bytes.Buffer

	s *status
}

func NewStatus(protocol, max int, desc string) *Status {
	var s status
	s.Version.Name, s.Version.Protocol = versionName(protocol), protocol
	s.Players.Max, s.Description = max, chat.NewMessage(desc)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	_ = enc.Encode(s)

	st := &Status{enc: enc, buf: &buf, s: &s}

	size := buf.Len() + 34 //34 for the favicon key and prepended info, including quotes and comma

	b := bytes.NewBuffer(nil)
	st.loadIcon(b)

	if size+b.Len() < math.MaxInt16 {
		st.s.Favicon = "data:image/png;base64," + base64.StdEncoding.EncodeToString(b.Bytes())
	}

	enc.Encode(s)
	return st
}

func (s *Status) loadIcon(buf *bytes.Buffer) {
	f, err := os.Open("server-icon.png")
	defer f.Close()
	if err != nil {
		return
	}

	_, _ = f.Seek(0, 0)
	m, err := png.Decode(f)
	if err != nil {
		return
	}

	var e png.Encoder
	e.CompressionLevel = png.DefaultCompression

	if err = e.Encode(buf, m); err != nil {
		fmt.Printf("%v compressiong server icon", err)
	}

}

func (s *Status) json() []byte {
	return s.buf.Bytes()
}

func versionName(protocol int) string {
	return map[int]string{
		757: "1.18.1",
		756: "1.17.1",
		755: "1.17",
	}[protocol]
}

// status represents the json response in struct for more performance
type status struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description chat.Message `json:"description"`
	Favicon     string       `json:"favicon"`
}
