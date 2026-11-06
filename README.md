# 🎨 NekoPixel

## Embedding in NekoBox

The Vue package exports `NekoPixelApp` for mounting inside an existing Vue app.
Configure its API base URL as `/service/pixel` on the NekoBox API origin, and
point the gateway's Pixel upstream at the NekoPixel server's `/api` path.

```ts
import { watchEffect } from 'vue'
import NekoPixel, { setAuthToken } from '@e99p1ant/neko-pixel'
import '@e99p1ant/neko-pixel/style.css'

app.use(NekoPixel, { baseURL: '/service/pixel' })
watchEffect(() => setAuthToken(authStore.isSignedIn ? authStore.sessionID : ''))
```

Render `<NekoPixelApp :editable="isSignedIn" />` where the canvas should appear.
The component sends the current NekoBox session in the `Authorization` header.

Keep the NekoPixel backend reachable only from the NekoBox gateway. It trusts
the `X-NekoBox-User-ID` header that the gateway sets after resolving the session.

## License

MIT License
