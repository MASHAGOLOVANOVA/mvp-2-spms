package people

import "fmt"

type Person struct {
	Name       string
	Surname    string
	Middlename string
}

func (p *Person) FullNameToString() string {
	fullname := fmt.Sprint(p.Name, " ", p.Surname)
	if p.Middlename != "" {
		fullname += fmt.Sprint(" ", p.Middlename)
	}
	return fullname
}
