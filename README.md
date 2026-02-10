# PAK PRO STUDIO 🚀

**PAK Pro Studio** es una estación de transmisión local de nivel profesional diseñada para ser ligera, rápida y estética, construida con **Go**, **Vue 3** y **PixiJS v8**.

---

## 💎 Estado Actual del Proyecto

La aplicación ha evolucionado de un prototipo básico a un estudio funcional con las siguientes capacidades:

### 🎨 Motor Visual (Frontend)
- **Resolución Nativa Full HD (1080p):** El lienzo de composición opera internamente a 1920x1080.
- **Composición por Capas:** Sistema robusto de 3 capas (Fondo, Video, Overlays) gestionado por PixiJS v8 con aceleración de hardware.
- **Layouts Profesionales:** 4 modos inteligentes (`Solo`, `Grid`, `Sidebar`, `PiP`) con bordes suavizados y márgenes cinematográficos.
- **Identidad Visual (Branding):**
    - Subida de logos personalizados con escalado dinámico.
    - Fondos de estudio personalizables.
    - **Banners y Tickers:** Sistema de persistencia (localStorage) para banners, incluyendo el modo "Cinta de Noticias" (Marquee) con bucle infinito.
- **UI Premium:** Interfaz oscura ("PAK Pro") con estilo glassmorphism y controles flotantes.

### ⚙️ Backend y Grabación (Go)
- **Orquestación con FFmpeg:** El backend gestiona procesos independientes para streaming (RTMP) y grabación local.
- **Grabación en el Servidor:** La captura de video se realiza en el backend vía WebRTC -> RTP -> FFmpeg, permitiendo grabaciones en MP4 Full HD con alta fidelidad (`crf 20`) sin sobrecargar el navegador.
- **Señalización WebRTC:** Servidor de señalización basado en WebSockets para conexión de invitados remotos.
- **Túneles Seguros:** Integración con `localhost.run` para permitir acceso a invitados externos vía HTTPS.

---

## ❌ Problemas Conocidos (Para arreglar próximamente)

### 1. Escalado del Layout de Vista Previa
- **Problema:** A pesar de ser Full HD, el contenedor de la vista previa a veces se visualiza más pequeño de lo esperado o no se expande para llenar todo el espacio central entre la barra lateral y el footer.
- **Estado:** Se han aplicado correcciones de CSS (`flex-1`, `w-full`), pero persiste una inconsistencia en el centrado vertical en algunas resoluciones de pantalla.

### 2. Persistencia de Archivos de Grabación
- **Problema:** Algunos procesos de grabación de FFmpeg no finalizan correctamente la escritura del archivo MP4 en la carpeta `recordings/`, o el stream RTP no se inicia a tiempo para ser capturado.
- **Estado:** Se han añadido logs en `/logs/ffmpeg_record.log` y archivos SDP temporales para debuguear la comunicación RTP entre el navegador y Go.

### 3. Audio en la Grabación
- **Problema:** Actualmente la grabación en el backend se centra en el stream de video del canvas.
- **Meta:** Implementar un mezclador de audio en el backend para capturar tanto al anfitrión como a los invitados en el archivo final.

### 4. OverconstrainedError (Cámara)
- **Problema:** En ciertos hardware, las restricciones de cámara muy estrictas bloquean el inicio del estudio.
- **Solución actual:** Se implementaron "fallbacks" automáticos de Full HD -> HD -> Básico para garantizar que la cámara siempre encienda.

---

## 🛠 Desarrollo y Uso

### 1. Preparar el Frontend
```bash
cd frontend
npm install
npm run build
```

### 2. Ejecutar el Servidor (Backend)
Desde la raíz del proyecto:
```bash
go run .
```
El estudio se abrirá automáticamente en `http://localhost:8080`.

### 📂 Estructura de Archivos Clave
- `recordings/`: Aquí se guardan los archivos MP4 generados.
- `logs/`: Logs de diagnóstico de FFmpeg para streaming y grabación.
- `frontend/src/stores/studio.js`: Estado global (Banners, Logos, Configuración).
- `app.go`: Lógica principal del backend y gestión de procesos FFmpeg.