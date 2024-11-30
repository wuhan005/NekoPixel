<template>
  <div class="uk-container uk-container-large">
    <div class="uk-card uk-card-default">
      <div class="uk-card-header" v-if="!NEKO_BOX_MODE">
        <h3 class="uk-card-title">NekoPixel</h3>
      </div>
      <div class="uk-card-body" v-if="!isMobile()">
        <div class="uk-margin">
          <div class="place-board" ref="board">
            <canvas ref="paintingCanvas" class="painting-canvas"
                    :width="size.width"
                    :height="size.height"
                    @wheel="handleWheel"
                    @mousedown="handleMouseDown"
                    @mouseup="handleMouseUp"
                    @mousemove="handleMouseMove"
            >
            </canvas>
          </div>
        </div>
        <div class="toolbar">
          <div class="color-picker">
            <div class="item palette-panel">
              <div class="palette-item" v-for="(color, index) in colors" v-bind:key="index"
                   :style="`background-color: #${color}`" @click="onSelectColor(color)"></div>
            </div>
            <div class="current-palette" :style="`background-color: #${currentColor}`"></div>
            <div class="paint">
              <div v-if="!isPainting">
                <button class="uk-button uk-button-primary uk-button-small" @click="startPainting">开始绘制</button>
              </div>
              <div v-else>
                <div class="uk-text-small">可绘制像素：{{ status.availablePixels }}</div>
                <div class="uk-text-small">已绘制像素：{{ paintedPixels.length }} 个</div>
                <div class="button">
                  <button class="uk-button uk-button-primary uk-button-small" @click="submit">提交</button>
                  <button class="uk-button uk-button-default uk-button-small" @click="stopPainting">清除</button>
                </div>
              </div>
            </div>
          </div>

          <div class="item range">
            <div class="slider">
              <input class="uk-range" type="range" v-model="ratio" min="1" :max="MAX_RATIO" :step="RATIO_STEP"
                     aria-label="Range">
            </div>
            <div class="text uk-text-small">{{ ratio }}x</div>
          </div>
        </div>
      </div>
      <div v-else class="uk-card-body">
        请使用电脑浏览器访问本页面，手机端暂不支持绘制。
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref, watch, nextTick} from 'vue'
import {canvasColor, getColors, getPixels, pixel, setPixels, getStatus, pixelStatus} from './api/pixel.ts'

const NEKO_BOX_MODE = import.meta.env.VITE_NEKO_BOX_MODE === 'true'

const WIDTH = 1280
const HEIGHT = 720
const MAX_RATIO = 10
const RATIO_STEP = 0.5

const paintingCanvas = ref<HTMLCanvasElement | null>(null)
const baseCanvas = ref<HTMLCanvasElement | null>(null)
const baseContext = ref<CanvasRenderingContext2D | null>(null)
const deltaX = ref<number>(0)
const deltaY = ref<number>(0)
const currentColor = ref<string>('000000')
const paintedPixels = ref<pixel[]>([])

const size = ref<{ width: number, height: number }>({
  width: WIDTH,
  height: HEIGHT
})
const ratio = ref<number>(1)
const ratioChanging = ref<boolean>(false);
const colors = ref<string[]>([])
const canvasPixels = ref<canvasColor>()
const isMoving = ref<boolean>(true)
const isPainting = ref<boolean>(false)
const status = ref<pixelStatus>({} as pixelStatus)

const isMobile = () => {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
}

const loadCanvasPixels = async () => {
  canvasPixels.value = (await getPixels()).data
}

const initBaseCanvas = () => {
  const width = canvasPixels.value?.width ?? 0
  const height = canvasPixels.value?.height ?? 0

  baseCanvas.value = document.createElement('canvas')
  baseCanvas.value.width = width
  baseCanvas.value.height = height

  baseContext.value = baseCanvas.value.getContext('2d');
  if (!baseContext.value) {
    return
  }

  const pixels = canvasPixels.value
  if (!pixels || !baseContext.value) {
    return
  }

  let colorMap = new Map<string, number[]>()
  Object.keys(pixels.colors).forEach((index: string) => {
    let hexString = pixels.colors[index]
    if (hexString.length === 3) {
      hexString = hexString.replace(/(.)/g, '$1$1')
    }
    const rgbArray = [
      parseInt(hexString.charAt(0) + "" + hexString.charAt(1), 16),
      parseInt(hexString.charAt(2) + "" + hexString.charAt(3), 16),
      parseInt(hexString.charAt(4) + "" + hexString.charAt(5), 16),
    ]

    colorMap.set(index, rgbArray)
  })

  const imageData = baseContext.value.createImageData(width, height)

  const arrayBuffer = new ArrayBuffer(imageData.data.length)
  const clampedArray = new Uint8ClampedArray(arrayBuffer)
  const uint32Array = new Uint32Array(arrayBuffer)
  for (let i = 0; i < pixels.canvas.length; i++) {
    const index = pixels.canvas[i]
    const color = colorMap.get(index) ?? [0, 0, 0]
    const pixelValue = (255 << 24) | (color[2] << 16) | (color[1] << 8) | color[0]; // 注意: 这里使用的是big-endian
    uint32Array[i] = pixelValue;
  }

  imageData.data.set(clampedArray)
  baseContext.value.putImageData(imageData, 0, 0)
}

const refreshCanvas = () => {
  if (!paintingCanvas.value || !baseCanvas.value || !baseContext.value) {
    return
  }

  const ctx = paintingCanvas.value.getContext('2d')
  if (!ctx) {
    return
  }

  ctx.imageSmoothingEnabled = false

  ctx.save()

  ctx.clearRect(0, 0, paintingCanvas.value.width, paintingCanvas.value.height)
  ctx.scale(ratio.value, ratio.value)
  ctx.translate(deltaX.value, deltaY.value)
  ctx.drawImage(baseCanvas.value, 0, 0)

  ctx.restore()
}

watch(ratio, () => {
  if (ratioChanging.value) {
    return
  }
  ratioChanging.value = true
  nextTick(() => {
    refreshCanvas()
    ratioChanging.value = false
  })
})

const startPainting = () => {
  paintedPixels.value = []
  isPainting.value = true
  isMoving.value = false
}

const stopPainting = () => {
  isPainting.value = false
  isMoving.value = true
  initBaseCanvas()
  refreshCanvas()
}

const submit = () => {
  if (paintedPixels.value.length === 0) {
    return
  }

  setPixels({pixels: paintedPixels.value}).then(() => {
    isPainting.value = false
    isMoving.value = true

    loadCanvasPixels().then(() => {
      initBaseCanvas()
      refreshCanvas()
    })
  })
}

const handleWheel = (event: WheelEvent) => {
  event.preventDefault()
  if (ratioChanging.value) {
    return
  }

  ratioChanging.value = true

  if (event.deltaY > 0) {
    if (ratio.value - RATIO_STEP >= 1) {
      ratio.value -= RATIO_STEP
    }
  } else if (ratio.value + RATIO_STEP <= MAX_RATIO) {
    ratio.value += RATIO_STEP
  }

  refreshCanvas()
  ratioChanging.value = false
}

const handleMouseDown = (event: MouseEvent) => {
  if (isMoving.value) {
    return
  }

  event.preventDefault()

  if (isPainting.value) {
    // #1 If the web page is zoomed in, we need to adjust the offset.
    const paintingCanvasRatio = paintingCanvas.value!.getBoundingClientRect().width / paintingCanvas.value!.width

    const x = Math.floor(event.offsetX / ratio.value / paintingCanvasRatio - deltaX.value)
    const y = Math.floor(event.offsetY / ratio.value / paintingCanvasRatio - deltaY.value)
    const pixel = {x, y, color: currentColor.value}

    if (paintedPixels.value.some(p => p.x === pixel.x && p.y === pixel.y)) {
      paintedPixels.value = paintedPixels.value.map(p => {
        if (p.x === pixel.x && p.y === pixel.y) {
          return pixel
        }
        return p
      })

    } else {
      if (paintedPixels.value.length >= status.value.availablePixels) {
        return
      }

      paintedPixels.value.push(pixel)
    }

    // Draw pixel on canvas.
    if (baseContext.value) {
      baseContext.value.fillStyle = `#${currentColor.value}`
      baseContext.value.fillRect(x, y, 1, 1)
      refreshCanvas()
    }
  }
}

const handleMouseUp = (event: MouseEvent) => {
  if (isMoving.value) {
    return
  }

  event.preventDefault()
}

const handleMouseMove = (event: MouseEvent) => {
  if (!isMoving.value) {
    return
  }

  // Move canvas with translate.
  if (event.buttons === 1) {
    deltaX.value += event.movementX / ratio.value
    deltaY.value += event.movementY / ratio.value

    if (deltaX.value > 0) {
      deltaX.value = 0
    }
    if (deltaY.value > 0) {
      deltaY.value = 0
    }

    if (baseContext.value) {
      refreshCanvas()
    }
  }
}

const onSelectColor = (selectedColor: string) => {
  currentColor.value = selectedColor
}

const getStatusInfo = async () => {
  status.value = (await getStatus()).data
}

onMounted(async () => {
  getColors().then(res => {
    colors.value = res.data
  })
  await getStatusInfo()

  loadCanvasPixels().then(() => {
    initBaseCanvas()
    refreshCanvas()
    // initEvents()
  }).catch((err) => {
    console.error(err)
  })
})
</script>

<style lang="scss">
@import "../node_modules/uikit/src/scss/mixins.scss";
@import "../node_modules/uikit/src/scss/mixins-theme.scss";
@import "../node_modules/uikit/src/scss/variables-theme.scss";
@import "../node_modules/uikit/src/scss/variables.scss";
@import "../node_modules/uikit/src/scss/uikit-theme.scss";
@import "../node_modules/uikit/src/scss/uikit.scss";

.painting-canvas {
  display: block;
  background: #fff;
  cursor: default;
  outline: none;
  -webkit-tap-highlight-color: rgba(255, 255, 255, 0);
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;

  .item {
    width: 100%;
  }

  .color-picker {
    flex: 4;
    display: flex;
    align-items: center;
    flex-direction: row;
    gap: 10px;
  }

  .palette-panel {
    border-radius: 10px;
    padding: 5px;
    width: 680px;

    .palette-item {
      display: inline-block;
      width: 28px;
      height: 28px;
      margin: 1px 4px;
      border-radius: 5px;
      cursor: pointer;
      transform-origin: 50% 50%;
      border: 3px solid #00000020;
    }
  }

  .current-palette {
    width: 65px;
    height: 65px;
    border-radius: 5px;
    border: 3px solid #00000020;
  }

  .paint {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    margin-left: 30px;
    gap: 10px;

    .button {
      display: flex;
      gap: 10px;
    }
  }

  .range {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .slider {
      flex-grow: 1;
    }

    .text {
      text-align: right;
      width: 120px;
    }
  }
}
</style>
