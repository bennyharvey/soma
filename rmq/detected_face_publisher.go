package rmq

import (
	"encoding/json"

	"github.com/assembla/cony"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/bennyharvey/soma/entity"
)

type DetectedFacePublisher struct {
	client    *cony.Client
	publisher *cony.Publisher
	log       *logrus.Entry
}

func NewDetectedFacePublisher(uri, exchange, passageID string) *DetectedFacePublisher {
	client := cony.NewClient(cony.URL(uri),
		cony.Backoff(cony.DefaultBackoff))

	client.Declare([]cony.Declaration{
		cony.DeclareExchange(cony.Exchange{
			Name: exchange,
			Kind: exchangeKind,
		}),
	})

	publisher := cony.NewPublisher(exchange, detectedFacesTopic(passageID))

	client.Publish(publisher)

	log := logrus.WithFields(logrus.Fields{
		"subsystem":  "rmq_detected_face_publisher",
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

	return &DetectedFacePublisher{
		client:    client,
		publisher: publisher,
		log:       log,
	}
}

func (p *DetectedFacePublisher) Stop() {
	p.publisher.Cancel()
	p.client.Close()
}

func (p *DetectedFacePublisher) PublishDetectedFace(df entity.DetectedFace) {
	fJSON, err := json.Marshal(df)
	if err != nil {
		p.log.WithError(err).Error("failed to JSON marshal detected face1")
		return
	}

	err = p.publisher.Publish(amqp.Publishing{
		Body: fJSON,
	})
	if err != nil {
		p.log.WithError(err).Error("failed to publish detected face1")
		return
	}
}
