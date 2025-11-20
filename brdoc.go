package brdoc

import (
	"fmt"
	"math/rand"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	CpfLength  = 11
	CnpjLength = 14

	IsDigit0 = "Rio Grande do Sul"
	IsDigit1 = "Federal District, Goiás, Mato Grosso do Sul, and Tocantins"
	IsDigit2 = "Pará, Amazonas, Acre, Amapá, Rondônia, and Roraima"
	IsDigit3 = "Ceará, Maranhão, and Piauí"
	IsDigit4 = "Pernambuco, Rio Grande do Norte, Paraíba, and Alagoas"
	IsDigit5 = "Bahia and Sergipe"
	IsDigit6 = "Minas Gerais"
	IsDigit7 = "Rio de Janeiro and Espírito Santo"
	IsDigit8 = "São Paulo"
	IsDigit9 = "Paraná and Santa Catarina"
)

var (
	notAcceptedCPF []string
	rng            *rand.Rand
)

// Conversion map for alphanumeric CNPJ (ASCII - 48)
var charToValue = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 17, 'B': 18, 'C': 19, 'D': 20, 'E': 21, 'F': 22, 'G': 23, 'H': 24, 'I': 25,
	'J': 26, 'K': 27, 'L': 28, 'M': 29, 'N': 30, 'O': 31, 'P': 32, 'Q': 33, 'R': 34,
	'S': 35, 'T': 36, 'U': 37, 'V': 38, 'W': 39, 'X': 40, 'Y': 41, 'Z': 42,
}

func init() {
	// Initialize random number generator
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize non-accepted CPFs (all digits equal)
	notAcceptedCPF = make([]string, 0, 10)

	for i := range 10 {
		value := strings.Repeat(strconv.Itoa(i), 11)
		notAcceptedCPF = append(notAcceptedCPF, value)
	}
}

// ============================================================================
// CPF - Individual Taxpayer Registry
// ============================================================================

// CPF represents a Brazilian individual tax ID validator
type CPF struct {
	cpfNumber []int
}

// NewCPF creates a new CPF validator instance
func NewCPF() *CPF {
	return &CPF{}
}

// Generate generates a valid random CPF with unformatting
func (c *CPF) Generate() string {
	number := []int{0, 0, 0, 0, 0, 0, 0, 0, 0}

	for i := range 9 {
		number[i] = rng.Intn(10)
	}

	number = append(number, c.calculateFirstDigit(number))
	number = append(number, c.calculateSecondDigit(number))

	var sb strings.Builder

	for _, item := range number {
		sb.WriteString(strconv.Itoa(item))
	}

	return c.digits(sb.String())
}

// Validate validates a CPF number (with or without formatting)
func (c *CPF) Validate(value string) bool {
	c.clean(value)

	return c.isAccepted(value) && c.length(c.cpfNumber) && c.validate(c.cpfNumber)
}

// Format formats a CPF string to the standard format XXX.XXX.XXX-XX
func (c *CPF) Format(value string) (string, error) {
	c.clean(value)

	if !c.isAccepted(value) {
		return "", fmt.Errorf("CPF is not valid")
	}

	if len(c.cpfNumber) != CpfLength {
		return "", fmt.Errorf("CPF must have %d digits, got: %d", CpfLength, len(c.cpfNumber))
	}

	return c.maskCPF(c.cpfNumber), nil
}

// CheckOrigin returns the Brazilian state/region where the CPF was issued
// based on the 9th digit
func (c *CPF) CheckOrigin(value string) string {
	c.clean(value)

	if len(c.cpfNumber) < 9 {
		return ""
	}

	switch c.cpfNumber[8] {
	case 0:
		return IsDigit0
	case 1:
		return IsDigit1
	case 2:
		return IsDigit2
	case 3:
		return IsDigit3
	case 4:
		return IsDigit4
	case 5:
		return IsDigit5
	case 6:
		return IsDigit6
	case 7:
		return IsDigit7
	case 8:
		return IsDigit8
	case 9:
		return IsDigit9
	default:
		return ""
	}
}

// Private CPF methods

func (c *CPF) maskCPF(value []int) string {
	var sb strings.Builder

	for _, item := range value {
		sb.WriteString(strconv.Itoa(item))
	}

	cpf := c.digits(sb.String())

	return fmt.Sprintf("%s.%s.%s-%s", cpf[:3], cpf[3:6], cpf[6:9], cpf[9:])
}

func (c *CPF) clean(value string) {
	// Always reset and parse fresh to avoid stale state across calls
	c.cpfNumber = c.cpfNumber[:0]
	for _, item := range c.digits(value) {
		digit, err := strconv.Atoi(string(item))
		if err == nil {
			c.cpfNumber = append(c.cpfNumber, digit)
		}
	}
}

func (c *CPF) digits(value string) string {
	return regexp.MustCompile(`[^0-9]`).ReplaceAllString(value, "")
}

func (c *CPF) calculateFirstDigit(value []int) int {
	sum := 0
	for i, v := range value {
		sum += v * (10 - i)
	}

	rest := (sum * 10) % 11
	if rest == 10 || rest == 11 {
		rest = 0
	}

	return rest
}

func (c *CPF) calculateSecondDigit(value []int) int {
	sum := 0
	for i, v := range value {
		sum += v * (11 - i)
	}

	rest := (sum * 10) % 11
	if rest == 10 || rest == 11 {
		rest = 0
	}

	return rest
}

func (c *CPF) validate(value []int) bool {
	if len(value) != CpfLength {
		return false
	}
	// Calculate using base slices: first 9 for DV1, first 10 for DV2
	dv1 := c.calculateFirstDigit(value[:9])
	dv2 := c.calculateSecondDigit(append(value[:9], dv1))

	return dv1 == value[9] && dv2 == value[10]
}

func (c *CPF) isAccepted(value string) bool {
	// Reject CPFs with all equal digits
	return !slices.Contains(notAcceptedCPF, c.digits(value))
}

func (c *CPF) length(value []int) bool {
	return len(value) == CpfLength
}

// ============================================================================
// CNPJ - National Registry of Legal Entities (Alphanumeric)
// Based on the SERPRO specification
// ============================================================================

// CNPJ represents a Brazilian company tax ID validator (alphanumeric format)
type CNPJ struct{}

// NewCNPJ creates a new CNPJ validator instance
func NewCNPJ() *CNPJ {
	return &CNPJ{}
}

// Generate generates a valid alphanumeric CNPJ
func (c *CNPJ) Generate() string {
	var sb strings.Builder

	// Generate the first 12 random characters (numbers or letters)
	for range 12 {
		if rng.Intn(2) == 0 {
			sb.WriteByte(byte('0' + rng.Intn(10))) // Number
		} else {
			sb.WriteByte(byte('A' + rng.Intn(26))) // Letter
		}
	}

	cnpjBase := sb.String()

	// Calculate the two check digits
	dv1, err := c.calculateDV(cnpjBase)
	if err != nil {
		return ""
	}

	dv2, err := c.calculateDV(cnpjBase + strconv.Itoa(dv1))
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s%d%d", cnpjBase, dv1, dv2)
}

// Validate verifies if an alphanumeric CNPJ is valid per SERPRO specification
func (c *CNPJ) Validate(value string) bool {
	// Remove formatting
	cleaned := c.digits(value)

	if len(cleaned) != CnpjLength {
		return false
	}

	// Ensure the last 2 characters are numeric
	dv1, err1 := strconv.Atoi(string(cleaned[12]))
	if err1 != nil {
		return false
	}

	dv2, err2 := strconv.Atoi(string(cleaned[13]))
	if err2 != nil {
		return false
	}

	base := cleaned[:12]

	dv1Calc, err := c.calculateDV(base)
	if err != nil {
		return false
	}

	dv2Calc, err := c.calculateDV(base + strconv.Itoa(dv1Calc))
	if err != nil {
		return false
	}

	return dv1Calc == dv1 && dv2Calc == dv2
}

// Format formats a CNPJ to the standard format XX.XXX.XXX/XXXX-XX
func (c *CNPJ) Format(value string) (string, error) {
	cleaned := c.digits(value)

	if len(cleaned) != CnpjLength {
		return "", fmt.Errorf("CNPJ must have 14 characters, got: %d", len(cleaned))
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
func (c *CNPJ) calculateDV(value string) (int, error) {
	weights := []int{2, 3, 4, 5, 6, 7, 8, 9}
	sum := 0
	j := 0

	// Iterate the CNPJ from right to left applying the weights
	for i := len(value) - 1; i >= 0; i-- {
		val, ok := charToValue[rune(value[i])]
		if !ok {
			return 0, fmt.Errorf("invalid character: %c at position %d", value[i], i)
		}

		sum += val * weights[j]
		j = (j + 1) % len(weights) // Restart weights after the 8th element
	}

	remainder := sum % 11

	// Specific rule: if remainder = 0 or 1, DV = 0
	if remainder == 0 || remainder == 1 {
		return 0, nil
	}

	return 11 - remainder, nil
}

// clean removes formatting from an alphanumeric CNPJ
func (c *CNPJ) digits(value string) string {
	// Keep only digits and uppercase letters (A-Z); strip formatting like .-/ and spaces
	s := strings.ToUpper(value)
	return regexp.MustCompile(`[^0-9A-Z]`).ReplaceAllString(s, "")
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
	if len(cleaned) == CpfLength {
		cpf := NewCPF()
		return "CPF", cpf.Validate(doc)
	} else if len(cleaned) == CnpjLength {
		cnpj := NewCNPJ()
		return "CNPJ", cnpj.Validate(doc)
	}

	return "UNKNOWN", false
}
