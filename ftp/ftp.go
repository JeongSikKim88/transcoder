package ftp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/org_transcoder/response"
	"github.com/org_transcoder/transcoder/preset"
)

func FileDownloader() {

	jobTemplate := "/usr/src/tr/transcode.json"
	data := preset.Loadjson(jobTemplate)

	truploadIP := data.UploadIP
	// uploadIP := data.ResultIP

	// filePath := data.FilePath // /home/atlas/test2/sample3.mp4
	// customerName := strings.Split(filePath, "/")[2]
	// cutHome := strings.ReplaceAll(filePath, "/home/"+customerName, "")
	// fileName := data.FileName
	// fmt.Println(cutHome)

	c, err := ftp.Dial(truploadIP+":2100", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Println(err)
	}

	// c, err := ftp.Connect(truploadIP + ":2100")
	// if err != nil {
	// 	log.Println(err)
	// }

	err = c.Login("atlas", "Qwpo1209")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("login success")
	}

	// res, err := c.
	file, err := os.Open(
		"/usr/src/tr/sample.mp4",
		// os.O_CREATE|os.O_RDWR|os.O_TRUNC, // 파일이 없으면 생성,
		// 읽기/쓰기, 파일을 연 뒤 내용 삭제
		// os.FileMode(0644), // 파일 권한은 644
	)

	r := bufio.NewReader(file)
	err = c.Stor("/home/atlas/test2/sample.mp4", r)
	if err != nil {
		log.Println(err)
	}

	// Do something with the FTP conn
	if err := c.Quit(); err != nil {
		log.Println(err)
	}
}

func CurlDownFtp() {

	jobTemplate := "/usr/src/tr/transcode.json"
	data := preset.Loadjson(jobTemplate)

	truploadIP := data.UploadIP
	port := "2100"
	// uploadIP := data.ResultIP

	filePath := data.FilePath // /home/atlas/test2/sample3.mp4
	customerName := strings.Split(filePath, "/")[2]
	cutHome := strings.ReplaceAll(filePath, "/home/"+customerName, "")
	fileName := data.FileName
	ftpUser := customerName + ":Qwpo1209"
	fmt.Println(cutHome)

	// curl --ftp-create-dirs -u "atlas:Qwpo1209" "ftp://trupload.myskcdn.net:2100/test2/sample3.mp4" -O
	cmd := exec.Command("curl", "--ftp-create-dirs", "-u", ftpUser, "ftp://"+truploadIP+":"+port+"/"+cutHome, "-o", "/usr/src/tr/"+fileName)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	} else {
		log.Println(fileName + " is downloaded")
	}

}

func CurlJsonUpFtp() {

	jobTemplate := "/usr/src/tr/transcode.json"
	data := preset.Loadjson(jobTemplate)
	res, err := json.Marshal(response.ResponseJson())
	if err != nil {
		log.Println(err)
	}

	filePath := data.FilePath // /home/atlas/test2/sample3.mp4
	customerName := strings.Split(filePath, "/")[2]
	ftpUser := customerName + ":Qwpo1209"
	uploadIP := data.ResultIP
	port := "2100"

	f2, _ := os.Create("/usr/src/tr/" + strings.Split(data.FileName, ".")[0] + ".json")
	defer f2.Close()
	n, err := f2.WriteString(string(res))
	// fmt.Println(strings.Split(data.FileName, ".")[0] + ".json" + " is created")
	log.Println(strings.Split(data.FileName, ".")[0] + ".json" + " is created")
	fmt.Println("file size : ", n)
	// fmt.Println(string(res))

	// transcoding result file
	jasonFile := "/usr/src/tr/" + strings.Split(data.FileName, ".")[0] + ".json"
	// minio upload file name
	objectFile := ".json/" + strings.Split(strings.ReplaceAll(data.FilePath, "/home/", ""), ".")[0] + ".json"

	// lenName := strings.SplitAfter(filePath, "/")
	// cutName := strings.Split(filePath, "/")[len(lenName)-1]
	// cutHome := strings.ReplaceAll(filePath, "/home/"+customerName, "")
	// upPath := strings.ReplaceAll(cutHome, cutName, "")

	// curl --ftp-create-dirs -u "atlas:Qwpo1209" "ftp://trupload.myskcdn.net:2100/test2/sample3.mp4" -O
	cmd := exec.Command("curl", "-T", jasonFile, "--ftp-create-dirs", "-u", ftpUser, "ftp://"+uploadIP+":"+port+"/"+objectFile)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		log.Println(cmd)
		log.Println(err)
	} else {
		log.Println(jasonFile + " is uploaded")
	}

}
