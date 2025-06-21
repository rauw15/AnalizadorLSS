package main

import (
	"fmt"
	"regexp"
	"strings"
)

func SyntaxAnalysis(code string) string {
	errors := []string{}
	lines := strings.Split(code, "\n")

	// Validar balance de llaves y paréntesis
	openBraces := strings.Count(code, "{")
	closeBraces := strings.Count(code, "}")
	if openBraces != closeBraces {
		errors = append(errors, "Error de sintaxis: Llaves {} no balanceadas")
	}
	openParens := strings.Count(code, "(")
	closeParens := strings.Count(code, ")")
	if openParens != closeParens {
		errors = append(errors, "Error de sintaxis: Paréntesis () no balanceados")
	}

	// Validar punto y coma al final de cada instrucción (excepto líneas con '{', '}', o vacías)
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasSuffix(trim, "{") || strings.HasSuffix(trim, "}") {
			continue
		}
		if !strings.HasSuffix(trim, ";") && !regexp.MustCompile(`^while\\s*\\(`).MatchString(trim) {
			errors = append(errors, "Línea "+itoa(i+1)+": Falta punto y coma al final de la instrucción")
		}
	}

	// Validar estructura do-while
	doWhilePattern := regexp.MustCompile(`do\\s*{[\\s\\S]*}\\s*while\\s*\\([\\s\\S]*\\)\\s*;`)
	if strings.Contains(code, "do") && !doWhilePattern.MatchString(code) {
		errors = append(errors, "Error de sintaxis: Estructura do-while incorrecta")
	}

	if len(errors) > 0 {
		return strings.Join(errors, "\n")
	}
	return "Sintaxis válida"
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
