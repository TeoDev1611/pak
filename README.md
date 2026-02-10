# CHQS Studio (Chasqui Quick Studio)

**CHQS** es una estación de transmisión local diseñada para ser ligera, rápida y estética, construida con **Go** y **Vue 3**.

## 🚀 Estado Actual
El proyecto utiliza un servidor de backend en Go para la señalización WebRTC y la orquestación de FFmpeg, y un frontend en Vue para la interfaz de control.

### Lo que funciona:
- [x] **Arquitectura:** Go + Vue 3 (Vite) + Tailwind CSS.
- [x] **Lobby de Entrada:** Control de nombre y permisos de hardware.
- [x] **WebRTC:** Conexión bidireccional entre Anfitrión e Invitado.
- [x] **Túneles:** Generación de links de invitado vía `localhost.run`.
- [x] **Layouts:** Solo, Grid (Pegado) y Zoom (Espaciado) con escalado 16:9 automático.

### 🛠 Pendientes:
1.  **Ajuste de Ratio:** Centrar el lienzo 16:9 sin letterboxing desproporcionado.
2.  **Motor de Composición:** Implementar canvas oculto para la captura de streaming.
3.  **Compartir Pantalla:** Integrar la fuente de pantalla en los layouts.
4.  **Streaming:** Conexión final con FFmpeg.

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
El estudio se abrirá automáticamente en tu navegador en `http://localhost:8080`.

## 📡 Señalización
El servidor de señales corre internamente en `/ws` para gestionar el intercambio de SDP y Candidatos entre los participantes.
