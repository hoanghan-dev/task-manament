-- cần extension để dùng gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 1. Insert 10 users
INSERT INTO users (
    user_id,
    email,
    password_hash,
    full_name,
    create_at
)
SELECT
    gen_random_uuid(),
    'user' || i || '@gmail.com',
    '$2a$12$YyB4H6tHrB17V148sIEY4uoAXIRWJnZysmARHem4ba/.LhZP2D/KG',
    'User ' || i,
    NOW()
FROM generate_series(1, 10) AS s(i);

-- 2. Mỗi user có 3 workspaces
INSERT INTO workspaces (
    workspace_id,
    name,
    description,
    owner_id,
    create_at
)
SELECT
    gen_random_uuid(),
    'Workspace ' || w || ' of ' || u.full_name,
    'Description for workspace ' || w || ' of ' || u.full_name,
    u.user_id,
    NOW()
FROM users u
CROSS JOIN generate_series(1, 3) AS s(w)
WHERE u.email LIKE 'user%@gmail.com';

-- 3. Mỗi workspace có 10 tasks
INSERT INTO tasks (
    task_id,
    title,
    description,
    status,
    assignee_id,
    workspace_id,
    create_at
)
SELECT
    gen_random_uuid(),
    'Task ' || t || ' - ' || w.name,
    'Description for task ' || t || ' in ' || w.name,
    CASE
        WHEN t % 4 = 1 THEN 'TODO'
        WHEN t % 4 = 2 THEN 'IN_PROGRESS'
        WHEN t % 4 = 3 THEN 'DONE'
        ELSE 'BLOCKED'
    END,
    w.owner_id,
    w.workspace_id,
    NOW()
FROM workspaces w
CROSS JOIN generate_series(1, 10) AS s(t)
WHERE w.name LIKE 'Workspace%';

