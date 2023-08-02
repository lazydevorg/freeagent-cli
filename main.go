package main

import (
	"context"
	"github.com/lazydevorg/freeagent-cli/internal/client"
	"github.com/lazydevorg/freeagent-cli/internal/client/timeslip"
)

func main() {
	//cmd.Execute()
	//auth.Authenticate(true)

	c := client.NewClient(context.Background())
	timeslips := timeslip.Timeslips(c)
	weekTimeslips, err := timeslips.GetWeek()
	if err != nil {
		panic(err)
	}
	related := make(map[string]string)
	_ = timeslips.GetRelated(weekTimeslips, related)
	timeslips.PrintTable(weekTimeslips, related)

	//data, err := client.GetActiveProjects()
	//data, err := client.GetActiveTasks()
	//if err != nil {
	//	panic(err)
	//}
	//cli.RenderEntityTable(data)
	//fmt.Printf("%+v\n", data)
	//client.SaveToken()
}
