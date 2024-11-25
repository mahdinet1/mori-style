package product

import (
	"fmt"
	"moori/entity"
	"moori/pkg/errormsg"
	"moori/pkg/richError"
)

type FilterProductsFields struct {
	ID        []entity.ProductID
	Keyword   string
	FromPrice entity.CurrentPrice
	EndPrice  entity.CurrentPrice
}

type Storage interface {
	CreateProduct(product entity.Product) (entity.ProductID, error)
	AddBulky(product []entity.Product) error
	Filter(req FilterProductsFields) ([]entity.Product, error)
}
type EmbeddedService interface {
	TextToVector(text string) ([]float32, error)
	ImageToVector(url string) ([][]float32, error)
	ImagesToVector(url []string) ([][]float32, error)
}

type VectorStorage interface {
	AddNewProduct(productId uint, vector []float32, imageId uint64) error
	SearchByTextVector(vector []float32) ([]entity.ProductID, error)
	AddBulkyProduct(data []ProductsVector) error
}
type ProductsVector struct {
	ProductId entity.ProductID `json:"product_id"`
	Vector    [][]float32      `json:"vector"`
}
type KeyWordSearchSvc interface {
	SearchByText(keyword string) ([]interface{}, error)
}
type Service struct {
	storage       Storage
	vectorStorage VectorStorage
	embedder      EmbeddedService
	keywordSearch KeyWordSearchSvc
}

type CreateProductRequest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	CurrentPrice float32  `json:"current_price"`
	OffPercent   uint     `json:"off_percent"`
	Images       []string `json:"images"`
}
type CreateProductResponse struct {
	Product entity.Product `json:"product"`
	Message string         `json:"message"`
}

func New(storage Storage, vectorStorage VectorStorage, embedder EmbeddedService, keywordSearch KeyWordSearchSvc) Service {
	return Service{
		storage:       storage,
		vectorStorage: vectorStorage,
		embedder:      embedder,
		keywordSearch: keywordSearch,
	}
}

func (s *Service) CreateProduct(req CreateProductRequest) (CreateProductResponse, error) {
	const op = "service.CreateProduct"

	// convert image string to image struct
	var images []entity.Image
	for _, image := range req.Images {
		images = append(images, entity.Image{
			ID:        0,
			Address:   image,
			ProductId: 0,
		})
	}

	productId, err := s.storage.CreateProduct(entity.Product{
		ID:           0,
		Name:         req.Name,
		Description:  req.Description,
		OffPercent:   req.OffPercent,
		CurrentPrice: req.CurrentPrice,
		Images:       images,
	})

	if err != nil {
		return CreateProductResponse{}, richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode)
	}

	var vectors [][]float32
	vectors, err = s.embedder.ImagesToVector(req.Images)
	if err != nil {
		return CreateProductResponse{}, richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode)
	}
	// TODO - fix image id logic
	fmt.Println(vectors)
	//err = s.vectorStorage.AddNewProduct(productId, vectors)
	//if err != nil {
	//	return CreateProductResponse{}, richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode)
	//}

	return CreateProductResponse{
		Product: entity.Product{
			ID:           productId,
			Name:         req.Name,
			Description:  req.Description,
			OffPercent:   req.OffPercent,
			CurrentPrice: req.CurrentPrice,
			Images:       images,
		},
		Message: "product created",
	}, nil
}

func (s *Service) ListProducts() ([]entity.Product, error) {
	//const op = "service.ListProducts"
	// Todo - fix this with wright query base on product id
	//products, err := s.storage.ListProducts()
	//if err != nil {
	//	return nil, richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode)
	//}

	return nil, nil

}

type SearchProductsRequest struct {
	Query string
}

type SearchProductsResponse struct {
	Products []entity.Product `json:"products"`
}

func (s *Service) SearchInProducts(req SearchProductsRequest) (SearchProductsResponse, error) {
	const op = "service.SearchProduct"

	vector, err := s.embedder.TextToVector(req.Query)
	if err != nil {
		return SearchProductsResponse{}, richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode)
	}
	productsIds, sErr := s.vectorStorage.SearchByTextVector(vector)
	if sErr != nil {
		fmt.Println(op, sErr)
		return SearchProductsResponse{}, richError.New(op).SetWrappedError(sErr).SetCode(richError.UnexpectedCode)
	}
	filterReq := FilterProductsFields{
		ID:        productsIds,
		Keyword:   req.Query,
		FromPrice: 0,
		EndPrice:  0,
	}
	filteredProducts, fErr := s.storage.Filter(filterReq)

	if fErr != nil {
		return SearchProductsResponse{}, richError.New(op).SetWrappedError(fErr).SetCode(richError.UnexpectedCode).SetMessage(errormsg.CantScanQueryResult)
	}
	semanticSearch, err := s.keywordSearch.SearchByText(req.Query)
	if err != nil {
		return SearchProductsResponse{}, richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode).SetMessage(errormsg.SemanticSearchError)
	}

	var response []entity.Product
	existenceCheckMap := make(map[uint]bool)

	for _, item := range semanticSearch {
		// Assert item as map[string]interface{}
		mappedItem, ok := item.(map[string]interface{})
		if !ok {
			fmt.Println("Failed to assert item as map[string]interface{}")
			continue
		}

		// Extract fields with type assertions
		id, ok := mappedItem["id"].(float64) // Adjust type as necessary (e.g., float64 for JSON numbers)
		if !ok {
			fmt.Println("Failed to assert 'id' as float64")
			continue
		}

		name, _ := mappedItem["name"].(string) // Optional fields can use default values
		description, _ := mappedItem["description"].(string)
		offPercent, _ := mappedItem["off_percent"].(float64)
		currentPrice, _ := mappedItem["current_price"].(float64)
		// Extract and process images
		var images []entity.Image
		if imageList, ok := mappedItem["images"].([]interface{}); ok {
			for _, img := range imageList {
				if imgURL, ok := img.(string); ok {
					images = append(images, entity.Image{
						Address: imgURL,
					})
				}
			}
		}
		existenceCheckMap[uint(id)] = true
		// Append the extracted data to the response
		response = append(response, entity.Product{
			ID:           entity.ProductID(id), // Convert to the expected type
			Name:         name,
			Description:  description,
			OffPercent:   uint(offPercent),
			CurrentPrice: float32(currentPrice),
			Images:       images, // Populate Images if applicable
		})
	}
	for _, item := range filteredProducts {
		if _, ok := existenceCheckMap[uint(item.ID)]; !ok {
			response = append(response, item)
		}
	}

	return SearchProductsResponse{Products: response}, nil

}

type AddNewProductRequest struct {
	Products []entity.Product `json:"products"`
}

func (s *Service) AddProductsBulky(req AddNewProductRequest) error {
	const op = "service.AddNewProductBulky"

	if err := s.storage.AddBulky(req.Products); err != nil {
		e, ok := err.(richError.RichError)
		if ok {
			fmt.Println(e.RetrieveOperation(), e.RetrieveAncestorMsg())
		}
		return richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode)
	}

	var productsVector []ProductsVector
	for _, pr := range req.Products {
		var images_url []string
		for _, img := range pr.Images {
			images_url = append(images_url, img.Address)
		}
		vectors, err := s.embedder.ImagesToVector(images_url)
		if err != nil {
			return richError.New(op).SetWrappedError(err).SetCode(richError.UnexpectedCode).SetMessage(errormsg.ModelServiceError)
		}
		if vectors != nil {
			productsVector = append(productsVector, ProductsVector{
				ProductId: pr.ID,
				Vector:    vectors,
			})
		}

		err = s.vectorStorage.AddBulkyProduct(productsVector)
		if err != nil {
			fmt.Println(err.Error())
			return richError.New(op).
				SetCode(richError.UnexpectedCode).
				SetWrappedError(err).
				SetMessage(errormsg.DataBaseInsertionError)
		}
	}
	return nil
}
