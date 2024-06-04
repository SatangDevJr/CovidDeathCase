package entity

type CovidCase struct {
	Year         int     `json:"year"`
	Weeknum      int     `json:"weeknum"`
	Province     string  `json:"province"`
	Age          string  `json:"age"`
	AgeRange     string  `json:"age_range"`
	Occupation   string  `json:"occupation"`
	Type         string  `json:"type"`
	DeathCluster *string `json:"death_cluster"`
	UpdateDate   string  `json:"update_date"`
}