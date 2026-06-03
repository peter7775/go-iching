# Contributing to I Ching API

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the code, not the person
- Help others learn and grow

## How to Contribute

### Reporting Bugs

1. Check existing [issues](https://github.com/peter7775/go-iching/issues) to avoid duplicates
2. Create a new issue with:
   - Clear, descriptive title
   - Detailed description of the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version, etc.)

### Suggesting Enhancements

1. Check existing [discussions](https://github.com/peter7775/go-iching/discussions)
2. Create a feature request issue explaining:
   - The use case
   - Why it's needed
   - Possible implementation approaches

### Code Contributions

#### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/peter7775/go-iching.git
cd go-iching

# Install dependencies
go mod download

# Install linters
make install-lint

# Run tests
go test ./...
```

#### Before Submitting a PR

1. Fork the repository
2. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. Make your changes following code standards:
   ```bash
   # Format code
   make fmt

   # Run linter
   make lint

   # Run tests
   make test

   # Full audit
   make audit
   ```

4. Commit with clear messages:
   ```bash
   git commit -m "type: brief description

   Longer explanation if needed. Reference issues with #123.
   ```

   Commit types:
   - `feat`: New feature
   - `fix`: Bug fix
   - `docs`: Documentation
   - `style`: Formatting
   - `refactor`: Code reorganization
   - `test`: Tests
   - `build`: Build system
   - `ci`: CI/CD
   - `chore`: Maintenance

5. Push and create a Pull Request:
   ```bash
   git push origin feature/your-feature-name
   ```

#### PR Guidelines

- Link to related issues
- Describe the change and why it's needed
- Ensure tests pass (`make audit`)
- Update documentation if needed
- Keep PRs focused and reasonably sized

### Documentation

Documentation improvements are valuable! Please:

1. Update README.md for user-facing changes
2. Update PACKAGING.md for distribution changes
3. Add comments for complex logic
4. Update CHANGELOG.md for significant changes

### Testing

- Write tests for new features
- Ensure existing tests pass
- Aim for meaningful test coverage
- Use table-driven tests when appropriate

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"simple", "test", "result", false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Development Workflow

### Local Testing

```bash
# Start PostgreSQL for testing
make docker-up

# Run tests with PostgreSQL
STORAGE=postgres POSTGRES_DSN="..." go test ./...

# Stop containers
make docker-down
```

### Building Releases

```bash
# Run full audit
make audit

# Create git tag
git tag -a v0.2.0 -m "Release v0.2.0"

# GitHub Actions automatically builds and releases
git push origin v0.2.0
```

## Code Style

- Follow Go conventions and idioms
- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Keep functions focused and single-purpose
- Write clear, descriptive variable names
- Add comments for non-obvious logic

## Licensing

By contributing, you agree that your contributions will be licensed under the project's [LICENSE](LICENSE).

## Questions?

- Check [discussions](https://github.com/peter7775/go-iching/discussions)
- Open an issue with the "question" label
- Reach out via email: petrstepanek99@proton.me

## Recognition

- Contributors will be recognized in releases
- Thank you for making this project better!

---

Happy contributing! 🙏
