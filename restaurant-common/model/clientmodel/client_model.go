package clientmodel

type EnterRestaurantRequest struct {
	ClientName string `json:"clientName"`
	ClientId   string `json:"clientId"`
}

type EnterRestaurantResponseStatus string

const (
	EnterRestaurantResponseStatusWelcome = "Welcome"
	EnterRestaurantResponseStatusReject  = "Rejected"
)

type EnterRestaurantResponse struct {
	ClientId    string                        `json:"clientId"`
	Message     string                        `json:"message"`
	Status      EnterRestaurantResponseStatus `json:"status"`
	TableNumber int                           `json:"tableNumber"`
}
