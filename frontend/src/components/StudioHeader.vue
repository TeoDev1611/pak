<script setup>
import { useStudioStore } from '../stores/studio';
const studio = useStudioStore();

defineProps({
  isStreaming: Boolean
});

defineEmits(['toggle-stream']);
</script>

<template>
  <header class="h-16 bg-[#020617] border-b border-white/5 flex items-center justify-between px-8 shrink-0 z-50 shadow-2xl">
    <div class="flex items-center gap-8">
      <div class="flex items-center gap-3 group cursor-default">
        <div class="w-9 h-9 bg-orange-500 rounded-xl flex items-center justify-center shadow-lg shadow-orange-500/20 group-hover:scale-110 transition-transform">
          <span class="text-white font-black text-sm">P</span>
        </div>
        <h1 class="font-black text-white tracking-tighter text-2xl">PAK <span class="text-orange-500 text-lg">PRO</span></h1>
      </div>
      
      <div class="h-6 w-[1px] bg-white/10 mx-2"></div>
      
      <div class="flex items-center">
        <input 
          v-model="studio.streamTitle" 
          class="bg-white/5 border-none focus:ring-1 ring-orange-500/50 font-bold text-xs text-slate-300 w-72 rounded-xl px-4 py-2 transition-all placeholder:text-slate-600 outline-none"
          placeholder="Escribe el título de tu transmisión..."
        />
      </div>
    </div>

    <div class="flex items-center gap-6">
      <div 
        v-if="isStreaming"
        class="flex items-center gap-2 bg-red-500/10 px-4 py-2 rounded-full border border-red-500/20"
      >
        <span class="w-2 h-2 bg-red-500 rounded-full animate-pulse shadow-[0_0_10px_rgba(239,68,68,0.8)]"></span>
        <span class="text-[10px] font-black text-red-500 uppercase tracking-widest">Live</span>
      </div>
      
      <button 
        @click="$emit('toggle-stream')"
        :class="[
          isStreaming 
            ? 'bg-slate-800 hover:bg-slate-700 text-white shadow-xl border border-white/5' 
            : 'bg-orange-500 hover:bg-orange-600 text-white shadow-xl shadow-orange-500/20'
        ]"
        class="px-10 py-3 rounded-xl font-black text-[10px] uppercase tracking-widest transition-all active:scale-95"
      >
        {{ isStreaming ? 'Finalizar' : 'Emitir en Vivo' }}
      </button>
    </div>
  </header>
</template>
