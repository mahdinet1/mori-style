package embedder

import (
	"fmt"
	pb "moori/delivery/grpc/protos"
	"moori/pkg/errormsg"
	"moori/pkg/richError"
)

type Model interface {
	TextToVectorInference(query string) ([]float32, error)
	ImgToVectorInference(url string) ([]*pb.Vector, error)
	BatchImgToVectorInference(urls []string) ([]*pb.Vector, error)
}

type Service struct {
	model Model
}

func New(model Model) *Service {
	return &Service{
		model: model,
	}
}

func (s *Service) TextToVector(query string) ([]float32, error) {
	const op = "embeddersvc.TextToVector."
	vector, err := s.model.TextToVectorInference(query)

	if err != nil {
		fmt.Println(op, err)
		return nil, richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode)
	}
	return vector, nil
}

func (s *Service) ImageToVector(url string) ([][]float32, error) {
	const op = "embeddersvc.ImgToVector"
	vector, err := s.model.ImgToVectorInference(url)
	if err != nil {
		return nil, richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode)
	}

	var result [][]float32
	for _, vec := range vector {
		result = append(result, vec.Vector)
	}
	return result, nil
}

func (s *Service) ImagesToVector(urls []string) ([][]float32, error) {
	const op = "embeddersvc.ImagesToVector"
	vector, err := s.model.BatchImgToVectorInference(urls)

	if err != nil {
		return nil, richError.New(op).
			SetWrappedError(err).
			SetCode(richError.UnexpectedCode).
			SetMessage(errormsg.ModelServiceError)
	}
	var result [][]float32
	for _, vec := range vector {
		result = append(result, vec.Vector)
	}
	return result, nil
}
