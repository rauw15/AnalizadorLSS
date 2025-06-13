package main

import (
	"strings"
)

func SemanticAnalysis(code string) string {
	var errores []string

	if strings.Contains(code, `import React from "react"`) == false {
		errores = append(errores, `React debe importarse como: import React from "react"`)
	}

	if strings.Contains(code, "useStae") {
		errores = append(errores, "Error semántico: se quiso usar useState, pero está mal escrito como 'useStae'")
	}

	if strings.Contains(code, "ReactDOM.render") && !strings.Contains(code, "<App />") {
		errores = append(errores, "ReactDOM.render debería renderizar <App />")
	}

	if strings.Contains(code, "setCount") && !strings.Contains(code, "useState") {
		errores = append(errores, "setCount se usa sin declarar useState")
	}

	if len(errores) == 0 {
		return "Análisis semántico válido"
	}
	return "Errores semánticos:\n- " + strings.Join(errores, "\n- ")
}
