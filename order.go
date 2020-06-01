package goshopee

import (
	"fmt"
	"time"
)

type OrderService interface {
	List(interface{}) ([]Order, error)
	ListWithPagination(sid uint64, offset, limit uint32, options interface{}) ([]Order, *Pagination, error)
	Count(interface{}) (int, error)
	Get(string, interface{}) (*Order, error)
	Create(Order) (*Order, error)
	Update(Order) (*Order, error)
	Cancel(sid uint64, ordersn, reason string, itemid uint64) error
	Delete(int64) error
}

// Order https://open.shopee.com/documents?module=4&type=1&id=397
type Order struct {
	OrderSN                      string            `json:"ordersn"`
	BuyerUsername                string            `json:"buyer_username"`
	RecipientAddress             *RecipientAddress `json:"recipient_address"`
	Status                       string            `json:"order_status"`
	Currency                     string            `json:"currency"`
	TrackingNo                   string            `json:"tracking_no"`
	EscrowAmount                 float64           `json:"escrow_amount"`
	TotalAmount                  float64           `json:"total_amount"`
	Country                      string            `json:"country"`
	ServiceCode                  string            `json:"service_code"`
	EstimatedShippingFee         string            `json:"estimated_shipping_fee"`
	PaymentMethod                string            `json:"payment_method"`
	ShippingCarrier              string            `json:"shipping_carrier"`
	COD                          bool              `json:"cod"`                  // This value indicates whether the order was a COD (cash on delivery) order.
	DaysToShip                   uint32            `json:"days_to_ship"`         // Shipping preparation time set by the seller when listing item on Shopee.
	ActualShippingCost           float64           `json:"actual_shipping_cost"` // The actual shipping cost of the order if available from external logistics partners.
	GoodsToDeclare               bool              `json:"goods_to_declare"`
	MessageToSeller              string            `json:"message_to_seller"`
	Note                         string            `json:"note"`
	NoteUpdateTime               time.Time         `json:"note_update_time"`
	CreateTime                   time.Time         `json:"create_time"`
	UpdateTime                   time.Time         `json:"update_time"`
	Items                        []OrderItem       `json:"items"`
	PayTime                      *uint32           `json:"pay_time"` // The time when the order status is updated from UNPAID to PAID. This value is NULL when order is not paid yet.
	DropShipper                  string            `json:"dropshipper"`
	CreditCardNumber             string            `json:"credit_card_number"`
	DropShipperPhone             string            `json:"dropshipper_phone"`
	ShipByDate                   uint32            `json:"ship_by_date"`
	IsSplitUp                    bool              `json:"is_split_up"`
	BuyerCancelReason            string            `json:"buyer_cancel_reason"`
	CancelBy                     string            `json:"cancel_by"`
	FmTN                         string            `json:"fm_tn"` // The first-mile tracking number.
	CancelReason                 string            `json:"cancel_reason"`
	EscrowTax                    float64           `json:"escrow_tax"`
	IsActualShippingFeeConfirmed bool              `json:"is_actual_shipping_fee_confirmed"`
}

type OrderItem struct {
	ItemID                   uint64  `json:"item_id"`
	ItemName                 string  `json:"item_name"`
	ItemSKU                  string  `json:"item_sku"`
	VariationID              uint64  `json:"variation_id"`
	VariationName            string  `json:"variation_name"`
	VariationSKU             string  `json:"variation_sku"`
	VariationQuantity        uint32  `json:"variation_quantity_purchased"`
	VariationDiscountedPrice string  `json:"variation_discounted_price"`
	VariationOriginalPrice   string  `json:"variation_original_price"`
	IsWholesale              bool    `json:"is_wholesale"`
	Weight                   float64 `json:"weight"`
	IsAddOnDeal              bool    `json:"is_add_on_deal"`
	IsMainItem               bool    `json:"is_main_item"`
	AddOnDealID              uint32  `json:"add_on_deal_id"`
	PromotionType            string  `json:"promotion_type"`
	PromotionID              uint32  `json:"promotion_id"`
}

type RecipientAddress struct {
	Town        string `json:"town"`
	City        string `json:"city"`
	Name        string `json:"name"`
	District    string `json:"district"`
	Country     string `json:"country"`
	Zipcode     string `json:"zipcode"`
	FullAddress string `json:"full_address"`
	Phone       string `json:"phone"`
	State       string `json:"state"`
}

// OrdersResponse Represents shopee.orders.GetOrdersList
// https://open.shopee.com/documents?module=4&type=1&id=399
type OrdersResponse struct {
	Orders    []Order `json:"orders"`
	More      bool    `json:"more"`
	RequestID string  `json:"request_id"`
}

// OrdersDetailResponse https://open.shopee.com/documents?module=4&type=1&id=397
type OrdersDetailResponse struct {
	Orders    []Order `json:"orders"`
	RequestID string  `json:"request_id"`
}

// Pagination of results
// type Pagination struct {
// 	Offset   uint32 `json:"offset"`
// 	PageSize uint32 `json:"page_size"`
// 	Total    uint32 `json:"total"`
// 	More     bool   `json:"more"`
// }

// OrderServiceOp handles communication with the product related methods of
// the Shopee API.
type OrderServiceOp struct {
	client *Client
}

// List xxx
func (s *OrderServiceOp) List(options interface{}) ([]Order, error) {
	// TODO:
	return nil, nil
}

// ListWithPagination https://open.shopee.com/documents?module=4&type=1&id=399
func (s *OrderServiceOp) ListWithPagination(sid uint64, offset, limit uint32, options interface{}) ([]Order, *Pagination, error) {
	path := "/orders/basics"
	wrappedData := map[string]interface{}{
		"pagination_offset":           offset,
		"pagination_entries_per_page": limit,
		"shopid":                      sid,
	}
	resource := new(OrdersResponse)
	err := s.client.Post(path, wrappedData, resource)
	page := &Pagination{
		Offset:   offset,
		PageSize: limit,
		More:     resource.More,
	}
	return resource.Orders, page, err
}

func (s *OrderServiceOp) Count(interface{}) (int, error) {
	return 0, nil
}
func (s *OrderServiceOp) Get(sid uint64, ordersn string, options interface{}) (*Order, error) {
	path := "/orders/detail"
	wrappedData := map[string]interface{}{
		"ordersn_list": []string{ordersn},
		"shopid":       sid,
	}
	resource := new(OrdersDetailResponse)
	err := s.client.Post(path, wrappedData, resource)
	if len(resource.Orders) == 0 {
		return nil, fmt.Errorf("no such order: [%s]", ordersn)
	}
	return &resource.Orders[0], err
}

// Create https://open.shopee.com/documents?module=2&type=1&id=365
func (s *OrderServiceOp) Create(Order) (*Order, error) {
	return nil, nil
}

// Update https://open.shopee.com/documents?module=2&type=1&id=376
func (s *OrderServiceOp) Update(Order) (*Order, error) {
	return nil, nil
}

type OrderCancelResponse struct {
	ModifiedTime uint32 `json:"modified_time"`
	RequestID    string `json:"request_id"`
}

// Cancel https://open.shopee.com/documents?module=4&type=1&id=395
func (s *OrderServiceOp) Cancel(sid uint64, ordersn, reason string, itemid uint64) error {
	path := "/orders/cancel"
	wrappedData := map[string]interface{}{
		"ordersn":       ordersn,
		"cancel_reason": reason,
		"item_id":       itemid,
		"shopid":        sid,
	}
	resource := new(OrderCancelResponse)
	err := s.client.Post(path, wrappedData, resource)
	return err
}

func (s *OrderServiceOp) Delete(int64) error {
	return nil
}
