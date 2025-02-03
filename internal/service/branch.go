package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ibookerke/choco_parser_go/internal/config"
	"github.com/ibookerke/choco_parser_go/internal/domain"
	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
)

type BranchService struct {
	branchRepo domain.BranchRepository
	authRepo   domain.AuthRepository
	trm        trm.Manager
	cfg        config.Choco
}

func NewBranchService(
	branchRepo domain.BranchRepository,
	authRepo domain.AuthRepository,
	trm trm.Manager,
	cfg config.Choco,
) *BranchService {
	return &BranchService{
		branchRepo: branchRepo,
		authRepo:   authRepo,
		trm:        trm,
		cfg:        cfg,
	}
}

type BranchResponseData struct {
	ErrorCode int          `json:"error_code"`
	Status    string       `json:"status"`
	Message   string       `json:"message"`
	Data      []RespBranch `json:"data"`
}

type RespBranch struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"type"`
	Token    string `json:"token"`
	Location struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		PartnerID   string `json:"partner_id"`
		PartnerName string `json:"partner_name"`
		PartnerLogo string `json:"partner_logo"`
	} `json:"location"`
}

func (bs *BranchService) FetchBranches(ctx context.Context) (string, error) {
	fmt.Println("fetching branches ")
	url := "https://api-proxy.choco.kz/acl/v3/staff/terminals?filter[terminal_types][]=main&filter[terminal_types][]=takeaway&filter[terminal_types][]=promotions&filter[terminal_types][]=waiterless&filter[terminal_types][]=special&filter[terminal_types][]=dr_delivery&filter[permission_name]=filial-customers"
	// Create an HTTP client and request
	client := &http.Client{}

	token := bs.cfg.ChocoToken
	//token, err := getAccessToken(ctx, bs.authRepo, bs.cfg)
	//if err != nil {
	//	return "", fmt.Errorf("failed to get access token: %w", err)
	//}

	body, err := sendChocoRequest(ctx, client, url, token)
	if err != nil && err.Error() == "unexpected status code: 401" {
		token, err = fetchNewAccessToken(ctx, bs.authRepo, client, bs.cfg)
		if err != nil {
			return "", fmt.Errorf("failed to fetch new access token: %w", err)
		}
		body, err = sendChocoRequest(ctx, client, url, token)
		if err != nil {
			return "", fmt.Errorf("failed to send choco request: %w", err)
		}
	}
	if err != nil {
		return "", fmt.Errorf("failed to send choco request: %w", err)
	}

	// Parse the JSON response
	var responseData BranchResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Map the temporary structs to your original struct
	var branches []string
	for _, tempBranch := range responseData.Data {
		branch := domain.Branch{
			ID:              domain.BranchId(tempBranch.ID),
			Name:            tempBranch.Name,
			Status:          tempBranch.Status,
			TypeID:          tempBranch.Type.ID,
			TypeName:        tempBranch.Type.Name,
			TypeDescription: tempBranch.Type.Description,
			Token:           tempBranch.Token,
			LocationID:      tempBranch.Location.ID,
			LocationName:    tempBranch.Location.Name,
			PartnerID:       tempBranch.Location.PartnerID,
			PartnerName:     tempBranch.Location.PartnerName,
			PartnerLogo:     tempBranch.Location.PartnerLogo,
		}

		branchExists, err := bs.branchRepo.CheckIfBranchExist(ctx, int64(branch.ID))
		if err != nil {
			return "", fmt.Errorf("failed to check if branch exists: %v", err)
		}

		if !branchExists {
			_, err := bs.branchRepo.Create(ctx, &branch)
			if err != nil {
				return "", fmt.Errorf("failed to create branch: %v", err)
			}
		}

		branches = append(branches, domain.BranchIdToStr(branch.ID))
	}

	return strings.Join(branches[:], ","), nil
}

func (bs *BranchService) GetBranchTerminals(ctx context.Context, companyName string) (string, error) {
	branches, err := bs.branchRepo.GetBranchesByCompanyName(ctx, companyName)
	if err != nil {
		return "", fmt.Errorf("failed to get branches: %v", err)
	}

	if companyName == "malatang" {
		extraBranches, err := bs.branchRepo.GetBranchesByCompanyName(ctx, "maratang")
		if err != nil {
			return "", fmt.Errorf("failed to get branches: %v", err)
		}

		branches = append(branches, extraBranches...)
	}

	// implode slice of branch ids to a string by comma
	branchIds := make([]string, len(branches))
	for i, branch := range branches {
		branchIds[i] = domain.BranchIdToStr(branch)
	}

	return strings.Join(branchIds[:], ","), nil
}
