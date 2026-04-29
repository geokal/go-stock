<script setup>
import {onBeforeMount, ref} from 'vue'
import {GetStockList, IndustryResearchReport,EMDictCode} from "../../wailsjs/go/main/App";
import {ArrowDownOutline, CaretDown, CaretUp, PulseOutline, Refresh, RefreshCircleSharp,} from "@vicons/ionicons5";

import {useMessage} from "naive-ui";
import {BrowserOpenURL} from "../../wailsjs/runtime";
import {useI18n} from 'vue-i18n'

const { t } = useI18n()

const message=useMessage()
const list  = ref([])

const options =  ref([])

function getIndustryResearchReport(value) {
  message.loading(t('stockResearchReportList.loadingData'))
  IndustryResearchReport(value).then(result => {
    console.log(result)
    list.value = result
  })
}

onBeforeMount(()=>{
  getIndustryResearchReport('');
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
function openWin(code) {
  BrowserOpenURL("https://pdf.dfcfw.com/pdf/H3_"+code+"_1.pdf?1749744888000.pdf")
}

function EMDictCodeList(keyVal){
  if (keyVal){
    EMDictCode('016').then(result => {
      console.log(result)
        options.value=result.filter((value,index,array) => value.bkName.includes(keyVal)||value.firstLetter.includes(keyVal)||value.bkCode.includes(keyVal)).map(item => {
          return {
            label: item.bkName+" - "+item.bkCode,
            value: item.bkCode
          }
        })
    })
  }else{
    getIndustryResearchReport('')
  }

}
function handleSearch(value) {
  getIndustryResearchReport(value)
}
</script>

<template>
  <n-card>
    <n-auto-complete  :options="options" :placeholder="t('stockResearchReportList.industryKeyword')"  clearable filterable  :on-select="handleSearch"   :on-update:value="EMDictCodeList" />
  </n-card>
  <n-table striped size="small">
    <n-thead>
      <n-tr>
<!--        <n-th>代码</n-th>-->
<!--        <n-th>名称</n-th>-->
        <n-th>{{ t('stock.industry') }}</n-th>
        <n-th>{{ t('stockResearchReportList.title') }}</n-th>
        <n-th>{{ t('stockResearchReportList.emRating') }}</n-th>
        <n-th>{{ t('stockResearchReportList.ratingChange') }}</n-th>
        <n-th>{{ t('stockResearchReportList.orgRating') }}</n-th>
        <n-th>{{ t('stockResearchReportList.analyst') }}</n-th>
        <n-th>{{ t('stockResearchReportList.organization') }}</n-th>
        <n-th> <n-flex justify="space-between">{{ t('stockResearchReportList.date') }}<n-icon @click="getIndustryResearchReport" color="#409EFF" :size="20"  :component="RefreshCircleSharp"/></n-flex></n-th>
      </n-tr>
    </n-thead>
    <n-tbody>
      <n-tr v-for="item in list" :key="item.infoCode">
<!--        <n-td>{{item.stockCode}}</n-td>-->
<!--        <n-td :title="item.stockCode">-->
<!--          <n-popover trigger="hover" placement="right">-->
<!--            <template #trigger>-->
<!--              <n-tag type="info"  :bordered="false">{{item.stockName}}</n-tag>-->
<!--            </template>-->
<!--            <k-line-chart style="width: 800px" :code="getmMarketCode(item.market,item.stockCode)" :chart-height="500" :name="item.stockName" :k-days="20" :dark-theme="true"></k-line-chart>-->
<!--          </n-popover>-->
<!--        </n-td>-->
        <n-td><n-tag type="info"  :bordered="false">{{item.industryName}}</n-tag></n-td>
        <n-td>
          <n-a type="info"  @click="openWin(item.infoCode)"><n-text type="success">{{item.title}}</n-text></n-a>
        </n-td>
        <n-td><n-text :type="item.emRatingName==='增持'?'error':'info'">
          {{item.emRatingName}}
        </n-text></n-td>
        <n-td><n-text :type="item.ratingChange===0?'error':'info'">{{ratingChangeName(item.ratingChange)}}</n-text></n-td>
        <n-td>{{item.sRatingName	}}</n-td>
        <n-td><n-ellipsis style="max-width: 120px">{{item.researcher}}</n-ellipsis></n-td>
        <n-td>{{item.orgSName}}</n-td>
        <n-td>{{item.publishDate.substring(0,10)}}</n-td>
      </n-tr>
    </n-tbody>
</n-table>
</template>

<style scoped>

</style>