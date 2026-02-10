package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gphr/internal/adapters"

	"github.com/pion/webrtc/v3"
)

// App struct
type App struct {
	ctx           context.Context
	ffmpegCmd     *exec.Cmd
	recordCmd     *exec.Cmd // Commando separado para grabación
	ffmpegIn      io.WriteCloser
	isStreaming   bool
	isRecording   bool
	tunnelCmd     *exec.Cmd
	webrtcAdapter *adapters.WebRTCAdapter
}

// NewApp creates a new App application struct
func NewApp() *App {
	adapter, _ := adapters.NewWebRTCAdapter()
	return &App{
		webrtcAdapter: adapter,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ToggleTunnel abre o cierra un túnel hacia internet usando localhost.run
func (a *App) ToggleTunnel() (string, error) {
	fmt.Println("ToggleTunnel: Iniciando llamada...")
	if a.tunnelCmd != nil && a.tunnelCmd.Process != nil {
		fmt.Println("ToggleTunnel: Cerrando túnel existente...")
		a.tunnelCmd.Process.Kill()
		a.tunnelCmd = nil
		return "", nil
	}

	fmt.Println("ToggleTunnel: Ejecutando commando SSH...")
	// -T deshabilita la asignación de terminal, -o para ignorar hosts desconocidos
	a.tunnelCmd = exec.Command("ssh", "-T", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-R", "80:localhost:8080", "nokey@localhost.run")

	stdout, err := a.tunnelCmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := a.tunnelCmd.Start(); err != nil {
		return "", err
	}

	// Leemos en un goroutine para no bloquear
	urlChan := make(chan string)
	go func() {
		buf := make([]byte, 4096)
		fullOutput := ""
		for {
			n, err := stdout.Read(buf)
			if err != nil || n == 0 {
				break
			}
			fullOutput += string(buf[:n])
			// Si ya tenemos la URL, no have falta leer más para la activación
			if 1 > 0 { // Placeholder para lógica de detección
				if len(fullOutput) > 50 && (contains(fullOutput, ".lhr.life") || contains(fullOutput, ".lhr.pro")) {
					break
				}
			}
		}
		urlChan <- fullOutput
	}()

	select {
	case output := <-urlChan:
		fmt.Printf("ToggleTunnel Output final: %s\n", output)
		return output, nil
	case <-time.After(8 * time.Second):
		return "timeout", nil
	}
}

// Helper simple para buscar strings
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))))
}

// GetLocalIP devuelve la IP privada del equipo para el link de invitado
func (a *App) GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// ProcessStudioOffer recibe la oferta SDP del frontend y devuelve la respuesta
func (a *App) ProcessStudioOffer(sdp string) (string, error) {
	fmt.Println("Recibiendo oferta WebRTC del Canvas...")

	return a.webrtcAdapter.ProcessOffer(sdp, func(track *webrtc.TrackRemote) {
		// Reenviamos a dos puertos: 5004 (Stream) y 5005 (Record)
		go func() {
			connStream, _ := net.Dial("udp", "127.0.0.1:5004")
			connRecord, _ := net.Dial("udp", "127.0.0.1:5005")
			defer func() {
				if connStream != nil { connStream.Close() }
				if connRecord != nil { connRecord.Close() }
			}()

			buf := make([]byte, 1500)
			for {
				n, _, err := track.Read(buf)
				if err != nil {
					return
				}
				if connStream != nil { connStream.Write(buf[:n]) }
				if connRecord != nil { connRecord.Write(buf[:n]) }
			}
		}()
	})
}

// StartRecording inicia la grabación local con FFmpeg
func (a *App) StartRecording(filename string) error {
	if a.isRecording {
		return fmt.Errorf("ya se está grabando")
	}

	// Asegurar carpetas
	os.MkdirAll("recordings", 0755)
	os.MkdirAll("logs", 0755)

	fullPath, _ := filepath.Abs(filepath.Join("recordings", filename))
	logPath, _ := filepath.Abs(filepath.Join("logs", "ffmpeg_record.log"))
	
	fmt.Printf("Backend: Grabando en %s\n", fullPath)

	// Crear archivo SDP temporal
	sdpContent := "v=0\no=- 0 0 IN IP4 127.0.0.1\ns=PAK\nc=IN IP4 127.0.0.1\nt=0 0\nm=video 5005 RTP/AVP 96\na=rtpmap:96 H264/90000"
	sdpPath := filepath.Join(os.TempDir(), "pak_record.sdp")
	os.WriteFile(sdpPath, []byte(sdpContent), 0644)

	// Pequeña espera para que el stream RTP se estabilice
	time.Sleep(1 * time.Second)

	args := []string{
		"-protocol_whitelist", "file,rtp,udp",
		"-i", sdpPath, 
		"-c:v", "libx264",            
		"-preset", "veryfast",       
		"-crf", "20",                 
		"-pix_fmt", "yuv420p",        
		"-s", "1920x1080",            
		"-y",                         
		fullPath,
	}

	a.recordCmd = exec.Command("ffmpeg", args...)
	
	// Redirigir errores a un log para inspección
	logFile, _ := os.Create(logPath)
	a.recordCmd.Stderr = logFile

	if err := a.recordCmd.Start(); err != nil {
		fmt.Printf("Error FFmpeg Start: %v\n", err)
		return err
	}

	a.isRecording = true
	return nil
}

// StopRecording detiene la grabación local
func (a *App) StopRecording() {
	if !a.isRecording {
		return
	}
	a.isRecording = false
	if a.recordCmd != nil && a.recordCmd.Process != nil {
		// Enviar señal de interrupción para cerrar el archivo MP4 correctamente
		a.recordCmd.Process.Signal(os.Interrupt)

		// Esperar un memento para que FFmpeg cierre el contenedor
		go func() {
			time.Sleep(1 * time.Second)
			if a.recordCmd.ProcessState == nil || !a.recordCmd.ProcessState.Exited() {
				a.recordCmd.Process.Kill()
			}
			fmt.Println("Grabación finalizada y guardada.")
		}()
	}
}

// StartStream inicia el proceso de FFmpeg escuchando en el puerto UDP
func (a *App) StartStream(rtmpURL string) error {
	if a.isStreaming {
		return fmt.Errorf("ya hay un stream activo")
	}

	os.MkdirAll("logs", 0755)
	logPath, _ := filepath.Abs(filepath.Join("logs", "ffmpeg_stream.log"))

	// Crear archivo SDP temporal para el stream
	sdpContent := "v=0\no=- 0 0 IN IP4 127.0.0.1\ns=PAKStream\nc=IN IP4 127.0.0.1\nt=0 0\nm=video 5004 RTP/AVP 96\na=rtpmap:96 H264/90000"
	sdpPath := filepath.Join(os.TempDir(), "pak_stream.sdp")
	os.WriteFile(sdpPath, []byte(sdpContent), 0644)

	// Esperar a que WebRTC envíe los primeros paquetes
	time.Sleep(1 * time.Second)

	// Commandos de FFmpeg para recibir RTP por UDP y emitir a RTMP
	args := []string{
		"-protocol_whitelist", "file,rtp,udp",
		"-i", sdpPath,
		"-c:v", "libx264",
		"-preset", "ultrafast", 
		"-tune", "zerolatency", 
		"-maxrate", "3000k",
		"-bufsize", "6000k",
		"-pix_fmt", "yuv420p",
		"-g", "60",
		"-threads", "0", 
		"-f", "flv",
		rtmpURL,
	}

	a.ffmpegCmd = exec.Command("ffmpeg", args...)
	
	logFile, _ := os.Create(logPath)
	a.ffmpegCmd.Stderr = logFile

	if err := a.ffmpegCmd.Start(); err != nil {
		return err
	}

	a.isStreaming = true
	fmt.Println("FFmpeg: Proceso iniciado hacia", rtmpURL)
	return nil
}

func (a *App) StopStream() {
	if !a.isStreaming {
		return
	}
	a.isStreaming = false
	if a.ffmpegCmd != nil && a.ffmpegCmd.Process != nil {
		a.ffmpegCmd.Process.Kill()
	}
	fmt.Println("FFmpeg: Proceso detenido")
}

// CheckDependencies checks if the required binaries are installed
func (a *App) CheckDependencies() map[string]bool {
	deps := []string{"ffmpeg"}
	results := make(map[string]bool)

	for _, dep := range deps {
		_, err := exec.LookPath(dep)
		results[dep] = err == nil
	}

	return results
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
