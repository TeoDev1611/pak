<script setup>
import { onMounted, ref, watch, onUnmounted } from 'vue';
import { useStudioStore } from '../stores/studio';
import * as PIXI from 'pixi.js';
import * as AppGo from '../services/api';

const studio = useStudioStore();
const canvasContainer = ref(null);

let app = null;
let mediaRecorder = null;
let recordedChunks = [];

const participants = new Map(); 
let backgroundLayer, videoLayer, overlayLayer;
let bgSprite = null;
let logoSprite = null;
let bannerContainer = null;
let tickerFn = null; // Guardar referencia para limpiar la animación

// --- GESTIÓN DE ELEMENTOS DE MARCA ---
const updateBrandElements = async () => {
  if (!app) return;

  // 1. Fondo
  if (studio.backgroundUrl) {
    try {
      if (!bgSprite) {
        bgSprite = new PIXI.Sprite();
        backgroundLayer.addChild(bgSprite);
      }
      const texture = await PIXI.Assets.load(studio.backgroundUrl);
      bgSprite.texture = texture;
      const scale = Math.max(1280 / texture.width, 720 / texture.height);
      bgSprite.scale.set(scale);
      bgSprite.x = (1280 - texture.width * scale) / 2;
      bgSprite.y = (720 - texture.height * scale) / 2;
    } catch (e) { console.error("Error background:", e); }
  } else if (bgSprite) {
    bgSprite.destroy();
    bgSprite = null;
  }

  // 2. Logo
  if (studio.logoUrl) {
    try {
      if (!logoSprite) {
        logoSprite = new PIXI.Sprite();
        overlayLayer.addChild(logoSprite);
      }
      const texture = await PIXI.Assets.load(studio.logoUrl);
      logoSprite.texture = texture;
      // Usar la escala personalizada del store
      logoSprite.scale.set(studio.logoScale);
      logoSprite.x = 1280 - (texture.width * studio.logoScale) - 40;
      logoSprite.y = 40;
    } catch (e) { console.error("Error logo:", e); }
  } else if (logoSprite) {
    logoSprite.destroy();
    logoSprite = null;
  }
};

const updateBanner = () => {
  if (!app) return;
  
  // Limpiar ticker anterior si existe
  if (tickerFn) {
    app.ticker.remove(tickerFn);
    tickerFn = null;
  }

  if (bannerContainer) {
    bannerContainer.destroy({ children: true });
    bannerContainer = null;
  }

  const activeBanner = studio.banners.find(b => b.active);
  if (!activeBanner) return;

  bannerContainer = new PIXI.Container();

  if (activeBanner.isTicker) {
    const font = studio.bannerFont === 'Mono' ? 'monospace' : studio.bannerFont;
    const bg = new PIXI.Graphics().beginFill(studio.brandColor).drawRect(0, 0, 1280, 50).endFill();
    
    // Crear el texto base una vez para medirlo
    const baseTextStr = activeBanner.text.toUpperCase() + (activeBanner.subtext ? ` • ${activeBanner.subtext.toUpperCase()}` : "") + "          ";
    
    const text = new PIXI.Text({ 
      text: baseTextStr.repeat(10), // Repetir suficientes veces para llenar el ancho + margen
      style: { 
        fontFamily: font, 
        fontSize: 24, 
        fill: 'white', 
        fontWeight: 'bold',
        letterSpacing: 2
      } 
    });
    text.y = 12;
    
    bannerContainer.addChild(bg, text);
    bannerContainer.y = 720 - 50;
    
    // Medir el ancho de una sola repetición para el reinicio perfecto
    const singleText = new PIXI.Text({ text: baseTextStr, style: text.style });
    const singleWidth = singleText.width;
    singleText.destroy();

    tickerFn = () => {
      if (text) {
        text.x -= 2.2;
        // Si el primer bloque ya salió de la pantalla, reseteamos la posición sutilmente
        if (text.x <= -singleWidth) {
          text.x = 0;
        }
      }
    };
    app.ticker.add(tickerFn);
  } else {
    // --- BANNER ESTÁTICO CON ESTILOS ---
    const font = studio.bannerFont === 'Mono' ? 'monospace' : studio.bannerFont;
    const bg = new PIXI.Graphics();
    
    if (studio.bannerStyle === 'classic') {
      bg.beginFill(studio.brandColor).drawRect(0, 0, 850, 90).endFill();
    } else if (studio.bannerStyle === 'neon') {
      bg.lineStyle(4, studio.brandColor, 1)
        .beginFill(0x000000, 0.8)
        .drawRoundedRect(0, 0, 800, 85, 4)
        .endFill();
    } else if (studio.bannerStyle === 'minimal') {
      bg.beginFill(0xffffff, 0.95).drawRect(0, 20, 10, 50).endFill();
    } else {
      bg.beginFill(studio.brandColor).drawRoundedRect(0, 0, 800, 85, 12).endFill();
    }
    
    const textColor = studio.bannerStyle === 'minimal' ? '#000000' : 'white';
    
    const text = new PIXI.Text({ 
      text: activeBanner.text.toUpperCase(), 
      style: { 
        fontFamily: font, 
        fontSize: 30, 
        fill: textColor, 
        fontWeight: '900',
        dropShadow: studio.bannerStyle === 'classic',
        dropShadowBlur: 4
      } 
    });
    text.x = studio.bannerStyle === 'minimal' ? 25 : 35; 
    text.y = 12;

    const sub = new PIXI.Text({ 
      text: activeBanner.subtext, 
      style: { 
        fontFamily: font, 
        fontSize: 20, 
        fill: textColor, 
        alpha: 0.7 
      } 
    });
    sub.x = text.x; sub.y = 48;

    bannerContainer.addChild(bg, text, sub);
    bannerContainer.x = studio.bannerX;
    bannerContainer.y = studio.bannerY;
  }
  
  overlayLayer.addChild(bannerContainer);
};

// --- MOTOR DE LAYOUTS (CÁLCULO GEOMÉTRICO) ---
const calculateLayout = (screenW, screenH, count, type) => {
  // Ajustes para 1080p
  const padding = 60; 
  const gap = 30;     
  const safeW = screenW - (padding * 2);
  const safeH = screenH - (padding * 2);
  
  const slots = [];
  if (count === 0) return slots;

  // 1. SOLO
  if (type === 'solo' || count === 1) {
    const w = safeW;
    const h = safeH;
    slots.push({ x: padding, y: padding, w, h, radius: 40 });
  } 
  
  // 2. GRID
  else if (type === 'grid') {
    if (count === 2) {
      const w = (safeW - gap) / 2;
      const h = safeH * 0.9; 
      const y = (screenH - h) / 2;
      slots.push({ x: padding, y, w, h, radius: 30 });
      slots.push({ x: padding + w + gap, y, w, h, radius: 30 });
    } else {
      const cols = Math.ceil(Math.sqrt(count));
      const rows = Math.ceil(count / cols);
      const w = (safeW - (gap * (cols - 1))) / cols;
      const h = (safeH - (gap * (rows - 1))) / rows;
      for (let i = 0; i < count; i++) {
        slots.push({
          x: padding + (i % cols) * (w + gap),
          y: padding + Math.floor(i / cols) * (h + gap),
          w, h, radius: 25
        });
      }
    }
  }

  // 3. SIDEBAR
  else if (type === 'sidebar') {
    const mainW = safeW * 0.78;
    const sideW = safeW * 0.22 - gap;
    slots.push({ x: padding, y: padding, w: mainW, h: safeH, radius: 40 });
    
    const sideCount = count - 1;
    const sideH = (safeH - (gap * (sideCount - 1))) / sideCount;
    for (let i = 0; i < sideCount; i++) {
      slots.push({
        x: padding + mainW + gap,
        y: padding + i * (sideH + gap),
        w: sideW, h: sideH, radius: 20
      });
    }
  }

  // 4. PIP
  else if (type === 'pip') {
    slots.push({ x: 0, y: 0, w: screenW, h: screenH, radius: 0 }); 
    const pipW = screenW * 0.25;
    const pipH = pipW * (9/16);
    for (let i = 1; i < count; i++) {
      slots.push({
        x: screenW - padding - pipW,
        y: screenH - padding - pipH,
        w: pipW, h: pipH, radius: 25
      });
    }
  }

  return slots;
};

// --- GESTIÓN DE PARTICIPANTES ---
const updateParticipant = (id, x, y, w, h, radius = 20) => {
  const p = participants.get(id);
  if (!p) return;

  // Animación suave de posición (opcional, pero mejora el look)
  p.container.x = x;
  p.container.y = y;

  const videoW = p.video.videoWidth || 1280;
  const videoH = p.video.videoHeight || 720;
  const scale = Math.max(w / videoW, h / videoH);
  p.sprite.scale.set(scale);
  p.sprite.x = (w - videoW * scale) / 2;
  p.sprite.y = (h - videoH * scale) / 2;

  // 1. Máscara con bordes perfectos
  p.mask.clear().beginFill(0xffffff).drawRoundedRect(0, 0, w, h, radius).endFill();
  
  // 2. Borde exterior elegante (Stroke)
  if (!p.border) {
    p.border = new PIXI.Graphics();
    p.container.addChild(p.border);
  }
  p.border.clear()
    .lineStyle(2, 0xffffff, 0.1) // Borde muy sutil
    .drawRoundedRect(0, 0, w, h, radius);

  // 3. Etiqueta de Nombre Profesional
  p.text.visible = studio.showParticipantNames;
  p.nameBg.visible = studio.showParticipantNames;

  if (studio.showParticipantNames) {
    const font = studio.bannerFont === 'Mono' ? 'monospace' : studio.bannerFont;
    p.text.style.fontFamily = font;
    p.text.style.fontSize = Math.max(14, w * 0.035); // Tamaño dinámico
    p.text.x = 20;
    p.text.y = h - p.text.height - 15;
    
    p.nameBg.clear()
      .beginFill(0x000000, 0.5)
      .drawRoundedRect(12, h - p.text.height - 20, p.text.width + 16, p.text.height + 8, 6)
      .endFill();
  }
};

const createParticipant = async (id, stream, name) => {
  const video = document.createElement('video');
  video.srcObject = stream;
  video.muted = true;
  video.autoplay = true;
  video.playsInline = true;
  video.style.position = 'absolute'; video.style.width = '1px'; video.style.height = '1px'; video.style.opacity = '0';
  document.body.appendChild(video);
  
  // Recalcular layout cuando el video tenga dimensiones reales
  video.onloadedmetadata = () => {
    syncStudio();
  };

  await video.play().catch(e => console.warn("Video play error", e));

  const container = new PIXI.Container();
  const texture = PIXI.Texture.from(video, { resourceOptions: { autoPlay: true } });
  const sprite = new PIXI.Sprite(texture);
  const mask = new PIXI.Graphics();
  const nameBg = new PIXI.Graphics();
  const text = new PIXI.Text({ text: name, style: { fontFamily: 'Inter', fontSize: 20, fill: 'white', fontWeight: 'bold' } });

  container.addChild(sprite, mask, nameBg, text);
  sprite.mask = mask;
  videoLayer.addChild(container);

  participants.set(id, { container, sprite, mask, text, nameBg, video });
};

const syncStudio = async () => {
  if (!app) return;

  await updateBrandElements();
  updateBanner();

  const activeParticipants = [];
  
  // Solo agregar si están en la lista de "On Stage"
  if (studio.localStream && studio.participantsOnStage.includes('local')) {
    activeParticipants.push({ id: 'local', stream: studio.localStream, name: studio.userName });
  }
  
  if (studio.guestConnected && studio.guestStream && studio.participantsOnStage.includes('guest')) {
    activeParticipants.push({ id: 'guest', stream: studio.guestStream, name: 'Invitado' });
  }

  // 1. Limpiar eliminados
  for (const [id, p] of participants) {
    if (!activeParticipants.find(ap => ap.id === id)) {
      videoLayer.removeChild(p.container);
      p.video.remove();
      participants.delete(id);
    }
  }

  // 2. Crear nuevos
  for (const ap of activeParticipants) {
    if (!participants.has(ap.id)) {
      await createParticipant(ap.id, ap.stream, ap.name);
    }
  }

  // 3. Aplicar Layout
  const slots = calculateLayout(1280, 720, activeParticipants.length, studio.layout);
  activeParticipants.forEach((ap, i) => {
    const slot = slots[i];
    if (slot) updateParticipant(ap.id, slot.x, slot.y, slot.w, slot.h, slot.radius);
  });
};

const initPixi = async () => {
  app = new PIXI.Application();
  await app.init({ 
    width: 1920, 
    height: 1080, 
    background: '#020617', 
    antialias: true,
    hello: false
  });
  canvasContainer.value.appendChild(app.canvas);
  
  // Ajuste de estilo para que el canvas se adapte manteniendo el aspect ratio
  app.canvas.style.width = '100%';
  app.canvas.style.height = '100%';
  app.canvas.style.objectFit = 'contain';

  // Inicializar Capas
  backgroundLayer = new PIXI.Container();
  videoLayer = new PIXI.Container();
  overlayLayer = new PIXI.Container();
  app.stage.addChild(backgroundLayer, videoLayer, overlayLayer);

  // Forzar actualización de texturas en cada frame
  app.ticker.add(() => {
    participants.forEach(p => {
      if (p.sprite && p.sprite.texture && p.sprite.texture.source) {
        p.sprite.texture.source.update();
      }
    });
  });

  // Escuchar redimensionamiento de ventana
  window.addEventListener('resize', () => {
    if (app && app.renderer) {
      // El CSS se encarga del tamaño visual
    }
  });
};

defineExpose({ 
  connectToStudio: async () => {
    const stream = app.canvas.captureStream(30);
    const pc = new RTCPeerConnection({ iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] });
    stream.getTracks().forEach(t => pc.addTrack(t, stream));
    const offer = await pc.createOffer();
    await pc.setLocalDescription(offer);
    const res = await AppGo.SendOffer(offer.sdp);
    if (res?.sdp) await pc.setRemoteDescription(new RTCSessionDescription({ type: 'answer', sdp: res.sdp }));
  },
  startRecording: async () => {
    // Aseguramos que haya conexión WebRTC activa, ya que el backend graba lo que recibe por ahí
    // Nota: Idealmente deberíamos comprobar si ya está conectado
    if (!studio.isStreaming) { 
        // Si no está stremeando, forzamos la conexión "invisible" para enviar video al backend
        // Reusamos la lógica de conexión
        const stream = app.canvas.captureStream(30);
        const pc = new RTCPeerConnection({ iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] });
        stream.getTracks().forEach(t => pc.addTrack(t, stream));
        const offer = await pc.createOffer();
        await pc.setLocalDescription(offer);
        const res = await AppGo.SendOffer(offer.sdp);
        if (res?.sdp) await pc.setRemoteDescription(new RTCSessionDescription({ type: 'answer', sdp: res.sdp }));
    }

    const filename = studio.recordingTitle || `grabacion-${Date.now()}`;
    await AppGo.StartRecording(filename + ".mp4");
    studio.isRecording = true;
  },
  stopRecording: async () => {
    await AppGo.StopRecording();
    studio.isRecording = false;
    alert("Grabación finalizada y guardada en el backend.");
  }
});

onMounted(async () => {
  await initPixi();
  watch(
    [
      () => studio.isInitialized, 
      () => studio.guestConnected, 
      () => studio.layout, 
      () => studio.localStream,
      () => studio.participantsOnStage
    ], 
    syncStudio, 
    { immediate: true, deep: true }
  );

  watch(() => studio.banners, updateBanner, { deep: true });
  watch([
    () => studio.logoUrl, 
    () => studio.logoScale,
    () => studio.backgroundUrl, 
    () => studio.brandColor,
    () => studio.bannerStyle,
    () => studio.bannerFont,
    () => studio.bannerX,
    () => studio.bannerY,
    () => studio.showParticipantNames,
    () => studio.showLiveBadge
  ], () => {
    updateBrandElements();
    updateBanner();
    syncStudio(); // Forzar actualización de nombres
  });
});

onUnmounted(() => {
  participants.forEach(p => p.video.remove());
  if (app) app.destroy(true);
});
</script>

<template>
  <div class="stage-wrapper h-full w-full flex items-center justify-center bg-black">
    <div ref="canvasContainer" class="pixi-container w-full h-full flex items-center justify-center overflow-hidden">
       <div v-if="studio.showLiveBadge" class="absolute top-6 left-6 z-10 bg-red-600 text-white px-2 py-0.5 rounded text-[8px] font-black animate-pulse flex items-center gap-1.5 shadow-xl">
         <span class="w-1.5 h-1.5 bg-white rounded-full"></span> LIVE
       </div>
    </div>
  </div>
</template>