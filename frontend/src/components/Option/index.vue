<template>
  <div class="custom-select" >
    <div class="select-list">
      <ul >
        <li
            v-for="o in options"
            :key="o.value"
            :value="o.value"
            @click="selectOption(o)"
        >
          {{ o.label }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted,onBeforeUnmount,defineComponent, ref, watch } from 'vue';
import {useRoute} from "vue-router";
import {Events,Window} from '@/bridge'
const options = ref<{ label: string; value: string }[]>([])
const route = useRoute();

onMounted(async ()=>{
  const name = await Window.Name()
  console.log('onMounted'+ name)
  Events.On(name,(event:any)=>{
    console.log(event)
    options.value = event.data[0]
  })
})
onBeforeUnmount(async ()=>{
  const name = await Window.Name()
  Events.Off(name)
})
const selectOption =async (option: { label: string; value: string }) => {
  await Window.Hide()
  console.log(option)
  const name = await Window.Name()
  Events.Emit({name:'get'+name, data:[option.label,option.value]})
};

</script>

<style lang="less" scoped>

.custom-select {
  height: 100vh;
  width: 100%; /* 控件宽度 */
  //border: 1px solid #ccc;
}

.select-list {
  display: block;
  position: absolute;
  background-color: #dcdcdc;
  border: 1px solid #ccc;
  z-index: 99;
  //max-height: 140px; /* 设置最大高度 */
  width: 100%;
  overflow-y: auto; /* 允许垂直滚动 */
  border-radius: 4px;
  //margin-top: 2%;  /* 去掉外边距 */
}
.select-list ul {
  font-size: 13px;
  list-style-type: none; /* 去掉左边的点 */
  padding: 0; /* 去掉内边距 */
  margin: 0;
}
.select-list li {
  background-color: #dcdcdc;
  margin: 3px 6px;
  padding-left: 8px;
  padding-bottom: 2px;
  padding-top: 2px;
  cursor: pointer;
  transition: background-color 0.3s, box-shadow 0.3s; /* 平滑过渡 */
}

.select-list li:hover {
  background-color: #c8c8c8; /* 悬停效果 */
  border-radius: 4px;
}
</style>
