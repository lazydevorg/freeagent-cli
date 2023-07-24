package client

import "time"

type BankTransationExplanation struct {
	Url             string     `json:"url"`
	BankTransaction string     `json:"bank_transaction"`
	BankAccount     string     `json:"bank_account"`
	Category        string     `json:"category"`
	DatedOn         string     `json:"dated_on"`
	Description     string     `json:"description"`
	GrossValue      string     `json:"gross_value"`
	Project         string     `json:"project"`
	RebillType      string     `json:"rebill_type"`
	RebillFactor    string     `json:"rebill_factor"`
	UpdatedAt       time.Time  `json:"updated_at"`
	SalesTaxStatus  string     `json:"sales_tax_status"`
	SalesTaxRate    string     `json:"sales_tax_rate"`
	SalesTaxValue   string     `json:"sales_tax_value"`
	IsDeletable     bool       `json:"is_deletable"`
	Attachment      Attachment `json:"attachment"`
}
