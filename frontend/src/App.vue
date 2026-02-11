<script setup>
import { ref, onMounted, watch, computed } from 'vue';
import * as AppGo from './services/api.js';
import { useStudioStore } from './stores/studio';

// Componentes
import VideoCanvas from './components/VideoCanvas.vue';
import StudioHeader from './components/StudioHeader.vue';
import StudioSidebar from './components/StudioSidebar.vue';
import StudioFooter from './components/StudioFooter.vue';
import UserCameraCard from './components/UserCameraCard.vue';
import QrcodeVue from 'qrcode.vue';

const studio = useStudioStore();
const videoCanvas = ref(null);

const isTunnelActive = ref(false);
const isTunnelLoading = ref(false);
const guestLink = ref('');
const isStreaming = ref(false);
const showInviteModal = ref(false);
const copyNotification = ref(false);

// LÓGICA DE BANNERS (OVERLAY HTML)
const activeBanner = computed(() => studio.banners.find(b => b.active));

const bannerClasses = computed(() => {
  if (!activeBanner.value) return "";
  return "absolute px-8 py-4 z-30 flex flex-col items-center justify-center text-center max-w-[90%] pointer-events-none select-none shadow-2xl overflow-hidden";
});

const bannerPositionStyle = computed(() => {
    if (!activeBanner.value) return {};
    const b = activeBanner.value;
    return {
        left: studio.bannerPositionX + '%',
        top: studio.bannerPositionY + '%',
        transform: 'translate(-50%, -50%)',
        backgroundColor: b.bgColor || studio.bannerBgColor,
        color: b.textColor || studio.bannerTextColor,
        borderRadius: (b.borderRadius ?? studio.bannerBorderRadius) + 'px',
        fontFamily: studio.bannerFont,
        border: `1px solid rgba(255,255,255,0.1)`
    }
});

const startStudio = async () => {
  await studio.initLocalStream();
};

const copyLink = () => {
  navigator.clipboard.writeText(guestLink.value || `http://${window.location.hostname}:8080/invitado`);
  copyNotification.value = true;
  setTimeout(() => copyNotification.value = false, 2000);
};

const toggleStreaming = async () => {
  if (isStreaming.value) {
    await AppGo.StopStream();
    isStreaming.value = false;
    return;
  }
  try {
    await AppGo.StartStream(studio.rtmpUrl);
    if (videoCanvas.value) await videoCanvas.value.connectToStudio();
    isStreaming.value = true;
  } catch (err) {
    console.error(err);
  }
};

const toggleTunnel = async () => {
  if (isTunnelActive.value) {
    await AppGo.ToggleTunnel();
    isTunnelActive.value = false;
    return;
  }
  isTunnelLoading.value = true;
  try {
    const output = await AppGo.ToggleTunnel();
    if (output && output !== 'timeout') {
      const match = output.match(/([a-z0-9-]+\.lhr\.(life|pro|rocks))/);
      if (match) {
        guestLink.value = `https://${match[1]}/invitado`;
        isTunnelActive.value = true;
      }
    }
  } finally {
    isTunnelLoading.value = false;
  }
};

onMounted(() => { studio.initLocalStream(); });
</script>

<template>
  <div class="flex flex-col h-screen bg-black text-white font-sans antialiased selection:bg-orange-500 selection:text-white">
    
    <!-- LOBBY (Minimalista) -->
    <div v-if="!studio.isInitialized" class="fixed inset-0 z-[1000] bg-black flex flex-col items-center justify-center p-6 text-center">
       <div class="flex flex-col md:flex-row gap-12 items-center max-w-5xl w-full">
         <div class="flex-1 aspect-video bg-[#0A0A0A] border border-neutral-800 rounded-md overflow-hidden relative">
           <video 
             v-if="studio.localStream" 
             :srcObject="studio.localStream" 
             autoplay 
             playsinline 
             muted 
             class="w-full h-full object-cover -scale-x-100"
             @mounted="(el) => el.srcObject = studio.localStream"
           ></video>
           <div v-else class="w-full h-full flex items-center justify-center">
              <div class="w-8 h-8 border-2 border-neutral-800 border-t-white rounded-full animate-spin"></div>
           </div>
         </div>
         <div class="w-full max-w-sm text-left space-y-6">
           <div class="flex items-center gap-4">
             <div class="w-12 h-12 bg-white text-black font-bold flex items-center justify-center rounded-md text-xl">P</div>
             <h1 class="text-4xl font-bold tracking-tighter text-white uppercase">Chasqui <span class="text-neutral-500">Pro</span></h1>
           </div>
           <div class="space-y-4">
             <input v-model="studio.userName" class="w-full bg-black border border-neutral-800 text-white rounded-md px-4 py-3 text-sm focus:border-white outline-none transition-colors" placeholder="Tu nombre..." />
             <button @click="startStudio" class="w-full bg-white text-black font-bold py-3 rounded-md text-sm hover:bg-neutral-200 transition-colors uppercase tracking-widest">Entrar al Estudio</button>
           </div>
         </div>
       </div>
    </div>

    <StudioHeader :is-streaming="isStreaming" @toggle-stream="toggleStreaming" />

    <div class="flex-1 flex overflow-hidden">
      <main class="flex-1 flex flex-col min-h-0 overflow-hidden bg-[#050505]">
        
        <!-- ZONA DE ESCENARIO: MAXIMIZADO (FULL BLEED) -->
        <div class="flex-1 relative bg-[#050505] overflow-hidden flex items-center justify-center">
                      <div class="aspect-video h-full w-full max-w-full max-h-full relative shadow-2xl bg-black flex items-center justify-center">
                      <VideoCanvas ref="videoCanvas" class="w-full h-full" />
                      
                      <!-- Badge LIVE Minimalista -->            <div v-if="isStreaming && studio.showLiveBadge" class="absolute top-4 left-4 flex items-center gap-2 bg-black/90 border border-white/10 px-3 py-1.5 rounded-md backdrop-blur-md z-10">
                <div class="w-2 h-2 rounded-full bg-red-500 animate-pulse"></div>
                <span class="text-[10px] font-bold tracking-[0.2em] text-white">STREAMING LIVE</span>
            </div>

            <!-- Badge REC -->
            <div v-if="studio.isRecording && studio.showLiveBadge" class="absolute top-4 right-4 flex items-center gap-2 bg-red-600/90 border border-white/10 px-3 py-1.5 rounded-md backdrop-blur-md z-10">
                <div class="w-2 h-2 rounded-full bg-white animate-pulse"></div>
                <span class="text-[10px] font-bold tracking-[0.2em] text-white uppercase">Grabando</span>
            </div>
          </div>
        </div>

        <!-- DOCK DE CONTROL INTEGRADO (LAYOUTS) -->
        <div class="h-10 border-t border-neutral-800 bg-black flex items-center justify-center gap-1 px-4 shrink-0 z-20">
          <button v-for="l in ['solo', 'grid', 'sidebar', 'pip']" :key="l" @click="studio.layout = l" 
            :class="[studio.layout === l ? 'bg-neutral-800 text-white' : 'text-neutral-500 hover:text-white hover:bg-neutral-900']"
            class="px-4 h-6 flex items-center justify-center rounded-sm transition-all uppercase text-[8px] font-bold tracking-widest border border-transparent hover:border-neutral-800"
          >{{ l }}</button>
        </div>

        <!-- BARRA DE ESTADO Y GUEST LIST (COMPACTA) -->
        <div class="h-14 bg-neutral-900 border-t border-neutral-800 flex items-center px-4 justify-between shrink-0 z-30">
          
          <!-- Izquierda: Host Info -->
          <div class="flex items-center gap-4">
             <div class="flex items-center gap-2 px-3 py-1.5 bg-black rounded-md border border-neutral-800">
                <div class="w-2 h-2 rounded-full" :class="studio.isMicOn ? 'bg-green-500' : 'bg-red-500'"></div>
                <span class="text-[10px] font-bold text-neutral-300 uppercase tracking-widest">{{ studio.userName }} (Host)</span>
             </div>
          </div>

          <!-- Centro: Guest List (Chips) -->
          <div class="flex items-center gap-2">
             <span class="text-[9px] font-bold text-neutral-600 uppercase tracking-widest mr-2">Participantes:</span>
             
             <!-- Chip de Invitado -->
             <div v-if="studio.guestConnected" class="flex items-center gap-3 px-3 py-1.5 bg-neutral-800 border border-neutral-700 rounded-full animate-in fade-in slide-in-from-bottom-2">
                <span class="text-[10px] font-bold text-white uppercase">👤 Invitado</span>
                <div class="w-1.5 h-1.5 rounded-full bg-green-500"></div>
                <button @click="studio.guestConnected = false" class="text-neutral-500 hover:text-white transition-colors ml-1 font-bold text-xs">×</button>
             </div>
             
             <div v-else class="text-[10px] text-neutral-600 italic font-medium">No hay invitados conectados</div>
          </div>

          <!-- Derecha: Acciones Rápidas -->
          <div class="flex items-center gap-3">
             <button @click="showInviteModal = true" class="px-4 py-2 bg-white text-black text-[9px] font-bold uppercase tracking-widest rounded-md hover:bg-neutral-200 transition-colors">
               Invitar
             </button>
          </div>

        </div>
      </main>

      <StudioSidebar 
        :is-tunnel-active="isTunnelActive" 
        :guest-link="guestLink"
        @toggle-tunnel="toggleTunnel" 
        @start-record="videoCanvas.startRecording()"
        @stop-record="videoCanvas.stopRecording()"
      />
    </div>

    <StudioFooter @open-invite="showInviteModal = true" />

    <!-- MODAL DE INVITACIÓN -->
    <div v-if="showInviteModal" class="fixed inset-0 z-[3000] flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/80 backdrop-blur-md" @click="showInviteModal = false"></div>
      
      <div class="relative bg-neutral-900 border border-neutral-700 w-full max-w-md rounded-xl overflow-hidden shadow-2xl animate-in zoom-in-95 duration-200">
        <div class="p-6 space-y-6">
          <div class="flex justify-between items-center">
            <h3 class="text-xl font-bold tracking-tight text-white">Invitar Participante</h3>
            <button @click="showInviteModal = false" class="text-neutral-500 hover:text-white text-2xl">&times;</button>
          </div>

          <div class="flex flex-col items-center gap-6 py-4">
            <div class="p-4 bg-white rounded-lg shadow-inner">
               <QrcodeVue :value="guestLink || `http://${window.location.hostname}:8080/invitado`" :size="200" level="H" />
            </div>
            <p class="text-[10px] text-neutral-400 uppercase font-bold tracking-[0.2em] text-center">Escanea para unirte desde el móvil</p>
          </div>

          <div class="space-y-3">
            <label class="text-[9px] font-bold text-neutral-500 uppercase tracking-widest">Enlace Directo</label>
            <div class="flex gap-2">
              <input readonly :value="guestLink || `http://${window.location.hostname}:8080/invitado`" class="flex-1 bg-black border border-neutral-700 rounded-md px-4 py-3 text-xs text-white outline-none focus:border-neutral-500" />
              <button @click="copyLink" class="bg-white text-black px-6 rounded-md font-bold text-[10px] uppercase tracking-widest hover:bg-neutral-200 transition-colors">
                {{ copyNotification ? '¡Copiado!' : 'Copiar' }}
              </button>
            </div>
          </div>
        </div>
        
        <div class="bg-neutral-900/50 p-4 border-t border-neutral-800 flex items-center gap-3">
           <div class="w-2 h-2 rounded-full bg-green-500"></div>
           <span class="text-[9px] font-bold text-neutral-400 uppercase tracking-widest">Estudio listo para recibir invitados</span>
        </div>
      </div>
    </div>

    <div v-if="isTunnelLoading" class="fixed inset-0 bg-black/90 z-[2000] flex flex-col items-center justify-center">
      <div class="w-8 h-8 border-2 border-white border-t-transparent rounded-full animate-spin mb-4"></div>
      <p class="text-[10px] font-bold text-white uppercase tracking-widest">Procesando Túnel...</p>
    </div>
  </div>
</template>

<style>

@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800;900&display=swap');

html, body, #app { height: 100%; margin: 0; padding: 0; overflow: hidden; background: black; font-family: 'Inter', sans-serif; }

.scrollbar-hide::-webkit-scrollbar { display: none; }



@keyframes marquee {

  0% { transform: translateX(100%); }

  100% { transform: translateX(-100%); }

}

.animate-marquee {

  animation: marquee 15s linear infinite;

}

</style>
