package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	rtpSignal     chan bool // Canal para avisar que hay datos fluyendo
	videoTrack    *webrtc.TrackRemote
	hasRTP        bool // Estado de si estamos recibiendo RTP
}

// NewApp creates a new App application struct
func NewApp() *App {
	adapter, _ := adapters.NewWebRTCAdapter()
	return &App{
		webrtcAdapter: adapter,
		rtpSignal:     make(chan bool, 10), // Aumentamos buffer
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
		a.videoTrack = track
		// Reenviamos a dos puertos: 5004 (Stream) y 5005 (Record)
		go func() {
			connStream, _ := net.Dial("udp", "127.0.0.1:5004")
			connRecord, _ := net.Dial("udp", "127.0.0.1:5005")
			defer func() {
				if connStream != nil {
					connStream.Close()
				}
				if connRecord != nil {
					connRecord.Close()
				}
			}()

			buf := make([]byte, 1500)
			for {
				n, _, err := track.Read(buf)
				if err != nil {
					a.hasRTP = false
					return
				}

				if n > 0 {
					if !a.hasRTP {
						a.hasRTP = true
						fmt.Println("Backend: ¡Primer paquete RTP recibido!")
					}
					// Notificamos que hay datos
					select {
					case a.rtpSignal <- true:
					default:
					}

					if connStream != nil {
						connStream.Write(buf[:n])
					}
					if connRecord != nil {
						connRecord.Write(buf[:n])
					}
				}
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
	os.MkdirAll("recordings", 0o755)
	os.MkdirAll("logs", 0o755)

	// Añadir timestamp para evitar sobrescritura
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	if base == "" {
		base = "grabacion"
	}
	finalName := fmt.Sprintf("%s_%s%s", base, timestamp, ".mp4")

	fullPath, _ := filepath.Abs(filepath.Join("recordings", finalName))
	logPath, _ := filepath.Abs(filepath.Join("logs", "ffmpeg_record.log"))

	fmt.Printf("Backend: Iniciando grabación en %s\n", fullPath)

	// 1. SOLICITAR KEYFRAME (PLI) INMEDIATAMENTE
	if a.videoTrack != nil {
		fmt.Println("Backend: Solicitando Keyframe (PLI) para sincronizar grabación...")
		a.webrtcAdapter.RequestKeyframe(uint32(a.videoTrack.SSRC()))
	}

	// Crear archivo SDP temporal simplificado
	sdpContent := "v=0\no=- 0 0 IN IP4 127.0.0.1\ns=PAK\nc=IN IP4 127.0.0.1\nt=0 0\nm=video 5005 RTP/AVP 96\na=rtpmap:96 H264/90000"
	sdpPath := filepath.Join(os.TempDir(), "pak_record.sdp")
	os.WriteFile(sdpPath, []byte(sdpContent), 0o644)

	// 2. ARGUMENTOS LIMPIOS (Sin -s que falla en SDP demuxer)
	args := []string{
		"-y",
		"-protocol_whitelist", "file,udp,rtp",
		"-f", "sdp",
		"-i", sdpPath,
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-crf", "22",
		"-pix_fmt", "yuv420p",
		"-movflags", "+faststart",
		fullPath,
	}

	a.recordCmd = exec.Command("ffmpeg", args...)
	stderr, _ := a.recordCmd.StderrPipe()

	// 3. LÓGICA DE ESPERA (PERMISIVA)
	go func() {
		fmt.Println("Backend: Preparando flujo para FFmpeg...")
		
		// Intentamos pedir keyframe
		if a.videoTrack != nil {
			a.webrtcAdapter.RequestKeyframe(uint32(a.videoTrack.SSRC()))
		}

		// Espera mínima para que el PLI surta efecto, pero sin abortar
		select {
		case <-a.rtpSignal:
			fmt.Println("Backend: ¡Datos RTP detectados!")
		case <-time.After(1500 * time.Millisecond):
			fmt.Println("Backend: Iniciando FFmpeg preventivamente (sin esperar señal)...")
		}

		if err := a.recordCmd.Start(); err != nil {
			fmt.Printf("Error crítico FFmpeg Start: %v\n", err)
			return
		}
		a.isRecording = true

		// Logs y Wait
		go func() {
			slurp, _ := io.ReadAll(stderr)
			if len(slurp) > 0 {
				os.WriteFile(logPath, slurp, 0o644)
			}
		}()

		err := a.recordCmd.Wait()
		a.isRecording = false

		// Verificación final
		time.Sleep(500 * time.Millisecond)
		if info, statErr := os.Stat(fullPath); statErr == nil {
			fmt.Printf("✅ ÉXITO: Grabación finalizada. %s (%.2f MB)\n", info.Name(), float64(info.Size())/1024/1024)
		} else {
			fmt.Printf("❌ ERROR: El archivo no se generó.\n")
		}

		if err != nil {
			fmt.Printf("Backend: FFmpeg terminó con error: %v\n", err)
		}
	}()

	return nil
}

// StopRecording detiene la grabación local enviando señal de interrupción
func (a *App) StopRecording() {
	if !a.isRecording {
		return
	}

	if a.recordCmd != nil && a.recordCmd.Process != nil {
		fmt.Println("Backend: Enviando señal Interrupt a FFmpeg...")
		// Enviar señal de interrupción para cerrar el archivo MP4 correctamente (átomo 'moov')
		a.recordCmd.Process.Signal(os.Interrupt)

		// Damos un tiempo a que FFmpeg cierre elegantemente antes de forzar
		go func() {
			time.Sleep(2 * time.Second)
			if a.isRecording {
				fmt.Println("Backend: FFmpeg no cerró a tiempo, forzando Kill...")
				a.recordCmd.Process.Kill()
				a.isRecording = false
			}
		}()
	}
}

// StartStream inicia el proceso de FFmpeg escuchando en el puerto UDP
func (a *App) StartStream(rtmpURL string) error {
	if a.isStreaming {
		return fmt.Errorf("ya hay un stream activo")
	}

	os.MkdirAll("logs", 0o755)
	logPath, _ := filepath.Abs(filepath.Join("logs", "ffmpeg_stream.log"))

	// Crear archivo SDP temporal para el stream
	sdpContent := "v=0\no=- 0 0 IN IP4 127.0.0.1\ns=PAKStream\nc=IN IP4 127.0.0.1\nt=0 0\nm=video 5004 RTP/AVP 96\na=rtpmap:96 H264/90000"
	sdpPath := filepath.Join(os.TempDir(), "pak_stream.sdp")
	os.WriteFile(sdpPath, []byte(sdpContent), 0o644)

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
