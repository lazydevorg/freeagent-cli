package client

import "time"

type BankAccount struct {
	Url            string    `json:"url"`
	OpeningBalance string    `json:"opening_balance"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	IsPersonal     bool      `json:"is_personal"`
	Status         string    `json:"status"`
	Currency       string    `json:"currency"`
	CurrentBalance string    `json:"current_balance"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedAt      time.Time `json:"created_at"`
}

func GetBankAccounts() ([]BankAccount, error) {
	bankAccounts, err := GetCollection[BankAccount]("bank_accounts", "bank_accounts", nil)
	if err != nil {
		return nil, err
	}
	return bankAccounts, nil
}

func GetBankAccount(id string) (*BankAccount, error) {
	bankAccount, err := GetEntity[BankAccount]("bank_accounts/"+id, "bank_account")
	if err != nil {
		return nil, err
	}
	return bankAccount, nil
}
