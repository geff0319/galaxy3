<script setup lang="ts" xmlns="http://www.w3.org/1999/html">
import { computed, ref } from 'vue'
import { useI18n, I18nT } from 'vue-i18n'
import { useMessage, useBool } from '@/hooks'
import { TencentTextTranslate } from '@/bridge'
import { useTranslateStore } from '@/stores'

const { t } = useI18n()
const { message } = useMessage()
const translateStore = useTranslateStore()


const onSubmit =async () => {
  if(translateStore.trans.originalText === ''){
    return
  }
  console.log(translateStore.trans.sourceLanguage)
  console.log(translateStore.trans.targetLanguage)
  console.log(translateStore.trans.originalText)
  try{
    translateStore.loading = true
    const res = await TencentTextTranslate(translateStore.trans.originalText,translateStore.trans.sourceLanguage,translateStore.trans.targetLanguage)
    console.log(res)
    translateStore.loading = false
    translateStore.trans.translationText = res
  } catch (error: any) {
    translateStore.loading = false
    console.error(error)
    message.error(error)
  }
}
const change = (v:any)=>{
  console.log(v)
}
</script>

<template>
  <div class="grid-list-header">
<!--    翻译接口：-->
<!--    <Select v-model="source" :options="sourceOptions" size="small" />-->
  <div class="select-language">
    <a-select
        ref="select"
        v-model:value=translateStore.trans.sourceLanguage
        style="width: 140px"
        size="small"
    >
      <a-select-option v-for="item in translateStore.sourceLanguageOptions()" :key="item.value" :value="item.value">
        <div style="display: flex;flex-direction: row;align-items: center">
          <Icon style="width: 20px; height: 20px;margin-right:5px;" :icon="item.value"/>
          {{ item.label }}
        </div>
      </a-select-option>
    </a-select>
    <Icon style="flex-shrink: 0" icon="exchange"/>
    <a-select
        ref="select"
        v-model:value=translateStore.trans.targetLanguage
        style="width: 140px"
        size="small"
        @change = "change"
    >
      <a-select-option v-for="item in translateStore.trans.languageOptions" :key="item.value" :value="item.value">
        <div style="display: flex;flex-direction: row;align-items: center">
          <Icon style="width: 20px; height: 20px;margin-right:5px;" :icon="item.value"/>
          {{ item.label }}
        </div>
      </a-select-option>
    </a-select>
    <div class="grid-list-header">
      <Button @click="onSubmit" type="primary" class="ml-auto" :loading="translateStore.loading">立即翻译</Button>
    </div>
  </div>
  </div>
  <textarea v-model="translateStore.trans.originalText" placeholder="请输入原文" class="textarea-item" style="background-color: rgb(243,243,243);cursor:text;"></textarea>
  <textarea v-model="translateStore.trans.translationText"  disabled class="textarea-item"></textarea>

</template>

<style lang="less" scoped>
.textarea-item {
  resize: none;
  transition:
      box-shadow 0.4s,
      background-color 0.4s;
  font-size: 14px;
  font-family: "幼圆", "Yu Yuan", sans-serif; /* 设置幼圆字体 */
  height: 50%;
  margin: 5px 0;
  background-color: #dcdcdc;
  width: 100%;
  padding: 7px;
  border: 1px solid  #ccc;
  border-radius: 4px;
  box-sizing: border-box;
  outline: none; /* 去除默认聚焦边框 */
}
.select-language {
  min-width: 350px;
  display: flex;
  justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
  margin-left: auto; /* 将左边距设置为自动 */
  margin-right: auto; /* 将右边距设置为自动，以确保元素水平居中 */
  /*margin-left: 80px;*/
  border-radius: 8px;
  font-size: 12px;
  /*background-color: #e6e6e6; !* 灰色背景 *!*/
}
</style>
