package main

import (
	"os/exec"
	"fmt"
	"time"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const path = "/var/lib/postgresql/"

func fullPath(filename string) string {
	return fmt.Sprintf("%s/%s", path, filename)
}

func uploadToS3(filename string) {
	fmt.Println("Enviando backup ao S3...")
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(fullPath(filename))

	if err != nil {
		println(err)
		return
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("backup"),
		Key:    aws.String(filename),
		Body:   file,
	})

	defer file.Close()

	if err != nil {
		println(err)
		return
	}

	fmt.Println(result)
}

func main() {
	cmd := "su"
	arg0 := "-"
	arg1 := "postgres"
	arg2 := "-c"

	app := "brasiliaio_5" // os.Getenv("PSQL-BACKUP-APP")
	db := "brasiliaio_production_5" // os.Getenv("PSQL-BACKUP-DB")

	filename := fmt.Sprintf("pg-%d.dump", time.Now().Unix())
	arg3 := fmt.Sprintf("pg_dump -U %s -w %s > %s", app, db, fullPath(filename))

	fmt.Println("Criando backup", filename)
	exec.Command(cmd, arg0, arg1, arg2, arg3).Output()

	uploadToS3(filename)
}