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
	styleWarnings := []string{}

	varDeclPattern := regexp.MustCompile(`^(int|String)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(=\s*([^;]+))?;`)
	assignPattern := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*([^;]+);`)
	ifCondPattern := regexp.MustCompile(`^if\s*\(([^)]+)\)\s*\{?`)
	printlnPattern := regexp.MustCompile(`System\.out\.println\s*\((.*)\);`)
	equalsPattern := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)\.equals\s*\((.*)\)`)

	// Advertencia: nombre de clase debe iniciar con mayúscula
	classPattern := regexp.MustCompile(`public\s+class\s+([a-zA-Z_][a-zA-Z0-9_]*)`)
	for _, line := range lines {
		if m := classPattern.FindStringSubmatch(line); m != nil {
			className := m[1]
			if len(className) > 0 && className[0] >= 'a' && className[0] <= 'z' {
				styleWarnings = append(styleWarnings, "Convención: El nombre de la clase debería iniciar con mayúscula ("+className+")")
			}
		}
	}
	// Advertencia: sugerir comentarios
	commentFound := false
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "//") || strings.HasPrefix(strings.TrimSpace(line), "/*") {
			commentFound = true
			break
		}
	}
	if !commentFound {
		styleWarnings = append(styleWarnings, "Sugerencia: Agrega comentarios para mejorar la documentación del código.")
	}
	// Advertencia: sugerir manejo de excepciones
	tryCatchFound := false
	for _, line := range lines {
		if strings.Contains(line, "try") && strings.Contains(line, "catch") {
			tryCatchFound = true
			break
		}
	}
	if !tryCatchFound {
		styleWarnings = append(styleWarnings, "Sugerencia: Considera manejar posibles excepciones con try-catch.")
	}

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
				if !validateAssignmentType(typeName, m[4], variables) {
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
				if !validateAssignmentType(v.Type, m[2], variables) {
					errors = append(errors, fmt.Sprintf("Línea %d: Asignación incompatible para '%s'", i+1, varName))
				}
			} else {
				errors = append(errors, fmt.Sprintf("Línea %d: Variable '%s' usada sin declarar", i+1, varName))
			}
			continue
		}
		// 3. Uso de variables en condiciones de if
		if m := ifCondPattern.FindStringSubmatch(trim); m != nil {
			cond := m[1]
			condNoStrings := removeStrings(cond)
			usedVars := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`).FindAllString(condNoStrings, -1)
			for _, v := range usedVars {
				if _, exists := variables[v]; !exists && v != "String" && v != "equals" {
					errors = append(errors, fmt.Sprintf("Línea %d: Variable '%s' usada sin declarar en condición", i+1, v))
				}
				if vv, exists := variables[v]; exists {
					vv.Used = true
				}
			}
			// Validar comparaciones tipo String == int, etc.
			if strings.Contains(condNoStrings, "==") || strings.Contains(condNoStrings, ">") || strings.Contains(condNoStrings, "<") {
				parts := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)\s*(==|!=|>|<|>=|<=)\s*([a-zA-Z_][a-zA-Z0-9_]*|[0-9]+)`).FindAllStringSubmatch(condNoStrings, -1)
				for _, p := range parts {
					left := p[1]
					right := p[3]
					if v, ok := variables[left]; ok && v.Type == "String" && (strings.Contains(condNoStrings, ">") || strings.Contains(condNoStrings, "<")) {
						errors = append(errors, fmt.Sprintf("Línea %d: Comparación inválida con String en condición", i+1))
					}
					if v, ok := variables[left]; ok && v.Type == "int" && strings.HasPrefix(right, "\"") {
						errors = append(errors, fmt.Sprintf("Línea %d: Comparación de int con string en condición", i+1))
					}
				}
			}
		}
		// 4. Validar uso correcto de .equals
		if m := equalsPattern.FindStringSubmatch(trim); m != nil {
			varName := m[1]
			if v, exists := variables[varName]; exists {
				if v.Type != "String" {
					errors = append(errors, fmt.Sprintf("Línea %d: .equals solo se debe usar con String", i+1))
				}
				v.Used = true
			} // No marcar error si no existe, ya lo hace la condición
		}
		// 5. Validar uso correcto de System.out.println
		if m := printlnPattern.FindStringSubmatch(trim); m != nil {
			arg := m[1]
			argNoStrings := removeStrings(arg)
			usedVars := regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`).FindAllString(argNoStrings, -1)
			for _, v := range usedVars {
				if _, exists := variables[v]; !exists && v != "String" {
					errors = append(errors, fmt.Sprintf("Línea %d: Variable '%s' usada sin declarar en println", i+1, v))
				}
				if vv, exists := variables[v]; exists {
					vv.Used = true
				}
			}
		}
	}
	// 6. Variables declaradas y no usadas
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
		if len(styleWarnings) > 0 {
			res += "\nSugerencias de estilo:\n- " + strings.Join(styleWarnings, "\n- ")
		}
		return res
	}
	return "❌ Errores semánticos:\n- " + strings.Join(errors, "\n- ") + "\nSugerencias de estilo:\n- " + strings.Join(styleWarnings, "\n- ")
}

// Valida que la asignación sea compatible con el tipo (solo int y String)
func validateAssignmentType(typeName, expr string, variables map[string]*Variable) bool {
	expr = strings.TrimSpace(expr)
	if typeName == "int" {
		if regexp.MustCompile(`^".*"$`).MatchString(expr) {
			return false
		}
		return regexp.MustCompile(`^[0-9]+$`).MatchString(expr) || (regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(expr) && variables[expr] != nil && variables[expr].Type == "int")
	}
	if typeName == "String" {
		return regexp.MustCompile(`^".*"$`).MatchString(expr) || (regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(expr) && variables[expr] != nil && variables[expr].Type == "String")
	}
	return true
}

func removeStrings(line string) string {
	// Elimina los literales de texto entre comillas dobles
	return regexp.MustCompile(`"[^"]*"`).ReplaceAllString(line, "")
}
