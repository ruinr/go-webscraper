package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

var (
	errMissingProductName     = errors.New("missing product name")
	errMissingProductCategory = errors.New("missing product category")
	errMissingTTLDuration     = errors.New("missing TTL duration")
	errEmptyProduct           = errors.New("product is empty")
)

//StoreProduct save AmazonProduct into Redis with key product:{ASIN}
func StoreProduct(c *redis.Client, scrapedProduct *AmazonProduct) (err error) {
	var product = make(map[string]interface{})
	if scrapedProduct.Asin == "" {
		return ErrMissingASIN
	}
	if scrapedProduct.Name == "" {
		return errMissingProductName
	}
	if len(scrapedProduct.Categories) == 0 {
		return errMissingProductCategory
	}
	product["asin"] = scrapedProduct.Asin
	product["name"] = scrapedProduct.Name
	product["categories"] = strings.Join(scrapedProduct.Categories, ";")
	product["ranks"] = strings.Join(scrapedProduct.Ranks, ";")
	product["dimensions"] = strings.Join(scrapedProduct.Dimensions, ";")
	product["created_at"] = scrapedProduct.CreatedAt

	err = c.HMSet("product:"+scrapedProduct.Asin, product).Err()
	if err != nil {
		return err
	}
	return nil
}

//FetchProduct get AmazonProduct from Redis
//Todo: Fetch product will be used in Phase 2 for limited caching
func FetchProduct(c *redis.Client, asin string) (product AmazonProduct, err error) {
	if asin == "" {
		err = ErrMissingASIN
		return
	}
	//Check if product key exists
	exist, err := c.Exists("product:" + asin).Result()
	if err != nil {
		return
	} else if exist == 0 {
		err = fmt.Errorf("product key: product:%s doesn't exist", asin)
		return
	}

	name, err := c.HGet("product:"+asin, "name").Result()
	if err != nil {
		return
	}
	categories, err := c.HGet("product:"+asin, "categories").Result()
	if err != nil {
		return
	}
	createdAt, err := c.HGet("product:"+asin, "created_at").Result()
	if err != nil {
		return
	}

	product.Asin = asin
	product.Name = name
	product.Categories = strings.Split(categories, ";")
	product.CreatedAt = createdAt

	//ranks are not required, and it can be nil
	ranks, _ := c.HGet("product:"+asin, "ranks").Result()
	if len(ranks) > 0 {
		product.Ranks = strings.Split(ranks, ";")
	}
	//dimensions are not required, and it can be nil
	dimensions, _ := c.HGet("product:"+asin, "dimensions").Result()
	if len(dimensions) > 0 {
		product.Dimensions = strings.Split(dimensions, ";")
	}

	return
}

//GetProductFromCache tries to grab cached product
//from cacheProduct with key cacheProduct:{ASIN}
func GetProductFromCache(c *redis.Client, asin string) (product AmazonProduct, err error) {
	val, err := c.Get("cacheProduct:" + asin).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(val), &product)
	if err != nil {
		return
	}
	return
}

//AddProductToCache cache AmazonProduct with key cacheProduct:{ASIN}
func AddProductToCache(c *redis.Client, scrapedProduct *AmazonProduct, duration time.Duration) (err error) {
	if scrapedProduct.Name == "" {
		err = errEmptyProduct
		return
	}
	productJSON, err := json.Marshal(scrapedProduct)
	if err != nil {
		return
	}
	if duration == 0 {
		err = errMissingTTLDuration
		return
	}
	//store value as JSON string
	err = c.Set("cacheProduct:"+scrapedProduct.Asin, string(productJSON), duration).Err()
	return
}
