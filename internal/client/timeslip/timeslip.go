package timeslip

import (
	"fmt"
	"github.com/lazydevorg/freeagent-cli/internal/client"
	"github.com/rodaine/table"
	"time"
)

type Timeslip struct {
	Url       string       `json:"url"`
	Task      string       `json:"task"`
	Project   string       `json:"project"`
	User      string       `json:"user"`
	DateOn    DatedOnField `json:"dated_on"`
	Hours     string       `json:"hours"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type DatedOnField time.Time

func (d *DatedOnField) String() string {
	return time.Time(*d).Format("Mon 2 Jan 2006")
}

func (d *DatedOnField) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("\"2006-01-02\"", string(data))
	if err != nil {
		return fmt.Errorf("error decoding dated_on field with value %s: %w", data, err)
	}
	*d = DatedOnField(t)
	return nil
}

func GetRelated(timeslips []Timeslip, related map[string]string) error {
	err := client.GetRelatedEntities(timeslips, "Project", related, func(entity map[string]interface{}) string {
		return entity["name"].(string)
	})
	if err != nil {
		return fmt.Errorf("error getting related projects: %w", err)
	}
	err = client.GetRelatedEntities(timeslips, "Task", related, func(entity map[string]interface{}) string {
		return entity["name"].(string)
	})
	if err != nil {
		return fmt.Errorf("error getting related tasks: %w", err)
	}
	err = client.GetRelatedEntities(timeslips, "User", related, func(entity map[string]interface{}) string {
		return fmt.Sprintf("%s %s", entity["first_name"], entity["last_name"])
	})
	if err != nil {
		return fmt.Errorf("error getting related users: %w", err)
	}
	return nil
}

func Create(timeslip *Timeslip) (*Timeslip, error) {
	timeslip, err := client.PostEntity("timeslips", "timeslip", timeslip)
	if err != nil {
		return nil, err
	}
	return timeslip, nil
}

func GetWeek() ([]Timeslip, error) {
	now := time.Now().AddDate(0, 0, -10)
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	firstDayOfWeek := now.AddDate(0, 0, offset)
	lastDayOfWeek := firstDayOfWeek.AddDate(0, 0, 6)
	params := map[string]string{
		"view":      "all",
		"from_date": firstDayOfWeek.Format("2006-01-02"),
		"to_date":   lastDayOfWeek.Format("2006-01-02"),
	}
	return client.GetCollection[Timeslip]("timeslips", "timeslips", params)
}

func PrintTable(timeslips []Timeslip, related map[string]string) {
	tbl := table.New("Project", "Task", "User", "Date", "Hours")
	for _, timeslip := range timeslips {
		tbl.AddRow(related[timeslip.Project], related[timeslip.Task], related[timeslip.User], timeslip.DateOn.String(), timeslip.Hours)
	}
	tbl.Print()
}
