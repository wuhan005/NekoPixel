import { httpClient } from './client'

export interface pixelStatus {
    availablePixels: number
}

export function getStatus() {
    return httpClient.get<pixelStatus>('/pixels/status')
}

export function getColors() {
    return httpClient.get<string[]>('/pixels/colors')
}

export interface canvasColor {
    width: number;
    height: number;
    colors: { [key: string]: string }
    canvas: string
}

export function getPixels() {
    return httpClient.get<canvasColor>('/pixels')
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
    return httpClient.post('/pixels', data)
}
