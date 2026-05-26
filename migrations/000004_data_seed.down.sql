-- Xóa tasks trước vì tasks phụ thuộc workspaces và users
DELETE FROM tasks
WHERE title LIKE 'Task % - Workspace %';

-- Xóa workspaces
DELETE FROM workspaces
WHERE name LIKE 'Workspace % of User %';

-- Xóa users
DELETE FROM users
WHERE email LIKE 'user%@gmail.com';