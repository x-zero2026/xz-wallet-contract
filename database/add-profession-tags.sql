-- Add profession_tags field to tasks table
-- Date: 2026-02-08

-- Step 1: Add profession_tags column
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS profession_tags TEXT[] DEFAULT '{}';

-- Step 2: Create GIN index for profession_tags (for efficient array queries)
CREATE INDEX IF NOT EXISTS idx_tasks_profession_tags ON tasks USING GIN(profession_tags);

-- Step 3: Add comment
COMMENT ON COLUMN tasks.profession_tags IS '任务所需的职业标签数组，最多5个标签，支持预定义标签和自定义标签';

-- Step 4: Verify the changes
DO $$
DECLARE
    tags_exists BOOLEAN;
BEGIN
    SELECT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'tasks' AND column_name = 'profession_tags'
    ) INTO tags_exists;
    
    IF tags_exists THEN
        RAISE NOTICE 'Migration successful! profession_tags column added to tasks table.';
    ELSE
        RAISE EXCEPTION 'Migration failed! Column not found.';
    END IF;
END $$;

-- Migration complete
SELECT 'Migration completed successfully. Tasks table now has profession_tags column.' AS status;
