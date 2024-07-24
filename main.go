package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	// "github.com/aws/aws-sdk-go-v2/aws"
	// "github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	//go:embed all:templates/*
	templateFS embed.FS

	//go:embed css/output.css
	css embed.FS

	//parsed templates
	html *template.Template
)

func main() {
	handleSigTerms()

	var err error
	html, err = TemplateParseFSRecursive(templateFS, ".html", true, nil)
	if err != nil {
		log.Fatal(err, "failed to parse templates")
	}

	router := http.NewServeMux()
	router.Handle("/css/output.css", http.FileServer(http.FS(css)))

	router.Handle("/company/add", Action(companyAdd))
	router.Handle("/company/add/", Action(companyAdd))

	router.Handle("/company/edit", Action(companyEdit))
	router.Handle("/company/edit/", Action(companyEdit))

	router.Handle("/company", Action(companies))
	router.Handle("/company/", Action(companies))

	router.Handle("/", Action(index))
	router.Handle("/index.html", Action(index))

	//logging/tracing
	var f *os.File
	log_dir := os.Getenv("LOG_DIR")
	if log_dir != "" {
		if _, err := os.Stat(log_dir); os.IsNotExist(err) {
			err = os.Mkdir(log_dir, 0755)
			if err != nil {
				log.Fatalf("error creating log dir: %v", err)
			}
		}

		logpath := getLogFilePath(log_dir)
		f, err = os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening log file: %v", err)
		}
		defer f.Close()
	} else {
		f = os.Stdout
	}
	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	logger := log.New(f, "http: ", log.LstdFlags)
	middleware := tracing(nextRequestID)(logging(logger)(router))

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	logger.Println("listening on port " + port)
	if err := http.ListenAndServe(":"+port, middleware); err != nil {
		logger.Println("http.ListenAndServe():", err)
		os.Exit(1)
	}
}

// exit process immediately upon sigterm
func handleSigTerms() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("received SIGTERM, exiting")
		os.Exit(1)
	}()
}

// 	cfg, err := config.LoadDefaultConfig(context.TODO())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	s3client := s3.NewFromConfig(cfg)
// 	logfilename := filepath.Base(logpath)
// 	uploadFile(s3client, logpath, "silas-s3-bucket", fmt.Sprintf("application-logs/%s", logfilename))

//	func uploadFile(s3Client *s3.Client, fileName string, bucketName string, objectKey string) error {
//		file, err := os.Open(fileName)
//		if err != nil {
//			log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
//		} else {
//			defer file.Close()
//			_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
//				Bucket: aws.String(bucketName),
//				Key:    aws.String(objectKey),
//				Body:   file,
//			})
//			if err != nil {
//				log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
//					fileName, bucketName, objectKey, err)
//			}
//		}
//		return err
//	}
