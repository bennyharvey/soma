package rmq

import (
	"encoding/json"
	"sync"

	"github.com/assembla/cony"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/bennyharvey/soma/entity"
)

type RecognizedFaceHandler interface {
	HandleRecognizedFace(entity.RecognizedFace)
}

type RecognizedFaceConsumer struct {
	conyClient   *cony.Client
	conyConsumer *cony.Consumer
	wg           sync.WaitGroup
}

func NewRecognizedFaceConsumer(uri, exchange, passageID string, rfh RecognizedFaceHandler) *RecognizedFaceConsumer {

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
			Key:      recognizedFacesTopic(passageID),
		}),
	})

	consumer := cony.NewConsumer(queue, cony.AutoTag(), cony.AutoAck(),
		cony.Qos(1))

	client.Consume(consumer)

	fc := &RecognizedFaceConsumer{
		conyClient:   client,
		conyConsumer: consumer,
	}

	log := logrus.WithFields(logrus.Fields{
		"subsystem":  "rmq_recognized_face_consumer",
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

				var rf entity.RecognizedFace

				err := json.Unmarshal(d.Body, &rf)
				if err != nil {
					log.WithError(err).Errorf(
						"failed to JSON unmarshal recognized face1")
					continue
				}

				rfh.HandleRecognizedFace(rf)

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

func (fc *RecognizedFaceConsumer) Stop() {
	fc.conyConsumer.Cancel()
	fc.conyClient.Close()
	fc.wg.Wait()
}
