# Arquitectura del Sistema - PAK

PAK es una aplicación de escritorio diseñada como un "Mini StreamYard Local". La arquitectura se basa en el patrón de **Cliente Pesado / Servidor Ligero**, donde el procesamiento visual ocurre en el frontend y la gestión de procesos en el backend (Go).

## 1. Visión General de los Componentes

### Capa de Presentación (Frontend - Vue.js & PixiJS)
- **Motor de Renderizado:** Utiliza **PixiJS** (WebGL) para componer escenas en tiempo real sobre un `<canvas>`.
- **Captura de Medios:** El flujo de video final se obtiene mediante `canvas.captureStream()`.
- **Comunicación WebRTC:** Envía el flujo capturado al backend mediante una conexión WebRTC local (usando la librería Pion en el backend).
- **Interfaz de Usuario:** Construida con **Vue 3** y **Tailwind CSS**.

### Capa de Aplicación (Backend - Go)
- **Servidor Web Nativo:** Go sirve los archivos estáticos del frontend (`dist`) y gestiona las rutas API a través de `http.ServeMux`.
- **API REST & WebSocket:** Proporciona endpoints para el control de la aplicación (iniciar stream, grabación, túneles) y un servidor de señalización para invitados.
- **WebRTC Adapter:** Implementado con **Pion**. Recibe los tracks de video/audio del frontend para su procesamiento.
- **Orquestador de Procesos:** Gestiona el ciclo de vida de binarios externos (**FFmpeg**, **SSH** para túneles).

### Infraestructura y Herramientas Externas
- **FFmpeg:** Realiza la codificación final, grabación en disco (MP4) y transmisión (RTMP).
- **Localhost.run (SSH Tunnel):** Permite exponer el servidor de señalización local a internet para que invitados remotos se unan.

---

## 2. Flujo de Datos y Transferencia

### A. Flujo de Control (Comandos)
1. El usuario interactúa con la UI (ej. clic en "Iniciar Transmisión").
2. El frontend realiza peticiones **HTTP POST** o se comunica vía **WebSocket** con el backend en Go.
3. El backend recibe la petición y ejecuta la lógica de negocio o lanza procesos externos.

### B. Flujo de Video (Pipeline de Media)
La transferencia de video sigue este camino crítico:

1. **Generación:** El frontend renderiza la escena en un Canvas.
2. **Transferencia Local:** 
   - El frontend inicia una oferta WebRTC hacia el endpoint `/api/stream/offer`.
   - El backend (`app.go`) procesa la oferta y establece un **Peer Connection** utilizando Pion.
   - El video viaja del Canvas al Backend vía WebRTC.
3. **Distribución Interna (UDP/RTP):** 
   - El backend recibe los paquetes RTP del video.
   - Estos paquetes se reenvían mediante sockets UDP locales a puertos específicos:
     - `127.0.0.1:5004` (Para Transmisión en vivo).
     - `127.0.0.1:5005` (Para Grabación local).
4. **Consumo por FFmpeg:**
   - FFmpeg lee los paquetes RTP utilizando archivos `.sdp` temporales como descriptores de entrada.
   - **Codificación:** FFmpeg codifica el flujo (usualmente libx264).
   - **Salida:** Envía el flujo a un servidor RTMP (YouTube/Twitch) o lo guarda en `recordings/*.mp4`.

### C. Flujo de Invitados (Signaling & WebRTC)
1. **Acceso:** El backend inicia un túnel SSH que devuelve una URL pública dinámica.
2. **Señalización:** El invitado se conecta a la URL pública -> El tráfico se redirige al WebSocket local (`/ws`).
3. **Negociación:** Frontend del Studio e Invitado intercambian SDP/ICE a través del servidor de señalización.
4. **P2P:** Se establece la conexión WebRTC directa entre el Invitado y el Studio para el envío de cámara/micrófono.

---

## 3. Estructura de Directorios

- `/frontend`: Código fuente de la interfaz Vue.js.
  - `/src/components`: Elementos visuales (Canvas, Sidebar, etc.).
  - `/src/stores`: Estado global (Pinia).
- `/internal`: Lógica de negocio y adaptadores siguiendo principios de arquitectura limpia.
  - `/adapters`: Implementaciones técnicas (WebRTC con Pion).
- `/cmd`: Punto de entrada de la aplicación.
- `main.go`: Configuración del servidor HTTP, rutas y orquestación inicial.
- `app.go`: Lógica principal de la aplicación y gestión de FFmpeg.
- `signaling.go`: Servidor WebSocket para coordinación de invitados.

---

## 4. Tecnologías Clave

| Componente | Tecnología |
| :--- | :--- |
| **Framework UI** | Vue 3 |
| **Renderizado 2D/3D** | PixiJS (WebGL) |
| **Lenguaje Backend** | Go (Golang) |
| **Servidor Web** | Net/HTTP (Standard Library) |
| **Protocolo Media** | WebRTC (Pion) / RTP |
| **Procesamiento Video** | FFmpeg |
| **Redes / Túneles** | SSH (localhost.run) |
