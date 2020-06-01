package goshopee

type Logistic struct {
	ID                   uint64  `json:"logistic_id"`
	Name                 string  `json:"logistic_name"`
	Enabled              bool    `json:"enabled"`
	ShippingFee          float64 `json:"shipping_fee"`
	SizeID               uint64  `json:"size_id"`
	IsFree               bool    `json:"is_free"`
	EstimatedShippingFee float64 `json:"estimated_shipping_fee"`
}
