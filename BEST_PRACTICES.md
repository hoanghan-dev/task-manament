# Best Practices Guide

This document outlines the best practices for developing task-manament to maintain clean, scalable code.

## Table of Contents
1. [General Principles](#general-principles)
2. [JavaScript/TypeScript Best Practices](#javascripttypescript-best-practices)
3. [Component Development](#component-development)
4. [Service Layer](#service-layer)
5. [Testing Best Practices](#testing-best-practices)
6. [Documentation](#documentation)
7. [Performance](#performance)
8. [Security](#security)

## General Principles

### 1. Keep It Simple (KISS)
- Write code that is easy to understand
- Avoid unnecessary complexity
- Choose clarity over cleverness

```javascript
// Bad: Too clever
const result = arr.reduce((a, b, i) => i % 2 ? a : [...a, [b]], []);

// Good: Clear and readable
const pairs = array.reduce((result, item, index) => {
  if (index % 2 === 0) {
    result.push([item]);
  } else {
    result[result.length - 1].push(item);
  }
  return result;
}, []);
```

### 2. DRY - Don't Repeat Yourself
- Extract reusable logic into functions
- Use shared utilities and helpers
- Create reusable components

```javascript
// Bad: Repeated validation
if (name && name.length > 0) { }
if (email && email.length > 0) { }
if (phone && phone.length > 0) { }

// Good: Create reusable validator
const isRequired = (value) => Boolean(value?.length > 0);

if (isRequired(name)) { }
if (isRequired(email)) { }
if (isRequired(phone)) { }
```

### 3. SOLID Principles

#### Single Responsibility
```javascript
// Bad: Multiple responsibilities
class UserManager {
  createUser(userData) { }
  sendEmail(email) { }
  validateEmail(email) { }
  updateDatabase(user) { }
}

// Good: Single responsibility
class UserService {
  createUser(userData) { }
  updateUser(userId, updates) { }
}

class EmailService {
  sendEmail(email) { }
  validateEmail(email) { }
}

class UserRepository {
  save(user) { }
  update(userId, updates) { }
}
```

#### Open/Closed Principle
```javascript
// Bad: Hard to extend
function formatOutput(data, format) {
  if (format === 'json') {
    return JSON.stringify(data);
  } else if (format === 'csv') {
    return convertToCSV(data);
  }
  // Must modify function for each new format
}

// Good: Open for extension
class Formatter {
  format(data) { throw new Error('Must implement'); }
}

class JSONFormatter extends Formatter {
  format(data) { return JSON.stringify(data); }
}

class CSVFormatter extends Formatter {
  format(data) { return convertToCSV(data); }
}
```

## JavaScript/TypeScript Best Practices

### 1. Variable Naming
```javascript
// Bad
let d = new Date();
let u = getUser();
let x = 5;

// Good
let currentDate = new Date();
let currentUser = getUser();
let maxRetries = 5;

// For boolean variables, use is/has prefix
let isActive = true;
let hasPermission = false;
let isLoading = false;
```

### 2. Function Length and Complexity
```javascript
// Bad: Too long, many responsibilities
function processTask(task) {
  // Validate
  if (!task.title || task.title.length === 0) throw new Error('Title required');
  if (!task.description || task.description.length === 0) throw new Error('Description required');
  
  // Transform
  const normalized = task.title.trim().toLowerCase();
  
  // Save
  const result = saveToDatabase(normalized);
  
  // Send notification
  sendEmail(task.assignee, 'New task created');
  
  // Log
  logger.info('Task created');
  
  return result;
}

// Good: Focused, single responsibility
function processTask(task) {
  validateTask(task);
  const normalizedTask = normalizeTask(task);
  const savedTask = saveTask(normalizedTask);
  notifyAssignee(task);
  return savedTask;
}

function validateTask(task) {
  if (!task.title?.length) throw new ValidationError('Title required');
  if (!task.description?.length) throw new ValidationError('Description required');
}

function normalizeTask(task) {
  return {
    ...task,
    title: task.title.trim().toLowerCase(),
  };
}
```

### 3. Error Handling
```javascript
// Bad: Silent failures
try {
  await saveTask(task);
} catch (e) {
  // Ignore
}

// Good: Proper error handling
try {
  await saveTask(task);
} catch (error) {
  logger.error('Failed to save task', { taskId: task.id, error });
  throw new TaskSaveError('Could not save task', { cause: error });
}
```

### 4. Async/Await
```javascript
// Bad: Callback hell
function loadTask(id) {
  getTask(id, (task) => {
    getRelatedTasks(task.id, (related) => {
      getUser(task.userId, (user) => {
        renderTask(task, related, user);
      });
    });
  });
}

// Good: Clean async/await
async function loadTask(id) {
  try {
    const task = await getTask(id);
    const [related, user] = await Promise.all([
      getRelatedTasks(task.id),
      getUser(task.userId),
    ]);
    renderTask(task, related, user);
  } catch (error) {
    handleError(error);
  }
}
```

### 5. Array/Object Methods
```javascript
// Good: Use array methods
const completedTasks = tasks.filter(task => task.isComplete);
const taskTitles = tasks.map(task => task.title);
const hasOverdueTasks = tasks.some(task => task.dueDate < today);
const allApproved = tasks.every(task => task.isApproved);

// Avoid imperative loops when declarative methods exist
// Bad
const results = [];
for (let i = 0; i < tasks.length; i++) {
  if (tasks[i].priority === 'high') {
    results.push(tasks[i]);
  }
}

// Good
const highPriorityTasks = tasks.filter(task => task.priority === 'high');
```

## Component Development

### 1. Component Structure
```javascript
// Good structure
export function TaskCard({ task, onEdit, onDelete }) {
  const [isExpanded, setIsExpanded] = useState(false);

  const handleDelete = useCallback(() => {
    if (confirm('Delete this task?')) {
      onDelete(task.id);
    }
  }, [task.id, onDelete]);

  return (
    <article className="task-card">
      <h3>{task.title}</h3>
      {isExpanded && <p>{task.description}</p>}
      <button onClick={() => setIsExpanded(!isExpanded)}>
        {isExpanded ? 'Collapse' : 'Expand'}
      </button>
      <button onClick={() => onEdit(task)}>Edit</button>
      <button onClick={handleDelete}>Delete</button>
    </article>
  );
}
```

### 2. Props Validation
```javascript
// TypeScript
interface TaskCardProps {
  task: Task;
  onEdit: (task: Task) => void;
  onDelete: (taskId: string) => void;
}

export function TaskCard({ task, onEdit, onDelete }: TaskCardProps) {
  // ...
}

// Or PropTypes (non-TypeScript)
import PropTypes from 'prop-types';

TaskCard.propTypes = {
  task: PropTypes.shape({
    id: PropTypes.string.isRequired,
    title: PropTypes.string.isRequired,
  }).isRequired,
  onEdit: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
};
```

### 3. Avoid Prop Drilling
```javascript
// Bad: Prop drilling
function Page(props) {
  return <Component1 user={props.user} />;
}

function Component1({ user }) {
  return <Component2 user={user} />;
}

function Component2({ user }) {
  return <div>{user.name}</div>;
}

// Good: Use context
const UserContext = React.createContext();

function Page(props) {
  return (
    <UserContext.Provider value={props.user}>
      <Component1 />
    </UserContext.Provider>
  );
}

function Component2() {
  const user = useContext(UserContext);
  return <div>{user.name}</div>;
}
```

## Service Layer

### 1. Service Structure
```javascript
// Good: Clean service interface
export const taskService = {
  async getAll() {
    const response = await api.get('/tasks');
    return response.data;
  },

  async getById(id) {
    const response = await api.get(`/tasks/${id}`);
    return response.data;
  },

  async create(taskData) {
    const validated = validateTask(taskData);
    const response = await api.post('/tasks', validated);
    return response.data;
  },

  async update(id, updates) {
    const validated = validateTask(updates);
    const response = await api.put(`/tasks/${id}`, validated);
    return response.data;
  },

  async delete(id) {
    await api.delete(`/tasks/${id}`);
  },
};
```

### 2. Error Handling in Services
```javascript
// Good: Custom error classes
class TaskServiceError extends Error {
  constructor(message, { cause, taskId } = {}) {
    super(message);
    this.name = 'TaskServiceError';
    this.cause = cause;
    this.taskId = taskId;
  }
}

export const taskService = {
  async create(taskData) {
    try {
      const response = await api.post('/tasks', taskData);
      return response.data;
    } catch (error) {
      throw new TaskServiceError('Failed to create task', { cause: error });
    }
  },
};
```

## Testing Best Practices

### 1. Test Structure (AAA Pattern)
```javascript
describe('TaskService', () => {
  describe('create', () => {
    it('should create a task with valid data', async () => {
      // Arrange
      const taskData = { title: 'New task', priority: 'high' };

      // Act
      const result = await taskService.create(taskData);

      // Assert
      expect(result).toHaveProperty('id');
      expect(result.title).toBe('New task');
      expect(result.createdAt).toBeDefined();
    });

    it('should throw error when title is empty', async () => {
      // Arrange
      const taskData = { title: '', priority: 'high' };

      // Act & Assert
      await expect(taskService.create(taskData))
        .rejects
        .toThrow('Title is required');
    });
  });
});
```

### 2. Mocking
```javascript
// Good: Mock API calls
jest.mock('../api');

describe('taskService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should call API with correct parameters', async () => {
    // Arrange
    const taskData = { title: 'Test' };
    api.post.mockResolvedValue({ data: { id: '1', ...taskData } });

    // Act
    const result = await taskService.create(taskData);

    // Assert
    expect(api.post).toHaveBeenCalledWith('/tasks', taskData);
    expect(result.id).toBe('1');
  });
});
```

### 3. Edge Cases
```javascript
// Always test edge cases
describe('formatDate', () => {
  it('should handle null input', () => {
    expect(formatDate(null)).toBe('Invalid date');
  });

  it('should handle undefined input', () => {
    expect(formatDate(undefined)).toBe('Invalid date');
  });

  it('should handle invalid date', () => {
    expect(formatDate(new Date('invalid'))).toBe('Invalid date');
  });

  it('should format valid date correctly', () => {
    const date = new Date('2026-05-14');
    expect(formatDate(date)).toMatch(/May 14, 2026/);
  });
});
```

## Documentation

### 1. JSDoc Comments
```javascript
/**
 * Validates a task object
 * @param {Object} task - The task to validate
 * @param {string} task.title - Task title (required)
 * @param {string} [task.description] - Task description (optional)
 * @param {string} [task.priority='medium'] - Task priority
 * @throws {ValidationError} If validation fails
 * @returns {Object} Validated task object
 */
export function validateTask(task) {
  // implementation
}
```

### 2. README Documentation
Each module should have clear documentation about:
- Purpose
- Usage examples
- API reference
- Common use cases

## Performance

### 1. Memoization
```javascript
// Prevent unnecessary re-renders
const TaskList = React.memo(function TaskList({ tasks }) {
  return (
    <ul>
      {tasks.map(task => (
        <TaskItem key={task.id} task={task} />
      ))}
    </ul>
  );
});

// Memoize callbacks
const handleEdit = useCallback((task) => {
  updateTask(task);
}, [updateTask]);
```

### 2. Lazy Loading
```javascript
// Load components only when needed
const TaskEditor = React.lazy(() => import('./TaskEditor'));

function TaskPage() {
  return (
    <Suspense fallback={<Spinner />}>
      <TaskEditor />
    </Suspense>
  );
}
```

## Security

### 1. Input Validation
```javascript
// Always validate user input
function validateTaskTitle(title) {
  if (!title || typeof title !== 'string') {
    throw new ValidationError('Invalid title');
  }
  if (title.length < 1 || title.length > 255) {
    throw new ValidationError('Title must be 1-255 characters');
  }
  return title.trim();
}
```

### 2. XSS Prevention
```javascript
// Bad: Vulnerable to XSS
<div dangerHTML={userInput} />

// Good: Safe
<div>{userInput}</div>

// If you must use HTML, sanitize it
import DOMPurify from 'dompurify';
<div dangerouslySetInnerHTML={{
  __html: DOMPurify.sanitize(userInput)
}} />
```

### 3. Secrets Management
```javascript
// Bad: Hardcoded secrets
const API_KEY = 'sk_live_abc123';

// Good: Environment variables
const API_KEY = process.env.API_KEY;

// Validate required secrets exist
if (!process.env.API_KEY) {
  throw new Error('Missing required environment variable: API_KEY');
}
```

---

## Summary

Following these best practices will help ensure:
- ✅ Clean, readable code
- ✅ Fewer bugs and issues
- ✅ Easier maintenance
- ✅ Better performance
- ✅ Enhanced security
- ✅ Improved team collaboration

**Remember**: Code is read much more often than it's written. Write for the next person who reads it.
