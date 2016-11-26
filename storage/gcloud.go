package storage


import (
	"bytes"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
	"log"
	"net/http"
)

const (
	PROJECT_ID  = "890260246160"
	BUCKET_NAME = "cloudftp"
	STORAGE_URL = "https://storage.googleapis.com/%s/%s/%s/%s"
)

type GCloudFileService struct {
	ctx     context.Context
	client  *http.Client
	service *storage.Service
}

func NewGCloudFileService() *GCloudFileService {
	gCloudService := &GCloudFileService{}
	gCloudService.ctx = context.Background()
	var err error
	gCloudService.client, err = google.DefaultClient(gCloudService.ctx, storage.CloudPlatformScope)
	if err != nil {
		log.Fatalf(" [ERROR] Error while initializing client for gstorage. %s", err)
	}

	gCloudService.service, err = storage.New(gCloudService.client)
	if err != nil {
		log.Fatalf(" [ERROR] Error while creating gstorage service %s", err)
	}
	return gCloudService
}

func (gcs *GCloudFileService) Save(fileContent []byte, dstFilePath string) error {
	return gcs.putFileInGStorage(BUCKET_NAME, fileContent, dstFilePath)
}

func (gcs *GCloudFileService) putFileInGStorage(bucketName string, fContent []byte, dstFPath string) error {
	fileObject := &storage.Object{Name: dstFPath}
	fReader := bytes.NewReader(fContent)
	_, err := gcs.service.Objects.Insert(bucketName, fileObject).Media(fReader).Do()
	return err
}