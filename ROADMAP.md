# ROADMAP DE RESCATE Y ESTABILIZACIÓN: GPHR Studio

Este documento define la nueva estrategia técnica para eliminar el lag, arreglar los layouts y profesionalizar la aplicación, separando la **Vista Previa** (lo que ve el usuario) del **Motor de Transmisión** (lo que se envía a internet).

## 🚨 Fase 1: Estabilización del Núcleo (PRIORIDAD URGENTE)
**Objetivo:** Eliminar el lag de la cámara local y del invitado.
*   [ ] **Refactor de UI:** Eliminar PixiJS de la capa de visualización (`VideoCanvas.vue`).
*   [ ] **Motor Nativo:** Usar etiquetas `<video>` HTML5 estándar para el preview. Esto garantiza **0 latencia** y aceleración de hardware nativa del navegador.
*   [ ] **Limpieza:** Eliminar bucles de renderizado (`tickers`) innecesarios en la interfaz.

## 📐 Fase 2: Layouts CSS (Reemplazo de Lógica Visual)
**Objetivo:** Que los modos Solo, Grid y Zoom se vean perfectos en cualquier pantalla.
*   [ ] **Implementar CSS Grid:** Crear contenedores CSS para cada layout:
    *   **Solo:** `w-full h-full object-cover`.
    *   **Grid:** `grid-cols-2 gap-0`.
    *   **Zoom:** `flex gap-4 items-center justify-center` con fondo estilizado.
*   [ ] **Responsive:** Asegurar que se adapte al móvil sin cálculos matemáticos complejos, usando las reglas nativas del navegador.

## ⚙️ Fase 3: Motor de Composición (Invisible)
**Objetivo:** Preparar la señal para FFmpeg sin afectar la UI.
*   [ ] **Canvas Oculto:** Re-introducir PixiJS pero **solo en memoria/oculto**.
*   [ ] **Clonación de Stream:** Copiar los streams de video al canvas oculto solo para "dibujar" la imagen que se enviará al servidor RTMP.
*   [ ] **Optimización:** Limitar los FPS del canvas oculto a 30fps fijos (estándar de transmisión) para no saturar la CPU, mientras la UI va fluida a 60fps.

## 📡 Fase 4: Conexión y Señalización Robustas
**Objetivo:** Que la conexión de invitados sea "a prueba de balas".
*   [ ] **Indicadores de Estado:** Mostrar claramente "Conectando...", "Reconectando", "Sin señal".
*   [ ] **Manejo de Errores:** Si el invitado se desconecta, volver automáticamente al layout "Solo".

## 🚀 Fase 5: Transmisión (Backend)
**Objetivo:** Enviar la señal final.
*   [ ] **Captura de Canvas:** Usar `canvas.captureStream()` del canvas oculto.
*   [x] **Fix Ghosting:** Eliminadas transiciones CSS en banners (instantáneo).
*   [ ] **Grabación Robusta:** Sincronización de Keyframes (SPS/PPS) en grabaciones iniciadas post-stream (Pendiente).
*   [ ] **Pipe a Go:** Conectar el stream capturado con el proceso FFmpeg en el backend.

---

### ¿Por qué este cambio?
Actualmente estamos intentando que el Canvas haga todo (UI + Procesamiento). Al separar la **Vista (HTML)** del **Procesamiento (Canvas)**, logramos:
1.  **Fluidez total** en la interfaz (el navegador optimiza los videos mejor que nadie).
2.  **Layouts perfectos** gracias a CSS Grid.
3.  **Menor consumo de CPU**, reservando la potencia para cuando activemos la transmisión real.