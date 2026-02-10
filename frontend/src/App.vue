<script setup>
import { ref, onMounted } from 'vue'
import * as AppGo from './services/api.js'
import VideoCanvas from './components/VideoCanvas.vue'
import { useStudioStore } from './stores/studio'
import QrcodeVue from 'qrcode.vue'

const studio = useStudioStore()
const isTunnelActive = ref(false)
const isTunnelLoading = ref(false)
const guestLink = ref('')

const startStudio = async () => {
  await studio.initLocalStream();
}

const toggleTunnel = async () => {
  if (isTunnelActive.value) {
    await AppGo.ToggleTunnel()
    isTunnelActive.value = false
    return
  }
  isTunnelLoading.value = true
  try {
    const output = await AppGo.ToggleTunnel()
    if (output && output !== 'timeout') {
      const match = output.match(/([a-z0-9-]+\.lhr\.(life|pro|rocks))/)
      if (match) {
        guestLink.value = `https://${match[1]}/invitado`
        isTunnelActive.value = true
      }
    }
  } catch (err) {
    console.error(err)
  } finally {
    isTunnelLoading.value = false
  }
}
</script>

<template>
  <div class="flex flex-col h-screen bg-[#020617] text-slate-200 overflow-hidden font-sans">
    
    <!-- LOBBY -->
    <div v-if="!studio.isInitialized" class="fixed inset-0 z-[200] bg-slate-950 flex flex-col items-center justify-center p-6 text-center">
      <div class="w-20 h-20 bg-orange-500 rounded-3xl flex items-center justify-center shadow-2xl shadow-orange-500/20 mb-8">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"/></svg>
      </div>
      <h1 class="text-4xl font-poppins font-bold text-white mb-2">GPHR Studio</h1>
      <p class="text-slate-400 mb-8 max-w-md text-sm">Prepara tu cámara y micrófono antes de entrar.</p>
      
      <div class="w-full max-w-sm bg-slate-900 border border-white/5 rounded-2xl p-6 text-left">
        <label class="block text-[10px] font-bold uppercase text-slate-500 mb-2">Tu nombre</label>
        <input v-model="studio.userName" class="w-full bg-slate-800 border border-white/10 rounded-xl px-4 py-3 text-white focus:ring-2 ring-orange-500 outline-none mb-4" />
        <button @click="startStudio" class="w-full bg-orange-500 hover:bg-orange-600 text-white font-bold py-4 rounded-xl shadow-xl transition-all active:scale-95">ENTRAR AL ESTUDIO</button>
      </div>
    </div>

    <!-- HEADER -->
    <header class="h-14 border-b border-white/5 flex items-center justify-between px-6 bg-[#0f172a] shrink-0">
      <div class="flex items-center gap-4">
        <div class="w-8 h-8 bg-orange-500 rounded flex items-center justify-center font-bold text-white">CH</div>
        <input class="bg-transparent border-none focus:ring-0 font-bold text-sm text-white" v-model="studio.streamTitle" />
      </div>
      <div class="flex items-center gap-3">
        <div class="bg-red-600/10 text-red-500 px-3 py-1 rounded border border-red-500/20 text-[9px] font-bold uppercase tracking-tighter animate-pulse">Offline</div>
        <button class="bg-red-600 hover:bg-red-700 text-white px-6 py-2 rounded-md font-bold text-[10px] uppercase shadow-lg">En Vivo</button>
      </div>
    </header>

    <!-- MAIN BODY -->
    <main class="flex-1 flex overflow-hidden">
      <div class="flex-1 flex flex-col overflow-hidden relative">
        
        <!-- NOTIFICATION -->
        <Transition name="fade">
          <div v-if="studio.guestConnected" class="absolute top-6 left-1/2 -translate-x-1/2 z-[150] bg-teal-500 text-white px-6 py-2 rounded-full shadow-2xl font-bold text-xs flex items-center gap-2">
            <span class="w-2 h-2 bg-white rounded-full animate-ping"></span>
            ¡Invitado Conectado!
          </div>
        </Transition>

        <div v-if="isTunnelLoading" class="absolute inset-0 bg-slate-950/90 backdrop-blur-md z-[100] flex flex-col items-center justify-center text-center">
          <div class="w-12 h-12 border-4 border-teal-500 border-t-transparent rounded-full animate-spin mb-4"></div>
          <h3 class="font-poppins text-lg font-bold text-teal-500">Abriendo túnel...</h3>
        </div>

        <!-- 2. EL ESCENARIO (EL CONTENEDOR QUE ESCALA) -->
        <div class="flex-1 min-h-0 relative flex items-center justify-center bg-[#020617]">
          <VideoCanvas />

          <!-- Selector de Layouts Flotante -->
          <div class="absolute bottom-10 left-1/2 -translate-x-1/2 bg-white/90 backdrop-blur-md shadow-2xl rounded-2xl p-1.5 flex gap-1.5 z-50 text-slate-400">
            <button @click="studio.setLayout('solo')" :class="[studio.layout === 'solo' ? 'ring-2 ring-orange-500 bg-orange-50 text-orange-600' : '']" class="w-11 h-11 flex items-center justify-center rounded-xl transition-all">
              <div class="w-6 h-4 border-2 border-current rounded-sm"></div>
            </button>
            <button @click="studio.setLayout('grid')" :class="[studio.layout === 'grid' ? 'ring-2 ring-orange-500 bg-orange-50 text-orange-600' : '']" class="w-11 h-11 flex items-center justify-center rounded-xl transition-all">
              <div class="w-6 h-4 flex"><div class="flex-1 border-2 border-current"></div><div class="flex-1 border-2 border-current border-l-0"></div></div>
            </button>
            <button @click="studio.setLayout('zoom')" :class="[studio.layout === 'zoom' ? 'ring-2 ring-teal-500 bg-teal-50 text-teal-600' : '']" class="w-11 h-11 flex items-center justify-center rounded-xl transition-all">
              <div class="w-6 h-4 flex gap-1"><div class="flex-1 border-2 border-current rounded-[2px]"></div><div class="flex-1 border-2 border-current rounded-[2px]"></div></div>
            </button>
          </div>
        </div>

        <!-- 3. EL DOCK -->
        <div class="h-28 flex gap-4 overflow-x-auto p-4 shrink-0 bg-[#0f172a]/50 border-t border-white/5">
          <div :class="studio.isCamOn ? 'border-orange-500' : 'border-white/10'" class="min-w-[140px] bg-slate-800 rounded-xl border-2 relative overflow-hidden shadow-xl">
            <div class="absolute bottom-0 left-0 right-0 bg-black/60 p-1.5 text-[8px] font-bold flex justify-between items-center">
              <span class="truncate">{{ studio.userName }}</span>
              <span class="text-orange-500 uppercase">Tú</span>
            </div>
          </div>

          <div class="flex gap-3">
            <button @click="studio.toggleScreenShare()" class="min-w-[100px] bg-slate-900/50 border-2 border-dashed border-white/10 rounded-xl flex flex-col items-center justify-center gap-1 hover:border-orange-500 transition-all text-slate-500 hover:text-orange-500">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/></svg>
              <span class="text-[8px] font-bold uppercase">Presentar</span>
            </button>
            <button v-if="!isTunnelActive" @click="toggleTunnel" class="min-w-[100px] bg-slate-900/50 border-2 border-dashed border-white/10 rounded-xl flex flex-col items-center justify-center gap-1 hover:border-teal-500 transition-all text-slate-500 hover:text-teal-500">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M12 6v6m0 0v6m0-6h6m-6 0H6"/></svg>
              <span class="text-[8px] font-bold uppercase">Invitar</span>
            </button>
            <div v-if="isTunnelActive" class="min-w-[200px] bg-slate-900 border border-teal-500/30 rounded-xl flex p-2 gap-3 items-center">
               <div class="flex-1 flex flex-col justify-between py-1 text-left">
                  <span class="text-[7px] font-bold text-teal-500 uppercase tracking-tighter">Invitación</span>
                  <button @click="toggleTunnel" class="bg-teal-600 text-white text-[8px] font-bold py-1 rounded">Copiar Link</button>
               </div>
               <div class="bg-white p-1 rounded-lg shrink-0"><QrcodeVue :value="guestLink" :size="45" /></div>
            </div>
          </div>
        </div>
      </div>

      <!-- SIDEBAR -->
      <aside class="w-72 bg-[#0f172a] border-l border-white/5 shrink-0 p-5">
        <h2 class="text-xs font-bold uppercase text-slate-400 mb-6 tracking-widest">Marca</h2>
        <!-- Selector de Color Simple -->
        <div class="space-y-2">
          <p class="text-[9px] font-bold text-slate-500 uppercase">Color de Marca</p>
          <input type="color" v-model="studio.brandColor" class="w-full h-8 bg-transparent border-none cursor-pointer" />
        </div>
      </aside>
    </main>

    <!-- FOOTER -->
    <footer class="h-16 border-t border-white/5 bg-[#0f172a] flex items-center justify-center gap-6 px-6 shrink-0 shadow-2xl">
      <button @click="studio.toggleMic()" :class="[studio.isMicOn ? 'bg-slate-800' : 'bg-red-500/20 text-red-500']" class="w-10 h-10 rounded-xl border border-white/5 flex items-center justify-center transition-all hover:scale-105"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z"/></svg></button>
      <button @click="studio.toggleCam()" :class="[studio.isCamOn ? 'bg-slate-800' : 'bg-red-500/20 text-red-500']" class="w-10 h-10 rounded-xl border border-white/5 flex items-center justify-center transition-all hover:scale-105"><svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"/></svg></button>
      <div class="h-8 w-[1px] bg-white/10 mx-2"></div>
      <button class="bg-[#1e293b] hover:bg-red-600 text-white px-8 py-2 rounded-lg font-bold text-[10px] uppercase shadow-lg transition-all active:scale-95">Abandonar</button>
    </footer>
  </div>
</template>

<style>
html, body, #app { height: 100%; margin: 0; padding: 0; overflow: hidden; background: #020617; font-family: 'Inter', sans-serif; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.5s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>