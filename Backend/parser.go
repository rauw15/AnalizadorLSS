package main

import (
	"strings"
)

func SyntaxAnalysis(code string) string {
	var errores []string

	if !strings.Contains(code, "function App()") {
		errores = append(errores, "No se definió la función principal 'App()'")
	}

	if strings.Contains(code, "return") && !strings.Contains(code, "{") {
		errores = append(errores, "La instrucción 'return' debe estar dentro de un bloque '{}'")
	}

	// Verifica uso básico de paréntesis y llaves balanceadas
	if !balanced(code, '(', ')') {
		errores = append(errores, "Paréntesis desbalanceados")
	}
	if !balanced(code, '{', '}') {
		errores = append(errores, "Llaves desbalanceadas")
	}

	if strings.Contains(code, "imprt") {
		errores = append(errores, "Uso incorrecto de la palabra clave 'import'")
	}
	if strings.Contains(code, "functin") {
		errores = append(errores, "Uso incorrecto de la palabra clave 'function'")
	}
	if strings.Contains(code, "defalt") {
		errores = append(errores, "Uso incorrecto de la palabra clave 'default'")
	}

	if len(errores) == 0 {
		return "Sintaxis válida"
	}
	return "Errores sintácticos:\n- " + strings.Join(errores, "\n- ")
}

func balanced(s string, open, close rune) bool {
	count := 0
	for _, r := range s {
		if r == open {
			count++
		} else if r == close {
			count--
			if count < 0 {
				return false
			}
		}
	}
	return count == 0
}
