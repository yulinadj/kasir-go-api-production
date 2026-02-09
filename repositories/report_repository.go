package repositories

import (
	"database/sql"
	"kasir-go-api/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetDailySalesReport() (*models.DailySalesReport, error) {
	today := time.Now().Format("2006-01-02")

	var totalRevenue, totalTransaksi int

	// Get total revenue and transaction count for today
	err := repo.db.QueryRow(`
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(*) as total_transaksi
		FROM transactions
		WHERE DATE(created_at) = $1
	`, today).Scan(&totalRevenue, &totalTransaksi)

	if err != nil {
		return nil, err
	}

	// Get best selling product for today
	var bestProduct *models.BestSellingProduct
	var nama string
	var qty int

	err = repo.db.QueryRow(`
		SELECT 
			p.name,
			SUM(td.quantity) as qty_terjual
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN product p ON td.product_id = p.id
		WHERE DATE(t.created_at) = $1
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`, today).Scan(&nama, &qty)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == nil {
		bestProduct = &models.BestSellingProduct{
			Nama:       nama,
			QtyTerjual: qty,
		}
	}

	return &models.DailySalesReport{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: bestProduct,
	}, nil
}

func (repo *ReportRepository) GetSalesReportByDateRange(startDate, endDate string) (*models.SalesReport, error) {
	var totalRevenue, totalTransaksi int

	// Get total revenue and transaction count for date range
	err := repo.db.QueryRow(`
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COUNT(*) as total_transaksi
		FROM transactions
		WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2
	`, startDate, endDate).Scan(&totalRevenue, &totalTransaksi)

	if err != nil {
		return nil, err
	}

	// Get best selling product for date range
	var bestProduct *models.BestSellingProduct
	var nama string
	var qty int

	err = repo.db.QueryRow(`
		SELECT 
			p.name,
			SUM(td.quantity) as qty_terjual
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN product p ON td.product_id = p.id
		WHERE DATE(t.created_at) >= $1 AND DATE(t.created_at) <= $2
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`, startDate, endDate).Scan(&nama, &qty)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == nil {
		bestProduct = &models.BestSellingProduct{
			Nama:       nama,
			QtyTerjual: qty,
		}
	}

	return &models.SalesReport{
		StartDate:      startDate,
		EndDate:        endDate,
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: bestProduct,
	}, nil
}
