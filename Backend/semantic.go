package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Variable struct {
	Name     string
	Type     string
	Declared bool
	Used     bool
}

func SemanticAnalysis(code string) string {
	lines := strings.Split(code, "\n")
	variables := map[string]*Variable{}
	errors := []string{}
	warnings := []string{}

	varDeclPattern := regexp.MustCompile(`^(int|float|double|char)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(=\s*([^;]+))?;`)
	assignPattern := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*([^;]+);`)
	whileCondPattern := regexp.MustCompile(`^while\s*\(([^)]+)\)`) // para validar condición de while

	// 1. Declaración de variables y duplicados
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		if m := varDeclPattern.FindStringSubmatch(trim); m != nil {
			typeName := m[1]
			varName := m[2]
			if _, exists := variables[varName]; exists {
				errors = append(errors, fmt.Sprintf("Línea %d: Variable '%s' ya declarada", i+1, varName))
			} else {
				variables[varName] = &Variable{Name: varName, Type: typeName, Declared: true, Used: false}
			}
			if m[4] != "" {
				// Validar tipo de asignación
				if !validateAssignmentType(typeName, m[4]) {
					errors = append(errors, fmt.Sprintf("Línea %d: Asignación incompatible para '%s'", i+1, varName))
				}
			}
			continue
		}
		// 2. Asignaciones a variables
		if m := assignPattern.FindStringSubmatch(trim); m != nil {
			varName := m[1]
			if v, exists := variables[varName]; exists {
				v.Used = true
				if !validateAssignmentType(v.Type, m[2]) {
					errors = append(errors, fmt.Sprintf("Línea %d: Asignación incompatible para '%s'", i+1, varName))
				}
			} else {
				errors = append(errors, fmt.Sprintf("Línea %d: Variable '%s' usada sin declarar", i+1, varName))
			}
			continue
		}
		// 3. Uso de variables en expresiones (básico)
		for name := range variables {
			if strings.Contains(trim, name) {
				variables[name].Used = true
			}
		}
		// 4. Validar condición de while
		if m := whileCondPattern.FindStringSubmatch(trim); m != nil {
			cond := m[1]
			if !validateCondition(cond, variables) {
				errors = append(errors, fmt.Sprintf("Línea %d: Condición de while inválida o usa variables no declaradas", i+1))
			}
		}
	}
	// 5. Variables declaradas y no usadas
	for _, v := range variables {
		if v.Declared && !v.Used {
			warnings = append(warnings, fmt.Sprintf("Variable '%s' declarada y no usada", v.Name))
		}
	}
	if len(errors) == 0 {
		res := "✅ Análisis semántico válido"
		if len(warnings) > 0 {
			res += "\nAdvertencias:\n- " + strings.Join(warnings, "\n- ")
		}
		return res
	}
	return "❌ Errores semánticos:\n- " + strings.Join(errors, "\n- ")
}

// Valida que la asignación sea compatible con el tipo (solo int y números por ahora)
func validateAssignmentType(typeName, expr string) bool {
	expr = strings.TrimSpace(expr)
	if typeName == "int" {
		return regexp.MustCompile(`^[0-9]+(\s*[+\-*/]\s*[a-zA-Z0-9_]+)*$`).MatchString(expr) || regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(expr)
	}
	// Puedes agregar más reglas para float, double, char...
	return true
}

// Valida que la condición del while use solo variables declaradas o números
func validateCondition(cond string, variables map[string]*Variable) bool {
	parts := regexp.MustCompile(`\W+`).Split(cond, -1)
	for _, p := range parts {
		if p == "" {
			continue
		}
		if regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(p) {
			if _, exists := variables[p]; !exists {
				return false
			}
		}
	}
	return true
}
