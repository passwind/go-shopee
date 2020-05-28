package goshopee

type ItemService interface {
	List(interface{}) ([]Item, error)
	ListWithPagination(sid uint64, offset, limit uint32, options interface{}) ([]Item, *Pagination, error)
	Count(interface{}) (int, error)
	Get(int64, interface{}) (*Item, error)
	Create(Item) (*Item, error)
	Update(Item) (*Item, error)
	Delete(int64) error
}

type Item struct {
	ItemID      uint64      `json:"item_id"`
	ShopID      uint64      `json:"shop_id"`
	UpdateTime  uint32      `json:"update_time"`
	Status      string      `json:"status"`
	CategoryID  uint64      `json:"category_id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Price       float64     `json:"price,omitempty"`
	Stock       uint32      `json:"stock,omitempty"`
	ItemSKU     string      `json:"item_sku,omitempty"`
	Weight      float64     `json:"weight,omitempty"`
	Variations  []Variation `json:"variations"`
	Images      []Image     `json:"images,omitempty"`
	Attributes  []Attribute `json:"attributes,omitempty"`
	Logistics   []Logistic  `json:"logistics,omitempty"`
	Is2tierItem bool        `json:"is_2tier_item"`
	Tenures     []uint32    `json:"tenures,omitempty"`
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
func (s *ItemServiceOp) Get(int64, interface{}) (*Item, error) {
	return nil, nil
}

// Create https://open.shopee.com/documents?module=2&type=1&id=365
func (s *ItemServiceOp) Create(Item) (*Item, error) {
	return nil, nil
}

// Update https://open.shopee.com/documents?module=2&type=1&id=376
func (s *ItemServiceOp) Update(Item) (*Item, error) {
	return nil, nil
}

// Delete https://open.shopee.com/documents?module=2&type=1&id=369
func (s *ItemServiceOp) Delete(int64) error {
	return nil
}
