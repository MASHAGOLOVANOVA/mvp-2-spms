package domainaggregate

type Student struct {
	Person
	//EnrollmentYear uint
	EducationalProgrammeId string
	Cource                 uint
}

type StudentAccount struct {
	Id         string
	Login      string
	StudentId  string
	University string
	EdProgName string
}
