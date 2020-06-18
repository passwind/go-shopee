package goshopee

type ItemCategoryService interface {
	List(sid uint64, options map[string]interface{}) ([]ItemCategory, error)
}

// ItemCategoryServiceOp handles communication with the product related methods of
// the Shopee API.
type ItemCategoryServiceOp struct {
	client *Client
}

type ItemCategoriesResponse struct {
	Categories []ItemCategory `json:"categories"`
	RequestID  string         `json:"request_id"`
}

// List xxx
func (s *ItemCategoryServiceOp) List(sid uint64, options map[string]interface{}) ([]ItemCategory, error) {
	path := "/item/categories/get"
	wrappedData := map[string]interface{}{
		"shopid": sid,
	}
	resource := new(ItemCategoriesResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Categories, err
}

type ItemCategory struct {
	ID               uint64            `json:"category_id"`
	ParentID         uint64            `json:"parent_id"`
	Name             string            `json:"category_name"`
	HasChildren      bool              `json:"has_children"`
	DaysToShipLimits *DaysToShipLimits `json:"days_to_ship_limits"`
}

type DaysToShipLimits struct {
	Min int `json:"min_limit"`
	Max int `json:"max_limit"`
}
