package goshopee

// Shop https://open.shopee.com/documents?module=6&type=1&id=410
/*
{
    "shop_id": 220006999,
    "shop_name": "roxieshopify.sg",
    "country": "SG",
    "shop_description": "",
    "videos": [],
    "images": [],
    "disable_make_offer": 0,
    "enable_display_unitno": false,
    "item_limit": 10000,
    "request_id": "541bb4a26f490b8244c9ea7b5a94cb5c",
    "status": "NORMAL",
    "installment_status": 0,
    "sip_a_shops": [],
    "is_cb": true,
    "non_pre_order_dts": 3,
    "auth_time": 1593507209,
    "expire_time": 1625043209
}
*/
type Shop struct {
	ID                  uint64             `json:"shop_id"`
	Name                string             `json:"shop_name"`
	Country             string             `json:"country"`
	Description         string             `json:"shop_description"`
	Videos              []string           `json:"videos"`
	Images              []string           `json:"images"`
	DisableMakeOffer    uint32             `json:"disable_make_offer"`
	EnableDisplayUnitNo bool               `json:"enable_display_unitno"`
	ItemLimit           uint32             `json:"item_limit"`
	RequestID           string             `json:"request_id"`
	Status              string             `json:"status"`
	InstallmentStatus   uint32             `json:"installment_status"`
	SIPAffiliateShops   []SIPAffiliateShop `json:"sip_a_shops"`
	IsCB                bool               `json:"is_cb"`
	NonPreOrderDTS      int                `json:"non_pre_order_dts"`
	AuthTime            int64              `json:"auth_time"`
	ExpireTime          int64              `json:"expire_time"`
}

type SIPAffiliateShop struct {
	AffiliateShopID string `json:"a_shop_id"`
	Country         string `json:"country"`
}

type ShopService interface {
	Get(sid uint64) (*Shop, error)
}

// ShopServiceOp handles communication with the product related methods of
// the Shopee API.
type ShopServiceOp struct {
	client *Client
}

func (s *ShopServiceOp) Get(sid uint64) (*Shop, error) {
	path := "/shop/get"
	wrappedData := map[string]interface{}{
		"shopid": sid,
	}
	resource := new(Shop)
	err := s.client.Post(path, wrappedData, resource)
	return resource, err
}
