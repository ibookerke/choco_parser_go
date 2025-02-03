package domain

import (
	"context"
	"strconv"
)

type BranchId int64

func BranchIdToStr(id BranchId) string {
	return strconv.FormatInt(int64(id), 10)
}

type Branch struct {
	ID              BranchId `json:"id"`
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	TypeID          int      `json:"type_id"`
	TypeName        string   `json:"type_name"`
	TypeDescription string   `json:"type_description"`
	Token           string   `json:"token"`
	LocationID      string   `json:"location_id"`
	LocationName    string   `json:"location_name"`
	PartnerID       string   `json:"partner_id"`
	PartnerName     string   `json:"partner_name"`
	PartnerLogo     string   `json:"partner_logo"`
}

type BranchRepository interface {
	Create(ctx context.Context, branch *Branch) (*Branch, error)
	CheckIfBranchExist(ctx context.Context, id int64) (bool, error)
	GetBranchesByCompanyName(ctx context.Context, companyName string) ([]BranchId, error)
}
