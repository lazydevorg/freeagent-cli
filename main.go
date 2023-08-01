package main

import (
	"github.com/lazydevorg/freeagent-cli/internal/client/timeslip"
)

func main() {
	//cmd.Execute()
	//auth.Authenticate(true)

	timeslips, err := timeslip.GetWeek()
	if err != nil {
		panic(err)
	}
	related := make(map[string]string)
	_ = timeslip.GetRelated(timeslips, related)
	timeslip.PrintTable(timeslips, related)

	//data, err := client.GetActiveProjects()
	//data, err := client.GetActiveTasks()
	//if err != nil {
	//	panic(err)
	//}
	//cli.RenderEntityTable(data)
	//fmt.Printf("%+v\n", data)
	//client.SaveToken()
}
