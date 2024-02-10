import axios from 'axios';

export interface pixelStatus {
    availablePixels: number
}

export function getStatus() {
    return axios.get<pixelStatus>('/api/pixels/status')
}

export function getColors() {
    return axios.get<string[]>('/api/pixels/colors')
}

export interface canvasColor {
    width: number;
    height: number;
    colors: { [key: string]: string }
    canvas: string
}

export function getPixels() {
    return axios.get<canvasColor>('/api/pixels')
}

export interface pixel {
    x: number;
    y: number;
    color: string;
}

export interface setPixelData {
    pixels: pixel[];
}

export function setPixels(data: setPixelData) {
    return axios.post('/api/pixels', data)
}
