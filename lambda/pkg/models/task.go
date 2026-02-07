package models

import (
	"time"
)

// Task represents a task in the database
type Task struct {
	TaskID          string    `json:"task_id"`
	ContractTaskID  *int64    `json:"contract_task_id,omitempty"`
	ProjectID       string    `json:"project_id"`
	CreatorDID      string    `json:"creator_did"`
	ExecutorDID     *string   `json:"executor_did,omitempty"`
	TaskName        string    `json:"task_name"`
	TaskDescription string    `json:"task_description"`
	AcceptanceCriteria string `json:"acceptance_criteria"`
	RewardAmount    string    `json:"reward_amount"`
	PaidAmount      string    `json:"paid_amount"`
	Visibility      string    `json:"visibility"`
	Status          string    `json:"status"`
	ProfessionTags  []string  `json:"profession_tags,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	CancelledAt     *time.Time `json:"cancelled_at,omitempty"`
}

// TaskBid represents a bid on a task
type TaskBid struct {
	BidID              string    `json:"bid_id"`
	TaskID             string    `json:"task_id"`
	BidderDID          string    `json:"bidder_did"`
	BidMessage         *string   `json:"bid_message,omitempty"`
	CreditScoreSnapshot int      `json:"credit_score_snapshot"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// TaskSubmission represents a work submission
type TaskSubmission struct {
	SubmissionID   string     `json:"submission_id"`
	TaskID         string     `json:"task_id"`
	SubmissionType string     `json:"submission_type"`
	Content        string     `json:"content"`
	FileURLs       []string   `json:"file_urls,omitempty"`
	Status         string     `json:"status"`
	RejectionReason *string   `json:"rejection_reason,omitempty"`
	SubmittedAt    time.Time  `json:"submitted_at"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
}

// CreditHistory represents credit score changes
type CreditHistory struct {
	HistoryID    string    `json:"history_id"`
	UserDID      string    `json:"user_did"`
	TaskID       *string   `json:"task_id,omitempty"`
	ChangeAmount int       `json:"change_amount"`
	Reason       string    `json:"reason"`
	BeforeScore  int       `json:"before_score"`
	AfterScore   int       `json:"after_score"`
	CreatedAt    time.Time `json:"created_at"`
}

// User represents a user with wallet info
type User struct {
	DID            string  `json:"did"`
	EthAddress     string  `json:"eth_address"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	CreditScore    int     `json:"credit_score"`
	TasksCompleted int     `json:"tasks_completed"`
	TasksCancelled int     `json:"tasks_cancelled"`
	XZTBalance     string  `json:"xzt_balance"`
	EscrowApproved bool    `json:"escrow_approved"`
	CreatedAt      time.Time `json:"created_at"`
}

// TaskStatus constants
const (
	TaskStatusPending                  = "pending"
	TaskStatusBidding                  = "bidding"
	TaskStatusAccepted                 = "accepted"
	TaskStatusDesignSubmitted          = "design_submitted"
	TaskStatusDesignApproved           = "design_approved"
	TaskStatusImplementationSubmitted  = "implementation_submitted"
	TaskStatusImplementationApproved   = "implementation_approved"
	TaskStatusFinalSubmitted           = "final_submitted"
	TaskStatusCompleted                = "completed"
	TaskStatusCancelled                = "cancelled"
)

// SubmissionType constants
const (
	SubmissionTypeDesign         = "design"
	SubmissionTypeImplementation = "implementation"
	SubmissionTypeFinal          = "final"
)

// BidStatus constants
const (
	BidStatusPending  = "pending"
	BidStatusAccepted = "accepted"
	BidStatusRejected = "rejected"
)

// Milestone payment percentages (in basis points, 10000 = 100%)
const (
	MilestoneDesign         = 3000 // 30%
	MilestoneImplementation = 5000 // 50% (cumulative 80%)
	MilestoneFinal          = 2000 // 20% (cumulative 100%)
)
