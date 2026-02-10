<script setup>
import { computed } from 'vue';
import { useStudioStore } from '../stores/studio';

const studio = useStudioStore();
const props = defineProps({
  id: String,
  stream: Object,
  name: String,
  isLocal: Boolean
});

const isOnStage = computed(() => studio.participantsOnStage.includes(props.id));
const isHidden = computed(() => props.isLocal && !studio.isCamOn);
</script>

<template>
  <div 
    class="min-w-[120px] max-w-[120px] aspect-video bg-slate-900 rounded-md border-2 transition-all overflow-hidden relative shadow-sm group cursor-pointer"
    :class="[isOnStage ? 'border-orange-500' : 'border-slate-200 hover:border-slate-300']"
    @click="studio.toggleParticipant(id)"
  >
    <!-- Video Preview -->
    <video 
      v-if="stream && !isHidden" 
      :srcObject="stream" 
      autoplay playsinline muted
      class="w-full h-full object-cover transition-opacity"
      :class="[isOnStage ? 'opacity-100' : 'opacity-40 group-hover:opacity-70', { '-scale-x-100': isLocal }]"
    ></video>
    
    <div v-else class="absolute inset-0 flex items-center justify-center bg-slate-800 text-slate-600">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" stroke-width="2"/></svg>
    </div>

    <!-- Badge de Estado Minimalista -->
    <div class="absolute top-1 left-1 px-1 py-0.5 rounded-[3px] text-[7px] font-black uppercase tracking-tighter shadow-sm"
      :class="isOnStage ? 'bg-orange-500 text-white' : 'bg-slate-700 text-slate-300'"
    >
      {{ isOnStage ? 'En Vivo' : 'Off' }}
    </div>

    <!-- Overlay de Acción Sutil -->
    <div class="absolute inset-0 flex items-center justify-center bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity">
       <div class="bg-white/90 backdrop-blur-sm text-slate-900 px-2 py-0.5 rounded text-[8px] font-black uppercase shadow-lg">
         {{ isOnStage ? 'Quitar' : 'Subir' }}
       </div>
    </div>

    <div class="absolute bottom-0 left-0 right-0 bg-black/40 p-0.5 text-[7px] font-bold text-white truncate px-1">
      {{ name }}
    </div>
  </div>
</template>
