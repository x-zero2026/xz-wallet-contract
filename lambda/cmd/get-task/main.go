package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/x-zero/xz-wallet/pkg/db"
	"github.com/x-zero/xz-wallet/pkg/models"
	"github.com/x-zero/xz-wallet/pkg/response"
)

type GetTaskResponse struct {
	Task        models.Task              `json:"task"`
	Creator     UserInfo                 `json:"creator"`
	Executor    *UserInfo                `json:"executor,omitempty"`
	Submissions []models.TaskSubmission  `json:"submissions"`
	Bids        []BidInfo                `json:"bids,omitempty"`
}

type UserInfo struct {
	DID         string `json:"did"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	CreditScore int    `json:"credit_score"`
}

type BidInfo struct {
	models.TaskBid
	BidderUsername  string `json:"bidder_username"`
	BidderCreditScore int  `json:"bidder_credit_score"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type,Authorization",
				"Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
			},
		}, nil
	}

	taskID := request.PathParameters["id"]
	if taskID == "" {
		return response.Error(400, "Missing task ID")
	}

	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	pool := db.GetPool()

	// Get task
	var task models.Task
	err := pool.QueryRow(ctx, `
		SELECT task_id, contract_task_id, project_id, creator_did, executor_did,
		       task_name, task_description, acceptance_criteria,
		       reward_amount, paid_amount, visibility, status,
		       created_at, updated_at, completed_at, cancelled_at
		FROM tasks WHERE task_id = $1
	`, taskID).Scan(
		&task.TaskID, &task.ContractTaskID, &task.ProjectID, &task.CreatorDID, &task.ExecutorDID,
		&task.TaskName, &task.TaskDescription, &task.AcceptanceCriteria,
		&task.RewardAmount, &task.PaidAmount, &task.Visibility, &task.Status,
		&task.CreatedAt, &task.UpdatedAt, &task.CompletedAt, &task.CancelledAt,
	)
	if err != nil {
		return response.Error(404, "Task not found")
	}

	// Get creator info
	var creator UserInfo
	err = pool.QueryRow(ctx, `
		SELECT did, username, email, credit_score FROM users WHERE did = $1
	`, task.CreatorDID).Scan(&creator.DID, &creator.Username, &creator.Email, &creator.CreditScore)
	if err != nil {
		return response.Error(500, "Failed to get creator info")
	}

	// Get executor info if exists
	var executor *UserInfo
	if task.ExecutorDID != nil {
		var exec UserInfo
		err = pool.QueryRow(ctx, `
			SELECT did, username, email, credit_score FROM users WHERE did = $1
		`, *task.ExecutorDID).Scan(&exec.DID, &exec.Username, &exec.Email, &exec.CreditScore)
		if err == nil {
			executor = &exec
		}
	}

	// Get submissions
	rows, err := pool.Query(ctx, `
		SELECT submission_id, task_id, submission_type, content, file_urls,
		       status, rejection_reason, submitted_at, reviewed_at
		FROM task_submissions WHERE task_id = $1 ORDER BY submitted_at DESC
	`, taskID)
	if err != nil {
		return response.Error(500, "Failed to get submissions")
	}
	defer rows.Close()

	submissions := []models.TaskSubmission{}
	for rows.Next() {
		var sub models.TaskSubmission
		err := rows.Scan(
			&sub.SubmissionID, &sub.TaskID, &sub.SubmissionType, &sub.Content, &sub.FileURLs,
			&sub.Status, &sub.RejectionReason, &sub.SubmittedAt, &sub.ReviewedAt,
		)
		if err == nil {
			submissions = append(submissions, sub)
		}
	}

	// Get bids (only for creator)
	bids := []BidInfo{}
	bidRows, err := pool.Query(ctx, `
		SELECT tb.bid_id, tb.task_id, tb.bidder_did, tb.bid_message,
		       tb.credit_score_snapshot, tb.status, tb.created_at, tb.updated_at,
		       u.username, u.credit_score
		FROM task_bids tb
		JOIN users u ON tb.bidder_did = u.did
		WHERE tb.task_id = $1
		ORDER BY tb.created_at DESC
	`, taskID)
	if err == nil {
		defer bidRows.Close()
		for bidRows.Next() {
			var bid BidInfo
			err := bidRows.Scan(
				&bid.BidID, &bid.TaskID, &bid.BidderDID, &bid.BidMessage,
				&bid.CreditScoreSnapshot, &bid.Status, &bid.CreatedAt, &bid.UpdatedAt,
				&bid.BidderUsername, &bid.BidderCreditScore,
			)
			if err == nil {
				bids = append(bids, bid)
			}
		}
	}

	return response.Success(GetTaskResponse{
		Task:        task,
		Creator:     creator,
		Executor:    executor,
		Submissions: submissions,
		Bids:        bids,
	})
}

func main() {
	lambda.Start(handler)
}
