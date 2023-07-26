package main

import (
	"github.com/lazydevorg/freeagent-cli/internal/cli"
	"github.com/lazydevorg/freeagent-cli/internal/client"
)

func main() {
	//cmd.Execute()
	//auth.Authenticate(false)
	//data, err := client.GetBankAccounts()
	//data, err := client.GetBankTransactionExplanations("https://api.sandbox.freeagent.com/v2/bank_accounts/21044")
	data, err := client.GetBankTransactions("1161419")
	if err != nil {
		panic(err)
	}
	cli.RenderCollectionTable(data)
	//fmt.Printf("%+v\n", data)
	//client.SaveToken()
}
