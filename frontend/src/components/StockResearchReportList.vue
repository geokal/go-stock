<script setup>
import {onBeforeMount, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {GetStockList, StockResearchReport} from "../../wailsjs/go/main/App";
import {ArrowDownOutline, CaretDown, CaretUp, PulseOutline, Refresh, RefreshCircleSharp,} from "@vicons/ionicons5";

import KLineChart from "./KLineChart.vue";
import MoneyTrend from "./moneyTrend.vue";
import {useMessage} from "naive-ui";
import {BrowserOpenURL} from "../../wailsjs/runtime";

const { t } = useI18n()
const {stockCode}=defineProps(
    {
      stockCode: {
        type: String,
        default: ''
      }
    }
)

const message=useMessage()
const list  = ref([])

const options =  ref([])

function getStockResearchReport(value) {
  StockResearchReport(value).then(result => {
    list.value = result
  })
}

onBeforeMount(()=>{
  getStockResearchReport(stockCode);
})

function ratingChangeName(ratingChange){
  if(ratingChange===0){
    return t('stockResearchReportList.adjustUp')
  }else if(ratingChange===1){
    return t('stockResearchReportList.adjustDown')
  }else if(ratingChange===2){
    return t('stockResearchReportList.first')
  }else if(ratingChange===3){
    return t('stockResearchReportList.maintain')
  }else if (ratingChange===4){
    return t('stockResearchReportList.noChange')
  }else{
    return ''
  }
}
function getmMarketCode(market,code) {
  if(market==="SHENZHEN"){
    return "sz"+code
  }else if(market==="SHANGHAI"){
    return "sh"+code
  }else if(market==="BEIJING"){
    return "bj"+code
  }else if(market==="HONGKONG"){
    return "hk"+code
  }else{
    return code
  }
}
function openWin(code) {
  BrowserOpenURL("https://pdf.dfcfw.com/pdf/H3_"+code+"_1.pdf?1749744888000.pdf")
}

function findStockList(query){
  if (query){
    GetStockList(query).then(result => {
      options.value=result.map(item => {
        return {
          label: item.name+" - "+item.ts_code,
          value: item.ts_code
        }
      })
    })
  }else{
    getStockResearchReport('')
  }
}
function handleSearch(value) {
  getStockResearchReport(value)
}
</script>

<template>
  <n-card>
    <n-auto-complete  :options="options" :placeholder="t('stockResearchReportList.enterStockName')"  clearable filterable  :on-select="handleSearch" :on-update:value="findStockList"  />
  </n-card>
  <n-table striped size="small">
    <n-thead>
      <n-tr>
        <n-th>{{ t('stockResearchReportList.name') }}</n-th>
        <n-th>{{ t('stockResearchReportList.industry') }}</n-th>
        <n-th>{{ t('stockResearchReportList.title') }}</n-th>
        <n-th>{{ t('stockResearchReportList.emRating') }}</n-th>
        <n-th>{{ t('stockResearchReportList.ratingChange') }}</n-th>
        <n-th>{{ t('stockResearchReportList.orgRating') }}</n-th>
        <n-th>{{ t('stockResearchReportList.analyst') }}</n-th>
        <n-th>{{ t('stockResearchReportList.organization') }}</n-th>
        <n-th> <n-flex justify="space-between">{{ t('stockResearchReportList.date') }}<n-icon @click="getStockResearchReport(stockCode)" color="#409EFF" :size="20"  :component="RefreshCircleSharp"/></n-flex></n-th>
      </n-tr>
    </n-thead>
    <n-tbody>
      <n-tr v-for="item in list" :key="item.infoCode">
        <n-td :title="item.stockCode">
          <n-popover trigger="hover" placement="right">
            <template #trigger>
              <n-tag type="info"  :bordered="false">{{item.stockName}}</n-tag>
            </template>
            <k-line-chart style="width: 800px" :code="getmMarketCode(item.market,item.stockCode)" :chart-height="500" :stockName="item.stockName" :k-days="20" :dark-theme="true"></k-line-chart>
          </n-popover>
        </n-td>
        <n-td><n-tag type="info"  :bordered="false">{{item.indvInduName}}</n-tag></n-td>
        <n-td>
          <n-a type="info"  @click="openWin(item.infoCode)">{{item.title}}</n-a>
        </n-td>
        <n-td><n-text :type="item.emRatingName===t('stockResearchReportList.increaseHold')?'error':'info'">
          {{item.emRatingName}}
        </n-text></n-td>
        <n-td><n-text :type="item.ratingChange===0?'error':'info'">{{ratingChangeName(item.ratingChange)}}</n-text></n-td>
        <n-td>{{item.sRatingName}}</n-td>
        <n-td>{{item.researcher}}</n-td>
        <n-td>{{item.orgSName}}</n-td>
        <n-td>{{item.publishDate.substring(0,10)}}</n-td>
      </n-tr>
    </n-tbody>
</n-table>
</template>

<style scoped>

</style>