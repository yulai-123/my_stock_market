package do

import (
	"context"
	"my_stock_market/model/po"
)

type FundBasic struct {
	ID            int64   `json:"id"`
	TSCode        string  `json:"ts_code"`
	Name          string  `json:"name"`
	Management    string  `json:"management"`
	Custodian     string  `json:"custodian"`
	FundType      string  `json:"fund_type"`
	FoundDate     string  `json:"found_date"`
	DueDate       string  `json:"due_date"`
	ListDate      string  `json:"list_date"`
	IssueDate     string  `json:"issue_date"`
	DelistDate    string  `json:"delist_date"`
	IssueAmount   float64 `json:"issue_amount"`
	MFee          float64 `json:"m_fee"`
	CFee          float64 `json:"c_fee"`
	DurationYear  float64 `json:"duration_year"`
	PValue        float64 `json:"p_value"`
	MinAmount     float64 `json:"min_amount"`
	ExpReturn     float64 `json:"exp_return"`
	Benchmark     string  `json:"benchmark"`
	Status        string  `json:"status"`
	InvestType    string  `json:"invest_type"`
	Type          string  `json:"type"`
	Trustee       string  `json:"trustee"`
	PurcStartdate string  `json:"purc_startdate"`
	RedmStartdate string  `json:"redm_startdate"`
	Market        string  `json:"market"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	DeletedAt int64 `json:"deleted_at"`
}

func FundBasicDO2PO(ctx context.Context, do *FundBasic) (*po.FundBasic, error) {
	if do == nil {
		return nil, nil
	}

	return &po.FundBasic{
		ID:            do.ID,
		TSCode:        do.TSCode,
		Name:          do.Name,
		Management:    do.Management,
		Custodian:     do.Custodian,
		FundType:      do.FundType,
		FoundDate:     do.FoundDate,
		DueDate:       do.DueDate,
		ListDate:      do.ListDate,
		IssueDate:     do.IssueDate,
		DelistDate:    do.DelistDate,
		IssueAmount:   do.IssueAmount,
		MFee:          do.MFee,
		CFee:          do.CFee,
		DurationYear:  do.DurationYear,
		PValue:        do.PValue,
		MinAmount:     do.MinAmount,
		ExpReturn:     do.ExpReturn,
		Benchmark:     do.Benchmark,
		Status:        do.Status,
		InvestType:    do.InvestType,
		Type:          do.Type,
		Trustee:       do.Trustee,
		PurcStartdate: do.PurcStartdate,
		RedmStartdate: do.RedmStartdate,
		Market:        do.Market,
		CreatedAt:     do.CreatedAt,
		UpdatedAt:     do.UpdatedAt,
		DeletedAt:     do.DeletedAt,
	}, nil
}

func FundBasicPO2DO(ctx context.Context, po *po.FundBasic) *FundBasic {
	if po == nil {
		return nil
	}

	return &FundBasic{
		ID:            po.ID,
		TSCode:        po.TSCode,
		Name:          po.Name,
		Management:    po.Management,
		Custodian:     po.Custodian,
		FundType:      po.FundType,
		FoundDate:     po.FoundDate,
		DueDate:       po.DueDate,
		ListDate:      po.ListDate,
		IssueDate:     po.IssueDate,
		DelistDate:    po.DelistDate,
		IssueAmount:   po.IssueAmount,
		MFee:          po.MFee,
		CFee:          po.CFee,
		DurationYear:  po.DurationYear,
		PValue:        po.PValue,
		MinAmount:     po.MinAmount,
		ExpReturn:     po.ExpReturn,
		Benchmark:     po.Benchmark,
		Status:        po.Status,
		InvestType:    po.InvestType,
		Type:          po.Type,
		Trustee:       po.Trustee,
		PurcStartdate: po.PurcStartdate,
		RedmStartdate: po.RedmStartdate,
		Market:        po.Market,
		CreatedAt:     po.CreatedAt,
		UpdatedAt:     po.UpdatedAt,
		DeletedAt:     po.DeletedAt,
	}
}
