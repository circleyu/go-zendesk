package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	JobCompleted = "completed"
	JobQueued    = "queued"
	JobWorking   = "working"
	JobFailed    = "failed"
	JobKilled    = "killed"
)

type JobStatus struct {
	ID       string `json:"id,omitempty"`
	Message  string `json:"message,omitempty"`
	Progress int64  `json:"progress,omitempty"`

	Results []struct {
		ID      int64  `json:"id,omitempty"`
		Index   int64  `json:"index,omitempty"`
		Action  string `json:"action,omitempty"`
		Status  string `json:"status,omitempty"`
		Success bool   `json:"success,omitempty"`
	} `json:"results,omitempty"`

	Status string `json:"status,omitempty"`
	Total  int64  `json:"total,omitempty"`
	URL    string `json:"url,omitempty"`
}

type JobAPI interface {
	ListJobStatuses(ctx context.Context) (JobStatus, error)
	ShowManyJobStatuses(ctx context.Context, ids []string) (JobStatus, error)
	ShowJobStatuses(ctx context.Context, id string) (JobStatus, error)
}

func (z *Client) ListJobStatuses(ctx context.Context) (JobStatus, error) {
	var result struct {
		JobStatus JobStatus `json:"job_status"`
	}

	body, err := z.get(ctx, "/job_statuses.json")
	if err != nil {
		return JobStatus{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return JobStatus{}, err
	}

	return result.JobStatus, nil
}

func (z *Client) ShowManyJobStatuses(ctx context.Context, ids []string) (JobStatus, error) {
	var result struct {
		JobStatus JobStatus `json:"job_status"`
	}

	var req struct {
		IDs string `url:"ids,omitempty"`
	}
	req.IDs = strings.Join(ids, ",")

	u, err := addOptions("/job_statuses/show_many.json", req)
	if err != nil {
		return JobStatus{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return JobStatus{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return JobStatus{}, err
	}

	return result.JobStatus, nil
}

func (z *Client) ShowJobStatuses(ctx context.Context, id string) (JobStatus, error) {
	var result struct {
		JobStatus JobStatus `json:"job_status"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/job_statuses/%s.json", id))
	if err != nil {
		return JobStatus{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return JobStatus{}, err
	}

	return result.JobStatus, nil
}
