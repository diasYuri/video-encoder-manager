package services_test

import (
	"github.com/diasYuri/encoder-go/application/repositories"
	"github.com/diasYuri/encoder-go/application/services"
	"github.com/diasYuri/encoder-go/domain"
	"github.com/diasYuri/encoder-go/framework/database"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"github.com/joho/godotenv"
)

func init(){
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb){
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.FilePath = "test.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}

	return video, repo
}

func TestVideoServiceDownload(t *testing.T)  {
	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("bucketTest")
	require.Nil(t, err)
}

func TestVideoServiceFragment(t *testing.T)  {
	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("bucketTest")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)
}