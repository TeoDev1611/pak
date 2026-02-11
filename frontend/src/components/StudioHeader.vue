<script setup>
import { useStudioStore } from '../stores/studio';
const studio = useStudioStore();

defineProps({
  isStreaming: Boolean
});

defineEmits(['toggle-stream']);
</script>

<template>
  <header class="h-14 bg-black border-b border-neutral-800 flex items-center justify-between px-4 shrink-0 z-50">
    <div class="flex items-center gap-4">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 bg-white text-black font-bold flex items-center justify-center rounded-md text-xs">P</div>
        <div class="h-6 w-[1px] bg-neutral-800"></div>
        <span class="text-sm font-medium text-white">Chasqui <span class="text-neutral-500">Pro Studio</span></span>
      </div>
      
      <div class="flex items-center ml-4">
        <input 
          v-model="studio.streamTitle" 
          class="bg-transparent border border-neutral-800 focus:border-neutral-600 font-medium text-xs text-neutral-300 w-64 rounded-md px-3 py-1.5 transition-all placeholder:text-neutral-700 outline-none"
          placeholder="Título del proyecto..."
        />
      </div>
    </div>

    <div class="flex items-center gap-4">
      <div 
        v-if="isStreaming"
        class="flex items-center gap-2 bg-red-600 px-2 py-0.5 rounded-sm"
      >
        <span class="text-[9px] font-bold text-white uppercase tracking-wider">LIVE</span>
      </div>
      
      <button 
        @click="$emit('toggle-stream')"
        :class="[
          isStreaming 
            ? 'bg-neutral-900 border border-neutral-800 text-neutral-400 hover:text-white' 
            : 'bg-white text-black hover:bg-neutral-200'
        ]"
        class="px-5 h-8 rounded-md font-bold text-[10px] uppercase tracking-widest transition-colors"
      >
        {{ isStreaming ? 'Detener' : 'Transmitir' }}
      </button>
    </div>
  </header>
</template>