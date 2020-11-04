package goshopee

type ItemService interface {
	List(interface{}) ([]Item, error)
	ListWithPagination(sid uint64, offset, limit uint32, options interface{}) ([]Item, *Pagination, error)
	Count(interface{}) (int, error)
	Get(uint64, uint64) (*Item, error)
	Create(newItem Item) (*ItemOper, error)
	Update(Item) (*ItemOper, error)
	UpdatePrice(sid, itemid uint64, price float64) (*ItemPriceOper, error)
	UpdateStock(sid, itemid uint64, stock uint32) (*ItemStockOper, error)
	Delete(sid, itemid uint64) error
	UnlistItem(sid, itemid uint64, unlist bool) ([]UnlistItemSuccess, []UnlistItemFailed, error)
	InitTierVariation(sid, itemid uint64, tierVariations []TierVariation, variations []TierVariationOperDef) ([]Variation, error)
	AddTierVariation(sid, itemid uint64, variations []TierVariationOperDef) ([]Variation, error)
	GetVariations(sid, itemid uint64) ([]TierVariation, []Variation, error)
	UpdateTierVariationList(sid, itemid uint64, tierVariations []TierVariation) error
	UpdateTierVariationIndex(sid, itemid uint64, variations []TierVariationIndexOperDef) error
}

// Item from https://open.shopee.com/documents?module=2&type=1&id=374
type ItemBase struct {
	ItemID                uint64      `json:"item_id"`
	ShopID                uint64      `json:"shopid"`
	ItemSKU               string      `json:"item_sku"`
	Status                string      `json:"status"`
	Name                  string      `json:"name"`
	Description           string      `json:"description"`
	Currency              string      `json:"currency"`
	HasVariation          bool        `json:"has_variation"`
	Price                 float64     `json:"price"`
	Stock                 uint32      `json:"stock"`
	CreateTime            uint32      `json:"create_time"`
	UpdateTime            uint32      `json:"update_time"`
	Weight                float64     `json:"weight"`
	CategoryID            uint64      `json:"category_id"`
	OriginalPrice         float64     `json:"original_price"`
	Variations            []Variation `json:"variations"`
	Attributes            []Attribute `json:"attributes"`
	Logistics             []Logistic  `json:"logistics"`
	Wholesales            []Wholesale `json:"wholesales,omitempty"`
	RatingStar            float64     `json:"rating_star"`
	CMTCount              uint32      `json:"cmt_count"`
	Sales                 uint32      `json:"sales"`
	Views                 uint32      `json:"views"`
	Likes                 uint32      `json:"likes"`
	PackageLength         float64     `json:"package_length"`
	PackageWidth          float64     `json:"package_width"`
	PackageHeight         float64     `json:"package_height"`
	DaysToShip            uint32      `json:"days_to_ship"`
	SizeChart             string      `json:"size_chart"`
	Condition             string      `json:"condition"`
	DiscountID            uint32      `json:"discount_id"`
	Is2tierItem           bool        `json:"is_2tier_item"`
	Tenures               []uint32    `json:"tenures"`
	ReservedStock         uint32      `json:"reserved_stock"`
	IsPreOrder            bool        `json:"is_pre_order"`
	InflatedPrice         float64     `json:"inflated_price"`
	InflatedOriginalPrice float64     `json:"inflated_original_price"`
	SipItemPrice          float64     `json:"sip_item_price"`
	PriceSource           string      `json:"price_source"`
}

type Item struct {
	ItemBase

	Images []string `json:"images"`
}

// ItemOper from https://open.shopee.com/documents?module=2&type=1&id=365
type ItemOper struct {
	ItemBase

	Images []Image `json:"images"`
}

type Wholesale struct {
	Min       uint32  `json:"min"`
	Max       uint32  `json:"max"`
	UnitPrice float64 `json:"unit_price"`
}

type ItemResponse struct {
	ItemID    uint32 `json:"item_id"`
	RequestID string `json:"request_id"`
}

// ItemsResponse Represents the result from the GetItemsList endpoint
// https://open.shopee.com/documents?module=2&type=1&id=375
type ItemsResponse struct {
	Items     []Item `json:"items"`
	More      bool   `json:"more"`
	Total     uint32 `json:"total"`
	RequestID string `json:"request_id"`
}

// Pagination of results
type Pagination struct {
	Offset   uint32 `json:"offset"`
	PageSize uint32 `json:"page_size"`
	Total    uint32 `json:"total"`
	More     bool   `json:"more"`
}

// ItemServiceOp handles communication with the product related methods of
// the Shopee API.
type ItemServiceOp struct {
	client *Client
}

func (s *ItemServiceOp) List(options interface{}) ([]Item, error) {
	// TODO:
	return nil, nil
}

// ListWithPagination https://open.shopee.com/documents?module=2&type=1&id=375
func (s *ItemServiceOp) ListWithPagination(sid uint64, offset, limit uint32, options interface{}) ([]Item, *Pagination, error) {
	path := "/items/get"
	wrappedData := map[string]interface{}{
		"pagination_offset":           offset,
		"pagination_entries_per_page": limit,
		"shopid":                      sid,
	}
	resource := new(ItemsResponse)
	err := s.client.Post(path, wrappedData, resource)
	page := &Pagination{
		Offset:   offset,
		PageSize: limit,
		Total:    resource.Total,
		More:     resource.More,
	}
	return resource.Items, page, err
}

func (s *ItemServiceOp) Count(interface{}) (int, error) {
	return 0, nil
}

type ItemDetailResponse struct {
	ItemID  uint64 `json:"item_id"`
	Item    *Item  `json:"item"`
	Warning string `json:"warning"`
}

func (s *ItemServiceOp) Get(sid, itemid uint64) (*Item, error) {
	path := "/item/get"
	wrappedData := map[string]interface{}{
		"item_id": itemid,
		"shopid":  sid,
	}
	resource := new(ItemDetailResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Item, err
}

type ItemOperResponse struct {
	ItemID    uint64   `json:"item_id"`
	Item      *Item    `json:"item"`
	SizeChart string   `json:"size_chart,omitempty"`
	Warning   string   `json:"warning"`
	FailImage []string `json:"fail_image,omitempty"`
}

// Create https://open.shopee.com/documents?module=2&type=1&id=365
func (s *ItemServiceOp) Create(newItem ItemOper) (*Item, error) {
	path := "/item/add"
	wrappedData, err := ToMapData(newItem)
	resource := new(ItemOperResponse)
	err = s.client.Post(path, wrappedData, resource)
	if resource == nil {
		return nil, err
	}
	return resource.Item, err
}

// Update https://open.shopee.com/documents?module=2&type=1&id=376
func (s *ItemServiceOp) Update(updItem ItemOper) (*Item, error) {
	path := "/item/update"
	wrappedData, err := ToMapData(updItem)
	resource := new(ItemOperResponse)
	err = s.client.Post(path, wrappedData, resource)
	if resource == nil {
		return nil, err
	}
	return resource.Item, err
}

type ItemDeleteResponse struct {
	ItemID    uint64 `json:"item_id"`
	Msg       string `json:"msg"`
	RequestID string `json:"request_id"`
}

// Delete https://open.shopee.com/documents?module=2&type=1&id=369
func (s *ItemServiceOp) Delete(sid, itemid uint64) error {
	path := "/item/delete"
	wrappedData := map[string]interface{}{
		"item_id": itemid,
		"shopid":  sid,
	}
	resource := new(ItemDeleteResponse)
	err := s.client.Post(path, wrappedData, resource)
	return err
}

type ItemPriceOper struct {
	ID            uint64  `json:"item_id"`
	ModifiedTime  uint32  `json:"modified_time"`
	Price         float64 `json:"price"`
	InflatedPrice float64 `json:"inflated_price"`
}

type ItemPriceOperResponse struct {
	Item      *ItemPriceOper `json:"item"`
	RequestID string         `json:"request_id"`
}

// UpdatePrice https://open.shopee.com/documents?module=2&type=1&id=377
func (s *ItemServiceOp) UpdatePrice(sid, itemid uint64, price float64) (*ItemPriceOper, error) {
	path := "/items/update_price"
	wrappedData := map[string]interface{}{
		"item_id": itemid,
		"price":   price,
		"shopid":  sid,
	}
	resource := new(ItemPriceOperResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Item, err
}

type ItemStockOper struct {
	ID           uint64 `json:"item_id"`
	ModifiedTime uint32 `json:"modified_time"`
	Stock        uint32 `json:"stock"`
}

type ItemStockOperResponse struct {
	Item      *ItemStockOper `json:"item"`
	RequestID string         `json:"request_id"`
}

// UpdateStock https://open.shopee.com/documents?module=2&type=1&id=378
func (s *ItemServiceOp) UpdateStock(sid, itemid uint64, stock uint32) (*ItemStockOper, error) {
	path := "/items/update_stock"
	wrappedData := map[string]interface{}{
		"item_id": itemid,
		"stock":   stock,
		"shopid":  sid,
	}
	resource := new(ItemStockOperResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Item, err
}

type UnlistItemFailed struct {
	ItemID           uint32 `json:"item_id"`
	ErrorDescription string `json:"error_description"`
}

type UnlistItemSuccess struct {
	ItemID uint32 `json:"item_id"`
	Unlist bool   `json:"unlist"`
}

type UnlistResponse struct {
	Failed    []UnlistItemFailed  `json:"failed"`
	Success   []UnlistItemSuccess `json:"success"`
	RequestID string              `json:"request_id"`
}

// UnlistItem https://open.shopee.com/documents?module=2&type=1&id=431
func (s *ItemServiceOp) UnlistItem(sid, itemid uint64, unlist bool) ([]UnlistItemSuccess, []UnlistItemFailed, error) {
	path := "/items/unlist"
	wrappedData := map[string]interface{}{
		"items": map[string]interface{}{
			"item_id": itemid,
			"unlist":  unlist,
		},
		"shopid": sid,
	}
	resource := new(UnlistResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Success, resource.Failed, err
}
