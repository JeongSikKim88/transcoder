package thumbnail

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/org_transcoder/minio"
	"github.com/org_transcoder/transcoder/preset"
)

// MakeThumbnail by using ffmpeg
func MakeThumbnail() {

	jobTemplate := "/usr/src/tr/transcode.json"
	data := preset.Loadjson(jobTemplate)

	resultIP := data.ResultIP

	orgFile := "/usr/src/tr/" + data.FileName
	thbName := strings.Split(data.FileName, ".")[0] + ".png"
	customerName := strings.Split(data.FilePath, "/")[2] // /home/woori/test/sample.mp4

	// lenName := strings.SplitAfter(data.FilePath, "/")
	// fileName := strings.Split(lenName, "/")[len(lenName)-1]

	replaceHome := strings.ReplaceAll(data.FilePath, "/home/"+customerName+"/", "") // woori/test2/sample2.mp4
	splitMp4 := data.FileName                                                       // sample2.mp4
	customerUpPath := strings.ReplaceAll(replaceHome, splitMp4, "")                 // test2/

	cmd := exec.Command("/usr/src/tr/ffmpeg", "-i", orgFile, "-ss", "5", "-vcodec", "png", "-vframes", "1", "-y", "/usr/src/tr/"+thbName)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	} else {
		log.Println(thbName + " thumbnail is created")
	}

	// minio upload file name
	uploadPath := customerName + "/" + "thb/" + customerUpPath + thbName
	// uploadPath := ".org/woori/test2/" + data.FileName
	// minio.FileUploader(customeName, orgFile, objectFile)
	minio.FileUploader(resultIP, customerName, "/usr/src/tr/"+thbName, uploadPath)

}
