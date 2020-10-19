package response

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/org_transcoder/transcoder/preset"
)

// ConverDataList is json file struct
type ConverDataList struct {
	RunTime       string `json:"runTime"`
	FileSize      string `json:"fileSize"`
	FilePath      string `json:"filePath"`
	ThumbnailPath string `json:"thumbnailPath"`
	Cntinfo       string `json:"cntinfo"`
	Md5           string `json:"md5"`
	PresetName    string `json:"presetName"`
}

type OrgMetadata struct {
}

// Response is json file struct
type Response struct {
	ConvertCount    int                    `json:"convertCount"`
	RunTime         string                 `json:"runTime"`
	FileName        string                 `json:"FileName"`
	FilePath        string                 `json:"FilePath"`
	FileSize        string                 `json:"FileSize"`
	OrgMetadata     map[string]interface{} `json:"orgMetadata"`
	Md5             string                 `json:"md5"`
	ConvertDataList []ConverDataList       `json:"converDataList"`
	RtMsg           string                 `json:"rtMsg"`
	Rt              string                 `json:"rt"`
}

// ResponseJson is Transcoding result json file
func ResponseJson() Response {

	jobTemplate := "/usr/src/tr/transcode.json"
	data := preset.Loadjson(jobTemplate)
	// fmt.Println("chk debug : response.go")

	res := Response{}

	var org map[string]interface{}

	orgfile := "/usr/src/tr/" + data.FileName
	// orgfile := data.FileName
	// fmt.Println(orgfile)

	filepath := data.FilePath
	res.FilePath = strings.ReplaceAll(filepath, "/home/woori", "")

	filename := data.FileName
	res.FileName = filename

	convercount := len(data.Transcodings)
	res.ConvertCount = convercount

	// ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 -sexagesimal /tmp/sample.mp4
	// runtime, _ := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=duration", "-of", "default=noprint_wrappers=1:nokey=1", "-sexagesimal", orgfile).Output()
	runtime, _ := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", "-sexagesimal", orgfile).Output()
	res.RunTime = strings.TrimRight(string(runtime), "\n")

	filesize, _ := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=size", "-of", "default=noprint_wrappers=1:nokey=1", orgfile).Output()
	res.FileSize = strings.TrimRight(string(filesize), "\n")

	md5, _ := exec.Command("md5sum", orgfile).Output()
	res.Md5 = strings.Split(strings.TrimRight(string(md5), "\n"), " ")[0]

	rtMsg := "OK"
	res.RtMsg = rtMsg

	rt := "100"
	res.Rt = rt

	// ffmpegConf := &ffmpeg.Config{
	// 	FfmpegBinPath:   "/usr/local/bin/ffmpeg",
	// 	FfprobeBinPath:  "/usr/local/bin/ffprobe",
	// 	ProgressEnabled: true,
	// }
	// f := ffmpeg.New(ffmpegConf)
	// &f.

	// ffprobe -v quiet -print_format json -show_format -i /tmp/sample.mp4
	os.Chdir("/usr/src/tr")
	metaFileName := data.FileName
	cmd, _ := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-i", metaFileName).Output()
	json.Unmarshal(cmd, &org)
	res.OrgMetadata = org["format"].(map[string]interface{})

	// ************ Make ConverDataList ************
	list := make([]ConverDataList, 0)

	// reference : https://trac.ffmpeg.org/wiki/FFprobeTips
	for i := 0; i < len(data.Transcodings); i++ {
		convertlist := ConverDataList{}
		output := "/usr/src/tr/" + data.Transcodings[i].Output

		runTime, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", "-sexagesimal", output).Output()
		convertlist.RunTime = strings.TrimRight(string(runTime), "\n")

		filesize, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "error", "-show_entries", "format=size", "-of", "default=noprint_wrappers=1:nokey=1", output).Output()
		convertlist.FileSize = strings.TrimRight(string(filesize), "\n")

		customerName := strings.Split(data.FilePath, "/")[2]
		replaceHome := strings.ReplaceAll(data.FilePath, "/home/"+customerName, "") // /home/woori/test/sample2.mp4 -> /test/sample2.mp4
		splitMp4 := data.FileName                                                   // sample2.mp4
		customerUpPath := strings.ReplaceAll(replaceHome, splitMp4, "")             // /test/sample2.mp4 -> /test/
		filepath := "http://wooribank-vod.skcdn.net" + customerUpPath + data.Transcodings[i].Output
		convertlist.FilePath = filepath

		// if data.DefaultThumbnail
		thbName := strings.Split(data.FileName, ".")[0] + ".png"
		thumbnailPath := "http://wooribank-vod.skcdn.net/thb" + customerUpPath + thbName
		convertlist.ThumbnailPath = thumbnailPath

		md5, _ := exec.Command("md5sum", output).Output()
		convertlist.Md5 = strings.Split(strings.TrimRight(string(md5), "\n"), " ")[0]

		presetname := data.Presets[i].PresetName
		convertlist.PresetName = presetname

		resolution, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", output).Output()
		cutResolution := strings.ReplaceAll(string(resolution), "\n", "")
		fps, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "0", "-of", "csv=p=0", "-select_streams", "v:0", "-show_entries", "stream=r_frame_rate", output).Output()
		cutFps := strings.ReplaceAll(string(fps), "\n", "")
		codec, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=codec_name", "-of", "default=noprint_wrappers=1:nokey=1", output).Output()
		cutCodec := strings.ReplaceAll(string(codec), "\n", "")
		bitrate, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "error", "-show_entries", "format=bit_rate", "-of", "default=noprint_wrappers=1:nokey=1", output).Output()
		cutBitrate := strings.ReplaceAll(string(bitrate), "\n", "")
		audiocodec, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "error", "-select_streams", "a:0", "-show_entries", "stream=codec_name", "-of", "default=noprint_wrappers=1:nokey=1", output).Output()
		cutAudioCodec := strings.ReplaceAll(string(audiocodec), "\n", "")
		convertlist.Cntinfo = "resolution=" + string(cutResolution) + " fps=" + cutFps + " codec=" + cutCodec + " bitrate=" + cutBitrate + " audio_codec=" + cutAudioCodec

		list = append(list, convertlist)

	}
	res.ConvertDataList = list
	// fmt.Println(list)
	// ************ Make ConverDataList ************

	return res

}

// ResponseJson2 is Transcoding result json file
func ResponseJson2() Response {

	jobTemplate := "/usr/src/tr/transcode.json"
	data := preset.Loadjson(jobTemplate)
	// fmt.Println("chk debug : response.go")

	res := Response{}

	var org map[string]interface{}

	orgfile := "/usr/src/tr/" + data.FileName
	// orgfile := data.FileName
	// fmt.Println(orgfile)

	filepath := data.FilePath
	res.FilePath = strings.ReplaceAll(filepath, "/home/woori", "")

	filename := data.FileName
	res.FileName = filename

	convercount := len(data.Transcodings)
	res.ConvertCount = convercount

	// ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 -sexagesimal /tmp/sample.mp4
	// runtime, _ := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=duration", "-of", "default=noprint_wrappers=1:nokey=1", "-sexagesimal", orgfile).Output()
	runtime, _ := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", "-sexagesimal", orgfile).Output()
	res.RunTime = strings.TrimRight(string(runtime), "\n")

	filesize, _ := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=size", "-of", "default=noprint_wrappers=1:nokey=1", orgfile).Output()
	res.FileSize = strings.TrimRight(string(filesize), "\n")

	md5, _ := exec.Command("md5sum", orgfile).Output()
	res.Md5 = strings.Split(strings.TrimRight(string(md5), "\n"), " ")[0]

	rtMsg := "OK"
	res.RtMsg = rtMsg

	rt := "100"
	res.Rt = rt

	// ffmpegConf := &ffmpeg.Config{
	// 	FfmpegBinPath:   "/usr/local/bin/ffmpeg",
	// 	FfprobeBinPath:  "/usr/local/bin/ffprobe",
	// 	ProgressEnabled: true,
	// }
	// f := ffmpeg.New(ffmpegConf)
	// &f.

	// ffprobe -v quiet -print_format json -show_format -i /tmp/sample.mp4
	os.Chdir("/usr/src/tr")
	metaFileName := data.FileName
	cmd, _ := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-i", metaFileName).Output()
	json.Unmarshal(cmd, &org)
	res.OrgMetadata = org["format"].(map[string]interface{})

	// ************ Make ConverDataList ************
	list := make([]ConverDataList, 0)

	// reference : https://trac.ffmpeg.org/wiki/FFprobeTips
	for i := 0; i < len(data.Transcodings); i++ {
		convertlist := ConverDataList{}
		output := "/usr/src/tr/" + data.Transcodings[i].Output

		runTime, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", "-sexagesimal", output).Output()
		convertlist.RunTime = strings.TrimRight(string(runTime), "\n")

		filesize, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "error", "-show_entries", "format=size", "-of", "default=noprint_wrappers=1:nokey=1", output).Output()
		convertlist.FileSize = strings.TrimRight(string(filesize), "\n")

		customerName := strings.Split(data.FilePath, "/")[2]
		replaceHome := strings.ReplaceAll(data.FilePath, "/home/"+customerName, "") // /home/woori/test/sample2.mp4 -> /test/sample2.mp4
		splitMp4 := data.FileName                                                   // sample2.mp4
		customerUpPath := strings.ReplaceAll(replaceHome, splitMp4, "")             // /test/sample2.mp4 -> /test/
		filepath := customerUpPath + data.Transcodings[i].Output
		convertlist.FilePath = filepath

		// if data.DefaultThumbnail
		thbName := strings.Split(data.FileName, ".")[0] + ".png"
		thumbnailPath := customerUpPath + thbName
		convertlist.ThumbnailPath = thumbnailPath

		md5, _ := exec.Command("md5sum", output).Output()
		convertlist.Md5 = strings.Split(strings.TrimRight(string(md5), "\n"), " ")[0]

		presetname := data.Presets[i].PresetName
		convertlist.PresetName = presetname

		resolution, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", output).Output()
		cutResolution := strings.ReplaceAll(string(resolution), "\n", "")
		fps, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "0", "-of", "csv=p=0", "-select_streams", "v:0", "-show_entries", "stream=r_frame_rate", output).Output()
		cutFps := strings.ReplaceAll(string(fps), "\n", "")
		codec, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=codec_name", "-of", "default=noprint_wrappers=1:nokey=1", output).Output()
		cutCodec := strings.ReplaceAll(string(codec), "\n", "")
		bitrate, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "error", "-show_entries", "format=bit_rate", "-of", "default=noprint_wrappers=1:nokey=1", output).Output()
		cutBitrate := strings.ReplaceAll(string(bitrate), "\n", "")
		audiocodec, _ := exec.Command("/usr/src/tr/ffprobe", "-v", "error", "-select_streams", "a:0", "-show_entries", "stream=codec_name", "-of", "default=noprint_wrappers=1:nokey=1", output).Output()
		cutAudioCodec := strings.ReplaceAll(string(audiocodec), "\n", "")
		convertlist.Cntinfo = "resolution=" + string(cutResolution) + " fps=" + cutFps + " codec=" + cutCodec + " bitrate=" + cutBitrate + " audio_codec=" + cutAudioCodec

		list = append(list, convertlist)

	}
	res.ConvertDataList = list
	// fmt.Println(list)
	// ************ Make ConverDataList ************

	return res

}
