package minio

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/org_transcoder/transcoder/preset"
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

// FileDownloader by using minio
func FileDownloader(uploadIP string, customerName string, minioPath string, downPath string) {

	jobTemplate := "/usr/src/tr/transcode.json"
	data := preset.Loadjson(jobTemplate)

	ctx := context.Background()
	endpoint := uploadIP + ":9000"
	// endpoint := "trupload.myskcdn.net:9000"
	accessKeyID := "admin"
	secretAccessKey := "!Asdf1209"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Println(err)
	} else {
		InfoLogger.Println("Successfully minio is initialized")
	}

	// Make a new bucket (/home)
	// bucketName := "woori"
	bucketName := customerName

	// Upload the file
	objectName := minioPath
	filePath := downPath

	// Upload the file with FPutObject
	err = minioClient.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		log.Println(data.FileName, err)
		ErrorLogger.Println(data.FileName, err)
	} else {
		InfoLogger.Println("Successfully download ", objectName)
		log.Printf("Successfully download %s\n", objectName)
	}
}

// FileUploader by using minio sdk (https://docs.min.io/docs/golang-client-quickstart-guide.html)
func FileUploader(resultIP string, customerName string, resultFile string, uploadPath string) {
	ctx := context.Background()
	endpoint := resultIP + ":9000"
	// endpoint := "upload.myskcdn.net:9000"
	accessKeyID := "admin"
	secretAccessKey := "!Asdf1209"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {

		defer func() {
			s := recover()
			fmt.Println(s)
		}()
		ErrorLogger.Println("minio server dead")
		panic("minio server dead")
		// log.Fatalln(err)
	}

	// Make a new bucket called mymusic.
	// bucketName := "woori"
	bucketName := customerName
	location := "asia-east-1"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Println(err)
			ErrorLogger.Println(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
		InfoLogger.Println("Successfully created", bucketName)
	}

	// Upload the zip file
	// var objectPath string
	var contentType string
	if strings.Split(uploadPath, ".")[1] == "mp4" {
		// objectPath = "woori/test2/" + uploadPath
		contentType = "video/mp4"
	} else {
		// objectPath = ".json/woori/test2/" + uploadPath
		contentType = "application/json"
	}
	objectPath := uploadPath
	filePath := resultFile
	// contentType := "application/vnd.apple.mpegurl"

	// Upload the file with FPutObject
	n, err := minioClient.FPutObject(ctx, bucketName, objectPath, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Println(err)
		ErrorLogger.Println(err)
	} else {
		// cmd := exec.Command("rm", "-f", resultFile)
		// cmd.Start()
		// log.Printf(resultFile + " was deleted")
	}

	// UploadInfo contains information about the
	// newly uploaded or copied object.
	// type UploadInfo struct {
	// 	Bucket       string
	// 	Key          string
	// 	ETag         string
	// 	Size         int64
	// 	LastModified time.Time
	// 	Location     string
	// 	VersionID    string

	// 	// Lifecycle expiry-date and ruleID associated with the expiry
	// 	// not to be confused with `Expires` HTTP header.
	// 	Expiration       time.Time
	// 	ExpirationRuleID string
	// }

	log.Printf("Successfully uploaded %s of size %d\n", objectPath, n.Size)
	InfoLogger.Println("Successfully uploaded", objectPath, "of size", n.Size)
}

// ImageDownloader by using minio sdk
func ImageDownloader(downIP string, customerName string, minioPath string, downPath string) {
	ctx := context.Background()
	endpoint := downIP + ":9000"
	// endpoint := "upload.myskcdn.net:9000"
	accessKeyID := "admin"
	secretAccessKey := "!Asdf1209"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Println(err)
	}

	// Make a new bucket (/home)
	bucketName := customerName

	// Upload the file
	objectName := minioPath
	filePath := downPath

	defer func() {
		s := recover()
		fmt.Println(s)
	}()

	// Upload the file with FPutObject
	err = minioClient.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		defer func() {
			s := recover()
			log.Println(s)
			// log.Fatalln(s)
		}()
		// log.Fatalln(err)
		panic(err)
	}

	log.Println("Successfully download: ", objectName)
}

// OrgFileUploader by using minio sdk (https://docs.min.io/docs/golang-client-quickstart-guide.html)
func OrgFileUploader(resultIP string, customerName string, orgFile string, uploadPath string) {
	ctx := context.Background()
	endpoint := resultIP + ":9000"
	// endpoint := "upload.myskcdn.net:9000"
	accessKeyID := "admin"
	secretAccessKey := "!Asdf1209"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called mymusic.
	// bucketName := "woori"
	bucketName := customerName
	location := "asia-east-1"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	objectPath := uploadPath
	contentType := "video/mp4"
	filePath := orgFile
	// contentType := "application/vnd.apple.mpegurl"

	// Upload the zip file with FPutObject
	n, err := minioClient.FPutObject(ctx, bucketName, objectPath, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	// UploadInfo contains information about the
	// newly uploaded or copied object.
	// type UploadInfo struct {
	// 	Bucket       string
	// 	Key          string
	// 	ETag         string
	// 	Size         int64
	// 	LastModified time.Time
	// 	Location     string
	// 	VersionID    string

	// 	// Lifecycle expiry-date and ruleID associated with the expiry
	// 	// not to be confused with `Expires` HTTP header.
	// 	Expiration       time.Time
	// 	ExpirationRuleID string
	// }

	log.Printf("Successfully uploaded %s of size %d\n", objectPath, n.Size)
}

func FileList() {
	ctx := context.Background()
	endpoint := "truploader.myskcdn.net:9000"
	accessKeyID := "admin"
	secretAccessKey := "!Asdf1209"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	bucketName := "woori"

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		// Prefix:    "myprefix",
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		fmt.Println(object.Key)
		// fmt.Println(object.ETag)
	}
}
