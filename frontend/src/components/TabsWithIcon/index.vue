<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import icons from '@/components/Icon/icons'

type TabItemType = {
  // icon:(typeof icons)[number]
  icon:string
  key: string
  tab: string
}

interface Props {
  activeKey: string
  items: TabItemType[]
  height?: string
}

const props = withDefaults(defineProps<Props>(), {
  height: ''
})

const emits = defineEmits(['update:activeKey'])

const { t } = useI18n()

const handleChange = (key: string) => emits('update:activeKey', key)

const isActive = ({ key }: TabItemType) => key === props.activeKey
</script>

<template>
  <div :style="{ height }" class="tabs">
    <div class="tab">
      <button
        v-for="tab in items"
        :key="tab.key"
        @click="handleChange(tab.key)"
        class="custom-btn"
        :class="[{'active-tab':isActive(tab)},{ custom_icon: !isActive(tab) }]"
      >
        <svg v-if="isActive(tab)" > <use :href="`#${tab.icon}-active`||`forbid`"></use></svg>
        <svg v-else > <use :href="`#${tab.icon}`||`forbid`"></use></svg>
<!--        <icon :icon="tab.icon "></icon>-->
        <span :style="{marginLeft: '8px'}">{{ t(tab.tab) }}</span>
      </button>

      <slot name="extra" />
    </div>

    <div class="slot">
      <slot :name="activeKey"></slot>
    </div>
  </div>
</template>

<style lang="less" scoped>
@import "@/assets/main";

.custom_icon{
  .icon_hover();
}
.tabs {
  display: flex;
}
.tab {
  width: 20%;
  display: flex;
  align-items: center;
  flex-direction: column;
}
.slot {
  width: 80%;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}
.custom-btn {
  width: 80%;
  padding: 8px 12px;
  margin: 3px 4px;
  font-size: 14px;
  border: 2px #ccc; /* 设置边框 */
  cursor: pointer; /* 鼠标悬停时显示指针 */
  transition: all 0.3s ease; /* 为所有变化添加过渡效果 */
  color: var(--btn-text-color);
  background-color: var(--btn-text-bg);
  border: none;
  display: flex;
  justify-content: flex-start; /* 左对齐 */
  align-items: center;
  border-radius: 6px;
  &:hover {
    color: var(--btn-text-hover-color);
    background-color: var(--btn-text-hover-bg);
  }
  &:active {
    color: var(--btn-text-active-color);
    background-color: var(--btn-text-active-bg);
  }
  svg {
    width: 22px;
    height: 22px;
  }
}
button.active-tab {
  font-weight: bold; /* 如果选中则加粗 */
  color: var(--btn-text-active-color);
  background-color: var(--btn-text-active-bg);
}
.custom-btn .span{
  font-family: "幼圆", "Yu Yuan", sans-serif;
}
</style>
