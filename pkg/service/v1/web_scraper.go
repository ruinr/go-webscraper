package v1

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/ptypes"
	v1 "github.com/rnidev/go-webscraper/pkg/api/v1"
	"github.com/rnidev/go-webscraper/pkg/logger"
)

type webScraperServer struct {
	redisdb *redis.Client
}

var (
	//ErrMissingASIN is shared error message, returns if asin is missing
	ErrMissingASIN = errors.New("missing ASIN in request")
	//defaultTTL is the default time duration for cached product
	defaultTTL = time.Duration(int64(20)) * time.Minute
)

//NewScraperServer takes a new redis client for scraper server
func NewScraperServer(client *redis.Client) v1.WebScraperServer {
	return &webScraperServer{redisdb: client}
}

//GetProduct returns GetProductResponse and error
func (s *webScraperServer) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.GetProductResponse, error) {
	//validation
	if req.Asin == "" {
		return &v1.GetProductResponse{}, ErrMissingASIN
	}
	//Get a new redis client with context
	var product v1.Product
	cachedProduct, err := GetProductFromCache(s.redisdb, req.Asin)
	if err != nil && err != redis.Nil {
		return &v1.GetProductResponse{}, err
	} else if cachedProduct.Name != "" {
		product, err = mapProduct(&cachedProduct)
		if err != nil {
			return &v1.GetProductResponse{}, err
		}
		return &v1.GetProductResponse{
			Product: &product,
		}, nil
	}

	//Start product scraping
	var scrapedProduct AmazonProduct
	scrapedProduct.Asin = req.Asin
	res, err := scrapedProduct.GetProductInfoByASIN()

	if err != nil {
		logger.Log.Info("on error response: ", zap.String("error:", string(res.Body)))
		return &v1.GetProductResponse{}, err
	}

	if len(scrapedProduct.Categories) > 0 {
		scrapedProduct.CreatedAt = time.Now().In(time.UTC).Format(time.RFC3339Nano)
		err = StoreProduct(s.redisdb, &scrapedProduct)
		if err != nil {
			return &v1.GetProductResponse{}, err
		}
		err = AddProductToCache(s.redisdb, &scrapedProduct, defaultTTL)
		if err != nil {
			return &v1.GetProductResponse{}, err
		}
	}

	product, err = mapProduct(&scrapedProduct)

	return &v1.GetProductResponse{
		Product: &product,
	}, err
}

func mapProduct(scrapedProduct *AmazonProduct) (product v1.Product, err error) {
	product.Asin = scrapedProduct.Asin
	product.Name = scrapedProduct.Name

	for index, category := range scrapedProduct.Categories {
		product.Categories = append(product.Categories, &v1.ProductCategory{
			Name:  category,
			Level: int64(index + 1),
		})
	}

	for index, rank := range scrapedProduct.Ranks {
		product.Ranks = append(product.Ranks, &v1.ProductRank{
			RankInfo: rank,
			Level:    int64(index + 1),
		})
	}

	product.Dimensions = scrapedProduct.Dimensions

	t, err := time.Parse(time.RFC3339Nano, scrapedProduct.CreatedAt)
	if err != nil {
		return
	}
	product.CreatedAt, err = ptypes.TimestampProto(t)
	if err != nil {
		return
	}

	return
}
