package logistic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
)

type rajaOngkirProvider struct {
	url    string
	apiKey string
}

func NewRajaOngkirProvider(config config.LogisticConfig) Provider {
	return &rajaOngkirProvider{
		url:    config.URL,
		apiKey: config.ApiKey,
	}
}

func (p *rajaOngkirProvider) Cost(cityOriginID, cityDestinationID int, weight, logisticName string) int {
	payload := strings.NewReader(fmt.Sprintf("origin=%d&destination=%d&weight=%s&courier=%s", cityOriginID, cityDestinationID, weight, "jne"))
	req, _ := http.NewRequest(http.MethodPost, p.url, payload)
	req.Header.Add("key", p.apiKey)
	req.Header.Add("content-type", appconst.ContentTypeFormUrlEncoded)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1
	}
	if res.StatusCode != http.StatusOK {
		return -1
	}
	defer res.Body.Close()
	var data struct {
		RajaOngkir dto.RajaOngkirResponse `json:"rajaongkir"`
	}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return -1
	}
	var cost int
	for _, result := range data.RajaOngkir.Results[0].Costs {
		if result.Service == appconst.LogisticServiceYES {
			cost = result.Cost[0].Value
			break
		}
	}
	return cost
}
