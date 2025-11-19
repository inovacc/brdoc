package brdoc

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	notAcceptedCPF []string
	rng            *rand.Rand
)

// Mapa de conversão para CNPJ alfanumérico (ASCII - 48)
var charToValue = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 17, 'B': 18, 'C': 19, 'D': 20, 'E': 21, 'F': 22, 'G': 23, 'H': 24, 'I': 25,
	'J': 26, 'K': 27, 'L': 28, 'M': 29, 'N': 30, 'O': 31, 'P': 32, 'Q': 33, 'R': 34,
	'S': 35, 'T': 36, 'U': 37, 'V': 38, 'W': 39, 'X': 40, 'Y': 41, 'Z': 42,
}

func init() {
	// Inicializa gerador de números aleatórios
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	// Inicializa CPFs não aceitos (todos dígitos iguais)
	notAcceptedCPF = make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		value := strings.Repeat(strconv.Itoa(i), 11)
		notAcceptedCPF = append(notAcceptedCPF, value)
	}
}

// ============================================================================
// CPF - Cadastro de Pessoas Físicas
// ============================================================================

// CPF represents a Brazilian individual tax ID validator
type CPF struct {
	cpfNumber []int
}

// CPFResponse contains CPF validation results
type CPFResponse struct {
	CPF     string
	IsValid bool
	Origin  string
}

// NewCPF creates a new CPF validator instance
func NewCPF() *CPF {
	return &CPF{}
}

// Generate generates a valid random CPF with formatting
func (c *CPF) Generate() string {
	number := make([]int, 9)
	for i := 0; i < 9; i++ {
		number[i] = rng.Intn(10)
	}

	number = append(number, c.calculateFirstDigit(number))
	number = append(number, c.calculateSecondDigit(number))
	return c.maskCPF(number)
}

// Validate validates a CPF number (with or without formatting)
func (c *CPF) Validate(values string) bool {
	c.clean(values)
	return c.isAccepted(values) && c.length(c.cpfNumber) && c.validate(c.cpfNumber)
}

// Format formats a CPF string to the standard format XXX.XXX.XXX-XX
func (c *CPF) Format(s string) (string, error) {
	c.clean(s)
	if !c.isAccepted(s) {
	}
	if len(c.cpfNumber) != 11 {
		return "", fmt.Errorf("CPF deve ter 11 dígitos, recebido: %d", len(c.cpfNumber))
	}

	return c.maskCPF(c.cpfNumber), nil
}

// CheckOrigin returns the Brazilian state/region where the CPF was issued
// based on the 9th digit
func (c *CPF) CheckOrigin(values string) string {
	c.clean(values)
	if len(c.cpfNumber) < 9 {
		return ""
	}

	switch c.cpfNumber[8] {
	case 0:
		return "Rio Grande do Sul"
	case 1:
		return "Distrito Federal, Goiás, Mato Grosso do Sul e Tocantins"
	case 2:
		return "Pará, Amazonas, Acre, Amapá, Rondônia e Roraima"
	case 3:
		return "Ceará, Maranhão e Piauí"
	case 4:
		return "Pernambuco, Rio Grande do Norte, Paraíba e Alagoas"
	case 5:
		return "Bahia e Sergipe"
	case 6:
		return "Minas Gerais"
	case 7:
		return "Rio de Janeiro e Espírito Santo"
	case 8:
		return "São Paulo"
	case 9:
		return "Paraná e Santa Catarina"
	default:
		return ""
	}
}

// Private CPF methods

func (c *CPF) maskCPF(values []int) string {
	cpf := ""
	for _, item := range values {
		cpf += strconv.Itoa(item)
	}
	cpf = strings.ReplaceAll(cpf, "-", "")
	return fmt.Sprintf("%s.%s.%s-%s", cpf[:3], cpf[3:6], cpf[6:9], cpf[9:])
}

func (c *CPF) clean(values string) {
	c.cpfNumber = nil
	values = strings.ReplaceAll(values, ".", "")
	values = strings.ReplaceAll(values, "-", "")
	for _, item := range values {
		digit, err := strconv.Atoi(string(item))
		if err == nil {
			c.cpfNumber = append(c.cpfNumber, digit)
		}
	}
}

func (c *CPF) calculateFirstDigit(values []int) int {
	sum := 0
	for i := 0; i < 9; i++ {
		sum += values[i] * (10 - i)
	}
	rest := (sum * 10) % 11
	if rest == 10 || rest == 11 {
		rest = 0
	}
	return rest
}

func (c *CPF) calculateSecondDigit(values []int) int {
	sum := 0
	for i := 0; i < 10; i++ {
		sum += values[i] * (11 - i)
	}
	rest := (sum * 10) % 11
	if rest == 10 || rest == 11 {
		rest = 0
	}
	return rest
}

func (c *CPF) validate(values []int) bool {
	if len(values) != 11 {
		return false
	}
	return c.calculateFirstDigit(values) == values[9] &&
		c.calculateSecondDigit(values) == values[10]
}

func (c *CPF) isAccepted(values string) bool {
	re := regexp.MustCompile(`[^0-9]`)
	cpf := re.ReplaceAllString(values, "")

	for _, notAccepted := range notAcceptedCPF {
		if cpf == notAccepted {
			return false
		}
	}

	return true
}

func (c *CPF) length(values []int) bool {
	return len(values) == 11
}

// ============================================================================
// CNPJ - Cadastro Nacional de Pessoa Jurídica (Alfanumérico)
// Baseado na especificação SERPRO
// ============================================================================

// CNPJ represents a Brazilian company tax ID validator (alphanumeric format)
type CNPJ struct{}

// CNPJResponse contains CNPJ validation results
type CNPJResponse struct {
	CNPJ      string
	IsValid   bool
	Formatted string
}

// NewCNPJ creates a new CNPJ validator instance
func NewCNPJ() *CNPJ {
	return &CNPJ{}
}

// Generate generates a valid alphanumeric CNPJ
func (cnpj *CNPJ) Generate() string {
	var sb strings.Builder

	// Gera os primeiros 12 caracteres aleatórios (números ou letras)
	for i := 0; i < 12; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte(byte('0' + rng.Intn(10))) // Número
		} else {
			sb.WriteByte(byte('A' + rng.Intn(26))) // Letra
		}
	}

	cnpjBase := sb.String()

	// Calcula os dois dígitos verificadores
	dv1, err := cnpj.calculateDV(cnpjBase)
	if err != nil {
		return ""
	}

	dv2, err := cnpj.calculateDV(cnpjBase + strconv.Itoa(dv1))
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s%d%d", cnpjBase, dv1, dv2)
}

// Validate verifies if an alphanumeric CNPJ is valid per SERPRO specification
func (cnpj *CNPJ) Validate(value string) bool {
	// Remove formatação
	cleaned := cnpj.clean(value)

	if len(cleaned) != 14 {
		return false
	}

	// Valida que os últimos 2 caracteres são numéricos
	dv1, err1 := strconv.Atoi(string(cleaned[12]))
	dv2, err2 := strconv.Atoi(string(cleaned[13]))
	if err1 != nil || err2 != nil {
		return false
	}

	base := cleaned[:12]
	dv1Calc, err := cnpj.calculateDV(base)
	if err != nil {
		return false
	}

	dv2Calc, err := cnpj.calculateDV(base + strconv.Itoa(dv1Calc))
	if err != nil {
		return false
	}

	return dv1Calc == dv1 && dv2Calc == dv2
}

// Format formats a CNPJ to the standard format XX.XXX.XXX/XXXX-XX
func (cnpj *CNPJ) Format(value string) (string, error) {
	cleaned := cnpj.clean(value)

	if len(cleaned) != 14 {
		return "", fmt.Errorf("CNPJ deve ter 14 caracteres, recebido: %d", len(cleaned))
	}

	return fmt.Sprintf("%s.%s.%s/%s-%s",
		cleaned[0:2],
		cleaned[2:5],
		cleaned[5:8],
		cleaned[8:12],
		cleaned[12:14],
	), nil
}

// Private CNPJ methods

// calculateDV calculates a check digit using modulo 11
// Official SERPRO algorithm for alphanumeric CNPJ
func (cnpj *CNPJ) calculateDV(value string) (int, error) {
	pesos := []int{2, 3, 4, 5, 6, 7, 8, 9}
	soma := 0
	j := 0

	// Percorre o CNPJ da direita para a esquerda aplicando os pesos
	for i := len(value) - 1; i >= 0; i-- {
		val, ok := charToValue[rune(value[i])]
		if !ok {
			return 0, fmt.Errorf("caractere inválido: %c na posição %d", value[i], i)
		}
		soma += val * pesos[j]
		j = (j + 1) % len(pesos) // Reinicia pesos após o 8º elemento
	}

	resto := soma % 11

	// Regra específica: se resto = 0 ou 1, DV = 0
	if resto == 0 || resto == 1 {
		return 0, nil
	}

	return 11 - resto, nil
}

// clean removes formatting from an alphanumeric CNPJ
func (cnpj *CNPJ) clean(value string) string {
	re := regexp.MustCompile(`[^0-9A-Z]`)

	return strings.ToUpper(re.ReplaceAllString(value, ""))
}

// ============================================================================
// Utility Functions
// ============================================================================

// ValidateDocument automatically identifies and validates CPF or CNPJ
func ValidateDocument(doc string) (docType string, isValid bool) {
	cleaned := strings.ReplaceAll(doc, ".", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "/", "")
	cleaned = strings.ToUpper(cleaned)

	// Identifica pelo tamanho
	if len(cleaned) == 11 {
		cpf := NewCPF()
		return "CPF", cpf.Validate(doc)
	} else if len(cleaned) == 14 {
		cnpj := NewCNPJ()
		return "CNPJ", cnpj.Validate(doc)
	}

	return "UNKNOWN", false
}
