<script setup lang="ts">
import {onBeforeMount, onUnmounted, ref} from 'vue'
import {useI18n} from 'vue-i18n'
const { t } = useI18n()
import {HotEvent} from "../../wailsjs/go/main/App";
const list  = ref([])

const task =ref()
onBeforeMount(async () => {
  list.value = await HotEvent(50)
  task.value=setInterval(async ()=>{
    list.value = await HotEvent(50)
  }, 1000*10)
})

onUnmounted(async ()=>{
  clearInterval(task.value)
})
</script>

<template>
  <n-list bordered>
    <template #header>
      {{ t('hotEvents.xueqiuHot') }}
    </template>
    <n-list-item v-for="(item, index) in list" :key="index">
        <n-thing :title="item.tag" :description="item.content"  >
          <template v-if="item.pic" #avatar>
            <n-avatar :src="item.pic" :size="60">
            </n-avatar>
          </template>
        </n-thing>
    </n-list-item>
  </n-list>
</template>

<style scoped>

</style>