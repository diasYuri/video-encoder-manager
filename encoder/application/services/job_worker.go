package services

import (
	"github.com/diasYuri/encoder-go/domain"
	"github.com/diasYuri/encoder-go/framework/utils"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job domain.Job
	Message *amqp.Delivery
	Error error
}

func JobWorker(messageChannel chan amqp.Delivery, returnChannel chan JobWorkerResult, jobService JobService, workerID int){

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