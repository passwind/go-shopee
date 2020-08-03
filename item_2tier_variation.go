package goshopee

// TierVariation 2-tier variation
// https://open.shopee.com/documents?module=2&type=1&id=422
type TierVariation struct {
	Name      string   `json:"name"`
	Options   []string `json:"options"`
	ImagesURL []string `json:"images_url"`
}

type TierVariationOperResponse struct {
	RequestID       string          `json:"request_id"`
	ItemID          uint32          `json:"item_id"`
	TierVariation   []TierVariation `json:"tier_variation,omitempty"`
	VariationIDList []Variation     `json:"variation_id_list"`
}

func (s *ItemServiceOp) InitTierVariation(sid, itemid uint64, tierVariations []TierVariation, variations []Variation) ([]Variation, error) {
	path := "/item/tier_var/init"
	wrappedData := map[string]interface{}{
		"item_id":        itemid,
		"shopid":         sid,
		"tier_variation": tierVariations,
		"variation":      variations,
	}
	resource := new(TierVariationOperResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource.VariationIDList, err
}

func (s *ItemServiceOp) AddTierVariation(sid, itemid uint64, variations []Variation) ([]Variation, error) {
	path := "/item/tier_var/add"
	wrappedData := map[string]interface{}{
		"item_id":   itemid,
		"shopid":    sid,
		"variation": variations,
	}
	resource := new(TierVariationOperResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource.VariationIDList, err
}

func (s *ItemServiceOp) GetVariations(sid, itemid uint64) ([]TierVariation, []Variation, error) {
	path := "/item/tier_var/get"
	wrappedData := map[string]interface{}{
		"item_id": itemid,
		"shopid":  sid,
	}
	resource := new(TierVariationOperResponse)
	err := s.client.Post(path, wrappedData, resource)
	return resource.TierVariation, resource.VariationIDList, err
}

func (s *ItemServiceOp) UpdateTierVariationList(sid, itemid uint64, tierVariations []TierVariation) error {
	path := "/item/tier_var/update_list"
	wrappedData := map[string]interface{}{
		"item_id":        itemid,
		"shopid":         sid,
		"tier_variation": tierVariations,
	}
	resource := new(ItemResponse)
	err := s.client.Post(path, wrappedData, resource)
	return err
}

func (s *ItemServiceOp) UpdateTierVariationIndex(sid, itemid uint64, variations []Variation) error {
	path := "/item/tier_var/update"
	wrappedData := map[string]interface{}{
		"item_id":   itemid,
		"shopid":    sid,
		"variation": variations,
	}
	resource := new(ItemResponse)
	err := s.client.Post(path, wrappedData, resource)
	return err
}
