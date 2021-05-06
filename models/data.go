package models

type Armourments struct{
	Id int64
	Title string
	Qty int64 /// in case its a deathstar
}

type Spacecraft struct{
	Id int64
	Name string
	Class string
	Status string
	Crew int32

	Armourments []Armourments
}

type ErrorAlreadyExists struct{
	details string
}

func NewErrorAlreadyExists(details string)ErrorAlreadyExists{
	return ErrorAlreadyExists{
		details: details,
	}
}

func (e ErrorAlreadyExists) Error() string {
	return "The entry already exists" + e.details
}

type ErrorNotExist struct{
	details string
}

func NewErrorNotExist(details string)ErrorAlreadyExists{
	return ErrorAlreadyExists{
		details: details,
	}
}

func (e ErrorNotExist) Error() string {
	return "The entry already exists" + e.details
}
