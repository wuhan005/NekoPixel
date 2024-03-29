import {mergeConfig} from 'vite';
import baseConfig from './vite.config.base';

export default mergeConfig(
    {
        mode: 'production',
        plugins: [],
        base: './',
        build: {
            rollupOptions: {
                output: {
                    manualChunks: {
                        vue: ['vue'],
                    },
                },
            },
            chunkSizeWarningLimit: 2000,
        },
    },
    baseConfig
);
