<script setup lang="ts">
import { useRoute,useRouter } from 'vue-router';
import {Events, Window} from "@/bridge";
import {message} from "ant-design-vue";
import * as Stores from "@/stores";

const route = useRoute();
const router = useRouter();
const wsClientStore = Stores.useWsClientStore()

Events.On('notify',(event:any)=>{
  router.push({path: '/ytdlp'})
  const w = Window.Get('MainWin')
  w.IsFocused().then((status)=>{
    console.log(status)
    if (!status){
      w.Show()
      w.Focus()
    }
  })
  console.log(event)
  switch (event.data[0]) {
    case "info":
      message.info(event.data[1]);
      break;
    case "error":
      message.error(event.data[1])
      break;
    case "warn":
      message.warn(event.data[1]);
      break;
    default:
      message.info(event.data[1])
      break;
  }
})
//ws
if(wsClientStore.ws.autoConnect){
  try {
    wsClientStore.connectWs()
  }catch (error){
    message.error("ws连接失败：" + error)
  }
}

</script>

<template>
  <div class="content">
    <RouterView v-slot="{ Component }">
<!--      <KeepAlive>-->
<!--        <component v-if="route.meta.keepAlive" :is="Component" />-->
<!--      </KeepAlive>-->
      <component v-if="!route.meta.keepAlive" :is="Component" />
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
