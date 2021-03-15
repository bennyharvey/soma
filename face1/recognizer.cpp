#include <dlib/dnn.h>
#include <dlib/opencv.h>
#include <opencv2/core/core.hpp>
#include <opencv2/core/types_c.h>

#include "recognizer.h"

template <template <int,template<typename>class,int,typename> class block, int N, template<typename>class BN, typename SUBNET>
using residual = dlib::add_prev1<block<N,BN,1,dlib::tag1<SUBNET>>>;

template <template <int,template<typename>class,int,typename> class block, int N, template<typename>class BN, typename SUBNET>
using residual_down = dlib::add_prev2<dlib::avg_pool<2,2,2,2,dlib::skip1<dlib::tag2<block<N,BN,2,dlib::tag1<SUBNET>>>>>>;

template <int N, template <typename> class BN, int stride, typename SUBNET>
using block  = BN<dlib::con<N,3,3,1,1,dlib::relu<BN<dlib::con<N,3,3,stride,stride,SUBNET>>>>>;

template <int N, typename SUBNET> using ares      = dlib::relu<residual<block,N,dlib::affine,SUBNET>>;
template <int N, typename SUBNET> using ares_down = dlib::relu<residual_down<block,N,dlib::affine,SUBNET>>;

template <typename SUBNET> using alevel0 = ares_down<256,SUBNET>;
template <typename SUBNET> using alevel1 = ares<256,ares<256,ares_down<256,SUBNET>>>;
template <typename SUBNET> using alevel2 = ares<128,ares<128,ares_down<128,SUBNET>>>;
template <typename SUBNET> using alevel3 = ares<64,ares<64,ares<64,ares_down<64,SUBNET>>>>;
template <typename SUBNET> using alevel4 = ares<32,ares<32,ares<32,SUBNET>>>;

using recognizer_cnn = dlib::loss_metric<dlib::fc_no_bias<128,dlib::avg_pool_everything<
                            alevel0<
                            alevel1<
                            alevel2<
                            alevel3<
                            alevel4<
                            dlib::max_pool<3,3,2,2,dlib::relu<dlib::affine<dlib::con<32,7,7,2,2,
                            dlib::input_rgb_image_sized<150>
                            >>>>>>>>>>>>;

static std::vector<dlib::matrix<dlib::rgb_pixel>> jitter_image(const dlib::matrix<dlib::rgb_pixel>& img, int count)
{
    thread_local dlib::rand rnd;

    std::vector<dlib::matrix<dlib::rgb_pixel>> crops;
    for (int i = 0; i < count; i++) {
        crops.push_back(dlib::jitter_image(img, rnd));
    }

    return crops;
}

class Recognizer {

public:
	Recognizer(const char* shaper_model_path, const char* recognizer_model_path, int jittering) {
	    dlib::deserialize(shaper_model_path) >> shaper;
		dlib::deserialize(recognizer_model_path) >> recognizer;
		jittering_ = jittering;
	}

	float* recognize(const dlib::matrix<dlib::rgb_pixel>& img) {
		auto shape = shaper(img, dlib::rectangle(0, 0, img.nr(), img.nc()));

		dlib::matrix<dlib::rgb_pixel> chip;

		dlib::extract_image_chip(img, dlib::get_face_chip_details(shape, size, padding), chip);

        std::lock_guard<std::mutex> lock(recognizer_mutex);

        dlib::matrix<float,0,1> descriptor_matrix;

        if (jittering_ > 0) {
            descriptor_matrix = dlib::mean(dlib::mat(recognizer(jitter_image(chip, jittering_))));
        } else {
            descriptor_matrix = recognizer(chip);
        }

        float* descriptor = (float*)calloc(descriptor_matrix.nr(), sizeof(float));

        for (int i = 0; i < descriptor_matrix.nr(); i++) {
            descriptor[i] = descriptor_matrix(i,0);
        }

        return descriptor;
	}

private:
    dlib::shape_predictor shaper;

	recognizer_cnn recognizer;
	std::mutex recognizer_mutex;

	int size = 150;
    double padding = 0.25;
    int jittering_ = 0;
};


RecognizerInitResult* recognizer_init(const char* shaper_model_path, const char* recognizer_model_path, int jittering) {
	RecognizerInitResult* result = (RecognizerInitResult*)malloc(sizeof(RecognizerInitResult));

	try {
		Recognizer* fr = new Recognizer(shaper_model_path, recognizer_model_path, jittering);
		result->recognizer = (void*)fr;
		result->err_str = NULL;
	} catch (std::exception& e) {
	    result->recognizer = NULL;
		result->err_str = strdup(e.what());
	}

	return result;
}

void recognizer_free(void* recognizer) {
    delete (Recognizer*)(recognizer);
}

RecognizerRecognizeResult* recognizer_recognize(void* recognizer, const void* image) {
    RecognizerRecognizeResult* result = (RecognizerRecognizeResult*)malloc(sizeof(RecognizerRecognizeResult));

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
	    result->descriptor = NULL;
		result->err_str = strdup(e.what());
		return result;
    }

	result->descriptor = ((Recognizer*)(recognizer))->recognize(img_mat);
	result->err_str = NULL;

	return result;
}