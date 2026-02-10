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
    await fetch(`${API_URL}/stream/start`, {
        method: 'POST',
        body: JSON.stringify({ url: rtmpUrl }),
        headers: { 'Content-Type': 'application/json' }
    });
};

export const StopStream = async () => {
    await fetch(`${API_URL}/stream/stop`, { method: 'POST' });
};

export const StartRecording = async (filename) => {
    await fetch(`${API_URL}/stream/record/start`, {
        method: 'POST',
        body: JSON.stringify({ filename }),
        headers: { 'Content-Type': 'application/json' }
    });
};

export const StopRecording = async () => {
    await fetch(`${API_URL}/stream/record/stop`, { method: 'POST' });
};

export const SendOffer = async (sdp) => {
    try {
        const res = await fetch(`${API_URL}/stream/offer`, {
            method: 'POST',
            body: JSON.stringify({ sdp }),
            headers: { 'Content-Type': 'application/json' }
        });
        return await res.json();
    } catch (e) {
        console.error("Error sending offer:", e);
        return null;
    }
};
