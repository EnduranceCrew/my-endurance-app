// Package validator fornece validação de CPF e e-mail sem dependências externas.
package validator

import (
	"regexp"
	"strconv"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail retorna true se o e-mail tiver formato válido.
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

// IsValidCPF valida o CPF usando o algoritmo oficial dos dígitos verificadores.
func IsValidCPF(cpf string) bool {
	// Remove formatação
	cpf = regexp.MustCompile(`\D`).ReplaceAllString(cpf, "")

	if len(cpf) != 11 {
		return false
	}

	// Rejeita sequências trivialmente inválidas
	if allSame(cpf) {
		return false
	}

	if !validDigit(cpf, 9) || !validDigit(cpf, 10) {
		return false
	}

	return true
}

// FormatCPF retorna o CPF no formato 000.000.000-00.
func FormatCPF(cpf string) string {
	cpf = regexp.MustCompile(`\D`).ReplaceAllString(cpf, "")
	if len(cpf) != 11 {
		return cpf
	}
	return cpf[:3] + "." + cpf[3:6] + "." + cpf[6:9] + "-" + cpf[9:]
}

// SanitizeCPF remove qualquer caractere não numérico.
func SanitizeCPF(cpf string) string {
	return regexp.MustCompile(`\D`).ReplaceAllString(cpf, "")
}

// ── helpers internos ─────────────────────────────────────────────────────────

func allSame(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

func validDigit(cpf string, pos int) bool {
	sum := 0
	weight := pos + 1
	for i := 0; i < pos; i++ {
		d, _ := strconv.Atoi(string(cpf[i]))
		sum += d * weight
		weight--
	}
	remainder := (sum * 10) % 11
	if remainder == 10 || remainder == 11 {
		remainder = 0
	}
	expected, _ := strconv.Atoi(string(cpf[pos]))
	return remainder == expected
}
