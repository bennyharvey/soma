package rmq

import (
	"encoding/json"

	"github.com/assembla/cony"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/bennyharvey/soma/entity"
)

type RecognizedFacePublisher struct {
	client    *cony.Client
	publisher *cony.Publisher
	log       *logrus.Entry
}

func NewRecognizedFacePublisher(uri, exchange, passageID string) *RecognizedFacePublisher {
	client := cony.NewClient(cony.URL(uri),
		cony.Backoff(cony.DefaultBackoff))

	client.Declare([]cony.Declaration{
		cony.DeclareExchange(cony.Exchange{
			Name: exchange,
			Kind: exchangeKind,
		}),
	})

	publisher := cony.NewPublisher(exchange, recognizedFacesTopic(passageID))

	client.Publish(publisher)

	log := logrus.WithFields(logrus.Fields{
		"subsystem":  "rmq_recognized_face_publisher",
		"passage_id": passageID,
	})

	go func() {
		for client.Loop() {
			select {
			case err, ok := <-client.Errors():
				if !ok {
					continue
				}
				if err == (*amqp.Error)(nil) {
					continue
				}
				log.WithError(err).Error("got cony client error")
			case blocked, ok := <-client.Blocking():
				if !ok {
					continue
				}
				log.WithField("reason", blocked.Reason).
					Warn("cony client is blocking")
			}
		}
	}()

	return &RecognizedFacePublisher{
		client:    client,
		publisher: publisher,
		log:       log,
	}
}

func (p *RecognizedFacePublisher) Stop() {
	p.publisher.Cancel()
	p.client.Close()
}

func (p *RecognizedFacePublisher) PublishRecognizedFace(rf entity.RecognizedFace) {
	fJSON, err := json.Marshal(rf)
	if err != nil {
		p.log.WithError(err).Error("failed to JSON marshal recognized face1")
		return
	}

	err = p.publisher.Publish(amqp.Publishing{
		Body: fJSON,
	})
	if err != nil {
		p.log.WithError(err).Error("failed to publish recognized face1")
		return
	}
}
