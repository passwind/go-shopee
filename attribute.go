package goshopee

type Attribute struct {
	ID             uint64 `json:"-"`
	AttributeID    uint64 `json:"attributes_id"`
	AttributeValue string `json:"value"`
	ProductID      uint64 `json:"-"`
}
