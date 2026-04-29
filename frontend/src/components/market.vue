<script setup>
import * as echarts from "echarts";
import {computed, h, nextTick, onBeforeMount, onBeforeUnmount, onMounted,onUnmounted, ref} from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
import {
  GetAIResponseResult,
  GetConfig,
  GetIndustryRank,
  GetPromptTemplates,
  GetTelegraphList,
  GlobalStockIndexes,
  IsTradingTime,
  IsHKTradingTime,
  IsUSTradingTime,
  ReFleshTelegraphList,
  SaveAIResponseResult,
  SaveAsMarkdown,
  ShareAnalysis,
  SummaryStockNews,
  GetAiConfigs,
} from "../../wailsjs/go/main/App";
import {EventsOff, EventsOn} from "../../wailsjs/runtime";
import NewsList from "./newsList.vue";
import KLineChart from "./KLineChart.vue";
import StockLightweightKlineChart from "./StockLightweightKlineChart.vue";
import { CaretDown, CaretUp, PulseOutline,} from "@vicons/ionicons5";
import {NAvatar, NButton, NFlex, NText, useMessage, useNotification} from "naive-ui";
import {MdPreview} from "md-editor-v3";
import {useRoute} from 'vue-router'
import RankTable from "./rankTable.vue";
import IndustryMoneyRank from "./industryMoneyRank.vue";
import StockResearchReportList from "./StockResearchReportList.vue";
import StockNoticeList from "./StockNoticeList.vue";
import LongTigerRankList from "./LongTigerRankList.vue";
import IndustryResearchReportList from "./IndustryResearchReportList.vue";
import HotStockList from "./HotStockList.vue";
import HotEvents from "./HotEvents.vue";
import HotTopics from "./HotTopics.vue";
import InvestCalendarTimeLine from "./InvestCalendarTimeLine.vue";
import ClsCalendarTimeLine from "./ClsCalendarTimeLine.vue";
import SelectStock from "./SelectStock.vue";
import Stockhotmap from "./stockhotmap.vue";

const route = useRoute()
const icon = ref('https://raw.githubusercontent.com/ArvinLovegood/go-stock/master/build/appicon.png');

const message = useMessage()
const notify = useNotification()
const panelHeight = ref(window.innerHeight - 240)

const telegraphList = ref([])
const sinaNewsList = ref([])
const foreignNewsList = ref([])
const common = ref([])
const america = ref([])
const europe = ref([])
const asia = ref([])
const other = ref([])
const euStocks = ref([])
const globalStockIndexes = ref(null)
const summaryModal = ref(false)
const summaryBTN = ref(true)
const darkTheme = ref(false)
const httpProxyEnabled = ref(false)
const theme = computed(() => {
  return darkTheme ? 'dark' : 'light'
})
const aiSummary = ref(``)
const aiSummaryTime = ref("")
const modelName = ref("")
const chatId = ref("")
const question = ref(``)
const aiConfigId = ref(null)
const sysPromptId = ref(null)
const loading = ref(true)
const aiConfigs = ref([])
const sysPromptOptions = ref([])
const userPromptOptions = ref([])
const promptTemplates = ref([])
const industryRanks = ref([])
const sort = ref("0")
const nowTab = ref("marketNews")
const indexInterval = ref(null)
const indexIndustryRank = ref(null)
const tradingCheckInterval = ref(null)
const mdPreviewRef = ref(null)
const stockCode= ref('')
const enableTools= ref(true)
const thinkingMode = ref(false)
const treemapRef = ref(null);
let treemapchart =null;

function getIndex() {
  GlobalStockIndexes().then((res) => {
    globalStockIndexes.value = res
    common.value = res["common"]
    america.value = res["america"]
    europe.value = res["europe"]
    asia.value = res["asia"]
    other.value = res["other"]
  })
}

onBeforeMount(() => {
  nowTab.value = route.query.name
  stockCode.value = route.query.stockCode
  GetConfig().then(result => {
    summaryBTN.value = result.openAiEnable
    darkTheme.value = result.darkTheme
    httpProxyEnabled.value = result.httpProxyEnabled
  })
  GetPromptTemplates("", "").then(res => {
    promptTemplates.value = res
    sysPromptOptions.value = promptTemplates.value.filter(item => item.type === t('promptTemplateList.systemPromptType'))
    userPromptOptions.value = promptTemplates.value.filter(item => item.type === t('promptTemplateList.userPromptType'))
  })

  GetAiConfigs().then(res=>{
    aiConfigs.value = res
    aiConfigId.value = res[0].ID
  })
  GetEuronextStocks().then(res=>{
    euStocks.value = res
  })
  GetTelegraphList(t('market.caixiangTelegraph')).then((res) => {
    telegraphList.value = res
  })
  GetTelegraphList(t('market.sinaFinance')).then((res) => {
    sinaNewsList.value = res
  })
  GetTelegraphList(t('market.foreignMedia')).then((res) => {
    foreignNewsList.value = res
  })
  getIndex();
  industryRank();
  startTradingTimers();

  tradingCheckInterval.value = setInterval(async () => {
    const [cn, hk, us] = await Promise.all([
      IsTradingTime().catch(() => false),
      IsHKTradingTime().catch(() => false),
      IsUSTradingTime().catch(() => false)
    ])
    const anyTrading = cn || hk || us
    if (anyTrading && !indexInterval.value) {
      startTradingTimers()
    } else if (!anyTrading && indexInterval.value) {
      stopTradingTimers()
    }
  }, 60000)
})


onBeforeUnmount(() => {
  EventsOff("changeMarketTab")
  EventsOff("newTelegraph")
  EventsOff("newSinaNews")
  EventsOff("summaryStockNews")
  stopTradingTimers()
  if (tradingCheckInterval.value) {
    clearInterval(tradingCheckInterval.value)
  }
})

function startTradingTimers() {
  stopTradingTimers()
  indexInterval.value = setInterval(() => {
    getIndex()
  }, 3000)
  indexIndustryRank.value = setInterval(() => {
    industryRank()
    ReFlesh(t('market.caixiangTelegraph'))
    ReFlesh(t('market.sinaFinance'))
    ReFlesh(t('market.foreignMedia'))
  }, 1000 * 10)
}

function stopTradingTimers() {
  if (indexInterval.value) {
    clearInterval(indexInterval.value)
    indexInterval.value = null
  }
  if (indexIndustryRank.value) {
    clearInterval(indexIndustryRank.value)
    indexIndustryRank.value = null
  }
}

onUnmounted(() => {

});
EventsOn("changeMarketTab", async (msg) => {
  //message.info(msg.name)
  console.log(msg.name)
  updateTab(msg.name)
})

EventsOn("newTelegraph", (data) => {
  if (data!=null) {
    for (let i = 0; i < data.length; i++) {
      telegraphList.value.pop()
    }
    telegraphList.value.unshift(...data)
  }
})
EventsOn("newSinaNews", (data) => {
  if (data!=null) {
  for (let i = 0; i < data.length; i++) {
    sinaNewsList.value.pop()
  }
  sinaNewsList.value.unshift(...data)
  }
})
EventsOn("tradingViewNews", (data) => {
  if (data!=null) {
    for (let i = 0; i < data.length; i++) {
      foreignNewsList.value.pop()
    }
    foreignNewsList.value.unshift(...data)
  }
})

//获取页面高度
window.onresize = () => {
  panelHeight.value = window.innerHeight - 240
}

function getAreaName(code) {
  switch (code) {
    case "america":
      return t('market.america') || 'Americas'
    case "europe":
      return t('market.europe') || 'Europe'
    case "asia":
      return t('market.asia') || 'Asia'
    case "common":
      return t('market.common') || 'Common'
    case "other":
      return t('market.other') || 'Other'
  }
}

function changeIndustryRankSort() {
  if (sort.value === "0") {
    sort.value = "1"
  } else {
    sort.value = "0"
  }
  industryRank()
}

function industryRank() {

  GetIndustryRank(sort.value, 150).then(result => {
    if (result.length > 0) {
      //console.log(result)
      industryRanks.value = result
    } else {
      message.info(t('market.noData'))
    }
  })
}

function reAiSummary() {
  aiSummary.value = ""
  summaryModal.value = true
  loading.value = true
  SummaryStockNews(question.value,aiConfigId.value, sysPromptId.value,enableTools.value,thinkingMode.value,"summaryStockNews","")
}

function getAiSummary() {
  summaryModal.value = true
  loading.value = true
  GetAIResponseResult("市场资讯").then(result => {
    loading.value = false
    if (result.content) {
      aiSummary.value = result.content
      question.value = result.question
      loading.value = false

      const date = new Date(result.CreatedAt);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      aiSummaryTime.value = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
      modelName.value = result.modelName
    } else {
      aiSummaryTime.value = ""
      aiSummary.value = ""
      modelName.value = ""
      //SummaryStockNews(question.value, sysPromptId.value,enableTools.value)
    }
  })
}

function updateTab(name) {
  summaryBTN.value = (name === "marketNews");
  nowTab.value = name
}

EventsOn("summaryStockNews", async (msg) => {
  loading.value = false
  ////console.log(msg)
  if (msg === "DONE") {
    await SaveAIResponseResult("市场资讯", "市场资讯", aiSummary.value, chatId.value, question.value,aiConfigId.value)
    message.info(t('market.aiAnalysisComplete'))
    message.destroyAll()
    loading.value = false
  } else {
    if (msg.chatId) {
      chatId.value = msg.chatId
    }
    if (msg.question) {
      question.value = msg.question
    }
    if (msg.content) {
      aiSummary.value = aiSummary.value + msg.content
    }
    if (msg.reasoning_content) {
      aiSummary.value = aiSummary.value + msg.reasoning_content
    }
    if (msg.extraContent) {
      aiSummary.value = aiSummary.value + msg.extraContent
    }
    if (msg.model) {
      modelName.value = msg.model
    }
    if (msg.time) {
      aiSummaryTime.value = msg.time
    }
    loading.value = true
    scrollToAiResultBottom()
  }
})

function scrollToAiResultBottom() {
  nextTick(() => {
    const previewEl = mdPreviewRef.value?.$el
    if (previewEl) {
      const scrollContainer = previewEl.querySelector('.md-editor-preview-wrapper') || 
                               previewEl.querySelector('.md-editor-preview') ||
                               previewEl
      if (scrollContainer) {
        scrollContainer.scrollTop = scrollContainer.scrollHeight
      }
    }
  })
}

async function copyToClipboard() {
  try {
    await navigator.clipboard.writeText(aiSummary.value);
    message.success(t('market.analysisResultCopied'));
  } catch (err) {
    message.error(t('market.copyFailed') + err);
  }
}

function saveAsMarkdown() {
  SaveAsMarkdown('市场资讯', '市场资讯').then(result => {
    message.success(result)
  })
}

function share() {
  ShareAnalysis('市场资讯', '市场资讯').then(msg => {
    //message.info(msg)
    notify.info({
      avatar: () =>
          h(NAvatar, {
            size: 'small',
            round: false,
            src: icon.value
          }),
      title: t('market.shareToCommunity'),
      duration: 1000 * 30,
      content: () => {
        return h('div', {
          style: {
            'text-align': 'left',
            'font-size': '14px',
          }
        }, {default: () => msg})
      },
    })
  })
}

function ReFlesh(source) {
  //console.log("ReFlesh:", source)
  ReFleshTelegraphList(source).then(res => {
    if (source === t('market.caixiangTelegraph')) {
      telegraphList.value = res
    }
    if (source === t('market.sinaFinance')) {
      sinaNewsList.value = res
    }
    if (source === t('market.foreignMedia')) {
      foreignNewsList.value = res
    }
  })
}
</script>

<template>
  <n-card>
    <n-tabs type="line" animated @update-value="updateTab" :value="nowTab" style="--wails-draggable:no-drag">
      <n-tab-pane name="marketNews" :tab="t('market.marketNews')">
        <n-grid :cols="1" :y-gap="0">
          <n-gi>
            <AnalyzeMartket :dark-theme="darkTheme" :chart-height="300" :kDays="1" :name="t('market.recent24hHotWords')" />
          </n-gi>
          <n-gi>
            <n-grid :cols="foreignNewsList.length?3:2" :y-gap="0">
              <n-gi>
                <news-list :newsList="telegraphList" :header-title="t('market.caixiangTelegraph')" @update:message="ReFlesh"></news-list>
              </n-gi>
              <n-gi>
                <news-list :newsList="sinaNewsList" :header-title="t('market.sinaFinance')" @update:message="ReFlesh"></news-list>
              </n-gi>
              <n-gi v-if="foreignNewsList.length>0">
                <news-list :newsList="foreignNewsList" :header-title="t('market.foreignMedia')" @update:message="ReFlesh"></news-list>
              </n-gi>

            </n-grid>
          </n-gi>
        </n-grid>

      </n-tab-pane>
      <n-tab-pane name="globalIndexes" :tab="t('market.globalIndexes')">
        <n-tabs type="segment" animated>
          <n-tab-pane name="globalIndices" :tab="t('market.globalIndices')">
            <n-grid :cols="5" :y-gap="0">
              <n-gi v-for="(val, key) in globalStockIndexes" :key="key">
                <n-list bordered>
                  <template #header>
                    {{ getAreaName(key) }}
                  </template>
                  <n-list-item v-for="item in val" :key="item.code">
                    <n-grid :cols="3" :y-gap="0">
                      <n-gi>

                        <n-text :type="item.zdf>0?'error':'success'">
                          <n-image :src="item.img" width="20"/> &nbsp;{{ item.name }}
                        </n-text>
                      </n-gi>
                      <n-gi>
                        <n-text :type="item.zdf>0?'error':'success'">{{ item.zxj }}</n-text>&nbsp;
                        <n-text :type="item.zdf>0?'error':'success'">
                          <n-number-animation :precision="2" :from="0" :to="item.zdf"/>
                          %
                        </n-text>

                      </n-gi>
                      <n-gi>
                        <n-text :type="item.state === 'open' ? 'success' : 'warning'">{{
                            item.state === 'open' ? t('market.open') : t('market.closed')
                          }}
                        </n-text>
                      </n-gi>
                    </n-grid>
                  </n-list-item>
                </n-list>
              </n-gi>
            </n-grid>
          </n-tab-pane>
          <n-tab-pane name="shangzhengIndex" :tab="t('market.shangzheng')">
            <k-line-chart code="sh000001" :chart-height="panelHeight" :stock-name="t('market.shangzheng')" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="shenzhengIndex" :tab="t('market.shenzheng')">
            <k-line-chart code="sz399001" :chart-height="panelHeight" :stock-name="t('market.shenzheng')" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="chuangyeIndex" :tab="t('market.chuangye')">
            <k-line-chart code="sz399006" :chart-height="panelHeight" :stock-name="t('market.chuangye')" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="hengshengIndex" :tab="t('market.hengsheng')">
            <k-line-chart code="hkHSI" :chart-height="panelHeight" :stock-name="t('market.hengsheng')" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="nasdaqIndex" :tab="t('market.nasdaq')">
            <k-line-chart code="us.IXIC" :chart-height="panelHeight" :stock-name="t('market.nasdaq')" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="djiaIndex" :tab="t('market.djia')">
            <k-line-chart code="us.DJI" :chart-height="panelHeight" :stock-name="t('market.djia')" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="sp500Index" :tab="t('market.sp500')">
            <k-line-chart code="us.INX" :chart-height="panelHeight" :stock-name="t('market.sp500')" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
      <n-tab-pane name="majorIndexes" :tab="t('market.majorIndexes')">
        <n-tabs type="segment" animated>

<!--          <n-tab-pane name="西部数据" tab="西部数据">-->
<!--            <StockLightweightKlineChart code="105.WDC" :chart-height="panelHeight" stock-name="西部数据"-->
<!--                                        :dark-theme="true"></StockLightweightKlineChart>-->
<!--          </n-tab-pane>-->

          <n-tab-pane name="shangzhengIndex" :tab="t('market.shangzheng')"  >
            <StockLightweightKlineChart code="000001.SH" :chart-height="panelHeight-130" :stock-name="t('market.shangzheng')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="shenzhengComposite" :tab="t('market.shenzheng')"  >
            <StockLightweightKlineChart code="399001.SZ" :chart-height="panelHeight-130" :stock-name="t('market.shenzheng')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="chuangyeIndex" :tab="t('market.chuangye')"  >
            <StockLightweightKlineChart code="399006.SZ" :chart-height="panelHeight-130" :stock-name="t('market.chuangye')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>

          <n-tab-pane name="hengshengIndex" :tab="t('market.hengsheng')">
            <StockLightweightKlineChart code="100.HSI" :chart-height="panelHeight" :stock-name="t('market.hengsheng')"
                                        :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="djiaIndex" :tab="t('market.djia')">
            <StockLightweightKlineChart code="100.DJIA" :chart-height="panelHeight" :stock-name="t('market.djia')"
                                        :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="sp500Index" :tab="t('market.sp500')">
            <StockLightweightKlineChart code="100.SPX" :chart-height="panelHeight" :stock-name="t('market.sp500')"
                                        :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="nasdaqIndex" :tab="t('market.nasdaq')">
            <StockLightweightKlineChart code="100.NDX" :chart-height="panelHeight" :stock-name="t('market.nasdaq')"
                                        :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>

          <n-tab-pane name="huShen300" :tab="t('market.huShen300')">
            <StockLightweightKlineChart code="000300.SH" :chart-height="panelHeight-130" :stock-name="t('market.huShen300')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="shangzheng50" :tab="t('market.shangzheng50')">
            <StockLightweightKlineChart code="000016.SH" :chart-height="panelHeight-130" :stock-name="t('market.shangzheng50')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="zhongzhengA500" :tab="t('market.zhongzhengA500')">
            <StockLightweightKlineChart code="000510.SH" :chart-height="panelHeight-130" :stock-name="t('market.zhongzhengA500')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="zhongzheng1000" :tab="t('market.zhongzheng1000')">
            <StockLightweightKlineChart code="000852.SH" :chart-height="panelHeight-130" :stock-name="t('market.zhongzheng1000')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>

          <n-tab-pane name="kechuang50" :tab="t('market.kechuang50')"  >
            <StockLightweightKlineChart code="000688.SH" :chart-height="panelHeight-130" :stock-name="t('market.kechuang50')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="kechuangChip" :tab="t('market.kechuangChip')"  >
            <StockLightweightKlineChart code="000685.SH" :chart-height="panelHeight-130" :stock-name="t('market.kechuangChip')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="securitiesLead" :tab="t('market.securitiesLead')"  >
            <StockLightweightKlineChart code="399437.SZ" :chart-height="panelHeight-130" :stock-name="t('market.securitiesLead')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="highEndEquipment" :tab="t('market.highEndEquipment')"  >
            <StockLightweightKlineChart code="399437.SZ" :chart-height="panelHeight-130" :stock-name="t('market.highEndEquipment')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="zhongzhengBank" :tab="t('market.zhongzhengBank')">
            <StockLightweightKlineChart code="399986.SZ" :chart-height="panelHeight-130" :stock-name="t('market.zhongzhengBank')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="shangzhengMedicine" :tab="t('market.shangzhengMedicine')">
            <StockLightweightKlineChart code="000037.SH" :chart-height="panelHeight-130" :stock-name="t('market.shangzhengMedicine')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="zhongzhengBaijiu" :tab="t('market.zhongzhengBaijiu')">
            <StockLightweightKlineChart code="399997.SZ" :chart-height="panelHeight-130" :stock-name="t('market.zhongzhengBaijiu')" :dark-theme="true"></StockLightweightKlineChart>
          </n-tab-pane>
          <n-tab-pane name="ftseChina3xLong" :tab="t('market.ftseChina3xLong')">
            <k-line-chart code="usYINN.AM" :chart-height="panelHeight" :stock-name="t('market.ftseChina3xLong')" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="vixFearIndex" :tab="t('market.vixFearIndex')">
            <k-line-chart code="usUVXY.AM" :chart-height="panelHeight" :stock-name="t('market.vixFearIndex')" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="european" :tab="t('market.european')">
            <n-grid :cols="2" :x-gap="12" :y-gap="12">
              <n-gi>
                <StockLightweightKlineChart code="eu:DAX" :chart-height="(panelHeight-130)/2" :stock-name="t('market.dax')" :dark-theme="true"></StockLightweightKlineChart>
              </n-gi>
              <n-gi>
                <StockLightweightKlineChart code="eu:CAC" :chart-height="(panelHeight-130)/2" :stock-name="t('market.cac40')" :dark-theme="true"></StockLightweightKlineChart>
              </n-gi>
              <n-gi>
                <StockLightweightKlineChart code="uk:FTSE" :chart-height="(panelHeight-130)/2" :stock-name="t('market.ftse100')" :dark-theme="true"></StockLightweightKlineChart>
              </n-gi>
              <n-gi>
                <StockLightweightKlineChart code="ch:SMI" :chart-height="(panelHeight-130)/2" :stock-name="t('market.smi')" :dark-theme="true"></StockLightweightKlineChart>
              </n-gi>
            </n-grid>
          </n-tab-pane>
          <n-divider>{{ t('market.europeanStocks') }}</n-divider>
          <n-tab-pane name="euronextStocks" :tab="t('market.euronext')">
            <n-grid :cols="4" :x-gap="12" :y-gap="12">
              <n-gi v-for="stock in euStocks.filter(s => s.exchange === 'Euronext')" :key="stock.code">
                <StockLightweightKlineChart :code="'eu:' + stock.code" :chart-height="panelHeight-180" :stock-name="stock.eName" :dark-theme="true"></StockLightweightKlineChart>
              </n-gi>
            </n-grid>
          </n-tab-pane>
          <n-tab-pane name="lseStocks" :tab="t('market.lse')">
            <n-grid :cols="4" :x-gap="12" :y-gap="12">
              <n-gi v-for="stock in euStocks.filter(s => s.exchange === 'LSE')" :key="stock.code">
                <StockLightweightKlineChart :code="'uk:' + stock.code" :chart-height="panelHeight-180" :stock-name="stock.eName" :dark-theme="true"></StockLightweightKlineChart>
              </n-gi>
            </n-grid>
          </n-tab-pane>
          <n-tab-pane name="swissStocks" :tab="t('market.swissExchange')">
            <n-grid :cols="4" :x-gap="12" :y-gap="12">
              <n-gi v-for="stock in euStocks.filter(s => s.exchange === 'Swiss')" :key="stock.code">
                <StockLightweightKlineChart :code="'ch:' + stock.code" :chart-height="panelHeight-180" :stock-name="stock.eName" :dark-theme="true"></StockLightweightKlineChart>
              </n-gi>
            </n-grid>
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
      <n-tab-pane name="industryRanking" :tab="t('market.industryRanking')">
        <n-tabs type="card" animated>
          <n-tab-pane name="industryRiseRanking" :tab="t('market.industryRiseRanking')">
            <n-table striped>
              <n-thead>
                <n-tr>
                  <n-th>{{ t('market.industryName') }}</n-th>
                  <n-th @click="changeIndustryRankSort">{{ t('market.industryRise') }}
                    <n-icon v-if="sort==='0'" :component="CaretDown"/>
                    <n-icon v-if="sort==='1'" :component="CaretUp"/>
                  </n-th>
                  <n-th>{{ t('market.industry5DayRise') }}</n-th>
                  <n-th>{{ t('market.industry20DayRise') }}</n-th>
                  <n-th>{{ t('market.leadingStock') }}</n-th>
                  <n-th>{{ t('market.rise') }}</n-th>
                  <n-th>{{ t('market.latestPrice') }}</n-th>
                </n-tr>
              </n-thead>
              <n-tbody>
                <n-tr v-for="item in industryRanks" :key="item.bd_code">
                  <n-td>
                    <n-tag :bordered=false type="info">{{ item.bd_name }}</n-tag>
                  </n-td>
                  <n-td>
                    <n-text :type="item.bd_zdf>0?'error':'success'">{{ item.bd_zdf }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.bd_zdf5>0?'error':'success'">{{ item.bd_zdf5 }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.bd_zdf20>0?'error':'success'">{{ item.bd_zdf20 }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.nzg_zdf>0?'error':'success'"> {{ item.nzg_name }}
                      <n-text type="info">{{ item.nzg_code }}</n-text>
                    </n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.nzg_zdf>0?'error':'success'"> {{ item.nzg_zdf }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.nzg_zdf>0?'error':'success'">{{ item.nzg_zxj }}</n-text>
                  </n-td>
                </n-tr>
              </n-tbody>
            </n-table>
            <n-table striped>
              <n-thead>
                <n-tr>
                  <n-th>{{ t('market.industryName') }}</n-th>
                  <n-th @click="changeIndustryRankSort">{{ t('market.industryRise') }}
                    <n-icon v-if="sort==='0'" :component="CaretDown"/>
                    <n-icon v-if="sort==='1'" :component="CaretUp"/>
                  </n-th>
                  <n-th>{{ t('market.industry5DayRise') }}</n-th>
                  <n-th>{{ t('market.industry20DayRise') }}</n-th>
                  <n-th>{{ t('market.leadingStock') }}</n-th>
                  <n-th>{{ t('market.rise') }}</n-th>
                  <n-th>{{ t('market.latestPrice') }}</n-th>
                </n-tr>
              </n-thead>
              <n-tbody>
                <n-tr v-for="item in industryRanks" :key="item.bd_code">
                  <n-td>
                    <n-tag :bordered=false type="info">{{ item.bd_name }}</n-tag>
                  </n-td>
                  <n-td>
                    <n-text :type="item.bd_zdf>0?'error':'success'">{{ item.bd_zdf }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.bd_zdf5>0?'error':'success'">{{ item.bd_zdf5 }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.bd_zdf20>0?'error':'success'">{{ item.bd_zdf20 }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.nzg_zdf>0?'error':'success'"> {{ item.nzg_name }}
                      <n-text type="info">{{ item.nzg_code }}</n-text>
                    </n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.nzg_zdf>0?'error':'success'"> {{ item.nzg_zdf }}%</n-text>
                  </n-td>
                  <n-td>
                    <n-text :type="item.nzg_zdf>0?'error':'success'">{{ item.nzg_zxj }}</n-text>
                  </n-td>
                </n-tr>
              </n-tbody>
            </n-table>
          </n-tab-pane>
          <n-tab-pane name="industryMoneyRanking" :tab="t('market.industryMoneyRanking')">
            <industryMoneyRank :fenlei="'0'" :header-title="t('market.industryMoneyRanking')" :sort="'netamount'"/>
          </n-tab-pane>
          <n-tab-pane name="csrcIndustryRanking" :tab="t('market.csrcIndustryRanking')">
            <industryMoneyRank :fenlei="'2'" :header-title="t('market.csrcIndustryRanking')" :sort="'netamount'"/>
          </n-tab-pane>
          <n-tab-pane name="conceptMoneyRanking" :tab="t('market.conceptMoneyRanking')">
            <industryMoneyRank :fenlei="'1'" :header-title="t('market.conceptMoneyRanking')" :sort="'netamount'"/>
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
      <n-tab-pane name="stockFundFlow" :tab="t('market.stockFundFlow')">
        <n-tabs type="card" animated>
          <n-tab-pane name="netamount" :tab="t('market.netInflowRanking')">
            <RankTable :header-title="t('market.netInflowRanking')" :sort="'netamount'"/>
          </n-tab-pane>
          <n-tab-pane name="outamount" :tab="t('market.outflowRanking')">
            <RankTable :header-title="t('market.outflowRanking')" :sort="'outamount'"/>
          </n-tab-pane>
          <n-tab-pane name="ratioamount" :tab="t('market.netInflowRateRanking')">
            <RankTable :header-title="t('market.netInflowRateRanking')" :sort="'ratioamount'"/>
          </n-tab-pane>
          <n-tab-pane name="r0_net" :tab="t('market.mainForceNetInflowRanking')">
            <RankTable :header-title="t('market.mainForceNetInflowRanking')" :sort="'r0_net'"/>
          </n-tab-pane>
          <n-tab-pane name="r0_out" :tab="t('market.mainForceOutflowRanking')">
            <RankTable :header-title="t('market.mainForceOutflowRanking')" :sort="'r0_out'"/>
          </n-tab-pane>
          <n-tab-pane name="r0_ratio" :tab="t('market.mainForceNetInflowRateRanking')">
            <RankTable :header-title="t('market.mainForceNetInflowRateRanking')" :sort="'r0_ratio'"/>
          </n-tab-pane>
          <n-tab-pane name="r3_net" :tab="t('market.retailNetInflowRanking')">
            <RankTable :header-title="t('market.retailNetInflowRanking')" :sort="'r3_net'"/>
          </n-tab-pane>
          <n-tab-pane name="r3_out" :tab="t('market.retailOutflowRanking')">
            <RankTable :header-title="t('market.retailOutflowRanking')" :sort="'r3_out'"/>
          </n-tab-pane>
          <n-tab-pane name="r3_ratio" :tab="t('market.retailNetInflowRateRanking')">
            <RankTable :header-title="t('market.retailNetInflowRateRanking')" :sort="'r3_ratio'"/>
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
      <n-tab-pane name="dragonTigerList" :tab="t('market.dragonTigerList')">
        <LongTigerRankList />
      </n-tab-pane>
      <n-tab-pane name="stockResearchReport" :tab="t('market.stockResearchReport')">
        <StockResearchReportList :stock-code="stockCode"/>
      </n-tab-pane>
      <n-tab-pane name="companyAnnouncement" :tab="t('market.companyAnnouncement')">
        <StockNoticeList :stock-code="stockCode" />
      </n-tab-pane>
      <n-tab-pane name="industryResearch" :tab="t('market.industryResearch')">
        <IndustryResearchReportList/>
      </n-tab-pane>
      <n-tab-pane name="currentHot" :tab="t('market.currentHot')">
        <n-tabs type="card" animated>
          <n-tab-pane name="global" :tab="t('market.global')">
            <HotStockList :market-type="'10'"/>
          </n-tab-pane>
          <n-tab-pane name="hushi" :tab="t('market.hushi')">
            <HotStockList :market-type="'12'"/>
          </n-tab-pane>
          <n-tab-pane name="ganggu" :tab="t('market.ganggu')">
            <HotStockList :market-type="'13'"/>
          </n-tab-pane>
          <n-tab-pane name="meigu" :tab="t('market.meigu')">
            <HotStockList :market-type="'11'"/>
          </n-tab-pane>
          <n-tab-pane name="hotTopics" :tab="t('market.hotTopics')">
            <n-grid :cols="1" :y-gap="10">
              <n-grid-item>
                <HotTopics/>
              </n-grid-item>
<!--              <n-grid-item>-->
<!--                <HotEvents/>-->
<!--              </n-grid-item>-->
            </n-grid>
          </n-tab-pane>
          <n-tab-pane name="majorEventsTimeline" :tab="t('market.majorEventsTimeline')">
            <InvestCalendarTimeLine />
          </n-tab-pane>
          <n-tab-pane name="financialCalendar" :tab="t('market.financialCalendar')">
            <ClsCalendarTimeLine />
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
      <n-tab-pane name="indicatorStock" :tab="t('market.indicatorStock')">
        <select-stock />
      </n-tab-pane>
      <n-tab-pane name="famousStations" :tab="t('market.famousStations')">
        <Stockhotmap />
      </n-tab-pane>
    </n-tabs>
  </n-card>
  <n-modal transform-origin="center" v-model:show="summaryModal" preset="card" style="width: 800px;"
           :title="t('market.aiSummaryTitle')">
    <n-spin size="small" :show="loading">
      <MdPreview ref="mdPreviewRef" style="height: 440px;text-align: left" :modelValue="aiSummary" :theme="theme"/>
    </n-spin>
    <template #footer>
      <n-flex justify="space-between" ref="tipsRef">
        <n-text type="info" v-if="aiSummaryTime">
          <n-tag v-if="modelName" type="warning" round :title="chatId" :bordered="false">{{ modelName }}</n-tag>
          {{ aiSummaryTime }}
        </n-text>
        <n-text type="error">* {{ t('market.disclaimer') }}</n-text>
      </n-flex>
    </template>
    <template #action>
      <n-flex justify="left" style="margin-bottom: 10px">
        <n-switch v-model:value="enableTools" :round="false">
          <template #checked>
            {{ t('market.toolCall') }}
          </template>
          <template #unchecked>
            {{ t('market.nonToolCall') }}
          </template>
        </n-switch>
        <n-switch v-model:value="thinkingMode" :round="false">
          <template #checked>
            {{ t('market.thinkingMode') }}
          </template>
          <template #unchecked>
            {{ t('market.nonThinkingMode') }}
          </template>
        </n-switch>


        <n-gradient-text type="error" style="margin-left: 10px">* {{ t('market.toolCallTip') }}</n-gradient-text>
      </n-flex>
      <n-flex justify="space-between" style="margin-bottom: 10px">
        <n-select style="width: 32%" v-model:value="aiConfigId" label-field="name" value-field="ID"
                  :options="aiConfigs" :placeholder="t('market.selectAiModel')"/>
        <n-select style="width: 32%" v-model:value="sysPromptId" label-field="name" value-field="ID"
                  :options="sysPromptOptions" :placeholder="t('market.selectSystemPrompt')"/>
        <n-select style="width: 32%" v-model:value="question" label-field="name" value-field="content"
                  :options="userPromptOptions" :placeholder="t('market.selectUserPrompt')"/>
      </n-flex>
      <n-flex justify="right">
        <n-input v-model:value="question" style="text-align: left" clearable
                 type="textarea"
                 :show-count="true"
                 :placeholder="t('market.enterQuestion')"
                 :autosize="{
              minRows: 2,
              maxRows: 5
            }"
        />
        <n-button size="tiny" type="warning" @click="reAiSummary">{{ t('market.againSummary') }}</n-button>
        <n-button size="tiny" type="success" @click="copyToClipboard">{{ t('market.copyToClipboard') }}</n-button>
        <n-button size="tiny" type="primary" @click="saveAsMarkdown">{{ t('market.saveAsMarkdown') }}</n-button>
        <n-button size="tiny" type="error" @click="share">{{ t('market.shareToCommunity') }}</n-button>
      </n-flex>
    </template>
  </n-modal>

  <div style="position: fixed;bottom: 18px;right:25px;z-index: 10;" v-if="summaryBTN">
    <n-input-group>
      <n-button type="primary" @click="getAiSummary">
        <n-icon :component="PulseOutline"/> &nbsp;{{ t('market.aiSummaryTitle') }}
      </n-button>
    </n-input-group>
  </div>



</template>
<style scoped>
</style>