package mq

import (
	"context"

	"github.com/hibiken/asynq"
)

type Producer interface {
	ProduceVerificationTask(ctx context.Context, payload *VerificationPayload, opts ...asynq.Option) error
	ProducePasswordResetTask(ctx context.Context, payload *PasswordResetPayload, opts ...asynq.Option) error
	ProducePharmacistCredentialsTask(ctx context.Context, payload *PharmacistCredentialsPayload, opts ...asynq.Option) error
	ProduceApprovePaymentProof(ctx context.Context, payload *ApprovePaymentProofPayload, opts ...asynq.Option) error
	ProduceConfirmUserOrder(ctx context.Context, payload *ConfirmUserOrderPayload, opts ...asynq.Option) error
	ProduceUpdateDaysAndHours(ctx context.Context, payload *UpdatePartnerDaysAndHoursPayload, opts ...asynq.Option) error
}
