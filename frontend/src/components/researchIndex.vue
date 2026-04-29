<script setup>
import {computed, h, onBeforeMount, onBeforeUnmount, onMounted,onUnmounted, ref,reactive} from 'vue'
import {GetAIResponseResultList} from "../../wailsjs/go/main/App";
import {NButton, NEllipsis, NText} from "naive-ui";
import ResearchReport from "./researchReport.vue";
import AiRecommendStocksList from "./aiRecommendStocksList.vue";
import PromptTemplateList from "./promptTemplateList.vue";
import AllStockList from "./allStockList.vue";
import AllStockInfoList from "./allStockInfoList.vue";
import CronTaskManager from "./cron-task-manager.vue";
import TradingRecordManager from "./TradingRecordManager.vue";
import StockChangesMonitor from "./stockChangesMonitor.vue";
import MCPServiceManager from "./mcp-server-manager.vue";
import SkillManager from "./skill-manager.vue";
import UplimitLadder from "./uplimitLadder.vue";
import PromptPlaza from "./promptPlaza.vue";
import {useI18n} from 'vue-i18n'
import {EventsOff, EventsOn} from "../../wailsjs/runtime";
import {useRoute} from 'vue-router'

const { t } = useI18n()

const nowTab = ref("aiAnalysisReport")
const route = useRoute()
onBeforeMount(() => {
  nowTab.value = route.query.name
})

onBeforeUnmount(() => {
  EventsOff("changeResearchTab")
})

onUnmounted(() => {

});

EventsOn("changeResearchTab", async (msg) => {
  console.log("changeResearchTab", msg)
  updateTab(msg.name)
})
function updateTab(name) {
  nowTab.value = name
}
</script>

<template>
  <n-card>
    <n-tabs type="line" animated @update-value="updateTab" :value="nowTab" style="--wails-draggable:no-drag">
      <n-tab-pane name="aiAnalysisReport" :tab="t('menu.aiAnalysisReport')">
        <ResearchReport/>
      </n-tab-pane>
      <n-tab-pane name="stockRecommendRecord" :tab="t('menu.stockRecommendRecord')">
        <AiRecommendStocksList/>
      </n-tab-pane>
      <n-tab-pane name="stockChangesMonitor" :tab="t('menu.stockChangesMonitor')">
        <StockChangesMonitor/>
      </n-tab-pane>
      <n-tab-pane name="uplimitLadder" :tab="t('menu.uplimitLadder')">
        <UplimitLadder/>
      </n-tab-pane>
      <n-tab-pane name="promptTemplate" :tab="t('menu.promptTemplate')">
        <PromptTemplateList/>
      </n-tab-pane>
      <n-tab-pane name="promptPlaza" :tab="t('menu.promptPlaza')">
        <PromptPlaza/>
      </n-tab-pane>
      <n-tab-pane name="stockInfoFilter" :tab="t('menu.stockInfoFilter')">
        <AllStockList/>
      </n-tab-pane>
      <n-tab-pane name="scheduledTask" :tab="t('menu.scheduledTask')">
        <CronTaskManager />
      </n-tab-pane>
      <n-tab-pane name="tradingLog" :tab="t('menu.tradingLog')">
        <TradingRecordManager />
      </n-tab-pane>
<!--      <n-tab-pane name="全部股票信息">-->
<!--        <AllStockInfoList/>-->
<!--      </n-tab-pane>-->
      <n-tab-pane name="mcpService" :tab="t('menu.mcpService')">
        <MCPServiceManager/>
      </n-tab-pane>
<!--      <n-tab-pane name="技能管理">-->
<!--        <SkillManager/>-->
<!--      </n-tab-pane>-->
    </n-tabs>
  </n-card>
</template>

<style scoped>
</style>
