package brdoc

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// ============================================================================
// CPF Tests
// ============================================================================

func TestCPF_Generate(t *testing.T) {
	cpf := NewCPF()

	for range 10 {
		generated := cpf.Generate()

		if !cpf.Validate(generated) {
			t.Errorf("Generated CPF is invalid: %s", generated)
		}

		_, _ = fmt.Fprintf(os.Stdout, "Generated CPF: %s | Origin: %s\n", generated, cpf.CheckOrigin(generated))
	}
}

func TestCPF_Validate(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{"Valid formatted CPF", "123.456.789-09", true},
		{"Valid unformatted CPF", "12345678909", true},
		{"Invalid CPF - wrong check digit", "123.456.789-00", false},
		{"Invalid CPF - all zeros", "000.000.000-00", false},
		{"Invalid CPF - all equal digits", "111.111.111-11", false},
		{"Invalid CPF - wrong length", "123.456.789", false},
	}

	cpf := NewCPF()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cpf.Validate(tt.cpf)
			if result != tt.expected {
				t.Errorf("Validate(%s) = %v, expected %v", tt.cpf, result, tt.expected)
			}
		})
	}
}

func TestCPF_Format(t *testing.T) {
	cpf := NewCPF()

	input := "12345678909"
	expected := "123.456.789-09"

	result, err := cpf.Format(input)
	require.NoError(t, err)

	if result != expected {
		t.Errorf("Format(%s) = %s, expected %s", input, result, expected)
	}
}

func TestCPF_CheckOrigin(t *testing.T) {
	tests := []struct {
		cpf      string
		expected string
	}{
		{"123.456.780-09", IsDigit0},
		{"123.456.788-09", IsDigit8},
		{"123.456.789-09", IsDigit9},
	}

	cpf := NewCPF()
	for _, tt := range tests {
		t.Run(tt.cpf, func(t *testing.T) {
			result := cpf.CheckOrigin(tt.cpf)
			if result != tt.expected {
				t.Errorf("CheckOrigin(%s) = %s, expected %s", tt.cpf, result, tt.expected)
			}
		})
	}
}

// Additional real-world valid CPF samples (provided) — all must validate
func TestCPF_Validate_ProvidedValid(t *testing.T) {
	samples := []string{
		"013.723.737-56",
		"260.808.754-03",
		"205.117.448-20",
		"213.872.640-10",
		"722.628.653-02",
		"747.356.416-10",
		"486.158.855-32",
	}

	cpf := NewCPF()

	for _, s := range samples {
		t.Run(s, func(t *testing.T) {
			if !cpf.Validate(s) {
				t.Fatalf("Expected provided CPF to be valid: %s", s)
			}

			formatted, err := cpf.Format(s)
			if err != nil {
				t.Fatalf("Unexpected formatting error for %s: %v", s, err)
			}

			if !cpf.Validate(formatted) {
				t.Fatalf("Formatted CPF should be valid: %s", formatted)
			}
		})
	}
}

// Ensure unformatted variants of provided CPFs are also valid and format back correctly
func TestCPF_Format_ProvidedValid_Unformatted(t *testing.T) {
	strip := func(s string) string {
		out := make([]rune, 0, len(s))
		for _, r := range s {
			switch r {
			case '.', '-', ' ':
				// skip
			default:
				out = append(out, r)
			}
		}

		return string(out)
	}

	samples := []string{
		"013.723.737-56",
		"260.808.754-03",
		"205.117.448-20",
		"213.872.640-10",
		"722.628.653-02",
		"747.356.416-10",
		"486.158.855-32",
	}

	cpf := NewCPF()

	for _, formatted := range samples {
		unformatted := strip(formatted)
		t.Run(unformatted, func(t *testing.T) {
			if !cpf.Validate(unformatted) {
				t.Fatalf("Expected unformatted provided CPF to be valid: %s", unformatted)
			}

			got, err := cpf.Format(unformatted)
			if err != nil {
				t.Fatalf("Unexpected formatting error for %s: %v", unformatted, err)
			}

			if got != formatted {
				t.Fatalf("Format(%s) = %s, expected %s", unformatted, got, formatted)
			}
		})
	}
}

// ============================================================================
// Alphanumeric CNPJ Tests
// ============================================================================

func TestCNPJ_ValidateExampleFromPDF(t *testing.T) {
	cnpj := NewCNPJ()

	// Example from SERPRO documentation
	example := "12ABC34501DE35"

	if !cnpj.Validate(example) {
		t.Errorf("SERPRO example CNPJ should be valid: %s", example)
	}

	// Also test with formatting
	formattedExample := "12.ABC.345/01DE-35"
	if !cnpj.Validate(formattedExample) {
		t.Errorf("Formatted CNPJ should be valid: %s", formattedExample)
	}
}

// Additional real-world valid CNPJ samples (provided) — all must validate
func TestCNPJ_Validate_ProvidedValid(t *testing.T) {
	samples := []string{
		"HR.YUP.H8D/0001-02",
		"48.175.226/0001-50",
		"SE.URZ.76B/0001-02",
		"37.077.670/0001-16",
		"52.311.151/0001-64",
		"64.814.243/0001-46",
		"Z7.BM3.7VE/0001-93",
		"V2.P0M.NVE/0001-07",
	}

	cnpj := NewCNPJ()

	for _, s := range samples {
		t.Run(s, func(t *testing.T) {
			if !cnpj.Validate(s) {
				t.Fatalf("Expected provided CNPJ to be valid: %s", s)
			}

			// Round-trip formatting should keep the same visual representation (uppercase)
			formatted, err := cnpj.Format(s)
			if err != nil {
				t.Fatalf("Unexpected formatting error for %s: %v", s, err)
			}

			// Validate formatted too
			if !cnpj.Validate(formatted) {
				t.Fatalf("Formatted CNPJ should be valid: %s", formatted)
			}
		})
	}
}

// Ensure unformatted variants of provided CNPJs are also valid and format back correctly
func TestCNPJ_Format_ProvidedValid_Unformatted(t *testing.T) {
	strip := func(s string) string {
		// remove formatting characters .-/ and spaces
		out := make([]rune, 0, len(s))
		for _, r := range s {
			switch r {
			case '.', '-', '/', ' ':
				// skip
			default:
				out = append(out, r)
			}
		}

		return string(out)
	}

	samples := []string{
		"HR.YUP.H8D/0001-02",
		"48.175.226/0001-50",
		"SE.URZ.76B/0001-02",
		"37.077.670/0001-16",
		"52.311.151/0001-64",
		"64.814.243/0001-46",
		"Z7.BM3.7VE/0001-93",
		"V2.P0M.NVE/0001-07",
	}

	cnpj := NewCNPJ()

	for _, formatted := range samples {
		unformatted := strip(formatted)
		t.Run(unformatted, func(t *testing.T) {
			if !cnpj.Validate(unformatted) {
				t.Fatalf("Expected unformatted provided CNPJ to be valid: %s", unformatted)
			}

			got, err := cnpj.Format(unformatted)
			if err != nil {
				t.Fatalf("Unexpected formatting error for %s: %v", unformatted, err)
			}

			// The formatter normalizes to uppercase and standard mask; compare after normalizing expected to uppercase
			expected := strings.ToUpper(formatted)
			if got != expected {
				t.Fatalf("Format(%s) = %s, expected %s", unformatted, got, expected)
			}
		})
	}
}

func TestCNPJ_Generate(t *testing.T) {
	cnpj := NewCNPJ()

	for range 10 {
		generated := cnpj.Generate()

		if !cnpj.Validate(generated) {
			t.Errorf("Generated CNPJ is invalid: %s", generated)
		}

		formatted, err := cnpj.Format(generated)
		if err != nil {
			t.Errorf("Error formatting CNPJ: %v", err)
		}

		_, _ = fmt.Fprintf(os.Stdout, "Generated CNPJ: %s | Formatted: %s\n", generated, formatted)
	}
}

func TestCNPJ_Validate(t *testing.T) {
	tests := []struct {
		name     string
		cnpj     string
		expected bool
	}{
		{"Valid CNPJ - SERPRO example", "12ABC34501DE35", true},
		{"Valid formatted CNPJ", "12.ABC.345/01DE-35", true},
		{"Invalid CNPJ - wrong check digits", "12ABC34501DE00", false},
		{"Invalid CNPJ - wrong length", "12ABC345", false},
		{"Invalid CNPJ - non-numeric check digits", "12ABC34501DEAA", false},
	}

	cnpj := NewCNPJ()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cnpj.Validate(tt.cnpj)
			if result != tt.expected {
				t.Errorf("Validate(%s) = %v, expected %v", tt.cnpj, result, tt.expected)
			}
		})
	}
}

func TestCNPJ_Format(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			"Valid CNPJ without formatting",
			"12ABC34501DE35",
			"12.ABC.345/01DE-35",
			false,
		},
		{
			"CNPJ already formatted",
			"12.ABC.345/01DE-35",
			"12.ABC.345/01DE-35",
			false,
		},
		{
			"CNPJ invalid length",
			"12ABC345",
			"",
			true,
		},
	}

	cnpj := NewCNPJ()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cnpj.Format(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for input: %s", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if result != tt.expected {
					t.Errorf("Format(%s) = %s, expected %s", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestCNPJ_CalculateDV_Manual(t *testing.T) {
	cnpj := NewCNPJ()

	// Manual test of SERPRO example: 12ABC34501DE
	base := "12ABC34501DE"

	dv1, err := cnpj.calculateDV(base)
	if err != nil {
		t.Fatalf("Error calculating DV1: %v", err)
	}

	if dv1 != 3 {
		t.Errorf("DV1 calculated: %d, expected: 3", dv1)
	}

	dv2, err := cnpj.calculateDV(base + "3")
	if err != nil {
		t.Fatalf("Error calculating DV2: %v", err)
	}

	if dv2 != 5 {
		t.Errorf("DV2 calculated: %d, expected: 5", dv2)
	}

	_, _ = fmt.Fprintf(os.Stdout, "✓ Check digits calculated correctly: %d%d\n", dv1, dv2)
}

// ============================================================================
// Utility Functions Tests
// ============================================================================

func TestValidateDocument(t *testing.T) {
	tests := []struct {
		name    string
		doc     string
		docType string
		isValid bool
	}{
		{"Valid CPF", "123.456.789-09", "CPF", true},
		{"Valid CNPJ", "12.ABC.345/01DE-35", "CNPJ", true},
		{"Invalid CPF", "123.456.789-00", "CPF", false},
		{"Invalid CNPJ", "12.ABC.345/01DE-00", "CNPJ", false},
		{"Unknown document", "12345", "UNKNOWN", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docType, isValid := ValidateDocument(tt.doc)

			if docType != tt.docType {
				t.Errorf("Expected type: %s, got: %s", tt.docType, docType)
			}

			if isValid != tt.isValid {
				t.Errorf("Validation expected: %v, got: %v", tt.isValid, isValid)
			}
		})
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkCPF_Generate(b *testing.B) {
	cpf := NewCPF()
	for b.Loop() {
		_ = cpf.Generate()
	}
}

func BenchmarkCPF_Validate(b *testing.B) {
	cpf := NewCPF()

	testCPF := "123.456.789-09"
	for b.Loop() {
		cpf.Validate(testCPF)
	}
}

func BenchmarkCNPJ_Generate(b *testing.B) {
	cnpj := NewCNPJ()
	for b.Loop() {
		_ = cnpj.Generate()
	}
}

func BenchmarkCNPJ_Validate(b *testing.B) {
	cnpj := NewCNPJ()

	testCNPJ := "12.ABC.345/01DE-35"
	for b.Loop() {
		cnpj.Validate(testCNPJ)
	}
}
