import {mergeConfig} from 'vite';
import baseConfig from './vite.config.base';

export default mergeConfig(
    {
        mode: 'development',
        server: {
            proxy: {
                "/api": {
                    target: "http://localhost:8080",
                    changeOrigin: true,
                },
            },
            open: true,
            fs: {
                strict: true,
            },
        },
        plugins: [],
    },
    baseConfig
);
