package domain_test

import (
	"github.com/diasYuri/encoder-go/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestValidateIfVideoIsEmpty(t *testing.T){
	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)
	require.Empty(t, video)
}

func TestVideoIdIsNotAUuid(t *testing.T){
	video := domain.NewVideo()

	video.ID = "abc"

	err := video.Validate()
	require.Error(t, err)
}

func TestVideoValidation(t *testing.T){
	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.CreatedAt = time.Now()
	video.ResourceID = "a"
	video.FilePath = "path"

	err := video.Validate()
	require.Nil(t, err)
}