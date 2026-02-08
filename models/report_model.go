package models

type ResportToday struct {
	TotalRevenue      int               `json:"total_revenue"`
	TotalTransaction  int               `json:"total_transaksi"`
	ProductBestSeller ProductBestSeller `json:"product_best_seller"`
}

type ProductBestSeller struct {
	Name     string `json:"name"`
	QtySold  int    `json:"qty_sold"`
}
