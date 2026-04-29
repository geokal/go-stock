<script setup>
import {computed, h, onBeforeMount, onBeforeUnmount, onMounted, onUnmounted, ref, reactive} from 'vue'
import {GetStockChanges, GetConfig, GetStockChangeHistory, SaveStockChangesToHistory, GetAllStockChangesWithPaging} from "../../wailsjs/go/main/App";
import {NButton, NTag, NText, useMessage, useNotification} from "naive-ui";
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const notify = useNotification()
const message = useMessage()
const loadingRef = ref(true)
const dataRef = ref([])
const autoRefresh = ref(true)
const refreshInterval = ref(null)
const refreshSeconds = ref(10)
const countdown = ref(10)

const viewMode = ref('realtime')
const isTrading = ref(false)
const marketStatus = ref('')

const paginationReactive = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 50,
  itemCount: 0,
  keyword: "",
  range: null,
  startTime: null,
  endTime: null,
  minVolume: null,
  minAmount: null,
  minChangeRate: null,
  maxChangeRate: null,
  industry: "",
  concept: "",
  prefix({ itemCount }) {
    return `${itemCount} ${t('common.records')}`
  }
})

const volumeOptions = [
  { label: t('stockChangesMonitor.unlimited'), value: null },
  { label: t('stockChangesMonitor.over100lots'), value: 10000 },
  { label: t('stockChangesMonitor.over200lots'), value: 20000 },
  { label: t('stockChangesMonitor.over500lots'), value: 50000 },
  { label: t('stockChangesMonitor.over1000lots'), value: 100000 },
  { label: t('stockChangesMonitor.over2000lots'), value: 200000 },
  { label: t('stockChangesMonitor.over5000lots'), value: 500000 },
  { label: t('stockChangesMonitor.over10000lots'), value: 1000000 },
  { label: t('stockChangesMonitor.over20000lots'), value: 2000000 },
  { label: t('stockChangesMonitor.over50000lots'), value: 5000000 },
]

const amountOptions = [
  { label: t('stockChangesMonitor.unlimited'), value: null },
  { label: t('stockChangesMonitor.over100w'), value: 1000000 },
  { label: t('stockChangesMonitor.over500w'), value: 5000000 },
  { label: t('stockChangesMonitor.over1000w'), value: 10000000 },
  { label: t('stockChangesMonitor.over2000w'), value: 20000000 },
  { label: t('stockChangesMonitor.over5000w'), value: 50000000 },
  { label: t('stockChangesMonitor.over100m'), value: 100000000 },
  { label: t('stockChangesMonitor.over200m'), value: 200000000 },
  { label: t('stockChangesMonitor.over500m'), value: 500000000 },
]

const changeRateOptions = [
  { label: t('stockChangesMonitor.unlimited'), value: null },
  { label: t('stockChangesMonitor.over3pct'), value: 3 },
  { label: t('stockChangesMonitor.over5pct'), value: 5 },
  { label: t('stockChangesMonitor.over7pct'), value: 7 },
  { label: t('stockChangesMonitor.over9pct'), value: 9 },
  { label: t('stockChangesMonitor.overLimitUp'), value: 9.9 },
]

const bullishTypes = [
  {label: t('stockChangesMonitor.rocketLaunch'), value: '8201'},
  {label: t('stockChangesMonitor.quickRebound'), value: '8202'},
  {label: t('stockChangesMonitor.largeBuy'), value: '8193'},
  {label: t('stockChangesMonitor.sealLimitUp'), value: '4'},
  {label: t('stockChangesMonitor.openLimitDown'), value: '32'},
  {label: t('stockChangesMonitor.hasBigBuy'), value: '64'},
  {label: t('stockChangesMonitor.auctionRise'), value: '8207'},
  {label: t('stockChangesMonitor.aboveMA5'), value: '8209'},
  {label: t('stockChangesMonitor.upGap'), value: '8211'},
  {label: t('stockChangesMonitor.newHigh60D'), value: '8213'},
  {label: t('stockChangesMonitor.bigRise60D'), value: '8215'},
  {label: t('stockChangesMonitor.openLimitUp'), value: '16'},
]

const bearishTypes = [
  {label: t('stockChangesMonitor.acceleratedFall'), value: '8204'},
  {label: t('stockChangesMonitor.highDive'), value: '8203'},
  {label: t('stockChangesMonitor.largeSell'), value: '8194'},
  {label: t('stockChangesMonitor.sealLimitDown'), value: '8'},
  {label: t('stockChangesMonitor.hasBigSell'), value: '128'},
  {label: t('stockChangesMonitor.auctionFall'), value: '8208'},
  {label: t('stockChangesMonitor.belowMA5'), value: '8210'},
  {label: t('stockChangesMonitor.downGap'), value: '8212'},
  {label: t('stockChangesMonitor.newLow60D'), value: '8214'},
  {label: t('stockChangesMonitor.bigFall60D'), value: '8216'},
]

const allTypeValues = [...bullishTypes, ...bearishTypes].map(t => t.value)
const selectedTypes = ref(allTypeValues)

const columnsRef = ref([
  {
    title: t('stockChangesMonitor.date'),
    key: 'changeDate',
    width: 100,
    render(row) {
      const date = row.changeDate || row.ChangeDate
      const time = row.changeTime || row.ChangeTime || row.time
      if (date) {
        return h(NText, {type: 'info'}, {default: () => date + ' ' + time})
      }
      return h(NText, {type: 'info'}, {default: () => time})
    }
  },
  {
    title: t('stockChangesMonitor.code'),
    key: 'code',
    width: 100,
    render(row) {
      const code = row.stockCode || row.StockCode || row.code
      return h(NText, {type: 'info', style: 'cursor: pointer', onClick: () => copyCode(code)}, {default: () => code})
    }
  },
  {
    title: t('stockChangesMonitor.name'),
    key: 'name',
    width: 100,
    render(row) {
      return row.stockName || row.StockName || row.name
    }
  },
  {
    title: t('stockChangesMonitor.changeType'),
    key: 'typeName',
    width: 120,
    render(row) {
      const typeName = row.typeName || row.TypeName
      const bullishSet = new Set([t('stockChangesMonitor.rocketLaunch'), t('stockChangesMonitor.quickRebound'), t('stockChangesMonitor.largeBuy'), t('stockChangesMonitor.sealLimitUp'), t('stockChangesMonitor.openLimitDown'), t('stockChangesMonitor.hasBigBuy'), t('stockChangesMonitor.auctionRise'), t('stockChangesMonitor.aboveMA5'), t('stockChangesMonitor.upGap'), t('stockChangesMonitor.newHigh60D'), t('stockChangesMonitor.bigRise60D'), t('stockChangesMonitor.openLimitUp')])
      const bearishSet = new Set([t('stockChangesMonitor.acceleratedFall'), t('stockChangesMonitor.highDive'), t('stockChangesMonitor.largeSell'), t('stockChangesMonitor.sealLimitDown'), t('stockChangesMonitor.hasBigSell'), t('stockChangesMonitor.auctionFall'), t('stockChangesMonitor.belowMA5'), t('stockChangesMonitor.downGap'), t('stockChangesMonitor.newLow60D'), t('stockChangesMonitor.bigFall60D')])
      
      let tagType = 'default'
      if (bullishSet.has(typeName)) {
        tagType = 'error'
      } else if (bearishSet.has(typeName)) {
        tagType = 'success'
      }
      return h(NTag, {type: tagType, size: 'small'}, {default: () => typeName})
    }
  },
  {
    title: t('stockChangesMonitor.price'),
    key: 'price',
    width: 80,
    render(row) {
      const price = row.price || row.Price
      if (price > 0) {
        return price.toFixed(2)
      }
      return '-'
    }
  },
  {
    title: t('stockChangesMonitor.changeRatePct'),
    key: 'changeRate',
    width: 100,
    render(row) {
      const changeRate = row.changeRate || row.ChangeRate
      if (changeRate !== 0) {
        const color = changeRate > 0 ? '#dc2626' : '#16a34a'
        const prefix = changeRate > 0 ? '+' : ''
        return h('span', {style: {color: color, fontWeight: '500'}}, prefix + changeRate.toFixed(2) + '%')
      }
      return '-'
    }
  },
  {
    title: t('stockChangesMonitor.volume'),
    key: 'volume',
    width: 100,
    render(row) {
      const volume = row.volume || row.Volume
      if (volume > 0) {
        return formatVolume(volume)
      }
      return '-'
    }
  },
  {
    title: t('stockChangesMonitor.amount'),
    key: 'amount',
    width: 100,
    render(row) {
      const amount = row.amount || row.Amount
      if (amount > 0) {
        return formatAmount(amount)
      }
      return '-'
    }
  },
  {
    title: t('stockChangesMonitor.industry'),
    key: 'industry',
    width: 100,
    ellipsis: {
      tooltip: true
    },
    render(row) {
      return row.industry || row.Industry || '-'
    }
  },
  {
    title: t('stockChangesMonitor.concept'),
    key: 'concept',
    width: 150,
    ellipsis: {
      tooltip: true
    },
    render(row) {
      return row.concept || row.Concept || '-'
    }
  },
])

function formatVolume(vol) {
  const lots = vol / 100
  if (lots >= 100000000) {
    return (lots / 100000000).toFixed(2) + t('stockChangesMonitor.hundredMillionLots')
  } else if (lots >= 10000) {
    return (lots / 10000).toFixed(2) + t('stockChangesMonitor.tenThousandLots')
  } else if (lots >= 1) {
    return lots.toFixed(0) + t('stockChangesMonitor.lots')
  }
  return vol + t('stockChangesMonitor.shares')
}

function formatAmount(amount) {
  if (amount >= 100000000) {
    return (amount / 100000000).toFixed(2) + t('stockChangesMonitor.hundredMillion')
  } else if (amount >= 10000) {
    return (amount / 10000).toFixed(2) + t('stockChangesMonitor.tenThousand')
  }
  return amount.toFixed(2)
}

function copyCode(code) {
  navigator.clipboard.writeText(code).then(() => {
    message.success(t('stockChangesMonitor.copied') + code)
  })
}

function checkTradingTime() {
  const now = new Date()
  const day = now.getDay()
  const hour = now.getHours()
  const minute = now.getMinutes()
  const currentTime = hour * 100 + minute

  if (day === 0 || day === 6) {
    isTrading.value = false
    marketStatus.value = t('stockChangesMonitor.closedWeekend')
    return
  }

  const morningStart = 915
  const morningEnd = 1130
  const afternoonStart = 1257
  const afternoonEnd = 1500

  if (currentTime >= morningStart && currentTime <= morningEnd) {
    isTrading.value = true
    if (currentTime < 930) {
      marketStatus.value = t('stockChangesMonitor.auction')
    } else {
      marketStatus.value = t('stockChangesMonitor.morningTrading')
    }
  } else if (currentTime >= afternoonStart && currentTime <= afternoonEnd) {
    isTrading.value = true
    if (currentTime < 1300) {
      marketStatus.value = t('stockChangesMonitor.noonAuction')
    } else {
      marketStatus.value = t('stockChangesMonitor.afternoonTrading')
    }
  } else if (currentTime > morningEnd && currentTime < afternoonStart) {
    isTrading.value = false
    marketStatus.value = t('stockChangesMonitor.lunchBreak')
  } else if (currentTime > afternoonEnd) {
    isTrading.value = false
    marketStatus.value = t('stockChangesMonitor.closed')
  } else {
    isTrading.value = false
    marketStatus.value = t('stockChangesMonitor.notOpened')
  }
}

async function fetchRealtimeData() {
  loadingRef.value = true
  try {
    const types = selectedTypes.value.map(t => parseInt(t))
    const result = await GetStockChanges(types, 0, paginationReactive.pageSize)
    if (result) {
      dataRef.value = result.data || []
      paginationReactive.itemCount = result.totalCount || 0
    }
  } catch (e) {
    console.error(t('stockChangesMonitor.fetchFailed'), e)
  } finally {
    loadingRef.value = false
  }
}

async function fetchHistoryData() {
  loadingRef.value = true
  try {
    const query = {
      page: paginationReactive.page,
      pageSize: paginationReactive.pageSize,
    }
    if (paginationReactive.range && paginationReactive.range.length === 2) {
      query.startDate = formatDate(paginationReactive.range[0])
      query.endDate = formatDate(paginationReactive.range[1])
    }
    if (paginationReactive.startTime) {
      query.startTime = formatTime(paginationReactive.startTime)
    }
    if (paginationReactive.endTime) {
      query.endTime = formatTime(paginationReactive.endTime)
    }
    if (paginationReactive.minVolume) {
      query.minVolume = paginationReactive.minVolume
    }
    if (paginationReactive.minAmount) {
      query.minAmount = paginationReactive.minAmount
    }
    if (paginationReactive.minChangeRate) {
      query.minChangeRate = paginationReactive.minChangeRate
    }
    if (paginationReactive.maxChangeRate) {
      query.maxChangeRate = paginationReactive.maxChangeRate
    }
    if (paginationReactive.industry.trim()) {
      query.industry = paginationReactive.industry.trim()
    }
    if (paginationReactive.concept.trim()) {
      query.concept = paginationReactive.concept.trim()
    }
    if (paginationReactive.keyword.trim()) {
      const keyword = paginationReactive.keyword.trim()
      if (/^\d+$/.test(keyword)) {
        query.stockCode = keyword
      } else {
        query.stockName = keyword
      }
    }
    if (selectedTypes.value.length > 0) {
      query.changeTypes = selectedTypes.value.map(t => parseInt(t))
    }
    const result = await GetStockChangeHistory(query)
    if (result) {
      dataRef.value = result.list || []
      paginationReactive.itemCount = result.total || 0
      paginationReactive.pageCount = result.totalPages || 1
    }
  } catch (e) {
    console.error(t('stockChangesMonitor.fetchHistoryFailed'), e)
  } finally {
    loadingRef.value = false
  }
}

function formatDate(dateValue) {
  if (!dateValue) return ''
  const date = new Date(dateValue)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function formatTime(timeValue) {
  if (!timeValue) return ''
  const date = new Date(timeValue)
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${hours}:${minutes}:${seconds}`
}

async function fetchData() {
  if (viewMode.value === 'realtime') {
    await fetchRealtimeData()
  } else {
    await fetchHistoryData()
  }
}

async function saveCurrentData() {
  const types = selectedTypes.value.map(t => parseInt(t))
  const result = await SaveStockChangesToHistory(types)
  message.info(result)
}

function startAutoRefresh() {
  stopAutoRefresh()
  countdown.value = refreshSeconds.value
  refreshInterval.value = setInterval(() => {
    checkTradingTime()
    countdown.value--
    if (countdown.value <= 0) {
      if (viewMode.value === 'realtime' && !isTrading.value) {
        countdown.value = refreshSeconds.value
        return
      }
      fetchData()
      countdown.value = refreshSeconds.value
    }
  }, 1000)
}

function stopAutoRefresh() {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
    refreshInterval.value = null
  }
}

function toggleAutoRefresh() {
  if (autoRefresh.value) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

function selectAllBullish() {
  const bullishValues = bullishTypes.map(t => t.value)
  const currentBearish = selectedTypes.value.filter(t => bearishTypes.some(b => b.value === t))
  selectedTypes.value = [...bullishValues, ...currentBearish]
  fetchData()
}

function selectAllBearish() {
  const bearishValues = bearishTypes.map(t => t.value)
  const currentBullish = selectedTypes.value.filter(t => bullishTypes.some(b => b.value === t))
  selectedTypes.value = [...currentBullish, ...bearishValues]
  fetchData()
}

function selectAllTypes() {
  selectedTypes.value = allTypeValues
  fetchData()
}

function clearAllTypes() {
  selectedTypes.value = []
  fetchData()
}

function handleViewModeChange() {
  paginationReactive.page = 1
  if (viewMode.value === 'realtime') {
    checkTradingTime()
    if (isTrading.value) {
      autoRefresh.value = true
      startAutoRefresh()
    } else {
      autoRefresh.value = false
    }
  } else {
    autoRefresh.value = false
    stopAutoRefresh()
  }
  fetchData()
}

function handlePageChange(currentPage) {
  if (!loadingRef.value) {
    loadingRef.value = true
    paginationReactive.page = currentPage
    fetchHistoryData()
  }
}

function handleSearch() {
  paginationReactive.page = 1
  fetchHistoryData()
}

function handleSearchKeyup(e) {
  if (e.key === 'Enter') {
    handleSearch()
  }
}

async function fetchAllCurrentData() {
  loadingRef.value = true
  try {
    const result = await GetAllStockChangesWithPaging(500)
    if (result) {
      dataRef.value = result.data || []
      paginationReactive.itemCount = result.totalCount || 0
      message.success(t('stockChangesMonitor.changeDataFetched', { count: result.data?.length || 0 }))
    }
  } catch (e) {
    console.error(t('stockChangesMonitor.fetchAllFailed'), e)
    message.error(t('stockChangesMonitor.fetchAllFailed'))
  } finally {
    loadingRef.value = false
  }
}

onBeforeMount(() => {
  GetConfig().then(result => {
    if (result.darkTheme) {
      document.documentElement.classList.add('dark')
    }
  })
})

onMounted(() => {
  checkTradingTime()
  fetchData()
  if (viewMode.value === 'realtime' && isTrading.value && autoRefresh.value) {
    startAutoRefresh()
  }
})

onBeforeUnmount(() => {
  stopAutoRefresh()
})
</script>

<template>
  <n-card>
    <template #header>
      <n-space vertical>
        <n-space justify="space-between" align="center">
          <n-space align="center">
            <n-text strong>{{ t('stockChangesMonitor.title') }}</n-text>
            <n-tag v-if="viewMode === 'realtime'" :type="isTrading ? 'success' : 'warning'" size="small">
              {{ marketStatus }}
            </n-tag>
            <n-tag :type="viewMode === 'realtime' ? 'error' : 'info'" size="small">
              {{ viewMode === 'realtime' ? t('stockChangesMonitor.realtimeData') : t('stockChangesMonitor.historyData') }}
            </n-tag>
            <n-text depth="3" style="font-size: 12px">{{ t('stockChangesMonitor.recordCount', { count: paginationReactive.itemCount }) }}</n-text>
          </n-space>
          <n-space align="center">
            <n-radio-group v-model:value="viewMode" @update:value="handleViewModeChange">
              <n-radio-button value="realtime">{{ t('stockChangesMonitor.realtimeData') }}</n-radio-button>
              <n-radio-button value="history">{{ t('stockChangesMonitor.historyData') }}</n-radio-button>
            </n-radio-group>
            
            <template v-if="viewMode === 'realtime'">
              <n-text v-if="autoRefresh && isTrading" depth="3" style="font-size: 12px">
                {{ countdown }}{{ t('stockChangesMonitor.secondsUntilRefresh') }}
              </n-text>
              <n-text v-else-if="!isTrading" depth="3" style="font-size: 12px">
                {{ t('stockChangesMonitor.pauseRefreshNonTrading') }}
              </n-text>
              <n-switch v-model:value="autoRefresh" @update:value="toggleAutoRefresh" :disabled="!isTrading">
                <template #checked>{{ t('stockChangesMonitor.autoRefresh') }}</template>
                <template #unchecked>{{ t('stockChangesMonitor.manualRefresh') }}</template>
              </n-switch>
            </template>
            
            <n-button @click="fetchData" :loading="loadingRef" type="primary" size="small">
              {{ t('stockChangesMonitor.refresh') }}
            </n-button>
            
            <n-button v-if="viewMode === 'realtime'" @click="saveCurrentData" size="small">
              {{ t('stockChangesMonitor.saveToHistory') }}
            </n-button>
          </n-space>
        </n-space>

        <n-alert v-if="viewMode === 'realtime' && !isTrading" type="info" size="small">
          {{ t('stockChangesMonitor.nonTradingTimeAlert') }}
        </n-alert>

        <n-space v-if="viewMode === 'history'" align="center">
          <n-input 
            v-model:value="paginationReactive.keyword" 
            placeholder="{{ t('stockChangesMonitor.enterStockCodeOrName') }}" 
            clearable 
            style="width: 200px"
            @keyup="handleSearchKeyup"
          />
          <n-date-picker v-model:value="paginationReactive.range" type="daterange" clearable />
          <n-time-picker 
            v-model:value="paginationReactive.startTime" 
            placeholder="{{ t('stockChangesMonitor.startTime') }}" 
            clearable 
            format="HH:mm:ss"
            style="width: 120px"
          />
          <n-time-picker 
            v-model:value="paginationReactive.endTime" 
            placeholder="{{ t('stockChangesMonitor.endTime') }}" 
            clearable 
            format="HH:mm:ss"
            style="width: 120px"
          />
          <n-button type="primary" @click="handleSearch" :loading="loadingRef">
            {{ t('stockChangesMonitor.search') }}
          </n-button>
          <n-button @click="fetchAllCurrentData" :loading="loadingRef">
            {{ t('stockChangesMonitor.getTodayAllData') }}
          </n-button>
        </n-space>
        <n-space v-if="viewMode === 'history'" align="center" style="margin-top: 8px">
          <n-select
            v-model:value="paginationReactive.minVolume"
            :options="volumeOptions"
            :placeholder="t('stockChangesMonitor.volumeFilter')"
            style="width: 120px"
            clearable
          />
          <n-select
            v-model:value="paginationReactive.minAmount"
            :options="amountOptions"
            :placeholder="t('stockChangesMonitor.amountFilter')"
            style="width: 120px"
            clearable
          />
          <n-select
            v-model:value="paginationReactive.minChangeRate"
            :options="changeRateOptions"
            :placeholder="t('stockChangesMonitor.changeRateFilter')"
            style="width: 120px"
            clearable
          />
          <n-input
            v-model:value="paginationReactive.industry"
            :placeholder="t('stockChangesMonitor.industryKeyword')"
            clearable
            style="width: 120px"
          />
          <n-input
            v-model:value="paginationReactive.concept"
            :placeholder="t('stockChangesMonitor.conceptKeyword')"
            clearable
            style="width: 120px"
          />
        </n-space>
        
        <n-space align="center" style="margin-top: 8px">
          <n-text depth="3">{{ t('stockChangesMonitor.changeTypeFilter') }}</n-text>
          <n-button size="tiny" type="primary" @click="selectAllTypes">{{ t('stockChangesMonitor.selectAll') }}</n-button>
          <n-button size="tiny" @click="clearAllTypes">{{ t('stockChangesMonitor.clear') }}</n-button>
          <n-text depth="3" style="font-size: 12px">{{ t('stockChangesMonitor.selectedCount') }}</n-text>
        </n-space>

        <n-space vertical>
          <n-space align="center">
            <n-text style="color: #dc2626; font-weight: 500;">{{ t('stockChangesMonitor.bullishChanges') }}</n-text>
            <n-button size="tiny" @click="selectAllBullish">{{ t('stockChangesMonitor.selectAllBullish') }}</n-button>
          </n-space>
          <n-checkbox-group v-model:value="selectedTypes" @update:value="fetchData">
            <n-space>
              <n-checkbox v-for="item in bullishTypes" :key="item.value" :value="item.value" :label="item.label">
                <template #default>
                  <n-text :style="{color: selectedTypes.includes(item.value) ? '#dc2626' : undefined}">{{ item.label }}</n-text>
                </template>
              </n-checkbox>
            </n-space>
          </n-checkbox-group>
        </n-space>

        <n-space vertical>
          <n-space align="center">
            <n-text style="color: #16a34a; font-weight: 500;">{{ t('stockChangesMonitor.bearishChanges') }}</n-text>
            <n-button size="tiny" @click="selectAllBearish">{{ t('stockChangesMonitor.selectAllBearish') }}</n-button>
          </n-space>
          <n-checkbox-group v-model:value="selectedTypes" @update:value="fetchData">
            <n-space>
              <n-checkbox v-for="item in bearishTypes" :key="item.value" :value="item.value" :label="item.label">
                <template #default>
                  <n-text :style="{color: selectedTypes.includes(item.value) ? '#16a34a' : undefined}">{{ item.label }}</n-text>
                </template>
              </n-checkbox>
            </n-space>
          </n-checkbox-group>
        </n-space>
      </n-space>
    </template>

    <n-data-table
        remote
        :columns="columnsRef"
        :data="dataRef"
        :loading="loadingRef"
        :pagination="viewMode === 'history' ? paginationReactive : false"
        :bordered="false"
        :max-height="500"
        :scroll-x="1300"
        striped
        size="small"
        @update:page="handlePageChange"
    />
  </n-card>
</template>

<style scoped>
</style>
