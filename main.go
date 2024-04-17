package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/stinkyfingers/shoppinglistapi/server"
	"github.com/stinkyfingers/shoppinglistapi/source"
)

var (
	store = flag.String("store", "festival", "store to search")
	query = flag.String("query", "orange", "search query")
)

const port = ":8084"

func main() {
	fmt.Print("Running. \n")
	s, err := server.NewServer("jds")
	if err != nil {
		log.Fatalln(err)
	}
	rh, err := server.NewMux(s)
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(port, rh)
	if err != nil {
		log.Print(err)
	}
}

/*
deprecated, but keep for reference
*/

func cli() {
	source, err := source.GetSource(*store)
	if err != nil {
		log.Fatal("source ", err)
	}
	url, err := source.BuildURL(*query)
	if err != nil {
		log.Fatal("url ", err)
	}
	productMap, err := source.Search(url)
	if err != nil {
		log.Fatal("search ", err)
	}
	items, err := source.GetItems(productMap)
	if err != nil {
		log.Fatal("get items ", err)
	}
	fmt.Println(items)
	for _, item := range items {
		fmt.Println(item)
	}
}

const (
	appKey             = "festival_foods"
	fields             = "id,identifier,attribution_token,reference_id,reference_ids,upc,name,store_id,department_id,size,cover_image,price,sale_price,sale_price_md,sale_start_date,sale_finish_date,price_disclaimer,sale_price_disclaimer,is_favorite,relevance,popularity,shopper_walkpath,fulfillment_walkpath,quantity_step,quantity_minimum,quantity_initial,quantity_label,quantity_label_singular,varieties,quantity_size_ratio_description,status,status_id,sale_configuration_type_id,fulfillment_type_id,fulfillment_type_ids,other_attributes,clippable_offer,slot_message,call_out,has_featured_offer,tax_class_label,promotion_text,sale_offer,store_card_required,average_rating,review_count,like_code,shelf_tag_ids,offers,is_place_holder_cover_image,video_config,enforce_product_inventory,disallow_adding_to_cart,substitution_type_ids,unit_price,offer_sale_price,canonical_url,offered_together,sequence"
	limit              = 24
	relevanceSort      = "asc"
	sort               = "relevance"
	storeID            = 1955
	token              = "90d48a1ad035dcc14f910d451ad62b43"
	festivalSearchURL  = `https://api.freshop.com/1/products?app_key=festival_foods&fields=id,identifier,attribution_token,reference_id,reference_ids,upc,name,store_id,department_id,size,cover_image,price,sale_price,sale_price_md,sale_start_date,sale_finish_date,price_disclaimer,sale_price_disclaimer,is_favorite,relevance,popularity,shopper_walkpath,fulfillment_walkpath,quantity_step,quantity_minimum,quantity_initial,quantity_label,quantity_label_singular,varieties,quantity_size_ratio_description,status,status_id,sale_configuration_type_id,fulfillment_type_id,fulfillment_type_ids,other_attributes,clippable_offer,slot_message,call_out,has_featured_offer,tax_class_label,promotion_text,sale_offer,store_card_required,average_rating,review_count,like_code,shelf_tag_ids,offers,is_place_holder_cover_image,video_config,enforce_product_inventory,disallow_adding_to_cart,substitution_type_ids,unit_price,offer_sale_price,canonical_url,offered_together,sequence&include_offered_together=true&limit=24&q=orange&relevance_sort=asc&render_id=1712166093798&sort=relevance&store_id=1955&token=9a876428a89d78f56d6a350eb0667bd3`
	festivalProductURL = `https://api.freshop.com/1/products/78457?app_key=festival_foods&render_id=1712173153242&store_id=2238&token=9a876428a89d78f56d6a350eb0667bd3`

	// Schnucks has a allow origin header
	schnucksSearchTermURL = `https://api.schnucks.com/item-catalog-api/v1/item-search-terms?query=oranges`
	schnucksSearchURL     = `https://api.schnucks.com/item-catalog-api/v1/items?fulfillmentType=SELF&query=oranges&store=748&page=0&size=40`
)

type FestivalProduct struct {
	ID                      string   `json:"id"`
	Identifier              string   `json:"identifier"`
	AttributionToken        string   `json:"attribution_token"`
	ReferenceID             string   `json:"reference_id"`
	ReferenceIDs            []string `json:"reference_ids"`
	UPC                     string   `json:"upc"`
	Name                    string   `json:"name"`
	StoreID                 string   `json:"store_id"`
	DepartmentID            []string `json:"department_id"`
	Size                    string   `json:"size"`
	CoverImage              string   `json:"cover_image"`
	Price                   float64  `json:"price"`
	SalePrice               float64  `json:"sale_price"`
	SalePriceMD             float64  `json:"sale_price_md"`
	SaleStartDate           string   `json:"sale_start_date"`
	SaleFinishDate          string   `json:"sale_finish_date"`
	PriceDisclaimer         string   `json:"price_disclaimer"`
	SalePriceDisclaimer     string   `json:"sale_price_disclaimer"`
	IsFavorite              bool     `json:"is_favorite"`
	Relevance               float64  `json:"relevance"`
	Popularity              float64  `json:"popularity"`
	ShopperWalkpath         string   `json:"shopper_walkpath"`
	FulfillmentWalkpath     string   `json:"fulfillment_walkpath"`
	QuantityStep            float64  `json:"quantity_step"`
	QuantityMinimum         float64  `json:"quantity_minimum"`
	QuantityInitial         float64  `json:"quantity_initial"`
	QuantityLabel           string   `json:"quantity_label"`
	QuantityLabelSingular   string   `json:"quantity_label_singular"`
	Varieties               string   `json:"varieties"`
	QuantitySizeRatioDesc   string   `json:"quantity_size_ratio_description"`
	Status                  string   `json:"status"`
	StatusID                string   `json:"status_id"`
	SaleConfigurationTypeID int64    `json:"sale_configuration_type_id"`
	FulfillmentTypeID       string   `json:"fulfillment_type_id"`
	FulfillmentTypeIDs      []string `json:"fulfillment_type_ids"`
	OtherAttributes         string   `json:"other_attributes"`
	ClippableOffer          bool     `json:"clippable_offer"`
	SlotMessage             string   `json:"slot_message"`
	CallOut                 string   `json:"call_out"`
	HasFeaturedOffer        bool     `json:"has_featured_offer"`
	TaxClassLabel           string   `json:"tax_class_label"`
	PromotionText           string   `json:"promotion_text"`
	SaleOffer               bool     `json:"sale_offer"`
	StoreCardRequired       bool     `json:"store_card_required"`
	AverageRating           float64  `json:"average_rating"`
	ReviewCount             int64    `json:"review_count"`
	LikeCode                string   `json:"like_code"`
	ShelfTagIDs             []string `json:"shelf_tag_ids"`
	Offers                  string   `json:"offers"`
	IsPlaceHolderCoverImage bool     `json:"is_place_holder_cover_image"`
	VideoConfig             string   `json:"video_config"`
	EnforceProductInventory bool     `json:"enforce_product_inventory"`
	DisallowAddingToCart    bool     `json:"disallow_adding_to_cart"`
	SubstitutionTypeIDs     []string `json:"substitution_type_ids"`
	UnitPrice               float64  `json:"unit_price"`
	OfferSalePrice          float64  `json:"offer_sale_price"`
	CanonicalURL            string   `json:"canonical_url"`
	OfferedTogether         string   `json:"offered_together"`
	Sequence                int64    `json:"sequence"`
}

type FestivalResponse struct {
	Total int64             `json:"total"`
	Items []FestivalProduct `json:"items"`
}

func festivalSearch(query string) (*FestivalResponse, error) {
	url := fmt.Sprintf("https://api.freshop.com/1/products?app_key=%s&fields=%s&limit=%d&q=%s&relevance_sort=%s&sort=%s&store_id=%d&token=%s", appKey, fields, limit, query, relevanceSort, sort, storeID, token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response FestivalResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

type SchnucksProduct struct {
	UPCID   int64  `json:"upcid"`
	UPC     int64  `json:"upc"`
	FullUPC string `json:"fullUpc"`
	Name    string `json:"name"`
	// TODO
}

type SchnucksResponse struct {
	Codes    []string `json:"codes"`
	Messages []string `json:"messages"`
	PageInfo struct {
		Page     int64 `json:"page"`
		PageSize int64 `json:"pageSize"`
		Count    int64 `json:"count"`
	}
	Data []SchnucksProduct `json:"data"`
}

func listSchnucksStores(query string) (*SchnucksResponse, error) {
	url := fmt.Sprintf("https://api.schnucks.com/store-info-api/v1/stores/748")
	cli := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Origin", "https://www.schnucks.com")
	req.Header.Set("Referer", "https://www.schnucks.com/")
	req.Header.Set("Authorization", "QkkI38rzQT3tOIUSHa4n1CO71Ru1mrcr")
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(b))
	var response SchnucksResponse
	//if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
	//	return nil, err
	//}
	return &response, nil
}

func schnucksSearch(query string) (*SchnucksResponse, error) {
	url := fmt.Sprintf("https://api.schnucks.com/item-catalog-api/v1/items?fulfillmentType=SELF&query=%s&store=748&page=0&size=40", query)
	cli := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Origin", "https://www.schnucks.com")
	req.Header.Set("Referer", "https://www.schnucks.com/")
	req.Header.Set("Authorization", "QkkI38rzQT3tOIUSHa4n1CO71Ru1mrcr")
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response SchnucksResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
