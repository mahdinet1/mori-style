package embedderClient

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "moori/delivery/grpc/protos"
	"moori/pkg/errormsg"
	"moori/pkg/richError"
	"time"
)

type Client struct {
	address     string
	serviceConn pb.EmbedderClient
	context     context.Context
}
type Vector struct {
	vector []float32
}

func New(address string) *Client {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	c := pb.NewEmbedderClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Hour)
	//defer cancel()

	return &Client{
		address:     address,
		serviceConn: c,
		context:     ctx,
	}
}

func (client *Client) TextToVectorInference(query string) ([]float32, error) {
	const op = "embedderClient.TextToVectorInference"

	r, err := client.serviceConn.ReturnTextVector(client.context, &pb.TextToVectorRequest{Query: query})
	if err != nil {
		fmt.Println(op, err.Error())
		return nil, richError.New(op).SetCode(richError.UnexpectedCode).SetWrappedError(err)
	}
	return r.Vector, nil
}

func (client *Client) ImgToVectorInference(url string) ([]*pb.Vector, error) {
	const op = "embedderClient.ImgToVectorInference"
	urls := []string{url}
	r, err := client.serviceConn.ReturnImageVector(client.context, &pb.ImageVectorRequest{ImageUrl: urls})
	if err != nil {
		return nil, richError.New(op).SetCode(richError.UnexpectedCode).SetWrappedError(err)
	}
	return r.Vectors, nil
}

func (client *Client) BatchImgToVectorInference(urls []string) ([]*pb.Vector, error) {
	const op = "embedderClient.BatchImgToVectorInference"
	r, err := client.serviceConn.ReturnImageVector(client.context, &pb.ImageVectorRequest{ImageUrl: urls})
	if err != nil {
		log.Println(op, err.Error())
		return nil, richError.New(op).SetCode(richError.UnexpectedCode).SetWrappedError(err).SetMessage(errormsg.ModelServiceError)
	}
	return r.Vectors, nil
}
