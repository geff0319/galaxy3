<template>
<!--  <div class="mini-window"-->
<!--       :class="{ 'expanded': isExpanded }"-->
<!--       @mousedown="handleMouseDown"-->
<!--       @mouseenter="handleMouseEnter"-->
<!--       @mouseleave="handleMouseLeave">-->
  <div class="mini-window"
       :class="{ 'expanded': isExpanded }"
       @mouseenter="handleMouseEnter"
       @mouseleave="handleMouseLeave"
       @click="toggleExpand">
    <!-- 折叠状态 -->
    <div class="mini-content">
      <span class="mini-bg"></span>
    </div>

    <!-- 展开状态 -->
    <div class="expanded-content" @click.stop>
      <div class="actions">
        <div class="action-item" @click="handleAction('restore')">
          <icon icon="home"></icon>
          <span>收起</span>
        </div>
        <div class="action-item" @click="handleAction('settings')">
          <icon icon="home"></icon>
          <span>隐藏</span>
        </div>
        <div class="action-item" @click="handleAction('dashboard')">
          <icon icon="home"></icon>
          <span>仪表盘</span>
        </div>
      </div>
    </div>
  </div>
<!--  <div style="background: #29b280">123123</div>-->
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
let mouseEnterTime = null
let mouseLeaveTime:any = null
let windowInitialX = 0
let windowInitialY = 0

// 处理鼠标按下事件
const handleMouseDown = (e: MouseEvent) => {
  // if (isExpanded.value) return // 展开状态不允许拖动
  //
  // isDragging = false
  // initialMouseX = e.screenX // 使用screenX/screenY获取相对于屏幕的坐标
  // initialMouseY = e.screenY
  // console.log(initialMouseX)
  // console.log(initialMouseY)
  // const BallMenuWin = Window.Get("BallMenuWin")
  // BallMenuWin.Position().then((p)=>{
  //   console.log(11111)
  //   console.log('SetPosition===>',p.x,p.y)
  // })

  // mouseDownTime = Date.now()
  // 获取窗口初始位置
  // ipcRenderService.invoke('app:window:get-position').then(([x, y]: [number, number]) => {
  //   windowInitialX = x
  //   windowInitialY = y
  //
  //   document.addEventListener('mousemove', handleMouseMove)
  //   document.addEventListener('mouseup', handleMouseUp)
  // })
  document.addEventListener('mouseup', handleMouseUp)
}

// 处理鼠标移动事件
// const handleMouseMove = (e: MouseEvent) => {
//   const deltaX = e.screenX - initialMouseX
//   const deltaY = e.screenY - initialMouseY
//
//   // 判断是否达到拖动阈值
//   if (!isDragging && (Math.abs(deltaX) > 5 || Math.abs(deltaY) > 5)) {
//     isDragging = true
//   }
//
//   if (isDragging) {
//     // 计算新位置
//     const newX = windowInitialX + deltaX
//     const newY = windowInitialY + deltaY
//
//     // 发送新位置到主进程
//     ipcRenderService.send('app:window:set-position', { x: newX, y: newY })
//   }
// }

const handleMouseUp = (e: MouseEvent) => {
  // document.removeEventListener('mousemove', handleMouseMove)
  document.removeEventListener('mouseup', handleMouseUp)
  initialMouseX = e.clientX // 使用screenX/screenY获取相对于屏幕的坐标
  initialMouseY = e.clientY
  // 如果不是拖拽且点击时间小于200ms，则触发展开/收起
  // if (!isDragging && (Date.now() - mouseDownTime < 200)) {
  //   toggleExpand()
  // }
  toggleExpand()
}

const handleMouseEnter = async () => {
  console.log("handleMouseEnter")
  // if (mouseLeaveTime !== null) {
  //   clearTimeout(mouseLeaveTime);
  // }
}
function sleep(ms:any) {
  return new Promise(resolve => setTimeout(resolve, ms));
}
const handleMouseLeave = async (e:MouseEvent) => {
  console.log("handleMouseLeave")
  if(isExpanded.value){
    return
    // isExpanded.value = false
    // await sleep(300);
  }
  const BallMenuWin = Window.Get("BallMenuWin")
  const SingleBallWin = Window.Get("SingleBallWin")
  const p = await BallMenuWin.Position()
  await SingleBallWin.SetPosition(p.x,p.y)

  await BallMenuWin.Hide()
  await SingleBallWin.Show().then(()=>SingleBallWin.Focus().then(()=>SingleBallWin.SetAlwaysOnTop(true)))

  // mouseLeaveTime = setTimeout(async() => {
  //   const BallMenuWin = Window.Get("BallMenuWin")
  //   const SingleBallWin = Window.Get("SingleBallWin")
  //   const p = await BallMenuWin.Position()
  //   await SingleBallWin.SetPosition(p.x,p.y)
  //   await BallMenuWin.Hide()
  //   await SingleBallWin.Show()
  // }, 100);  // 延迟1秒
}

const toggleExpand = () => {

  isExpanded.value = !isExpanded.value
}

const handleAction = (action: string) => {
  switch (action) {
    case 'restore':
      isExpanded.value = false
      console.log('restore')
      break
    case 'dashboard':
      // ipcRenderService.send('app:window:restore-main', { route: 'INDEX' })
      console.log('dashboard')
      break
    case 'settings':
      Window.Get("BallMenuWin").Hide()
      console.log('settings')
      break
  }
  isExpanded.value = false
}

</script>

<style lang="less" scoped>
.mini-window {
  --wails-draggable: drag;
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
  cursor: pointer;

  &.expanded {
    width: 160px;
    height: 150px;
    border-radius: 12px;
    transform: translate(-110px, 0);

    .mini-content {
      opacity: 0;
      pointer-events: none;
    }

    .expanded-content {
      opacity: 1;
      pointer-events: auto;
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

  .expanded-content {
    position: absolute;
    top: 0;
    right: 0;
    width: 160px;
    height: 150px;
    opacity: 0;
    padding: 9px 12px;
    pointer-events: none;
    transition: opacity 0.3s;

    .actions {
      display: flex;
      flex-direction: column; // 移除 reverse
      gap: 8px;

      .action-item {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 10px 12px;
        border-radius: 8px;
        cursor: pointer;
        transition: all 0.2s ease;
        color: #3498db;

        span {
          font-size: 14px;
        }

        &:hover {
          background-color: #2a75a5;
          transform: scale(1.06);
          outline: 1px solid #005c8e;
        }
      }
    }
  }
}
</style>