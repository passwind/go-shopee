package goshopee

type Variation struct {
	ID                    uint64  `json:"variation_id"`
	Name                  string  `json:"name"`
	Stock                 uint32  `json:"stock"`
	Price                 float64 `json:"price"`
	VariationSKU          string  `json:"variation_sku"`
	Status                string  `json:"status"`
	CreateTime            uint32  `json:"create_time,omitempty"`
	UpdateTime            uint32  `json:"update_time,omitempty"`
	OriginalPrice         float64 `json:"original_price,omitempty"`
	InflatedOriginalPrice float64 `json:"inflated_original_price,omitempty"`
	InflatedPrice         float64 `json:"inflated_price,omitempty"`
	DiscountID            uint32  `json:"discount_id"`
}
