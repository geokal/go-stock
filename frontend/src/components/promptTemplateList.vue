<script setup>
import {computed, h, onBeforeMount, onMounted, ref, reactive} from 'vue'
import {
  GetPromptTemplateList,
  GetConfig,
  AddPromptTemplate,
  DeletePromptTemplate,
  UpdatePromptTemplate
} from "../../wailsjs/go/main/App";
import { EventsEmit } from "../../wailsjs/runtime";
import {NButton, NInput, NTag, NText, NSwitch, useMessage, useNotification,useDialog, NModal, NCard, NForm, NFormItem, NSpace, NPopover} from "naive-ui";
import { MdEditor, MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const notify = useNotification()
const message = useMessage()
const dialog = useDialog()
const editorDataRef = reactive({
  darkTheme: false
})
const editorTheme = ref('light')

onBeforeMount(() => {
  GetConfig().then(result => {
    if (result.darkTheme) {
      editorDataRef.darkTheme = true
      editorTheme.value = 'dark'
    }
  })
})

onMounted(() => {
  query({
    page: 1,
    pageSize: paginationReactive.pageSize
  }).then((data) => {
    dataRef.value = data.data
    paginationReactive.page = 1
    paginationReactive.pageCount = data.totalPages
    paginationReactive.itemCount = data.total
    loadingRef.value = false
  })
})

const dataRef = ref([])
const loadingRef = ref(true)

const columnsRef = ref([
  {
    title: t('promptTemplateList.name'),
    key: 'name',
    render(row) {
      if (row.type === '模型系统Prompt') {
        return h(NText, { type: "success" }, { default: () => row.name })
      }else{
        return h(NText, { type: "info" }, { default: () => row.name })
      }
    }
  },
  {
    title: t('promptTemplateList.type'),
    key: 'type',
    render(row) {
      const label = row.type === '模型系统Prompt' ? t('promptTemplateList.systemPromptType') : t('promptTemplateList.userPromptType')
      if (row.type === '模型系统Prompt') {
        return h(NTag, { type: "success" }, { default: () => label })
      }else{
        return h(NTag, { type: "info" }, { default: () => label })
      }
    }
  },
  {
    title: t('promptTemplateList.createdAt'),
    key: 'CreatedAt',
    render(row) {
      return row.CreatedAt.substring(0, 19).replace('T', ' ')
    }
  },
  {
    title: t('promptTemplateList.updatedAt'),
    key: 'UpdatedAt',
    render(row) {
      return row.UpdatedAt.substring(0, 19).replace('T', ' ')
    }
  },
  {
    title: t('promptTemplateList.content'),
    key: 'content',
    width: 200,
    render(row) {
      return h(NPopover, {
        trigger: 'hover',
        placement: 'left',
        showArrow: true,
        style: 'max-width: 800px; max-height: 400px; overflow: hidden',
        scrollable: true
      }, {
        trigger: () => h('span', {
          style: 'display: inline-block; max-width: 180px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; cursor: pointer;'
        }, row.content),
        default: () => h(MdPreview, {
          style:'text-align: left;',
          modelValue: row.content,
          theme: editorTheme.value
        })
      })
    }
  },
  {
    title: t('promptTemplateList.operation'),
    width: 260,
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-right: 5px',
            onClick: () => showEditModal(row)
          },
          { default: () => t('promptTemplateList.edit') }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'info',
            style: 'margin-right: 5px',
            onClick: () => showShareModal(row)
          },
          { default: () => t('promptTemplateList.share') }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            onClick: () => deletePromptTemplate(row.ID)
          },
          { default: () => t('promptTemplateList.delete') }
        )
      ]
    }
  }
])

const paginationReactive = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 12,
  itemCount: 0,
  prefix({ itemCount }) {
    return t('promptTemplateList.recordCount', { count: itemCount })
  }
})

const modalDataRef = reactive({
  visible: false,
  isEdit: false,
  formData: {
    ID: 0,
    name: '',
    type: '',
    content: ''
  }
})

const shareDataRef = reactive({
  visible: false,
  title: '',
  content: '',
  description: '',
  category: '',
  tags: '',
  isPublic: true,
  loading: false
})

const promptPlazaApiBase = ref('http://go-stock.sparkmemory.top:1918/api')

function query({ page, pageSize = 10, name = "", type = "", content = "" }) {
  return new Promise((resolve) => {
    GetPromptTemplateList({
      "page": page,
      "pageSize": pageSize,
      "name": name,
      "type": type,
      "content": content
    }).then((res) => {
      resolve({
        data: res.list,
        total: res.total,
        totalPages: res.totalPages
      })
    })
  })
}

function handlePageChange(currentPage) {
  if (!loadingRef.value) {
    loadingRef.value = true
    query({
      page: currentPage,
      pageSize: paginationReactive.pageSize,
      name: searchFormRef.name,
      type: searchFormRef.type,
      content: searchFormRef.content
    }).then((data) => {
      dataRef.value = data.data
      paginationReactive.page = currentPage
      paginationReactive.pageCount = data.totalPages
      paginationReactive.itemCount = data.total
      loadingRef.value = false
    })
  }
}
const promptTypeOptions = [
  {label: t('promptTemplateList.systemPromptType'), value: '模型系统Prompt'},
  {label: t('promptTemplateList.userPromptType'), value: '模型用户Prompt'},]
const searchFormRef = reactive({
  name: "",
  type: null,
  content: ""
})

function handleSearch() {
  if (!loadingRef.value) {
    loadingRef.value = true
    query({
      page: paginationReactive?.page ?? 1,
      pageSize: paginationReactive.pageSize,
      name: searchFormRef.name,
      type: searchFormRef.type,
      content: searchFormRef.content
    }).then((data) => {
      dataRef.value = data.data
      paginationReactive.page = data.page
      paginationReactive.pageCount = data.totalPages
      paginationReactive.itemCount = data.total
      loadingRef.value = false
    })
  }
}

function showAddModal() {
  modalDataRef.isEdit = false
  modalDataRef.formData = {
    ID: 0,
    name: '',
    type: '',
    content: ''
  }
  modalDataRef.visible = true
}

function showEditModal(row) {
  modalDataRef.isEdit = true
  modalDataRef.formData = {
    ID: row.ID,
    name: row.name,
    type: row.type,
    content: row.content
  }
  modalDataRef.visible = true
}

function savePromptTemplate() {
  if (!modalDataRef.formData.name || !modalDataRef.formData.type || !modalDataRef.formData.content) {
    message.warning(t('promptTemplateList.fillCompleteInfo') )
    return
  }

  const apiCall = modalDataRef.isEdit ? UpdatePromptTemplate : AddPromptTemplate
  apiCall(modalDataRef.formData).then((res) => {
    message.info( res )
    modalDataRef.visible = false
    handleSearch()
    EventsEmit('promptTemplatesChanged')
  })
}

function deletePromptTemplate(id) {

  dialog.warning({
    title: t('promptTemplateList.tip'),
    content: t('promptTemplateList.deleteConfirm'),
    positiveText: t('promptTemplateList.confirm'),
    negativeText: t('promptTemplateList.cancel'),
    onPositiveClick: () => {
      DeletePromptTemplate(id).then((res) => {
        message.info( res )
        handleSearch()
        EventsEmit('promptTemplatesChanged')
      })
    }
  })
}

function showShareModal(row) {
  shareDataRef.title = row.name || ''
  shareDataRef.content = row.content || ''
  shareDataRef.description = ''
  shareDataRef.category = row.type || ''
  shareDataRef.tags = ''
  shareDataRef.isPublic = true
  shareDataRef.visible = true
  GetConfig().then(result => {
    if (result.promptPlazaApiBase) {
      promptPlazaApiBase.value = result.promptPlazaApiBase
    }
  })
}

async function handleShare() {
  if (!shareDataRef.title || !shareDataRef.content) {
    message.warning(t('promptTemplateList.titleContentRequired'))
    return
  }
  const token = localStorage.getItem('promptPlazaToken')
  if (!token) {
    message.warning(t('promptTemplateList.loginFirst'))
    return
  }
  shareDataRef.loading = true
  try {
    const resp = await fetch(promptPlazaApiBase.value + '/prompts', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        title: shareDataRef.title,
        content: shareDataRef.content,
        description: shareDataRef.description,
        category: shareDataRef.category,
        tags: shareDataRef.tags,
        isPublic: shareDataRef.isPublic
      })
    })
    const json = await resp.json()
    if (json.code !== 0) {
      if (json.code === 401) {
        message.error(t('promptTemplateList.loginExpired'))
      } else {
        message.error(t('promptTemplateList.shareFailed') + (json.message || t('promptTemplateList.unknownError')))
      }
      return
    }
    message.success(t('promptTemplateList.shareSuccess'))
    shareDataRef.visible = false
  } catch (e) {
    message.error(t('promptTemplateList.shareFailed') + e.message)
  } finally {
    shareDataRef.loading = false
  }
}
</script>

<template>
  <div>
    <!-- 搜索区域 -->
    <n-space vertical style="margin-bottom: 16px">
      <n-space>
        <n-input v-model:value="searchFormRef.name" :placeholder="t('promptTemplateList.namePlaceholder')" clearable />
        <n-select style="width: 200px" v-model:value="searchFormRef.type" :options="promptTypeOptions" :placeholder="t('promptTemplateList.typePlaceholder')" clearable/>
        <n-input v-model:value="searchFormRef.content" :placeholder="t('promptTemplateList.contentPlaceholder')" clearable />
        <n-button type="success" @click="handleSearch">{{ t('promptTemplateList.search') }}</n-button>
        <n-button type="warning" @click="showAddModal">{{ t('promptTemplateList.addTemplate') }}</n-button>
      </n-space>
    </n-space>

    <!-- 数据表格 -->
    <n-data-table
      remote
      size="small"
      :columns="columnsRef"
      :data="dataRef"
      :loading="loadingRef"
      :pagination="paginationReactive"
      :row-key="(rowData) => rowData.ID"
      @update:page="handlePageChange"
      flex-height
      style="height: calc(100vh - 250px)"
    />

    <!-- 编辑/新增模态框 -->
    <n-modal v-model:show="modalDataRef.visible" preset="card" style="width: 1100px;text-align: left" :title="modalDataRef.formData.ID>0? t('promptTemplateList.modify') : t('promptTemplateList.add') + t('promptTemplateList.promptTemplate')">
      <n-form :model="modalDataRef.formData" label-placement="left" label-width="80">
        <n-form-item :label="t('promptTemplateList.name')" required>
          <n-input v-model:value="modalDataRef.formData.name" :placeholder="t('promptTemplateList.nameInputPlaceholder')" />
        </n-form-item>
        <n-form-item :label="t('promptTemplateList.type')" required>
          <n-select v-model:value="modalDataRef.formData.type" :options="promptTypeOptions" :placeholder="t('promptTemplateList.typeSelectPlaceholder')"/>
        </n-form-item>
        <n-form-item :label="t('promptTemplateList.content')" required>
          <MdEditor
            v-model="modalDataRef.formData.content"
            style="height: 400px"
            :theme="editorTheme"
            :preview="true"
            :toolbarsExclude="['github', 'htmlPreview', 'catalog', 'save']"
            :placeholder="t('promptTemplateList.contentInputPlaceholder')"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="modalDataRef.visible = false">{{ t('promptTemplateList.cancel') }}</n-button>
          <n-button type="primary" @click="savePromptTemplate">{{ t('promptTemplateList.save') }}</n-button>
        </n-space>
      </template>
    </n-modal>

    <n-modal v-model:show="shareDataRef.visible" preset="card" style="width: 700px;text-align: left" :title="t('promptTemplateList.shareToPlaza')">
      <n-form :model="shareDataRef" label-placement="left" label-width="80">
        <n-form-item :label="t('promptTemplateList.title')" required>
          <n-input v-model:value="shareDataRef.title" :placeholder="t('promptTemplateList.titleInputPlaceholder')" />
        </n-form-item>
        <n-space :size="8">
          <n-form-item :label="t('promptTemplateList.category')" label-placement="left" style="width: 300px">
            <n-input v-model:value="shareDataRef.category" :placeholder="t('promptTemplateList.categoryPlaceholder')" />
          </n-form-item>
          <n-form-item :label="t('promptTemplateList.tags')" label-placement="left" style="width: 300px">
            <n-input v-model:value="shareDataRef.tags" :placeholder="t('promptTemplateList.tagsPlaceholder')" />
          </n-form-item>
        </n-space>
        <n-form-item :label="t('promptTemplateList.description')">
          <n-input v-model:value="shareDataRef.description" type="textarea" :rows="2" :placeholder="t('promptTemplateList.descriptionPlaceholder')" />
        </n-form-item>
        <n-form-item :label="t('promptTemplateList.content')" required>
          <n-input v-model:value="shareDataRef.content" type="textarea" :rows="6" :placeholder="t('promptTemplateList.contentInputPlaceholder')" />
        </n-form-item>
        <n-form-item :label="t('promptTemplateList.public')">
          <n-switch v-model:value="shareDataRef.isPublic" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="shareDataRef.visible = false">{{ t('promptTemplateList.cancel') }}</n-button>
          <n-button type="primary" :loading="shareDataRef.loading" @click="handleShare">{{ t('promptTemplateList.share') }}</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<style scoped>
:deep(.md-editor) {
  text-align: left;
}
:deep(.n-popover .md-editor-preview) {
  padding: 8px 12px;
}
:deep(.n-popover .md-editor-preview-wrapper) {
  padding: 0;
}
</style>
