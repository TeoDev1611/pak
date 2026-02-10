// api.js - Adaptador para comunicarse con el backend Go vía HTTP/WebSocket

const API_URL = `http://${window.location.hostname}:8080/api`;
const WS_URL = `ws://${window.location.hostname}:8080/stream-ws`;

let streamWs = null;

export const CheckDependencies = async () => {
    try {
        const res = await fetch(`${API_URL}/deps`);
        return await res.json();
    } catch (e) {
        console.error("Error checking deps:", e);
        return { ffmpeg: false };
    }
};

export const ToggleTunnel = async () => {
    try {
        const res = await fetch(`${API_URL}/tunnel`, { method: 'POST' });
        const data = await res.json();
        return data.url || "timeout";
    } catch (e) {
        console.error("Error toggle tunnel:", e);
        return "";
    }
};

export const StartStream = async (rtmpUrl) => {
    // Iniciar conexión WebSocket para el video
    streamWs = new WebSocket(WS_URL);
    streamWs.binaryType = "arraybuffer";
    
    streamWs.onopen = async () => {
        console.log("WS Stream conectado");
        await fetch(`${API_URL}/stream/start`, {
            method: 'POST',
            body: JSON.stringify({ url: rtmpUrl }),
            headers: { 'Content-Type': 'application/json' }
        });
    };
};

export const StopStream = async () => {
    if (streamWs) {
        streamWs.close();
        streamWs = null;
    }
    await fetch(`${API_URL}/stream/stop`, { method: 'POST' });
};

export const PushVideoChunk = (chunk) => {
    // chunk es un Array (viene de App.vue), convertir a Uint8Array
    if (streamWs && streamWs.readyState === WebSocket.OPEN) {
        streamWs.send(new Uint8Array(chunk));
    }
};
