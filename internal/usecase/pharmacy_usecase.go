package usecase

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/logistic"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
	"github.com/shopspring/decimal"
)

type PharmacyUseCase interface {
	ListLogistics(ctx context.Context, addressID, pharmacyID int, weight string) (*dto.ListPharmacyLogistics, error)
}

type pharmacyUseCase struct {
	store            repository.Store
	logisticProvider logistic.Provider
}

func NewPharmacyUseCase(store repository.Store, logisticProvider logistic.Provider) PharmacyUseCase {
	return &pharmacyUseCase{
		store:            store,
		logisticProvider: logisticProvider,
	}
}

func (u *pharmacyUseCase) ListLogistics(ctx context.Context, addressID, pharmacyID int, weight string) (*dto.ListPharmacyLogistics, error) {
	logistics := []dto.Logistic{}

	repository := u.store.PharmacyLogistic()
	pharmacyLogistics, err := repository.FindAllByPharmacyID(ctx, addressID, pharmacyID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, pharmaLogistic := range pharmacyLogistics {
		wg.Add(1)
		go func(i int, pharmaLogistic model.PharmacyLogistic, wg *sync.WaitGroup) {
			defer wg.Done()
			var price decimal.Decimal
			if pharmaLogistic.Logistics.Name == appconst.LogisticNameOfficial {
				price = pharmaLogistic.Logistics.PricePerKilometers.Decimal.Mul(pharmaLogistic.DistanceKM.RoundUp(-2))
			} else {
				price = decimal.NewFromInt(int64(u.logisticProvider.Cost(pharmaLogistic.CustomerCityID, pharmaLogistic.PharmacyCityID, weight, pharmaLogistic.Logistics.Name)))
			}
			eda := time.Now().Add(time.Hour * 24 * time.Duration(pharmaLogistic.Logistics.EDA))
			mu.Lock()
			logistics = append(logistics, dto.Logistic{
				ID:            pharmaLogistic.Logistics.ID,
				Name:          pharmaLogistic.Logistics.Name + " " + pharmaLogistic.Logistics.Service,
				LogoUrl:       pharmaLogistic.Logistics.LogoUrl,
				Service:       pharmaLogistic.Logistics.Service,
				Estimation:    eda.Format("Mon, 02 Jan 2006"),
				Price:         price.RoundCeil(-2),
				IsRecommended: false,
			})
			mu.Unlock()
		}(i, pharmaLogistic, &wg)
	}
	wg.Wait()
	sort.Slice(logistics, func(i, j int) bool {
		return logistics[i].Price.LessThan(logistics[j].Price)
	})
	logistics[0].IsRecommended = true
	return &dto.ListPharmacyLogistics{
		PharmacyID:   pharmacyID,
		PharmacyName: pharmacyLogistics[0].Pharmacies.Name,
		Logistics:    logistics,
	}, nil
}
