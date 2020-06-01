package goshopee

type ItemService interface {
	List(interface{}) ([]Item, error)
	ListWithPagination(sid uint64, offset, limit uint32, options interface{}) ([]Item, *Pagination, error)
	Count(interface{}) (int, error)
	Get(uint64, uint64) (*ItemOper, error)
	Create(newItem Item) (*ItemOper, error)
	Update(Item) (*ItemOper, error)
	UpdatePrice(sid, itemid uint64, price float64) (*ItemPriceOper, error)
	UpdateStock(sid, itemid uint64, stock uint32) (*ItemStockOper, error)
	Delete(sid, itemid uint64) error
	UnlistItem(sid, itemid uint64, unlist bool) ([]UnlistItemSuccess, []UnlistItemFailed, error)
}

type Item struct {
	ItemID        uint64      `json:"item_id"`
	ShopID        uint64      `json:"shop_id"`
	UpdateTime    uint32      `json:"update_time"`
	Status        string      `json:"status"`
	CategoryID    uint64      `json:"category_id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	Price         float64     `json:"price"`
	Stock         uint32      `json:"stock"`
	ItemSKU       string      `json:"item_sku"`
	Weight        float64     `json:"weight"`
	PackageLength uint32      `json:"package_length"`
	PackageWidth  uint32      `json:"package_width"`
	PackageHight  uint32      `json:"package_height"`
	DaysToShip    uint32      `json:"days_to_ship"`
	Wholesales    []Wholesale `json:"wholesales,omitempty"`
	Variations    []Variation `json:"variations"`
	Images        []Image     `json:"images"`
	Attributes    []Attribute `json:"attributes"`
	Logistics     []Logistic  `json:"logistics"`
	Is2tierItem   bool        `json:"is_2tier_item"`
	Tenures       []uint32    `json:"tenures"`
	SizeChart     string      `json:"size_chart"`
	Condition     string      `json:"condition"`
	IsPreOrder    bool        `json:"is_pre_order"`
}

type Wholesale struct {
	Min       uint32  `json:"min"`
	Max       uint32  `json:"max"`
	UnitPrice float64 `json:"unit_price"`
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

func (s *ItemServiceOp) Get(sid, itemid uint64) (*ItemOper, error) {
	path := "/item/get"
	wrappedData := map[string]interface{}{
		"item_id": itemid,
		"shopid":  sid,
	}
	resource := new(ItemOperResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Item, err
}

type ItemOper struct {
	ShopID        uint64           `json:"shop_id"`
	ItemSKU       string           `json:"item_sku"`
	Status        string           `json:"status"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Images        []string         `json:"images"`
	Currency      string           `json:"currency"`
	HasVariation  bool             `json:"has_variation"`
	Price         float64          `json:"price"`
	Stock         uint32           `json:"stock"`
	CreateTime    uint32           `json:"create_time"`
	UpdateTime    uint32           `json:"update_time"`
	Weight        float64          `json:"weight"`
	CategoryID    uint64           `json:"category_id"`
	OriginalPrice float64          `json:"original_price"`
	Variations    []Variation      `json:"variations"`
	Attributes    []AttributeAdded `json:"attributes"`
	Logistics     []Logistic       `json:"logistics"`
	Wholesales    []Wholesale      `json:"wholesales"`
	Sales         uint32           `json:"sales"`
	Views         uint32           `json:"views"`
	Likes         uint32           `json:"likes"`
	PackageLength uint32           `json:"package_length"`
	PackageWidth  uint32           `json:"package_width"`
	PackageHight  uint32           `json:"package_height"`
	DaysToShip    uint32           `json:"days_to_ship"`
	RatingStar    float64          `json:"rating_star"`
	CmtCount      uint32           `json:"cmt_count"`
	Condition     string           `json:"condition"`
	DiscountID    uint32           `json:"discount_id"`
	IsPreOrder    bool             `json:"is_pre_order"`
}

type ItemOperResponse struct {
	ItemID    uint64    `json:"item_id"`
	Item      *ItemOper `json:"item"`
	SizeChart string    `json:"size_chart,omitempty"`
	Warning   string    `json:"warning"`
	FailImage []string  `json:"fail_image,omitempty"`
}

// Create https://open.shopee.com/documents?module=2&type=1&id=365
func (s *ItemServiceOp) Create(newItem Item) (*ItemOper, error) {
	path := "/item/add"
	wrappedData, err := ToMapData(newItem)
	resource := new(ItemOperResponse)
	err = s.client.Post(path, wrappedData, resource)
	return resource.Item, err
}

// Update https://open.shopee.com/documents?module=2&type=1&id=376
func (s *ItemServiceOp) Update(updItem Item) (*ItemOper, error) {
	path := "/item/update"
	wrappedData, err := ToMapData(updItem)
	resource := new(ItemOperResponse)
	err = s.client.Post(path, wrappedData, resource)
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
