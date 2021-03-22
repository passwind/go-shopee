package goshopee

// DiscountService: https://open.shopee.com/documents?module=1&type=1&id=361&version=1
type DiscountService interface {
	AddDiscount(uint64, map[string]interface{}) (*DiscountResponse,error)
	DeleteDiscount(uint64, uint64) (*DiscountActionResponse,error)
	AddDiscountItem(uint64,uint64, map[string]interface{}) (*DiscountResponse,error)
	DeleteDiscountItem(uint64, uint64, uint64, uint64) (*DiscountActionResponse,error)
	UpdateDiscount(uint64, uint64, map[string]interface{}) (*DiscountActionResponse,error)
	UpdateDiscountItems(uint64, uint64, map[string]interface{}) (*DiscountResponse,error)
}

type DiscountResponse struct {
	DiscountID uint64 `json:"discount_id"`
	Count uint32 `json:"count"`
	Warning string `json:"warning"`
	RequestID string `json:"request_id"`
	Errors []DiscountResponseError `json:"errors"`
}

type DiscountResponseError struct {
	ItemID uint64 `json:"item_id"`
	VariationID uint64 `json:"variation_id"`
	ErrorMsg string `json:"error_msg"`
}

type DiscountActionResponse struct {
	DiscountID uint64 `json:"discount_id"`
	RequestID string `json:"request_id"`
	ItemID uint64 `json:"item_id"`
	VariationID uint64 `json:"variation_id"`
	ModifyTime int64 `json:"modify_time"`
}

type Discount struct {
	ID uint64 `json:"discount_id"`
	Name string `json:"discount_name"`
	StartTime int64 `json:"start_time"`
	EndTime int64 `json:"end_time"`
	Status string `json:"status"`
}

type DiscountItem struct {
	ID uint64 `json:"item_id"`
	Name string `json:"item_name"`
	PurchaseLimit uint32 `json:"purchase_limit"`
	OriginalPrice float64 `json:"item_original_price"`
	PromotionPrice float64 `json:"item_promotion_price"`
	Stock uint32 `json:"stock"`
	InflatedOriginalPrice float64 `json:"item_inflated_original_price"`
	InflatedPromotionPrice float64 `json:"item_inflated_promotion_price"`
	Variations []DiscountVariation `json:"variations"`
}

type DiscountVariation struct {
	ID uint64 `json:"variation_id"`
	Name string `json:"variation_name"`
	OriginalPrice float64 `json:"variation_original_price"`
	PromotionPrice float64 `json:"variation_promotion_price"`
	Stock uint32 `json:"variation_stock"`
	InflatedOriginalPrice float64 `json:"variation_inflated_original_price"`
	InflatedPromotionPrice float64 `json:"variation_inflated_promotion_price"`
}

// DiscountServiceOp handles communication with the product related methods of
// the Shopee API.
type DiscountServiceOp struct {
	client *Client
}

func (s *DiscountServiceOp) AddDiscount(sid uint64, req map[string]interface{}) (*DiscountResponse,error) {
	path := "/discount/add"
	wrappedData := map[string]interface{}{
		"shopid":  sid,
	}
	for k,v:=range req {
		wrappedData[k]=v
	}
	resource := new(DiscountResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource, err
}


func (s *DiscountServiceOp) DeleteDiscount(sid, discountID uint64) (*DiscountActionResponse,error) {
	path := "/discount/delete"
	wrappedData := map[string]interface{}{
		"shopid":  sid,
		"discount_id": discountID,
	}
	
	resource := new(DiscountActionResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource, err
}

func (s *DiscountServiceOp) AddDiscountItem(sid, discountID uint64, req map[string]interface{}) (*DiscountResponse,error){
	path := "/discount/items/add"
	wrappedData := map[string]interface{}{
		"shopid":  sid,
		"discount_id": discountID,
	}
	for k,v:=range req {
		wrappedData[k]=v
	}
	resource := new(DiscountResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource, err
}

func (s *DiscountServiceOp) DeleteDiscountItem(sid, discountID, itemID, variationID uint64) (*DiscountActionResponse,error) {
	path := "/discount/item/delete"
	wrappedData := map[string]interface{}{
		"shopid":  sid,
		"discount_id": discountID,
		"item_id": itemID,
	}
	if variationID>0{
		wrappedData["variation_id"]=variationID
	}
	
	resource := new(DiscountActionResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource, err
}

func (s *DiscountServiceOp) UpdateDiscount(sid,discountID uint64, req map[string]interface{}) (*DiscountActionResponse,error) {
	path := "discount/update"
	wrappedData := map[string]interface{}{
		"shopid":  sid,
		"discount_id": discountID,
	}
	for k,v:=range req {
		wrappedData[k]=v
	}
	
	resource := new(DiscountActionResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource, err
}

func (s *DiscountServiceOp) UpdateDiscountItems(sid, discountID uint64, req map[string]interface{}) (*DiscountResponse,error) {
	path := "discount/items/update"
	wrappedData := map[string]interface{}{
		"shopid":  sid,
		"discount_id": discountID,
	}
	for k,v:=range req {
		wrappedData[k]=v
	}
	
	resource := new(DiscountResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource, err
}