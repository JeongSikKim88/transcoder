package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"

	"github.com/org_transcoder/ftp"
	"github.com/org_transcoder/minio"
	"github.com/org_transcoder/response"
	"github.com/org_transcoder/transcoder/ffmpeg"
	"github.com/org_transcoder/transcoder/preset"
	"github.com/org_transcoder/transcoder/thumbnail"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	file, err := os.OpenFile("/var/log/transcoder.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	// InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {

	app := fiber.New()
	// app.Use(logger.New())
	app.Use(middleware.Logger())
	app.Use(middleware.Logger("KST"))
	// app.Use(recover.New())

	// logwriter, e := syslog.New(syslog.LOG_NOTICE, "Transcoder")
	// if e == nil {
	// 	log.SetOutput(logwriter)
	// }

	// InfoLogger.Println("Starting the application...")
	// InfoLogger.Println("Something noteworthy happened")
	// WarningLogger.Println("There is something you should know about")
	// ErrorLogger.Println("Something went wrong")

	app.Get("/", func(c *fiber.Ctx) {

		// jobTemplate := "/usr/src/tr/transcode.json"
		// data := preset.Loadjson(jobTemplate)

		// customerName := strings.Split(filePath, "/")[2]
		// ftp.FileDownloader()
		// ftp.CurlDownFtp()

		// lenName := strings.SplitAfter(filePath, "/")
		// cutName := strings.Split(filePath, "/")[len(lenName)-1]
		// cutHome := strings.ReplaceAll(filePath, "/home/"+customerName, "")
		// upPath := strings.ReplaceAll(cutHome, cutName, "")

		// objectFile := ".json/" + strings.Split(strings.ReplaceAll(data.FilePath, "/home/", ""), ".")[0] + ".json"

		// ftp.FileDownloader()
		// filePath := data.FilePath // /home/atlas/test2/sample3.mp4
		// customerName := strings.Split(filePath, "/")[2]
		// ftpUser := customerName + ":Qwpo1209"

		// jasonFile := "/usr/src/tr/" + strings.Split(data.FileName, ".")[0] + ".json"
		// // minio upload file name
		// objectFile := ".json/" + strings.Split(strings.ReplaceAll(data.FilePath, "/home/", ""), ".")[0] + ".json"
		// uploadIP := data.ResultIP
		// port := "2100"

		// lenName := strings.SplitAfter(filePath, "/")
		// cutName := strings.Split(filePath, "/")[len(lenName)-1]
		// cutHome := strings.ReplaceAll(filePath, "/home/"+customerName, "")
		// upPath := strings.ReplaceAll(cutHome, cutName, "")

		// curl --ftp-create-dirs -u "atlas:Qwpo1209" "ftp://trupload.myskcdn.net:2100/test2/sample3.mp4" -O
		// cmd := exec.Command("curl", "-T", jasonFile, "--ftp-create-dirs", "-u", ftpUser, "ftp://"+uploadIP+":"+port+"/"+objectFile)
		// var out bytes.Buffer
		// var stderr bytes.Buffer
		// cmd.Stdout = &out
		// cmd.Stderr = &stderr
		// err := cmd.Run()
		// if err != nil {
		// 	fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		// } else {
		// 	log.Println(jasonFile + " is uploaded")
		// }

		c.Send("Hello, ATLAS Transcoder 👋!")
		// c.Send(objectFile)
	})

	app.Post("/tr1/job", func(c *fiber.Ctx) {

		result := c.Body()
		defer func() {
			s := recover()
			fmt.Println(s)
		}()

		f2, _ := os.Create("/usr/src/tr/transcode.json")
		defer f2.Close()
		n, err := f2.WriteString(result)
		if err != nil {
			ErrorLogger.Println(err)
		} else {
			InfoLogger.Println("JobTeplate is downloaded")
			InfoLogger.Println("transcode.json file size : ", n)
		}
		// fmt.Println(result)

		c.Send("Job Received\n", result)
		// c.Redirect("/tr1/transcode")
	})

	app.Get("/tr1/options", func(c *fiber.Ctx) {

		body := c.Body()
		fmt.Println(body)
		data := preset.Presets{}
		json.Unmarshal([]byte(body), &data)
		fmt.Println(data)
		chkOptions := make([]string, 0)
		// options := make([]string, 0)
		var videoOpt map[string]interface{}
		videoOpt = (data.Video)
		var audioOpt map[string]interface{}
		audioOpt = (data.Audio)
		var opts ffmpeg.Options
		width := data.Video["Width"]
		height := data.Video["Height"]
		resolution := width.(string) + "x" + height.(string)
		for chk := range videoOpt {
			switch chk {
			case "Codec":
				optValue := videoOpt[chk].(string)
				if optValue != "" {
					opts.VideoCodec = &(optValue)
				}
			case "Bitrate":
				optValue := videoOpt[chk].(string)
				if optValue != "" {
					opts.VideoBitRate = &(optValue)
				}
			case "Width":
				optValue := resolution
				if optValue != "" {
					opts.Resolution = &(optValue)
				}
			case "FrameRate":
				optValue := videoOpt[chk].(string)
				if optValue != "" {
					opts.FrameRate = &(optValue)
				}
			case "Profile":
				optValue := videoOpt[chk].(string)
				if optValue != "" {
					opts.VideoProfile = &(optValue)
				}
			case "Level":
				optValue := videoOpt[chk].(string)
				if optValue != "" {
					opts.ProfileLevel = &(optValue)
				}
			case "Crf":
				optValue := videoOpt[chk].(string)
				if optValue != "" {
					opts.Crf = &(optValue)
				}
			case "GOP":
				optValue, _ := videoOpt[chk].(string)
				if optValue != "" {
					v, _ := strconv.Atoi(optValue)
					opts.KeyframeInterval = &(v)
				}
			case "Bframe":
				optValue, _ := videoOpt[chk].(string)
				if optValue != "" {
					v, _ := strconv.Atoi(optValue)
					opts.Bframe = &(v)
				}
			case "Format":
				optValue := videoOpt[chk].(string)
				if optValue != "" {
					opts.OutputFormat = &(optValue)
				}
			}
		}
		for chk := range audioOpt {
			switch chk {
			case "Codec":
				optValue := audioOpt[chk].(string)
				if optValue != "" {
					opts.AudioCodec = &(optValue)
				}
			case "Bitrate":
				optValue := audioOpt[chk].(string)
				if optValue != "" {
					opts.AudioBitrate = &(optValue)
				}
			case "Audiorate":
				optValue := audioOpt[chk].(string)
				if optValue != "" {
					opts.AudioRate = &(optValue)
				}
			}
		}
		inputPath := "Sample.mp4"
		OverlayImg := data.Video["OverlayImage"].(string)
		if OverlayImg == "" {
			args := append([]string{"-i", inputPath, "-threads", "0"}, opts.GetStrArguments()...)
			chkOptions = append(args)
			fmt.Println(chkOptions)
			//output := data.Transcodings[i].Output
			//fmt.Println("ffmpeg", chkOptions)
			c.Send("ffmpeg ", args, " ")
		} else {
			args := append([]string{"-i", inputPath, "-i", OverlayImg, "-threads", "0"}, opts.GetStrArguments()...)
			chkOptions = append(args)
			// fmt.Println(chkOptions)
			//output := data.Transcodings[i].Output
			//fmt.Println("ffmpeg", chkOptions)
			c.Send("ffmpeg ", args, " ")
		}
		// options = append(chkOptions, output)
		defer func() {
			s := recover()
			fmt.Println(s)
		}()
		//+ data.FileName
		// /usr/src/tr/$FILE_NAME
	})

	app.Get("/tr1/2020-10-14/transcode", func(c *fiber.Ctx) {

		jobTemplate := "/usr/src/tr/transcode.json"
		data := preset.Loadjson(jobTemplate)
		defer func() {
			s := recover()
			fmt.Println(s)
		}()

		// ************ org file download Begin ************
		customerName := strings.Split(data.FilePath, "/")[2] // /home/woori/test/sample.mp4
		minioPath := strings.ReplaceAll(data.FilePath, "/home/"+customerName+"/", "")

		downPath := "/usr/src/tr/" + data.FileName // sample.mp4
		uploadIP := data.UploadIP
		// customerName := data.CustomerName
		// fmt.Println(orgFile)
		// file 위치 filepath로 이용해서
		minio.FileDownloader(uploadIP, customerName, minioPath, downPath)

		// cmd := exec.Command()
		// ************ org file download End ************

		if data.DefaultThumbnail["ThumbNum"] != nil {
			thumbnail.MakeThumbnail()
		}

		for i := 0; i < len(data.Presets); i++ {

			if data.Presets[i].Video["Pass"] == nil {
				InfoLogger.Println("Transcoding Start")
				transPass1()
			} else if data.Presets[i].Video["Pass"] == "1" {
				transPass1()
			} else {
				transPass2()
			}
		}

		// ************ result json file upload Begin ************
		InfoLogger.Println("result json file upload Begin")
		res, err := json.Marshal(response.ResponseJson2())
		if err != nil {
			ErrorLogger.Println(err)
			c.Send(err.Error())
		}

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
		resultIP := data.ResultIP

		minio.FileUploader(resultIP, customerName, jasonFile, objectFile)
		InfoLogger.Println("result json file upload End")
		// ************ result json file upload End ************

		// ************ result file upload Begin ************
		for i := 0; i < len(data.Transcodings); i++ {
			defer func() {
				s := recover()
				fmt.Println(s)
			}()

			// customeName := data.CustomerName
			// transcoding result file
			resultFile := "/usr/src/tr/" + data.Transcodings[i].Output
			// minio upload file name
			// full path : /home/woori/sample2.mp4
			replaceHome := strings.ReplaceAll(data.FilePath, "/home/", "")  // woori/test2/sample2.mp4
			splitMp4 := data.FileName                                       // sample2.mp4
			customerUpPath := strings.ReplaceAll(replaceHome, splitMp4, "") // woori/test2/
			uploadPath := customerUpPath + data.Transcodings[i].Output
			// minio.FileUploader(customeName, resultFile, objectFile)
			minio.FileUploader(resultIP, customerName, resultFile, uploadPath)
		}
		// ************ result file upload End************

		// ************ org file upload Begin************
		orgFile := "/usr/src/tr/" + data.FileName
		// minio upload file name
		uploadPath := ".org/" + strings.ReplaceAll(data.FilePath, "/home/", "")

		// uploadPath := ".org/woori/test2/" + data.FileName
		// minio.FileUploader(customeName, orgFile, objectFile)
		minio.OrgFileUploader(resultIP, customerName, orgFile, uploadPath)
		// ************ org file upload End************

		defer func() {
			s := recover()
			fmt.Println(s)
		}()

		// syslog 출력되도록 추가
		InfoLogger.Println("Transcoding End")

		c.Send("complete")
	})

	app.Get("/tr1/2020-10-15/transcode", func(c *fiber.Ctx) {

		jobTemplate := "/usr/src/tr/transcode.json"
		data := preset.Loadjson(jobTemplate)
		defer func() {
			s := recover()
			fmt.Println(s)
		}()

		// ************ result file download Begin ************
		customerName := strings.Split(data.FilePath, "/")[2]                          // /home/woori/test/sample.mp4
		minioPath := strings.ReplaceAll(data.FilePath, "/home/"+customerName+"/", "") // test/sample.mp4
		// sample.mp4
		downPath := "/usr/src/tr/" + data.FileName
		uploadIP := data.UploadIP
		// customerName := data.CustomerName
		// fmt.Println(orgFile)
		minio.FileDownloader(uploadIP, customerName, minioPath, downPath)

		// ftp.CurlDownFtp()
		// ************ result file download End ************

		// fmt.Println("thumbnail : ", data.DefaultThumbnail)
		// if data.DefaultThumbnail != nil {
		if data.DefaultThumbnail["ThumbNum"] != nil {
			thumbnail.MakeThumbnail()
		}
		// }

		for i := 0; i < len(data.Presets); i++ {

			if data.Presets[i].Video["Pass"] == nil {
				fmt.Println("pass nil")
				transPass3()
			} else if data.Presets[i].Video["Pass"] == "1" {
				fmt.Println("pass 1")
				transPass3()
			} else {
				fmt.Println("pass 2")
				transPass2()
			}
		}

		// ************ result json file upload Begin ************
		ftp.CurlJsonUpFtp()
		// ************ result json file upload End ************

		// ************ result file upload Begin ************
		for i := 0; i < len(data.Transcodings); i++ {
			defer func() {
				s := recover()
				fmt.Println(s)
			}()

			resultIP := data.ResultIP
			// customeName := data.CustomerName
			// transcoding result file
			resultFile := "/usr/src/tr/" + data.Transcodings[i].Output
			// minio upload file name
			// full path : /home/woori/sample2.mp4
			replaceHome := strings.ReplaceAll(data.FilePath, "/home/", "")  // woori/test2/sample2.mp4
			splitMp4 := data.FileName                                       // sample2.mp4
			customerUpPath := strings.ReplaceAll(replaceHome, splitMp4, "") // woori/test2/
			uploadPath := customerUpPath + data.Transcodings[i].Output      // woori/test2/sample_1080p.mp4
			// minio.FileUploader(customeName, resultFile, objectFile)
			minio.FileUploader(resultIP, customerName, resultFile, uploadPath)
		}
		// ************ result file upload End************

		// ************ org file upload Begin************
		resultIP := data.ResultIP
		orgFile := "/usr/src/tr/" + data.FileName
		// minio upload file name
		uploadPath := ".org/" + strings.ReplaceAll(data.FilePath, "/home/", "")
		// uploadPath := ".org/woori/test2/" + data.FileName
		// minio.FileUploader(customeName, orgFile, objectFile)
		minio.OrgFileUploader(resultIP, customerName, orgFile, uploadPath)
		// ************ org file upload End************

		defer func() {
			s := recover()
			fmt.Println(s)
		}()

		c.Send("complete")
	})

	app.Listen(3333)
	// app.Listen(2750)

	// func transcoding (ffmpegConf string, inputPath string, output string, opts string){
	// 	progress, err := ffmpeg.
	// 							New(ffmpegConf).
	// 							Input(inputPath).
	// 							Output(output).
	// 							WithOptions(opts).
	// 							Start(opts)

	// 						if err != nil {

	// 							log.Fatal(err)
	// 						}

	// 						for msg := range progress {
	// 							log.Printf(output+"%+v", msg)
	// 						}
	// }
}

func transPass1() {

	jobTemplate := "/usr/src/tr/transcode.json"
	data := preset.Loadjson(jobTemplate)
	// ************ transcoding file Begin ************
	var wg sync.WaitGroup
	chkPreset := make([]string, 0)
	// if data.Presets.Video != null {
	for i := 0; i < len(data.Presets); i++ {
		chkPreset = append(chkPreset, data.Presets[i].PresetName)
		fmt.Println(chkPreset[i])
		if data.Presets[i].Video["OverlayImage"] != nil {
			OverlayImg := data.Presets[i].Video["OverlayImage"].(string) // /home/atlas/ovrimg/test/sample.png

			if data.Presets[i].Video["OverlaySize"] != nil {

				size := data.Presets[i].Video["OverlaySize"].(string)

				customerName := strings.Split(data.FilePath, "/")[2]
				// parseFolder := strings.ReplaceAll(OverlayImg, "/home/"+customerName+"/ovrimg/", "") // /home/atlas/ovrimg/test/sample.mp4 -> test/sample.png
				minioPath := strings.ReplaceAll(OverlayImg, "/home/"+customerName+"/", "")
				// downPath := "/usr/src/tr/" + "atlas_logo.png"

				lenName := strings.SplitAfter(OverlayImg, "/")
				cutName := strings.Split(OverlayImg, "/")[len(lenName)-1]

				downPath := "/usr/src/tr/" + cutName

				fmt.Println("overlayimg downloader")
				downIP := data.ResultIP

				// bucketname := data.CustomerName
				// fmt.Println(orgFile)
				minio.ImageDownloader(downIP, customerName, minioPath, downPath)
				cmd := exec.Command("/usr/src/tr/ffmpeg", "-i", downPath, "-vf", "scale="+size, "-y", "/usr/src/tr/"+data.Presets[i].PresetName+"_img.png")
				cmd.Run()
				fmt.Println(cmd)
			}
		}

	}

	wg.Add(len(data.Transcodings))
	for i := 0; i < len(data.Transcodings); i++ {
		presetName := data.Transcodings[i].PresetName
		output := "/usr/src/tr/" + data.Transcodings[i].Output

		go func() {
			defer wg.Done()
			for i := 0; i < len(chkPreset); i++ {

				fmt.Println(presetName, chkPreset[i])
				if presetName == chkPreset[i] {

					var videoOpt map[string]interface{}
					videoOpt = (data.Presets[i].Video)

					var audioOpt map[string]interface{}
					audioOpt = (data.Presets[i].Audio)

					var opts ffmpeg.Options

					width := data.Presets[i].Video["Width"]
					height := data.Presets[i].Video["Height"]

					resolution := width.(string) + "x" + height.(string)

					// imgLocation := data.Presets[i].Video["OverlayXY"].(string)
					// imgSize := data.Presets[i].Video["OverlaySize"]

					// imgOverlay := imgLocation.(string) + imgSize.(string)

					for chk := range videoOpt {
						// fmt.Println(i)
						switch chk {
						case "Codec":
							optValue := videoOpt[chk].(string)
							opts.VideoCodec = &(optValue)
						case "Bitrate":
							optValue := videoOpt[chk].(string)
							opts.VideoBitRate = &(optValue)
						case "Width":
							optValue := resolution
							opts.Resolution = &(optValue)
						case "FrameRate":
							optValue := videoOpt[chk].(string)
							opts.FrameRate = &(optValue)
						case "Profile":
							optValue := videoOpt[chk].(string)
							opts.VideoProfile = &(optValue)
						case "Level":
							optValue := videoOpt[chk].(string)
							opts.ProfileLevel = &(optValue)
						case "Crf":
							optValue := videoOpt[chk].(string)
							opts.Crf = &(optValue)
						case "OverlayXY":
							optValue := videoOpt[chk].(string)
							opts.Overlay = &(optValue)
						// case "PixelAspect":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.Aspect = &(optValue)
						// case "RateControl":
						// 	optValue := videoOpt[chk].(string)
						// 	if optValue == "CBR" {
						// 		opts.VideoMaxBitRate = &(optValue)
						// 	}
						// 	opts.VideoProfile = &(optValue)
						// case "Pass":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.Pass = &(optValue)
						case "GOP":
							optValue, _ := strconv.Atoi(videoOpt[chk].(string))
							opts.KeyframeInterval = &(optValue)
						case "Bframe":
							optValue, _ := strconv.Atoi(videoOpt[chk].(string))
							opts.Bframe = &(optValue)
						// case "Iframe":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.VideoProfile = &(optValue)
						// case "Interlacing":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.VideoProfile = &(optValue)
						// case "Deinterlacing":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.VideoProfile = &(optValue)
						case "Format":
							optValue := videoOpt[chk].(string)
							opts.OutputFormat = &(optValue)
						}
					}

					for chk := range audioOpt {
						switch chk {
						case "Codec":
							optValue := audioOpt[chk].(string)
							opts.AudioCodec = &(optValue)
						case "Bitrate":
							optValue := audioOpt[chk].(string)
							opts.AudioBitrate = &(optValue)
						}
					}

					overwrite := true
					opts.Overwrite = &overwrite

					inputPath := "/usr/src/tr/" + data.FileName

					ffmpegConf := &ffmpeg.Config{
						FfmpegBinPath:   "/usr/src/tr/ffmpeg",
						FfprobeBinPath:  "/usr/src/tr/ffprobe",
						ProgressEnabled: true,
					}
					// for i := 0; i < len(data.Transcodings); i++ {

					// output := "/usr/src/tr/" + data.Transcodings[i].Output

					if opts.Overlay == nil {
						progress, err := ffmpeg.
							New(ffmpegConf).
							Input(inputPath).
							Output(output).
							WithOptions(opts).
							Start(opts)

						if err != nil {

							log.Fatal(err)
						}

						for msg := range progress {
							log.Printf(output+"%+v", msg)
						}

					} else {
						OverlayImg := "/usr/src/tr/" + data.Presets[i].PresetName + "_img.png"
						progress, err := ffmpeg.
							New(ffmpegConf).
							Input(inputPath).
							InputImage(OverlayImg).
							Output(output).
							WithOptions(opts).
							Start(opts)

						if err != nil {

							log.Println(err)
						}

						for msg := range progress {
							fmt.Printf(output+"%+v", msg)
						}
					}

				}
			}
		}()

	}

	wg.Wait()
	// ************ transcoding file End ************

}

func transPass2() {

	jobTemplate := "/usr/src/tr/transcode.json"
	data := preset.Loadjson(jobTemplate)
	// ************ transcoding file Begin ************
	var wg sync.WaitGroup
	chkPreset := make([]string, 0)

	// if data.Presets.Video != null {
	for i := 0; i < len(data.Presets); i++ {
		chkPreset = append(chkPreset, data.Presets[i].PresetName)
		// fmt.Println(chkPreset[i])

		if data.Presets[i].Video["OverlayImage"] != nil {
			OverlayImg := data.Presets[i].Video["OverlayImage"].(string)

			size := data.Presets[i].Video["OverlaySize"].(string)

			minioPath := OverlayImg
			// sample.mp4
			// bucketname := data.CustomerName
			// fmt.Println(orgFile)

			customerName := strings.Split(data.FilePath, "/")[2]
			lenName := strings.SplitAfter(OverlayImg, "/")
			cutName := strings.Split(OverlayImg, "/")[len(lenName)-1]

			downPath := "/usr/src/tr/" + cutName
			// bucketname := data.CustomerName
			// fmt.Println(orgFile)
			resultIP := data.ResultIP

			minio.ImageDownloader(resultIP, customerName, minioPath, downPath)
			cmd := exec.Command("/usr/src/tr/ffmpeg", "-i", downPath, "-vf", "scale="+size, "-y", "/usr/src/tr/"+data.Presets[i].PresetName+"_img.png")
			cmd.Run()
			fmt.Println(cmd)
		}
	}

	wg.Add(len(data.Transcodings))
	for i := 0; i < len(data.Transcodings); i++ {
		presetName := data.Transcodings[i].PresetName
		passOutput := "/usr/src/tr/pass_" + data.Transcodings[i].Output
		output := "/usr/src/tr/" + data.Transcodings[i].Output

		go func() {
			defer wg.Done()
			for i := 0; i < len(chkPreset); i++ {

				fmt.Println(presetName, chkPreset[i])
				if presetName == chkPreset[i] {

					var videoOpt map[string]interface{}
					videoOpt = (data.Presets[i].Video)

					var audioOpt map[string]interface{}
					audioOpt = (data.Presets[i].Audio)

					var opts ffmpeg.Options

					width := data.Presets[i].Video["Width"]
					height := data.Presets[i].Video["Height"]

					resolution := width.(string) + "x" + height.(string)

					// imgLocation := data.Presets[i].Video["OverlayXY"].(string)
					// imgSize := data.Presets[i].Video["OverlaySize"]

					// imgOverlay := imgLocation.(string) + imgSize.(string)

					// setPreset(videoOpt, audioOpt, resolution, imgLocation)

					for chk := range videoOpt {
						// fmt.Println(i)
						switch chk {
						case "Codec":
							optValue := videoOpt[chk].(string)
							opts.VideoCodec = &(optValue)
						case "Bitrate":
							optValue := videoOpt[chk].(string)
							opts.VideoBitRate = &(optValue)
						case "Width":
							optValue := resolution
							opts.Resolution = &(optValue)
						case "FrameRate":
							optValue := videoOpt[chk].(string)
							opts.FrameRate = &(optValue)
						case "Profile":
							optValue := videoOpt[chk].(string)
							opts.VideoProfile = &(optValue)
						case "level":
							optValue := videoOpt[chk].(string)
							opts.ProfileLevel = &(optValue)
						case "Crf":
							optValue := videoOpt[chk].(string)
							opts.Crf = &(optValue)
						case "OverlayXY":
							optValue := videoOpt[chk].(string)
							opts.Overlay = &(optValue)
						// case "PixelAspect":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.Aspect = &(optValue)
						// case "RateControl":
						// 	optValue := videoOpt[chk].(string)
						// 	if optValue == "CBR" {
						// 		opts.VideoMaxBitRate = &(optValue)
						// 	}
						// 	opts.VideoProfile = &(optValue)
						// case "Pass":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.Pass = &(optValue)
						case "GOP":
							optValue, _ := strconv.Atoi(videoOpt[chk].(string))
							opts.KeyframeInterval = &(optValue)
						case "Bframe":
							optValue, _ := strconv.Atoi(videoOpt[chk].(string))
							opts.Bframe = &(optValue)
						// case "Iframe":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.VideoProfile = &(optValue)
						// case "Interlacing":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.VideoProfile = &(optValue)
						// case "Deinterlacing":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.VideoProfile = &(optValue)
						case "Format":
							optValue := videoOpt[chk].(string)
							opts.OutputFormat = &(optValue)
						}
					}

					for chk := range audioOpt {
						switch chk {
						case "codec":
							optValue := audioOpt[chk].(string)
							opts.AudioCodec = &(optValue)
						case "Bitrate":
							optValue := audioOpt[chk].(string)
							opts.AudioBitrate = &(optValue)
						}
					}

					overwrite := true
					opts.Overwrite = &overwrite

					inputPath := "/usr/src/tr/" + data.FileName

					ffmpegConf := &ffmpeg.Config{
						FfmpegBinPath:   "/usr/src/tr/ffmpeg",
						FfprobeBinPath:  "/usr/src/tr/ffprobe",
						ProgressEnabled: true,
					}
					// for i := 0; i < len(data.Transcodings); i++ {

					Pass1 := "1"
					opts.Pass = &Pass1
					// output := "/usr/src/tr/" + data.Transcodings[i].Output

					if opts.Overlay == nil {
						progress, err := ffmpeg.
							New(ffmpegConf).
							Input(inputPath).
							Output(passOutput).
							WithOptions(opts).
							Start(opts)

						if err != nil {

							log.Fatal(err)
						}

						for msg := range progress {
							log.Printf(output+"%+v", msg)
						}

					} else {
						OverlayImg := "/usr/src/tr/" + data.Presets[i].PresetName + "_img.png"
						progress, err := ffmpeg.
							New(ffmpegConf).
							Input(inputPath).
							InputImage(OverlayImg).
							Output(passOutput).
							WithOptions(opts).
							Start(opts)

						if err != nil {

							log.Fatal(err)
						}

						for msg := range progress {
							log.Printf(output+"%+v", msg)
						}
					}

					Pass2 := "2"
					opts.Pass = &Pass2

					opts.Overlay = (nil)

					// setPreset2Pass(videoOpt, audioOpt, resolution)

					progress, err := ffmpeg.
						New(ffmpegConf).
						Input(passOutput).
						// InputImage(OverlayImg).
						Output(output).
						WithOptions(opts).
						Start(opts)

					if err != nil {

						log.Fatal(err)
					}

					for msg := range progress {
						log.Printf(output+"%+v", msg)
					}

				}
			}
		}()

	}

	wg.Wait()
	// ************ transcoding file End ************

}

func setPreset(videoOpt map[string]interface{}, audioOpt map[string]interface{}, resolution string, imgLocation string) {

	var opts ffmpeg.Options

	for chk := range videoOpt {
		// fmt.Println(i)
		switch chk {
		case "Codec":
			optValue := videoOpt[chk].(string)
			opts.VideoCodec = &(optValue)
		case "Bitrate":
			optValue := videoOpt[chk].(string)
			opts.VideoBitRate = &(optValue)
		case "Width":
			optValue := resolution
			opts.Resolution = &(optValue)
		case "FrameRate":
			optValue := videoOpt[chk].(string)
			opts.FrameRate = &(optValue)
		case "Profile":
			optValue := videoOpt[chk].(string)
			opts.VideoProfile = &(optValue)
		case "level":
			optValue := videoOpt[chk].(string)
			opts.ProfileLevel = &(optValue)
		case "Crf":
			optValue := videoOpt[chk].(string)
			opts.Crf = &(optValue)
		case "OverlayXY":
			optValue := imgLocation
			opts.Overlay = &(optValue)
		// case "PixelAspect":
		// 	optValue := videoOpt[chk].(string)
		// 	opts.Aspect = &(optValue)
		// case "RateControl":
		// 	optValue := videoOpt[chk].(string)
		// 	if optValue == "CBR" {
		// 		opts.VideoMaxBitRate = &(optValue)
		// 	}
		// 	opts.VideoProfile = &(optValue)
		// case "Pass":
		// 	optValue := videoOpt[chk].(string)
		// 	opts.Pass = &(optValue)
		case "GOP":
			optValue, _ := strconv.Atoi(videoOpt[chk].(string))
			opts.KeyframeInterval = &(optValue)
		case "Bframe":
			optValue, _ := strconv.Atoi(videoOpt[chk].(string))
			opts.Bframe = &(optValue)
		// case "Iframe":
		// 	optValue := videoOpt[chk].(string)
		// 	opts.VideoProfile = &(optValue)
		// case "Interlacing":
		// 	optValue := videoOpt[chk].(string)
		// 	opts.VideoProfile = &(optValue)
		// case "Deinterlacing":
		// 	optValue := videoOpt[chk].(string)
		// 	opts.VideoProfile = &(optValue)
		case "Format":
			optValue := videoOpt[chk].(string)
			opts.OutputFormat = &(optValue)
		}
	}

	for chk := range audioOpt {
		switch chk {
		case "codec":
			optValue := audioOpt[chk].(string)
			opts.AudioCodec = &(optValue)
		case "Bitrate":
			optValue := audioOpt[chk].(string)
			opts.AudioBitrate = &(optValue)
		}
	}

	// return opts
}

func setPreset2Pass(videoOpt map[string]interface{}, audioOpt map[string]interface{}, resolution string) {

	var opts ffmpeg.Options

	for chk := range videoOpt {
		// fmt.Println(i)
		switch chk {
		case "Codec":
			optValue := videoOpt[chk].(string)
			opts.VideoCodec = &(optValue)
		case "Bitrate":
			optValue := videoOpt[chk].(string)
			opts.VideoBitRate = &(optValue)
		case "Width":
			optValue := resolution
			opts.Resolution = &(optValue)
		case "FrameRate":
			optValue := videoOpt[chk].(string)
			opts.FrameRate = &(optValue)
		case "Profile":
			optValue := videoOpt[chk].(string)
			opts.VideoProfile = &(optValue)
		case "level":
			optValue := videoOpt[chk].(string)
			opts.ProfileLevel = &(optValue)
		case "Crf":
			optValue := videoOpt[chk].(string)
			opts.Crf = &(optValue)
		// case "PixelAspect":
		// 	optValue := videoOpt[chk].(string)
		// 	opts.Aspect = &(optValue)
		// case "RateControl":
		// 	optValue := videoOpt[chk].(string)
		// 	if optValue == "CBR" {
		// 		opts.VideoMaxBitRate = &(optValue)
		// 	}
		// 	opts.VideoProfile = &(optValue)
		// case "Pass":
		// 	optValue := videoOpt[chk].(string)
		// 	opts.Pass = &(optValue)
		case "GOP":
			optValue, _ := strconv.Atoi(videoOpt[chk].(string))
			opts.KeyframeInterval = &(optValue)
		case "Bframe":
			optValue, _ := strconv.Atoi(videoOpt[chk].(string))
			opts.Bframe = &(optValue)
		// case "Iframe":
		// 	optValue := videoOpt[chk].(string)
		// 	opts.VideoProfile = &(optValue)
		// case "Interlacing":
		// 	optValue := videoOpt[chk].(string)
		// 	opts.VideoProfile = &(optValue)
		// case "Deinterlacing":
		// 	optValue := videoOpt[chk].(string)
		// 	opts.VideoProfile = &(optValue)
		case "Format":
			optValue := videoOpt[chk].(string)
			opts.OutputFormat = &(optValue)
		}
	}

	for chk := range audioOpt {
		switch chk {
		case "codec":
			optValue := audioOpt[chk].(string)
			opts.AudioCodec = &(optValue)
		case "Bitrate":
			optValue := audioOpt[chk].(string)
			opts.AudioBitrate = &(optValue)
		}
	}
}

func transPass3() {

	jobTemplate := "/usr/src/tr/transcode.json"
	data := preset.Loadjson(jobTemplate)
	// ************ transcoding file Begin ************
	var wg sync.WaitGroup
	chkPreset := make([]string, 0)
	// if data.Presets.Video != null {
	for i := 0; i < len(data.Presets); i++ {
		chkPreset = append(chkPreset, data.Presets[i].PresetName)
		fmt.Println(chkPreset[i])
		// if data.Presets[i].Video["OverlayImage"] != nil {
		// 	OverlayImg := data.Presets[i].Video["OverlayImage"].(string) // /home/atlas/ovrimg/test/sample.png

		// 	size := data.Presets[i].Video["OverlaySize"].(string)

		// 	customerName := strings.Split(data.FilePath, "/")[2]
		// 	// parseFolder := strings.ReplaceAll(OverlayImg, "/home/"+customerName+"/ovrimg/", "") // /home/atlas/ovrimg/test/sample.mp4 -> test/sample.png
		// 	minioPath := strings.ReplaceAll(OverlayImg, "/home/"+customerName+"/", "")
		// 	// downPath := "/usr/src/tr/" + "atlas_logo.png"

		// 	lenName := strings.SplitAfter(OverlayImg, "/")
		// 	cutName := strings.Split(OverlayImg, "/")[len(lenName)-1]

		// 	downPath := "/usr/src/tr/" + cutName

		// 	fmt.Println("overlayimg downloader")

		// 	// bucketname := data.CustomerName
		// 	// fmt.Println(orgFile)
		// 	minio.ImageDownloader(customerName, minioPath, downPath)
		// 	cmd := exec.Command("/usr/src/tr/ffmpeg", "-i", downPath, "-vf", "scale="+size, "-y", "/usr/src/tr/"+data.Presets[i].PresetName+"_img.png")
		// 	cmd.Run()
		// 	fmt.Println(cmd)
		// }

	}

	wg.Add(len(data.Transcodings))
	for i := 0; i < len(data.Transcodings); i++ {
		presetName := data.Transcodings[i].PresetName
		output := "/usr/src/tr/" + data.Transcodings[i].Output

		go func() {
			defer wg.Done()
			for i := 0; i < len(chkPreset); i++ {

				fmt.Println(presetName, chkPreset[i])
				if presetName == chkPreset[i] {

					var videoOpt map[string]interface{}
					videoOpt = (data.Presets[i].Video)

					var audioOpt map[string]interface{}
					audioOpt = (data.Presets[i].Audio)

					var opts ffmpeg.Options

					width := data.Presets[i].Video["Width"]
					height := data.Presets[i].Video["Height"]

					resolution := width.(string) + "x" + height.(string)

					// imgLocation := data.Presets[i].Video["OverlayXY"].(string)
					// imgSize := data.Presets[i].Video["OverlaySize"]

					// imgOverlay := imgLocation.(string) + imgSize.(string)

					for chk := range videoOpt {
						// fmt.Println(i)
						switch chk {
						case "Codec":
							optValue := videoOpt[chk].(string)
							opts.VideoCodec = &(optValue)
						case "Bitrate":
							optValue := videoOpt[chk].(string)
							opts.VideoBitRate = &(optValue)
						case "Width":
							optValue := resolution
							opts.Resolution = &(optValue)
						case "FrameRate":
							optValue := videoOpt[chk].(string)
							opts.FrameRate = &(optValue)
						case "Profile":
							optValue := videoOpt[chk].(string)
							opts.VideoProfile = &(optValue)
						case "Level":
							optValue := videoOpt[chk].(string)
							opts.ProfileLevel = &(optValue)
						case "Crf":
							optValue := videoOpt[chk].(string)
							opts.Crf = &(optValue)
						case "OverlayXY":
							optValue := videoOpt[chk].(string)
							opts.Overlay = &(optValue)
						// case "PixelAspect":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.Aspect = &(optValue)
						// case "RateControl":
						// 	optValue := videoOpt[chk].(string)
						// 	if optValue == "CBR" {
						// 		opts.VideoMaxBitRate = &(optValue)
						// 	}
						// 	opts.VideoProfile = &(optValue)
						// case "Pass":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.Pass = &(optValue)
						case "GOP":
							optValue, _ := strconv.Atoi(videoOpt[chk].(string))
							opts.KeyframeInterval = &(optValue)
						case "Bframe":
							optValue, _ := strconv.Atoi(videoOpt[chk].(string))
							opts.Bframe = &(optValue)
						// case "Iframe":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.VideoProfile = &(optValue)
						// case "Interlacing":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.VideoProfile = &(optValue)
						// case "Deinterlacing":
						// 	optValue := videoOpt[chk].(string)
						// 	opts.VideoProfile = &(optValue)
						case "Format":
							optValue := videoOpt[chk].(string)
							opts.OutputFormat = &(optValue)
						}
					}

					for chk := range audioOpt {
						switch chk {
						case "Codec":
							optValue := audioOpt[chk].(string)
							opts.AudioCodec = &(optValue)
						case "Bitrate":
							optValue := audioOpt[chk].(string)
							opts.AudioBitrate = &(optValue)
						}
					}

					overwrite := true
					opts.Overwrite = &overwrite

					inputPath := "/usr/src/tr/" + data.FileName

					ffmpegConf := &ffmpeg.Config{
						FfmpegBinPath:   "/usr/src/tr/ffmpeg",
						FfprobeBinPath:  "/usr/src/tr/ffprobe",
						ProgressEnabled: true,
					}
					// for i := 0; i < len(data.Transcodings); i++ {

					// output := "/usr/src/tr/" + data.Transcodings[i].Output

					if opts.Overlay == nil {
						progress, err := ffmpeg.
							New(ffmpegConf).
							Input(inputPath).
							Output(output).
							WithOptions(opts).
							Start(opts)

						if err != nil {

							log.Fatal(err)
						}

						for msg := range progress {
							log.Printf(output+"%+v", msg)
						}

					} else {
						OverlayImg := "/usr/src/tr/" + data.Presets[i].PresetName + "_img.png"
						progress, err := ffmpeg.
							New(ffmpegConf).
							Input(inputPath).
							InputImage(OverlayImg).
							Output(output).
							WithOptions(opts).
							Start(opts)

						if err != nil {

							log.Fatal(err)
						}

						for msg := range progress {
							log.Printf(output+"%+v", msg)
						}
					}

				}
			}
		}()

	}

	wg.Wait()
	// ************ transcoding file End ************

}
