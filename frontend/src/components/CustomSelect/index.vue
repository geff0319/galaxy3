<template>
  <div class="custom-select" >
    <div :style="{position: 'relative', display: 'inline-block','pointer-events': disable?'none':'auto'}" @click.stop="toggleOptions" >
      <input
          type="text"
          class="select-input"
          :value="selectedOption"
          placeholder="请选择"
          readonly
      />
      <div>
        <Icon class="div-icon" icon="foldDown" />
<!--        <Icon class="div-icon" v-show="showOptions" icon="foldUp" />-->
      </div>
    </div>

    <div class="select-list" v-show="showOptions" :style="'max-height:'+ maxHeight + 'px'">
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

interface Props {
  options: { label: string; value: string }[]
  maxHeight? :number
  disable?:boolean
}

const props = withDefaults(defineProps<Props>(), {
  options: () => [],
  maxHeight:100,
  disable:false,
})

const emits = defineEmits(['change'])

const showOptions = ref(false);
const selectedOption = ref('');


const toggleOptions = () => {
  showOptions.value = !showOptions.value;
};

const selectOption = (option: { label: string; value: string }) => {
  selectedOption.value = option.label;
  emits('change', option.value)
  showOptions.value = false;
};

const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement;
  const customSelect = document.querySelector('.custom-select');
  if (customSelect && !customSelect.contains(target)) {
    showOptions.value = false;
  }
};

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside);
});

// const init =()=>{
//   console.log("select init")
//   console.log(props.options)
//   selectedOption.value = props.options[0].label;
//   emits('change', props.options[0].value)
// }
// init()
watch(
    ()=>props.options,
    (newValue) => {
        selectedOption.value = newValue[0].label;
        emits('change', newValue[0].value)
        console.log(newValue[0].label) // 更新 currentValue
    }
);
</script>

<style lang="less" scoped>
.custom-select {
  position: relative;
  width: 100%; /* 控件宽度 */
}
.div-icon{
  height: 95%;
  position: absolute;
  right: 1px;
  top: 50%;
  transform: translateY(-50%);
  border-radius: 0 4px 4px 0; /* 圆角 */
  cursor: pointer;
  background-color: #dcdcdc;
  transition: background-color 0.3s; /* 背景颜色过渡效果 */
  &:hover {
    background-color: #c8c8c8; // 悬浮时改变背景颜色
  }
}

.select-input {
  //cursor: not-allowed; /* 禁用点击时的光标样式 */
  background-color: #dcdcdc;
  width: 100%; /* 控件宽度 */
  padding: 7px;
  border: 1px solid  #ccc;
  border-radius: 4px;
  box-sizing: border-box;
  cursor: pointer;
  outline: none; /* 去除默认聚焦边框 */
}

.select-list {
  display: block;
  position: absolute;
  background-color: #dcdcdc;
  border: 1px solid #ccc;
  z-index: 99;
  //max-height: 90px; /* 设置最大高度 */
  width: 100%;
  overflow-y: auto; /* 允许垂直滚动 */

  border-radius: 4px;
  margin-top: 2%;  /* 去掉外边距 */

  //transition: max-height 0.5s ease, opacity 0.5s ease; /* 添加过渡效果 */
}
.select-list ul {
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
