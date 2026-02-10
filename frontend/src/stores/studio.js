import { defineStore } from 'pinia';

export const useStudioStore = defineStore('studio', {
  state: () => ({
    isInitialized: false,
    layout: 'grid', // solo, grid, sidebar, pip
    userName: 'Anfitrión',
    
    // Streams
    localStream: null,
    guestStream: null,
    guestConnected: false,
    
    // Controles de Medios
    isCamOn: true,
    isMicOn: true,
    
    // Inventario de Dispositivos
    videoDevices: [],
    audioDevices: [],
    currentVideoDeviceId: localStorage.getItem('lastVideoDevice') || undefined,
    currentAudioDeviceId: localStorage.getItem('lastAudioDevice') || undefined,
    
    // Estética
    brandColor: '#f97316',
    logoUrl: null,
    logoScale: 0.15,
    backgroundUrl: null,
    bannerStyle: 'modern', // modern, classic, neon, minimal
    bannerFont: 'Outfit',   // Outfit, Bebas Neue, Space Grotesk, Inter, Mono
    bannerY: 590,          // Posición vertical del banner
    bannerX: 240,          // Posición horizontal
    showParticipantNames: true,
    showLiveBadge: true,
    recordingTitle: 'Grabacion-PAK',
    isRecording: false,
    
    // Banners (Lower Thirds)
    banners: JSON.parse(localStorage.getItem('studio_banners')) || [
      { id: 1, text: '¡Bienvenidos a PAK STUDIO!', subtext: 'Transmisión en vivo', active: false, isTicker: false }
    ],
    
    // Gestión del Escenario
    participantsOnStage: [],
  }),

  actions: {
    addBanner(text, subtext, isTicker = false) {
      this.banners.push({ id: Date.now(), text, subtext, active: false, isTicker });
      this.saveBanners();
    },
    toggleBanner(id) {
      this.banners.forEach(b => {
        if (b.id === id) b.active = !b.active;
        else b.active = false;
      });
      this.saveBanners();
    },
    removeBanner(id) {
      this.banners = this.banners.filter(b => b.id !== id);
      this.saveBanners();
    },
    saveBanners() {
      localStorage.setItem('studio_banners', JSON.stringify(this.banners));
    },
    setLogo(url) { this.logoUrl = url; },
    setBackground(url) { this.backgroundUrl = url; },
    toggleParticipant(id) {
      if (this.participantsOnStage.includes(id)) {
        this.participantsOnStage = this.participantsOnStage.filter(p => p !== id);
      } else {
        this.participantsOnStage.push(id);
      }
    },
    async fetchDevices() {
      try {
        const devices = await navigator.mediaDevices.enumerateDevices();
        this.videoDevices = devices.filter(d => d.kind === 'videoinput');
        this.audioDevices = devices.filter(d => d.kind === 'audioinput');
      } catch (e) {
        console.error("Error listando dispositivos:", e);
      }
    },

    async initLocalStream(vDeviceId = this.currentVideoDeviceId, aDeviceId = this.currentAudioDeviceId) {
      if (this.localStream) {
        this.localStream.getTracks().forEach(track => track.stop());
      }

      // Intentar primero con calidad ideal, luego fallback a básico
      const configs = [
        {
          audio: aDeviceId ? { deviceId: { ideal: aDeviceId } } : true,
          video: {
            deviceId: vDeviceId ? { ideal: vDeviceId } : undefined,
            width: { ideal: 1280 },
            height: { ideal: 720 },
            frameRate: { ideal: 30 }
          }
        },
        { video: true, audio: true }
      ];

      for (const constraints of configs) {
        try {
          console.log("Admin: Intentando capturar con:", constraints);
          const stream = await navigator.mediaDevices.getUserMedia(constraints);
          
          const vTrack = stream.getVideoTracks()[0];
          const aTrack = stream.getAudioTracks()[0];
          
          if (vTrack) {
            this.currentVideoDeviceId = vTrack.getSettings().deviceId;
            localStorage.setItem('lastVideoDevice', this.currentVideoDeviceId);
          }
          if (aTrack) {
            this.currentAudioDeviceId = aTrack.getSettings().deviceId;
            localStorage.setItem('lastAudioDevice', this.currentAudioDeviceId);
          }

          this.localStream = stream;
          this.isInitialized = true;
          
          // Auto-subir al escenario si está vacío
          if (this.participantsOnStage.length === 0) {
            this.participantsOnStage.push('local');
          }

          this.applyMediaStates();
          
          if (this.videoDevices.length === 0) await this.fetchDevices();
          console.log("Admin: Cámara iniciada con éxito");
          return stream;
        } catch (e) {
          console.warn("Admin: Fallo en intento de cámara:", e);
        }
      }

      this.isInitialized = false; // Asegurar que no pasamos si falla todo
      alert("No se pudo acceder a la cámara del Admin. Revisa los permisos.");
      throw new Error("Cámara no disponible");
    },

    toggleCam() {
      this.isCamOn = !this.isCamOn;
      this.applyMediaStates();
    },

    toggleMic() {
      this.isMicOn = !this.isMicOn;
      this.applyMediaStates();
    },

    applyMediaStates() {
      if (this.localStream) {
        this.localStream.getVideoTracks().forEach(t => t.enabled = this.isCamOn);
        this.localStream.getAudioTracks().forEach(t => t.enabled = this.isMicOn);
      }
    },

    setLayout(l) {
      this.layout = l;
    }
  }
});
