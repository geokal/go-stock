<script setup>
import {h, onBeforeMount, onMounted, ref, reactive} from 'vue'
import {useI18n} from 'vue-i18n'
const { t } = useI18n()
import {
  GetAllStockInfoList,
  GetAllStocks,
  GetConfig, GetSponsorInfo
} from "../../wailsjs/go/main/App";
import {NButton, NInput, NTag, NText, useMessage, useNotification, NDataTable, NSpace, NPagination} from "naive-ui";
import sparkLine from "./stockSparkLine.vue"
import klineChart from "./KLineChart.vue"
import KLineChart from "./KLineChart.vue";
import {format} from "date-fns";

const notify = useNotification()
const message = useMessage()

const editorDataRef = reactive({
  darkTheme: false
})

onBeforeMount(() => {
  GetConfig().then(result => {
    if (result.darkTheme) {
      editorDataRef.darkTheme = true
    }
  })

  GetSponsorInfo().then((res) => {
    // console.log(res)
    vipLevel.value = res.vipLevel;
    vipStartTime.value = res.vipStartTime;
    vipEndTime.value = res.vipEndTime;
    //判断时间是否到期
    if (res.vipLevel) {
      if (res.vipEndTime < format(new Date(), 'yyyy-MM-dd HH:mm:ss')) {
        //notify.warning({content: 'VIP已到期'})
        expired.value = true;
      }
    }else{
      //notify.success({content: '未开通VIP'})
    }
    isValidVip.value = !(vipLevel.value === "" || Number(vipLevel.value) <= 0);

  })

})

onMounted(() => {
  console.log('stock-list mounted')
  loadStocks(1, paginationReactive.pageSize)
})

const dataRef = ref([])
const loadingRef = ref(false)
const vipLevel=ref("");
const vipStartTime=ref("");
const vipEndTime=ref("");
const expired=ref(false)
const isValidVip=ref(false) // 是否是会员
const columnsRef = ref([
  // {
  //   title: '数据时间',
  //   key: 'MAX_TRADE_DATE',
  //   width: 120,
  // },
  {
    title: t('stock.code'),
    key: 'SECUCODE',
    width: 100,
    render(row) {
      return h(NText, { type: "info" }, { default: () => row.SECUCODE })
    }
  },
  {
    title: t('stock.name'),
    key: 'SECURITY_NAME_ABBR',
    width: 100,
    render(row) {
      return h(NText, { type: "success" }, { default: () => row.SECURITY_NAME_ABBR })
    }
  },
  {
    title: t('stock.price'),
    key: 'NEW_PRICE',
    width: 100,
    render(row) {
      const price = row.NEW_PRICE
      return h(NText, { type: "info" }, { default: () => isNumeric(price) ? price : '-' })
    }
  },
  {
    title: t('stock.changeRate'),
    key: 'CHANGE_RATE',
    width: 100,
    render(row) {
      const rate = toNumber(row.CHANGE_RATE, 0)
      const type = rate >= 0 ? 'error' : 'success'
      const sign = rate >= 0 ? '+' : ''
      return h(NText, { type: type }, { default: () => `${sign}${rate.toFixed(2)}%` })
    }
  },
  {
    title: t('stock.sparkline'),
    key: 'sparkline',
    width: 120,
    render(row) {
      return h(sparkLine, {
        idSuffix: row.SECUCODE,
        stockName: row.SECURITY_NAME_ABBR,
        stockCode: row.SECUCODE,
        lastPrice: row.NEW_PRICE,
        openPrice: row.PRE_CLOSE_PRICE,
        tooltip: true
      })
    }
  },
  {
    title: t('stock.highPrice'),
    key: 'HIGH_PRICE',
    width: 100,
    render(row) {
      const price = row.HIGH_PRICE
      return h(NText, { type: "info" }, { default: () => isNumeric(price) ? price : '-' })
    }
  },
  {
    title: t('stock.lowPrice'),
    key: 'LOW_PRICE',
    width: 100,
    render(row) {
      const price = row.LOW_PRICE
      return h(NText, { type: "info" }, { default: () => isNumeric(price) ? price : '-' })
    }
  },
  // {
  //   title: '前收价',
  //   key: 'PRE_CLOSE_PRICE',
  //   width: 100,
  //   render(row) {
  //     return h(NText, { type: "info" }, { default: () => row.PRE_CLOSE_PRICE.toFixed(2) })
  //   }
  // },
  {
    title: t('stock.volume'),
    key: 'VOLUME',
    width: 120,
    render(row) {
      const volume = toNumber(row.VOLUME, 0)
      let displayVolume = volume
      if (volume >= 100000000) {
        displayVolume = (volume / 100000000).toFixed(2) + '亿'
      } else if (volume >= 10000) {
        displayVolume = (volume / 10000).toFixed(2) + '万'
      }
      return h(NText, { type: "info" }, { default: () => displayVolume })
    }
  },
  {
    title: t('stock.turnover'),
    key: 'DEAL_AMOUNT',
    width: 120,
    render(row) {
      const amount = toNumber(row.DEAL_AMOUNT, 0)
      let displayAmount = amount
      if (amount >= 100000000) {
        displayAmount = (amount / 100000000).toFixed(2) + '亿'
      } else if (amount >= 10000) {
        displayAmount = (amount / 10000).toFixed(2) + '万'
      }
      return h(NText, { type: "info" }, { default: () => displayAmount })
    }
  },
  {
    title: t('stock.turnoverRate'),
    key: 'TURNOVERRATE',
    width: 80,
    render(row) {
      const rate = row.TURNOVERRATE
      return h(NText, { type: "info" }, { default: () => isNumeric(rate) ? rate : '-' })
    }
  },
  {
    title: t('stock.volumeRatio'),
    key: 'VOLUME_RATIO',
    width: 80,
    render(row) {
      const ratio = row.VOLUME_RATIO
      return h(NText, { type: "info" }, { default: () => isNumeric(ratio) ? ratio : '-' })
    }
  },
  {
    title: t('stock.industry'),
    key: 'INDUSTRY',
    width: 100,
    render(row) {
      return h(NTag, { type: "primary", size: "small" }, { default: () => row.INDUSTRY })
    }
  },
  {
    title: t('stock.concept'),
    key: 'CONCEPT',
    width: 100,
    ellipsis: {
      tooltip: true
    },
    render(row) {
      if(typeof row.CONCEPT === 'string'){
        return h(NTag, { type: "info", size: "small" ,style: "margin-right: 4px;" }, { default: () => row.CONCEPT })
      }else{
        if (!row.CONCEPT || row.CONCEPT.length === 0) {
          return h(NText, { type: "secondary" }, { default: () => '-' })
        }
        return row.CONCEPT.map(concept =>
            h(NTag, { type: "info", size: "small", style: "margin-right: 4px;" }, { default: () => concept })
        )
      }
    }
  },
  // {
  //   title: '交易所',
  //   key: 'MARKET',
  //   width: 100,
  //   render(row) {
  //     return h(NTag, { type: "warning", size: "small" }, { default: () => row.MARKET })
  //   }
  // },
  {
    title: t('common.edit'),
    render(row, index) {
      return [h(
          NButton,
          {
            secondary: true,
            size: 'small',
            type: 'warning', // 橙色按钮
            onClick: () => showKline(row)
          },
          { default: () => t('allStockList.dayKLine') }
      ),]
    }
  },
])

const paginationReactive = reactive({
  keyword:"",
  page: 1,
  pageCount: 1,
  pageSize: 9,
  itemCount: 0,
  prefix({ itemCount }) {
    return t('stock.stockCount', { count: itemCount })
  }
})
const optionsReactive= reactive([
  {
    label: t('common.all'),
    value: ''
  },
 ])

function loadStocks(page, pageSize) {
  if((vipLevel.value===""|| Number(vipLevel.value) <=0)){
    handleReset()
  }
  if (!loadingRef.value) {
    loadingRef.value = true
    GetAllStocks(page, pageSize, paginationReactive.keyword, technicalIndicatorReactive).then((res) => {
      console.log(res)
      if (res && res.result && res.result.data) {
        dataRef.value = res.result.data
        paginationReactive.page = page
        paginationReactive.pageCount = Math.ceil(res.result.count / pageSize)
        paginationReactive.itemCount = res.result.count
      } else {
        dataRef.value = []
        paginationReactive.page = 1
        paginationReactive.pageCount = 1
        paginationReactive.itemCount = 0
        message.error(t('allStockList.getStockDataFailed'))
      }
      loadingRef.value = false
    }).catch(err => {
      message.error(t('allStockList.getStockDataFailed') + ': ' + err.message)
      loadingRef.value = false
    })
  }
}
function handleCheckedChange(checked) {

  if(checked&&(vipLevel.value===""|| Number(vipLevel.value) <=0)){
    handleReset()
    message.warning(t('allStockList.vipExpiredNotice'))
  }
}
function handlePageChange(currentPage) {
  loadStocks(currentPage, paginationReactive.pageSize)
}

function handlePageSizeChange(pageSize) {
  paginationReactive.pageSize = pageSize
  loadStocks(1, pageSize)
}
function handleSearch() {
  loadStocks(1, paginationReactive.pageSize)
}
function handleUpdateVal(value) {
  console.log('handleUpdateVal', value)
  if (value === '') {
    optionsReactive.splice(1, optionsReactive.length - 1)
  } else {
    GetAllStockInfoList({
      searchKeyWord: value
    }).then((res) => {
      console.log('GetAllStockInfoList result:', res)
      if (res  && res.list) {
        optionsReactive.splice(1, optionsReactive.length - 1)
        optionsReactive.push(...res.list.map(item => {
          return {
            label: item.SECURITY_NAME_ABBR,
            value: item.SECURITY_NAME_ABBR,
            obj: item,
          }
        }))
      }
    }).catch(err => {
      message.error(t('allStockList.getStockDataFailed') + ': ' + err.message)
    })
  }
}
const modalDataRef = reactive({
  visible: false,
  title: "",
  content: "",
  riskRemarks: "",
  stockCode: "",
  stockName: "",
  remarks: "",
})
function showKline(row) {
  console.log('showKline', row)
  modalDataRef.title = row.SECURITY_NAME_ABBR
  modalDataRef.stockCode = getStockCode(row.SECUCODE)
  modalDataRef.stockName = row.SECURITY_NAME_ABBR
  modalDataRef.visible = true
}
function getStockCode(stockCode) {
  if(stockCode.indexOf( ".")>0){
    stockCode=stockCode.split(".")[1]+stockCode.split(".")[0]
  }
  //转化为小写
  stockCode=stockCode.toLowerCase()
  return stockCode

}
const technicalIndicatorReactive = reactive({
  MACD_GOLDEN_FORK: false,
  KDJ_GOLDEN_FORK: false,
  BREAK_THROUGH: false,
  LOW_FUNDS_INFLOW: false,
  HIGH_FUNDS_OUTFLOW: false,
  BREAKUP_MA_5DAYS: false,
  LONG_AVG_ARRAY: false,
  SHORT_AVG_ARRAY: false,
  UPPER_LARGE_VOLUME: false,
  DOWN_NARROW_VOLUME: false,
  ONE_DAYANG_LINE: false,
  TWO_DAYANG_LINES: false,
  RISE_SUN: false,
  POWER_FULGUN: false,
  RESTORE_JUSTICE: false,
  DOWN_7DAYS: false,
  UPPER_8DAYS:false,
  UPPER_9DAYS:false,
  UPPER_4DAYS:false,
  HEAVEN_RULE:false,
  UPSIDE_VOLUME: false,
  BEARISH_ENGULFING: false,
  REVERSING_HAMMER: false,
  SHOOTING_STAR: false,
  EVENING_STAR: false,
  FIRST_DAWN: false,
  PREGNANT: false,
  BLACK_CLOUD_TOPS: false,
  MORNING_STAR: false,
  NARROW_FINISH: false,
  UPP_DAYS:0,
  CONCERN_RANK_7DAYS:0,
  UPNDAY:0,
  DOWNNDAY :0,
})

function handleReset(){
  technicalIndicatorReactive.MACD_GOLDEN_FORK = false
  technicalIndicatorReactive.KDJ_GOLDEN_FORK = false
  technicalIndicatorReactive.BREAK_THROUGH = false
  technicalIndicatorReactive.LOW_FUNDS_INFLOW = false
  technicalIndicatorReactive.HIGH_FUNDS_OUTFLOW = false
  technicalIndicatorReactive.BREAKUP_MA_5DAYS = false
  technicalIndicatorReactive.LONG_AVG_ARRAY = false
  technicalIndicatorReactive.SHORT_AVG_ARRAY = false
  technicalIndicatorReactive.UPPER_LARGE_VOLUME = false
  technicalIndicatorReactive.DOWN_NARROW_VOLUME = false
  technicalIndicatorReactive.ONE_DAYANG_LINE = false
  technicalIndicatorReactive.TWO_DAYANG_LINES = false
  technicalIndicatorReactive.RISE_SUN = false
  technicalIndicatorReactive.POWER_FULGUN = false
  technicalIndicatorReactive.RESTORE_JUSTICE = false
  technicalIndicatorReactive.DOWN_7DAYS = false
  technicalIndicatorReactive.UPPER_8DAYS=false
  technicalIndicatorReactive.UPPER_9DAYS=false
  technicalIndicatorReactive.UPPER_4DAYS=false
  technicalIndicatorReactive.HEAVEN_RULE=false
  technicalIndicatorReactive.ONE_DAYANG_LINE=false
  technicalIndicatorReactive.TWO_DAYANG_LINES= false
  technicalIndicatorReactive.RISE_SUN=false
  technicalIndicatorReactive.POWER_FULGUN=false
  technicalIndicatorReactive.RESTORE_JUSTICE=false
  technicalIndicatorReactive.DOWN_7DAYS=false
  technicalIndicatorReactive.UPPER_8DAYS=false
  technicalIndicatorReactive.UPPER_9DAYS=false
  technicalIndicatorReactive.UPPER_4DAYS=false
  technicalIndicatorReactive.HEAVEN_RULE=false
  technicalIndicatorReactive.UPSIDE_VOLUME=false
  technicalIndicatorReactive.BEARISH_ENGULFING=false
  technicalIndicatorReactive.REVERSING_HAMMER=false
  technicalIndicatorReactive.SHOOTING_STAR=false
  technicalIndicatorReactive.EVENING_STAR=false
  technicalIndicatorReactive.FIRST_DAWN=false
  technicalIndicatorReactive.PREGNANT=false
  technicalIndicatorReactive.BLACK_CLOUD_TOPS=false
  technicalIndicatorReactive.MORNING_STAR=false
  technicalIndicatorReactive.NARROW_FINISH=false
  technicalIndicatorReactive.UPP_DAYS=0
  technicalIndicatorReactive.CONCERN_RANK_7DAYS=0
  technicalIndicatorReactive.UPNDAY=0
  technicalIndicatorReactive.DOWNNDAY=0
}

// 判断是否是数字
const isNumeric = (value) => {
  if (value === null || value === undefined || value === '') {
    return false
  }
  return !isNaN(Number(value))
}

// 安全转换数字
const toNumber = (value, defaultValue = 0) => {
  const num = Number(value)
  return isNaN(num) ? defaultValue : num
}

</script>

<template>
    <n-space justify="start">
      <n-card size="small" :bordered="false"  style="text-align: left">
        <n-checkbox   @update:checked="handleCheckedChange" v-model:checked="technicalIndicatorReactive.MACD_GOLDEN_FORK">
        {{ t('allStockList.macdGoldenFork') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.KDJ_GOLDEN_FORK">
        {{ t('allStockList.kdjGoldenFork') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.BREAK_THROUGH">
        {{ t('allStockList.volumeBreakout') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.LOW_FUNDS_INFLOW">
        {{ t('allStockList.lowFundsInflow') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.HIGH_FUNDS_OUTFLOW">
        {{ t('allStockList.highFundsOutflow') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.BREAKUP_MA_5DAYS">
        {{ t('allStockList.breakupMA5Days') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.LONG_AVG_ARRAY">
        {{ t('allStockList.longAvgArray') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.SHORT_AVG_ARRAY">
        {{ t('allStockList.shortAvgArray') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.UPPER_LARGE_VOLUME">
        {{ t('allStockList.upperLargeVolume') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.DOWN_NARROW_VOLUME">
        {{ t('allStockList.downNarrowVolume') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.ONE_DAYANG_LINE">
        {{ t('allStockList.oneDayangLine') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.TWO_DAYANG_LINES">
        {{ t('allStockList.twoDayangLines') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"     v-model:checked="technicalIndicatorReactive.RISE_SUN">
        {{ t('allStockList.riseSun') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.POWER_FULGUN">
        {{ t('allStockList.powerfulgun') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.RESTORE_JUSTICE">
        {{ t('allStockList.restoreJustice') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"     v-model:checked="technicalIndicatorReactive.DOWN_7DAYS">
        {{ t('allStockList.down7Days') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.UPPER_8DAYS">
        {{ t('allStockList.upper8Days') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.UPPER_9DAYS">
        {{ t('allStockList.upper9Days') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.UPPER_4DAYS">
        {{ t('allStockList.upper4Days') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.HEAVEN_RULE">
        {{ t('allStockList.heavenRule') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.UPSIDE_VOLUME">
        {{ t('allStockList.upsideVolume') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.BEARISH_ENGULFING">
        {{ t('allStockList.bearishEngulfing') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.REVERSING_HAMMER">
        {{ t('allStockList.reversingHammer') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.SHOOTING_STAR">
        {{ t('allStockList.shootingStar') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.EVENING_STAR">
        {{ t('allStockList.eveningStar') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.FIRST_DAWN">
        {{ t('allStockList.firstDawn') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.PREGNANT">
        {{ t('allStockList.pregnant') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.BLACK_CLOUD_TOPS">
        {{ t('allStockList.blackCloudTops') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.MORNING_STAR">
        {{ t('allStockList.morningStar') }}
      </n-checkbox>
      <n-checkbox  @update:checked="handleCheckedChange"    v-model:checked="technicalIndicatorReactive.NARROW_FINISH">
        {{ t('allStockList.narrowFinish') }}
      </n-checkbox>
      </n-card>
      <n-card size="small" :bordered="false"  style="text-align: left">
        <n-radio-group size="small"  @update:checked="handleCheckedChange" name="UPP_DAYS"   v-model:value="technicalIndicatorReactive.UPP_DAYS">
          <n-radio :value="3">{{ t('allStockList.popularityRankRise3Days') }}</n-radio>
          <n-radio :value="5">{{ t('allStockList.popularityRankRise5Days') }}</n-radio>
          <n-radio :value="7">{{ t('allStockList.popularityRankRise7Days') }}</n-radio>
        </n-radio-group>
        <n-divider vertical/>
        <n-radio-group  size="small" @update:checked="handleCheckedChange"  name="CONCERN_RANK_7DAYS"  v-model:value="technicalIndicatorReactive.CONCERN_RANK_7DAYS">
          <n-radio :value="10">{{ t('allStockList.attentionRank7DaysTop10') }}</n-radio>
          <n-radio :value="50">{{ t('allStockList.attentionRank7DaysTop50') }}</n-radio>
          <n-radio :value="100">{{ t('allStockList.attentionRank7DaysTop100') }}</n-radio>
        </n-radio-group>

        <n-radio-group  size="small" @update:checked="handleCheckedChange" name="UPNDAY"  v-model:value="technicalIndicatorReactive.UPNDAY">
          <n-radio :value="3">{{ t('allStockList.consecutiveRiseDays3Days') }}</n-radio>
          <n-radio :value="5">{{ t('allStockList.consecutiveRiseDays5Days') }}</n-radio>
          <n-radio :value="8">{{ t('allStockList.consecutiveRiseDays8Days') }}</n-radio>
        </n-radio-group>
        <n-divider vertical/>
        <n-radio-group  size="small" @update:checked="handleCheckedChange" name="DOWNNDAY"  v-model:value="technicalIndicatorReactive.DOWNNDAY">
          <n-radio :value="3">{{ t('allStockList.consecutiveDeclineDays3Days') }}</n-radio>
          <n-radio :value="5">{{ t('allStockList.consecutiveDeclineDays5Days') }}</n-radio>
          <n-radio :value="8">{{ t('allStockList.consecutiveDeclineDays8Days') }}</n-radio>
          <n-radio :value="10">{{ t('allStockList.consecutiveDeclineDays10Days') }}</n-radio>
          <n-radio :value="14">{{ t('allStockList.consecutiveDeclineDays14Days') }}</n-radio>
        </n-radio-group>
      </n-card>

    </n-space>
    <n-input-group>
<!--    <n-input clearable placeholder="输入股票名称" v-model:value="paginationReactive.keyword"/>-->
      <n-auto-complete
          v-model:value="paginationReactive.keyword"
          :input-props="{
            autocomplete: 'disabled',
          }"
          :options="optionsReactive"
          :placeholder="t('allStockList.enterKeyword')"
          clearable
          @input="handleUpdateVal"
          @select="(value) => {
            paginationReactive.keyword = value
            handleSearch()
          }"
      />
    <n-button type="primary" ghost @click="handleSearch"  @input="handleSearch">
      {{ t('common.search') }}
    </n-button>
      <n-button @click="handleReset">{{ t('common.reset') }}</n-button>

    </n-input-group>
    <!-- 数据表格 -->
    <n-data-table
      remote
      size="small"
      :columns="columnsRef"
      :data="dataRef"
      :loading="loadingRef"
      :pagination="paginationReactive"
      :row-key="(rowData) => rowData.SECUCODE"
      flex-height
      style="height: calc(100vh - 380px);margin-top: 10px"
      @update:page="handlePageChange"
    />
    
    <!-- 分页控件 -->
<!--    <div style="margin-top: 16px; display: flex; justify-content: center;">-->
<!--      <n-pagination-->
<!--        v-model:page="paginationReactive.page"-->
<!--        v-model:page-size="paginationReactive.pageSize"-->
<!--        :page-count="paginationReactive.pageCount"-->
<!--        :item-count="paginationReactive.itemCount"-->
<!--        :page-sizes="[10, 20, 50, 100]"-->
<!--        show-size-picker-->
<!--        show-quick-jumper-->
<!--        @update:page="handlePageChange"-->
<!--        @update:page-size="handlePageSizeChange"-->
<!--      />-->
<!--    </div>-->

  <n-modal v-model:show="modalDataRef.visible" :title="modalDataRef.title" preset="card" style="width: 850px;">
    <n-card size="small">
      <KLineChart style="width: 800px" :code="getStockCode(modalDataRef.stockCode)" :chart-height="500" :stock-name="modalDataRef.stockName" :k-days="30" :dark-theme="editorDataRef.darkTheme"></KLineChart>
    </n-card>
  </n-modal>
</template>

<style scoped>
</style>
