<template>
  <div class="mini-window"
       :class="{ 'expanded': isExpanded }"
       @mouseenter="handleMouseEnter"
       @mouseleave="handleMouseLeave">
    <!-- 折叠状态 -->
    <div class="mini-content">
      <span class="mini-bg"></span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {Window} from "@/bridge";
// import { ipcRenderService } from '@/render/services/ipcService'
document.body.setAttribute('theme-mode', 'widget')
const isExpanded = ref(false)
let isDragging = false
let initialMouseX = 0
let initialMouseY = 0
let mouseDownTime = 0
let mouseEnterTime:any = null
let mouseLeaveTime = null

const handleMouseEnter = async() => {
  console.log("ball handleMouseEnter")
  // mouseEnterTime = setTimeout(async() => {
  //   const BallMenuWin = Window.Get("BallMenuWin")
  //   const SingleBallWin = Window.Get("SingleBallWin")
  //   const p = await SingleBallWin.Position()
  //   await BallMenuWin.SetPosition(p.x,p.y)
  //   await SingleBallWin.Hide()
  //   await BallMenuWin.Show()
  // }, 100);  // 延迟1秒
  const BallMenuWin = Window.Get("BallMenuWin")
  const SingleBallWin = Window.Get("SingleBallWin")
  const p = await SingleBallWin.Position()
  await BallMenuWin.SetPosition(p.x,p.y)
  await SingleBallWin.Hide()
  await BallMenuWin.Show().then(()=>BallMenuWin.Focus().then(()=>BallMenuWin.SetAlwaysOnTop(true)))
  // await BallMenuWin.OpenDevTools()
}

const handleMouseLeave = () => {
  // if (mouseEnterTime !== null) {
  //   clearTimeout(mouseEnterTime);
  // }
}

</script>

<style lang="less" scoped>
.mini-window {
  //--wails-draggable: drag;
  position: relative;
  margin-left: 195px;
  margin-top: 5px;
  width: 50px;
  height: 50px;
  border-radius: 25px;
  background: #ffffff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
  overflow: hidden;
  user-select: none;
  //cursor: pointer;

  &.expanded {
    width: 160px;
    height: 150px;
    border-radius: 12px;
    transform: translate(-110px, -100px);

    .mini-content {
      opacity: 0;
      pointer-events: none;
    }
  }

  .mini-content {
    position: absolute;
    bottom: 1px;
    right: 5px;
    opacity: 1;
    transition: opacity 0.3s;
    .mini-bg {
      cursor: pointer;
      display: inline-block;
      background: #8b2876;
      width: 40px;
      height: 40px;
      border-radius: 20px;
    }
  }
}
</style>