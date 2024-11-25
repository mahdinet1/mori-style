package meilisearch

import (
	"github.com/meilisearch/meilisearch-go"
	"log"
	"moori/pkg/richError"
)

type Config struct {
	Address  string `json:"address"`
	ApiKey   string `json:"api_key"`
	Index    string `json:"index"`
	Document []map[string]interface{}
}

type Client struct {
	config Config
	client meilisearch.ServiceManager
}

func New(cfg Config) *Client {
	client := meilisearch.New(cfg.Address, meilisearch.WithAPIKey(cfg.ApiKey))
	_, err := client.Index(cfg.Index).UpdateIndex("id")
	if err != nil {
		log.Println("meilisearch update index has error", err.Error())
		return nil
	}

	_, err = client.Index(cfg.Index).AddDocuments(cfg.Document)

	if err != nil {
		panic(err)
	}
	return &Client{
		config: cfg,
		client: client,
	}
}

func (cli *Client) Search(keyword string) ([]interface{}, error) {
	const op = "meili.search"
	search, err := cli.client.Index(cli.config.Index).Search(keyword, &meilisearch.SearchRequest{})
	if err != nil {
		log.Printf("meili search error %+v \n", err)
		return nil, richError.New(op).
			SetWrappedError(err).
			SetCode(richError.UnexpectedCode).
			SetMessage("meili search error!")
	}

	return search.Hits, nil
}
