<script setup>

import {CaretDown, CaretUp, RefreshCircleOutline} from "@vicons/ionicons5";
import {NText,useMessage} from "naive-ui";
import {onBeforeUnmount, onMounted, onUnmounted, ref} from "vue";
import {useI18n} from 'vue-i18n'
import {GetMoneyRankSina} from "../../wailsjs/go/main/App";
import KLineChart from "./KLineChart.vue";

const { t } = useI18n()
const props = defineProps({
  headerTitle: {
    type: String,
    default: 'netInflowRanking'
  },
  sort: {
    type: String,
    default: 'netamount'
  },
})
const message = useMessage()
const dataList= ref([])
const sort = ref(props.sort)
const interval = ref(null)
onMounted(()=>{
  sort.value=props.sort
  GetMoneyRankSinaData()
  interval.value=setInterval(()=>{
    GetMoneyRankSinaData()
  },1000*60)
})
onBeforeUnmount(()=>{
  clearInterval(interval.value)
})
function GetMoneyRankSinaData(){
  message.loading(t('market.refreshingData'))
  GetMoneyRankSina(sort.value).then(result => {
    if(result.length>0){
      dataList.value = result
    }
  })
}
</script>

<template>
  <n-table striped size="small">
    <n-thead>
      <n-tr>
        <n-th>{{ t('market.code') }}</n-th>
        <n-th>{{ t('market.name') }}</n-th>
        <n-th>{{ t('market.latestPrice') }}</n-th>
        <n-th>{{ t('market.changeRate') }}</n-th>
        <n-th>{{ t('market.turnoverRate') }}</n-th>
        <n-th>{{ t('market.billboardDeal') }}</n-th>
        <n-th>{{ t('market.fundOutflow') }}</n-th>
        <n-th>{{ t('market.fundInflow') }}</n-th>
        <n-th>{{ t('market.netInflow') }}</n-th>
        <n-th>{{ t('market.netInflowRate') }}</n-th>
        <n-th v-if="sort === 'r0_net'||sort==='r0_out'">{{ t('market.mainForceOutflow') }}</n-th>
        <n-th v-if="sort === 'r0_net'">{{ t('market.mainForceNetInflow') }}</n-th>
        <n-th v-if="sort === 'r0_net'">{{ t('market.mainForceNetInflow') }}</n-th>
        <n-th >{{ t('market.mainForceNetInflowRate') }}</n-th>
        <n-th v-if="sort === 'r3_net'||sort==='r3_out'">{{ t('market.retailOutflow') }}</n-th>
        <n-th v-if="sort === 'r3_net'">{{ t('market.retailNetInflow') }}</n-th>
        <n-th v-if="sort === 'r3_net'">{{ t('market.retailNetInflow') }}</n-th>
        <n-th >{{ t('market.retailNetInflowRate') }}</n-th>
      </n-tr>
    </n-thead>
    <n-tbody>
      <n-tr v-for="item in dataList" :key="item.symbol">
        <n-td><n-tag :bordered=false type="info">{{ item.symbol }}</n-tag></n-td>
        <n-td>
          <n-popover trigger="hover" placement="right">
            <template #trigger>
              <n-button tag="a"  text :type="item.changeratio>0?'error':'success'" :bordered=false >{{ item.name }}</n-button>
            </template>
            <k-line-chart style="width: 800px" :code="item.symbol" :chart-height="500" :stockName="item.name" :k-days="20" :dark-theme="true"></k-line-chart>
          </n-popover>
        </n-td>
        <n-td><n-text :type="item.changeratio>0?'error':'success'">{{item.trade}}</n-text></n-td>
        <n-td><n-text :type="item.changeratio>0?'error':'success'">{{(item.changeratio*100).toFixed(2)}}%</n-text></n-td>
        <n-td><n-text :type="item.turnover>500?'error':'info'">{{(item.turnover/100).toFixed(2)}}%</n-text></n-td>
        <n-td><n-text type="info">{{(item.amount/10000).toFixed(2)}}</n-text></n-td>
        <n-td><n-text type="info"> {{(item.outamount/10000).toFixed(2)}}</n-text></n-td>
        <n-td><n-text type="info"> {{(item.inamount/10000).toFixed(2)}}</n-text></n-td>
        <n-td><n-text type="info"> {{(item.netamount/10000).toFixed(2)}}</n-text></n-td>
        <n-td><n-text :type="item.ratioamount>0?'error':'success'"> {{(item.ratioamount*100).toFixed(2)}}%</n-text></n-td>
        <n-td v-if="sort === 'r0_net'||sort==='r0_out'"><n-text  type="success"> {{(item.r0_out/10000).toFixed(2)}}</n-text></n-td>
        <n-td v-if="sort === 'r0_net'"><n-text  type="error"> {{(item.r0_in/10000).toFixed(2)}}</n-text></n-td>
        <n-td v-if="sort === 'r0_net'"><n-text :type="item.r0_net>0?'error':'success'"> {{(item.r0_net/10000).toFixed(2)}}</n-text></n-td>
        <n-td ><n-text :type="item.r0_ratio>0?'error':'success'"> {{(item.r0_ratio*100).toFixed(2)}}%</n-text></n-td>
        <n-td v-if="sort === 'r3_net'||sort==='r3_out'"><n-text  type="success"> {{(item.r3_out/10000).toFixed(2)}}</n-text></n-td>
        <n-td v-if="sort === 'r3_net'"><n-text  type="error"> {{(item.r3_in/10000).toFixed(2)}}</n-text></n-td>
        <n-td v-if="sort === 'r3_net'"><n-text :type="item.r3_net>0?'error':'success'"> {{(item.r3_net/10000).toFixed(2)}}</n-text></n-td>
        <n-td ><n-text :type="item.r3_ratio>0?'error':'success'"> {{(item.r3_ratio*100).toFixed(2)}}%</n-text></n-td>
      </n-tr>
    </n-tbody>
  </n-table>
</template>

<style scoped>

</style>