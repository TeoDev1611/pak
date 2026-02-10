<script setup>
import { onMounted, ref, watch, onUnmounted } from 'vue';
import { useStudioStore } from '../stores/studio';

const studio = useStudioStore();
const localVideo = ref(null);
const guestVideo = ref(null);

let pc = null, ws = null;

const initWebRTC = () => {
  if (ws) ws.close();
  const protocol = location.protocol === 'https:' ? 'wss' : 'ws';
  ws = new WebSocket(`${protocol}://${location.host}/ws`);
  pc = new RTCPeerConnection({ iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] });

  if (studio.localStream) {
    studio.localStream.getTracks().forEach(t => pc.addTrack(t, studio.localStream));
  }

  pc.ontrack = (e) => {
    if (guestVideo.value) {
      guestVideo.value.srcObject = e.streams[0];
      studio.guestConnected = true;
    }
  };

  pc.onicecandidate = (e) => {
    if (e.candidate) ws.send(JSON.stringify({ type: 'candidate', candidate: e.candidate }));
  };

  ws.onmessage = async (e) => {
    const msg = JSON.parse(e.data);
    if (msg.type === 'offer') {
      await pc.setRemoteDescription(new RTCSessionDescription(msg));
      const ans = await pc.createAnswer(); await pc.setLocalDescription(ans);
      ws.send(JSON.stringify({ type: 'answer', sdp: ans.sdp }));
    } else if (msg.type === 'answer') {
      await pc.setRemoteDescription(new RTCSessionDescription(msg));
    } else if (msg.type === 'candidate') {
      await pc.addIceCandidate(new RTCIceCandidate(msg.candidate));
    }
  };
};

onMounted(() => {
  watch(() => studio.isInitialized, (val) => {
    if (val && studio.localStream) {
      localVideo.value.srcObject = studio.localStream;
      initWebRTC();
    }
  }, { immediate: true });
});

onUnmounted(() => {
  if (ws) ws.close();
  if (pc) pc.close();
});
</script>

<template>
  <div class="stage-wrapper">
    <!-- EL LIENZO 16:9 (ESTILO YOUTUBE) -->
    <div class="youtube-canvas" :class="`layout-${studio.layout}`">
      
      <!-- CONTENEDOR CÁMARA LOCAL -->
      <div class="video-slot local-slot" :class="{ 'cam-off': !studio.isCamOn }">
        <video ref="localVideo" autoplay playsinline muted class="video-el"></video>
        <div class="name-tag">{{ studio.userName }}</div>
      </div>

      <!-- CONTENEDOR CÁMARA INVITADO -->
      <div v-show="studio.guestConnected" class="video-slot guest-slot">
        <video ref="guestVideo" autoplay playsinline class="video-el"></video>
        <div class="name-tag">Invitado</div>
      </div>

      <!-- OVERLAYS -->
      <div class="absolute inset-0 pointer-events-none z-50">
        <img :src="studio.activeLogo" class="watermark" />
        <div class="live-badge">720p HD</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.stage-wrapper {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  padding: 2rem;
}

/* EL SECRETO: Un contenedor que siempre mantiene 16:9 y escala al padre */
.youtube-canvas {
  width: 100%;
  max-height: 100%;
  aspect-ratio: 16 / 9;
  background: #000;
  position: relative;
  overflow: hidden;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.video-slot {
  position: absolute;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  background: #0f172a;
}

.video-el {
  width: 100%;
  height: 100%;
  object-fit: cover; /* Para que se vean "como gente" */
}

.cam-off .video-el { opacity: 0; }

.name-tag {
  position: absolute;
  bottom: 10px;
  left: 10px;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 10px;
  font-weight: bold;
  color: white;
  border: 1px solid rgba(255,255,255,0.1);
  z-index: 20;
}

/* --- LOGICA DE LAYOUTS --- */

/* MODO SOLO */
.layout-solo .local-slot { width: 100%; height: 100%; left: 0; top: 0; }
.layout-solo .guest-slot { width: 0; height: 0; opacity: 0; }

/* MODO GRID (Pegado) */
.layout-grid .local-slot { width: 50%; height: 100%; left: 0; top: 0; border-right: 1px solid #000; }
.layout-grid .guest-slot { width: 50%; height: 100%; left: 50%; top: 0; opacity: 1; }

/* MODO ZOOM (Espaciado con fondo) */
.layout-zoom { background: radial-gradient(circle at center, #1e293b 0%, #020617 100%); }
.layout-zoom .local-slot { 
  width: 44%; height: 56%; left: 4%; top: 22%; 
  border-radius: 12px; border: 1px solid rgba(255,255,255,0.1);
}
.layout-zoom .guest-slot { 
  width: 44%; height: 56%; left: 52%; top: 22%; opacity: 1;
  border-radius: 12px; border: 1px solid rgba(255,255,255,0.1);
}

/* EXTRAS */
.watermark { position: absolute; top: 20px; right: 20px; height: 10%; opacity: 0.8; }
.live-badge { position: absolute; top: 20px; left: 20px; background: rgba(255,0,0,0.2); border: 1px solid rgba(255,0,0,0.4); color: #ff4444; padding: 2px 8px; border-radius: 4px; font-size: 9px; font-weight: bold; letter-spacing: 1px; }
</style>