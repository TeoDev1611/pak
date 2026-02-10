<script setup>
import { ref } from 'vue';
import { useStudioStore } from '../stores/studio';
import QrcodeVue from 'qrcode.vue';

const studio = useStudioStore();
const activeTab = ref('brand');

// Estado para nuevo banner
const newBannerText = ref('');
const newBannerSub = ref('');
const isNewTicker = ref(false);

defineProps({
  isTunnelActive: Boolean,
  guestLink: String
});

defineEmits(['toggle-tunnel']);

// Manejo de archivos
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
    studio.addBanner(newBannerText.value, newBannerSub.value, isNewTicker.value);
    newBannerText.value = '';
    newBannerSub.value = '';
    isNewTicker.value = false;
  }
};
</script>

<template>
  <aside class="w-80 bg-[#020617] border-l border-white/5 shrink-0 flex flex-col h-full overflow-hidden shadow-2xl">
    <!-- Pestañas Simuladas -->
    <div class="flex border-b border-white/5 bg-black/20">
      <button 
        @click="activeTab = 'brand'"
        :class="[activeTab === 'brand' ? 'border-orange-500 text-orange-500 bg-white/5' : 'border-transparent text-slate-500 hover:text-slate-300']"
        class="flex-1 py-5 text-[10px] font-black uppercase tracking-widest border-b-2 transition-all"
      >
        Marca
      </button>
      <button 
        @click="activeTab = 'banners'"
        :class="[activeTab === 'banners' ? 'border-orange-500 text-orange-500 bg-white/5' : 'border-transparent text-slate-500 hover:text-slate-300']"
        class="flex-1 py-5 text-[10px] font-black uppercase tracking-widest border-b-2 transition-all"
      >
        Banners
      </button>
      <button 
        @click="activeTab = 'record'"
        :class="[activeTab === 'record' ? 'border-orange-500 text-orange-500 bg-white/5' : 'border-transparent text-slate-500 hover:text-slate-300']"
        class="flex-1 py-5 text-[10px] font-black uppercase tracking-widest border-b-2 transition-all"
      >
        Grabación
      </button>
    </div>

    <div class="p-8 space-y-10 flex-1 overflow-y-auto scrollbar-hide">
      <!-- Sección Grabación -->
      <div v-if="activeTab === 'record'" class="space-y-8 animate-in fade-in duration-300">
        <div class="p-6 bg-white/5 border border-white/5 rounded-3xl space-y-6">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-red-500/20 flex items-center justify-center">
              <div class="w-3 h-3 bg-red-500 rounded-full" :class="{'animate-pulse': studio.isRecording}"></div>
            </div>
            <div>
              <h3 class="text-xs font-black text-white uppercase tracking-widest">Grabador Local</h3>
              <p class="text-[9px] text-slate-500 uppercase font-bold">{{ studio.isRecording ? 'Grabando ahora...' : 'Listo para grabar' }}</p>
            </div>
          </div>

          <div class="space-y-3">
            <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest">Nombre del Archivo</label>
            <input 
              v-model="studio.recordingTitle" 
              placeholder="Nombre del video..." 
              class="w-full bg-black/40 border border-white/10 rounded-xl px-4 py-3 text-xs text-white outline-none focus:ring-1 ring-orange-500/50"
            />
          </div>

          <button 
            @click="$emit(studio.isRecording ? 'stop-record' : 'start-record')"
            :class="studio.isRecording ? 'bg-white text-black' : 'bg-red-600 text-white'"
            class="w-full py-4 rounded-xl font-black text-[10px] uppercase tracking-[0.2em] transition-all shadow-xl"
          >
            {{ studio.isRecording ? 'Detener y Descargar' : 'Iniciar Grabación MP4' }}
          </button>
        </div>

        <div class="p-6 border border-white/5 rounded-3xl space-y-4">
          <h4 class="text-[9px] font-black text-slate-500 uppercase tracking-widest">Información Técnica</h4>
          <ul class="space-y-2">
            <li class="flex justify-between text-[9px] font-bold"><span class="text-slate-600">Formato:</span> <span class="text-slate-400">MP4 (H.264 / 10-bit)</span></li>
            <li class="flex justify-between text-[9px] font-bold"><span class="text-slate-600">Resolución:</span> <span class="text-slate-400">1920 x 1080 (Full HD)</span></li>
            <li class="flex justify-between text-[9px] font-bold"><span class="text-slate-600">FPS:</span> <span class="text-slate-400">30 fps</span></li>
          </ul>
        </div>
      </div>
      <!-- Sección Marca -->
      <div v-if="activeTab === 'brand'" class="space-y-10 animate-in fade-in duration-300">
        <div class="space-y-4">
          <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest">Color de Marca</label>
          <div class="grid grid-cols-6 gap-3">
            <button 
              v-for="color in ['#f97316', '#3b82f6', '#10b981', '#ef4444', '#8b5cf6', '#ffffff']" 
              :key="color"
              @click="studio.brandColor = color"
              :style="{ backgroundColor: color }"
              class="w-8 h-8 rounded-xl border-2 transition-all hover:scale-110"
              :class="[studio.brandColor === color ? 'border-white shadow-[0_0_15px_rgba(255,255,255,0.3)]' : 'border-transparent']"
            ></button>
          </div>
        </div>

        <div class="space-y-4">
          <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest">Logo del Canal</label>
          <div class="space-y-4">
            <div class="flex gap-2">
              <label class="flex-1 flex items-center justify-center p-3 border border-white/5 bg-white/5 rounded-xl hover:bg-white/10 transition-colors text-[9px] font-bold text-slate-300 uppercase cursor-pointer">
                Subir Logo
                <input type="file" class="hidden" accept="image/*" @change="e => handleFileUpload(e, 'logo')" />
              </label>
              <button @click="studio.setLogo(null)" class="px-4 border border-white/5 bg-white/5 rounded-xl hover:bg-white/10 transition-colors text-[9px] font-bold text-slate-500 uppercase">Quitar</button>
            </div>
            
            <!-- Control de Tamaño -->
            <div v-if="studio.logoUrl" class="space-y-2">
              <div class="flex justify-between text-[8px] font-bold text-slate-500 uppercase">
                <span>Tamaño del Logo</span>
                <span>{{ Math.round(studio.logoScale * 500) }}%</span>
              </div>
              <input 
                type="range" min="0.05" max="0.5" step="0.01" 
                v-model.number="studio.logoScale" 
                class="w-full h-1.5 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-orange-500"
              />
            </div>

            <div v-if="studio.logoUrl" class="p-2 bg-white/5 rounded-xl border border-white/5 flex justify-center overflow-hidden">
              <img :src="studio.logoUrl" :style="{ transform: `scale(${studio.logoScale * 2})` }" class="h-12 object-contain transition-transform" />
            </div>
          </div>
        </div>

        <div class="space-y-4">
          <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest">Fondo de Estudio</label>
          <div class="space-y-3">
            <div class="flex gap-2">
              <label class="flex-1 flex items-center justify-center p-3 border border-white/5 bg-white/5 rounded-xl hover:bg-white/10 transition-colors text-[9px] font-bold text-slate-300 uppercase cursor-pointer">
                Subir Fondo
                <input type="file" class="hidden" accept="image/*" @change="e => handleFileUpload(e, 'bg')" />
              </label>
              <button @click="studio.setBackground(null)" class="px-4 border border-white/5 bg-white/5 rounded-xl hover:bg-white/10 transition-colors text-[9px] font-bold text-slate-500 uppercase">Limpio</button>
            </div>
            <div v-if="studio.backgroundUrl" class="aspect-video bg-cover bg-center rounded-xl border border-white/5 shadow-lg" :style="{ backgroundImage: `url(${studio.backgroundUrl})` }"></div>
          </div>
        </div>

        <div class="space-y-4">
          <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest">Estilo Visual</label>
          <div class="grid grid-cols-2 gap-2">
            <button v-for="s in ['modern', 'classic', 'neon', 'minimal']" :key="s"
              @click="studio.bannerStyle = s"
              class="px-3 py-2 border border-white/5 bg-white/5 rounded-xl text-[9px] font-bold uppercase transition-all"
              :class="[studio.bannerStyle === s ? 'text-orange-500 border-orange-500/50 bg-orange-500/5' : 'text-slate-400']"
            >{{ s }}</button>
          </div>
        </div>

        <div class="space-y-4">
          <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest">Tipografía Pro</label>
          <div class="relative group">
            <select v-model="studio.bannerFont" class="w-full bg-white/5 border border-white/10 rounded-2xl px-5 py-4 text-xs text-white outline-none appearance-none focus:ring-2 ring-orange-500/30 transition-all cursor-pointer hover:bg-white/10">
              <option value="Outfit" class="bg-[#020617] text-white">Outfit (Modern Sans)</option>
              <option value="Bebas Neue" class="bg-[#020617] text-white">Bebas Neue (Impact)</option>
              <option value="Space Grotesk" class="bg-[#020617] text-white">Space Grotesk (Tech)</option>
              <option value="Montserrat" class="bg-[#020617] text-white">Montserrat (Clean)</option>
              <option value="Playfair Display" class="bg-[#020617] text-white">Playfair (Elegant Serif)</option>
              <option value="Mono" class="bg-[#020617] text-white">Roboto Mono (Developer)</option>
            </select>
            <div class="absolute right-5 top-1/2 -translate-y-1/2 pointer-events-none text-slate-500 group-hover:text-orange-500 transition-colors">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M19 9l-7 7-7-7" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
            </div>
          </div>
        </div>

        <div class="space-y-4">
          <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest">Opciones de Pantalla</label>
          <div class="space-y-3">
            <div class="flex items-center justify-between p-3 bg-white/5 rounded-xl border border-white/5">
              <span class="text-[10px] font-bold text-slate-400 uppercase">Mostrar Nombres</span>
              <button @click="studio.showParticipantNames = !studio.showParticipantNames" 
                class="w-10 h-5 rounded-full transition-all relative"
                :class="studio.showParticipantNames ? 'bg-orange-500' : 'bg-slate-800'"
              >
                <div class="absolute top-1 w-3 h-3 bg-white rounded-full transition-all" :class="studio.showParticipantNames ? 'left-6' : 'left-1'"></div>
              </button>
            </div>
            <div class="flex items-center justify-between p-3 bg-white/5 rounded-xl border border-white/5">
              <span class="text-[10px] font-bold text-slate-400 uppercase">Mostrar Badge LIVE</span>
              <button @click="studio.showLiveBadge = !studio.showLiveBadge" 
                class="w-10 h-5 rounded-full transition-all relative"
                :class="studio.showLiveBadge ? 'bg-orange-500' : 'bg-slate-800'"
              >
                <div class="absolute top-1 w-3 h-3 bg-white rounded-full transition-all" :class="studio.showLiveBadge ? 'left-6' : 'left-1'"></div>
              </button>
            </div>
          </div>
        </div>

        <div class="pt-10 border-t border-white/5">
          <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest mb-6 block">Acceso Invitados</label>
          <div v-if="isTunnelActive" class="space-y-6">
            <div class="bg-white p-4 rounded-3xl flex flex-col items-center gap-4 shadow-2xl">
              <QrcodeVue :value="guestLink" :size="140" level="H" class="rounded-xl shadow-sm" />
              <button @click="$emit('toggle-tunnel')" class="text-[9px] font-black text-red-500 hover:text-red-600 uppercase tracking-widest">Cerrar Puerta</button>
            </div>
            <div class="bg-black/40 border border-white/5 p-4 rounded-2xl text-[9px] text-slate-400 font-mono break-all leading-relaxed shadow-inner">
              {{ guestLink }}
            </div>
          </div>
          <button 
            v-else 
            @click="$emit('toggle-tunnel')"
            class="w-full bg-white text-[#020617] text-[10px] font-black py-5 rounded-2xl uppercase tracking-[0.2em] transition-all hover:bg-slate-200 shadow-xl"
          >
            Abrir Invitaciones
          </button>
        </div>
      </div>

      <!-- Sección Banners -->
      <div v-if="activeTab === 'banners'" class="space-y-8 animate-in slide-in-from-right duration-300">
        
        <!-- Controles de Posición (Solo si hay uno activo y no es ticker) -->
        <div v-if="studio.banners.find(b => b.active && !b.isTicker)" class="p-6 bg-orange-500/10 border border-orange-500/20 rounded-3xl space-y-4 shadow-lg">
           <label class="text-[10px] font-black text-orange-500 uppercase tracking-widest">Ajustar Posición</label>
           
           <div class="space-y-2">
              <div class="flex justify-between text-[8px] font-bold text-slate-400 uppercase"><span>Vertical</span><span>{{ studio.bannerY }}px</span></div>
              <input type="range" min="0" max="720" v-model.number="studio.bannerY" class="w-full h-1 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-orange-500" />
           </div>

           <div class="space-y-2">
              <div class="flex justify-between text-[8px] font-bold text-slate-400 uppercase"><span>Horizontal</span><span>{{ studio.bannerX }}px</span></div>
              <input type="range" min="0" max="1280" v-model.number="studio.bannerX" class="w-full h-1 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-orange-500" />
           </div>
        </div>

        <div class="space-y-4">
          <div v-for="banner in studio.banners" :key="banner.id" 
            @click="studio.toggleBanner(banner.id)"
            class="p-5 rounded-2xl border-2 transition-all cursor-pointer relative group overflow-hidden"
            :class="[banner.active ? 'border-orange-500 bg-orange-500/5 shadow-lg shadow-orange-500/10' : 'border-white/5 bg-white/[0.02] hover:bg-white/5']"
          >
            <div class="font-bold text-xs pr-8 text-white flex items-center gap-2">
              <span v-if="banner.active" class="w-2 h-2 bg-orange-500 rounded-full animate-pulse"></span>
              {{ banner.text }}
            </div>
            <div class="text-[10px] text-slate-400 mt-1">{{ banner.subtext }}</div>
            <div v-if="banner.isTicker" class="text-[8px] text-orange-500 font-black uppercase mt-2 tracking-widest">Modo Ticker</div>
            
            <button @click.stop="studio.removeBanner(banner.id)" class="absolute top-4 right-4 text-slate-600 hover:text-red-500 transition-colors">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path d="M6 18L18 6M6 6l12 12" stroke-width="2"/></svg>
            </button>
          </div>
        </div>

        <div class="p-6 bg-black/40 rounded-3xl border border-white/5 space-y-4 shadow-inner">
          <label class="text-[10px] font-black text-slate-500 uppercase tracking-widest">Nuevo Banner</label>
          <input v-model="newBannerText" placeholder="Título principal..." class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-xs text-white outline-none focus:ring-1 ring-orange-500/50" />
          <input v-model="newBannerSub" placeholder="Subtítulo o detalle..." class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-xs text-white outline-none focus:ring-1 ring-orange-500/50" />
          
          <div class="flex items-center gap-3 py-2 px-1">
            <button @click="isNewTicker = !isNewTicker" 
              class="w-10 h-5 rounded-full transition-all relative shrink-0"
              :class="isNewTicker ? 'bg-orange-500' : 'bg-slate-800'"
            >
              <div class="absolute top-1 w-3 h-3 bg-white rounded-full transition-all" :class="isNewTicker ? 'left-6' : 'left-1'"></div>
            </button>
            <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest cursor-pointer" @click="isNewTicker = !isNewTicker">Modo Cinta</label>
          </div>
          
          <button @click="addNewBanner" class="w-full bg-orange-500 text-white text-[10px] font-black py-4 rounded-xl uppercase tracking-[0.2em] shadow-lg shadow-orange-500/20 hover:bg-orange-600 transition-all">Añadir</button>
        </div>
      </div>
    </div>
  </aside>
</template>
