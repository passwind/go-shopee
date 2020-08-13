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
	Init(uint64, string, map[string]interface{}) error
	GetParameterForInit(sid uint64, ordersn string) (*GetParameterForInitResponse, error)
	GetLogisticInfo(sid uint64, ordersn string) (*GetLogisticInfoResponse, error)
	List(uint64) ([]Logistic, error)
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

type ListReponse struct {
	Logistics []Logistic `json:"logistics"`
	RequestID string     `json:"request_id"`
}

// List https://open.shopee.com/documents?module=3&type=1&id=384
func (s *LogisticServiceOp) List(sid uint64) ([]Logistic, error) {
	path := "/logistics/channel/get"
	wrappedData := map[string]interface{}{
		"shopid": sid,
	}
	resource := new(ListReponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Logistics, err
}

// Init https://open.shopee.com/documents?module=3&type=1&id=389
func (s *LogisticServiceOp) Init(sid uint64, ordersn string, params map[string]interface{}) (string, error) {
	path := "/logistics/init"
	wrappedData := map[string]interface{}{
		"ordersn": ordersn,
		"shopid":  sid,
	}
	for k, v := range params {
		wrappedData[k] = v
	}
	resource := new(LogisticInitResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource.TrackingNumber, err
}

type GetParameterForInitResponse struct {
	Pickup        []string `json:"pickup"`
	Dropoff       []string `json:"dropoff"`
	NonIntegrated []string `json:"non_integrated"`
	RequestID     string   `json:"request_id"`
}

// GetParameterForInit https://open.shopee.com/documents?module=3&type=1&id=386
func (s *LogisticServiceOp) GetParameterForInit(sid uint64, ordersn string) (*GetParameterForInitResponse, error) {
	path := "/logistics/init_parameter/get"
	wrappedData := map[string]interface{}{
		"ordersn": ordersn,
		"shopid":  sid,
	}
	resource := new(GetParameterForInitResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource, err
}

type GetLogisticInfoResponsePickup struct {
	AddressList []Address `json:"address_list"`
}

type GetLogisticInfoResponseDropoff struct {
	BranchList []Branch `json:"branch_list"`
}

type GetLogisticInfoResponse struct {
	Pickup     GetLogisticInfoResponsePickup  `json:"pickup"`
	Dropoff    GetLogisticInfoResponseDropoff `json:"dropoff"`
	InfoNeeded map[string][]string            `json:"info_needed"`
	RequestID  string                         `json:"request_id"`
}

// GetLogisticInfo https://open.shopee.com/documents?module=3&type=1&id=417
func (s *LogisticServiceOp) GetLogisticInfo(sid uint64, ordersn string) (*GetLogisticInfoResponse, error) {
	path := "/logistics/init_info/get"
	wrappedData := map[string]interface{}{
		"ordersn": ordersn,
		"shopid":  sid,
	}
	resource := new(GetLogisticInfoResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource, err
}
