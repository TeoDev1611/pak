package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"
	"time"
)

// App struct
type App struct {
	ctx         context.Context
	ffmpegCmd   *exec.Cmd
	ffmpegIn    io.WriteCloser
	isStreaming bool
	tunnelCmd   *exec.Cmd
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
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

	fmt.Println("ToggleTunnel: Ejecutando comando SSH...")
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
			// Si ya tenemos la URL, no hace falta leer más para la activación
			if (1 > 0) { // Placeholder para lógica de detección
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

// StartStream inicia el proceso de FFmpeg preparado para recibir video por stdin
func (a *App) StartStream(rtmpURL string) error {
	if a.isStreaming {
		return fmt.Errorf("ya hay un stream activo")
	}

	// Comandos de FFmpeg para recibir WebM/VP8 (de MediaRecorder) y convertir a RTMP/H264
	// Ajustamos para que lea de 'stdin' y emita a la URL proporcionada
	args := []string{
		"-i", "pipe:0",           // Lee de stdin
		"-c:v", "libx264",        // Convierte a H.264
		"-preset", "veryfast",    // Rapidez sobre compresión
		"-maxrate", "3000k",      // Bitrate máximo
		"-bufsize", "6000k",
		"-pix_fmt", "yuv420p",    // Formato compatible con YouTube/Twitch
		"-g", "60",               // Keyframes cada 2 segundos (a 30fps)
		"-f", "flv",              // Formato para RTMP
		rtmpURL,
	}

	a.ffmpegCmd = exec.Command("ffmpeg", args...)

	var err error
	a.ffmpegIn, err = a.ffmpegCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := a.ffmpegCmd.Start(); err != nil {
		return err
	}

	a.isStreaming = true
	fmt.Println("FFmpeg: Proceso iniciado hacia", rtmpURL)
	return nil
}

// PushVideoChunk recibe un trozo de video del frontend y lo escribe en FFmpeg
func (a *App) PushVideoChunk(data []byte) {
	if a.isStreaming && a.ffmpegIn != nil {
		_, err := a.ffmpegIn.Write(data)
		if err != nil {
			fmt.Println("Error escribiendo a FFmpeg:", err)
			a.StopStream()
		}
	}
}

func (a *App) StopStream() {
	if !a.isStreaming {
		return
	}
	a.isStreaming = false
	if a.ffmpegIn != nil {
		a.ffmpegIn.Close()
	}
	if a.ffmpegCmd != nil {
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
