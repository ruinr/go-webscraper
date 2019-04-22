package v1

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	v1 "github.com/rnidev/go-webscraper/pkg/service/v1"
)

// newTestRedis returns a redis.Cmdable.
func newTestRedis() *redis.Client {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	return client
}

//Test if Redis server is available
func RedisIsAvailable(client redis.Cmdable) bool {
	return client.Ping().Err() == nil
}

func TestStoreProduct(t *testing.T) {
	c := newTestRedis()
	c.FlushDB()
	//Setup test data of AmazonProducts
	tests := []struct {
		subject   string
		product   v1.AmazonProduct
		expect    int64
		expectErr bool
		err       error
	}{
		{
			subject: "Test Success",
			product: v1.AmazonProduct{
				Asin: "B07FSH5L52",
				Name: "Longwu Women's Loose Casual Front Tie Short Sleeve Bandage Party Dress",
				Categories: []string{
					"Clothing, Shoes & Jewelry", "Novelty & More",
					"Clothing", "Novelty", "Women", "Dresses",
				},
				Ranks: []string{
					"#2,680 in Clothing, Shoes & Jewelry", "#9 in Women's Novelty Dresses",
					"#166 in Women's Dresses", "#1573 in Women's Shops",
				},
				CreatedAt: "2019-04-22T01:04:16.292932Z",
			},
			expectErr: false,
		},
		{
			subject: "Test missing ASIN",
			product: v1.AmazonProduct{
				Name: "Longwu Women's Loose Casual Front Tie Short Sleeve Bandage Party Dress",
				Categories: []string{
					"Clothing, Shoes & Jewelry", "Novelty & More",
					"Clothing", "Novelty", "Women", "Dresses",
				},
			},
			expectErr: true,
			err:       errors.New("missing ASIN in request"),
		},
		{
			subject: "Test missing name",
			product: v1.AmazonProduct{
				Asin: "B07FSH5L52",
			},
			expectErr: true,
			err:       errors.New("missing product name"),
		},
		{
			subject: "Test missing category",
			product: v1.AmazonProduct{
				Asin: "B07FSH5L52",
				Name: "Longwu Women's Loose Casual Front Tie Short Sleeve Bandage Party Dress",
			},
			expectErr: true,
			err:       errors.New("missing product category"),
		},
	}
	//Loop through and run tests
	for _, test := range tests {
		t.Run(test.subject, func(t *testing.T) {
			err := v1.StoreProduct(c, &test.product)
			if (err != nil && !test.expectErr) || !reflect.DeepEqual(err, test.err) {
				t.Errorf("StoreProduct() error = %v, expect Err %v", err, test.err)
				return
			}
			if err == nil && test.expectErr {
				t.Errorf("StoreProduct() error is nil, error expected: %t", test.expectErr)
				return
			}
		})
	}
}

func TestFetchProduct(t *testing.T) {
	c := newTestRedis()
	c.FlushDB()

	product := v1.AmazonProduct{
		Asin: "B07FSH5L52",
		Name: "Longwu Women's Loose Casual Front Tie Short Sleeve Bandage Party Dress",
		Categories: []string{
			"Clothing, Shoes & Jewelry", "Novelty & More",
			"Clothing", "Novelty", "Women", "Dresses",
		},
		Ranks: []string{
			"#2,680 in Clothing, Shoes & Jewelry", "#9 in Women's Novelty Dresses",
			"#166 in Women's Dresses", "#1573 in Women's Shops",
		},
		CreatedAt: "2019-04-22T01:04:16.292932Z",
	}
	v1.StoreProduct(c, &product)
	//Setup test data of AmazonProducts
	tests := []struct {
		subject   string
		asin      string
		expect    v1.AmazonProduct
		expectErr bool
		err       error
	}{
		{
			subject:   "Test Success",
			asin:      "B07FSH5L52",
			expect:    product,
			expectErr: false,
		},
		{
			subject:   "Test key doesn't exist",
			asin:      "notexist",
			expectErr: true,
			err:       fmt.Errorf("product key: product:%s doesn't exist", "notexist"),
		},
		{
			subject:   "Test missing asin",
			asin:      "",
			expectErr: true,
			err:       errors.New("missing ASIN in request"),
		},
	}

	for _, test := range tests {
		t.Run(test.subject, func(t *testing.T) {
			response, err := v1.FetchProduct(c, test.asin)
			if (err != nil && !test.expectErr) || !reflect.DeepEqual(err, test.err) {
				t.Errorf("v1.FetchProduct() error = %v, expect Err %v", err, test.err)
				return
			}
			if err == nil && (response.Asin != test.expect.Asin ||
				response.Name != test.expect.Name ||
				!reflect.DeepEqual(response.Categories, test.expect.Categories) ||
				!reflect.DeepEqual(response.Ranks, test.expect.Ranks)) {
				t.Errorf("v1.FetchProduct() = %v, expect %v", response, test.expect)
				return
			}
		})
	}
}

func TestAddProductToCache(t *testing.T) {
	c := newTestRedis()
	c.FlushDB()

	product := v1.AmazonProduct{
		Asin: "B07FSH5L52",
		Name: "Longwu Women's Loose Casual Front Tie Short Sleeve Bandage Party Dress",
		Categories: []string{
			"Clothing, Shoes & Jewelry", "Novelty & More",
			"Clothing", "Novelty", "Women", "Dresses",
		},
		Ranks: []string{
			"#2,680 in Clothing, Shoes & Jewelry", "#9 in Women's Novelty Dresses",
			"#166 in Women's Dresses", "#1573 in Women's Shops",
		},
		CreatedAt: "2019-04-22T01:04:16.292932Z",
	}

	tests := []struct {
		subject   string
		req       v1.AmazonProduct
		duration  time.Duration
		expectErr bool
		err       error
	}{
		{
			subject:   "Test Success",
			req:       product,
			duration:  time.Duration(int64(20)) * time.Second,
			expectErr: false,
			err:       nil,
		},
		{
			subject:   "Test missing TTL duration",
			req:       product,
			duration:  0,
			expectErr: true,
			err:       errors.New("missing TTL duration"),
		},
		{
			subject:   "Test mising product",
			req:       v1.AmazonProduct{},
			duration:  time.Duration(int64(20)) * time.Second,
			expectErr: true,
			err:       errors.New("product is empty"),
		},
	}
	for _, test := range tests {
		t.Run(test.subject, func(t *testing.T) {
			err := v1.AddProductToCache(c, &test.req, test.duration)
			if (err != nil && !test.expectErr) || !reflect.DeepEqual(err, test.err) {
				t.Errorf("v1.AddProductToCache() error = %v, expect Err %v", err, test.err)
				return
			}
		})
	}
}

func TestGetProductFromCache(t *testing.T) {
	c := newTestRedis()
	c.FlushDB()

	product := v1.AmazonProduct{
		Asin: "B07FSH5L52",
		Name: "Longwu Women's Loose Casual Front Tie Short Sleeve Bandage Party Dress",
		Categories: []string{
			"Clothing, Shoes & Jewelry", "Novelty & More",
			"Clothing", "Novelty", "Women", "Dresses",
		},
		Ranks: []string{
			"#2,680 in Clothing, Shoes & Jewelry", "#9 in Women's Novelty Dresses",
			"#166 in Women's Dresses", "#1573 in Women's Shops",
		},
		CreatedAt: "2019-04-22T01:04:16.292932Z",
	}

	v1.AddProductToCache(c, &product, time.Duration(int64(20))*time.Second)

	tests := []struct {
		subject   string
		req       string
		expect    v1.AmazonProduct
		expectErr bool
		err       error
	}{
		{
			subject:   "Test Success",
			req:       product.Asin,
			expect:    product,
			expectErr: false,
			err:       nil,
		}, {
			subject:   "Test key doesn't exist in cache",
			req:       "somekeydoesntexist",
			expectErr: true,
			err:       redis.Nil,
		},
	}
	for _, test := range tests {
		t.Run(test.subject, func(t *testing.T) {
			response, err := v1.GetProductFromCache(c, test.req)
			if (err != nil && !test.expectErr) || !reflect.DeepEqual(err, test.err) {
				t.Errorf("v1.GetProductFromCache() error = %v, expect Err %v", err, test.err)
				return
			}
			if err == nil && (response.Asin != test.expect.Asin ||
				response.Name != test.expect.Name ||
				!reflect.DeepEqual(response.Categories, test.expect.Categories) ||
				!reflect.DeepEqual(response.Ranks, test.expect.Ranks)) {
				t.Errorf("v1.FetchProduct() = %v, expect %v", response, test.expect)
				return
			}
		})
	}
}
