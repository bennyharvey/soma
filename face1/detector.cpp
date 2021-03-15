#include <dlib/dnn.h>
#include <dlib/opencv.h>
#include <opencv2/core/core.hpp>
#include <opencv2/core/types_c.h>

#include "detector.h"

template <long num_filters, typename SUBNET> using con5d = dlib::con<num_filters,5,5,2,2,SUBNET>;
template <long num_filters, typename SUBNET> using con5  = dlib::con<num_filters,5,5,1,1,SUBNET>;

template <typename SUBNET> using downsampler  = dlib::relu<dlib::affine<con5d<32, dlib::relu<dlib::affine<con5d<32, dlib::relu<dlib::affine<con5d<16,SUBNET>>>>>>>>>;
template <typename SUBNET> using rcon5  = dlib::relu<dlib::affine<con5<45,SUBNET>>>;

using detector_cnn = dlib::loss_mmod<dlib::con<1,9,9,1,1,rcon5<rcon5<rcon5<downsampler<dlib::input_rgb_image_pyramid<dlib::pyramid_down<6>>>>>>>>;

class Detector {

public:
	Detector(const char* model_path) {
		dlib::deserialize(model_path) >> detector;
	}

	std::vector<Detection> detect(const dlib::matrix<dlib::rgb_pixel>& img) {
        std::lock_guard<std::mutex> lock(detector_mutex);
        auto ds = detector(img);

        std::vector<Detection> detections(ds.size());

        for(long unsigned int i = 0; i < ds.size(); i++) {
            detections[i].rectangle.left = ds[i].rect.left();
            detections[i].rectangle.right = ds[i].rect.right();
            detections[i].rectangle.top = ds[i].rect.top();
            detections[i].rectangle.bottom = ds[i].rect.bottom();
            detections[i].confidence = ds[i].detection_confidence;
        }

        return detections;
	}

	std::vector<std::vector<Detection>> batch_detect(const std::vector<dlib::matrix<dlib::rgb_pixel>>& imgs) {
        std::lock_guard<std::mutex> lock(detector_mutex);
        auto ds = detector(imgs);

        std::vector<std::vector<Detection>> detections(ds.size());

        for(long unsigned int i = 0; i < ds.size(); i++) {
            std::vector<Detection> img_detections(ds[i].size());

            for(long unsigned int j = 0; j < ds[i].size(); j++) {
                img_detections[j].rectangle.left = ds[i][j].rect.left();
                img_detections[j].rectangle.right = ds[i][j].rect.right();
                img_detections[j].rectangle.top = ds[i][j].rect.top();
                img_detections[j].rectangle.bottom = ds[i][j].rect.bottom();
                img_detections[j].confidence = ds[i][j].detection_confidence;
            }

            detections[i] = img_detections;
        }

        return detections;
    }

private:
	detector_cnn detector;
	std::mutex detector_mutex;
};

// Plain C interface for Go.

DetectorInitResult* detector_init(const char* model_path) {
	DetectorInitResult* result = (DetectorInitResult*)malloc(sizeof(DetectorInitResult));

	try {
		Detector* fd = new Detector(model_path);
		result->detector = (void*)fd;
		result->err_str = NULL;
	} catch (std::exception& e) {
	    result->detector = NULL;
		result->err_str = strdup(e.what());
	}

	return result;
}

void detector_free(void* detector) {
    delete (Detector*)(detector);
}

DetectorDetectResult* detector_detect(void* detector, const void* image) {
	DetectorDetectResult* result = (DetectorDetectResult*)malloc(sizeof(DetectorDetectResult));

	dlib::matrix<dlib::rgb_pixel> img_mat;

	try {
	    cv::Mat* opencv_img = (cv::Mat*)image;

        if (opencv_img->channels() > 1) {
		    dlib::cv_image<dlib::bgr_pixel> dlib_img(*opencv_img);
            dlib::assign_image(img_mat, dlib_img);
        } else {
            dlib::cv_image<uchar> dlib_img(*opencv_img);
            dlib::assign_image(img_mat, dlib_img);
        }
	} catch (std::exception& e) {
	    result->detections_count = 0;
	    result->detections = NULL;
		result->err_str = strdup(e.what());
		return result;
    }

	auto detections = ((Detector*)(detector))->detect(img_mat);

	Detection* detections_array = (Detection*)calloc(detections.size(), sizeof(Detection));

    for(std::vector<Detection>::size_type i = 0; i < detections.size(); i++) {
        detections_array[i] = detections[i];
    }

    result->detections_count = detections.size();
	result->detections = detections_array;
	result->err_str = NULL;

	return result;
}

DetectorBatchDetectResult* detector_batch_detect(void* detector, const void* images, const int images_count) {
    DetectorBatchDetectResult* result = (DetectorBatchDetectResult*)malloc(sizeof(DetectorBatchDetectResult));

    std::vector<dlib::matrix<dlib::rgb_pixel>> imgs_mats(images_count);

    try {
        cv::Mat** images_arr = (cv::Mat**)images;

        for (int i = 0; i < images_count; i++) {
            cv::Mat* opencv_img = images_arr[i];
            dlib::matrix<dlib::rgb_pixel> img_mat;

            if (opencv_img->channels() > 1) {
                dlib::cv_image<dlib::bgr_pixel> dlib_img(*opencv_img);
                dlib::assign_image(img_mat, dlib_img);
            } else {
                dlib::cv_image<uchar> dlib_img(*opencv_img);
                dlib::assign_image(img_mat, dlib_img);
            }

            imgs_mats[i] = img_mat;
        }
    } catch (std::exception& e) {
        result->detections = NULL;
        result->err_str = strdup(e.what());
        return result;
    }

    auto batch_detections = ((Detector*)(detector))->batch_detect(imgs_mats);

    BatchDetection* batch_detections_array = (BatchDetection*)calloc(batch_detections.size(), sizeof(BatchDetection));

    for(unsigned long i = 0; i < batch_detections.size(); i++) {
        Detection* detections_array = (Detection*)calloc(batch_detections[i].size(), sizeof(Detection));

        for(unsigned long j = 0; j < batch_detections[i].size(); j++) {
            detections_array[j] = batch_detections[i][j];
        }

        batch_detections_array[i].detections_count = batch_detections[i].size();
        batch_detections_array[i].detections = detections_array;
    }

    result->detections = batch_detections_array;
    result->err_str = NULL;

    return result;
}