package main

import (
	"fmt"
	"regexp"
	"strings"
)

func SyntaxAnalysis(code string) string {
	errors := []string{}
	lines := strings.Split(code, "\n")

	// 1. Validar declaración de clase
	classPattern := regexp.MustCompile(`^\s*public\s+class\s+[A-Z][a-zA-Z0-9_]*\s*\{`)
	classFound := false
	for _, line := range lines {
		if classPattern.MatchString(strings.TrimSpace(line)) {
			classFound = true
			break
		}
	}
	if !classFound {
		errors = append(errors, "No se encontró declaración de clase válida (public class NombreClase {...})")
	}

	// 2. Validar método main
	mainPattern := regexp.MustCompile(`^\s*public\s+static\s+void\s+main\s*\(String\s*\[\]\s*[a-zA-Z_][a-zA-Z0-9_]*\)\s*\{`)
	mainFound := false
	for _, line := range lines {
		if mainPattern.MatchString(strings.TrimSpace(line)) {
			mainFound = true
			break
		}
	}
	if !mainFound {
		errors = append(errors, "No se encontró método main válido (public static void main(String[] args) {...})")
	}

	// 3. Validar balance de llaves y paréntesis
	openBraces := strings.Count(code, "{")
	closeBraces := strings.Count(code, "}")
	if openBraces != closeBraces {
		errors = append(errors, "Llaves { } no balanceadas")
	}
	openParens := strings.Count(code, "(")
	closeParens := strings.Count(code, ")")
	if openParens != closeParens {
		errors = append(errors, "Paréntesis ( ) no balanceados")
	}

	// 4. Validar punto y coma al final de declaraciones y sentencias
	stmtPattern := regexp.MustCompile(`^(int|String)\s+[a-zA-Z_][a-zA-Z0-9_]*\s*(=.+)?;$|^System\.out\.println\s*\(.+\);$|^if\s*\(.+\)\s*\{?$|^\}$`)
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		if !stmtPattern.MatchString(trim) && !strings.HasSuffix(trim, "{") && !strings.HasSuffix(trim, "}") && !strings.HasPrefix(trim, "//") {
			errors = append(errors, fmt.Sprintf("Línea %d: Sintaxis inválida o falta punto y coma", i+1))
		}
	}

	// 5. Validar estructura de if
	ifPattern := regexp.MustCompile(`^if\s*\(.+\)\s*\{?$`)
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "if") && !ifPattern.MatchString(trim) {
			errors = append(errors, fmt.Sprintf("Línea %d: Estructura de if inválida", i+1))
		}
	}

	// 6. Validar llamadas a métodos
	printlnPattern := regexp.MustCompile(`System\.out\.println\s*\(.+\);`)
	equalsPattern := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*\.equals\s*\(.+\)`)
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.Contains(trim, ".println") && !printlnPattern.MatchString(trim) {
			errors = append(errors, fmt.Sprintf("Línea %d: Llamada a System.out.println inválida", i+1))
		}
		if strings.Contains(trim, ".equals") && !equalsPattern.MatchString(trim) {
			errors = append(errors, fmt.Sprintf("Línea %d: Llamada a .equals inválida", i+1))
		}
	}

	// 7. Validar comillas en strings
	stringPattern := regexp.MustCompile(`"[^"]*"`)
	for i, line := range lines {
		if strings.Contains(line, "\"") && len(stringPattern.FindAllString(line, -1)) == 0 {
			errors = append(errors, fmt.Sprintf("Línea %d: Error de comillas en string", i+1))
		}
	}

	// 8. Detectar condiciones fuera de if o con palabra mal escrita
	condPattern := regexp.MustCompile(`^\s*([a-zA-Z_][a-zA-Z0-9_]*)?\s*\([^\)]*\)\s*\{`)
	for i, line := range lines {
		if condPattern.MatchString(line) {
			m := condPattern.FindStringSubmatch(line)
			word := m[1]
			if word == "" {
				// Caso: solo (condición) { sin palabra antes
				errors = append(errors, fmt.Sprintf("Línea %d: Condición sin 'if' detectada, falta palabra clave 'if'", i+1))
			} else if word != "if" {
				// Caso: palabra mal escrita antes de la condición
				errors = append(errors, fmt.Sprintf("Línea %d: Palabra clave desconocida '%s' antes de condición, ¿quizá quisiste escribir 'if'?", i+1, word))
			}
		}
	}

	if len(errors) > 0 {
		return strings.Join(errors, "\n")
	}
	return "Sintaxis válida"
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
