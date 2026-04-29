<script setup lang="ts">
import {h} from 'vue'
import {NTag,NImage} from 'naive-ui'
import EmbeddedUrl from "./EmbeddedUrl.vue";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <n-tabs  type="line" animated>
    <n-tab-pane name="selectStock" :tab="t('stockHotmap.selectStock')">
      <embedded-url url="https://xuangutong.com.cn"  :height="'calc(100vh - 252px)'"/>
    </n-tab-pane>
    <n-tab-pane name="baidu" :tab="t('stockHotmap.baidu')">
      <embedded-url url="https://gushitong.baidu.com"  :height="'calc(100vh - 252px)'"/>
    </n-tab-pane>
    <n-tab-pane name="eastmoney" :tab="t('stockHotmap.eastmoney')">
      <embedded-url url="https://quote.eastmoney.com/stockhotmap/"  :height="'calc(100vh - 252px)'"/>
    </n-tab-pane>
    <n-tab-pane name="tophub" :tab="t('stockHotmap.tophub')">
      <embedded-url url="https://tophub.today/c/finance"  :height="'calc(100vh - 252px)'"/>
    </n-tab-pane>
    <n-tab-pane name="moyu" :tab="t('stockHotmap.moyu')">
      <embedded-url url="https://996.ninja/"  :height="'calc(100vh - 252px)'"/>
    </n-tab-pane>
    <n-tab-pane name="cls" :tab="t('stockHotmap.cls')">
      <embedded-url url="https://www.cls.cn/quotation"  :height="'calc(100vh - 252px)'"/>
    </n-tab-pane>
    <n-tab-pane name="msgWall" :tab="t('stockHotmap.msgWall')">
      <embedded-url url="https://go-stock.sparkmemory.top:16667/go-stock"  :height="'calc(100vh - 252px)'"/>
    </n-tab-pane>
    <n-tab-pane name="gold" :tab="t('stockHotmap.gold')">
      <embedded-url url="https://www.tradinghero.com/chart?symbol=XAUUSD.GOODS"  :height="'calc(100vh - 252px)'"/>
    </n-tab-pane>

    <n-tab-pane name="jin10Futures" :tab="t('stockHotmap.jin10Futures')">
      <embedded-url url="https://qihuo.jin10.com/"  :height="'calc(100vh - 252px)'"/>
    </n-tab-pane>
    <n-tab-pane name="jin10Data" :tab="t('stockHotmap.jin10Data')">
      <embedded-url url="https://datas.jin10.com/#/category/52005"  :height="'calc(100vh - 252px)'"/>
    </n-tab-pane>


    <n-tab-pane name="recommendMore" :tab="t('stockHotmap.recommendMore')">
    </n-tab-pane>
<!--    <n-tab-pane name="自在量化" tab="自在量化">-->
<!--      <embedded-url url="https://quant.zizizaizai.com/home"/>-->
<!--    </n-tab-pane>-->
  </n-tabs>
</template>

<style scoped>

</style>