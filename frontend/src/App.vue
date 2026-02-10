<script setup>
import { ref, onMounted, watch } from 'vue';
import * as AppGo from './services/api.js';
import { useStudioStore } from './stores/studio';

// Componentes
import VideoCanvas from './components/VideoCanvas.vue';
import StudioHeader from './components/StudioHeader.vue';
import StudioSidebar from './components/StudioSidebar.vue';
import StudioFooter from './components/StudioFooter.vue';
import UserCameraCard from './components/UserCameraCard.vue';

const studio = useStudioStore();
const videoCanvas = ref(null);

const isTunnelActive = ref(false);
const isTunnelLoading = ref(false);
const guestLink = ref('');
const isStreaming = ref(false);

const startStudio = async () => {
  await studio.initLocalStream();
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
  <div class="flex flex-col h-screen bg-[#020617] text-slate-200 overflow-hidden font-sans select-none">
    
    <!-- LOBBY -->
    <div v-if="!studio.isInitialized" class="fixed inset-0 z-[1000] bg-[#020617] flex flex-col items-center justify-center p-6 text-center">
       <div class="flex flex-col md:flex-row gap-12 items-center max-w-5xl w-full">
         <div class="flex-1 aspect-video bg-black rounded-[2rem] overflow-hidden border-8 border-slate-900 shadow-2xl relative">
           <video v-if="studio.localStream" :srcObject="studio.localStream" autoplay playsinline muted class="w-full h-full object-cover -scale-x-100"></video>
         </div>
         <div class="w-full max-w-sm text-left space-y-8">
           <h1 class="text-6xl font-black tracking-tighter text-white">PAK <span class="text-orange-500">PRO</span></h1>
           <div class="space-y-4">
             <input v-model="studio.userName" class="w-full bg-slate-900 border border-slate-800 text-white rounded-2xl px-6 py-4 font-bold outline-none ring-orange-500/20 focus:ring-4" />
             <button @click="startStudio" class="w-full bg-orange-500 text-white font-black py-5 rounded-2xl shadow-xl shadow-orange-500/30 uppercase text-xs tracking-widest hover:bg-orange-600 transition-all">Entrar al Estudio</button>
           </div>
         </div>
       </div>
    </div>

    <StudioHeader :is-streaming="isStreaming" @toggle-stream="toggleStreaming" />

    <div class="flex-1 flex overflow-hidden">
      <main class="flex-1 flex flex-col overflow-hidden bg-[#0f172a]">
        
        <!-- PROGRAM VIEW (Escenario) -->
        <div class="flex-1 min-h-0 relative flex items-center justify-center bg-black overflow-hidden group">
          <div class="w-full h-full relative flex items-center justify-center">
            <VideoCanvas ref="videoCanvas" class="w-full h-full" />
            
            <!-- Layout Selector Flotante -->
            <div class="absolute bottom-8 left-1/2 -translate-x-1/2 bg-slate-900/80 backdrop-blur-xl shadow-2xl rounded-2xl p-2 flex gap-1 z-40 border border-white/10 opacity-0 group-hover:opacity-100 transition-all">
              <button v-for="l in ['solo', 'grid', 'sidebar', 'pip']" :key="l" @click="studio.layout = l" 
                :class="[studio.layout === l ? 'bg-orange-500 text-white shadow-lg shadow-orange-500/20' : 'text-slate-400 hover:bg-white/5']"
                class="px-4 h-9 flex items-center justify-center rounded-xl transition-all uppercase text-[9px] font-black tracking-widest"
              >{{ l }}</button>
            </div>
          </div>
        </div>

        <!-- BACKSTAGE DOCK (Fuentes de Video) -->
        <div class="h-40 bg-[#020617] border-t border-white/5 flex flex-col shadow-2xl">
          <div class="px-8 py-2 border-b border-white/5 flex justify-between items-center">
            <span class="text-[9px] font-black text-slate-500 uppercase tracking-widest">Fuentes de Video</span>
            <div class="flex gap-4">
              <span class="text-[9px] font-bold text-orange-500 animate-pulse">{{ studio.guestConnected ? '1 Invitado Conectado' : 'Esperando Invitados' }}</span>
            </div>
          </div>
          <div class="flex-1 flex gap-4 overflow-x-auto p-4 px-10 items-center scrollbar-hide">
            <UserCameraCard 
              id="local"
              :stream="studio.localStream" 
              :name="studio.userName" 
              :is-local="true" 
            />
            
            <UserCameraCard 
              v-if="studio.guestConnected"
              id="guest"
              :stream="studio.guestStream"
              :name="'Invitado'" 
            />
            
            <div v-else class="min-w-[160px] aspect-video border-2 border-dashed border-white/5 rounded-xl flex flex-col items-center justify-center text-slate-600 gap-2 italic text-[10px] bg-white/[0.02] hover:bg-white/[0.04] transition-colors cursor-help">
               <svg class="w-5 h-5 opacity-20" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M12 4v16m8-8H4" stroke-width="2"/></svg>
               <span>Esperando invitado...</span>
            </div>
          </div>
        </div>
      </main>

      <StudioSidebar 
        class="border-l border-white/5"
        :is-tunnel-active="isTunnelActive" 
        :guest-link="guestLink"
        @toggle-tunnel="toggleTunnel" 
        @start-record="videoCanvas.startRecording()"
        @stop-record="videoCanvas.stopRecording()"
      />
    </div>

    <StudioFooter />

    <div v-if="isTunnelLoading" class="fixed inset-0 bg-black/90 backdrop-blur-xl z-[2000] flex flex-col items-center justify-center">
      <div class="w-12 h-12 border-4 border-orange-500 border-t-transparent rounded-full animate-spin mb-6"></div>
      <p class="text-[10px] font-black text-white uppercase tracking-widest animate-pulse">Configurando Túnel Seguro...</p>
    </div>
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@400;700;900&family=Bebas+Neue&family=Space+Grotesk:wght@400;700&family=Inter:wght@400;600;900&family=Montserrat:wght@400;700;900&display=swap');
html, body, #app { height: 100%; margin: 0; padding: 0; overflow: hidden; background: white; font-family: 'Inter', sans-serif; }
.scrollbar-hide::-webkit-scrollbar { display: none; }
</style>