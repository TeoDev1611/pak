Analizando la arquitectura de gigantes como **StreamYard** y **Google Meet**, y
contrastándolo con tu stack (Go + Vue 3), aquí tienes el "secreto" de cómo
logran estabilidad y qué paquetes deberías integrar en **GPHR** para dar ese
salto de calidad.

### 1. El modelo "StreamYard": Composición en el Cliente (Client-Side Compositing)

StreamYard no envía 4 videos separados al servidor para que este los una (eso
consume muchísima CPU en el backend). Lo que hacen es **unir todo en el
navegador del Host**.

- **Cómo funciona:**

1. Crean un elemento `<canvas>` invisible en el DOM.
2. Dibujan el fondo, las cámaras de los invitados y los textos (overlays) en ese
   canvas 30 o 60 veces por segundo.
3. Usan `canvas.captureStream()` para convertir ese dibujo en un `MediaStream`.
4. Ese **único stream** ya mezclado es el que se envía al servidor para
   transmitir a YouTube/Twitch.

**Para GPHR (Vue 3):**

- **No uses solo CSS para posicionar.** CSS es visual para el usuario, pero para
  el _stream_, usa el API de **Canvas**.
- **Paquete recomendado:** Si el Canvas nativo te resulta complejo de manejar
  (coordenadas X, Y), usa **Konva.js** o **Fabric.js** (wrappers de Canvas muy
  potentes). Vue tiene integración con Konva (`vue-konva`).
- _Ventaja:_ Lo que ves en el Canvas es EXACTAMENTE lo que sale en vivo. Se
  acabó el problema de "se ve bien en mi pantalla pero mal en el stream".

---

### 2. El modelo "Google Meet": La potencia de WebRTC en el Backend

Google Meet usa una arquitectura SFU (Selective Forwarding Unit). No mezclan el
video, solo lo reenvían inteligentemente. Para tu backend en **Go**, estás
usando WebSockets para señalización, pero para manejar el video real (si decides
procesarlo en el servidor) necesitas algo más robusto que un pipe a FFmpeg.

**El paquete estándar de oro en Go: `pion/webrtc**` Si estás haciendo WebRTC en
Go, **Pion** es obligatorio. Es la librería que usan la mayoría de proyectos
serious de Go.

- **¿Por qué Pion?**
- No depende de CGO (es Go puro, compila rápido y fácil).
- Te permite acceso a bajo nivel a los paquetes RTP (para corregir la latencia
  en Linux).
- Maneja mejor el "Jitter Buffer" (el culpable de que el video se acelere o se
  congele).

**Mejora para GPHR:** En lugar de `Browser -> WebSocket -> Go -> FFmpeg`, la
ruta professional es:

1. **Browser:** Envía video vía WebRTC (usando `RTCPeerConnection`) al servidor
   Go.
2. **Go (Pion):** Recibe el track de video RTP.
3. **Go (Pion):** Reenvía esos paquetes RTP a un proceso FFmpeg local o escribe
   en un archivo.

---

### 3. Paquetes y Tecnologías Específicas a Incorporar

Aquí está la "lista de compras" técnica basada en lo que usan los grandes:

#### A. Backend (Go)

1. **`github.com/pion/webrtc/v3`**:

- _Uso:_ Para recibir el video en el servidor de forma nativa, no como "datos
  crudos" por WebSocket, sino como stream de video real. Esto arregla tus
  problemas de sincronización.

2. **`github.com/pion/rtp` & `github.com/pion/interceptor**`:

- _Uso:_ Para monitorear la calidad de la red (RTCP) y pedir al navegador que
  baje la calidad si la red es mala (como have Meet cuando te ves pixelado).

#### B. Frontend (Vue 3)

1. **`@vueuse/core`**:

- _Uso:_ Tiene hooks maravillosos como `useUserMedia` y `useDisplayMedia` que
  simplifican la gestión de permisos y errores de cámara/micrófono.

2. **`pinia` (Gestión de Estado)**:

- _Uso:_ StreamYard maneja el estado de "quién está en pantalla" de forma
  global. Necesitas saber si el `guest_1` está `active`, `muted` o `hidden`.
  Pinia es perfecto para esto.

3. **`adapter.js` (webrtc-adapter)**:

- _Uso:_ Es un polyfill official de Google. Normaliza las diferencias de WebRTC
  entre Chrome, Firefox y Safari. Esto solucionará muchos de tus bugs de "uno ve
  al otro pero no vice-versa".

---

### 4. Soluciones a tus Problemas Específicos

#### ❌ El problema del Ratio y "YouTube Feel"

- **Solución:** No pelees con CSS puro.
- **Implementación:** Define un contenedor con `aspect-ratio: 16/9`. Dentro, usa
  CSS Grid.

```css
.stage {
  aspect-ratio: 16 / 9;
  background: #000;
  display: grid;
  /* Si hay 2 personas, divide en 2 columnas */
  grid-template-columns: repeat(auto-fit, minmax(0, 1fr));
}
```

#### ❌ Latencia en Linux

- **Causa probable:** El navegador en Linux a veces no usa decodificación por
  hardware para WebRTC por defecto.
- **Mejora:** Al crear la oferta SDP en el frontend, puedes manipular el string
  para preferir el códec **VP8** sobre H.264 (o vice-versa, dependiendo de qué
  soporte mejor tu distro de Linux). `pion` te permite negociar esto
  automáticamente.

#### ❌ Sincronización (Audio antes que video)

- **Solución:** Google Meet usa **Opus** para audio y lo prioriza. Asegúrate de
  que tu `MediaStream` tenga sincronización de pistas. Si usas el método de
  **Canvas Compositing** (StreamYard), la sincronización es automática porque el
  audio y el video salen del mismo "lugar" (el mezclador del navegador) hacia el
  servidor.

### ¿Cuál es el siguiente paso lógico para GPHR?

Te sugiero **cambiar la estrategia de Streaming**: No intentes mezclar el video
en el servidor con Go + FFmpeg (es muy difícil hacerlo bien). Implementa la
**Composición en Canvas** en el frontend (Vue).

1. Dibuja los videos en un Canvas.
2. Captura el stream del Canvas.
3. Envía ESE stream único al backend.

¿Te interesa ver un ejemplo de cómo capturar un Canvas en Vue 3 y enviarlo?
