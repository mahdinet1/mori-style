package qdrantDb

import (
	"github.com/qdrant/go-client/qdrant"
	"moori/pkg/richError"
)

type Config struct {
	Host   string
	Port   int
	ApiKey string
	UseTLS bool
}

type Db struct {
	config Config
	db     *qdrant.Client
}

func (m *Db) Connect() *qdrant.Client {
	return m.db
}

func New(config Config) (*Db, error) {
	const op = "qdrant.New"
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   config.Host,
		Port:   config.Port,
		APIKey: config.ApiKey,
		UseTLS: config.UseTLS, // uses default config with minimum TLS version set to 1.3
		// TLSConfig: &tls.Config{...},
		// GrpcOptions: []grpc.DialOption{},
	})
	if err != nil {
		return nil, richError.New(op).
			SetMessage("qdrant connection has failed!!").
			SetCode(richError.UnexpectedCode).
			SetWrappedError(err)
	}
	return &Db{config: config, db: client}, nil
}

func (m *Db) AddNewProduct() {}
