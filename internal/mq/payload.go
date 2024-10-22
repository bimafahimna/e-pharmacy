package mq

type VerificationPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type PasswordResetPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type PharmacistCredentialsPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ApprovePaymentProofPayload struct {
	PaymentID int    `json:"payment_id"`
	Status    string `json:"status"`
}

type ConfirmUserOrderPayload struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
}

type UpdatePartnerDaysAndHoursPayload struct {
	ID               int
	ActiveDays       []int  `json:"active_days"`
	OperationalStart string `json:"operational_start"`
	OperationalStop  string `json:"operational_stop"`
}
