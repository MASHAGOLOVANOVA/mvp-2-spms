package domainaggregate

type Student struct {
	Person
	//EnrollmentYear uint
	EducationalProgramme string
	Course               uint
	University           string
}

type StudentAccount struct {
	Id        string
	Login     string
	StudentId string
}
