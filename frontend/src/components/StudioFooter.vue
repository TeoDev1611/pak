<script setup>
import { useStudioStore } from '../stores/studio';
const studio = useStudioStore();

defineEmits(['toggle-screen-share', 'exit']);
</script>

<template>
  <footer class="h-16 bg-black border-t border-neutral-800 flex items-center justify-center gap-4 px-8 shrink-0 z-50">
    <!-- Micrófono -->
    <button 
      @click="studio.toggleMic()" 
      :class="[studio.isMicOn ? 'bg-black text-neutral-400 border-neutral-800 hover:border-neutral-600 hover:text-white' : 'bg-red-600 text-white border-red-600']"
      class="w-10 h-10 rounded-md flex items-center justify-center transition-all border"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </button>

    <!-- Cámara -->
    <button 
      @click="studio.toggleCam()" 
      :class="[studio.isCamOn ? 'bg-black text-neutral-400 border-neutral-800 hover:border-neutral-600 hover:text-white' : 'bg-red-600 text-white border-red-600']"
      class="w-10 h-10 rounded-md flex items-center justify-center transition-all border"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </button>

    <div class="w-[1px] h-6 bg-neutral-800 mx-2"></div>

    <!-- Compartir Pantalla -->
    <button 
      @click="$emit('toggle-screen-share')"
      class="bg-black border border-neutral-800 text-neutral-400 hover:text-white hover:border-neutral-600 px-6 h-10 rounded-md flex items-center gap-2 transition-all font-bold text-[10px] uppercase tracking-widest"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      Compartir
    </button>

    <!-- Invitar -->
    <button 
      @click="$emit('open-invite')"
      class="bg-black border border-neutral-800 text-neutral-400 hover:text-orange-500 hover:border-orange-900/50 px-6 h-10 rounded-md flex items-center gap-2 transition-all font-bold text-[10px] uppercase tracking-widest"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      Invitar
    </button>

    <!-- Salir -->
    <button 
      @click="$emit('exit')"
      class="bg-black border border-neutral-800 text-neutral-600 hover:text-red-500 hover:border-red-900 px-6 h-10 rounded-md flex items-center gap-2 transition-all font-bold text-[10px] uppercase tracking-widest"
    >
      Salir
    </button>
  </footer>
</template>