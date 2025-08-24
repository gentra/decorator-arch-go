# Contributing to Decorator Architecture Go

Thank you for your interest in contributing to this project! This document provides guidelines and information for contributors.

## Development Setup

1. **Prerequisites**
   - Go 1.24.5 or later
   - Git

2. **Clone the repository**
   ```bash
   git clone https://github.com/gentra/decorator-arch-go.git
   cd decorator-arch-go
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Run tests**
   ```bash
   go test ./...
   ```

## Code Style and Standards

- Follow Go's official [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for code formatting
- Run `golangci-lint` before submitting PRs
- Write unit tests for new functionality
- Use table-driven tests with Gherkin syntax for test names

## Architectural Rules

This project follows strict architectural rules:

1. **Single Interface Per Domain**: Each domain has ONLY ONE interface called "Service"
2. **Implementation Folder Constraint**: Implementation folders can ONLY contain files that implement the domain's Service interface
3. **Interface Extraction Rule**: New interfaces become separate domain folders
4. **Factory-Managed Composition**: Strategy and decorator patterns are handled in domain factory folders

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following the architectural rules
4. Add tests for new functionality
5. Ensure all tests pass
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Testing Guidelines

- Write unit tests in table-driven style
- Use Gherkin syntax for test names (Given, When, Then)
- Use `stretchr/testify` for mocks and assertions
- Aim for high test coverage

## Questions or Need Help?

- Open an issue for bugs or feature requests
- Use the discussion tab for questions
- Check existing issues and PRs first

Thank you for contributing!
