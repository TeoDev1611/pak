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
    class="min-w-[140px] aspect-video bg-black rounded-sm border transition-all overflow-hidden relative group cursor-pointer"
    :class="[isOnStage ? 'border-white' : 'border-neutral-800 hover:border-neutral-600']"
    @click="studio.toggleParticipant(id)"
  >
    <!-- Video Preview -->
    <video 
      v-if="stream && !isHidden" 
      :srcObject="stream" 
      autoplay playsinline muted
      class="w-full h-full object-cover transition-opacity"
      :class="[isOnStage ? 'opacity-100' : 'opacity-30 group-hover:opacity-50', { '-scale-x-100': isLocal }]"
    ></video>
    
    <div v-else class="absolute inset-0 flex items-center justify-center bg-neutral-900 text-neutral-700">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" stroke-width="2"/></svg>
    </div>

    <!-- Badge de Estado Industrial -->
    <div class="absolute top-1.5 left-1.5 px-1 py-0.5 rounded-[1px] text-[7px] font-bold uppercase tracking-widest"
      :class="isOnStage ? 'bg-white text-black' : 'bg-neutral-800 text-neutral-500'"
    >
      {{ isOnStage ? 'LIVE' : 'IDLE' }}
    </div>

    <!-- Label Nombre -->
    <div class="absolute bottom-1.5 left-1.5 right-1.5 flex justify-between items-center">
      <span class="text-[7px] font-bold text-white uppercase tracking-wider truncate max-w-[80%] drop-shadow-md">{{ name }}</span>
      <div v-if="props.isLocal" class="w-1.5 h-1.5 rounded-full bg-blue-500 shadow-[0_0_5px_rgba(59,130,246,0.5)]"></div>
    </div>
  </div>
</template>