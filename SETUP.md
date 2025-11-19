# 🚀 brdoc - Setup Guide

## 📦 Package Contents

This ZIP contains the complete **brdoc** library for validating, generating, and formatting Brazilian fiscal documents (CPF and CNPJ).

### 📂 Project Structure

```
brdoc-project/
├── .github/
│   └── workflows/
│       └── ci.yml              # GitHub Actions CI/CD pipeline
├── examples/
│   ├── go.mod                  # Examples module configuration
│   └── main.go                 # Usage examples
├── .gitignore                  # Git ignore rules
├── CHANGELOG.md                # Version history
├── CONTRIBUTING.md             # Contribution guidelines
├── LICENSE                     # MIT License
├── README.md                   # Complete documentation
├── doc.go                      # Package documentation
├── go.mod                      # Module configuration
├── validator.go                # Main implementation
└── validator_test.go           # Test suite
```

## 🔧 Installation Steps

### 1. Upload to GitHub

```bash
# Extract the ZIP
unzip brdoc-v0.1.0.zip
cd brdoc-project

# Initialize git (if not already a git repo)
git init

# Add all files
git add .

# Commit
git commit -m "Initial commit: brdoc v0.1.0"

# Add remote (your repository)
git remote add origin https://github.com/inovacc/brdoc.git

# Push to GitHub
git push -u origin main
```

### 2. Verify Installation

```bash
# Run tests
go test -v

# Run benchmarks
go test -bench=. -benchmem

# Check coverage
go test -cover
```

### 3. Try the Examples

```bash
cd examples
go run main.go
```

Expected output:

```
=== 🇧🇷 Brazilian Document Validator Demo ===

📋 CPF Examples
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1️⃣  Generating random CPF:
   Generated: XXX.XXX.XXX-XX
   Origin: São Paulo

2️⃣  Validating CPFs:
   ✅ 123.456.789-09
   ❌ 111.111.111-11
   ❌ 000.000.000-00

...
```

## 📚 Usage in Your Project

### Install the package

```bash
go get github.com/inovacc/brdoc
```

### Import and use

```go
package main

import (
  "fmt"
  "github.com/inovacc/brdoc"
)

func main() {
  // CPF
  cpf := brdoc.NewCPF()
  fmt.Println(cpf.Validate("123.456.789-09")) // true or false

  // CNPJ
  cnpj := brdoc.NewCNPJ()
  fmt.Println(cnpj.Validate("12.ABC.345/01DE-35")) // true or false
}
```

## 🧪 Quality Checks

### Run all tests

```bash
go test -v ./...
```

### Check test coverage

```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run linter

```bash
golangci-lint run
```

### Run benchmarks

```bash
go test -bench=. -benchmem
```

## 📖 Documentation

### Generate godoc

```bash
godoc -http=:6060
# Visit: http://localhost:6060/pkg/github.com/inovacc/brdoc/
```

### View online

After pushing to GitHub:

- https://pkg.go.dev/github.com/inovacc/brdoc

## 🔐 Security

### Report vulnerabilities

```bash
# Check for known vulnerabilities
go list -json -m all | nancy sleuth
```

## 🚀 Release Process

### Creating a new release

1. Update CHANGELOG.md
2. Update version in go.mod (if needed)
3. Commit changes:
   ```bash
   git commit -m "Release v0.1.0"
   ```
4. Create and push tag:
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```
5. Create release on GitHub with release notes

## 🎯 Next Steps

1. ✅ Upload to GitHub: `https://github.com/inovacc/brdoc`
2. ✅ Enable GitHub Actions (CI will run automatically)
3. ✅ Add repository description and topics
4. ✅ Create first release (v0.1.0)
5. 📝 Share on social media/communities
6. 📊 Monitor usage via pkg.go.dev

## 📊 Badges Setup

Add these to your README (after first release):

```markdown
[![Go Reference](https://pkg.go.dev/badge/github.com/inovacc/brdoc.svg)](https://pkg.go.dev/github.com/inovacc/brdoc)
[![Go Report Card](https://goreportcard.com/badge/github.com/inovacc/brdoc)](https://goreportcard.com/report/github.com/inovacc/brdoc)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
```

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

## 📞 Support

- 🐛 Issues: https://github.com/inovacc/brdoc/issues
- 💬 Discussions: https://github.com/inovacc/brdoc/discussions

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Made with ❤️ by INOVACLOUD CONSULTORIA LTDA**

Repository: https://github.com/inovacc/brdoc
