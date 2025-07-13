package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AnalysisRequest struct {
	Code string `json:"code"`
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}

func main() {
	http.HandleFunc("/analyze", analyzeHandler)
	fmt.Println("[INFO] Servidor iniciado en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("[ERROR] Error al iniciar el servidor:", err)
	}
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[INFO] Nueva solicitud recibida:", r.Method, r.URL.Path)
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		fmt.Println("[ERROR] Método no permitido:", r.Method)
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req AnalysisRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("[ERROR] Error al decodificar la solicitud:", err)
		http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
		return
	}

	fmt.Println("[INFO] Código recibido:", req.Code)

	// Análisis léxico
	lexico := ResumenLexico(req.Code)

	// Análisis sintáctico
	arbol, erroresSintacticos := AnalisisSintactico(req.Code)
	sintactico := map[string]interface{}{
		"arbol":   arbol,
		"errores": erroresSintacticos,
	}

	// Análisis semántico
	semantico := AnalisisSemantico(req.Code)

	resp := map[string]interface{}{
		"lexico":     lexico,
		"sintactico": sintactico,
		"semantico":  semantico,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println("[ERROR] Error al codificar la respuesta:", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	fmt.Println("[INFO] Respuesta enviada exitosamente")
}
