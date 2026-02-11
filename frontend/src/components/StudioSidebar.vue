<script setup>
import { ref } from 'vue';
import { useStudioStore } from '../stores/studio';
import QrcodeVue from 'qrcode.vue';

const studio = useStudioStore();
const activeTab = ref('brand');

const newBannerText = ref('');
const newBannerSub = ref('');
const isNewMarquee = ref(false);
const isNewTicker = ref(false);

defineProps({
  isTunnelActive: Boolean,
  guestLink: String
});

defineEmits(['toggle-tunnel']);

const handleFileUpload = (event, type) => {
  const file = event.target.files[0];
  if (!file) return;
  const reader = new FileReader();
  reader.onload = (e) => {
    if (type === 'logo') studio.setLogo(e.target.result);
    if (type === 'bg') studio.setBackground(e.target.result);
  };
  reader.readAsDataURL(file);
};

const addNewBanner = () => {
  if (newBannerText.value) {
    studio.addBanner(newBannerText.value, newBannerSub.value, isNewTicker.value, 'custom', isNewMarquee.value);
    newBannerText.value = '';
    newBannerSub.value = '';
    isNewMarquee.value = false;
    isNewTicker.value = false;
  }
};
</script>

<template>
  <aside class="w-80 bg-black border-l border-neutral-800 shrink-0 flex flex-col h-full overflow-hidden">
    <!-- Pestañas Minimalistas Industriales -->
    <div class="flex border-b border-neutral-800">
      <button v-for="t in ['brand', 'banners', 'record']" :key="t"
        @click="activeTab = t"
        :class="[activeTab === t ? 'text-white border-white' : 'text-neutral-500 border-transparent hover:text-neutral-300']"
        class="flex-1 py-3 text-[9px] font-bold uppercase tracking-[0.2em] border-b-2 transition-all"
      >
        {{ t === 'brand' ? 'Marca' : t === 'banners' ? 'Capas' : 'Rec' }}
      </button>
    </div>

    <div class="p-4 space-y-8 flex-1 overflow-y-auto scrollbar-hide">
      
      <!-- SECCIÓN: MARCA -->
      <div v-if="activeTab === 'brand'" class="space-y-8 animate-in fade-in duration-300">
        <div class="space-y-4">
          <label class="text-[9px] font-bold text-neutral-500 uppercase tracking-widest">Acento de Color</label>
          <div class="grid grid-cols-5 gap-2">
            <button v-for="color in ['#f97316', '#3b82f6', '#10b981', '#ef4444', '#ffffff', '#a855f7', '#ec4899', '#06b6d4', '#84cc16', '#eab308']" :key="color"
              @click="studio.brandColor = color"
              :style="{ backgroundColor: color }"
              class="w-full aspect-square rounded-sm border transition-all"
              :class="[studio.brandColor === color ? 'border-white scale-110 shadow-lg shadow-white/20' : 'border-neutral-800']"
            ></button>
          </div>
        </div>

        <div class="space-y-4">
          <label class="text-[9px] font-bold text-neutral-500 uppercase tracking-widest">Logo del Canal</label>
          <div class="space-y-4">
            <div class="flex gap-2">
              <label class="flex-1 flex items-center justify-center py-2 border border-neutral-800 rounded-sm hover:bg-neutral-900 transition-colors text-[9px] font-bold text-neutral-400 uppercase cursor-pointer">
                Subir Logo
                <input type="file" class="hidden" accept="image/*" @change="e => handleFileUpload(e, 'logo')" />
              </label>
              <button @click="studio.setLogo(null)" class="px-3 border border-neutral-800 rounded-sm hover:bg-neutral-900 text-[9px] font-bold text-neutral-600 uppercase">×</button>
            </div>
            
            <div v-if="studio.logoUrl" class="space-y-4">
              <div class="grid grid-cols-2 gap-2">
                <button v-for="pos in ['top-left', 'top-right', 'bottom-left', 'bottom-right']" :key="pos"
                  @click="studio.logoPosition = pos"
                  :class="[studio.logoPosition === pos ? 'bg-white text-black' : 'bg-neutral-900 text-neutral-500 border-neutral-800']"
                  class="py-1.5 border rounded-sm text-[8px] font-bold uppercase"
                >
                  {{ pos.replace('-', ' ') }}
                </button>
              </div>
              <input type="range" min="0.05" max="0.5" step="0.01" v-model.number="studio.logoScale" class="w-full h-1 bg-neutral-800 rounded-lg appearance-none cursor-pointer accent-white" />
            </div>
          </div>
        </div>

        <div class="space-y-4 pt-4 border-t border-neutral-800">
           <label class="text-[9px] font-bold text-neutral-400 uppercase tracking-widest block mb-4">Visibilidad en Pantalla</label>
           
           <div class="flex items-center justify-between">
              <span class="text-[10px] text-neutral-300 uppercase font-bold tracking-widest">Indicador LIVE/REC</span>
              <button @click="studio.showLiveBadge = !studio.showLiveBadge" class="w-8 h-4 rounded-full border border-neutral-700 relative transition-colors" :class="studio.showLiveBadge ? 'bg-orange-500' : 'bg-black'">
                <div class="absolute top-0.5 w-2.5 h-2.5 rounded-full transition-all" :class="studio.showLiveBadge ? 'bg-white right-0.5' : 'bg-neutral-700 left-0.5'"></div>
              </button>
           </div>

           <div class="flex items-center justify-between">
              <span class="text-[10px] text-neutral-300 uppercase font-bold tracking-widest">Nombres Participantes</span>
              <button @click="studio.showParticipantNames = !studio.showParticipantNames" class="w-8 h-4 rounded-full border border-neutral-700 relative transition-colors" :class="studio.showParticipantNames ? 'bg-orange-500' : 'bg-black'">
                <div class="absolute top-0.5 w-2.5 h-2.5 rounded-full transition-all" :class="studio.showParticipantNames ? 'bg-white right-0.5' : 'bg-neutral-700 left-0.5'"></div>
              </button>
           </div>
        </div>
      </div>

      <!-- SECCIÓN: CAPAS (BANNERS) -->
      <div v-if="activeTab === 'banners'" class="space-y-6 animate-in slide-in-from-right duration-300">
        
        <div class="p-4 bg-neutral-900 border border-neutral-700 rounded-sm space-y-4">
           <label class="text-[9px] font-bold text-neutral-400 uppercase tracking-widest">Ajustes Globales</label>
           <div class="space-y-3">
              <div class="flex flex-col gap-1.5">
                <span class="text-[8px] text-neutral-500 uppercase font-bold">Fuente</span>
                <select v-model="studio.bannerFont" class="w-full bg-black border border-neutral-700 rounded-sm px-2 py-1.5 text-[10px] text-white outline-none focus:border-neutral-500 font-bold uppercase tracking-widest">
                  <option value="Inter">Inter (Sleek)</option>
                  <option value="Roboto">Roboto (Clean)</option>
                  <option value="Oswald">Oswald (Impact)</option>
                  <option value="Courier Prime">Courier (Retro)</option>
                </select>
              </div>

              <div class="flex flex-col gap-1.5">
                <span class="text-[8px] text-neutral-500 uppercase font-bold">Color de Fondo</span>
                <div class="flex gap-2 items-center">
                  <input type="color" v-model="studio.bannerBgColor" class="w-8 h-8 bg-transparent border-none cursor-pointer" />
                  <input type="text" v-model="studio.bannerBgColor" class="flex-1 bg-black border border-neutral-700 rounded-sm px-2 py-1 text-[10px] text-white font-mono" />
                </div>
              </div>

              <div class="flex flex-col gap-1.5">
                <span class="text-[8px] text-neutral-500 uppercase font-bold">Color de Texto</span>
                <div class="flex gap-2 items-center">
                  <input type="color" v-model="studio.bannerTextColor" class="w-8 h-8 bg-transparent border-none cursor-pointer" />
                  <input type="text" v-model="studio.bannerTextColor" class="flex-1 bg-black border border-neutral-700 rounded-sm px-2 py-1 text-[10px] text-white font-mono" />
                </div>
              </div>

              <div class="flex flex-col gap-1.5">
                <div class="flex justify-between text-[8px] font-bold text-neutral-500 uppercase"><span>Redondeado</span><span>{{ studio.bannerBorderRadius }}px</span></div>
                <input type="range" min="0" max="50" v-model.number="studio.bannerBorderRadius" class="w-full h-1 bg-neutral-800 appearance-none cursor-pointer accent-white" />
              </div>
           </div>
        </div>

        <!-- Controles de Posición -->
        <div class="p-4 bg-neutral-900 border border-neutral-700 rounded-sm space-y-4">
           <div class="flex justify-between items-center">
             <label class="text-[9px] font-bold text-neutral-400 uppercase tracking-widest">Posicionamiento</label>
             <button @click="studio.bannerPositionX = 50; studio.bannerPositionY = 85" class="text-[8px] text-neutral-500 hover:text-white uppercase font-bold">Reset</button>
           </div>
           <div class="space-y-3">
              <div class="flex justify-between text-[8px] font-bold text-neutral-500 uppercase"><span>Vertical</span><span>{{ studio.bannerPositionY }}%</span></div>
              <input type="range" min="0" max="100" v-model.number="studio.bannerPositionY" class="w-full h-1 bg-neutral-800 appearance-none cursor-pointer accent-orange-500" />
              <div class="flex justify-between text-[8px] font-bold text-neutral-500 uppercase"><span>Horizontal</span><span>{{ studio.bannerPositionX }}%</span></div>
              <input type="range" min="0" max="100" v-model.number="studio.bannerPositionX" class="w-full h-1 bg-neutral-800 appearance-none cursor-pointer accent-orange-500" />
           </div>
        </div>

        <div class="space-y-3">
          <div v-for="banner in studio.banners" :key="banner.id" 
            @click="studio.toggleBanner(banner.id)"
            class="p-4 rounded-sm border transition-all cursor-pointer relative group"
            :class="[banner.active ? 'border-white bg-neutral-900' : 'border-neutral-800 bg-black hover:bg-neutral-950']"
          >
            <div class="font-bold text-[11px] text-white">{{ banner.text }}</div>
            <div class="text-[9px] text-neutral-500 mt-0.5" v-if="banner.subtext">{{ banner.subtext }}</div>
            <div class="flex gap-2 mt-2">
              <div v-if="banner.isMarquee" class="text-[7px] text-neutral-400 font-bold uppercase border border-neutral-800 px-1">Animado</div>
              <div v-if="banner.isTicker" class="text-[7px] text-neutral-400 font-bold uppercase border border-neutral-800 px-1">Ticker</div>
            </div>
            
            <button @click.stop="studio.removeBanner(banner.id)" class="absolute top-2 right-2 text-neutral-700 hover:text-white transition-colors">×</button>
          </div>
        </div>

        <div class="p-4 border border-neutral-700 rounded-sm space-y-4 bg-black">
          <label class="text-[9px] font-bold text-neutral-400 uppercase tracking-widest text-center block">Nuevo Banner</label>
          <input v-model="newBannerText" placeholder="Texto principal..." class="w-full bg-neutral-900 border border-neutral-700 rounded-sm px-3 py-2 text-xs text-white outline-none focus:border-neutral-500" />
          <input v-model="newBannerSub" placeholder="Subtexto (opcional)..." class="w-full bg-neutral-900 border border-neutral-700 rounded-sm px-3 py-2 text-xs text-white outline-none focus:border-neutral-500" />
          
          <div class="flex flex-col gap-3 pt-2">
            <div class="flex items-center gap-3">
              <button @click="isNewMarquee = !isNewMarquee" class="w-8 h-4 rounded-sm border border-neutral-700 relative" :class="isNewMarquee ? 'bg-orange-500' : 'bg-black'">
                <div class="absolute top-0.5 w-2.5 h-2.5 transition-all" :class="isNewMarquee ? 'bg-white right-0.5' : 'bg-neutral-700 left-0.5'"></div>
              </button>
              <label class="text-[9px] font-bold text-neutral-400 uppercase">Movimiento (Marquee)</label>
            </div>
            <div class="flex items-center gap-3">
              <button @click="isNewTicker = !isNewTicker" class="w-8 h-4 rounded-sm border border-neutral-700 relative" :class="isNewTicker ? 'bg-orange-500' : 'bg-black'">
                <div class="absolute top-0.5 w-2.5 h-2.5 transition-all" :class="isNewTicker ? 'bg-white right-0.5' : 'bg-neutral-700 left-0.5'"></div>
              </button>
              <label class="text-[9px] font-bold text-neutral-400 uppercase">Modo Ticker</label>
            </div>
          </div>
          
          <button @click="addNewBanner" class="w-full bg-white text-black text-[9px] font-bold py-2.5 rounded-sm uppercase tracking-widest hover:bg-neutral-200 transition-colors">Añadir al Inventario</button>
        </div>
      </div>

      <!-- SECCIÓN: GRABACIÓN -->
      <div v-if="activeTab === 'record'" class="space-y-6 animate-in fade-in duration-300">
        <div class="p-4 border border-neutral-800 rounded-sm space-y-4">
          <div class="flex items-center gap-2">
            <div class="w-1.5 h-1.5 rounded-full" :class="studio.isRecording ? 'bg-red-600 animate-pulse' : 'bg-neutral-800'"></div>
            <h3 class="text-[10px] font-bold text-white uppercase tracking-widest">Grabador</h3>
          </div>
          <input v-model="studio.recordingTitle" class="w-full bg-transparent border border-neutral-800 rounded-sm px-3 py-2 text-xs text-white outline-none focus:border-neutral-600" />
          <button @click="$emit(studio.isRecording ? 'stop-record' : 'start-record')"
            :class="studio.isRecording ? 'bg-white text-black' : 'bg-red-600 text-white'"
            class="w-full py-2.5 rounded-sm font-bold text-[10px] uppercase tracking-widest transition-colors"
          >
            {{ studio.isRecording ? 'Terminar' : 'Iniciar' }}
          </button>
        </div>
      </div>

    </div>
  </aside>
</template>