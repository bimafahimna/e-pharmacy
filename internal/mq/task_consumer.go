package mq

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/logger"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/mailer"
	"github.com/hibiken/asynq"
)

type taskConsumer struct {
	server *asynq.Server
	mailer mailer.Mailer
	config *config.Config
}

func NewTaskConsumer(opt asynq.RedisClientOpt, mailer mailer.Mailer, config *config.Config) Consumer {
	server := asynq.NewServer(opt, asynq.Config{
		Queues: map[string]int{
			appconst.QueueDefault:  30,
			appconst.QueueCritical: 15,
			appconst.QueueLow:      5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			logger.Log.WithFields(map[string]interface{}{
				"ERROR":   err,
				"TYPE":    task.Type(),
				"PAYLOAD": task.Payload(),
			}).Info("FAILED TO PROCESS TASK")
		}),
		Logger: logger.Log,
	})

	return &taskConsumer{
		server: server,
		mailer: mailer,
		config: config,
	}
}

func (c *taskConsumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			mux := asynq.NewServeMux()

			mux.HandleFunc(taskSendVerificationEmail, c.SendVerificationEmail)
			mux.HandleFunc(taskSendPasswordResetEmail, c.SendPasswordResetEmail)
			mux.HandleFunc(taskSendPharmacistCredentials, c.SendPharmacistCredentials)
			mux.HandleFunc(taskApprovePaymentProof, c.ApprovePaymentProof)
			mux.HandleFunc(taskConfirmUserOrder, c.ConfirmUserOrder)
			mux.HandleFunc(taskUpdatePartnerDaysAndHours, c.UpdatePartnerDaysAndHours)

			return c.server.Start(mux)
		}
	}
}

func (c *taskConsumer) SendEmail(ctx context.Context, task *asynq.Task) error {
	if err := c.mailer.Send(); err != nil {
		return apperror.ErrInternalServerError
	}

	logger.Log.WithFields(map[string]interface{}{
		"TYPE":    task.Type(),
		"PAYLOAD": string(task.Payload()),
	}).Info("CONSUMED TASK")

	return nil
}

func (c *taskConsumer) SendVerificationEmail(ctx context.Context, task *asynq.Task) error {
	var payload VerificationPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return apperror.ErrBadRequest
	}

	subject := "[Puxing] Verify your account"
	htmlBody := fmt.Sprintf(`
    <html>
    <body>
        <p>Please verify your account through this link:</p>
        <a href="%s/auth/verify-account?verif_token=%s">Verify Account</a>
    </body>
    </html>`, c.mailer.Host(), payload.Token)

	c.mailer.Set(payload.Email, subject, htmlBody)

	return c.SendEmail(ctx, task)
}

func (c *taskConsumer) SendPasswordResetEmail(ctx context.Context, task *asynq.Task) error {
	var payload PasswordResetPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return apperror.ErrBadRequest
	}

	subject := "[Puxing] Reset your password"
	htmlBody := fmt.Sprintf(`
    <html>
    <body>
        <p>Please follow this link to reset your password:</p>
        <a href="%s/auth/reset-password?reset_token=%s">Reset Password</a>
    </body>
    </html>`, c.mailer.Host(), payload.Token)

	c.mailer.Set(payload.Email, subject, htmlBody)

	return c.SendEmail(ctx, task)
}

func (c *taskConsumer) SendPharmacistCredentials(ctx context.Context, task *asynq.Task) error {
	var payload PharmacistCredentialsPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return apperror.ErrBadRequest
	}

	subject := "[Puxing] Pharmacist credentials"
	htmlBody := fmt.Sprintf(`
    <html>
    <body>
        <p>Your pharmacist credentials:</p>
        <p>Email: %s</p>
        <p>Password: %s</p>
    </body>
    </html>`, payload.Email, payload.Password)

	c.mailer.Set(payload.Email, subject, htmlBody)

	return c.SendEmail(ctx, task)
}

func (c *taskConsumer) ApprovePaymentProof(ctx context.Context, task *asynq.Task) error {
	var payload ApprovePaymentProofPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return apperror.ErrBadRequest
	}

	url := fmt.Sprintf("%s/orders/payments/%d", c.config.URL.Backend, payload.PaymentID)

	jsonPayload, err := json.Marshal(map[string]string{"status": payload.Status})
	if err != nil {
		return apperror.ErrInternalServerError
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return apperror.ErrInternalServerError
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Secret-Key", base64.RawURLEncoding.EncodeToString([]byte("worker-secret-key")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return apperror.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apperror.ErrInternalServerError
	}

	logger.Log.WithFields(map[string]interface{}{
		"TYPE":    task.Type(),
		"PAYLOAD": string(task.Payload()),
	}).Info("CONSUMED TASK")

	return nil
}

func (c *taskConsumer) ConfirmUserOrder(ctx context.Context, task *asynq.Task) error {
	var payload ConfirmUserOrderPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return apperror.ErrBadRequest
	}

	url := fmt.Sprintf("%s/orders/%d", c.config.URL.Backend, payload.OrderID)

	jsonPayload, err := json.Marshal(map[string]string{"status": payload.Status})
	if err != nil {
		return apperror.ErrInternalServerError
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return apperror.ErrInternalServerError
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Secret-Key", base64.RawURLEncoding.EncodeToString([]byte("worker-secret-key")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return apperror.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apperror.ErrInternalServerError
	}

	logger.Log.WithFields(map[string]interface{}{
		"TYPE":    task.Type(),
		"PAYLOAD": string(task.Payload()),
	}).Info("CONSUMED TASK")

	return nil
}

func (c *taskConsumer) UpdatePartnerDaysAndHours(ctx context.Context, task *asynq.Task) error {
	var payload UpdatePartnerDaysAndHoursPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return apperror.ErrBadRequest
	}

	url := fmt.Sprintf("%s/worker/partners/%d", c.config.URL.Backend, payload.ID)

	jsonPayload, err := json.Marshal(map[string]interface{}{
		"active_days":       payload.ActiveDays,
		"operational_start": payload.OperationalStart,
		"operational_stop":  payload.OperationalStop,
	})
	if err != nil {
		return apperror.ErrInternalServerError
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return apperror.ErrInternalServerError
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Secret-Key", base64.RawURLEncoding.EncodeToString([]byte("worker-secret-key")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return apperror.ErrInternalServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apperror.ErrInternalServerError
	}

	logger.Log.WithFields(map[string]interface{}{
		"TYPE":    task.Type(),
		"PAYLOAD": string(task.Payload()),
	}).Info("CONSUMED TASK")

	return nil
}

func (c *taskConsumer) Shutdown() {
	c.server.Shutdown()
}
