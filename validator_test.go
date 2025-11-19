package brdoc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// ============================================================================
// Testes de CPF
// ============================================================================

func TestCPF_Generate(t *testing.T) {
	cpf := NewCPF()

	for i := 0; i < 10; i++ {
		generated := cpf.Generate()

		if !cpf.Validate(generated) {
			t.Errorf("CPF gerado inválido: %s", generated)
		}

		fmt.Printf("CPF gerado: %s | Origem: %s\n",
			generated,
			cpf.CheckOrigin(generated))
	}
}

func TestCPF_Validate(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{"CPF válido formatado", "123.456.789-09", true},
		{"CPF válido sem formatação", "12345678909", true},
		{"CPF inválido - dígito errado", "123.456.789-00", false},
		{"CPF inválido - todos zeros", "000.000.000-00", false},
		{"CPF inválido - todos iguais", "111.111.111-11", false},
		{"CPF inválido - tamanho errado", "123.456.789", false},
	}

	cpf := NewCPF()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cpf.Validate(tt.cpf)
			if result != tt.expected {
				t.Errorf("Validate(%s) = %v, esperado %v",
					tt.cpf, result, tt.expected)
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
		t.Errorf("Format(%s) = %s, esperado %s", input, result, expected)
	}
}

func TestCPF_CheckOrigin(t *testing.T) {
	tests := []struct {
		cpf      string
		expected string
	}{
		{"123.456.780-09", "Rio Grande do Sul"},
		{"123.456.788-09", "São Paulo"},
		{"123.456.789-09", "Paraná e Santa Catarina"},
	}

	cpf := NewCPF()
	for _, tt := range tests {
		t.Run(tt.cpf, func(t *testing.T) {
			result := cpf.CheckOrigin(tt.cpf)
			if result != tt.expected {
				t.Errorf("CheckOrigin(%s) = %s, esperado %s",
					tt.cpf, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Testes de CNPJ Alfanumérico
// ============================================================================

func TestCNPJ_ValidateExampleFromPDF(t *testing.T) {
	cnpj := NewCNPJ()

	// Exemplo do documento SERPRO
	exemplo := "12ABC34501DE35"

	if !cnpj.Validate(exemplo) {
		t.Errorf("CNPJ do exemplo SERPRO deveria ser válido: %s", exemplo)
	}

	// Testa com formatação
	exemploFormatado := "12.ABC.345/01DE-35"
	if !cnpj.Validate(exemploFormatado) {
		t.Errorf("CNPJ formatado deveria ser válido: %s", exemploFormatado)
	}
}

func TestCNPJ_Generate(t *testing.T) {
	cnpj := NewCNPJ()

	for i := 0; i < 10; i++ {
		generated := cnpj.Generate()

		if !cnpj.Validate(generated) {
			t.Errorf("CNPJ gerado inválido: %s", generated)
		}

		formatted, err := cnpj.Format(generated)
		if err != nil {
			t.Errorf("Erro ao formatar CNPJ: %v", err)
		}

		fmt.Printf("CNPJ gerado: %s | Formatado: %s\n", generated, formatted)
	}
}

func TestCNPJ_Validate(t *testing.T) {
	tests := []struct {
		name     string
		cnpj     string
		expected bool
	}{
		{"CNPJ válido - exemplo SERPRO", "12ABC34501DE35", true},
		{"CNPJ válido formatado", "12.ABC.345/01DE-35", true},
		{"CNPJ inválido - DV errado", "12ABC34501DE00", false},
		{"CNPJ inválido - tamanho errado", "12ABC345", false},
		{"CNPJ inválido - caractere inválido no DV", "12ABC34501DEAA", false},
	}

	cnpj := NewCNPJ()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cnpj.Validate(tt.cnpj)
			if result != tt.expected {
				t.Errorf("Validate(%s) = %v, esperado %v",
					tt.cnpj, result, tt.expected)
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
			"CNPJ válido sem formatação",
			"12ABC34501DE35",
			"12.ABC.345/01DE-35",
			false,
		},
		{
			"CNPJ já formatado",
			"12.ABC.345/01DE-35",
			"12.ABC.345/01DE-35",
			false,
		},
		{
			"CNPJ tamanho inválido",
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
					t.Errorf("Esperava erro para input: %s", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Erro inesperado: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Format(%s) = %s, esperado %s",
						tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestCNPJ_CalculateDV_Manual(t *testing.T) {
	cnpj := NewCNPJ()

	// Teste manual do exemplo SERPRO: 12ABC34501DE
	base := "12ABC34501DE"

	dv1, err := cnpj.calculateDV(base)
	if err != nil {
		t.Fatalf("Erro ao calcular DV1: %v", err)
	}

	if dv1 != 3 {
		t.Errorf("DV1 calculado: %d, esperado: 3", dv1)
	}

	dv2, err := cnpj.calculateDV(base + "3")
	if err != nil {
		t.Fatalf("Erro ao calcular DV2: %v", err)
	}

	if dv2 != 5 {
		t.Errorf("DV2 calculado: %d, esperado: 5", dv2)
	}

	fmt.Printf("✓ DVs calculados corretamente: %d%d\n", dv1, dv2)
}

// ============================================================================
// Testes de Funções Utilitárias
// ============================================================================

func TestValidateDocument(t *testing.T) {
	tests := []struct {
		name    string
		doc     string
		docType string
		isValid bool
	}{
		{"CPF válido", "123.456.789-09", "CPF", true},
		{"CNPJ válido", "12.ABC.345/01DE-35", "CNPJ", true},
		{"CPF inválido", "123.456.789-00", "CPF", false},
		{"CNPJ inválido", "12.ABC.345/01DE-00", "CNPJ", false},
		{"Documento desconhecido", "12345", "UNKNOWN", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docType, isValid := ValidateDocument(tt.doc)

			if docType != tt.docType {
				t.Errorf("Tipo esperado: %s, recebido: %s", tt.docType, docType)
			}

			if isValid != tt.isValid {
				t.Errorf("Validação esperada: %v, recebida: %v", tt.isValid, isValid)
			}
		})
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkCPF_Generate(b *testing.B) {
	cpf := NewCPF()
	for i := 0; i < b.N; i++ {
		cpf.Generate()
	}
}

func BenchmarkCPF_Validate(b *testing.B) {
	cpf := NewCPF()
	testCPF := "123.456.789-09"
	for i := 0; i < b.N; i++ {
		cpf.Validate(testCPF)
	}
}

func BenchmarkCNPJ_Generate(b *testing.B) {
	cnpj := NewCNPJ()
	for i := 0; i < b.N; i++ {
		_ = cnpj.Generate()
	}
}

func BenchmarkCNPJ_Validate(b *testing.B) {
	cnpj := NewCNPJ()
	testCNPJ := "12.ABC.345/01DE-35"
	for i := 0; i < b.N; i++ {
		cnpj.Validate(testCNPJ)
	}
}
