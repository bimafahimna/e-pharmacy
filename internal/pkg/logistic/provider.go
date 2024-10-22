package logistic

type Provider interface {
	Cost(cityOriginID, cityDestinationID int, weight, logisticName string) int
}
