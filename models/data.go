package models

type Armourments struct{
	Id int64 `json:"id"`
	Title string `json:"title"`
	Qty int64 `json:"qty"`/// in case its a deathstar
}

type Spacecraft struct{
	Id int64 `json:"id"`
	Name string `json:"name,omitempty"`
	Class string `json:"class,omitempty"`
	Status string `json:"status,omitempty"`
	Crew int32 `json:"crew,omitempty"`

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
