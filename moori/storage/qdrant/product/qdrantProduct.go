package qdrantProduct

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
	"log"
	"moori/entity"
	"moori/pkg/richError"
	"moori/service/product"
	qdrantDb "moori/storage/qdrant"
)

type ProductDB struct {
	conn *qdrantDb.Db
}

func New(conn *qdrantDb.Db) *ProductDB {

	return &ProductDB{conn: conn}
}
func (db *ProductDB) Conn() *qdrant.Client {
	return db.conn.Connect()
}
func CreateProductCollection(client *qdrant.Client) {
	collections, err := client.ListCollections(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(collections)
	log.Println("Creating product collection")
	exist, err := client.CollectionExists(context.Background(), "product")
	if err != nil {
		log.Fatal("CollectionExists has failed!: ", err, exist)
	}
	if !exist {
		err = client.CreateCollection(context.Background(), &qdrant.CreateCollection{
			CollectionName: "product",
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     512,
				Distance: qdrant.Distance_Cosine,
			}),
		})
		if err != nil {
			log.Fatalln("CreateCollection has failed!", err)
		}
	}

}
func (db *ProductDB) AddNewProduct(productId uint, vector []float32, imageId uint64) error {
	const op = "qdrant.ProductDB.AddNewProduct"

	_, err := db.conn.Connect().Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "product",
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(imageId),
				Payload: qdrant.NewValueMap(map[string]any{"product_id": productId}),
				Vectors: qdrant.NewVectorsDense(vector),
			},
		},
	})
	if err != nil {
		return richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode)
	}
	return nil
}
func (db *ProductDB) SearchByTextVector(vector []float32) ([]entity.ProductID, error) {
	const op = "qdrant.ProductDB.SearchByTextVector"
	query, err := db.conn.Connect().Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "product",
		WithPayload:    qdrant.NewWithPayload(true),
		Query:          qdrant.NewQueryDense(vector),
	})
	if err != nil {
		fmt.Println(op, err.Error())
		return nil, richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode)
	}

	var products []entity.ProductID
	uniqueIDs := make(map[entity.ProductID]struct{}) // Map to track unique product IDs
	for _, point := range query {

		payload := point.Payload
		prId := entity.ProductID(payload["product_id"].GetIntegerValue())

		// Check if the ID is already in the map
		if _, exists := uniqueIDs[prId]; !exists {
			// Add the ID to the map and append to products
			uniqueIDs[prId] = struct{}{}
			products = append(products, prId)
		}
	}

	return products, nil
}

func (db *ProductDB) AddBulkyProduct(pr []product.ProductsVector) error {
	const op = "qdrant.ProductDB.AddBulkyProduct"
	var points []*qdrant.PointStruct
	for _, productsVector := range pr {
		for _, vec := range productsVector.Vector {
			fmt.Println("product_id", productsVector.ProductId)
			point := &qdrant.PointStruct{
				Id:      qdrant.NewIDUUID(uuid.New().String()),
				Payload: qdrant.NewValueMap(map[string]any{"product_id": productsVector.ProductId}),
				Vectors: qdrant.NewVectorsDense(vec),
			}
			points = append(points, point)
		}

	}
	_, err := db.conn.Connect().Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "product",
		Points:         points,
	})
	if err != nil {
		fmt.Println(op, err.Error())
		return richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode)
	}
	return nil
}
