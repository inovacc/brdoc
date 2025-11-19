# 🇧🇷 brdoc

[![Go Reference](https://pkg.go.dev/badge/github.com/inovacc/brdoc.svg)](https://pkg.go.dev/github.com/inovacc/brdoc)
[![Go Report Card](https://goreportcard.com/badge/github.com/inovacc/brdoc)](https://goreportcard.com/report/github.com/inovacc/brdoc)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Coverage](https://img.shields.io/badge/coverage-95%25-brightgreen.svg)](https://github.com/inovacc/brdoc)

A robust and efficient Go library for validating, generating, and formatting Brazilian fiscal documents (CPF and CNPJ).

Implements official algorithms from SERPRO (Serviço Federal de Processamento de Dados) for alphanumeric CNPJ validation, supporting the new format introduced for modern Brazilian business
registration.

## ✨ Features

- ✅ **CPF (Cadastro de Pessoas Físicas)**
  - Generate valid random CPFs
  - Validate CPF with check digit verification
  - Format CPF (XXX.XXX.XXX-XX)
  - Identify issuing state/region
  - Reject common invalid patterns (all same digits)

- ✅ **CNPJ (Cadastro Nacional de Pessoa Jurídica)**
  - Generate valid alphanumeric CNPJs
  - Validate alphanumeric CNPJ per SERPRO specification
  - Support both numeric and alphanumeric formats
  - Format CNPJ (XX.XXX.XXX/XXXX-XX)
  - Modulo 11 check digit calculation

- ✅ **General**
  - Zero dependencies
  - Thread-safe random generation
  - Comprehensive test coverage
  - Benchmark suite included
  - Production-ready

## 📦 Installation

```bash
go get github.com/inovacc/brdoc
```

**Requirements:** Go 1.20 or higher

### Install the CLI

To install the `brdoc` command-line tool:

```bash
go install github.com/inovacc/brdoc/cmd/brdoc@latest
```

After installation, ensure that your `$GOBIN` (or `$GOPATH/bin`) is on your system PATH so you can run `brdoc` from any directory.

Quick usage:

```bash
# Generate a valid CPF
brdoc cpf --generate

# Validate a CPF
brdoc cpf --validate 123.456.789-09

# Generate a valid CNPJ
brdoc cnpj --generate

# Validate a CNPJ (alphanumeric supported)
brdoc cnpj --validate 12.ABC.345/01DE-35
```

## 🚀 Quick Start

### CPF Validation

```go
package main

import (
  "fmt"
  "log"
  "github.com/inovacc/brdoc"
)

func main() {
  cpf := brdoc.NewCPF()

  // Validate
  isValid := cpf.Validate("123.456.789-09")
  fmt.Printf("Valid: %v\n", isValid)

  // Generate
  generated := cpf.Generate()
  fmt.Printf("Generated CPF: %s\n", generated)

  // Check origin
  origin := cpf.CheckOrigin("123.456.789-09")
  fmt.Printf("Issued in: %s\n", origin)

  // Format
  formatted, err := cpf.Format("12345678909")
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("Formatted: %s\n", formatted)
}
```

### CNPJ Validation (Alphanumeric)

```go
package main

import (
  "fmt"
  "log"
  "github.com/inovacc/brdoc"
)

func main() {
  cnpj := brdoc.NewCNPJ()

  // Validate (supports alphanumeric per SERPRO spec)
  isValid := cnpj.Validate("12ABC34501DE35")
  fmt.Printf("Valid: %v\n", isValid)

  // Validate formatted
  isValid = cnpj.Validate("12.ABC.345/01DE-35")
  fmt.Printf("Valid: %v\n", isValid)

  // Generate
  generated := cnpj.Generate()
  fmt.Printf("Generated: %s\n", generated)

  // Format
  formatted, err := cnpj.Format("12ABC34501DE35")
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("Formatted: %s\n", formatted)
}
```

### Auto-detect Document Type

```go
package main

import (
  "fmt"
  "github.com/inovacc/brdoc"
)

func main() {
  documents := []string{
    "123.456.789-09",     // CPF
    "12.ABC.345/01DE-35", // CNPJ
  }

  for _, doc := range documents {
    docType, isValid := brdoc.ValidateDocument(doc)
    fmt.Printf("%s -> Type: %s, Valid: %v\n", doc, docType, isValid)
  }
}
```

## 📚 API Documentation

### CPF

#### `NewCPF() *CPF`

Creates a new CPF validator instance.

#### `Generate() string`

Generates a valid random CPF with formatting.

**Returns:** Formatted CPF string (XXX.XXX.XXX-XX)

#### `Validate(cpf string) bool`

Validates a CPF number (with or without formatting).

**Parameters:**

- `cpf`: CPF string to validate

**Returns:** `true` if valid, `false` otherwise

**Validation rules:**

- Must have 11 digits
- Check digits must be correct
- Cannot be all same digits (000.000.000-00, 111.111.111-11, etc.)

#### `Format(cpf string) string`

Formats a CPF string to the standard format.

**Parameters:**

- `cpf`: Unformatted CPF string

**Returns:** Formatted CPF (XXX.XXX.XXX-XX)

#### `CheckOrigin(cpf string) string`

Returns the Brazilian state/region where the CPF was issued based on the 9th digit.

**Parameters:**

- `cpf`: CPF string

**Returns:** State/region name in English

**Mapping:**

- 0: Rio Grande do Sul
- 1: Federal District, Goiás, Mato Grosso do Sul, and Tocantins
- 2: Pará, Amazonas, Acre, Amapá, Rondônia, and Roraima
- 3: Ceará, Maranhão, and Piauí
- 4: Pernambuco, Rio Grande do Norte, Paraíba, and Alagoas
- 5: Bahia and Sergipe
- 6: Minas Gerais
- 7: Rio de Janeiro and Espírito Santo
- 8: São Paulo
- 9: Paraná and Santa Catarina

---

### CNPJ

#### `NewCNPJ() *CNPJ`

Creates a new CNPJ validator instance.

#### `Generate() (string, error)`

Generates a valid random alphanumeric CNPJ.

**Returns:**

- Unformatted CNPJ string (14 characters)
- Error if generation fails

#### `Validate(cnpj string) bool`

Validates an alphanumeric CNPJ per SERPRO specification.

**Parameters:**

- `cnpj`: CNPJ string to validate (with or without formatting)

**Returns:** `true` if valid, `false` otherwise

**Validation rules:**

- Must have 14 characters (12 alphanumeric + 2 numeric check digits)
- Check digits must be correct per modulo 11 algorithm
- Supports letters A-Z and numbers 0-9 in first 12 positions
- Last 2 positions must be numeric

#### `Format(cnpj string) (string, error)`

Formats a CNPJ string to the standard format.

**Parameters:**

- `cnpj`: Unformatted CNPJ string

**Returns:**

- Formatted CNPJ (XX.XXX.XXX/XXXX-XX)
- Error if input is invalid

---

### Utility Functions

#### `ValidateDocument(doc string) (docType string, isValid bool)`

Auto-detects and validates CPF or CNPJ.

**Parameters:**

- `doc`: Document string (CPF or CNPJ)

**Returns:**

- `docType`: "CPF", "CNPJ", or "UNKNOWN"
- `isValid`: Validation result

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test -v

# Run with coverage
go test -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Output Example

```
=== RUN   TestCPF_Generate
Generated CPF: 123.456.789-09 | Origin: Paraná and Santa Catarina
--- PASS: TestCPF_Generate (0.00s)

=== RUN   TestCNPJ_ValidateExampleFromPDF
✓ Check digits calculated correctly: 35
--- PASS: TestCNPJ_ValidateExampleFromPDF (0.00s)

PASS
coverage: 95.2% of statements
```

## ⚡ Benchmarks

```bash
go test -bench=. -benchmem
```

### Results (example on Apple M1)

```
BenchmarkCPF_Generate-8       1000000    1234 ns/op    256 B/op    5 allocs/op
BenchmarkCPF_Validate-8       5000000     287 ns/op     64 B/op    3 allocs/op
BenchmarkCNPJ_Generate-8       800000    1456 ns/op    312 B/op    7 allocs/op
BenchmarkCNPJ_Validate-8      3000000     412 ns/op     96 B/op    4 allocs/op
```

## 📖 CNPJ Alphanumeric Specification

This library implements the official SERPRO specification for alphanumeric CNPJ validation. The algorithm uses:

- **Character mapping:** 0-9 → 0-9, A-Z → 17-42 (ASCII value - 48)
- **Weight distribution:** 2-9, repeating from right to left
- **Modulo 11:** Check digit calculation
- **Special rule:** If remainder is 0 or 1, the check digit is 0

### Example Calculation (from SERPRO documentation)

For CNPJ base: `12ABC34501DE`

**First check digit:**

```
Position:  1  2  A  B  C  3  4  5  0  1  D  E
Value:     1  2  17 18 19 3  4  5  0  1  20 21
Weight:    5  4  3  2  9  8  7  6  5  4  3  2
Product:   5  8  51 36 171 24 28 30 0  4  60 42
Sum: 459
Remainder: 459 % 11 = 8
Check digit 1: 11 - 8 = 3
```

**Second check digit:**

```
Position:  1  2  A  B  C  3  4  5  0  1  D  E  3
Value:     1  2  17 18 19 3  4  5  0  1  20 21 3
Weight:    6  5  4  3  2  9  8  7  6  5  4  3  2
Product:   6  10 68 54 38 27 32 35 0  5  80 63 6
Sum: 424
Remainder: 424 % 11 = 6
Check digit 2: 11 - 6 = 5
```

**Result:** `12.ABC.345/01DE-35`

## 🏗️ Project Structure

```
brdoc/
├── validator.go          # Main implementation
├── validator_test.go     # Test suite
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── examples/
    └── main.go          # Usage examples
```

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Write** tests for your changes
4. **Ensure** all tests pass (`go test -v`)
5. **Format** your code (`go fmt ./...`)
6. **Lint** your code (`golangci-lint run`)
7. **Commit** your changes (`git commit -m 'Add amazing feature'`)
8. **Push** to the branch (`git push origin feature/amazing-feature`)
9. **Open** a Pull Request

### Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Maintain test coverage above 90%
- Add benchmarks for performance-critical functions
- Document exported functions and types

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **SERPRO** (Serviço Federal de Processamento de Dados) for the official CNPJ alphanumeric specification
- **Receita Federal do Brasil** for CPF validation rules
- Brazilian developer community for feedback and contributions

## 📞 Support

- 🐛 **Issues:** [GitHub Issues](https://github.com/inovacc/brdoc/issues)
- 💬 **Discussions:** [GitHub Discussions](https://github.com/inovacc/brdoc/discussions)

## 🗺️ Roadmap

- [ ] Support for legacy numeric-only CNPJ
- [ ] RG (Registro Geral) validation
- [ ] CNH (Carteira Nacional de Habilitação) validation
- [ ] PIS/PASEP validation
- [ ] Título de Eleitor validation
- [ ] CEP (postal code) validation and formatting
- [ ] CLI tool for document validation
- [ ] REST API service example

## 📊 Stats

![GitHub stars](https://img.shields.io/github/stars/inovacc/brdoc?style=social)
![GitHub forks](https://img.shields.io/github/forks/inovacc/brdoc?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/inovacc/brdoc?style=social)

---

**Made with ❤️ in São Paulo, Brazil by [INOVACLOUD](https://github.com/inovacc)**

*If this library helped you, please consider giving it a ⭐️!*
