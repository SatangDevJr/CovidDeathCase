package deathcase

import (
	"covid/src/pkg/entity"
	"covid/src/pkg/external/ddc"
	covidError "covid/src/pkg/utils/error"
	"covid/src/pkg/utils/logger"
	"fmt"
	"sort"
)

type UseCase interface {
	GetTopThreeProvice()([]entity.TopThreeDeathCase,*covidError.ErrorCode)
}

type Service struct {
	UseCase
	DDCService ddc.UseCase
	Logs       logger.Logger
}

func NewService(ddcService ddc.UseCase, logs logger.Logger) *Service {
	service := &Service{
		DDCService: ddcService,
		Logs:       logs,
	}
	service.UseCase = service
	return service
}

func (service *Service) GetTopThreeProvice()([]entity.TopThreeDeathCase,*covidError.ErrorCode){

	//get data assets from external service
	resGetDeathCaseRound4, errGetDeathCaseRound4 := service.DDCService.GetDeathCaseRound4();
	if errGetDeathCaseRound4 != nil {
		go service.Logs.Error("", "deathcase_service_get_death_case_round_4_error", "", errGetDeathCaseRound4)
		return nil, errGetDeathCaseRound4
	}

	//key string value int64 to init store case
	provinceCount := make(map[string]int64)

	//count case
	for _, count := range resGetDeathCaseRound4 {
		provinceCount[count.Province]++
	}

	//mapping key value to enity model
	var provinceCounts []entity.TopThreeDeathCase
	for key, value := range provinceCount {
		provinceCounts = append(provinceCounts, entity.TopThreeDeathCase{Province: key, Cases: value})
	}

	//sort
	sort.Slice(provinceCounts, func(indexStart, indexNextCompare int) bool {
		return provinceCounts[indexStart].Cases > provinceCounts[indexNextCompare].Cases
	})

	//cut only top threee
	top3Provinces := provinceCounts
	if len(provinceCounts) > 3 {
		top3Provinces = provinceCounts[:3]
	}

	//println on pod for show
	fmt.Println("Top 3 Provinces with the highest number of cases:")
	for _, entry := range top3Provinces {
		fmt.Printf("%s: %d cases\n", entry.Province, entry.Cases)
	}

	return top3Provinces,nil
}