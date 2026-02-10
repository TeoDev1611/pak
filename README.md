# GPHR Studio (Gopher Quick Studio)

**GPHR** es una estación de transmisión local diseñada para ser ligera, rápida y estética, construida con **Go** y **Vue 3**.

---

## 🧠 Cómo Funciona (Arquitectura)

### 1. El Servidor (Backend - Go)
- **Servidor HTTP:** Sirve los archivos estáticos del frontend (Vue) y la página del invitado.
- **Señalización WebRTC:** Usa WebSockets (`/ws`) para permitir que el Anfitrión y el Invitado intercambien sus ofertas (SDP) y candidatos (ICE). Es un broadcast simple: lo que recibe de uno, lo envía al otro.
- **Orquestador FFmpeg:** (En desarrollo) Recibirá los fragmentos de video del frontend vía WebSockets para enviarlos a una URL RTMP.

### 2. El Escenario (Frontend - Vue 3)
- **Vista Previa:** Actualmente usa etiquetas `<video>` nativas de HTML5 con `object-fit: cover`. Se migró de PixiJS a HTML nativo para intentar eliminar la latencia.
- **Sistema de Layouts:** Utiliza clases CSS reactivas y posicionamiento absoluto para acomodar las cámaras dentro de un contenedor que intenta mantener una relación de aspecto 16:9.
- **WebRTC:** Cada cliente captura su `localStream` y crea una `RTCPeerConnection` para enviar su señal y recibir la del otro.

---

## ❌ Lo que NO funciona (o tiene problemas)

### 1. El Ratio del Canvas (El "YouTube Feel")
- **Problema:** El contenedor 16:9 no siempre se centra o escala correctamente. En lugar de encogerse para caber en el hueco disponible entre los paneles, a veces se sale de los bordes o deja letterboxing (barras negras) desproporcionadas.
- **Meta:** Lograr que el rectángulo negro 16:9 sea el "jefe" y escale dinámicamente manteniendo su centro.

### 2. Latencia en Linux
- **Problema:** A pesar de usar video nativo, persiste un lag notable en algunos entornos.
- **Posible Causa:** Sincronización de frames en el driver de video de Linux o falta de aceleración de hardware activa en el navegador.

### 3. Sincronización WebRTC
- **Problema:** La conexión a veces es unidireccional (uno ve al otro pero no viceversa) si los tiempos de carga del WebSocket no son perfectos.
- **Mejora Necesaria:** Implementar un sistema de "re-intento" o "re-negociación" automático.

### 4. Flujo de Transmisión (Streaming)
- **Estado:** El botón "En Vivo" es puramente visual. Falta conectar el flujo del video (posiblemente volviendo a un canvas oculto) hacia el backend de Go para que FFmpeg procese la señal.

---

## 🛠 Desarrollo

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
