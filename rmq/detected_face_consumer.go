package rmq

import (
	"encoding/json"
	"sync"

	"github.com/assembla/cony"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/bennyharvey/soma/entity"
)

type DetectedFaceHandler interface {
	HandleDetectedFace(entity.DetectedFace)
}

type DetectedFaceConsumer struct {
	conyClient   *cony.Client
	conyConsumer *cony.Consumer
	wg           sync.WaitGroup
}

func NewDetectedFaceConsumer(uri, exchange, passageID string, dfh DetectedFaceHandler) *DetectedFaceConsumer {

	client := cony.NewClient(cony.URL(uri), cony.Backoff(cony.DefaultBackoff))

	conyExchange := cony.Exchange{
		Name: exchange,
		Kind: exchangeKind,
	}

	queue := &cony.Queue{
		Name:       "",
		Durable:    false,
		AutoDelete: true,
		Exclusive:  true,
	}

	client.Declare([]cony.Declaration{
		cony.DeclareExchange(conyExchange),
		cony.DeclareQueue(queue),
		cony.DeclareBinding(cony.Binding{
			Queue:    queue,
			Exchange: conyExchange,
			Key:      detectedFacesTopic(passageID),
		}),
	})

	consumer := cony.NewConsumer(queue, cony.AutoTag(), cony.AutoAck(),
		cony.Qos(1))

	client.Consume(consumer)

	fc := &DetectedFaceConsumer{
		conyClient:   client,
		conyConsumer: consumer,
	}

	log := logrus.WithFields(logrus.Fields{
		"subsystem":  "rmq_detected_face_consumer",
		"passage_id": passageID,
	})

	fc.wg.Add(1)
	go func() {
		defer fc.wg.Done()
		for client.Loop() {
			select {

			case d, ok := <-fc.conyConsumer.Deliveries():
				if !ok {
					continue
				}

				var df entity.DetectedFace

				err := json.Unmarshal(d.Body, &df)
				if err != nil {
					log.WithError(err).Errorf(
						"failed to JSON unmarshal detected face1")
					continue
				}

				dfh.HandleDetectedFace(df)

			case err, ok := <-consumer.Errors():
				if !ok {
					continue
				}
				if err != nil {
					log.WithError(err).Error("got consumer error")
				}

			case err, ok := <-client.Errors():
				if !ok {
					continue
				}
				if err != (*amqp.Error)(nil) {
					log.WithError(err).Error("got client error")
				}
			}
		}
	}()

	return fc
}

func (fc *DetectedFaceConsumer) Stop() {
	fc.conyConsumer.Cancel()
	fc.conyClient.Close()
	fc.wg.Wait()
}
