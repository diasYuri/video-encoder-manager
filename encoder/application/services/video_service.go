package services

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/diasYuri/encoder-go/application/repositories"
	"github.com/diasYuri/encoder-go/domain"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
)

type VideoService struct {
	Video *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService{
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(v.Video.FilePath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}

	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	log.Printf("video %v has benn stored", v.Video.ID)

	return nil
}

func (v *VideoService) Fragment() error  {

	err := os.Mkdir(os.Getenv("localStoragePath")+"/"+v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+v.Video.ID+".frag")
	cmdArgs = append(cmdArgs, "--user-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+v.Video.ID)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error{

	err := os.Remove(os.Getenv(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4"))
	if err != nil {
		log.Printf("Error removing mp4 ", v.Video.ID, ".mp4")
		return err
	}

	err = os.Remove(os.Getenv(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag"))
	if err != nil {
		log.Printf("Error removing frag ", v.Video.ID, ".mp4")
		return err
	}

	err = os.RemoveAll(os.Getenv(os.Getenv("localStoragePath") + "/" + v.Video.ID))
	if err != nil {
		log.Printf("Error removing folder ", v.Video.ID)
		return err
	}


	log.Printf("Files have been removed: "+ v.Video.ID)
	return nil
}

func (v *VideoService) InsertVideo() error {
	_, err := v.VideoRepository.Insert(v.Video)
	if err != nil {
		return err
	}

	return nil
}


func printOutput(out []byte){
	if len(out) > 0{
		log.Printf("=====> Output: %s\n", string(out))
	}
}