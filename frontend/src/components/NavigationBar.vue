<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import rawRoutes from '@/router/routes'

const { t } = useI18n()

const routes = rawRoutes.filter((r) => !r.meta?.hidden)
</script>

<template>
  <div class="nav">
    <div v-for="r in routes" :key="r.path">
      <RouterLink :to="r.path" custom #default="{ navigate, isActive }">
        <div @click="navigate" :title="(r.meta && t(r.meta.name)) || r.name as string" :class="[{'bar-active':isActive},'bar',{'bar-icon_hover':!isActive}]">
          <svg v-if="isActive"> <use :href="`#${r?.meta?.icon}-active`||`forbid`"></use></svg>
          <svg v-else > <use :href="`#${r.meta?.icon}`||`forbid`"></use></svg>
        </div>
<!--        <Button @click="navigate" :type="isActive ? 'link' : 'text'">-->
<!--          {{ (r.meta && t(r.meta.name)) || r.name }}-->
<!--        </Button>-->
      </RouterLink>
    </div>
  </div>
</template>

<style lang="less" scoped>
@import "@/assets/main";
.nav {
  display: flex;
  justify-content: center; /* 水平居中 */
  //align-items: center;     /* 垂直居中 */
}
.bar {
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center;     /* 垂直居中 */
  padding: 4px 8px;
  margin: 0 8px;
  &:hover {
    border-radius: 10px; /* 所有角都是10px的圆角 */
    color: var(--btn-text-hover-color);
    background-color: var(--btn-text-hover-bg);
  }
}
.bar-icon_hover {
  .icon_hover();
}
.bar-active {
  border-radius: 10px; /* 所有角都是10px的圆角 */
  color: var(--btn-text-active-color);
  background-color: var(--btn-text-active-bg);
}
svg {
  width: 24px;
  height: 24px;
}
</style>
