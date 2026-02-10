<script setup>
import { useStudioStore } from '../stores/studio';
const studio = useStudioStore();

defineEmits(['toggle-screen-share', 'exit']);
</script>

<template>
  <footer class="h-20 bg-[#020617] border-t border-white/5 flex items-center justify-center gap-6 px-8 shrink-0 shadow-2xl z-50">
    <!-- Micrófono -->
    <button 
      @click="studio.toggleMic()" 
      :class="[studio.isMicOn ? 'bg-white/5 text-slate-300 hover:bg-white/10' : 'bg-red-500 text-white shadow-lg shadow-red-500/20']"
      class="w-14 h-14 rounded-2xl flex items-center justify-center transition-all border border-white/5"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </button>

    <!-- Cámara -->
    <button 
      @click="studio.toggleCam()" 
      :class="[studio.isCamOn ? 'bg-white/5 text-slate-300 hover:bg-white/10' : 'bg-red-500 text-white shadow-lg shadow-red-500/20']"
      class="w-14 h-14 rounded-2xl flex items-center justify-center transition-all border border-white/5"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </button>

    <div class="w-[1px] h-8 bg-white/10 mx-2"></div>

    <!-- Compartir Pantalla -->
    <button 
      @click="$emit('toggle-screen-share')"
      class="bg-white/5 text-slate-300 hover:bg-white/10 px-8 h-14 rounded-2xl flex items-center gap-3 transition-all border border-white/5 font-black text-[10px] uppercase tracking-[0.2em]"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      Compartir
    </button>

    <!-- Salir -->
    <button 
      @click="$emit('exit')"
      class="bg-transparent border-2 border-white/5 text-slate-600 hover:border-red-500/50 hover:text-red-500 px-10 h-14 rounded-2xl flex items-center gap-2 transition-all font-black text-[10px] uppercase tracking-[0.2em]"
    >
      Salir
    </button>
  </footer>
</template>
