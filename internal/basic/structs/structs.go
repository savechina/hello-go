package structs

import (
	"fmt"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("basic", "structs", Run)
}

func main() {
	Run()
}

// Run executes the structs chapter examples.
func Run() {
	examples := []string{
		"1) struct literals => " + buildProfile("Alice", 30, "Taipei"),
		"2) methods => " + celebrateBirthday("Bob", 27),
		"3) embedding => " + describePromotion("Carol", 32, "Kaohsiung", "Platform", "Engineer"),
	}

	for _, example := range examples {
		fmt.Println(example)
	}
}

type address struct {
	city string
}

type profile struct {
	name    string
	age     int
	address address
}

func (p profile) summary() string {
	return fmt.Sprintf("%s is %d years old and lives in %s", p.name, p.age, p.address.city)
}

func (p *profile) haveBirthday() {
	p.age++
}

type employee struct {
	profile
	department string
	title      string
}

func (e employee) badge() string {
	return fmt.Sprintf("%s [%s] %s", e.name, e.department, e.title)
}

func (e *employee) promote(newTitle string) {
	e.title = newTitle
}

func buildProfile(name string, age int, city string) string {
	p := profile{
		name: name,
		age:  age,
		address: address{
			city: city,
		},
	}

	return p.summary()
}

func celebrateBirthday(name string, age int) string {
	p := profile{
		name: name,
		age:  age,
		address: address{
			city: "Taichung",
		},
	}

	p.haveBirthday()
	return p.summary()
}

func describePromotion(name string, age int, city string, department string, title string) string {
	e := employee{
		profile: profile{
			name: name,
			age:  age,
			address: address{
				city: city,
			},
		},
		department: department,
		title:      title,
	}

	e.promote("Senior " + title)
	return fmt.Sprintf("%s -> %s", e.badge(), e.summary())
}
