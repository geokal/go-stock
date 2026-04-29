<script setup lang="ts">
import {onBeforeMount, onUnmounted, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {InvestCalendarTimeLine} from "../../wailsjs/go/main/App";
import { addMonths, format ,parse} from 'date-fns';
import { zhCN, enUS } from 'date-fns/locale';

import {useMessage} from 'naive-ui'
import {Star48Filled} from "@vicons/fluent";
const { t, locale } = useI18n()
const today = new Date();
const year = today.getFullYear();
const month = String(today.getMonth() + 1).padStart(2, '0');
const day = String(today.getDate()).padStart(2, '0');

const formattedDate = `${year}-${month}-${day}`;
const formattedYM = `${year}-${month}`;
const list  = ref([])
const message=useMessage()

function goBackToday() {
  setTimeout(() => {
    nextTick(
        () => {
          const elementById = document.getElementById(formattedDate);
          if (elementById) {
            elementById.scrollIntoView({
              behavior: 'auto',
              block: 'start'
            })
          }
        }
    )
  }, 500)
}

onBeforeMount(() => {
  InvestCalendarTimeLine(formattedYM).then(res => {
    list.value = res
    goBackToday();
  })
})
onMounted(()=>{

})
function loadMore(){
  if (list.value.length>0){
    let day=parse(list.value[list.value.length-1].date, 'yyyy-MM-dd', new Date())
    let nextMonth=addMonths(day,1)
    let ym = format(nextMonth, 'yyyy-MM');
    console.log(ym)
    InvestCalendarTimeLine(ym).then(res => {
      if (res.length==0){
        message.warning(t('calendar.noMoreData'))
        return
      }
      list.value.push( ...res)
    })
  }
}
function getweekday(date){
  let day=parse(date, 'yyyy-MM-dd', new Date())
  return format(day, 'EEEE', {locale: locale.value === 'en' ? enUS : zhCN})
}
</script>

<template>
    <n-list bordered   style="max-height: calc(100vh - 230px);text-align: left;">
      <n-scrollbar style="max-height: calc(100vh - 230px);" >
      <n-list-item v-for="(item, index) in list" :id="item.date" :key="item.date">
          <n-thing :title="item.date+' '+getweekday(item.date)">
            <n-list :bordered="false" hoverable>
              <n-list-item v-for="(l,i ) in item.list" :key="l.article_id	">
                <n-flex justify="space-between">
                <n-text :type="item.date===formattedDate?'warning':'info'">{{i+1}}# {{l.title}}</n-text>
                <n-rate v-if="l.like_count>0" readonly :default-value="l.like_count" :count="l.like_count" >
                  <n-icon :component="Star48Filled"/>
                </n-rate>
                </n-flex>
              </n-list-item>
            </n-list>
          </n-thing>
      </n-list-item>
        <n-list-item v-if="list.length==0">
          <n-text type="info">{{ t('common.noData') }}</n-text>
        </n-list-item>
        <n-list-item v-else style="text-align: center;">
          <n-button-group>
            <n-button  strong secondary type="info" @click="loadMore">{{ t('common.loadMore') }}</n-button>
            <n-button  strong secondary  type="warning" @click="goBackToday">{{ t('common.backToToday') }}</n-button>
          </n-button-group>
        </n-list-item>
      </n-scrollbar>
    </n-list>
</template>

<style scoped>

</style>