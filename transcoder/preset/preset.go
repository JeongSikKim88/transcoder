package preset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Audio is audio encoding info
type Audio struct {
	codec          string `json:"codec"`
	profile        string `json:"profile"`
	BitrateControl string `json:"BitrateControl"`
	Coding         string `json:"Coding"`
	Bitrate        string `json:"Bitrate"`
	Sample         string `json:"Sample"`
}

// Video is video encoding info
type Video struct {
	Codec         string `json:"Codec"`
	Bitrate       string `json:"Bitrate"`
	Width         string `json:"Width"`
	Height        string `jsong:"Height"`
	FrameRate     string `json:"FrameRate"`
	Profile       string `json:"Profile"`
	Level         string `json:"Level"`
	PixelAspect   string `json:"PixelAspect"`
	RateControl   string `json:"RateControl"`
	Pass          string `json:"Pass"`
	GOP           string `json:"GOP"`
	Bframe        string `json:"Bframe"`
	Iframe        string `json:"Iframe"`
	Interlacing   string `json:"Interlacing"`
	Deinterlacing string `json:"Deinterlacing"`
	Format        string `json:"Format"`
	OverlayXY     string `json:"OverlayXY"`
	OverlaySize   string `json:"OverlaySize"`
	OverlayImage  string `json:"OverlayImage"`
}

// Presets is preset in subtask
type Presets struct {
	PresetName string                 `json:"PresetName"`
	Uploader   string                 `json:"Uploader"`
	Creater    string                 `json:"Creater"`
	Video      map[string]interface{} `json:"Video"`
	Audio      map[string]interface{} `json:"Audio"`
}

// Transcodings is subtask in Info struct
type Transcodings struct {
	PresetName      string `json:"PresetName"`
	Output          string `json:"Output"`
	OverlayImage    string `json:"OverlayImage"`
	OverlayPosition string `json:"overlayPosition"`
}

// DefaultThumbnail is extracting thumbnail in Info struct
type DefaultThumbnail struct {
	ThumbNum int `json:"ThumbNum"`
}

// Info is json file struct
type Info struct {
	JobID            string                 `json:"JobID"`
	FileName         string                 `json:"FileName"`
	CustomerName     string                 `json:"CustomerName"`
	CustomerID       string                 `json:"CustomerID"`
	FilePath         string                 `json:"FilePath"`
	UploadIP         string                 `json:"UploadIP"`
	ResultIP         string                 `json:"ResultIP"`
	TemplateName     string                 `json:"TemplateName"`
	DefaultThumbnail map[string]interface{} `json:"DefaultThumbnail"`
	GridThumbnail    string                 `json:"GridThumbnail"`
	Callback         string                 `json:"Callback"`
	Transcodings     []Transcodings         `json:"Transcodings"`
	Presets          []Presets              `json:"Presets"`
}

// Loadjson is Transcoding json file
func Loadjson(jobTemplate string) Info {

	defer func() {
		s := recover()
		fmt.Println(s)
	}()

	// jsonFile, err := ioutil.ReadFile("/tmp/transcode.json")
	jsonFile, err := ioutil.ReadFile(jobTemplate)
	if err != nil {
		fmt.Println(err)
	}

	var info Info
	// save json content in data
	json.Unmarshal(jsonFile, &info)
	return info
}
