package services

import (
	"encoding/json"
	"github.com/diasYuri/encoder-go/domain"
	"github.com/diasYuri/encoder-go/framework/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"os"
	"sync"
	"time"
)

type JobWorkerResult struct {
	Job domain.Job
	Message *amqp.Delivery
	Error error
}

var Mutex = &sync.Mutex{}

func JobWorker(messageChannel chan amqp.Delivery, returnChannel chan JobWorkerResult, jobService JobService, job domain.Job, workerID int){

	//{
	//	"resource_id": "id_de_controle",
	//	"file_path": "exemplo.mp4"
	//}

	for message := range messageChannel {
		err := utils.IsJson(string(message.Body))
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		Mutex.Lock()
		err = json.Unmarshal(message.Body, &jobService.VideoService.Video)
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.VideoService.Video.ID = uuid.NewV4().String()
		Mutex.Unlock()

		err = jobService.VideoService.Video.Validate()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		Mutex.Lock()
		err = jobService.VideoService.InsertVideo()
		Mutex.Unlock()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}


		job.Video = jobService.VideoService.Video
		job.OutputBucketPath = os.Getenv("outputBucketName")
		job.ID = uuid.NewV4().String()
		job.Status = "STARTING"
		job.CreatedAt = time.Now()


		Mutex.Lock()
		_, err = jobService.JobRepository.Insert(&job)
		Mutex.Unlock()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.Job = &job

		err = jobService.Start()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		returnChannel <- returnJobResult(job, message, nil)
	}

}

func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult  {
	result := JobWorkerResult{
		Job: job,
		Message: &message,
		Error: err,
	}

	return result
}