#pragma once

typedef struct {
	void* detector;
	const char* err_str;
} DetectorInitResult;

typedef struct {
    long left;
    long right;
    long top;
    long bottom;
} Rectangle;

typedef struct {
    Rectangle rectangle;
    double confidence;
} Detection;

typedef struct {
	int detections_count;
	Detection* detections;
	const char* err_str;
} DetectorDetectResult;

typedef struct {
    int detections_count;
	Detection* detections;
} BatchDetection;

typedef struct {
	BatchDetection* detections;
	const char* err_str;
} DetectorBatchDetectResult;

#ifdef __cplusplus
extern "C" {
#endif

DetectorInitResult* detector_init(const char* model_path);
void detector_free(void* detector);
DetectorDetectResult* detector_detect(void* detector, const void* image);
DetectorBatchDetectResult* detector_batch_detect(void* detector, const void* images, const int images_count);

#ifdef __cplusplus
}
#endif