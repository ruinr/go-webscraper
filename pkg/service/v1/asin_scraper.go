package v1

import (
	"html"
	"strings"

	"github.com/gocolly/colly"
)

//AmazonProduct is the default product struct for ASIN service
type AmazonProduct struct {
	Asin       string   `json:"asin"`
	Name       string   `json:"name"`
	Categories []string `json:"categories"`
	Ranks      []string `json:"ranks"`
	Dimensions []string `json:"dimensions"`
	CreatedAt  string   `json:"created_at"`
}

//GetProductInfoByASIN takes asin, build target url, and returns product info
func (product *AmazonProduct) GetProductInfoByASIN() (res *colly.Response, err error) {
	//ToDo: Take domain as a request for Phase 2
	domain := "www.amazon.com"
	var productURL string
	productURL = "https://" + domain + "/dp/" + product.Asin
	// Instantiate default collector
	c := colly.NewCollector(
		//Only allow whitelisted domains to be visited
		colly.AllowedDomains(domain),
	)

	// Error Handling
	c.OnError(func(r *colly.Response, rerr error) {
		res = r
		err = rerr
		return
	})

	// Start scraping product information
	/*
		Target: Product Categories, multiple
		Callback when find a element match the following attributes
		"#wayfinding-breadcrumbs_feature_div ul li span.a-list-item"
		"#wayfinding-breadcrumbs_feature_div" is unique to the div that
		contains information about the product categories
	*/
	c.OnHTML("#wayfinding-breadcrumbs_feature_div ul li span.a-list-item",
		func(e *colly.HTMLElement) {
			category := html.UnescapeString(e.ChildText(".a-link-normal"))
			if category != "" {
				product.Categories = append(product.Categories, category)
			}
		})

	c.OnHTML("#titleSection h1#title",
		func(e *colly.HTMLElement) {
			product.Name = html.UnescapeString(e.ChildText("span#productTitle"))
		})

	//There are three different layouts for product details: Ranks and Dimensions
	//Try Table View
	/*
		Target: Product Dimensions, can have one, multiple, or no dimension at all
		"#prodDetails" is unique to the div that
		contains information about the dimensions
	*/
	c.OnHTML("#prodDetails .wrapper .col1 .techD .content .attrG .pdTab table tbody tr",
		func(e *colly.HTMLElement) {
			if e.ChildText("td[class=label]") == "Product Dimensions" {
				dimensions := html.UnescapeString(e.ChildText("td[class=value]"))
				product.Dimensions = strings.Split(dimensions, ";")
			}
		})

	//Target: Product Main Rank
	c.OnHTML("#SalesRank td[class=value]",
		func(e *colly.HTMLElement) {
			result := html.UnescapeString(e.Text)
			resultSlice := strings.Split(strings.TrimSpace(result), "(")
			if len(resultSlice) > 0 {
				product.Ranks = append(product.Ranks, resultSlice[0])
			}
		})

	//Target: Product Subcategory Ranks
	c.OnHTML("#SalesRank td[class=value] ul.zg_hrsr li.zg_hrsr_item",
		func(e *colly.HTMLElement) {
			subRank := html.UnescapeString(e.ChildText("span.zg_hrsr_rank"))
			subCategory := html.UnescapeString(e.ChildText("span.zg_hrsr_ladder"))
			if subRank != "" && subCategory != "" {
				product.Ranks = append(product.Ranks, subRank+" "+subCategory)
			}
		})

	//Try bullet view

	//Target: Prodcut Main Rank
	c.OnHTML("#dpx-amazon-sales-rank_feature_div",
		func(e *colly.HTMLElement) {
			result := html.UnescapeString(e.ChildText("li#SalesRank")
			resultSlice := strings.Split(result, ":")
			var mainRank string
			if len(resultSlice) > 1 {
				mainRank = strings.TrimSpace(strings.Split(resultSlice[1], "(")[0])
			}
			if len(mainRank) > 1 {
				product.Ranks = append(product.Ranks, mainRank)
			}
		})

	// Try another bullet view for main rank
	c.OnHTML("#detail-bullets table tbody tr .bucket .content ul",
		func(e *colly.HTMLElement) {
			result := html.UnescapeString(e.ChildText("li#SalesRank"))
			resultSlice := strings.Split(result, ":")
			var mainRank string
			if len(results) > 1 {
				mainRank = strings.TrimSpace(strings.Split(resultSlice[1], "(")[0])
			}
			if len(mainRank) > 1 {
				product.Ranks = append(product.Ranks, mainRank)
			}
		})

	//Target: Ranks in subcategories
	c.OnHTML("li#SalesRank ul.zg_hrsr li.zg_hrsr_item",
		func(e *colly.HTMLElement) {
			subRank := html.UnescapeString(e.ChildText("span.zg_hrsr_rank"))
			subCategory := html.UnescapeString(e.ChildText("span.zg_hrsr_ladder"))
			if subRank != "" && subCategory != "" {
				product.Ranks = append(product.Ranks, subRank+" "+subCategory)
			}
		})

	//Target: Product Dimensions
	c.OnHTML("#detail-bullets table tbody tr td.bucket .content ul li",
		func(e *colly.HTMLElement) {
			//Target: Product Dimensions
			if e.ChildText("b") == "Product Dimensions:" {
				result := html.UnescapeString(e.Text)
				resultSlice := strings.Split(result, ":")
				if len(resultSlice) > 1 {
					product.Dimensions = strings.Split(strings.TrimSpace(resultSlice[1]), ";")
				}
			}
		})

	// A different bullet view for dimensions
	c.OnHTML("#detailBullets_feature_div ul li span",
		func(e *colly.HTMLElement) {
			if e.ChildText("span.a-text-bold") == "Product Dimensions:" {
				result := html.UnescapeString(e.Text)
				resultSlice := strings.Split(result, ":")
				if len(resultSlice) > 1 {
					product.Dimensions = strings.Split(strings.TrimSpace(resultSlice[1]), ";")
				}
			}
		})

	c.Visit(productURL)

	return
}
