<script setup lang="ts">
import { useRoute,useRouter } from 'vue-router';
import {Events, Window} from "@/bridge";
import {message as antmessage} from "ant-design-vue/es/components";
import * as Stores from "@/stores";


</script>

<template>
  <div class="content">
    <RouterView v-slot="{ Component }">
<!--      <KeepAlive >-->
<!--        <component v-if="route.meta.keepAlive" :is="Component" />-->
<!--      </KeepAlive>-->
      <component :is="Component" />
    </RouterView>
  </div>
</template>

<style scoped>
.content {
  overflow-y: auto;
  margin-top: 8px;
  padding: 0 8px;
  height: 100%;
  display: flex;
  flex-direction: column;
}
</style>
