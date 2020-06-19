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

type LogisticService interface {
	Init(uint64, string, string) error
}

// LogisticServiceOp handles communication with the logistics related methods of
// the Shopee API.
type LogisticServiceOp struct {
	client *Client
}

type LogisticInitResponse struct {
	TrackingNumber string `json:"tracking_number"`
	RequestID      string `json:"request_id"`
}

// Init https://open.shopee.com/documents?module=3&type=1&id=389
func (s *LogisticServiceOp) Init(sid uint64, ordersn, trackingNo string) error {
	path := "/logistics/init"
	wrappedData := map[string]interface{}{
		"ordersn": ordersn,
		"shopid":  sid,
		"non_integrated": map[string]interface{}{
			"tracking_no": trackingNo,
		},
	}
	resource := new(LogisticInitResponse)
	err := s.client.Post(path, wrappedData, resource)
	return err
}
