export {};

declare global {
    interface Window {
        NEKO_CONFIG: {
            pixelBaseURL: string;
        };
    }
}
