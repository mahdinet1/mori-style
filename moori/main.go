package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"io"
	"log"
	"moori/config"
	embedderClient "moori/delivery/grpc/embedder"
	echoServer "moori/delivery/httpserver"
	meiliCaller "moori/delivery/meilisearch"
	"moori/service/embedder"
	"moori/service/product"
	"moori/service/semanticKeyWordSearch"
	"moori/storage/mysql"
	"moori/storage/mysql/mysqlProduct"
	qdrantDb "moori/storage/qdrant"
	qdrantProduct "moori/storage/qdrant/product"
	"os"
)

// TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func main() {

	mysqlCfg, qdrantCfg, meiliCfg := config.New()
	embClient := embedderClient.New(os.Getenv("EmbedderAddress"))
	embSvc := embedder.New(embClient)
	vectorStorage, err := qdrantDb.New(qdrantCfg)
	if err != nil {
		panic(err)
	}
	qdrantProductStorage := qdrantProduct.New(vectorStorage)

	storage := mysql.New(mysqlCfg)
	mysql.AutoMigrate(storage.Connect())
	qdrantProduct.CreateProductCollection(qdrantProductStorage.Conn())

	jsonFile, _ := os.Open("./cleanData.json")
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var products []map[string]interface{}
	jErr := json.Unmarshal(byteValue, &products)
	if jErr != nil {
		panic(jErr)
	}
	meiliCfg.Document = products
	meiliSearch := meiliCaller.New(meiliCfg)

	semanticSvc := semanticKeyWordSearch.NewService(meiliSearch)

	productSvc := product.New(mysqlProduct.New(storage), qdrantProductStorage, embSvc, semanticSvc)
	server := echoServer.NewServer(echoServer.Config{Port: os.Getenv("EchoServerPort")}, productSvc)
	server.Start()
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
