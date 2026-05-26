-- 1. Insert 10 users
INSERT INTO users (
    user_id,
    email,
    password_hash,
    full_name,
    create_at
)
VALUES
('11111111-1111-1111-1111-111111111111', 'user1@gmail.com', '$2a$12$YyB4H6tHrB17V148sIEY4uoAXIRWJnZysmARHem4ba/.LhZP2D/KG', 'User 1', NOW()),
('22222222-2222-2222-2222-222222222222', 'user2@gmail.com', '$2a$12$YyB4H6tHrB17V148sIEY4uoAXIRWJnZysmARHem4ba/.LhZP2D/KG', 'User 2', NOW()),
('33333333-3333-3333-3333-333333333333', 'user3@gmail.com', '$2a$12$YyB4H6tHrB17V148sIEY4uoAXIRWJnZysmARHem4ba/.LhZP2D/KG', 'User 3', NOW()),
('44444444-4444-4444-4444-444444444444', 'user4@gmail.com', '$2a$12$YyB4H6tHrB17V148sIEY4uoAXIRWJnZysmARHem4ba/.LhZP2D/KG', 'User 4', NOW()),
('55555555-5555-5555-5555-555555555555', 'user5@gmail.com', '$2a$12$YyB4H6tHrB17V148sIEY4uoAXIRWJnZysmARHem4ba/.LhZP2D/KG', 'User 5', NOW()),
('66666666-6666-6666-6666-666666666666', 'user6@gmail.com', '$2a$12$YyB4H6tHrB17V148sIEY4uoAXIRWJnZysmARHem4ba/.LhZP2D/KG', 'User 6', NOW()),
('77777777-7777-7777-7777-777777777777', 'user7@gmail.com', '$2a$12$YyB4H6tHrB17V148sIEY4uoAXIRWJnZysmARHem4ba/.LhZP2D/KG', 'User 7', NOW()),
('88888888-8888-8888-8888-888888888888', 'user8@gmail.com', '$2a$12$YyB4H6tHrB17V148sIEY4uoAXIRWJnZysmARHem4ba/.LhZP2D/KG', 'User 8', NOW()),
('99999999-9999-9999-9999-999999999999', 'user9@gmail.com', '$2a$12$YyB4H6tHrB17V148sIEY4uoAXIRWJnZysmARHem4ba/.LhZP2D/KG', 'User 9', NOW()),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'user10@gmail.com', '$2a$12$YyB4H6tHrB17V148sIEY4uoAXIRWJnZysmARHem4ba/.LhZP2D/KG', 'User 10', NOW());

-- 2. Mỗi user có đúng 1 workspace
INSERT INTO workspaces (
    workspace_id,
    name,
    description,
    owner_id,
    create_at
)
VALUES
('aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Workspace of User 1', 'Description for workspace of User 1', '11111111-1111-1111-1111-111111111111', NOW()),
('aaaaaaa2-aaaa-aaaa-aaaa-aaaaaaaaaaa2', 'Workspace of User 2', 'Description for workspace of User 2', '22222222-2222-2222-2222-222222222222', NOW()),
('aaaaaaa3-aaaa-aaaa-aaaa-aaaaaaaaaaa3', 'Workspace of User 3', 'Description for workspace of User 3', '33333333-3333-3333-3333-333333333333', NOW()),
('aaaaaaa4-aaaa-aaaa-aaaa-aaaaaaaaaaa4', 'Workspace of User 4', 'Description for workspace of User 4', '44444444-4444-4444-4444-444444444444', NOW()),
('aaaaaaa5-aaaa-aaaa-aaaa-aaaaaaaaaaa5', 'Workspace of User 5', 'Description for workspace of User 5', '55555555-5555-5555-5555-555555555555', NOW()),
('aaaaaaa6-aaaa-aaaa-aaaa-aaaaaaaaaaa6', 'Workspace of User 6', 'Description for workspace of User 6', '66666666-6666-6666-6666-666666666666', NOW()),
('aaaaaaa7-aaaa-aaaa-aaaa-aaaaaaaaaaa7', 'Workspace of User 7', 'Description for workspace of User 7', '77777777-7777-7777-7777-777777777777', NOW()),
('aaaaaaa8-aaaa-aaaa-aaaa-aaaaaaaaaaa8', 'Workspace of User 8', 'Description for workspace of User 8', '88888888-8888-8888-8888-888888888888', NOW()),
('aaaaaaa9-aaaa-aaaa-aaaa-aaaaaaaaaaa9', 'Workspace of User 9', 'Description for workspace of User 9', '99999999-9999-9999-9999-999999999999', NOW()),
('aaaaaa10-aaaa-aaaa-aaaa-aaaaaaaaaa10', 'Workspace of User 10', 'Description for workspace of User 10', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NOW());

-- 3. Insert tasks
INSERT INTO tasks (
    task_id,
    title,
    description,
    status,
    assignee_id,
    workspace_id,
    create_at
)
VALUES
('bbbbbbb1-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 'Task 1 - Workspace of User 1', 'Description for task 1 in Workspace of User 1', 'TODO', '11111111-1111-1111-1111-111111111111', 'aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1', NOW()),
('bbbbbbb2-bbbb-bbbb-bbbb-bbbbbbbbbbb2', 'Task 2 - Workspace of User 1', 'Description for task 2 in Workspace of User 1', 'IN_PROGRESS', '11111111-1111-1111-1111-111111111111', 'aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1', NOW()),
('bbbbbbb3-bbbb-bbbb-bbbb-bbbbbbbbbbb3', 'Task 3 - Workspace of User 1', 'Description for task 3 in Workspace of User 1', 'DONE', '11111111-1111-1111-1111-111111111111', 'aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1', NOW()),
('bbbbbbb4-bbbb-bbbb-bbbb-bbbbbbbbbbb4', 'Task 4 - Workspace of User 1', 'Description for task 4 in Workspace of User 1', 'BLOCKED', '11111111-1111-1111-1111-111111111111', 'aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaa1', NOW()),

('bbbbbbb5-bbbb-bbbb-bbbb-bbbbbbbbbbb5', 'Task 1 - Workspace of User 2', 'Description for task 1 in Workspace of User 2', 'TODO', '22222222-2222-2222-2222-222222222222', 'aaaaaaa2-aaaa-aaaa-aaaa-aaaaaaaaaaa2', NOW()),
('bbbbbbb6-bbbb-bbbb-bbbb-bbbbbbbbbbb6', 'Task 2 - Workspace of User 2', 'Description for task 2 in Workspace of User 2', 'IN_PROGRESS', '22222222-2222-2222-2222-222222222222', 'aaaaaaa2-aaaa-aaaa-aaaa-aaaaaaaaaaa2', NOW()),
('bbbbbbb7-bbbb-bbbb-bbbb-bbbbbbbbbbb7', 'Task 3 - Workspace of User 2', 'Description for task 3 in Workspace of User 2', 'DONE', '22222222-2222-2222-2222-222222222222', 'aaaaaaa2-aaaa-aaaa-aaaa-aaaaaaaaaaa2', NOW()),
('bbbbbbb8-bbbb-bbbb-bbbb-bbbbbbbbbbb8', 'Task 4 - Workspace of User 2', 'Description for task 4 in Workspace of User 2', 'BLOCKED', '22222222-2222-2222-2222-222222222222', 'aaaaaaa2-aaaa-aaaa-aaaa-aaaaaaaaaaa2', NOW()); 