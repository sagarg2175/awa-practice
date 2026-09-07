package handler

import (
	"fmt"
	"io"
	"main/domain"
	"main/repo"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

type AwsProjHandler struct {
	svc        *repo.AwsProjRepo
	uploader   *s3manager.Uploader
	bucketName string
	region     string
}

func NewAwsProjHandler(svc *repo.AwsProjRepo, uploader *s3manager.Uploader, bucketName, region string) *AwsProjHandler {
	return &AwsProjHandler{
		svc,
		uploader,
		bucketName,
		region,
	}
}

type CreateRequest struct {
	Name   string `json:"name" validate:"required,customName"`
	Domain string `json:"domain" validate:"required,min=8"`
}

func (ah *AwsProjHandler) AwaProjCreate(ctx *gin.Context) {

	var req CreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}

	batch := domain.Master{
		Name:   &req.Name,
		Domain: &req.Domain,
	}

	err := ah.svc.CreateAwsProj(ctx, &batch)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": []string{"User created successfully"},
	})

}

func (ah *AwsProjHandler) AwsProjUpload(ctx *gin.Context) {

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": []string{"Error occurred while processing form data"},
		})
		return
	}

	files := form.File["files"]

	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": []string{"No files provided"},
		})
		return
	}

	var uploadedFiles []string
	var errors []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		url, err := saveFile(
			ah.uploader,
			ah.bucketName,
			ah.region,
			file,
			fileHeader,
		)

		defer file.Close()
		if err != nil {
			errors = append(errors, "error Uploading"+fileHeader.Filename+" : "+err.Error())
		}
		uploadedFiles = append(uploadedFiles, url)
	}
	if len(errors) > 0 {
		ctx.JSON(http.StatusPartialContent, gin.H{
			"success": false,
			"message": []string{"Some files could not be uploaded"},
			"files":   uploadedFiles,
			"errors":  errors,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": []string{"Files uploaded successfully"},
		"files":   uploadedFiles,
	})
}

func saveFile(uploader *s3manager.Uploader, bucketName string, region string,
	fileReader io.Reader,
	fileHeader *multipart.FileHeader) (string, error) {

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileHeader.Filename),
		Body:        fileReader,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		bucketName,
		region,
		fileHeader.Filename,
	)

	return url, nil

}
