package goshopee

type ItemAttributeService interface {
	List(cid uint64, options map[string]interface{}) ([]ItemAttribute, error)
}

type ItemAttribute struct {
	ID        uint64               `json:"attribute_id"`
	Name      string               `json:"attribute_name"`
	Mandatory bool                 `json:"is_mandatory"`
	Type      string               `json:"attribute_type"`
	InputType string               `json:"input_type"`
	Options   []string             `json:"options"`
	Values    []ItemAttributeValue `json:"values"`
}

type ItemAttributeValue struct {
	OriginalValue  string `json:"original_value"`
	TranslateValue string `json:"translate_value"`
}

// ItemAttributeServiceOp handles communication with the product related methods of
// the Shopee API.
type ItemAttributeServiceOp struct {
	client *Client
}

type ItemAttributesResponse struct {
	Attributes []ItemAttribute `json:"attributes"`
	RequestID  string          `json:"request_id"`
}

// List xxx
func (s *ItemAttributeServiceOp) List(cid uint64, options map[string]interface{}) ([]ItemAttribute, error) {
	path := "/item/attributes/get"
	wrappedData := map[string]interface{}{
		"category_id": cid,
	}
	for k, v := range options {
		wrappedData[k] = v
	}
	resource := new(ItemAttributesResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Attributes, err
}
