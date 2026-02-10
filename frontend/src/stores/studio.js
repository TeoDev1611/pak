import { defineStore } from 'pinia';

export const useStudioStore = defineStore('studio', {
  state: () => ({
    isInitialized: false,
    layout: 'solo', 
    userName: 'Anfitrión',
    streamTitle: 'Mi Transmisión',
    isCamOn: true,
    isMicOn: true,
    localStream: null,
    screenStream: null,
    guestConnected: false,
    brandColor: '#f97316',
    activeLogo: 'https://vitejs.dev/logo.svg',
  }),
  actions: {
    setLayout(l) { this.layout = l; },
    setLogo(logo) { this.activeLogo = logo; },
    
    async initLocalStream() {
      if (this.localStream) return this.localStream;
      console.log("Solicitando permisos...");
      try {
        const stream = await navigator.mediaDevices.getUserMedia({ 
          video: { width: { ideal: 1280 }, height: { ideal: 720 }, frameRate: { ideal: 30, max: 30 } }, 
          audio: true 
        });
        this.localStream = stream;
        this.isInitialized = true;
        return stream;
      } catch (e) {
        console.error("Error de cámara:", e);
        alert("No se pudo acceder a la cámara. Por favor, asegúrate de dar permisos en el navegador.");
      }
    },

    toggleCam() {
      this.isCamOn = !this.isCamOn;
      if (this.localStream) {
        this.localStream.getVideoTracks().forEach(t => t.enabled = this.isCamOn);
      }
    },
    
    toggleMic() {
      this.isMicOn = !this.isMicOn;
      if (this.localStream) {
        this.localStream.getAudioTracks().forEach(t => t.enabled = this.isMicOn);
      }
    }
  }
});