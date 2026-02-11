import { defineStore } from 'pinia';

export const useStudioStore = defineStore('studio', {
  state: () => ({
    isInitialized: false,
    layout: 'grid', 
    userName: localStorage.getItem('user_name') || 'Anfitrión',
    streamTitle: 'Mi Transmisión',
    
    // Streams
    localStream: null,
    guestStream: null,
    guestConnected: false,
    
    // Controles de Medios
    isCamOn: true,
    isMicOn: true,
    
    // Mezcla de Audio
    audioContext: null,
    audioDestination: null,
    localAudioSource: null,
    guestAudioSource: null,
    
    // Inventario de Dispositivos
    videoDevices: [],
    audioDevices: [],
    currentVideoDeviceId: localStorage.getItem('lastVideoDevice') || undefined,
    currentAudioDeviceId: localStorage.getItem('lastAudioDevice') || undefined,
    
    // Estética Industrial
    brandColor: '#f97316',
    logoUrl: null,
    logoScale: 0.15,
    logoPosition: 'top-right', // 'top-left', 'top-right', 'bottom-left', 'bottom-right'
    backgroundUrl: null,
    bannerFont: 'Inter', // 'Inter', 'Roboto', 'Oswald', 'Courier Prime'
    bannerY: 590,          
    bannerX: 240,          
    bannerPositionX: 50,
    bannerPositionY: 85,
    bannerBgColor: '#f97316',
    bannerTextColor: '#ffffff',
    bannerBorderRadius: 4,
    bannerPadding: 20,
    showParticipantNames: true,
    showLiveBadge: true,
    recordingTitle: 'Grabacion-PAK',
    isRecording: false,
    rtmpUrl: 'rtmp://localhost/live/stream',
    
    // Banners
    banners: JSON.parse(localStorage.getItem('studio_banners')) || [
      { id: 1, text: '¡BIENVENIDOS A CHASQUI!', subtext: 'PRO STUDIO', active: false, isTicker: false, style: 'modern', isMarquee: false }
    ],
    
    // Gestión del Escenario
    participantsOnStage: [],
  }),

  actions: {
    addBanner(text, subtext, isTicker = false, style = 'custom', isMarquee = false) {
      this.banners.push({ 
        id: Date.now(), 
        text, 
        subtext, 
        active: false, 
        isTicker, 
        style, 
        isMarquee,
        bgColor: this.bannerBgColor,
        textColor: this.bannerTextColor,
        borderRadius: this.bannerBorderRadius
      });
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
      } catch (e) { console.error(e); }
    },

    // ROBUSTEZ DE CÁMARA (Waterfall Fallback)
    async initLocalStream() {
      if (this.localStream) {
        this.localStream.getTracks().forEach(track => track.stop());
      }

      const configs = [
        { video: { width: 1920, height: 1080 }, audio: true },
        { video: { width: 1280, height: 720 }, audio: true },
        { video: true, audio: true }
      ];

      for (const constraints of configs) {
        try {
          const stream = await navigator.mediaDevices.getUserMedia(constraints);
          this.localStream = stream;
          this.isInitialized = true;
          
          if (this.participantsOnStage.length === 0) {
            this.participantsOnStage.push('local');
          }

          this.setupAudioMixing(stream);
          this.applyMediaStates();
          return stream;
        } catch (e) {
          console.warn("Intento de cámara fallido:", e.name);
        }
      }
      alert("No se pudo acceder a la cámara.");
      throw new Error("Cámara no disponible");
    },

    // IMPLEMENTAR MEZCLA DE AUDIO
    setupAudioMixing(localStream) {
      if (!this.audioContext) {
        this.audioContext = new (window.AudioContext || window.webkitAudioContext)();
        this.audioDestination = this.audioContext.createMediaStreamDestination();
      }

      if (this.localAudioSource) this.localAudioSource.disconnect();
      if (localStream.getAudioTracks().length > 0) {
        this.localAudioSource = this.audioContext.createMediaStreamSource(localStream);
        this.localAudioSource.connect(this.audioDestination);
      }
    },

    mixGuestAudio(guestStream) {
      if (this.audioContext && guestStream.getAudioTracks().length > 0) {
        if (this.guestAudioSource) this.guestAudioSource.disconnect();
        this.guestAudioSource = this.audioContext.createMediaStreamSource(guestStream);
        this.guestAudioSource.connect(this.audioDestination);
      }
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

    setLayout(l) { this.layout = l; }
  }
});