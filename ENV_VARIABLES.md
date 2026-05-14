# Environment Variables

This file documents all environment variables used in the project.

## Required Variables

### API Configuration
- `API_URL` - Base URL for API calls
- `API_KEY` - Authentication key for API
- `API_TIMEOUT` - Request timeout in milliseconds (default: 30000)

### Database
- `DATABASE_URL` - Database connection string
- `DATABASE_POOL_SIZE` - Connection pool size (default: 10)

### Authentication
- `JWT_SECRET` - Secret key for JWT signing
- `JWT_EXPIRY` - JWT token expiry time (default: 24h)
- `REFRESH_TOKEN_EXPIRY` - Refresh token expiry (default: 7d)

### Email Service
- `SMTP_HOST` - SMTP server host
- `SMTP_PORT` - SMTP server port
- `SMTP_USER` - SMTP authentication user
- `SMTP_PASSWORD` - SMTP authentication password
- `EMAIL_FROM` - Default sender email address

### Logging
- `LOG_LEVEL` - Logging level (error, warn, info, debug)
- `LOG_FILE` - Path to log file

### Environment
- `NODE_ENV` - Environment (development, staging, production)
- `DEBUG` - Debug mode (true/false)

## Optional Variables

### Application
- `APP_NAME` - Application name (default: task-manament)
- `APP_VERSION` - Application version
- `SUPPORT_EMAIL` - Support contact email

### Features
- `ENABLE_EMAIL_NOTIFICATIONS` - Enable email notifications (default: true)
- `ENABLE_ANALYTICS` - Enable analytics tracking (default: false)
- `ENABLE_CACHING` - Enable response caching (default: true)

## Example .env

```
# API
API_URL=http://localhost:3000
API_KEY=your-api-key-here
API_TIMEOUT=30000

# Database
DATABASE_URL=postgres://user:password@localhost:5432/task-manament
DATABASE_POOL_SIZE=10

# Authentication
JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=7d

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@example.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM=noreply@task-manament.com

# Logging
LOG_LEVEL=debug
LOG_FILE=logs/app.log

# Environment
NODE_ENV=development
DEBUG=true

# Application
APP_NAME=task-manament
SUPPORT_EMAIL=support@task-manament.com

# Features
ENABLE_EMAIL_NOTIFICATIONS=true
ENABLE_ANALYTICS=false
ENABLE_CACHING=true
```

## Development (.env.local)

For local development, create a `.env.local` file with test values. This file is gitignored.

## Security Notes

- **Never commit .env files** - Use .env.example as template
- **Never share secrets** - Each developer maintains their own .env
- **Rotate secrets regularly** - Update API keys and secrets periodically
- **Validate on startup** - Check required variables exist on app start
- **Use strong values** - Generate strong random values for secrets

## Accessing Variables

In your code:

```javascript
// Node.js
const apiUrl = process.env.API_URL;
const jwtSecret = process.env.JWT_SECRET;

// With validation
const apiUrl = process.env.API_URL;
if (!apiUrl) {
  throw new Error('Missing required environment variable: API_URL');
}
```

## CI/CD Environment Variables

For GitHub Actions and deployment, set environment variables in:
1. GitHub Secrets (Settings > Secrets and variables > Actions)
2. Deployment platform settings (Heroku, AWS, etc.)

Never paste secrets in workflow files.
