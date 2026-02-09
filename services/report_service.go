package services

import (
	"fmt"
	"kasir-go-api/models"
	"kasir-go-api/repositories"
	"time"
)

type ReportService struct {
	reportRepo *repositories.ReportRepository
}

func NewReportService(reportRepo *repositories.ReportRepository) *ReportService {
	return &ReportService{reportRepo: reportRepo}
}

func (s *ReportService) GetDailySalesReport() (*models.DailySalesReport, error) {
	return s.reportRepo.GetDailySalesReport()
}

func (s *ReportService) GetSalesReportByDateRange(startDate, endDate string) (*models.SalesReport, error) {
	// Validate date format
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format. Use YYYY-MM-DD")
	}

	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format. Use YYYY-MM-DD")
	}

	// Validate date range
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	if end.Before(start) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}

	return s.reportRepo.GetSalesReportByDateRange(startDate, endDate)
}
