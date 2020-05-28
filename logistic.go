package goshopee

type Logistic struct {
	ID         uint64 `json:"-"`
	LogisticID uint64 `json:"logistic_id"`
	Enabled    bool   `json:"enabled"`
	ProductID  uint64 `json:"-"`
}
