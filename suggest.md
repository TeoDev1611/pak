El tamaño del video ya es correcto, pero la interfaz se ve "fea" y vacía en los
lados. Parece un parche negro gigante. Necesito mejorar la **Estética del
Escenario (Stage)** para que parezca un software professional.

**TU MISIÓN DE DISEÑO (CSS/Tailwind):**

1. **El Fondo del Estudio (Background):** En lugar de un negro plano
   (`bg-[#0d1117]`), quiero que el contenedor padre (el que está detrás del
   video) tenga profundidad.
   - Usa un degradado sutil: `bg-gradient-to-b from-gray-900 to-black`.
   - (Opcional) Añade un patrón sutil de puntos o grilla si es possible con CSS,
     o simplemente mantén el degradado elegante.

2. **El Lienzo de Video (Canvas Wrapper):** Ahora mismo se funde con el fondo.
   Necesito que **resalte**. Aplica estas clases al `div` que envuelve al
   `<VideoCanvas>`:
   - **Sombra Profunda:** `shadow-[0_0_50px_rgba(0,0,0,0.5)]` (para que parezca
     que flota).
   - **Borde Sutil:** `border border-white/10` (un borde fino para delimitar
     dónde termina el video).
   - **Esquinas:** `rounded-xl` (para suavizar los bordes).
   - **Color de Relleno:** `bg-black` (el canvas en sí debe set negro puro).

3. **Barra de Herramientas Flotante (Layout Switcher):** Actualmente flota de
   forma extraña.
   - Haz que parezca un **"Dock" de crystal (Glassmorphism)**.
   - Clases sugeridas:
     `bg-black/40 backdrop-blur-md border border-white/10 rounded-full px-4 py-2`.
   - Colócalo unos píxeles _debajo_ del video o superpuesto en la parte inferior
     central con margen.

**RESUMEN DE ESTRUCTURA EN `App.vue`:**

```html
<div
  class="flex-1 relative flex items-center justify-center bg-gradient-to-b from-slate-900 to-black overflow-hidden p-6"
>
  <div
    ref="stageContainer"
    class="aspect-video w-full max-h-full bg-black rounded-xl shadow-2xl border border-white/10 relative overflow-hidden flex items-center justify-center"
  >
    <VideoCanvas ref="videoCanvas" />

    <div
      class="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2 bg-black/60 backdrop-blur-md p-2 rounded-2xl border border-white/10 z-50"
    >
    </div>
  </div>
</div>
```

Por favor, actualiza la sección del `<main>` en `App.vue` con estos estilos
visuals mejorados.

### ¿Qué va a cambiar visualmente?

1. **Profundidad:** Al poner un degradado gris oscuro de fondo y el video negro
   encima con borde y sombra, creas **contraste**. El cerebro entenderá: _"Ah,
   eso gris es la mesa de trabajo y eso negro brillante es la pantalla"_.
2. **Límites:** El borde `border-white/10` es clave. Aunque sea sutil, dibuja
   una línea que dice "aquí termina el video", eliminando la sensación de vacío
   infinito.
3. **Elegancia:** El efecto _glassmorphism_ (fondo borroso) en los botones los
   have ver modernos en lugar de botones pegados con paint.
