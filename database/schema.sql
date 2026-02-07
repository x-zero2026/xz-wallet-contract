-- ============================================
-- XZ Wallet Database Schema
-- ============================================

-- Extend users table with wallet and credit info
ALTER TABLE users 
    ADD COLUMN IF NOT EXISTS credit_score INT DEFAULT 5000,
    ADD COLUMN IF NOT EXISTS tasks_completed INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS tasks_cancelled INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS xzt_balance DECIMAL(20, 8) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS escrow_approved BOOLEAN DEFAULT FALSE;

-- Create index for credit score queries
CREATE INDEX IF NOT EXISTS idx_users_credit_score ON users(credit_score);

-- ============================================
-- Tasks Table
-- ============================================
CREATE TABLE IF NOT EXISTS tasks (
    task_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_task_id BIGINT UNIQUE,
    
    -- Relationships
    project_id UUID NOT NULL REFERENCES projects(project_id) ON DELETE CASCADE,
    creator_did VARCHAR(66) NOT NULL REFERENCES users(did),
    executor_did VARCHAR(66) REFERENCES users(did),
    
    -- Task details
    task_name VARCHAR(255) NOT NULL,
    task_description TEXT NOT NULL,
    acceptance_criteria TEXT NOT NULL,
    
    -- Financial
    reward_amount DECIMAL(20, 8) NOT NULL CHECK (reward_amount > 0),
    paid_amount DECIMAL(20, 8) DEFAULT 0 CHECK (paid_amount >= 0),
    
    -- Settings
    visibility VARCHAR(20) NOT NULL CHECK (visibility IN ('project', 'global')),
    profession_tags TEXT[] DEFAULT '{}',
    
    -- Status tracking
    status VARCHAR(30) NOT NULL DEFAULT 'pending' CHECK (status IN (
        'pending',
        'bidding',
        'accepted',
        'design_submitted',
        'design_approved',
        'implementation_submitted',
        'implementation_approved',
        'final_submitted',
        'completed',
        'cancelled'
    )),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP
);

-- Indexes for tasks
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_creator ON tasks(creator_did);
CREATE INDEX IF NOT EXISTS idx_tasks_executor ON tasks(executor_did);
CREATE INDEX IF NOT EXISTS idx_tasks_project ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_visibility ON tasks(visibility);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_profession_tags ON tasks USING GIN(profession_tags);

-- ============================================
-- Task Bids Table
-- ============================================
CREATE TABLE IF NOT EXISTS task_bids (
    bid_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(task_id) ON DELETE CASCADE,
    bidder_did VARCHAR(66) NOT NULL REFERENCES users(did),
    
    -- Bid details
    bid_message TEXT,
    credit_score_snapshot INT NOT NULL,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN (
        'pending',
        'accepted',
        'rejected'
    )),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint: one bid per user per task
    UNIQUE(task_id, bidder_did)
);

-- Indexes for bids
CREATE INDEX IF NOT EXISTS idx_bids_task ON task_bids(task_id);
CREATE INDEX IF NOT EXISTS idx_bids_bidder ON task_bids(bidder_did);
CREATE INDEX IF NOT EXISTS idx_bids_status ON task_bids(status);

-- ============================================
-- Task Submissions Table
-- ============================================
CREATE TABLE IF NOT EXISTS task_submissions (
    submission_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(task_id) ON DELETE CASCADE,
    
    -- Submission details
    submission_type VARCHAR(30) NOT NULL CHECK (submission_type IN (
        'design',
        'implementation',
        'final'
    )),
    content TEXT NOT NULL,
    file_urls TEXT[],
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN (
        'pending',
        'approved',
        'rejected'
    )),
    rejection_reason TEXT,
    
    -- Timestamps
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP
);

-- Indexes for submissions
CREATE INDEX IF NOT EXISTS idx_submissions_task ON task_submissions(task_id);
CREATE INDEX IF NOT EXISTS idx_submissions_type ON task_submissions(submission_type);
CREATE INDEX IF NOT EXISTS idx_submissions_status ON task_submissions(status);

-- ============================================
-- Credit History Table
-- ============================================
CREATE TABLE IF NOT EXISTS credit_history (
    history_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_did VARCHAR(66) NOT NULL REFERENCES users(did),
    task_id UUID REFERENCES tasks(task_id),
    
    -- Change details
    change_amount INT NOT NULL,
    reason VARCHAR(50) NOT NULL CHECK (reason IN (
        'task_completed',
        'task_cancelled_design',
        'task_cancelled_implementation',
        'manual_adjustment'
    )),
    
    -- Scores
    before_score INT NOT NULL,
    after_score INT NOT NULL,
    
    -- Timestamp
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for credit history
CREATE INDEX IF NOT EXISTS idx_credit_user ON credit_history(user_did, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_credit_task ON credit_history(task_id);

-- ============================================
-- Transaction History Table (Optional)
-- ============================================
CREATE TABLE IF NOT EXISTS xzt_transactions (
    tx_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Blockchain info
    tx_hash VARCHAR(66) UNIQUE,
    block_number BIGINT,
    
    -- Transaction details
    from_address VARCHAR(66) NOT NULL,
    to_address VARCHAR(66) NOT NULL,
    amount DECIMAL(20, 8) NOT NULL,
    
    -- Type
    tx_type VARCHAR(30) NOT NULL CHECK (tx_type IN (
        'transfer',
        'task_lock',
        'milestone_payment',
        'task_refund'
    )),
    
    -- Related task
    task_id UUID REFERENCES tasks(task_id),
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN (
        'pending',
        'confirmed',
        'failed'
    )),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    confirmed_at TIMESTAMP
);

-- Indexes for transactions
CREATE INDEX IF NOT EXISTS idx_tx_hash ON xzt_transactions(tx_hash);
CREATE INDEX IF NOT EXISTS idx_tx_from ON xzt_transactions(from_address);
CREATE INDEX IF NOT EXISTS idx_tx_to ON xzt_transactions(to_address);
CREATE INDEX IF NOT EXISTS idx_tx_task ON xzt_transactions(task_id);
CREATE INDEX IF NOT EXISTS idx_tx_type ON xzt_transactions(tx_type);
CREATE INDEX IF NOT EXISTS idx_tx_created_at ON xzt_transactions(created_at DESC);

-- ============================================
-- Update Triggers
-- ============================================

-- Auto-update updated_at timestamp for tasks
CREATE OR REPLACE FUNCTION update_tasks_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_tasks_updated_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_tasks_updated_at();

-- Auto-update updated_at timestamp for task_bids
CREATE OR REPLACE FUNCTION update_bids_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_bids_updated_at
    BEFORE UPDATE ON task_bids
    FOR EACH ROW
    EXECUTE FUNCTION update_bids_updated_at();

-- ============================================
-- Views for Common Queries
-- ============================================

-- Active tasks with bid counts
CREATE OR REPLACE VIEW v_active_tasks AS
SELECT 
    t.*,
    u_creator.username as creator_username,
    u_executor.username as executor_username,
    COUNT(DISTINCT tb.bid_id) as bid_count,
    COUNT(DISTINCT ts.submission_id) as submission_count
FROM tasks t
LEFT JOIN users u_creator ON t.creator_did = u_creator.did
LEFT JOIN users u_executor ON t.executor_did = u_executor.did
LEFT JOIN task_bids tb ON t.task_id = tb.task_id AND tb.status = 'pending'
LEFT JOIN task_submissions ts ON t.task_id = ts.task_id
WHERE t.status NOT IN ('completed', 'cancelled')
GROUP BY t.task_id, u_creator.username, u_executor.username;

-- User statistics
CREATE OR REPLACE VIEW v_user_stats AS
SELECT 
    u.did,
    u.username,
    u.credit_score,
    u.tasks_completed,
    u.tasks_cancelled,
    u.xzt_balance,
    COUNT(DISTINCT t_created.task_id) as tasks_created,
    COUNT(DISTINCT t_executing.task_id) as tasks_executing,
    COUNT(DISTINCT tb.bid_id) as active_bids
FROM users u
LEFT JOIN tasks t_created ON u.did = t_created.creator_did
LEFT JOIN tasks t_executing ON u.did = t_executing.executor_did AND t_executing.status NOT IN ('completed', 'cancelled')
LEFT JOIN task_bids tb ON u.did = tb.bidder_did AND tb.status = 'pending'
GROUP BY u.did, u.username, u.credit_score, u.tasks_completed, u.tasks_cancelled, u.xzt_balance;

-- ============================================
-- Sample Data (Optional - for testing)
-- ============================================

-- Uncomment to insert sample data
/*
-- Update admin user with initial XZT balance
UPDATE users 
SET xzt_balance = 10000.00, credit_score = 5000
WHERE username = 'admin';
*/
