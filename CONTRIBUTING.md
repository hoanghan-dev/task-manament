# Contributing to task-manament

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

Be respectful and professional in all interactions. We're building a welcoming community.

## Getting Started

### Prerequisites
- Node.js (v16 or higher)
- npm or yarn
- Git

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/hoanghan-dev/task-manament.git
cd task-manament

# Install dependencies
npm install

# Install pre-commit hooks
npx husky install

# Verify setup
npm run lint
npm run test
```

## Development Workflow

### 1. Create a Feature Branch
```bash
git checkout -b feature/your-feature-name
```

Use one of these prefixes:
- `feature/` - New functionality
- `bugfix/` - Bug fixes
- `hotfix/` - Urgent production fixes
- `refactor/` - Code improvements (no functional change)
- `docs/` - Documentation updates
- `test/` - Test improvements
- `chore/` - Maintenance tasks

### 2. Make Your Changes
- Keep commits small and focused
- Follow the coding style (see below)
- Write meaningful commit messages
- Add or update tests as needed

### 3. Run Quality Checks

```bash
# Format code
npm run format

# Lint code
npm run lint

# Run tests
npm run test

# Check test coverage
npm run test:coverage
```

### 4. Commit Your Changes

Follow conventional commit format:
```
type(scope): short description

Longer description explaining the change, why it was made,
and any relevant context.

Closes #123
```

**Types**: feat, fix, docs, style, refactor, perf, test, chore

**Example**:
```
feat(task): add ability to delete tasks

- Implement delete endpoint
- Add confirmation dialog
- Update task service
- Add unit tests

Closes #42
```

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then open a PR on GitHub. Include:
- Description of changes
- Related issue numbers
- Screenshots (if UI changes)
- Testing notes

## Coding Standards

### General Principles
- **Keep it simple**: Avoid unnecessary complexity
- **Be explicit**: Write clear, self-documenting code
- **DRY**: Don't repeat yourself
- **SOLID**: Follow SOLID principles
- **Clean Code**: Aim for readable, maintainable code

### Naming Conventions

```javascript
// Constants: SCREAMING_SNAKE_CASE
const MAX_RETRIES = 3;
const API_ENDPOINT = 'https://api.example.com';

// Functions & methods: camelCase
function getUserById(userId) {}
const createTask = () => {};

// Classes: PascalCase
class TaskManager {}
class UserService {}

// Variables: camelCase
let currentUser = null;
const taskList = [];

// Booleans: is/has prefix
const isCompleted = true;
const hasPermission = false;
```

### File Naming
- Components: PascalCase (e.g., `TaskList.jsx`)
- Utilities: camelCase (e.g., `dateUtils.js`)
- Tests: Same as source + `.test.js` (e.g., `taskService.test.js`)
- Styles: kebab-case (e.g., `task-list.css`)

### Code Style

#### Length Limits
- Max line length: 100 characters
- Max function length: 25 lines (excluding comments)
- Max file size: 300 lines

#### Comments
- Explain **why**, not what
- Use single-line comments: `//`
- Keep comments up-to-date with code
- Avoid obvious comments

```javascript
// Good: Explains the reasoning
// We sort by creation date descending to show newest first
const sortedTasks = tasks.sort((a, b) => b.createdAt - a.createdAt);

// Bad: Restates the code
// Sort tasks
const sortedTasks = tasks.sort((a, b) => b.createdAt - a.createdAt);
```

#### Formatting
- Use 2-space indentation
- No trailing spaces
- Final newline in all files
- Automatic formatting via Prettier

## Testing

### Test Requirements
- Write tests for new features
- Update tests when changing code
- Aim for 80%+ code coverage
- Focus on meaningful tests, not just coverage

### Test Structure
```javascript
describe('TaskService', () => {
  describe('createTask', () => {
    it('should create a new task with valid input', () => {
      // Arrange
      const taskData = { title: 'New Task', description: 'Test' };
      
      // Act
      const result = taskService.createTask(taskData);
      
      // Assert
      expect(result).toHaveProperty('id');
      expect(result.title).toBe('New Task');
    });

    it('should throw error when title is empty', () => {
      expect(() => {
        taskService.createTask({ title: '' });
      }).toThrow(ValidationError);
    });
  });
});
```

## Pull Request Process

1. **Before Submitting**
   - [ ] Code is formatted (`npm run format`)
   - [ ] No linting errors (`npm run lint`)
   - [ ] All tests pass (`npm run test`)
   - [ ] Test coverage is adequate
   - [ ] No console.log/debugger statements
   - [ ] Commit messages are clear
   - [ ] Documentation updated

2. **PR Description Should Include**
   - What changes were made
   - Why the changes were made
   - How to test the changes
   - Any breaking changes
   - Related issue numbers

3. **Code Review**
   - At least one approval required
   - Address review comments
   - Keep discussions professional

4. **Merging**
   - Squash commits if requested
   - Delete feature branch after merge
   - Verify deployment in staging

## Common Issues & Solutions

### Pre-commit hook fails
```bash
# Check what's wrong
npm run lint
npm run format

# Fix issues
npm run format -- --write .

# If still failing, fix manually then try again
git add .
git commit -m "fix: resolve linting issues"
```

### Tests failing
```bash
# Run tests with verbose output
npm run test -- --verbose

# Run specific test file
npm run test -- taskService.test.js

# Run tests in watch mode
npm run test -- --watch
```

### Merge conflicts
```bash
# Update your branch
git fetch origin
git rebase origin/main

# Resolve conflicts in your editor, then
git add .
git rebase --continue
```

## Questions?

- Check existing issues and documentation
- Open a new issue with your question
- Contact the maintainers

## License

By contributing, you agree that your contributions will be licensed under the project's license.

---

Thank you for helping make task-manament better! 🚀
