import {mergeConfig} from 'vite';
import baseConfig from './vite.config.base';

export default mergeConfig(
    {
        mode: 'umd',
        plugins: [],
        base: './',
        build: {
            lib: {
                entry: 'src/main.ts',
                name: 'NekoPixel',
                fileName: (format) => `neko-pixel-app.${format}.js`
            },
            rollupOptions: {
                output: {
                    globals: {
                        vue: 'Vue'
                    }
                }
            },
        },
    },
    baseConfig
);
