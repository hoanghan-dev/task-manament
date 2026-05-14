"# task-manament 📝

A modern, clean task management application built with best practices in mind.

## 📋 Overview

task-manament is a task management system designed to help users organize, prioritize, and track their tasks efficiently. The project emphasizes clean code, maintainability, and developer experience.

## ✨ Features

- Create, read, update, and delete tasks
- Priority and status management
- Task organization and filtering
- User authentication and authorization
- Real-time updates
- Responsive design

## 🚀 Quick Start

### Prerequisites
- Node.js 16+ 
- npm or yarn

### Installation

```bash
# Clone the repository
git clone https://github.com/hoanghan-dev/task-manament.git
cd task-manament

# Install dependencies
npm install

# Setup pre-commit hooks
npm run prepare

# Verify everything works
npm run lint
npm test
```

### Development

```bash
# Start development server
npm run dev

# Run tests in watch mode
npm run test:watch

# Format code
npm run format

# Lint code
npm run lint:fix
```

## 📁 Project Structure

```
src/               # Source code
tests/             # Test files
docs/              # Documentation
config/            # Configuration files
public/            # Static assets
```

See [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) for detailed structure information.

## 🏗️ Architecture

This project follows clean architecture principles:

- **Separation of Concerns**: Clear separation between components, services, and utilities
- **Single Responsibility**: Each module has a single, well-defined purpose
- **Dependency Injection**: Services are injected, not instantiated
- **Testability**: Code is written to be easily testable

## 📚 Documentation

- [CODE_REVIEW.md](./CODE_REVIEW.md) - Comprehensive code review and guidelines
- [BEST_PRACTICES.md](./BEST_PRACTICES.md) - Best practices and patterns
- [CONTRIBUTING.md](./CONTRIBUTING.md) - Contributing guidelines
- [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) - Detailed project structure
- [docs/](./docs/) - Additional documentation

## 🧪 Testing

```bash
# Run all tests
npm test

# Run tests with coverage
npm test:coverage

# Run tests in watch mode
npm run test:watch
```

Target test coverage: **80%+**

## 🔍 Code Quality

This project maintains high code quality through:

### Linting
```bash
npm run lint        # Check for issues
npm run lint:fix    # Fix issues automatically
```

**Tools**: ESLint with strict rules

### Formatting
```bash
npm run format      # Format all files
npm run format:check # Check if formatting needed
```

**Tools**: Prettier for consistent formatting

### Pre-commit Hooks
Automatic checks before each commit:
- Linting on changed files
- Code formatting
- Test validation

## 🔐 Code Standards

We follow these principles:

- **Clean Code**: Write code that is easy to understand and maintain
- **DRY**: Don't Repeat Yourself
- **SOLID**: Design principles for maintainable code
- **Conventional Commits**: Clear commit message format
- **Meaningful Naming**: Clear variable and function names
- **Error Handling**: Proper error handling throughout

See [CODE_REVIEW.md](./CODE_REVIEW.md) for comprehensive guidelines.

## 🤝 Contributing

We welcome contributions! Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for:
- Development setup
- Coding standards
- Pull request process
- Testing requirements

### Quick Contribution Checklist
- [ ] Code follows style guide (`npm run lint:fix`)
- [ ] All tests pass (`npm test`)
- [ ] New tests added for new features
- [ ] Meaningful commit messages
- [ ] PR description is clear and complete

## 🔄 Workflow

1. **Branch**: Create a feature branch (`feature/task-name`)
2. **Code**: Write code following guidelines
3. **Test**: Add and run tests
4. **Lint**: Fix any linting issues
5. **Commit**: Make meaningful commits
6. **Push**: Push to your branch
7. **PR**: Create a pull request for review
8. **Review**: Address review comments
9. **Merge**: Merge to main after approval

## 📊 Project Status

- **Status**: 🎯 In Development
- **Node Version**: 16+
- **License**: MIT
- **Last Updated**: May 2026

## 🛠️ Tools & Technologies

### Development
- **Language**: JavaScript/TypeScript
- **Testing**: Jest
- **Linting**: ESLint
- **Formatting**: Prettier
- **Pre-commit**: Husky + lint-staged

### Quality Assurance
- Automated linting on every commit
- Unit and integration tests
- Code coverage tracking
- Security vulnerability scanning

## 📋 Checklist for New Developers

Getting started? Follow this checklist:

- [ ] Read this README
- [ ] Read [CONTRIBUTING.md](./CONTRIBUTING.md)
- [ ] Review [CODE_REVIEW.md](./CODE_REVIEW.md)
- [ ] Review [BEST_PRACTICES.md](./BEST_PRACTICES.md)
- [ ] Understand [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md)
- [ ] Install dependencies: `npm install`
- [ ] Run tests: `npm test`
- [ ] Read existing code in `src/`
- [ ] Ask questions in discussions

## 🐛 Reporting Issues

Found a bug? Please create an issue with:
- Clear title and description
- Steps to reproduce
- Expected vs actual behavior
- Environment details

## 💬 Questions or Discussions?

Open a discussion in the GitHub repo for:
- Questions about the codebase
- Design discussions
- Feature suggestions
- General help

## 📝 License

This project is licensed under the MIT License - see [LICENSE](./LICENSE) file for details.

## 👥 Contributors

Thanks to all contributors! See [CONTRIBUTING.md](./CONTRIBUTING.md) to get started.

---

**Made with ❤️ by the task-manament team**" 
