package rmq

import "fmt"

const exchangeKind = "topic"

func detectedFacesTopic(passageID string) string {
	return fmt.Sprintf("streams.%s.detected_faces", passageID)
}

func recognizedFacesTopic(passageID string) string {
	return fmt.Sprintf("streams.%s.recognized_faces", passageID)
}
