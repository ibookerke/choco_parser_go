package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ibookerke/choco_parser_go/internal/config"
	"github.com/ibookerke/choco_parser_go/internal/domain"
	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
)

type CompanyCustomersService struct {
	companyCustomerRepo domain.CompanyCustomerRepository
	trm                 trm.Manager
	cfg                 config.Choco
}

func NewCompanyCustomersService(
	companyCustomerRepo domain.CompanyCustomerRepository,
	trm trm.Manager,
	cfg config.Choco,
) *CompanyCustomersService {
	return &CompanyCustomersService{
		companyCustomerRepo: companyCustomerRepo,
		trm:                 trm,
		cfg:                 cfg,
	}
}

type CustomersResponse struct {
	Meta struct {
		Page struct {
			CurrentPage int `json:"currentPage"`
			LastPage    int `json:"lastPage"`
		} `json:"page"`
	} `json:"meta"`
	Data []struct {
		ID         int64 `json:"id"`
		Attributes struct {
			UserID        int64   `json:"user_id"`
			Phone         string  `json:"phone"`
			Turnover      float64 `json:"turnover"`
			FullName      string  `json:"full_name"`
			LastVisitDate string  `json:"last_visit_date"`
			VisitsCount   int64   `json:"visits_count"`
			AverageBill   float64 `json:"average_bill"`
		} `json:"attributes"`
	} `json:"data"`
}

func (s *CompanyCustomersService) FetchCompanyCustomers(ctx context.Context, terminals, companyName string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	page := 1

	fetchFrom := "2010-01-01+00:00:00"
	today := time.Now().Format("2006-01-02") + "+23:59:59"
	baseUrl := fmt.Sprintf(
		"https://api-proxy.choco.kz/analytics/v1/customers?terminals=%s&sort=turnover&filter[start_date]=%s&filter[end_date]=%s",
		terminals,
		fetchFrom,
		today,
	)

	for {
		// Construct the paginated URL
		url := fmt.Sprintf("%s&page=%d", baseUrl, page)
		token := s.cfg.ChocoToken
		body, err := sendChocoRequest(ctx, client, url, token)
		if err != nil {
			return fmt.Errorf("failed to send choco request: %w", err)
		}

		// Process the response body
		var response CustomersResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		for _, item := range response.Data {
			fmt.Println("Processing customer ID: ", item.ID)

			customer := domain.CompanyCustomer{
				ID:            item.ID,
				Company:       companyName,
				UserID:        item.Attributes.UserID,
				FullName:      item.Attributes.FullName,
				Phone:         item.Attributes.Phone,
				Turnover:      item.Attributes.Turnover,
				LastVisitDate: item.Attributes.LastVisitDate,
				VisitsCount:   item.Attributes.VisitsCount,
				AverageBill:   item.Attributes.AverageBill,
			}

			if err := s.companyCustomerRepo.Store(ctx, &customer); err != nil {
				return fmt.Errorf("failed to store customer: %w", err)
			}
		}

		// Check if there are more pages
		if page >= response.Meta.Page.LastPage {
			break
		}
		page++
	}

	return nil
}
