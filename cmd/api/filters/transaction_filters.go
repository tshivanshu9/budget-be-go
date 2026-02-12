package filters

import (
	"errors"
	"time"

	"github.com/jinzhu/now"
	"gorm.io/gorm"
)

type TransactionFilter struct {
	FromDate   string `query:"from_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate    string `query:"end_date" validate:"omitempty,datetime=2006-01-02"`
	CategoryId uint   `query:"category_id" validate:"omitempty,number,min=1"`
	WalletId   uint   `query:"wallet_id" validate:"omitempty,number,min=1"`
	Type       string `query:"type" validate:"omitempty,oneof=income expense"`
	Month      uint8  `query:"month" validate:"omitempty,number,min=1,max=12"`
	Year       uint16 `query:"year" validate:"omitempty,number,min=1"`
}

func (filter *TransactionFilter) ApplyFilters(query *gorm.DB) *gorm.DB {
	query = filter.filterByWalletId(query)
	query = filter.filterByCategoryId(query)
	query = filter.filterByType(query)
	query = filter.filterByMonth(query)
	query = filter.filterByYear(query)
	query = filter.filterByDateRange(query)
	return query
}

func (filter *TransactionFilter) filterByWalletId(query *gorm.DB) *gorm.DB {
	if filter.WalletId == 0 {
		return query
	}

	return query.Where("wallet_id = ?", filter.WalletId)
}

func (filter *TransactionFilter) filterByCategoryId(query *gorm.DB) *gorm.DB {
	if filter.CategoryId == 0 {
		return query
	}
	return query.Where("category_id = ?", filter.CategoryId)
}

func (filer *TransactionFilter) filterByType(query *gorm.DB) *gorm.DB {
	if filer.Type == "" {
		return query
	}
	return query.Where("type = ?", filer.Type)
}

func (filer *TransactionFilter) filterByMonth(query *gorm.DB) *gorm.DB {
	if filer.Month == 0 {
		return query
	}
	return query.Where("month = ?", filer.Month)
}

func (filer *TransactionFilter) filterByYear(query *gorm.DB) *gorm.DB {
	if filer.Year == 0 {
		return query
	}
	return query.Where("year = ?", filer.Year)
}

func (filter *TransactionFilter) filterByDateRange(query *gorm.DB) *gorm.DB {
	startDate := now.BeginningOfMonth()
	endDate := now.EndOfMonth()
	if filter.FromDate == "" && filter.EndDate == "" {
		filter.FromDate = startDate.Format("2006-01-02")
		filter.EndDate = endDate.Format("2006-01-02")
	} else {
		parsedFromDate, err := time.Parse(time.DateOnly, filter.FromDate)
		if err != nil {
			return query
		}
		parsedEndDate, err := time.Parse(time.DateOnly, filter.EndDate)
		if err != nil {
			return query
		}
		filter.FromDate = parsedFromDate.Format("2006-01-02")
		filter.EndDate = parsedEndDate.Format("2006-01-02")
	}
	return query.Where("date BETWEEN ? AND ?", filter.FromDate, filter.EndDate)
}

func (f *TransactionFilter) ValidateDates() error {
	if f.FromDate == "" && f.EndDate != "" || f.FromDate != "" && f.EndDate == "" {
		return errors.New("both from_date and end_date must be provided together")
	}
	if f.FromDate != "" && f.EndDate != "" {
		parsedFromDate, err := time.Parse(time.DateOnly, f.FromDate)
		if err != nil {
			return errors.New("invalid from_date format, expected YYYY-MM-DD")
		}
		parsedEndDate, err := time.Parse(time.DateOnly, f.EndDate)
		if err != nil {
			return errors.New("invalid end_date format, expected YYYY-MM-DD")
		}
		if parsedFromDate.After(parsedEndDate) {
			return errors.New("from_date cannot be after end_date")
		}
	}
	return nil
}
