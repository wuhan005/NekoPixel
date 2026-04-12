import { mergeConfig } from 'vite';
import dts from 'vite-plugin-dts';
import baseConfig from './vite.config.base';

export default mergeConfig(
    {
        mode: 'umd',
        plugins: [
            dts({
                include: ['src/**/*.ts', 'src/**/*.vue'],
                outDir: 'dist',
                rollupTypes: true,
                tsconfigPath: './tsconfig.json',
                insertTypesEntry: true,
            }),
        ],
        base: './',
        build: {
            lib: {
                entry: 'src/lib.ts',
                name: 'NekoPixel',
                fileName: (format) => `neko-pixel-app.${format}.js`,
                formats: ['umd', 'es'],
            },
            rollupOptions: {
                external: ['vue'],
                output: {
                    globals: {
                        vue: 'Vue'
                    },
                    exports: 'named',
                }
            },
        },
    },
    baseConfig
);
