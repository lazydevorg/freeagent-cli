package client

import "time"

type BankTransaction struct {
	Url                         string                      `json:"url"`
	Amount                      string                      `json:"amount"`
	BankAccount                 string                      `json:"bank_account"`
	DatedOn                     string                      `json:"dated_on"`
	Description                 string                      `json:"description"`
	FullDescription             string                      `json:"full_description"`
	UploadedAt                  time.Time                   `json:"uploaded_at"`
	UnexplainedAmount           string                      `json:"unexplained_amount"`
	IsManual                    bool                        `json:"is_manual"`
	TransactionId               string                      `json:"transaction_id"`
	CreatedAt                   time.Time                   `json:"created_at"`
	UpdatedAt                   time.Time                   `json:"updated_at"`
	MatchingTransactionsCount   int                         `json:"matching_transactions_count"`
	BankTransactionExplanations []BankTransationExplanation `json:"bank_transaction_explanations"`
}

func (c *Client) GetBankTransactions(bankAccountId string) ([]BankTransaction, error) {
	return GetCollection[BankTransaction](c, "bank_transactions?bank_account="+bankAccountId, "bank_transactions", nil)
}

func (c *Client) GetBankTransaction(id string) (*BankTransaction, error) {
	return GetEntity[BankTransaction](c, "bank_transactions/"+id, "bank_transaction")
}
