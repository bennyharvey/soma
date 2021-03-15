#pragma once

typedef struct {
	void* recognizer;
	const char* err_str;
} RecognizerInitResult;

const int FACE_DESCRIPTOR_SIZE = 128;

typedef struct {
	float* descriptor;
	const char* err_str;
} RecognizerRecognizeResult;

#ifdef __cplusplus
extern "C" {
#endif

RecognizerInitResult* recognizer_init(const char* shaper_model_path, const char* recognizer_model_path, int jittering);
void recognizer_free(void* recognizer);
RecognizerRecognizeResult* recognizer_recognize(void* recognizer, const void* image);

#ifdef __cplusplus
}
#endif