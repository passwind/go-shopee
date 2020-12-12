package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"g.wizardcloud.cn/service/shopee-ca-api/ca/models"
	db "g.wizardcloud.cn/service/shopee-ca-api/models"
	"g.wizardcloud.cn/service/shopee-ca-api/pkg/utils"
	shopeedb "g.wizardcloud.cn/service/shopee-ca-api/shopee/models"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	goshopee "github.com/passwind/go-shopee"
)

const (
	AsyncTaskGetOrders = "GetOrders"
)

func OrderStatusShopee2CA(shopeeStatus string) string {
	switch shopeeStatus {
	case "UNPAID":
		return "Pending"
	case "READY_TO_SHIP":
		return "ReleasedForShipment"
	case "CANCELLED":
		return "Canceled"
	case "COMPLETED":
		return "Shipped"
	}

	return "Pending"
}

func OrderStatusCA2Shopee(caStatus string) string {
	switch caStatus {
	case "ReleasedForShipment":
		return "READY_TO_SHIP"
	case "Shipped":
		return "COMPLETED"
	case "Pending":
		return "UNPAID"
	case "Canceled":
		return "CANCELLED"
	}

	return ""
}

// CancelReasonCA2Shopee ca to shopee
// ca https://access.channeladvisor.com/ApiDocumentation/Reference/Models/CancellationReason
// 1	Other
// 100	GeneralAdjustment
// 101	ItemNotAvailable
// 102	CustomerReturnedItem
// 103	CouldNotShip
// 104	AlternateItemProvided
// 105	BuyerCanceled
// 106	CustomerExchange
// 107	MerchandiseNotReceived
// 108	ShippingAddressUndeliverable
// shopee https://open.shopee.com/documents?module=4&type=1&id=395
// The reason seller want to cancel this order.
// Applicable values: OUT_OF_STOCK, CUSTOMER_REQUEST, UNDELIVERABLE_AREA, COD_NOT_SUPPORTED.
func CancelReasonCA2Shopee(caStatus string) string {
	switch caStatus {
	case "ShippingAddressUndeliverable":
		fallthrough
	case "CouldNotShip":
		return "UNDELIVERABLE_AREA"
	case "ItemNotAvailable":
		return "OUT_OF_STOCK"
		// case "GeneralAdjustment":
		// case "CustomerReturnedItem":
		// case "AlternateItemProvided":
		// case "BuyerCanceled":
		// case "CustomerExchange":
		// case "MerchandiseNotReceived":
		// case "Other":
	}

	return "CUSTOMER_REQUEST"
}

// AdjustOrderStatus 调整shopee的订单状态支持CA操作
func AdjustOrderStatus(id, origStatus string) string {
	status := origStatus
	cadborders, err := db.GetOrder(id)
	if err != nil {
		logs.Error("error to get local order: %s [%s]", err, id)
	} else {
		if cadborders.Status == "AcknowledgedBySeller" && origStatus == "ReleasedForShipment" {
			status = "AcknowledgedBySeller"
		} else if cadborders.Status == "Canceled" {
			status = "Canceled"
		}
	}

	return status
}

func OrderAddressShopee2CA(orig *goshopee.RecipientAddress) models.Address {
	if orig == nil {
		return models.Address{}
	}
	return models.Address{
		EmailAddress:    "",
		FirstName:       orig.Name,
		LastName:        orig.Name,
		AddressLine1:    orig.FullAddress,
		City:            orig.City,
		Country:         orig.Country,
		PostalCode:      orig.Zipcode,
		StateOrProvince: orig.State,
		AddressLine2:    "",
		CompanyName:     "",
		DaytimePhone:    orig.Phone,
	}
}

// ProduceCAOrderFromShopeeOrder TODO
func ProduceCAOrderFromShopeeOrder(sid uint64, upd *goshopee.Order) (db.Order, error) {
	status := OrderStatusShopee2CA(upd.Status)
	res := db.Order{
		ID:     upd.OrderSN,
		ShopID: sid,
		Status: status,
	}
	order, _ := ProduceOrderFromShopeeOrder(upd)
	byts, err := json.Marshal(order)
	if err != nil {
		return res, err
	}
	res.Content = string(byts)
	return res, nil
}

func ProduceOrderItemFromShopeeOrderItem(orig goshopee.OrderItem) (models.Item, error) {
	sku := strings.TrimSpace(orig.VariationSKU)
	if sku == "" {
		sku = fmt.Sprintf("%s-%s", orig.ItemName, orig.VariationName)
	}
	return models.Item{
		ID:        sku,
		SellerSKU: sku,
		UnitPrice: utils.Atof(orig.VariationDiscountedPrice),
		Quantity:  orig.VariationQuantity,
	}, nil
}
func ProduceOrderFromShopeeOrder(upd *goshopee.Order) (models.Order, error) {
	status := OrderStatusShopee2CA(upd.Status)
	address := OrderAddressShopee2CA(upd.RecipientAddress)
	items := []models.Item{}
	for _, oi := range upd.Items {
		ii, _ := ProduceOrderItemFromShopeeOrderItem(oi)
		items = append(items, ii)
	}
	order := models.Order{
		ID:                      upd.OrderSN,
		OrderDateUtc:            utils.FormatISO8601UTC(upd.CreateTime),
		OrderStatus:             status,
		BuyerAddress:            &address,
		ShippingAddress:         &address,
		RequestedShippingMethod: upd.ShippingCarrier,
		DeliverByDateUtc:        utils.FormatISO8601UTC(upd.ShipByDate),
		ShippingLabelURL:        "", //TODO
		TotalPrice:              utils.Atof(upd.TotalAmount),
		TotalTaxPrice:           utils.Atof(upd.EscrowTax),
		TotalShippingPrice:      utils.Atof(upd.EstimatedShippingFee),
		TotalShippingTaxPrice:   0.00, //TODO
		TotalGiftOptionPrice:    0.00,
		TotalGiftOptionTaxPrice: 0.00,
		TotalOrderDiscount:      0.00,
		TotalShippingDiscount:   0.00,
		OtherFees:               0.00,
		Currency:                upd.Currency,
		VatInclusive:            false, // TODO
		Items:                   items,
		SpecialInstructions:     upd.Note,
		PrivateNotes:            upd.Note,
		PaymentMethod:           upd.PaymentMethod,
		PaymentTransactionID:    "",
	}
	return order, nil
}

func SaveOrder(sid uint64, order goshopee.Order, origCA *db.Order) {
	byts, err := json.Marshal(&order)
	if err != nil {
		logs.Error("error to encode order: %s", err)
	}
	dborder := shopeedb.Orders{
		OrderSN:     order.OrderSN,
		ShopID:      sid,
		OrderStatus: order.Status,
		UpdateTime:  order.UpdateTime,
		Data:        string(byts),
	}
	_, err = shopeedb.AddOrUpdateOrder(order.OrderSN, dborder)
	if err != nil {
		logs.Error("error to save order: %s", err)
	}

	dbcaorder, err := ProduceCAOrderFromShopeeOrder(sid, &order)
	if err != nil {
		logs.Error("produce ca order error: %s", err)
	}

	if origCA != nil && origCA.Status == "AcknowledgedBySeller" &&
		dbcaorder.Status != "Shipped" && dbcaorder.Status != "Canceled" {
		dbcaorder.Status = "AcknowledgedBySeller"
	}

	_, err = db.AddOrUpdateOrder(dbcaorder.ID, dbcaorder)
	if err != nil {
		logs.Error("error to save ca order: %s", err)
	}
}

func FetchOrdersFromShopee(sid uint64, status string, limit int) ([]models.Order, error) {
	// TODO:
	if limit > 20 {
		limit = 20
	}
	logs.Info("FetchOrdersFromShopee: %d - %d", sid, limit)
	defer func() {
		logs.Info("FetchOrdersFromShopee: %d end", sid)
	}()

	options := make(map[string]interface{})

	shopeestatus := OrderStatusCA2Shopee(status)
	if shopeestatus != "" {
		options["order_status"] = shopeestatus
	}

	orders, _, err := shopeeClient.Order.ListWithPagination(sid, 0, uint32(limit), options)
	if err != nil {
		return nil, fmt.Errorf("fetch orders from shopee error: %s", err)
	}
	// logs.Debug("page is :%#v", page)

	results := []models.Order{}

	if len(orders) == 0 {
		return results, nil
	}

	sns := []string{}
	for _, o := range orders {
		sns = append(sns, o.OrderSN)
	}

	detailOrders, errlist, err := shopeeClient.Order.GetMulti(sid, sns)
	if err != nil {
		return nil, fmt.Errorf("fetch orders from shopee error: %s", err)
	}
	if len(errlist) > 0 {
		logs.Warn("error to fetch orders: %d %v - #v", sid, sns, errlist)
	}

	for _, order := range detailOrders {
		origCA, _ := db.GetOrder(order.OrderSN)
		status := OrderStatusShopee2CA(order.Status)
		if origCA != nil && origCA.Status == "AcknowledgedBySeller" &&
			status != "Shipped" && status != "Canceled" {
			status = "AcknowledgedBySeller"
		}

		go SaveOrder(sid, order, origCA)

		caorder, _ := ProduceOrderFromShopeeOrder(&order)
		caorder.OrderStatus = status

		results = append(results, caorder)
	}
	return results, nil
}

// GetOrders get orders
func GetOrders(sid uint64, status string, limit int) (*models.ActionResponse, error) {
	// return GetRemoteOrders(sid, status, limit)
	results, err := FetchOrdersFromShopee(sid, status, limit)
	if err != nil {
		resp := models.ErrorReponse(models.ErrSystem(fmt.Errorf("error to create task: %s", err)))
		return resp, nil
	}
	orders := []models.Order{}
	for _, o := range results {
		o.OrderStatus = AdjustOrderStatus(o.ID, o.OrderStatus)
		orders = append(orders, o)
	}
	return models.ObjectsResponse(orders), nil
}

func HandleGetOrdersTask(t *AsyncTask) {
	params := strings.Split(t.TaskRequest, ":")
	sid, _ := strconv.ParseUint(params[1], 10, 64)
	limit, _ := strconv.Atoi(params[3])
	results, err := FetchOrdersFromShopee(sid, "", limit)
	if err != nil {
		resp := models.ErrorReponse(models.ErrSystem(fmt.Errorf("error to create task: %s", err)))
		UpdateTask(t.ID, AsyncTask{
			TaskResponse: resp,
			TaskStatus:   TaskStatusComplete,
		})
		return
	}
	// logs.Debug("api result: %#v", results)
	// TODO: adjust local status
	orders := []models.Order{}
	for _, o := range results {
		o.OrderStatus = AdjustOrderStatus(o.ID, o.OrderStatus)
		orders = append(orders, o)
	}
	if t, err := UpdateTask(t.ID, AsyncTask{
		TaskResponse: models.ObjectsResponse(orders),
		TaskStatus:   TaskStatusComplete,
	}); err != nil {
		logs.Error("error to end task: %s [%#v]", err, t)
	}

}

func GetRemoteOrders(sid uint64, status string, limit int) (*models.ActionResponse, error) {
	task, err := NewTask(AsyncTaskGetOrders, AsyncTask{
		TaskStatus:  TaskStatusPending,
		TaskRequest: fmt.Sprintf("AsyncGetOrders:%d:%s:%d", sid, status, limit),
	})
	if err != nil {
		resp := models.ErrorReponse(models.ErrSystem(fmt.Errorf("error to create task: %s", err)))
		return resp, err
	}

	go HandleAsyncTasks(AsyncTaskGetOrders)

	resp := models.PendingResponse(task.URI())
	return resp, nil
}

func GetLocalOrders(sid uint64, status string, limit int) (*models.ActionResponse, error) {
	dborders, err := db.GetOrders(sid, status, limit)
	if err != nil {
		return nil, err
	}
	results := []models.Order{}
	for _, v := range dborders {
		var order models.Order
		if err := json.Unmarshal([]byte(v.Content), &order); err != nil {
			return nil, err
		}
		order.OrderStatus = v.Status
		results = append(results, order)
	}

	return models.ObjectsResponse(results), nil
}

func FetchOrder(sid uint64, sn string) (*models.Order, error) {
	sorder, err := shopeeClient.Order.Get(sid, sn)
	if err != nil {
		return nil, err
	}
	order, err := ProduceOrderFromShopeeOrder(sorder)
	return &order, err
}

func GetOrder(sid uint64, id string) (*models.ActionResponse, *models.Error) {
	order, err := FetchOrder(sid, id)
	if err != nil {
		logs.Error("error to fetch order: %s [%d - %s]", err, sid, id)
		return nil, models.ErrOrderNotFound() // models.ErrSystem(err)
	}
	order.OrderStatus = AdjustOrderStatus(id, order.OrderStatus)
	return models.ObjectsResponse(order), nil
}

func GetLocalOrder(id string) (*models.ActionResponse, *models.Error) {
	dbOrder, err := db.GetOrder(id)
	if err == gorm.ErrRecordNotFound {
		return nil, models.ErrOrderNotFound()
	} else if err != nil {
		return nil, models.ErrSystem(err)
	}
	var order models.Order
	if err := json.Unmarshal([]byte(dbOrder.Content), &order); err != nil {
		return nil, models.ErrSystem(err)
	}
	order.OrderStatus = dbOrder.Status
	return models.ObjectsResponse(order), nil
}

func AcknowledgeOrder(sid uint64, id string) *models.Error {
	// check order existed?
	_, err1 := GetOrder(sid, id)
	if err1 != nil {
		return err1
	}

	if err := db.AcknowledgeOrder(id); err != nil {
		return models.ErrSystem(err)
	}
	return nil
}

func CancelOrder(sid uint64, id string, req models.OrderCancellation) *models.Error {
	// check order existed?
	_, err1 := GetOrder(sid, id)
	if err1 != nil {
		return err1
	}

	reason := CancelReasonCA2Shopee(req.Items[0].Reason)
	err := shopeeClient.Order.Cancel(sid, id, reason, nil)
	if err != nil {
		return models.ErrSystem(err)
	}
	return nil
}

func RefundOrder(sid uint64, id string) *models.Error {
	// check order existed?
	_, err1 := GetOrder(sid, id)
	if err1 != nil {
		return err1
	}

	if err := db.RefundOrder(id); err != nil {
		return models.ErrSystem(err)
	}
	return nil
}

func AddOrUpdateFulfullments(sid uint64, req []models.OrderFulfillment) ([]models.OrderFulfillmentResult, error) {
	return InitShopeeLogistics(sid, req)
}

func InitShopeeLogistics(sid uint64, req []models.OrderFulfillment) ([]models.OrderFulfillmentResult, error) {
	results := []models.OrderFulfillmentResult{}
	errs := []*models.Error{}
	for _, v := range req {
		for j, item := range v.Items {
			if j == 0 {
				// params := map[string]interface{}{
				// 	"non_integrated": map[string]interface{}{
				// 		"tracking_no": item.TrackingNumber,
				// 	},
				// }
				params := map[string]interface{}{
					"dropoff": map[string]interface{}{
						"branch_id":        123,
						"sender_real_name": "Mike Zhu",
						"tracking_no":      item.TrackingNumber,
					},
				}
				if _, err := shopeeClient.Logistic.Init(sid, v.ID, params); err != nil {
					err1 := models.ErrShipmentFailed()
					errs = append(errs, err1)
				}
				break
			}
		}

		res := models.OrderFulfillmentResult{
			ID:     v.ID,
			Result: "Success",
		}
		if len(errs) > 0 {
			res.Result = "Fail"
			res.Errors = errs
		}
		results = append(results, res)
	}
	return results, nil
}

func AddOrUpdateLocalFulfullments(sid uint64, req []models.OrderFulfillment) ([]models.OrderFulfillmentResult, error) {
	results := []models.OrderFulfillmentResult{}
	errs := []*models.Error{}
	for _, v := range req {
		for _, item := range v.Items {
			byts, err := json.Marshal(item)
			if err != nil {
				err1 := models.ErrShipmentFailed()
				errs = append(errs, err1)
				continue
			}
			newFF := db.OrderFulfillment{
				ShopID:         sid,
				OrderID:        v.ID,
				OrderItemID:    item.OrderItemID,
				SellerSKU:      item.SellerSKU,
				TrackingNumber: item.TrackingNumber,
				Content:        string(byts),
			}
			_, err = db.AddOrUpdateOrderFulfillment(sid, v.ID, item.OrderItemID, item.SellerSKU, newFF)
			if err != nil {
				err1 := models.ErrShipmentFailed()
				errs = append(errs, err1)
				continue
			}
		}
		err := db.UpdateOrder(v.ID, map[string]interface{}{
			"status": "Shipped",
		})
		if err != nil {
			err1 := models.ErrShipmentFailed()
			errs = append(errs, err1)
			continue
		}
		res := models.OrderFulfillmentResult{
			ID:     v.ID,
			Result: "Success",
		}
		if len(errs) > 0 {
			res.Result = "Fail"
			res.Errors = errs
		}
		results = append(results, res)
	}
	return results, nil
}
