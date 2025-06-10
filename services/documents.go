package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"io"
	"net/http"

	"time"

	"github.com/go-ping/ping"
	"github.com/hirochachacha/go-smb2"
)

type Response struct {
	Message string `json:"message"`
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	// Create JSON response
	response := Response{
		Message: "Document service endpoint",
	}

	// Set response headers and send JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func FetchPDFBase64(w http.ResponseWriter, r *http.Request) {
	// rutaSMB := r.URL.Query().Get("ruta")

	rutaSMB := "smb://SRVFACELE11/Timbrado3.3/Invoice One/InterfazIME/Configuracion/IES161108I36/XmlRecibido/2025-05/IES161108I36_QB_14166.pdf"

	if rutaSMB == "" {
		http.Error(w, "Falta el parámetro 'ruta'", http.StatusBadRequest)
		fmt.Println("Falta el parámetro 'ruta'")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Message": "Falta el parámetro 'ruta'",
		})
		return
	}

	rutaLimpia := strings.TrimPrefix(rutaSMB, "smb://")
	partes := strings.SplitN(rutaLimpia, "/", 2)
	if len(partes) != 2 {
		http.Error(w, "Ruta SMB no válida", http.StatusBadRequest)
		fmt.Println("Ruta SMB no válida")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Message": "Ruta SMB no válida",
		})
		return
	}
	host := partes[0]
	rutaArchivo := partes[1]

	rutaPartes := strings.SplitN(rutaArchivo, "/", 2)
	share := rutaPartes[0]
	pathInShare := rutaPartes[1]

	fmt.Println(host)
	fmt.Println(share)
	fmt.Println(pathInShare)

	fmt.Println("Conectando por SMB a: ", host+":445 ...")
	conn, err := net.Dial("tcp", host+":445")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error conectando SMB: %v", err), http.StatusInternalServerError)
		fmt.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Message": "Error conectando SMB",
			"Error":   err,
		})
		return
	}

	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Conectado a " + host + ":445 atraves de SMB"))
	fmt.Println("Conectado a " + host + ":445 atraves de SMB 🎉🎈")

	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     "", // agrega aquí si tienes usuario
			Password: "",
			Domain:   "",
		},
	}

	session, err := d.Dial(conn)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error iniciando sesión SMB: %v", err), http.StatusInternalServerError)
		fmt.Println("Error iniciando sesión SMB: ", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Message": "Error iniciando sesión SMB",
			"Error":   err,
		})
		return
	}
	defer session.Logoff()

	fs, err := session.Mount(share)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error montando recurso: %v", err), http.StatusInternalServerError)
		fmt.Println("Error montando recurso: ", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Message": "Error montando recurso",
			"Error":   err,
		})
		return
	}
	defer fs.Umount()

	file, err := fs.Open(pathInShare)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error abriendo archivo: %v", err), http.StatusInternalServerError)
		fmt.Println("Error abriendo archivo: ", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Message": "Error abriendo archivo",
			"Error":   err,
		})
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error leyendo archivo: %v", err), http.StatusInternalServerError)
		fmt.Println("Error leyendo archivo: ", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Message": "Error leyendo archivo",
			"Error":   err,
		})
		return
	}

	encoded := base64.StdEncoding.EncodeToString(content)

	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"base64":"%s"}`, encoded)

}

func MakePing(w http.ResponseWriter, r *http.Request) {
	pinger, err := ping.NewPinger("google.com")
	if err != nil {
		panic(err)
	}

	pinger.Count = 3
	pinger.Timeout = time.Second * 5
	pinger.SetPrivileged(false) // Cambiado a false para no requerir privilegios

	fmt.Println("Haciendo ping a:", pinger.Addr())

	err = pinger.Run() // Bloqueante
	if err != nil {
		fmt.Printf("Error al hacer ping: %v\n", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Message": "Ping Fallido",
			"Error":   err,
		})
		// Continuamos con la ejecución aunque falle el ping
	} else {
		stats := pinger.Statistics()
		fmt.Printf("Resultados del Ping: %d paquetes enviados, %d recibidos\n", stats.PacketsSent, stats.PacketsRecv)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Message":            "Resultados del Ping",
			"Paquetes_Enviados":  stats.PacketsSent,
			"Paquetes_Recibidos": stats.PacketsRecv,
		})
	}
}
