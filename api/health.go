package api

import (
	"context"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type HealthCheckRequest struct{}
type HealthCheckResponse struct {
	Status healthStatus
}

type healthStatus string

const (
	Up   healthStatus = "UP"
	Down healthStatus = "DOWN"
)

func (ctrl *controller) HealthUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req HealthCheckRequest,
			resp *HealthCheckResponse,
		) error {
			resp.Status = Up
			return nil
		},
	)

	u.SetName("health")
	u.SetTitle("Health Check")
	u.SetDescription("Health Check for FlowG Node")
	u.SetTags("health")

	u.SetExpectedErrors(status.Internal)

	return u
}
