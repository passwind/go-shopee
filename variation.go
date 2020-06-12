package goshopee

type Variation struct {
	ID                    uint64  `json:"variation_id"`
	Name                  string  `json:"name"`
	Stock                 uint32  `json:"stock"`
	Price                 float64 `json:"price"`
	VariationSKU          string  `json:"variation_sku"`
	Status                string  `json:"status"`
	CreateTime            uint32  `json:"create_time"`
	UpdateTime            uint32  `json:"update_time"`
	OriginalPrice         float64 `json:"original_price"`
	InflatedOriginalPrice float64 `json:"inflated_original_price"`
	InflatedPrice         float64 `json:"inflated_price"`
	DiscountID            uint32  `json:"discount_id"`
	ModifiedTime          uint32  `json:"modified_time"`
	ItemID                uint64  `json:"item_id"`
}

type VariationService interface {
	Create(uint64, uint64, Variation) (*Variation, error)
	UpdateVariationPrice(uint64, uint64, Variation) (*Variation, error)
	UpdateVariationStock(uint64, uint64, Variation) (*Variation, error)
}

// VariationServiceOp handles communication with the product related methods of
// the Shopee API.
type VariationServiceOp struct {
	client *Client
}

type AddVariationsRequest struct {
	ItemID     uint64      `json:"item_id"`
	Variations []Variation `json:"variations"`
	ShopID     uint64      `json:"shop_id"`
}

type AddVariationsReponse struct {
	ItemID       uint64      `json:"item_id"`
	ModifiedTime uint32      `json:"modified_time"`
	Variations   []Variation `json:"variations"`
	RequestID    string      `json:"request_id"`
}

// Create https://open.shopee.com/documents?module=2&type=1&id=368
func (s *VariationServiceOp) Create(sid, itemID uint64, newItem Variation) (*Variation, error) {
	path := "/item/add_variations"
	req := AddVariationsRequest{
		ItemID: itemID,
		ShopID: sid,
		Variations: []Variation{
			newItem,
		},
	}
	wrappedData, err := ToMapData(req)
	resource := new(AddVariationsReponse)
	err = s.client.Post(path, wrappedData, resource)
	if len(resource.Variations) == 0 {
		return nil, err
	}
	return &resource.Variations[0], err
}

type UpdateVariationPriceRequest struct {
	ItemID      uint64  `json:"item_id"`
	VariationID uint64  `json:"variation_id"`
	Price       float64 `json:"price"`
	ShopID      uint64  `json:"shop_id"`
}

type UpdateVariationPriceResponse struct {
	Variation Variation `json:"item"`
	RequestID string    `json:"request_id"`
}

// UpdateVariationPrice https://open.shopee.com/documents?module=2&type=1&id=379
func (s *VariationServiceOp) UpdateVariationPrice(sid, itemID uint64, updItem Variation) (*Variation, error) {
	path := "/items/update_variation_price"
	req := UpdateVariationPriceRequest{
		ItemID:      itemID,
		VariationID: updItem.ID,
		Price:       updItem.Price,
		ShopID:      sid,
	}
	wrappedData, err := ToMapData(req)
	resource := new(UpdateVariationPriceResponse)
	err = s.client.Post(path, wrappedData, resource)
	return &resource.Variation, err
}

type UpdateVariationStockRequest struct {
	ItemID      uint64 `json:"item_id"`
	VariationID uint64 `json:"variation_id"`
	Stock       uint32 `json:"stock"`
	ShopID      uint64 `json:"shop_id"`
}

type UpdateVariationStockResponse struct {
	Variation Variation `json:"item"`
	RequestID string    `json:"request_id"`
}

// UpdateVariationStock https://open.shopee.com/documents?module=2&type=1&id=380
func (s *VariationServiceOp) UpdateVariationStock(sid, itemID uint64, updItem Variation) (*Variation, error) {
	path := "/items/update_variation_stock"
	req := UpdateVariationStockRequest{
		ItemID:      itemID,
		VariationID: updItem.ID,
		Stock:       updItem.Stock,
		ShopID:      sid,
	}
	wrappedData, err := ToMapData(req)
	resource := new(UpdateVariationStockResponse)
	err = s.client.Post(path, wrappedData, resource)
	return &resource.Variation, err
}
