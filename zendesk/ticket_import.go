package zendesk

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type TicketImport struct {
	ID              int64         `json:"id,omitempty"`
	URL             string        `json:"url,omitempty"`
	ExternalID      string        `json:"external_id,omitempty"`
	Type            string        `json:"type,omitempty"`
	Subject         string        `json:"subject,omitempty"`
	RawSubject      string        `json:"raw_subject,omitempty"`
	Description     string        `json:"description,omitempty"`
	Priority        string        `json:"priority,omitempty"`
	Status          string        `json:"status,omitempty"`
	Recipient       string        `json:"recipient,omitempty"`
	RequesterID     int64         `json:"requester_id,omitempty"`
	SubmitterID     int64         `json:"submitter_id,omitempty"`
	AssigneeID      int64         `json:"assignee_id,omitempty"`
	OrganizationID  int64         `json:"organization_id,omitempty"`
	GroupID         int64         `json:"group_id,omitempty"`
	CollaboratorIDs []int64       `json:"collaborator_ids,omitempty"`
	FollowerIDs     []int64       `json:"follower_ids,omitempty"`
	EmailCCIDs      []int64       `json:"email_cc_ids,omitempty"`
	ForumTopicID    int64         `json:"forum_topic_id,omitempty"`
	ProblemID       int64         `json:"problem_id,omitempty"`
	HasIncidents    bool          `json:"has_incidents,omitempty"`
	DueAt           time.Time     `json:"due_at,omitempty"`
	Tags            []string      `json:"tags,omitempty"`
	CustomFields    []CustomField `json:"custom_fields,omitempty"`

	EmailCCs []struct {
		ID     int64  `json:"user_id,omitempty"`
		Email  string `json:"user_email,omitempty"`
		Name   string `json:"user_name,omitempty"`
		Action string `json:"action,omitempty"`
	} `json:"email_ccs,omitempty"`

	Via *Via `json:"via,omitempty"`

	SharingAgreementIDs []int64   `json:"sharing_agreement_ids,omitempty"`
	FollowupIDs         []int64   `json:"followup_ids,omitempty"`
	ViaFollowupSourceID int64     `json:"via_followup_source_id,omitempty"`
	MacroIDs            []int64   `json:"macro_ids,omitempty"`
	TicketFormID        int64     `json:"ticket_form_id,omitempty"`
	BrandID             int64     `json:"brand_id,omitempty"`
	AllowChannelback    bool      `json:"allow_channelback,omitempty"`
	AllowAttachments    bool      `json:"allow_attachments,omitempty"`
	IsPublic            bool      `json:"is_public,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`

	// Collaborators is POST only
	Collaborators Collaborators `json:"collaborators,omitempty"`

	// Comment is POST only and required
	Comments []TicketComment `json:"comments,omitempty"`

	// Requester is POST only and can be used to create a ticket for a nonexistent requester
	Requester Requester `json:"requester,omitempty"`
}

type TicketImportOptions struct {
	ArchiveImmediately bool `url:"archive_immediately,omitempty"`
}

type TicketImportAPI interface {
	ImportTicket(ctx context.Context, ticket TicketImport, opts *TicketImportOptions) (Ticket, error)
	BatchImportTickets(ctx context.Context, tickets []TicketImport, opts *TicketImportOptions) (JobStatus, error)
}

func (z *Client) ImportTicket(ctx context.Context, ticket TicketImport, opts *TicketImportOptions) (Ticket, error) {
	var data struct {
		Ticket TicketImport `json:"ticket"`
	}
	data.Ticket = ticket

	var result struct {
		Ticket Ticket `json:"ticket"`
	}

	tmp := opts
	if tmp == nil {
		tmp = &TicketImportOptions{
			ArchiveImmediately: false,
		}
	}

	u, err := addOptions("/imports/tickets.json", tmp)
	if err != nil {
		return Ticket{}, err
	}

	body, err := z.post(ctx, u, data)

	if err != nil {
		return Ticket{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Ticket{}, err
	}
	return result.Ticket, nil
}

func (z *Client) BatchImportTickets(ctx context.Context, tickets []TicketImport, opts *TicketImportOptions) (JobStatus, error) {
	var data struct {
		Tickets []TicketImport `json:"tickets"`
	}
	data.Tickets = tickets

	var result struct {
		JobStatus JobStatus `json:"job_status"`
	}

	tmp := opts
	if tmp == nil {
		tmp = &TicketImportOptions{
			ArchiveImmediately: false,
		}
	}

	u, err := addOptions("/imports/tickets/create_many.json", tmp)
	if err != nil {
		return JobStatus{}, err
	}

	body, err := z.postWithStatus(ctx, u, data, http.StatusOK)

	if err != nil {
		return JobStatus{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return JobStatus{}, err
	}
	return result.JobStatus, nil
}
