package models

type DailySalesReport struct {
	TotalRevenue   int                 `json:"total_revenue"`
	TotalTransaksi int                 `json:"total_transaksi"`
	ProdukTerlaris *BestSellingProduct `json:"produk_terlaris"`
}

type BestSellingProduct struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

type SalesReport struct {
	StartDate      string              `json:"start_date,omitempty"`
	EndDate        string              `json:"end_date,omitempty"`
	TotalRevenue   int                 `json:"total_revenue"`
	TotalTransaksi int                 `json:"total_transaksi"`
	ProdukTerlaris *BestSellingProduct `json:"produk_terlaris"`
}
