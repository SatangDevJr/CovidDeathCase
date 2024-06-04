package ddc

import "covid/src/pkg/entity"

type ResponseFromAPICovid struct {
	Data []entity.CovidCase `json:"data"`
}