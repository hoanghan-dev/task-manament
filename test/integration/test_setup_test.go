package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	authHandler "dev/task-management/internal/modules/auth/handler"
	userRepo "dev/task-management/internal/modules/auth/repositories"
	authService "dev/task-management/internal/modules/auth/services"
	commentHandler "dev/task-management/internal/modules/comment/handler"
	commentRepo "dev/task-management/internal/modules/comment/repositories"
	commentService "dev/task-management/internal/modules/comment/services"
	taskHandler "dev/task-management/internal/modules/task/handler"
	taskRepo "dev/task-management/internal/modules/task/repositories"
	taskService "dev/task-management/internal/modules/task/services"
	workspaceHandler "dev/task-management/internal/modules/workspace/handler"
	workspaceRepo "dev/task-management/internal/modules/workspace/repositories"
	workspaceService "dev/task-management/internal/modules/workspace/services"
	"dev/task-management/internal/realtime"
	"dev/task-management/internal/router"
	"dev/task-management/pkg/cache"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// testDB holds the test database connection, shared across all tests in this package.
var testDB *sql.DB

// testRedis holds the test Redis client.
var testRedis *redis.Client

// testRouter holds the fully initialized gin.Engine for HTTP testing.
var testRouter *gin.Engine

// TestMain is the entry point for all tests in this package.
// It sets up the test environment (DB, Redis, Router) and tears it down after tests complete.
func TestMain(m *testing.M) {
	// Load .env.test if present (for local development).
	// Try multiple paths since CWD can vary depending on how tests are invoked.
	_ = godotenv.Load(".env.test")                  // if CWD is test/integration/
	_ = godotenv.Load("test/integration/.env.test") // if CWD is project root

	// Setup
	var err error

	testDB, err = setupTestDB()
	if err != nil {
		log.Fatalf("[INTEGRATION TEST] Failed to connect to test database: %v", err)
	}

	testRedis, err = setupTestRedis()
	if err != nil {
		log.Fatalf("[INTEGRATION TEST] Failed to connect to test Redis: %v", err)
	}

	// Run migrations
	if err := runMigrations(testDB); err != nil {
		log.Fatalf("[INTEGRATION TEST] Failed to run migrations: %v", err)
	}

	// Setup router with real dependencies
	testRouter = setupRouter(testDB, testRedis)

	// Run tests
	code := m.Run()

	// Teardown
	cleanupDatabase(testDB)
	testDB.Close()
	testRedis.Close()

	os.Exit(code)
}

// setupTestDB creates a connection to the test PostgreSQL database.
// It reads from TEST_DB_* environment variables (or falls back to DB_* variables).
func setupTestDB() (*sql.DB, error) {
	host := getEnvWithFallback("TEST_DB_HOST", "DB_HOST", "localhost")
	port := getEnvWithFallback("TEST_DB_PORT", "DB_PORT", "5432")
	user := getEnvWithFallback("TEST_DB_USER", "DB_USER", "postgres")
	password := getEnvWithFallback("TEST_DB_PASSWORD", "DB_PASSWORD", "postgres")
	dbName := getEnvWithFallback("TEST_DB_NAME", "DB_NAME", "task_management_test")
	sslMode := getEnvWithFallback("TEST_DB_SSLMODE", "DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping failed (is the test database '%s' created?): %w", dbName, err)
	}

	log.Printf("[INTEGRATION TEST] Connected to test database: %s@%s:%s/%s", user, host, port, dbName)
	return db, nil
}

// setupTestRedis creates a connection to the test Redis instance.
func setupTestRedis() (*redis.Client, error) {
	addr := getEnvWithFallback("REDIS_ADDR", "", "localhost:6379")
	password := getEnvWithFallback("REDIS_PASSWORD", "", "")
	dbStr := getEnvWithFallback("REDIS_DB", "", "1") // Use DB=1 for test isolation

	db := 1
	if _, err := fmt.Sscanf(dbStr, "%d", &db); err != nil {
		db = 1
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis.Ping failed: %w", err)
	}

	log.Printf("[INTEGRATION TEST] Connected to test Redis: %s (DB=%d)", addr, db)
	return rdb, nil
}

// setupRouter builds the gin.Engine with full production dependency injection,
// using the test database and test Redis connections.
// This mirrors the production NewApp() in internal/app/app.go.
func setupRouter(db *sql.DB, redisClient *redis.Client) *gin.Engine {
	gin.SetMode(gin.TestMode)

	redisService := cache.NewRedisCacheService(redisClient)

	hubConnectionWS := realtime.NewHubConnectionWS()
	subscriber := realtime.NewNotificationSubscriber(hubConnectionWS, redisService)
	wsHandler := realtime.NewWSHandler(hubConnectionWS, subscriber)

	wpRepo := workspaceRepo.NewWorkspaceRepository(db)
	wpService := workspaceService.NewWorkspaceService(wpRepo)
	wpHandler := workspaceHandler.NewWorkspaceHandler(wpService)

	tRepo := taskRepo.NewTaskRepository(db)
	tService := taskService.NewTaskService(tRepo, wpService, redisService)
	tHandler := taskHandler.NewTaskHandler(tService)

	uRepo := userRepo.NewUserRepository(db)
	aService := authService.NewAuthService(uRepo, wpService, redisService)
	aHandler := authHandler.NewAuthHandler(aService)

	cRepo := commentRepo.NewCommentRepository(db)
	cService := commentService.NewCommentService(cRepo, redisService)
	cHandler := commentHandler.NewCommentHandler(cService)

	r := router.SetupRouter(router.RouterDependencies{
		TaskHandler:      tHandler,
		AuthHander:       aHandler,
		WorkspaceHandler: wpHandler,
		CommentHandler:   cHandler,
		WSHandler:        wsHandler,
	})

	return r
}

// runMigrations executes DDL statements to create all required tables.
// This uses inline SQL matching the project's migration files
// so the test is self-contained (no dependency on migrate CLI).
func runMigrations(db *sql.DB) error {
	migrations := []string{
		// 000001 - users
		`CREATE TABLE IF NOT EXISTS users (
			user_id UUID PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			full_name TEXT NOT NULL,
			create_at TIMESTAMP NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,

		// 000002 - workspaces
		`CREATE TABLE IF NOT EXISTS workspaces (
			workspace_id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			owner_id UUID NOT NULL,
			create_at TIMESTAMP NOT NULL,
			CONSTRAINT fk_workspaces_users
				FOREIGN KEY (owner_id) REFERENCES users(user_id) ON DELETE CASCADE
		)`,

		// 000003 - tasks
		`CREATE TABLE IF NOT EXISTS tasks (
			task_id UUID PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'TODO',
			assignee_id UUID NULL,
			workspace_id UUID NOT NULL,
			create_at TIMESTAMP NOT NULL,
			CONSTRAINT task_status_contraint
				CHECK (status IN ('TODO', 'IN_PROGRESS', 'DONE', 'BLOCKED')),
			CONSTRAINT fk_tasks_workspaces
				FOREIGN KEY (workspace_id) REFERENCES workspaces(workspace_id) ON DELETE CASCADE,
			CONSTRAINT fk_tasks_users
				FOREIGN KEY (assignee_id) REFERENCES users(user_id) ON DELETE CASCADE
		)`,

		// 000005 - notifications
		`CREATE TABLE IF NOT EXISTS notifications (
			notification_id UUID PRIMARY KEY,
			sender_id UUID NOT NULL,
			receiver_id UUID NOT NULL,
			task_id UUID NOT NULL,
			message TEXT NOT NULL,
			create_at TIMESTAMP NOT NULL DEFAULT NOW(),
			CONSTRAINT fk_notifications_sender
				FOREIGN KEY (sender_id) REFERENCES users(user_id) ON DELETE CASCADE,
			CONSTRAINT fk_notifications_receiver
				FOREIGN KEY (receiver_id) REFERENCES users(user_id) ON DELETE CASCADE,
			CONSTRAINT fk_notifications_task
				FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
		)`,

		// 000006 - comments
		`CREATE TABLE IF NOT EXISTS comments (
			comment_id UUID PRIMARY KEY,
			task_id UUID NOT NULL,
			user_id UUID NOT NULL,
			content TEXT NOT NULL,
			create_at TIMESTAMP NOT NULL DEFAULT NOW(),
			update_at TIMESTAMP NULL,
			CONSTRAINT fk_comments_tasks
				FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE,
			CONSTRAINT fk_comments_users
				FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_task_id ON comments(task_id)`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration #%d failed: %w", i+1, err)
		}
	}

	log.Println("[INTEGRATION TEST] Migrations completed successfully")
	return nil
}

// cleanupDatabase truncates all test tables with RESTART IDENTITY CASCADE.
// Call this before each test to ensure a clean state.
func cleanupDatabase(db *sql.DB) {
	tables := []string{"comments", "notifications", "tasks", "workspaces", "users"}
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		if _, err := db.Exec(query); err != nil {
			log.Printf("[INTEGRATION TEST] Warning: failed to truncate %s: %v", table, err)
		}
	}
}

// getEnvWithFallback tries primary env var, then fallback env var, then default value.
func getEnvWithFallback(primary, fallback, defaultVal string) string {
	if val := os.Getenv(primary); val != "" {
		return val
	}
	if fallback != "" {
		if val := os.Getenv(fallback); val != "" {
			return val
		}
	}
	return defaultVal
}
