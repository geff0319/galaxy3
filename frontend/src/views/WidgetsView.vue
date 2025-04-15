<script setup lang="ts">
import FlipClock from "./FlipClockView/FlipClock.vue";
import YtdlpWidget from "./YtdlpView/YtdlpWidget.vue";
import OptionView from '@/components/Option/index.vue'
import BallView from '@/views/BallView/index.vue'

import { useRoute } from 'vue-router';
const route = useRoute();
</script>

<template>
<!--  <div v-if="route.path ==='/widgets'" style="&#45;&#45;wails-draggable: drag;&#45;&#45;custom-contextmenu: WidgetsWinMenu; width: 380px;height: 70px">-->
<!--    <FlipClock />-->
<!--  </div>-->
<!--  <div v-else-if="route.path ==='/widgets/ball'">-->
<!--    <BallView />-->
<!--  </div>-->
<!--  <div v-else-if="route.path ==='/ytdlpWidgets'">-->
<!--    <YtdlpWidget />-->
<!--  </div>-->
<!--  <div v-else-if="route.path ==='/option'">-->
<!--    <OptionView />-->
<!--  </div>-->
  <div>
    <RouterView v-slot="{ Component }">
      <component :is="Component" />
    </RouterView>
  </div>
</template>

<style lang="less" scoped>

</style>
