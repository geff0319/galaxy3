<script setup lang="ts">
import FlipClock from "./FlipClockView/FlipClock.vue";
import YtdlpWidget from "./YtdlpView/YtdlpWidget.vue";

import { useRoute } from 'vue-router';
document.body.setAttribute('theme-mode', 'widget')
const route = useRoute();
console.log(route)

</script>

<template>
    <div v-if="route.path ==='/widgets'" style="--wails-draggable: drag;--custom-contextmenu: WidgetsWinMenu; width: 380px;height: 70px">
      <FlipClock />
    </div>
    <div v-if="route.path ==='/ytdlpWidgets'" style="--wails-draggable: drag; width: 600px;height: 230px">
      <YtdlpWidget />
    </div>

</template>

<style lang="less" scoped>

</style>
