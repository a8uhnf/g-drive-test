package filepath

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFileNameFunc(t *testing.T) {
	path := "/home/hanifa/Downloads/CrazyHD.com-Batman Mystery Of The Batwoman 2003 BrRip 720p x264 YIFY (1).torrent"
	filename := Filename(path)
	assert.Equal(t, "CrazyHD.com-Batman Mystery Of The Batwoman 2003 BrRip 720p x264 YIFY (1).torrent", filename)
}