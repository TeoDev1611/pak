# CHQS Studio (Chasqui Quick Studio)

**CHQS** es una estación de transmisión local "Mini StreamYard" diseñada para ser ligera, rápida y estética.

## 🚀 Estado Actual
El proyecto se encuentra en **Fase de Estabilización de UI**. Se ha migrado de un motor de renderizado de Canvas (PixiJS) a una arquitectura de **Vista Previa HTML5 Nativa** para eliminar la latencia en entornos Linux.

### Lo que funciona:
- [x] **Lobby de Entrada:** Control de nombre y permisos de hardware.
- [x] **WebRTC:** Conexión bidireccional entre Anfitrión e Invitado.
- [x] **Túneles:** Generación de links de invitado vía `localhost.run`.
- [x] **Layouts:** Solo, Grid (Pegado) y Zoom (Espaciado) con escalado 16:9 automático.
- [x] **Estética:** Dark mode profesional con acentos Naranja/Teal.

### 🛠 Pendientes (Para Mañana):
1.  **Arreglar el Ratio del Canvas:** Asegurar que el contenedor 16:9 se centre perfectamente sin barras negras laterales (Letterboxing dinámico).
2.  **Motor de Composición Oculto:** Re-implementar el Canvas (oculto) para capturar la imagen final para streaming.
3.  **Compartir Pantalla:** Integrar el stream de pantalla como una fuente adicional en los layouts.
4.  **Integración FFmpeg:** Conectar el canvas capturado con el binario de FFmpeg en el backend.

## 🛠 Desarrollo
Para ejecutar en modo desarrollo:
```bash
wails dev
```

El servidor de señalización corre internamente en el puerto `8080`.