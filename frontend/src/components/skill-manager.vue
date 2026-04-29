<template>
  <n-space vertical style="margin-bottom: 12px">
    <n-space>
      <n-input
        v-model:value="searchKeyword"
        :placeholder="t('skill.searchPlaceholder')"
        style="width: 200px"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <n-icon :component="SearchOutline" />
        </template>
      </n-input>

      <n-select
        v-model:value="filterCategory"
        :options="categoryOptions"
        :placeholder="t('skill.category')"
        style="width: 120px"
        clearable
        filterable
      />

      <n-select
        v-model:value="filterEnable"
        :options="enableOptions"
        :placeholder="t('skill.enableStatus')"
        style="width: 100px"
        clearable
      />

      <n-button type="primary" @click="handleSearch">
        {{ t('common.search') }}
      </n-button>

      <n-button type="warning" @click="handleCreate">
        <template #icon>
          <n-icon :component="AddOutline" />
        </template>
        {{ t('skill.addSkill') }}
      </n-button>
    </n-space>

    <n-data-table
      remote
      :columns="columns"
      :data="tableData"
      :pagination="pagination"
      :loading="loading"
      :row-key="row => row.id"
      @update:page="handlePageChange"
    />
  </n-space>

  <n-modal
    v-model:show="showCreateModal"
    preset="card"
    :title="editingSkill ? t('skill.editSkill') : t('skill.addSkill')"
    style="width: 900px; max-height: 85vh"
    :mask-closable="false"
  >
    <n-scrollbar style="max-height: calc(85vh - 120px)">
    <n-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-placement="top"
      label-align="left"
    >
      <n-grid :cols="4" :x-gap="16">
        <n-form-item-gi :label="t('skill.name')" path="name" :span="2">
          <n-input v-model:value="formData.name" :placeholder="t('skill.enterName')" clearable />
        </n-form-item-gi>

        <n-form-item-gi :label="t('skill.category')" path="category">
          <n-select
            v-model:value="formData.category"
            :options="categoryOptions"
            :placeholder="t('skill.selectOrEnterCategory')"
            clearable
            filterable
            tag
          />
        </n-form-item-gi>

        <n-form-item-gi :label="t('skill.sortOrder')" path="sortOrder">
          <n-input-number v-model:value="formData.sortOrder" :min="0" :max="999" style="width: 100%" />
        </n-form-item-gi>
      </n-grid>

      <n-grid :cols="4" :x-gap="16">
        <n-form-item-gi :label="t('skill.enable')" path="enable" :span="1">
          <n-switch v-model:value="formData.enable" />
        </n-form-item-gi>

        <n-form-item-gi :label="t('skill.triggerKeywords')" path="triggerKeywords" :span="3">
          <n-input
            v-model:value="formData.triggerKeywords"
            :placeholder="t('skill.triggerKeywordsPlaceholder')"
            clearable
          />
        </n-form-item-gi>
      </n-grid>

      <n-form-item :label="t('skill.description')" path="description">
        <n-input
          v-model:value="formData.description"
          type="textarea"
          :autosize="{ minRows: 1, maxRows: 3 }"
          :placeholder="t('skill.enterDescription')"
          show-count
          maxlength="500"
        />
      </n-form-item>

      <n-form-item :label="t('skill.bindMcp')" path="mcpServerIds">
        <n-select
          v-model:value="formData.mcpServerIds"
          :options="mcpServerOptions"
          :placeholder="t('skill.selectMcpServer')"
          multiple
          clearable
        />
      </n-form-item>

      <n-form-item :label="t('skill.systemPrompt')" path="systemPrompt">
        <MdEditor
          v-model="formData.systemPrompt"
          style="height: 200px"
          :theme="editorTheme"
          :preview="true"
          :toolbarsExclude="['github', 'htmlPreview', 'catalog', 'save']"
          :placeholder="t('skill.systemPromptPlaceholder')"
        />
      </n-form-item>

      <n-form-item :label="t('skill.examples')" path="examples">
        <MdEditor
          v-model="formData.examples"
          style="height: 160px"
          :theme="editorTheme"
          :preview="true"
          :toolbarsExclude="['github', 'htmlPreview', 'catalog', 'save']"
          :placeholder="t('skill.examplesPlaceholder')"
        />
      </n-form-item>
    </n-form>
    </n-scrollbar>

    <template #footer>
      <n-space justify="end">
        <n-button @click="showCreateModal = false">{{ t('common.cancel') }}</n-button>
        <n-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ editingSkill ? t('common.save') : t('common.add') }}
        </n-button>
      </n-space>
    </template>
  </n-modal>

</template>

<script setup>
import { ref, reactive, h, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NButton, NSpace, NInput, NDataTable, NModal, NForm, NFormItem,
  NFormItemGi, NGrid, NTag, NSwitch, NIcon, NSelect, NInputNumber, NPopconfirm, NScrollbar, useMessage
} from 'naive-ui'
import { SearchOutline, AddOutline, TrashOutline, CreateOutline, FlashOutline } from '@vicons/ionicons5'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { CreateSkill, UpdateSkill, DeleteSkill, GetSkillList, EnableSkill, GetSkillByID, GetAllSkills } from '../../wailsjs/go/main/App.js'
import { GetMCPServerList, GetConfig } from '../../wailsjs/go/main/App.js'

const { t } = useI18n()
const message = useMessage()
const loading = ref(false)
const submitting = ref(false)
const searchKeyword = ref('')
const filterCategory = ref(null)
const filterEnable = ref(null)
const showCreateModal = ref(false)
const editingSkill = ref(false)
const formRef = ref(null)
const tableData = ref([])
const mcpServerOptions = ref([])
const editorTheme = ref('light')


const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  prefix: ({ itemCount }) => `${t('common.total')} ${itemCount} ${t('common.records')}`,
  onChange: (page) => {
    handlePageChange(page)
  },
  onUpdatePageSize: (size) => {
    pageSize.value = size
    pagination.pageSize = size
    currentPage.value = 1
    pagination.page = 1
    loadData()
  }
})

const formData = reactive({
  id: null,
  name: '',
  description: '',
  category: null,
  systemPrompt: '',
  examples: '',
  triggerKeywords: '',
  mcpServerIds: [],
  enable: true,
  sortOrder: 0
})

const formRules = {
  name: { required: true, message: t('skill.enterName'), trigger: ['input', 'blur'] }
}

const categoryOptions = [
  { label: t('skill.categoryStockAnalysis'), value: t('skill.categoryStockAnalysis') },
  { label: t('skill.categoryTechAnalysis'), value: t('skill.categoryTechAnalysis') },
  { label: t('skill.categoryFundamentalAnalysis'), value: t('skill.categoryFundamentalAnalysis') },
  { label: t('skill.categoryQuantStrategy'), value: t('skill.categoryQuantStrategy') },
  { label: t('skill.categoryRiskManagement'), value: t('skill.categoryRiskManagement') },
  { label: t('skill.categoryInfoResearch'), value: t('skill.categoryInfoResearch') },
  { label: t('skill.categoryGeneral'), value: t('skill.categoryGeneral') }
]

const enableOptions = [
  { label: t('skill.enabled'), value: true },
  { label: t('skill.disabled'), value: false }
]

const columns = [
  {
    title: t('skill.id'),
    key: 'id',
    width: 50
  },
  {
    title: t('skill.name'),
    key: 'name',
    width: 120,
    ellipsis: { tooltip: true }
  },
  {
    title: t('skill.category'),
    key: 'category',
    width: 90,
    render(row) {
      if (!row.category) return h(NTag, { type: 'default' }, { default: () => t('skill.uncategorized') })
      return h(NTag, { type: 'info' }, { default: () => row.category })
    }
  },
  {
    title: t('skill.description'),
    key: 'description',
    width: 200,
    ellipsis: { tooltip: { style: { maxWidth: '400px', wordBreak: 'break-all' } } }
  },
  {
    title: t('skill.bindMcp'),
    key: 'mcpServerIds',
    width: 100,
    render(row) {
      if (!row.mcpServerIds) return h(NTag, { type: 'default' }, { default: () => t('skill.none') })
      const ids = row.mcpServerIds.split(',').filter(s => s.trim())
      return h(NTag, { type: 'info' }, { default: () => `${ids.length} ${t('skill.count')}` })
    }
  },
  {
    title: t('skill.sortOrder'),
    key: 'sortOrder',
    width: 60
  },
  {
    title: t('skill.enable'),
    key: 'enable',
    width: 70,
    render(row) {
      return h(NSwitch, {
        value: row.enable,
        onUpdateValue: (val) => handleEnable(row, val)
      })
    }
  },
  {
    title: t('common.actions'),
    key: 'actions',
    width: 140,
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, {
            size: 'small', type: 'info', quaternary: true,
            onClick: () => handleEdit(row)
          }, {
            icon: () => h(NIcon, null, { default: () => h(CreateOutline) }),
            default: () => t('common.edit')
          }),
          h(NPopconfirm, {
            onPositiveClick: () => handleDelete(row)
          }, {
            trigger: () => h(NButton, {
              size: 'small', type: 'error', quaternary: true
            }, {
              icon: () => h(NIcon, null, { default: () => h(TrashOutline) }),
              default: () => t('common.delete')
            }),
            default: () => t('skill.confirmDelete')
          })
        ]
      })
    }
  }
]

const loadData = async () => {
  loading.value = true
  try {
    const result = await GetSkillList({
      page: currentPage.value,
      pageSize: pageSize.value,
      name: searchKeyword.value,
      category: filterCategory.value,
      enable: filterEnable.value
    })
    if (result) {
      tableData.value = result.data || []
      total.value = result.total || 0
      pagination.itemCount = total.value
      pagination.pageCount = Math.ceil(total.value / pageSize.value) || 1
    }
  } catch (error) {
    message.error(t('skill.loadFailed') + error)
  } finally {
    loading.value = false
  }
}

const loadMCPServers = async () => {
  try {
    const result = await GetMCPServerList({
      page: 1,
      pageSize: 100,
      name: '',
      status: '',
      enable: true
    })
    if (result && result.data) {
      mcpServerOptions.value = result.data.map(s => ({
        label: s.name,
        value: String(s.id)
      }))
    }
  } catch (error) {
    console.error(t('skill.loadMcpFailed'), error)
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
  pagination.page = page
  loadData()
}

const handleSearch = () => {
  currentPage.value = 1
  pagination.page = 1
  loadData()
}

const handleCreate = () => {
  editingSkill.value = false
  resetForm()
  showCreateModal.value = true
}

const handleEdit = async (row) => {
  editingSkill.value = true
  try {
    const skill = await GetSkillByID(row.id)
    if (skill) {
      resetForm()
      formData.id = skill.id
      formData.name = skill.name
      formData.description = skill.description
      formData.category = skill.category
      formData.systemPrompt = skill.systemPrompt
      formData.examples = skill.examples
      formData.triggerKeywords = skill.triggerKeywords
      formData.mcpServerIds = skill.mcpServerIds ? skill.mcpServerIds.split(',').filter(s => s.trim()) : []
      formData.enable = skill.enable
      formData.sortOrder = skill.sortOrder
      showCreateModal.value = true
    }
  } catch (error) {
    message.error(t('skill.getSkillFailed') + error)
  }
}

const handleDelete = async (row) => {
  try {
    const result = await DeleteSkill(row.id)
    if (result.includes(t('common.success'))) {
      message.success(result)
      loadData()
    } else {
      message.error(result)
    }
  } catch (error) {
    message.error(t('common.delete') + t('common.failed') + error)
  }
}

const handleEnable = async (row, enable) => {
  try {
    const result = await EnableSkill(row.id, enable)
    if (result.includes(t('common.success')) || result.includes(t('skill.enabled')) || result.includes(t('skill.disabled'))) {
      message.success(result)
      loadData()
    } else {
      message.error(result)
    }
  } catch (error) {
    message.error(t('common.operation') + t('common.failed') + error)
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  submitting.value = true
  try {
    const skillData = {
      name: formData.name,
      description: formData.description,
      category: formData.category || '',
      systemPrompt: formData.systemPrompt,
      examples: formData.examples,
      triggerKeywords: formData.triggerKeywords,
      mcpServerIds: formData.mcpServerIds.join(','),
      enable: formData.enable,
      sortOrder: formData.sortOrder
    }

    let result
    if (editingSkill.value) {
      skillData.id = formData.id
      result = await UpdateSkill(skillData)
    } else {
      result = await CreateSkill(skillData)
    }

    if (result.includes(t('common.success'))) {
      message.success(result)
      showCreateModal.value = false
      loadData()
    } else {
      message.error(result)
    }
  } catch (error) {
    message.error(t('common.operation') + t('common.failed') + error)
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  Object.assign(formData, {
    id: null,
    name: '',
    description: '',
    category: null,
    systemPrompt: '',
    examples: '',
    triggerKeywords: '',
    mcpServerIds: [],
    enable: true,
    sortOrder: 0
  })
  if (formRef.value) {
    formRef.value.restoreValidation()
  }
}

onMounted(() => {
  loadData()
  loadMCPServers()
  GetConfig().then(result => {
    if (result.darkTheme) {
      editorTheme.value = 'dark'
    }
  })
})
</script>

<style scoped>
:deep(.md-editor) {
  text-align: left;
}
</style>
