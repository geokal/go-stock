<template>
  <n-space vertical style="margin-bottom: 12px">
    <n-space>
      <n-input
        v-model:value="searchKeyword"
        :placeholder="t('mcpServer.searchPlaceholder')"
        style="width: 200px"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <n-icon :component="SearchOutline" />
        </template>
      </n-input>

      <n-select
        v-model:value="filterStatus"
        :options="statusOptions"
        :placeholder="t('mcpServer.serverStatus')"
        style="width: 120px"
        clearable
      />

      <n-button type="primary" @click="handleSearch">
        {{ t('common.search') }}
      </n-button>

      <n-button type="warning" @click="handleCreate">
        <template #icon>
          <n-icon :component="AddOutline" />
        </template>
        {{ t('mcpServer.createServer') }}
      </n-button>
    </n-space>
  </n-space>

  <n-data-table
    remote
    size="small"
    :columns="columns"
    :data="serverList"
    :loading="loading"
    :pagination="pagination"
    :row-key="(rowData) => rowData.id"
    :expanded-row-keys="expandedKeys"
    :row-props="rowProps"
    @update:expanded-row-keys="handleExpand"
    @update:page="handlePageChange"
    flex-height
    style="height: calc(100vh - 210px); margin-top: 10px"
  />

  <n-modal
    v-model:show="showCreateModal"
    :title="editingServer ? t('mcpServer.modifyServer') : t('mcpServer.createNewServer')"
    preset="dialog"
    :style="{ width: '750px' }"
    @close="resetForm"
    :z-index="2000"
    to="body"
  >
    <n-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-placement="left"
      label-width="130px"
      require-mark-placement="right-hanging"
    >
      <n-form-item :label="t('mcpServer.serverName')" path="name">
        <n-input v-model:value="formData.name" :placeholder="t('mcpServer.enterServerName')" clearable />
      </n-form-item>

      <n-form-item :label="t('mcpServer.description')" path="description">
        <n-input
          v-model:value="formData.description"
          type="textarea"
          :rows="2"
          :placeholder="t('mcpServer.enterServerDescription')"
          show-count
          maxlength="500"
        />
      </n-form-item>

      <n-form-item label="URL" path="url">
        <n-input v-model:value="formData.url" :placeholder="t('mcpServer.urlPlaceholder')" clearable />
      </n-form-item>

      <n-form-item :label="t('mcpServer.envVars')" path="env">
        <n-input
          v-model:value="formData.env"
          type="textarea"
          :rows="3"
          :placeholder="t('mcpServer.envPlaceholder')"
          show-count
        />
      </n-form-item>

      <n-form-item :label="t('mcpServer.enableStatus')" path="enable">
        <n-switch v-model:value="formData.enable" size="large">
          <template #checked>
            <n-icon :component="PlayCircleOutline" />
            {{ t('mcpServer.enable') }}
          </template>
          <template #unchecked>
            <n-icon :component="StopCircleOutline" />
            {{ t('mcpServer.disable') }}
          </template>
        </n-switch>
      </n-form-item>
    </n-form>

    <template #action>
      <n-button @click="showCreateModal = false">{{ t('common.cancel') }}</n-button>
      <n-button type="primary" @click="handleSubmit" :loading="submitting">
        <template #icon>
          <n-icon :component="CheckmarkCircleOutline" />
        </template>
        {{ editingServer ? t('mcpServer.modifyServer') : t('mcpServer.createNewServer') }}
      </n-button>
    </template>
  </n-modal>

  <n-modal
    v-model:show="showToolDetailModal"
    :title="t('mcpServer.toolParamsDetail')"
    preset="card"
    style="width: 850px"
    :z-index="2000"
  >
    <template v-if="currentTool">
      <n-descriptions bordered label-placement="top" :column="1" style="margin-bottom: 16px" content-style="text-align: left">
        <n-descriptions-item :label="t('mcpServer.toolName')">
          <n-text code>{{ currentTool.toolName }}</n-text>
        </n-descriptions-item>
        <n-descriptions-item :label="t('mcpServer.description')">{{ currentTool.description || t('mcpServer.noDescription') }}</n-descriptions-item>
      </n-descriptions>

      <template v-if="parsedParams.length > 0">
        <n-text strong style="margin-bottom: 8px; display: block">{{ t('mcpServer.paramsList') }}</n-text>
        <n-data-table
          :columns="paramDetailColumns"
          :data="parsedParams"
          size="small"
          bordered
          :pagination="false"
        />
      </template>
      <n-text v-else depth="3">{{ t('mcpServer.noParamsNeeded') }}</n-text>

      <n-collapse style="margin-top: 12px" v-if="currentTool.paramsSchema">
        <n-collapse-item :title="t('mcpServer.rawJsonSchema')" name="raw">
          <VueJsonPretty
            :data="parseJSON(currentTool.paramsSchema)"
            :deep="3"
            show-length
            show-line
            collapsed-on-click-bracket
          />
        </n-collapse-item>
      </n-collapse>
    </template>
  </n-modal>
</template>

<script setup>
import { ref, reactive, onMounted, computed, h } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NButton,
  NIcon,
  NTag,
  NSpace,
  NPopconfirm,
  NDescriptions,
  NDescriptionsItem,
  NText,
  NDataTable,
  NCollapse,
  NCollapseItem,
  useMessage
} from 'naive-ui'
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'
import {
  SearchOutline,
  AddOutline,
  PlayOutline,
  PauseOutline,
  TrashOutline,
  CreateOutline,
  PlayCircleOutline,
  StopCircleOutline,
  CheckmarkCircleOutline,
  FlashOutline,
  EyeOutline
} from '@vicons/ionicons5'
import {
  CreateMCPServer,
  UpdateMCPServer,
  DeleteMCPServer,
  GetMCPServerByID,
  GetMCPServerList,
  EnableMCPServer,
  TestMCPServer,
  GetMCPToolsByServerID,
  GetAllMCPTools
} from '../../wailsjs/go/main/App'

const { t } = useI18n()
const message = useMessage()

const formRef = ref(null)

const loading = ref(false)
const submitting = ref(false)
const showCreateModal = ref(false)
const showToolDetailModal = ref(false)
const editingServer = ref(false)
const searchKeyword = ref('')
const filterStatus = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const expandedKeys = ref([])
const currentTool = ref(null)

const formData = reactive({
  id: null,
  name: '',
  description: '',
  url: '',
  command: '',
  args: '',
  env: '',
  enable: true,
  status: 'untested'
})

const formRules = {
  name: { required: true, message: t('mcpServer.enterServerName'), trigger: ['input', 'blur'] },
  command: { required: true, message: t('mcpServer.enterCommand'), trigger: ['input', 'blur'] }
}

const statusOptions = [
  { label: t('mcpServer.available'), value: 'available' },
  { label: t('mcpServer.untested'), value: 'untested' },
  { label: t('mcpServer.unavailable'), value: 'unavailable' }
]

const getStatusLabel = (status) => {
  switch (status) {
    case 'available':
      return t('mcpServer.available')
    case 'untested':
      return t('mcpServer.untested')
    case 'unavailable':
      return t('mcpServer.unavailable')
    default:
      return status
  }
}

const formatJSON = (str) => {
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

const parseJSON = (str) => {
  try {
    return JSON.parse(str)
  } catch {
    return str
  }
}

const handleExpand = (keys) => {
  expandedKeys.value = keys
}

const rowProps = (row) => {
  return {
    style: 'cursor: pointer',
    onClick: (e) => {
      if (e.target.closest('.n-button, .n-popconfirm, .n-switch, a')) return
      const key = row.id
      const idx = expandedKeys.value.indexOf(key)
      if (idx === -1) {
        expandedKeys.value = [...expandedKeys.value, key]
      } else {
        expandedKeys.value = expandedKeys.value.filter(k => k !== key)
      }
    }
  }
}

const handleViewToolDetail = (tool) => {
  currentTool.value = tool
  showToolDetailModal.value = true
}

const parsedParams = computed(() => {
  if (!currentTool.value || !currentTool.value.paramsSchema) return []
  try {
    const schema = JSON.parse(currentTool.value.paramsSchema)
    const props = schema.properties || {}
    const required = schema.required || []
    return Object.entries(props).map(([name, prop]) => ({
      name,
      type: prop.type || '-',
      required: required.includes(name),
      description: prop.description || prop.desc || '-',
      enum: prop.enum ? prop.enum.join(', ') : '',
      default: prop.default !== undefined ? String(prop.default) : ''
    }))
  } catch {
    return []
  }
})

const paramDetailColumns = [
  {
    title: t('mcpServer.paramName'),
    key: 'name',
    width: 180,
    ellipsis: { tooltip: { style: { maxWidth: '300px' } } }
  },
  {
    title: t('mcpServer.type'),
    key: 'type',
    width: 80
  },
  {
    title: t('mcpServer.required'),
    key: 'required',
    width: 60,
    render(row) {
      return h(NTag, { type: row.required ? 'error' : 'default', size: 'small' }, {
        default: () => row.required ? t('mcpServer.yes') : t('mcpServer.no')
      })
    }
  },
  {
    title: t('mcpServer.description'),
    key: 'description',
    ellipsis: { tooltip: { style: { maxWidth: '400px', wordBreak: 'break-all' } } }
  },
  {
    title: t('mcpServer.enumValues'),
    key: 'enum',
    width: 120,
    ellipsis: { tooltip: { style: { maxWidth: '300px' } } }
  },
  {
    title: t('mcpServer.defaultValue'),
    key: 'default',
    width: 80,
    ellipsis: { tooltip: true }
  }
]

const columns = [
  {
    type: 'expand',
    renderExpand: (row) => {
      const tools = row.tools
      if (!tools || tools.length === 0) {
        return h(NText, { depth: 3, style: 'padding: 8px 16px' }, { default: () => t('mcpServer.noToolInfo') })
      }

      const toolColumns = [
        {
          title: t('mcpServer.toolName'),
          key: 'toolName',
          width: 200,
          ellipsis: { tooltip: true }
        },
        {
          title: t('mcpServer.description'),
          key: 'description',
          ellipsis: { tooltip: { style: { maxWidth: '400px', wordBreak: 'break-all' } } }
        },
        {
          title: t('mcpServer.params'),
          key: 'paramsSchema',
          width: 260,
          render(toolRow) {
            if (!toolRow.paramsSchema) {
              return h(NTag, { type: 'default', size: 'small' }, { default: () => t('mcpServer.noParams') })
            }
            try {
              const schema = JSON.parse(toolRow.paramsSchema)
              const props = schema.properties || {}
              const required = schema.required || []
              const paramNames = Object.keys(props)
              if (paramNames.length === 0) {
                return h(NTag, { type: 'default', size: 'small' }, { default: () => t('mcpServer.noParams') })
              }
              const tags = paramNames.map(name => {
                const prop = props[name]
                const isReq = required.includes(name)
                return h(NTag, {
                  size: 'small',
                  type: isReq ? 'info' : 'default',
                  style: 'margin: 2px'
                }, {
                  default: () => name + (prop.type ? ':' + prop.type : '') + (isReq ? '*' : '')
                })
              })
              return h('div', { style: 'display: flex; flex-wrap: wrap; align-items: center; gap: 0' }, [
                ...tags,
                h(NButton, {
                  size: 'tiny',
                  type: 'info',
                  quaternary: true,
                  style: 'margin-left: 4px',
                  onClick: () => handleViewToolDetail(toolRow)
                }, {
                  icon: () => h(NIcon, { component: EyeOutline })
                })
              ])
            } catch {
              return h(NButton, {
                size: 'tiny',
                type: 'info',
                quaternary: true,
                onClick: () => handleViewToolDetail(toolRow)
              }, {
                icon: () => h(NIcon, { component: EyeOutline }),
                default: () => t('mcpServer.viewParams')
              })
            }
          }
        }
      ]

      return h('div', { style: 'padding: 8px 16px' }, [
        h(NDataTable, {
          columns: toolColumns,
          data: tools,
          size: 'small',
          bordered: false,
          rowKey: (t) => t.id,
          pagination: false
        })
      ])
    }
  },
  {
    title: 'ID',
    key: 'id',
    width: 60,
    ellipsis: { tooltip: true }
  },
  {
    title: t('mcpServer.serverName'),
    key: 'name',
    width: 180,
    ellipsis: { tooltip: true }
  },
  {
    title: t('mcpServer.description'),
    key: 'description',
    width: 200,
    ellipsis: { tooltip: { style: { maxWidth: '400px', wordBreak: 'break-all' } } }
  },
  {
    title: 'URL',
    key: 'url',
    width: 180,
    ellipsis: { tooltip: true },
    render(row) {
      if (!row.url) return h(NText, { depth: 3 }, { default: () => '-' })
      return h(NText, { code: true, depth: 2 }, { default: () => row.url })
    }
  },
  {
    title: t('mcpServer.enable'),
    key: 'enable',
    width: 70,
    render(row) {
      return h(NTag, { type: row.enable ? 'success' : 'error' }, {
        default: () => (row.enable ? t('mcpServer.yes') : t('mcpServer.no'))
      })
    }
  },
  {
    title: t('mcpServer.status'),
    key: 'status',
    width: 80,
    render(row) {
      const typeMap = {
        available: 'success',
        untested: 'default',
        unavailable: 'error'
      }
      return h(NTag, { type: typeMap[row.status] || 'default' }, {
        default: () => getStatusLabel(row.status)
      })
    }
  },
  {
    title: t('mcpServer.testResult'),
    key: 'testResult',
    width: 200,
    ellipsis: { tooltip: { style: { maxWidth: '400px', wordBreak: 'break-all' } } },
    render(row) {
      if (!row.testResult) return h(NText, { depth: 3 }, { default: () => '-' })
      const isSuccess = row.status === 'available'
      return h(NText, { type: isSuccess ? 'success' : 'error' }, { default: () => row.testResult })
    }
  },
  {
    title: t('common.actions'),
    key: 'actions',
    width: 280,
    fixed: 'right',
    render(row) {
      return h(NSpace, {}, {
        default: () => [
          h(
            NButton,
            {
              size: 'tiny',
              type: 'info',
              onClick: () => handleTest(row)
            },
            {
              icon: () => h(NIcon, { component: FlashOutline }),
              default: () => t('mcpServer.test')
            }
          ),
          h(
            NButton,
            {
              size: 'tiny',
              type: row.enable ? 'warning' : 'info',
              onClick: () => handleToggleEnable(row)
            },
            {
              icon: () => h(NIcon, { component: row.enable ? PauseOutline : PlayOutline }),
              default: () => (row.enable ? t('mcpServer.disable') : t('mcpServer.enable'))
            }
          ),
          h(
            NButton,
            {
              size: 'tiny',
              type: 'primary',
              onClick: () => handleEdit(row)
            },
            {
              icon: () => h(NIcon, { component: CreateOutline }),
              default: () => t('common.edit')
            }
          ),
          h(
            NPopconfirm,
            {
              onPositiveClick: () => handleDelete(row.id)
            },
            {
              trigger: () =>
                h(
                  NButton,
                  {
                    size: 'tiny',
                    type: 'error'
                  },
                  {
                    icon: () => h(NIcon, { component: TrashOutline }),
                    default: () => t('common.delete')
                  }
                ),
              default: () => t('mcpServer.confirmDeleteServer', { name: row.name })
            }
          )
        ]
      })
    }
  }
]

const pagination = computed(() => ({
  page: currentPage.value,
  pageSize: pageSize.value,
  itemCount: total.value,
  pageCount: Math.ceil(total.value / pageSize.value) || 1,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  prefix: ({ itemCount }) => `${t('common.total')} ${itemCount} ${t('common.records')}`,
  onChange: handlePageChange,
  onUpdatePageSize: handlePageSizeChange
}))

const loadServerList = async () => {
  loading.value = true
  try {
    const query = {
      page: currentPage.value,
      pageSize: pageSize.value,
      name: searchKeyword.value,
      status: filterStatus.value
    }

    const result = await GetMCPServerList(query)
    if (result) {
      const servers = result.data || []
      let allTools = []
      try {
        allTools = await GetAllMCPTools() || []
      } catch (error) {
        console.error(t('mcpServer.loadToolsFailed'), error)
      }
      const toolsMap = {}
      for (const t of allTools) {
        if (!toolsMap[t.mcpServerId]) toolsMap[t.mcpServerId] = []
        toolsMap[t.mcpServerId].push(t)
      }
      for (const server of servers) {
        server.tools = toolsMap[server.id] || []
      }
      serverList.value = servers
      total.value = result.total || 0
    }
  } catch (error) {
    console.error(t('mcpServer.loadServerListFailed'), error)
    message.error(t('mcpServer.loadServerListFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = async () => {
  currentPage.value = 1
  await loadServerList()
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadServerList()
}

const handlePageSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadServerList()
}

const handleTest = async (row) => {
  try {
    const result = await TestMCPServer(row.id)
    message.success(result)
    await loadServerList()
  } catch (error) {
    message.error(t('mcpServer.testFailed') + error.message)
  }
}

const handleToggleEnable = async (row) => {
  try {
    const newEnable = !row.enable
    const result = await EnableMCPServer(row.id, newEnable)
    message.success(result)
    await loadServerList()
  } catch (error) {
    message.error(t('common.operation') + t('common.failed') + error.message)
  }
}

const handleCreate = () => {
  editingServer.value = false
  resetForm()
  showCreateModal.value = true
}

const handleEdit = async (row) => {
  editingServer.value = true
  try {
    const server = await GetMCPServerByID(row.id)
    if (server) {
      resetForm()
      formData.id = server.id
      formData.name = server.name
      formData.description = server.description
      formData.url = server.url
      formData.command = server.command
      formData.args = server.args
      formData.env = server.env
      formData.enable = server.enable
      formData.status = server.status
      showCreateModal.value = true
    }
  } catch (error) {
    message.error(t('mcpServer.getServerDetailFailed') + error.message)
  }
}

const handleDelete = async (id) => {
  try {
    const result = await DeleteMCPServer(id)
    message.success(result)
    await loadServerList()
  } catch (error) {
    message.error(t('common.delete') + t('common.failed') + error.message)
  }
}

const handleSubmit = async () => {
  try {
    if (formRef.value) {
      try {
        await formRef.value.validate()
      } catch {
        return
      }
    }

    submitting.value = true
    const submitData = { ...formData }

    let result
    if (formData.id) {
      result = await UpdateMCPServer(submitData)
    } else {
      result = await CreateMCPServer(submitData)
    }

    if (result.includes(t('common.success'))) {
      message.success(result)
      showCreateModal.value = false
      await loadServerList()
    } else {
      message.error(result)
    }
  } catch (error) {
    message.error(t('common.operation') + t('common.failed') + error.message)
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  Object.assign(formData, {
    id: null,
    name: '',
    description: '',
    url: '',
    command: '',
    args: '',
    env: '',
    enable: true,
    status: 'stopped'
  })
  if (formRef.value) {
    formRef.value.restoreValidation()
  }
}

const serverList = ref([])

onMounted(async () => {
  await loadServerList()
})
</script>

<style scoped>
</style>
