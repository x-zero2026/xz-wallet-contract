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

type ListTasksResponse struct {
	Tasks []TaskWithDetails `json:"tasks"`
	Total int               `json:"total"`
}

type TaskWithDetails struct {
	models.Task
	CreatorUsername  string  `json:"creator_username"`
	ExecutorUsername *string `json:"executor_username,omitempty"`
	BidCount         int     `json:"bid_count"`
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

	if err := db.InitDB(); err != nil {
		return response.Error(500, fmt.Sprintf("Database error: %v", err))
	}

	pool := db.GetPool()

	// Get query parameters
	visibility := request.QueryStringParameters["visibility"]
	projectID := request.QueryStringParameters["project_id"]
	status := request.QueryStringParameters["status"]
	executorDID := request.QueryStringParameters["executor_did"]
	creatorDID := request.QueryStringParameters["creator_did"]
	bidderDID := request.QueryStringParameters["bidder_did"]

	// Build query
	query := `
		SELECT 
			t.task_id, t.contract_task_id, t.project_id, t.creator_did, t.executor_did,
			t.task_name, t.task_description, t.acceptance_criteria,
			t.reward_amount, t.paid_amount, t.visibility, t.status, t.profession_tags,
			t.created_at, t.updated_at,
			u_creator.username as creator_username,
			u_executor.username as executor_username,
			COUNT(DISTINCT tb.bid_id) as bid_count
		FROM tasks t
		LEFT JOIN users u_creator ON t.creator_did = u_creator.did
		LEFT JOIN users u_executor ON t.executor_did = u_executor.did
		LEFT JOIN task_bids tb ON t.task_id = tb.task_id AND tb.status = 'pending'
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if visibility != "" {
		query += fmt.Sprintf(" AND t.visibility = $%d", argCount)
		args = append(args, visibility)
		argCount++
	}
	if projectID != "" {
		query += fmt.Sprintf(" AND t.project_id = $%d", argCount)
		args = append(args, projectID)
		argCount++
	}
	if status != "" {
		query += fmt.Sprintf(" AND t.status = $%d", argCount)
		args = append(args, status)
		argCount++
	}
	if executorDID != "" {
		query += fmt.Sprintf(" AND t.executor_did = $%d", argCount)
		args = append(args, executorDID)
		argCount++
	}
	if creatorDID != "" {
		query += fmt.Sprintf(" AND t.creator_did = $%d", argCount)
		args = append(args, creatorDID)
		argCount++
	}
	if bidderDID != "" {
		// Filter tasks where user has placed a bid
		query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM task_bids WHERE task_id = t.task_id AND bidder_did = $%d)", argCount)
		args = append(args, bidderDID)
		argCount++
	}

	query += " GROUP BY t.task_id, u_creator.username, u_executor.username ORDER BY t.created_at DESC"

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return response.Error(500, fmt.Sprintf("Query failed: %v", err))
	}
	defer rows.Close()

	tasks := []TaskWithDetails{}
	for rows.Next() {
		var task TaskWithDetails
		err := rows.Scan(
			&task.TaskID, &task.ContractTaskID, &task.ProjectID, &task.CreatorDID, &task.ExecutorDID,
			&task.TaskName, &task.TaskDescription, &task.AcceptanceCriteria,
			&task.RewardAmount, &task.PaidAmount, &task.Visibility, &task.Status, &task.ProfessionTags,
			&task.CreatedAt, &task.UpdatedAt,
			&task.CreatorUsername, &task.ExecutorUsername, &task.BidCount,
		)
		if err != nil {
			return response.Error(500, fmt.Sprintf("Scan failed: %v", err))
		}
		tasks = append(tasks, task)
	}

	return response.Success(ListTasksResponse{
		Tasks: tasks,
		Total: len(tasks),
	})
}

func main() {
	lambda.Start(handler)
}
