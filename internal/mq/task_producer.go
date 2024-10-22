package mq

import (
	"context"
	"encoding/json"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/logger"
	"github.com/hibiken/asynq"
)

type taskProducer struct {
	client *asynq.Client
}

func NewTaskProducer(opt asynq.RedisClientOpt) Producer {
	client := asynq.NewClient(opt)
	return &taskProducer{
		client: client,
	}
}

func (p *taskProducer) ProduceTask(ctx context.Context, taskType string, payload interface{}, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return apperror.ErrBadRequest
	}

	task := asynq.NewTask(taskType, jsonPayload, opts...)

	info, err := p.client.EnqueueContext(ctx, task)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	logger.Log.WithFields(map[string]interface{}{
		"TYPE":      task.Type(),
		"PAYLOAD":   string(task.Payload()),
		"QUEUE":     info.Queue,
		"MAX_RETRY": info.MaxRetry,
	}).Info("PRODUCED TASK")

	return nil
}

func (p *taskProducer) ProduceVerificationTask(ctx context.Context, payload *VerificationPayload, opts ...asynq.Option) error {
	return p.ProduceTask(ctx, taskSendVerificationEmail, payload, opts...)
}

func (p *taskProducer) ProducePasswordResetTask(ctx context.Context, payload *PasswordResetPayload, opts ...asynq.Option) error {
	return p.ProduceTask(ctx, taskSendPasswordResetEmail, payload, opts...)
}

func (p *taskProducer) ProducePharmacistCredentialsTask(ctx context.Context, payload *PharmacistCredentialsPayload, opts ...asynq.Option) error {
	return p.ProduceTask(ctx, taskSendPharmacistCredentials, payload, opts...)
}

func (p *taskProducer) ProduceApprovePaymentProof(ctx context.Context, payload *ApprovePaymentProofPayload, opts ...asynq.Option) error {
	return p.ProduceTask(ctx, taskApprovePaymentProof, payload, opts...)
}

func (p *taskProducer) ProduceUpdateDaysAndHours(ctx context.Context, payload *UpdatePartnerDaysAndHoursPayload, opts ...asynq.Option) error {
	return p.ProduceTask(ctx, taskApprovePaymentProof, payload, opts...)
}

func (p *taskProducer) ProduceConfirmUserOrder(ctx context.Context, payload *ConfirmUserOrderPayload, opts ...asynq.Option) error {
	return p.ProduceTask(ctx, taskConfirmUserOrder, payload, opts...)
}
