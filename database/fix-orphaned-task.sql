-- Fix orphaned task (created on blockchain but not in database)
-- Task ID 1: 100 XZT task created by admin

-- First, check if task already exists
SELECT task_id, contract_task_id, task_name, status 
FROM tasks 
WHERE contract_task_id = 1;

-- If not exists, insert it
-- You need to provide: project_id, task_name, task_description, acceptance_criteria
-- Replace the values below with actual task details

INSERT INTO tasks (
    contract_task_id,
    project_id,
    creator_did,
    task_name,
    task_description,
    acceptance_criteria,
    reward_amount,
    visibility,
    status,
    created_at,
    updated_at
) VALUES (
    1,  -- contract_task_id from blockchain
    '240e7d39-2d0a-4c15-bf89-f2aaabd40733',  -- system project_id
    '0x3070deb1c17432b094d30509ccbfd598fb2793435efdca9273dfbc558bc040ca',  -- admin DID
    'Recovered Task #1',  -- task_name (you should update this)
    'This task was created on blockchain but failed to save in database',  -- task_description
    'N/A - Please update acceptance criteria',  -- acceptance_criteria
    '100.00',  -- reward_amount (100 XZT)
    'global',  -- visibility
    'pending',  -- status
    NOW(),
    NOW()
)
ON CONFLICT (contract_task_id) DO NOTHING
RETURNING task_id, contract_task_id, task_name;

-- Verify the insert
SELECT 
    task_id,
    contract_task_id,
    task_name,
    reward_amount,
    status,
    created_at
FROM tasks 
WHERE contract_task_id = 1;
