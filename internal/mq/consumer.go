package mq

import (
	"context"

	"github.com/hibiken/asynq"
)

type Consumer interface {
	Start(ctx context.Context) error
	SendVerificationEmail(ctx context.Context, task *asynq.Task) error
	SendPasswordResetEmail(ctx context.Context, task *asynq.Task) error
	SendPharmacistCredentials(ctx context.Context, task *asynq.Task) error
	ApprovePaymentProof(ctx context.Context, task *asynq.Task) error
	ConfirmUserOrder(ctx context.Context, task *asynq.Task) error
	UpdatePartnerDaysAndHours(ctx context.Context, task *asynq.Task) error
	Shutdown()
}
