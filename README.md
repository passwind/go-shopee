# go-shopee

shopee open api with golang

https://open.shopee.com/documents?module=63&type=2&id=51

## How to use

Initialize Client And request order list

```
  app := goshopee.App{
		PartnerID:  xxxxxxxx,
		PartnerKey: yyyyyyyyyy,
		APIURL:     https://api.shopee.com,
	}

	client := goshopee.NewClient(app)

  // fetch order list
  client.Order.ListWithPagination(sid, 0, 100, nil)
```