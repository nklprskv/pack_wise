package health

import "context"

type Status struct {
	Status string `json:"status" example:"ok" doc:"Application health status"`
}

type GetHealthOutput struct {
	Body Status
}

func (c *Controller) HealthCheck(_ context.Context, _ *struct{}) (*GetHealthOutput, error) {
	return &GetHealthOutput{Body: Status{Status: "ok"}}, nil
}
