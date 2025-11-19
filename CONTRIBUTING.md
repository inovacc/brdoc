# Contributing to brdoc

First off, thank you for considering contributing to brdoc! 🎉

## 🤝 How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples**
- **Describe the behavior you observed and what you expected**
- **Include Go version and OS information**

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description of the suggested enhancement**
- **Explain why this enhancement would be useful**
- **List any similar features in other libraries**

### Pull Requests

1. Fork the repository
2. Create a new branch from `main`:
   ```bash
   git checkout -b feature/amazing-feature
   ```

3. Make your changes following our coding standards:

- Write clear, commented code
- Follow Go best practices and idioms
- Add tests for new functionality
- Update documentation as needed

4. Ensure all tests pass:
   ```bash
   go test -v ./...
   go test -race ./...
   ```

5. Format your code:
   ```bash
   go fmt ./...
   ```

6. Run linter:
   ```bash
   golangci-lint run
   ```

7. Commit your changes:
   ```bash
   git commit -m 'Add amazing feature'
   ```

8. Push to your fork:
   ```bash
   git push origin feature/amazing-feature
   ```

9. Open a Pull Request

## 📝 Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Write clear and concise comments
- Keep functions focused and small
- Use meaningful variable and function names

## 🧪 Testing

- Write unit tests for all new functionality
- Maintain test coverage above 90%
- Add benchmark tests for performance-critical code
- Test edge cases and error conditions

Example test structure:

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid case", "input", "output", false},
        {"error case", "bad", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Feature(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Feature() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("Feature() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## 📚 Documentation

- Update README.md for user-facing changes
- Add godoc comments for exported functions and types
- Include examples in documentation
- Update CHANGELOG.md (if exists)

## 🏗️ Project Structure

```
brdoc/
├── validator.go          # Main implementation
├── validator_test.go     # Test suite
├── examples/             # Usage examples
├── .github/              # GitHub specific files
│   └── workflows/        # CI/CD workflows
├── go.mod
├── LICENSE
└── README.md
```

## 🔍 Code Review Process

1. Maintainers will review your PR
2. Address any feedback or requested changes
3. Once approved, a maintainer will merge your PR

## 📋 Checklist

Before submitting your PR, ensure:

- [ ] Code follows project style guidelines
- [ ] Tests are added/updated and passing
- [ ] Documentation is updated
- [ ] Commit messages are clear and descriptive
- [ ] No unnecessary dependencies added
- [ ] Code is formatted with `go fmt`
- [ ] No linter warnings

## 🎯 Priority Areas

We're particularly interested in contributions for:

- Additional Brazilian document validators (RG, CNH, PIS, etc.)
- Performance improvements
- Documentation improvements
- Bug fixes
- Test coverage improvements

## ❓ Questions?

Feel free to open an issue for questions or reach out through:

- GitHub Issues
- GitHub Discussions

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to brdoc! 🚀
