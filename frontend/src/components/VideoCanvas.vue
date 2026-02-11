<script setup>
import { onMounted, ref, watch, onUnmounted } from 'vue';
import { useStudioStore } from '../stores/studio';
import * as PIXI from 'pixi.js';
import * as AppGo from '../services/api';

const studio = useStudioStore();
const canvasContainer = ref(null);

let app = null;
let resizeObserver = null;

const participants = new Map(); 
let backgroundLayer, videoLayer, overlayLayer;
let bgSprite = null;
let logoSprite = null;

// --- GESTIÓN DE ELEMENTOS DE MARCA ---
const updateBrandElements = async () => {
  if (!app) return;
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

  if (studio.logoUrl) {
    try {
      if (!logoSprite) {
        logoSprite = new PIXI.Sprite();
        overlayLayer.addChild(logoSprite);
      }
      const texture = await PIXI.Assets.load(studio.logoUrl);
      logoSprite.texture = texture;
      logoSprite.scale.set(studio.logoScale);
      
      const margin = 40;
      const w = texture.width * studio.logoScale;
      const h = texture.height * studio.logoScale;

      if (studio.logoPosition === 'top-left') {
        logoSprite.x = margin; logoSprite.y = margin;
      } else if (studio.logoPosition === 'top-right') {
        logoSprite.x = 1280 - w - margin; logoSprite.y = margin;
      } else if (studio.logoPosition === 'bottom-left') {
        logoSprite.x = margin; logoSprite.y = 720 - h - margin;
      } else if (studio.logoPosition === 'bottom-right') {
        logoSprite.x = 1280 - w - margin; logoSprite.y = 720 - h - margin;
      }
    } catch (e) { console.error("Error logo:", e); }
  } else if (logoSprite) {
    logoSprite.destroy();
    logoSprite = null;
  }
};

let bannerContainer = null;
let tickerFn = null;
let lastBannerId = null;

const updateBanner = () => {
  if (!app) return;
  
  const activeBanner = studio.banners.find(b => b.active);
  
  // 1. Si no hay banner activo, limpiar y salir
  if (!activeBanner) {
    if (tickerFn) { app.ticker.remove(tickerFn); tickerFn = null; }
    if (bannerContainer) { bannerContainer.destroy({ children: true }); bannerContainer = null; }
    lastBannerId = null;
    return;
  }

  // 2. Si el banner cambió (ID o Texto), lo recreamos completamente
  const bannerKey = `${activeBanner.id}-${activeBanner.text}-${activeBanner.subtext}-${studio.bannerBgColor}-${studio.bannerTextColor}-${studio.bannerBorderRadius}-${studio.bannerFont}-${activeBanner.isTicker}-${activeBanner.isMarquee}`;
  
  if (lastBannerId !== bannerKey) {
    if (tickerFn) { app.ticker.remove(tickerFn); tickerFn = null; }
    if (bannerContainer) { bannerContainer.destroy({ children: true }); bannerContainer = null; }
    
    bannerContainer = new PIXI.Container();
    const font = studio.bannerFont || 'Inter';
    const bgColor = new PIXI.Color(activeBanner.bgColor || studio.bannerBgColor).toNumber();
    const textColor = activeBanner.textColor || studio.bannerTextColor;
    const borderRadius = activeBanner.borderRadius ?? studio.bannerBorderRadius;

    if (activeBanner.isTicker) {
      const bg = new PIXI.Graphics().beginFill(bgColor).drawRect(0, 0, 1280, 50).endFill();
      const baseTextStr = activeBanner.text.toUpperCase() + (activeBanner.subtext ? ` • ${activeBanner.subtext.toUpperCase()}` : "") + "          ";
      const text = new PIXI.Text({ 
        text: baseTextStr.repeat(10), 
        style: { fontFamily: font, fontSize: 24, fill: textColor, fontWeight: 'bold', letterSpacing: 2 } 
      });
      text.y = 12;
      bannerContainer.addChild(bg, text);
      
      const singleText = new PIXI.Text({ text: baseTextStr, style: text.style });
      const singleWidth = singleText.width;
      singleText.destroy();

      tickerFn = () => {
        if (text) {
          text.x -= 2.2;
          if (text.x <= -singleWidth) text.x = 0;
        }
      };
      app.ticker.add(tickerFn);
      bannerContainer._isTickerMode = true;
    } else {
      const bg = new PIXI.Graphics();
      const padding = 40;
      const textContent = new PIXI.Container();

      const mainText = new PIXI.Text({ 
        text: activeBanner.text, 
        style: { fontFamily: font, fontSize: 42, fill: textColor, fontWeight: '900', letterSpacing: -1 } 
      });
      mainText.x = padding;
      mainText.y = activeBanner.subtext ? 15 : 28;

      const subText = activeBanner.subtext ? new PIXI.Text({ 
        text: activeBanner.subtext.toUpperCase(), 
        style: { fontFamily: font, fontSize: 16, fill: textColor, fontWeight: '700', letterSpacing: 2 } 
      }) : null;
      
      if (subText) {
        subText.alpha = 0.8;
        subText.x = mainText.x;
        subText.y = mainText.y + 48;
      }

      textContent.addChild(mainText);
      if (subText) textContent.addChild(subText);

      const bannerW = Math.max(800, textContent.width + (padding * 2));
      const bannerH = 110;
      
      bg.beginFill(bgColor).drawRoundedRect(0, 0, bannerW, bannerH, borderRadius).endFill();
      bg.lineStyle(1, 0xffffff, 0.1).drawRoundedRect(0, 0, bannerW, bannerH, borderRadius);
      
      bannerContainer.addChild(bg, textContent);
      bannerContainer.pivot.set(bannerW / 2, bannerH / 2);

      if (activeBanner.isMarquee) {
        const mask = new PIXI.Graphics().beginFill(0xffffff).drawRect(0, 0, bannerW, bannerH).endFill();
        textContent.mask = mask;
        bannerContainer.addChild(mask);
        tickerFn = () => {
          mainText.x -= 2;
          if (subText) subText.x = mainText.x;
          if (mainText.x < -mainText.width) {
              mainText.x = bannerW;
              if (subText) subText.x = bannerW;
          }
        };
        app.ticker.add(tickerFn);
      }
      bannerContainer._isTickerMode = false;
    }
    overlayLayer.addChild(bannerContainer);
    lastBannerId = bannerKey;
  }

  // 3. ACTUALIZACIÓN DE POSICIÓN (Instantánea, sin recrear)
  if (bannerContainer) {
    if (bannerContainer._isTickerMode) {
      bannerContainer.x = 0;
      bannerContainer.y = 720 - 50;
    } else {
      bannerContainer.x = (studio.bannerPositionX / 100) * 1280;
      bannerContainer.y = (studio.bannerPositionY / 100) * 720;
    }
  }
};

// --- GESTIÓN DE PARTICIPANTES ---
const calculateLayout = (screenW, screenH, count, type) => {
  const padding = 100; const gap = 40;     
  const safeW = screenW - (padding * 2); const safeH = screenH - (padding * 2);
  const slots = [];
  if (count === 0) return slots;

  if (type === 'solo' || count === 1) {
    // Modo Solo: Más pequeño y centrado para un look más "estudio"
    const soloW = safeW * 0.85;
    const soloH = soloW * (9/16);
    slots.push({ 
      x: (screenW - soloW) / 2, 
      y: (screenH - soloH) / 2, 
      w: soloW, 
      h: soloH, 
      radius: 2 
    });
  } else if (type === 'grid') {
    if (count === 2) {
      const w = (safeW - gap) / 2; const h = safeH * 0.9; const y = (screenH - h) / 2;
      slots.push({ x: padding, y, w, h, radius: 4 }, { x: padding + w + gap, y, w, h, radius: 4 });
    } else {
      const cols = Math.ceil(Math.sqrt(count)); const rows = Math.ceil(count / cols);
      const w = (safeW - (gap * (cols - 1))) / cols; const h = (safeH - (gap * (rows - 1))) / rows;
      for (let i = 0; i < count; i++) {
        slots.push({ x: padding + (i % cols) * (w + gap), y: padding + Math.floor(i / cols) * (h + gap), w, h, radius: 4 });
      }
    }
  } else if (type === 'sidebar') {
    const mainW = safeW * 0.78; const sideW = safeW * 0.22 - gap;
    slots.push({ x: padding, y: padding, w: mainW, h: safeH, radius: 4 });
    const sideCount = count - 1; const sideH = (safeH - (gap * (sideCount - 1))) / sideCount;
    for (let i = 0; i < sideCount; i++) {
      slots.push({ x: padding + mainW + gap, y: padding + i * (sideH + gap), w: sideW, h: sideH, radius: 4 });
    }
  } else if (type === 'pip') {
    slots.push({ x: 0, y: 0, w: screenW, h: screenH, radius: 0 }); 
    const pipW = screenW * 0.25; const pipH = pipW * (9/16);
    for (let i = 1; i < count; i++) {
      slots.push({ x: screenW - padding - pipW, y: screenH - padding - pipH, w: pipW, h: pipH, radius: 4 });
    }
  }
  return slots;
};

const updateParticipant = (id, x, y, w, h, radius = 4) => {
  const p = participants.get(id); if (!p) return;
  p.container.x = x; p.container.y = y;
  const videoW = p.video.videoWidth || 1280; const videoH = p.video.videoHeight || 720;
  const scale = Math.max(w / videoW, h / videoH);
  p.sprite.scale.set(scale);
  p.sprite.x = (w - videoW * scale) / 2; p.sprite.y = (h - videoH * scale) / 2;
  p.mask.clear().beginFill(0xffffff).drawRoundedRect(0, 0, w, h, radius).endFill();
  if (!p.border) { p.border = new PIXI.Graphics(); p.container.addChild(p.border); }
  p.border.clear().lineStyle(1, 0xffffff, 0.2).drawRoundedRect(0, 0, w, h, radius);
  p.text.visible = studio.showParticipantNames; p.nameBg.visible = studio.showParticipantNames;
  if (studio.showParticipantNames) {
    p.text.style.fontFamily = 'Inter'; p.text.style.fontSize = 12;
    p.text.x = 10; p.text.y = h - p.text.height - 10;
    p.nameBg.clear().beginFill(0x000000, 0.6).drawRect(5, h - p.text.height - 15, p.text.width + 10, p.text.height + 10).endFill();
  }
};

const createParticipant = async (id, stream, name) => {
  const video = document.createElement('video');
  video.srcObject = stream; video.muted = true; video.autoplay = true; video.playsInline = true;
  video.style.position = 'absolute'; video.style.width = '1px'; video.style.height = '1px'; video.style.opacity = '0';
  document.body.appendChild(video);
  video.onloadedmetadata = () => syncStudio();
  await video.play().catch(e => console.warn("Video play error", e));

  const container = new PIXI.Container();
  const texture = PIXI.Texture.from(video, { resourceOptions: { autoPlay: true } });
  const sprite = new PIXI.Sprite(texture);
  const mask = new PIXI.Graphics(); const nameBg = new PIXI.Graphics();
  const text = new PIXI.Text({ text: name, style: { fontFamily: 'Inter', fontSize: 12, fill: 'white', fontWeight: 'bold' } });
  container.addChild(sprite, mask, nameBg, text);
  sprite.mask = mask; videoLayer.addChild(container);
  participants.set(id, { container, sprite, mask, text, nameBg, video });
};

const syncStudio = async () => {
  if (!app) return;
  await updateBrandElements(); updateBanner();
  const activeParticipants = [];
  if (studio.localStream && studio.participantsOnStage.includes('local')) {
    activeParticipants.push({ id: 'local', stream: studio.localStream, name: studio.userName });
  }
  if (studio.guestConnected && studio.guestStream && studio.participantsOnStage.includes('guest')) {
    activeParticipants.push({ id: 'guest', stream: studio.guestStream, name: 'Invitado' });
  }
  for (const [id, p] of participants) {
    if (!activeParticipants.find(ap => ap.id === id)) {
      videoLayer.removeChild(p.container); p.video.remove(); participants.delete(id);
    }
  }
  for (const ap of activeParticipants) {
    if (!participants.has(ap.id)) await createParticipant(ap.id, ap.stream, ap.name);
  }
  const slots = calculateLayout(1280, 720, activeParticipants.length, studio.layout);
  activeParticipants.forEach((ap, i) => {
    const slot = slots[i]; if (slot) updateParticipant(ap.id, slot.x, slot.y, slot.w, slot.h, slot.radius);
  });
};

const initPixi = async () => {
  app = new PIXI.Application();
  await app.init({ width: 1280, height: 720, background: '#000000', antialias: true, hello: false });
  canvasContainer.value.appendChild(app.canvas);
  backgroundLayer = new PIXI.Container();
  videoLayer = new PIXI.Container();
  overlayLayer = new PIXI.Container();
  app.stage.addChild(backgroundLayer, videoLayer, overlayLayer);
  app.ticker.add(() => {
    participants.forEach(p => { if (p.sprite?.texture?.source) p.sprite.texture.source.update(); });
  });

  resizeObserver = new ResizeObserver(entries => {
    for (let entry of entries) {
      const cw = entry.contentRect.width; const ch = entry.contentRect.height;
      const targetRatio = 16 / 9; const containerRatio = cw / ch;
      let finalW, finalH;
      if (containerRatio > targetRatio) { finalH = ch; finalW = ch * targetRatio; }
      else { finalW = cw; finalH = cw / targetRatio; }
      if (app.canvas) { app.canvas.style.width = `${finalW}px`; app.canvas.style.height = `${finalH}px`; }
    }
  });
  resizeObserver.observe(canvasContainer.value);
};

// MEZCLA DE AUDIO EN EL STREAM
const getStudioStream = () => {
  const stream = app.canvas.captureStream(30);
  if (studio.audioDestination) {
    const audioTrack = studio.audioDestination.stream.getAudioTracks()[0];
    if (audioTrack) stream.addTrack(audioTrack);
  }
  return stream;
};

defineExpose({ 
  connectToStudio: async () => {
    const stream = getStudioStream();
    const pc = new RTCPeerConnection({ iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] });
    stream.getTracks().forEach(t => pc.addTrack(t, stream));
    const offer = await pc.createOffer();
    await pc.setLocalDescription(offer);
    const res = await AppGo.SendOffer(offer.sdp);
    if (res?.sdp) await pc.setRemoteDescription(new RTCSessionDescription({ type: 'answer', sdp: res.sdp }));
  },
  startRecording: async () => {
    if (!studio.isStreaming) {
      const stream = getStudioStream();
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
  }
});

onMounted(async () => {
  await initPixi();
  watch([() => studio.isInitialized, () => studio.guestConnected, () => studio.layout, () => studio.localStream, () => studio.participantsOnStage], syncStudio, { immediate: true, deep: true });
  // watch(() => studio.banners, updateBanner, { deep: true }); // Comentado: ahora manejado por Vue en App.vue
  watch([() => studio.logoUrl, () => studio.logoScale, () => studio.backgroundUrl, () => studio.brandColor, () => studio.bannerBgColor, () => studio.bannerTextColor, () => studio.bannerBorderRadius, () => studio.bannerFont, () => studio.bannerPositionX, () => studio.bannerPositionY, () => studio.showParticipantNames, () => studio.showLiveBadge], () => {
    updateBrandElements(); updateBanner(); syncStudio();
  });
});

onUnmounted(() => {
  if (resizeObserver) resizeObserver.disconnect();
  participants.forEach(p => p.video.remove());
  if (app) app.destroy(true);
});
</script>

<template>
  <div ref="canvasContainer" class="w-full h-full flex items-center justify-center bg-black overflow-hidden relative">
    <div v-if="studio.showLiveBadge && studio.isRecording" class="absolute top-2 left-2 z-10 bg-red-600 text-white px-1.5 py-0.5 rounded-sm text-[9px] font-bold animate-pulse pointer-events-none border border-white/20">
      REC
    </div>
  </div>
</template>