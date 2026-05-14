# Project Structure

This document describes the recommended project structure for task-manament.

## Directory Organization

```
task-manament/
├── .github/
│   ├── workflows/              # GitHub Actions CI/CD workflows
│   └── ISSUE_TEMPLATE/         # Issue templates
│
├── src/                        # Source code
│   ├── components/             # Reusable UI components
│   │   ├── TaskList/
│   │   ├── TaskForm/
│   │   └── ...
│   │
│   ├── services/               # Business logic and API calls
│   │   ├── taskService.js
│   │   ├── authService.js
│   │   └── ...
│   │
│   ├── models/                 # Data models
│   │   ├── Task.js
│   │   ├── User.js
│   │   └── ...
│   │
│   ├── utils/                  # Utility functions
│   │   ├── dateUtils.js
│   │   ├── validation.js
│   │   └── ...
│   │
│   ├── hooks/                  # Custom hooks (React)
│   │   └── useTask.js
│   │
│   ├── store/                  # State management
│   │   └── taskStore.js
│   │
│   ├── constants/              # Constants and enums
│   │   └── taskConstants.js
│   │
│   ├── styles/                 # Global styles
│   │   └── globals.css
│   │
│   ├── types/                  # TypeScript type definitions
│   │   └── task.ts
│   │
│   ├── config.js               # Configuration
│   └── index.js                # Application entry point
│
├── tests/                      # Test files
│   ├── unit/                   # Unit tests
│   │   ├── services/
│   │   ├── utils/
│   │   └── models/
│   │
│   ├── integration/            # Integration tests
│   │   └── taskFlow.test.js
│   │
│   └── e2e/                    # End-to-end tests
│       └── taskManagement.e2e.js
│
├── docs/                       # Documentation
│   ├── API.md                  # API documentation
│   ├── ARCHITECTURE.md         # System architecture
│   ├── DATABASE.md             # Database schema
│   └── SETUP.md                # Setup guide
│
├── public/                     # Static assets
│   ├── images/
│   ├── icons/
│   └── favicon.ico
│
├── config/                     # Configuration files
│   ├── database.js
│   ├── server.js
│   └── env.example
│
├── .env.example               # Environment variables template
├── .editorconfig              # Editor configuration
├── .eslintrc.json            # ESLint configuration
├── .gitignore                # Git ignore rules
├── .prettierrc                # Prettier configuration
├── commitlint.config.js       # Commit lint configuration
├── jest.config.js            # Jest configuration
├── package.json              # Project dependencies
├── README.md                 # Project overview
├── CONTRIBUTING.md           # Contributing guidelines
└── LICENSE                   # License file
```

## File Size Guidelines

Keep files focused and manageable:

- **Component files**: Max 200 lines
- **Service files**: Max 300 lines
- **Utility files**: Max 150 lines
- **Test files**: Usually larger is okay (comprehensive coverage)

If a file exceeds these limits, consider breaking it into smaller, focused pieces.

## Naming Conventions

### Components
```
src/components/
├── TaskList.jsx               # PascalCase, export default
├── TaskForm.jsx
└── TaskCard/
    ├── TaskCard.jsx           # Component file
    ├── TaskCard.module.css    # Scoped styles (optional)
    └── TaskCard.test.js       # Test file
```

### Services
```
src/services/
├── taskService.js             # camelCase
├── userService.js
└── api.js
```

### Utils
```
src/utils/
├── dateUtils.js               # camelCase
├── validation.js
└── formatters.js
```

### Tests
```
# Co-locate tests next to source or in tests/ folder
src/utils/dateUtils.js
src/utils/dateUtils.test.js

# or

tests/utils/dateUtils.test.js
```

## Import Organization

Organize imports in this order:

```javascript
// 1. External libraries
import React from 'react';
import axios from 'axios';

// 2. Internal components
import TaskList from '../components/TaskList';
import TaskForm from '../components/TaskForm';

// 3. Internal services
import { taskService } from '../services';

// 4. Utils and helpers
import { formatDate } from '../utils/dateUtils';

// 5. Styles
import './TaskPage.css';

// 6. Types (TypeScript)
import { Task } from '../types';
```

## Comments and Documentation

### Component Documentation
```javascript
/**
 * TaskList component
 * @component
 * 
 * @param {Object} props - Component props
 * @param {Array<Task>} props.tasks - List of tasks to display
 * @param {Function} props.onEdit - Callback for edit action
 * @param {Function} props.onDelete - Callback for delete action
 * 
 * @returns {React.ReactElement} Rendered task list
 * 
 * @example
 * <TaskList 
 *   tasks={tasks}
 *   onEdit={handleEdit}
 *   onDelete={handleDelete}
 * />
 */
function TaskList({ tasks, onEdit, onDelete }) {
  // implementation
}
```

### Function Documentation
```javascript
/**
 * Formats a date to a human-readable string
 * @param {Date} date - The date to format
 * @param {string} [format='en-US'] - The locale format
 * @returns {string} Formatted date string
 */
function formatDate(date, format = 'en-US') {
  // implementation
}
```

## Configuration Files

Each major area can have its own configuration:

```
config/
├── database.js        # Database connection config
├── server.js          # Server configuration
├── logger.js          # Logging configuration
└── cache.js           # Cache configuration
```

## Environment Variables

Always use `.env.example` as a template:

```
# .env.example (commit this)
DATABASE_URL=
API_KEY=
JWT_SECRET=
DEBUG=false
```

```
# .env (gitignored, for local development)
DATABASE_URL=postgres://user:pass@localhost/dbname
API_KEY=your-api-key-here
JWT_SECRET=your-secret-key
DEBUG=true
```

## Build Artifacts

Never commit build artifacts. Configure `.gitignore`:

```
dist/
build/
*.min.js
*.min.css
.next/
out/
coverage/
```

---

This structure promotes:
- ✅ Clear separation of concerns
- ✅ Easy navigation and refactoring
- ✅ Better code organization
- ✅ Scalability as project grows
- ✅ Consistency across team
