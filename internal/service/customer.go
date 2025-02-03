package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ibookerke/choco_parser_go/internal/config"
	"github.com/ibookerke/choco_parser_go/internal/domain"
	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
)

type CustomerService struct {
	customerRepo domain.CustomerRepository
	authRepo     domain.AuthRepository
	trm          trm.Manager
	cfg          config.Choco
}

func NewCustomerService(
	customerRepository domain.CustomerRepository,
	authRepository domain.AuthRepository,
	trm trm.Manager,
	cfg config.Choco,
) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepository,
		authRepo:     authRepository,
		trm:          trm,
		cfg:          cfg,
	}
}

func (c *CustomerService) FetchCustomerInfo(ctx context.Context, id domain.CustomerID, terminals string) (domain.Customer, error) {
	// Create customer
	exists, err := c.customerRepo.ExistsById(ctx, id)
	if err != nil {
		return domain.Customer{}, err
	}

	if exists {
		customer, err := c.customerRepo.FindById(ctx, id)
		if err != nil {
			return domain.Customer{}, err
		}
		return customer, nil
	}

	// fetch customerInfo
	customer, err := c.getFetchedCustomer(ctx, id, terminals)
	if err != nil {
		return domain.Customer{}, fmt.Errorf("failed to fetch customer info: %v", err)
	}

	_, err = c.customerRepo.Create(ctx, &customer)
	if err != nil {
		return domain.Customer{}, fmt.Errorf("failed to create customer: %v", err)
	}

	return customer, nil
}

type CustomerResponseData struct {
	Jsonapi struct {
		Version string `json:"version"`
	}
	Data struct {
		ID         int    `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			UserID         int         `json:"user_id"`
			UserAvatar     string      `json:"user_avatar"`
			Phone          string      `json:"phone"`
			Birthday       string      `json:"birthday"`
			DaysToBirthday interface{} `json:"days_to_birthday"`
			FullName       string      `json:"full_name"`
			/*
				"statistics": {
				                "turnover": 5790,
				                "orders_count": 1,
				                "average_bill": 5790,
				                "average_revenue": 1930,
				                "total_given_cashback": 0,
				                "total_payment_from_balance": 0
				            }
			*/
			Statistics struct {
				//Turnover                int64 `json:"turnover"`
				OrdersCount int64 `json:"orders_count"`
				//AverageBill             int64 `json:"average_bill"`
				//AverageRevenue          int64 `json:"average_revenue"`
				//TotalGivenCashback      int64 `json:"total_given_cashback"`
				//TotalPaymentFromBalance int64 `json:"total_payment_from_balance"`
			} `json:"statistics"`
		}
	}
}

func (c *CustomerService) getFetchedCustomer(ctx context.Context, id domain.CustomerID, terminals string) (domain.Customer, error) {
	// "https://api-proxy.choco.kz/analytics/v1/customer/12343106?terminals=9297,9341,9037,9038,9067,9066,8858"
	customerId := domain.CustomerIdToStr(id)

	url := "https://api-proxy.choco.kz/analytics/v1/customer/" + customerId + "?terminals=" + terminals
	// Create an HTTP client and request
	client := &http.Client{}

	//token, err := getAccessToken(ctx, c.authRepo, c.cfg)
	//if err != nil {
	//	return domain.Customer{}, fmt.Errorf("failed to get access token: %w", err)
	//}
	token := c.cfg.ChocoToken

	body, err := sendChocoRequest(ctx, client, url, token)
	if err != nil && err.Error() == "unexpected status code: 401" {
		token, err = fetchNewAccessToken(ctx, c.authRepo, client, c.cfg)
		if err != nil {
			return domain.Customer{}, fmt.Errorf("failed to fetch new access token: %w", err)
		}
		body, err = sendChocoRequest(ctx, client, url, token)
		if err != nil {
			return domain.Customer{}, fmt.Errorf("failed to send choco request: %w", err)
		}
	}
	if err != nil {
		return domain.Customer{}, fmt.Errorf("failed to send choco request: %w", err)
	}

	// Parse the JSON response
	var responseData CustomerResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		return domain.Customer{}, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	customer := domain.Customer{
		ID:         id,
		UserID:     responseData.Data.Attributes.UserID,
		FullName:   responseData.Data.Attributes.FullName,
		Phone:      responseData.Data.Attributes.Phone,
		Birthday:   responseData.Data.Attributes.Birthday,
		OrderCount: responseData.Data.Attributes.Statistics.OrdersCount,
	}

	return customer, nil
}
