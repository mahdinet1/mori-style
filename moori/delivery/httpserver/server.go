package httpServer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"moori/entity"
	"moori/service/product"
	"net/http"
	"reflect"
)

type Server struct {
	config     Config
	productSvc product.Service
}

type Config struct {
	Port string
}

func NewServer(config Config, productSvc product.Service) Server {
	return Server{
		config:     config,
		productSvc: productSvc,
	}
}
func (s *Server) Start() {
	e := echo.New()

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World this is a moori app echo server!")
	})

	// product group
	pr := e.Group("/product")
	pr.GET("/search", s.searchProduct)
	pr.POST("/group_add", s.addBulky)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", s.config.Port)))
}

func (s *Server) searchProduct(c echo.Context) error {
	searchQuery := c.QueryParam("query")
	fmt.Println(searchQuery)
	products, err := s.productSvc.SearchInProducts(product.SearchProductsRequest{Query: searchQuery})
	// TODO - map internal error to http errors
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, products)
}

func (s *Server) addBulky(c echo.Context) error {
	fmt.Println("add bulkyyyy")
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	// Read the file content
	all, err := io.ReadAll(src)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Parse the JSON into a slice of maps
	var products []map[string]interface{}
	if err := json.Unmarshal(all, &products); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Convert to a slice of Product structs
	var productsStruct []entity.Product
	for _, product := range products {
		// Parse images
		imagesUrl, ok := product["images"].([]interface{})
		var images []entity.Image
		if ok {
			for _, img := range imagesUrl {
				imgStr, valid := img.(string)
				if !valid {
					return echo.NewHTTPError(http.StatusBadRequest, "Invalid image URL format")
				}
				images = append(images, entity.Image{
					ID:        0,
					Address:   imgStr,
					ProductId: getUintValue(product, "id"),
				})
			}
		}

		// Parse other fields and construct the Product

		pr := entity.Product{
			ID:           getUintValue(product, "id"),
			Name:         getStringValue(product, "name"),
			Description:  getStringValue(product, "description"),
			Material:     getStringValue(product, "material"),
			ShopName:     getStringValue(product, "shop_name"),
			Link:         getStringValue(product, "link"),
			CategoryName: getStringValue(product, "category_name"),
			Region:       getStringValue(product, "getStringValue"),
			OffPercent:   uint(getFloatValue(product, "off_percent")),
			CurrentPrice: float32(getFloatValue(product, "current_price")),
			Images:       images,
		}
		productsStruct = append(productsStruct, pr)
	}

	err = s.productSvc.AddProductsBulky(product.AddNewProductRequest{
		Products: productsStruct[:100],
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "file uploaded")
}

func valueConverter(m map[string]interface{}, dest any) error {
	// Ensure dest is a pointer to a struct
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Struct {
		return errors.New("destination must be a pointer to a struct")
	}

	// Get the struct value and type
	destValue = destValue.Elem()

	for k, v := range m {
		// Get the struct field by name
		field := destValue.FieldByName(k)
		if !field.IsValid() {
			// If the field doesn't exist, skip it
			continue
		}

		if !field.CanSet() {
			// If the field cannot be set, skip it
			continue
		}

		// Convert the value to the correct type
		val := reflect.ValueOf(v)
		if val.Type().ConvertibleTo(field.Type()) {
			field.Set(val.Convert(field.Type()))
		} else {
			return fmt.Errorf("cannot convert value of field '%s' to type %s", k, field.Type().Name())
		}
	}

	return nil
}

// Helper function to safely extract string values from a map
func getStringValue(m map[string]interface{}, key string) string {

	if val, ok := m[key]; ok {
		if str, valid := val.(string); valid {
			return str
		}
	}
	return ""
}

func getUintValue(m map[string]interface{}, key string) entity.ProductID {
	if val, ok := m[key]; ok {

		if str, valid := val.(float64); valid {
			return entity.ProductID(str)
		}
	}
	return 0
}

// Helper function to safely extract float values from a map
func getFloatValue(m map[string]interface{}, key string) float64 {
	if val, ok := m[key]; ok {
		if num, valid := val.(float64); valid {
			return num
		}
	}
	return 0
}
