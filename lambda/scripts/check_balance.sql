SELECT u.username, u.eth_address 
FROM tasks t 
JOIN users u ON t.executor_did = u.did 
WHERE t.task_id = '231c59ac-f262-4624-b60e-d411386542f7';
