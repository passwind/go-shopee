package goshopee

type Attribute struct {
	ID    uint64 `json:"attributes_id"`
	Value string `json:"value"`
}

type AttributeAdded struct {
	ID          uint64 `json:"attribute_id"`
	Name        string `json:"attribute_name"`
	IsMandatory bool   `json:"is_mandatory"`
	Type        string `json:"attribute_type"`
	Value       string `json:"attribute_value"`
}
