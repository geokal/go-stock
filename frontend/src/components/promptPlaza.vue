<script setup>
import {computed, h, onBeforeMount, onMounted, ref, reactive} from 'vue'
import {useI18n} from 'vue-i18n'
import {GetConfig} from "../../wailsjs/go/main/App";
import {useMessage, useDialog} from "naive-ui";
import {MdPreview, MdEditor} from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import 'md-editor-v3/lib/style.css'

const message = useMessage()
const dialog = useDialog()
const {t} = useI18n()

const darkTheme = ref(false)
const editorTheme = ref('light')
const apiBase = ref('http://go-stock.sparkmemory.top:1918/api')
const token = ref(localStorage.getItem('promptPlazaToken') || '')
const currentUser = ref(null)
const categories = ref([])
const activeCategory = ref(null)
const activeSort = ref('')
const keyword = ref('')
const loading = ref(false)
const prompts = ref([])
const pagination = reactive({
  page: 1,
  pageSize: 12,
  itemCount: 0,
  pageCount: 1
})

const detailModal = reactive({
  show: false,
  data: null,
  comments: [],
  commentPage: 1,
  commentPageSize: 10,
  commentTotal: 0,
  commentLoading: false,
  newComment: '',
  replyTo: null
})

const loginModal = reactive({
  show: false,
  tab: 'login',
  username: localStorage.getItem('promptPlazaUsername') || '',
  password: localStorage.getItem('promptPlazaPassword') || '',
  nickname: ''
})

const createModal = reactive({
  show: false,
  title: '',
  content: '',
  description: '',
  category: '',
  tags: '',
  isPublic: true
})

const rankingModal = reactive({
  show: false,
  type: 'hot',
  range: 'all',
  list: [],
  loading: false
})

const editModal = reactive({
  show: false,
  id: 0,
  title: '',
  content: '',
  description: '',
  category: '',
  tags: '',
  isPublic: true,
  loading: false
})

const isLoggedIn = computed(() => !!token.value)

onBeforeMount(() => {
  GetConfig().then(result => {
    if (result.darkTheme) {
      darkTheme.value = true
      editorTheme.value = 'dark'
    }
    if (result.promptPlazaApiBase) {
      apiBase.value = result.promptPlazaApiBase
    }
  })
})

onMounted(() => {
  loadCategories()
  loadPrompts()
  if (token.value) {
    fetchCurrentUser()
  }
})

function getHeaders() {
  const headers = {'Content-Type': 'application/json'}
  if (token.value) {
    headers['Authorization'] = `Bearer ${token.value}`
  }
  return headers
}

async function apiGet(path, params = {}) {
  const url = new URL(apiBase.value + path)
  Object.entries(params).forEach(([k, v]) => {
    if (v !== null && v !== undefined && v !== '') {
      url.searchParams.set(k, v)
    }
  })
  const resp = await fetch(url.toString(), {headers: getHeaders()})
  const json = await resp.json()
  if (json.code !== 0) {
    throw new Error(json.message || t('promptPlaza.loadFailed'))
  }
  return json.data
}

async function apiPost(path, body = null) {
  const resp = await fetch(apiBase.value + path, {
    method: 'POST',
    headers: getHeaders(),
    body: body ? JSON.stringify(body) : null
  })
  const json = await resp.json()
  if (json.code !== 0) {
    throw new Error(json.message || t('promptPlaza.loadFailed'))
  }
  return json.data
}

async function apiPut(path, body) {
  const resp = await fetch(apiBase.value + path, {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify(body)
  })
  const json = await resp.json()
  if (json.code !== 0) {
    throw new Error(json.message || t('promptPlaza.loadFailed'))
  }
  return json.data
}

async function apiDelete(path) {
  const resp = await fetch(apiBase.value + path, {
    method: 'DELETE',
    headers: getHeaders()
  })
  const json = await resp.json()
  if (json.code !== 0) {
    throw new Error(json.message || t('promptPlaza.loadFailed'))
  }
  return json.data
}

async function loadCategories() {
  try {
    const data = await apiGet('/prompts/categories')
    categories.value = data || []
  } catch (e) {
    console.warn(t('promptPlaza.loadCategoriesFailed'), e)
  }
}

async function loadPrompts() {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize
    }
    if (activeCategory.value) params.category = activeCategory.value
    if (keyword.value) params.keyword = keyword.value
    if (activeSort.value) params.sort = activeSort.value
    const data = await apiGet('/prompts', params)
    prompts.value = data.list || []
    pagination.itemCount = data.total || 0
    pagination.pageCount = data.totalPages || 1
  } catch (e) {
    message.error(t('promptPlaza.loadFailed') + e.message)
  } finally {
    loading.value = false
  }
}

async function fetchCurrentUser() {
  try {
    const data = await apiGet('/user/me')
    currentUser.value = data
  } catch (e) {
    token.value = ''
    localStorage.removeItem('promptPlazaToken')
    currentUser.value = null
  }
}

async function handleLogin() {
  try {
    const data = await apiPost('/auth/login', {
      username: loginModal.username,
      password: loginModal.password
    })
    token.value = data.token
    localStorage.setItem('promptPlazaToken', data.token)
    localStorage.setItem('promptPlazaUsername', loginModal.username)
    localStorage.setItem('promptPlazaPassword', loginModal.password)
    currentUser.value = data.user
    loginModal.show = false
    message.success(t('promptPlaza.loginSuccess'))
    loadPrompts()
  } catch (e) {
    message.error(t('promptPlaza.loginFailed') + e.message)
  }
}

async function handleRegister() {
  try {
    const data = await apiPost('/auth/register', {
      username: loginModal.username,
      password: loginModal.password,
      nickname: loginModal.nickname
    })
    token.value = data.token
    localStorage.setItem('promptPlazaToken', data.token)
    localStorage.setItem('promptPlazaUsername', loginModal.username)
    localStorage.setItem('promptPlazaPassword', loginModal.password)
    currentUser.value = data.user
    loginModal.show = false
    loginModal.username = ''
    loginModal.password = ''
    loginModal.nickname = ''
    message.success(t('promptPlaza.registerSuccess'))
    loadPrompts()
  } catch (e) {
    message.error(t('promptPlaza.registerFailed') + e.message)
  }
}

function handleLogout() {
  dialog.warning({
    title: t('promptPlaza.loginRequired'),
    content: t('promptPlaza.logoutConfirm'),
    positiveText: t('promptPlaza.confirm'),
    negativeText: t('promptPlaza.cancel'),
    onPositiveClick: () => {
      token.value = ''
      localStorage.removeItem('promptPlazaToken')
      currentUser.value = null
      message.success(t('promptPlaza.logoutSuccess'))
      loadPrompts()
    }
  })
}

function handlePageChange(page) {
  pagination.page = page
  loadPrompts()
}

function handleSearch() {
  pagination.page = 1
  loadPrompts()
}

function handleCategoryFilter(cat) {
  activeCategory.value = cat || null
  pagination.page = 1
  loadPrompts()
}

function handleSortChange(sort) {
  activeSort.value = sort
  pagination.page = 1
  loadPrompts()
}

async function showDetail(id) {
  try {
    const data = await apiGet(`/prompts/${id}`)
    detailModal.data = data
    detailModal.show = true
    detailModal.newComment = ''
    detailModal.replyTo = null
    loadComments(id)
  } catch (e) {
    message.error(t('promptPlaza.loadDetailFailed') + e.message)
  }
}

async function loadComments(promptId) {
  detailModal.commentLoading = true
  try {
    const data = await apiGet(`/prompts/${promptId}/comments`, {
      page: detailModal.commentPage,
      pageSize: detailModal.commentPageSize
    })
    detailModal.comments = data.list || []
    detailModal.commentTotal = data.total || 0
  } catch (e) {
    console.warn(t('promptPlaza.loadCommentsFailed'), e)
  } finally {
    detailModal.commentLoading = false
  }
}

async function handleLike(prompt) {
  if (!isLoggedIn.value) {
    message.warning(t('promptPlaza.loginRequired'))
    loginModal.show = true
    return
  }
  try {
    const data = await apiPost(`/prompts/${prompt.id}/like`)
    prompt.isLiked = data.isLiked
    prompt.likesCount = data.likesCount
    if (detailModal.data && detailModal.data.id === prompt.id) {
      detailModal.data.isLiked = data.isLiked
      detailModal.data.likesCount = data.likesCount
    }
  } catch (e) {
    message.error(t('promptPlaza.operationFailed') + e.message)
  }
}

async function handleFavorite(prompt) {
  if (!isLoggedIn.value) {
    message.warning(t('promptPlaza.loginRequired'))
    loginModal.show = true
    return
  }
  try {
    const data = await apiPost(`/prompts/${prompt.id}/favorite`)
    prompt.isFavorited = data.isFavorited
    prompt.favoritesCount = data.favoritesCount
    if (detailModal.data && detailModal.data.id === prompt.id) {
      detailModal.data.isFavorited = data.isFavorited
      detailModal.data.favoritesCount = data.favoritesCount
    }
  } catch (e) {
    message.error(t('promptPlaza.operationFailed') + e.message)
  }
}

async function handleDownload(prompt) {
  try {
    const data = await apiGet(`/prompts/${prompt.id}/download`)
    const text = `${data.title}\n\n${data.content}\n\n${t('promptPlaza.category')}: ${data.category || '-'}\n${t('promptPlaza.tags') || 'Tags'}: ${data.tags || '-'}\n${t('promptPlaza.author') || 'Author'}: ${data.author?.nickname || data.author?.username || t('promptPlaza.anonymous')}\n${t('promptPlaza.createdAt') || 'Created'}: ${data.createdAt}`
    if (navigator.clipboard) {
      await navigator.clipboard.writeText(data.content)
      message.success(t('promptPlaza.copiedToClipboard'))
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = data.content
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
      message.success(t('promptPlaza.copiedToClipboard'))
    }
    prompt.downloadsCount = (prompt.downloadsCount || 0) + 1
  } catch (e) {
    message.error(t('promptPlaza.downloadFailed') + e.message)
  }
}

async function handleCopyContent(content) {
  try {
    if (navigator.clipboard) {
      await navigator.clipboard.writeText(content)
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = content
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    }
    message.success(t('promptPlaza.copiedToClipboard'))
  } catch (e) {
    message.error(t('promptPlaza.copyFailed'))
  }
}

async function submitComment() {
  if (!isLoggedIn.value) {
    message.warning(t('promptPlaza.loginRequired'))
    loginModal.show = true
    return
  }
  if (!detailModal.newComment.trim()) {
    message.warning(t('promptPlaza.commentContentRequired'))
    return
  }
  try {
    const body = {content: detailModal.newComment}
    if (detailModal.replyTo) {
      body.parentId = detailModal.replyTo.id
    }
    await apiPost(`/prompts/${detailModal.data.id}/comments`, body)
    detailModal.newComment = ''
    detailModal.replyTo = null
    detailModal.data.commentsCount = (detailModal.data.commentsCount || 0) + 1
    loadComments(detailModal.data.id)
    message.success(t('promptPlaza.commentSuccess'))
  } catch (e) {
    message.error(t('promptPlaza.commentFailed') + e.message)
  }
}

async function deleteComment(commentId) {
  dialog.warning({
    title: t('promptPlaza.loginRequired'),
    content: t('promptPlaza.deleteCommentConfirm'),
    positiveText: t('promptPlaza.confirm'),
    negativeText: t('promptPlaza.cancel'),
    onPositiveClick: async () => {
      try {
        await apiDelete(`/comments/${commentId}`)
        detailModal.data.commentsCount = Math.max(0, (detailModal.data.commentsCount || 1) - 1)
        loadComments(detailModal.data.id)
        message.success(t('promptPlaza.deleteSuccess'))
      } catch (e) {
        message.error(t('promptPlaza.deleteFailed') + e.message)
      }
    }
  })
}

function showEditModal(prompt) {
  editModal.id = prompt.id
  editModal.title = prompt.title || ''
  editModal.content = prompt.content || ''
  editModal.description = prompt.description || ''
  editModal.category = prompt.category || ''
  editModal.tags = prompt.tags || ''
  editModal.isPublic = prompt.isPublic !== false
  editModal.show = true
}

async function handleEdit() {
  if (!editModal.title || !editModal.content) {
    message.warning(t('promptPlaza.titleContentRequired'))
    return
  }
  editModal.loading = true
  try {
    await apiPut(`/prompts/${editModal.id}`, {
      title: editModal.title,
      content: editModal.content,
      description: editModal.description,
      category: editModal.category,
      tags: editModal.tags,
      isPublic: editModal.isPublic
    })
    editModal.show = false
    detailModal.show = false
    message.success(t('promptPlaza.modifySuccess'))
    loadPrompts()
    loadCategories()
  } catch (e) {
    message.error(t('promptPlaza.modifyFailed') + e.message)
  } finally {
    editModal.loading = false
  }
}

function handleDeletePrompt(prompt) {
  dialog.warning({
    title: t('promptPlaza.loginRequired'),
    content: t('promptPlaza.deletePromptConfirm'),
    positiveText: t('promptPlaza.confirm'),
    negativeText: t('promptPlaza.cancel'),
    onPositiveClick: async () => {
      try {
        await apiDelete(`/prompts/${prompt.id}`)
        detailModal.show = false
        message.success(t('promptPlaza.deleteSuccess'))
        loadPrompts()
        loadCategories()
      } catch (e) {
        message.error(t('promptPlaza.deleteFailed') + e.message)
      }
    }
  })
}

async function showCreateModal() {
  if (!isLoggedIn.value) {
    message.warning(t('promptPlaza.loginRequired'))
    loginModal.show = true
    return
  }
  createModal.title = ''
  createModal.content = ''
  createModal.description = ''
  createModal.category = ''
  createModal.tags = ''
  createModal.isPublic = true
  createModal.show = true
}

async function handleCreate() {
  if (!createModal.title || !createModal.content) {
    message.warning(t('promptPlaza.titleContentRequired'))
    return
  }
  try {
    await apiPost('/prompts', {
      title: createModal.title,
      content: createModal.content,
      description: createModal.description,
      category: createModal.category,
      tags: createModal.tags,
      isPublic: createModal.isPublic
    })
    createModal.show = false
    message.success(t('promptPlaza.publishSuccess'))
    loadPrompts()
    loadCategories()
  } catch (e) {
    message.error(t('promptPlaza.publishFailed') + e.message)
  }
}

async function showRanking(type = 'hot', range = 'all') {
  rankingModal.type = type
  rankingModal.range = range
  rankingModal.show = true
  rankingModal.loading = true
  try {
    const data = await apiGet('/prompts/ranking', {type, range, limit: 50})
    rankingModal.list = data.list || []
  } catch (e) {
    message.error(t('promptPlaza.loadRankingFailed') + e.message)
  } finally {
    rankingModal.loading = false
  }
}

function formatTime(timeStr) {
  if (!timeStr) return ''
  return timeStr.substring(0, 19).replace('T', ' ')
}

function timeAgo(timeStr) {
  if (!timeStr) return ''
  const now = new Date()
  const time = new Date(timeStr)
  const diff = Math.floor((now - time) / 1000)
  if (diff < 60) return t('promptPlaza.justNow')
  if (diff < 3600) return t('promptPlaza.minutesAgo', { n: Math.floor(diff / 60) })
  if (diff < 86400) return t('promptPlaza.hoursAgo', { n: Math.floor(diff / 3600) })
  if (diff < 2592000) return t('promptPlaza.daysAgo', { n: Math.floor(diff / 86400) })
  return formatTime(timeStr)
}
</script>

<template>
  <div style="padding: 0">
    <n-space vertical :size="12">
      <n-space justify="space-between" align="center">
        <n-space align="center">
          <n-input
            v-model:value="keyword"
            :placeholder="t('promptPlaza.searchPlaceholder')"
            clearable
            style="width: 260px"
            @keyup.enter="handleSearch"
          />
          <n-button type="primary" @click="handleSearch">{{ t('promptPlaza.search') }}</n-button>
          <n-button quaternary @click="showRanking('hot')">🏆 {{ t('promptPlaza.ranking') }}</n-button>
        </n-space>
        <n-space>
          <n-button type="success" @click="showCreateModal">✏️ {{ t('promptPlaza.publishPrompt') }}</n-button>
          <template v-if="isLoggedIn">
            <n-tag type="success" size="medium" round>
              {{ currentUser?.nickname || currentUser?.username || t('promptPlaza.loggedIn') }}
            </n-tag>
            <n-button size="small" quaternary @click="handleLogout">{{ t('promptPlaza.logout') }}</n-button>
          </template>
          <template v-else>
            <n-button type="info" size="small" @click="loginModal.show = true; loginModal.tab = 'login'">{{ t('promptPlaza.loginRegister') }}</n-button>
          </template>
        </n-space>
      </n-space>

      <n-space align="center" :size="8">
        <n-text depth="3" style="font-size: 13px">{{ t('promptPlaza.category') }}</n-text>
        <n-button
          :type="activeCategory === null ? 'primary' : 'default'"
          size="small"
          @click="handleCategoryFilter(null)"
        >{{ t('promptPlaza.all') }}</n-button>
        <n-button
          v-for="cat in categories"
          :key="cat"
          :type="activeCategory === cat ? 'primary' : 'default'"
          size="small"
          @click="handleCategoryFilter(cat)"
        >{{ cat }}</n-button>
        <n-divider vertical />
        <n-text depth="3" style="font-size: 13px">{{ t('promptPlaza.sort') }}</n-text>
        <n-button :type="activeSort === '' ? 'primary' : 'default'" size="small" @click="handleSortChange('')">{{ t('promptPlaza.latest') }}</n-button>
        <n-button :type="activeSort === 'likes' ? 'primary' : 'default'" size="small" @click="handleSortChange('likes')">{{ t('promptPlaza.hottest') }}</n-button>
        <n-button :type="activeSort === 'downloads' ? 'primary' : 'default'" size="small" @click="handleSortChange('downloads')">{{ t('promptPlaza.downloads') }}</n-button>
        <n-button :type="activeSort === 'comments' ? 'primary' : 'default'" size="small" @click="handleSortChange('comments')">{{ t('promptPlaza.comments') }}</n-button>
      </n-space>

      <n-spin :show="loading">
        <n-grid :cols="3" :x-gap="12" :y-gap="12" responsive="screen">
          <n-gi v-for="item in prompts" :key="item.id">
            <n-card
              hoverable
              size="small"
              style="cursor: pointer; height: 100%"
              @click="showDetail(item.id)"
            >
              <template #header>
                <n-text strong style="font-size: 15px">{{ item.title }}</n-text>
              </template>
              <template #header-extra>
                <n-tag v-if="item.category" size="small" type="info">{{ item.category }}</n-tag>
              </template>
              <n-ellipsis :line-clamp="2" :tooltip="false" style="color: var(--n-text-color-3); font-size: 13px; margin-bottom: 8px">
                {{ item.description || item.content }}
              </n-ellipsis>
              <template #footer>
                <n-space justify="space-between" align="center">
                  <n-text depth="3" style="font-size: 12px">
                    {{ item.user?.nickname || item.user?.username || t('promptPlaza.anonymous') }} · {{ timeAgo(item.createdAt) }}
                  </n-text>
                  <n-space :size="12" style="font-size: 12px">
                    <n-text :type="item.isLiked ? 'error' : 'default'" style="cursor: pointer" @click.stop="handleLike(item)">
                      {{ item.isLiked ? '❤️' : '🤍' }} {{ item.likesCount || 0 }}
                    </n-text>
                    <n-text :type="item.isFavorited ? 'warning' : 'default'" style="cursor: pointer" @click.stop="handleFavorite(item)">
                      {{ item.isFavorited ? '⭐' : '☆' }} {{ item.favoritesCount || 0 }}
                    </n-text>
                    <n-text depth="3">
                      💬 {{ item.commentsCount || 0 }}
                    </n-text>
                    <n-text depth="3">
                      ⬇️ {{ item.downloadsCount || 0 }}
                    </n-text>
                  </n-space>
                </n-space>
              </template>
              <template #action v-if="item.tags">
                <n-space :size="4">
                  <n-tag v-for="tag in item.tags.split(',').filter(t=>t).slice(0, 3)" :key="tag" size="tiny" round>{{ tag.trim() }}</n-tag>
                </n-space>
              </template>
            </n-card>
          </n-gi>
        </n-grid>
        <n-empty v-if="!loading && prompts.length === 0" :description="t('promptPlaza.noPrompts')" style="margin-top: 40px" />
      </n-spin>

      <n-space justify="center" style="margin-top: 12px" v-if="pagination.pageCount > 1">
        <n-pagination
          v-model:page="pagination.page"
          :page-count="pagination.pageCount"
          :page-size="pagination.pageSize"
          @update:page="handlePageChange"
        />
      </n-space>
    </n-space>

    <n-modal v-model:show="detailModal.show" preset="card" style="width: 1100px; max-width: 95vw" :title="detailModal.data?.title || t('promptPlaza.promptDetail')">
      <template v-if="detailModal.data">
        <n-space align="left" justify="space-between" style="margin-bottom: 12px">
          <n-space align="left" :size="8">
            <n-tag v-if="detailModal.data.category" type="info" size="small">{{ detailModal.data.category }}</n-tag>
            <n-text depth="3" style="font-size: 12px">
              {{ detailModal.data.user?.nickname || detailModal.data.user?.username || t('promptPlaza.anonymous') }} · {{ formatTime(detailModal.data.createdAt) }}
            </n-text>
            <n-text depth="3" style="font-size: 12px" v-if="detailModal.data.updatedAt && detailModal.data.updatedAt !== detailModal.data.createdAt">
              · {{ t('promptPlaza.updatedAt') }} {{ formatTime(detailModal.data.updatedAt) }}
            </n-text>
          </n-space>
          <n-space :size="8">
            <n-button
              v-if="currentUser && detailModal.data.userId === currentUser.id"
              size="tiny"
              type="warning"
              @click="showEditModal(detailModal.data)"
            >
              ✏️ {{ t('promptPlaza.edit') }}
            </n-button>
            <n-button
              v-if="currentUser && detailModal.data.userId === currentUser.id"
              size="tiny"
              type="error"
              @click="handleDeletePrompt(detailModal.data)"
            >
              🗑️ {{ t('promptPlaza.delete') }}
            </n-button>
            <n-button
              :type="detailModal.data.isLiked ? 'error' : 'default'"
              size="tiny"
              @click="handleLike(detailModal.data)"
            >
              {{ detailModal.data.isLiked ? '❤️ ' + t('promptPlaza.liked') : '🤍 ' + t('promptPlaza.like') }} {{ detailModal.data.likesCount || 0 }}
            </n-button>
            <n-button
              :type="detailModal.data.isFavorited ? 'warning' : 'default'"
              size="tiny"
              @click="handleFavorite(detailModal.data)"
            >
              {{ detailModal.data.isFavorited ? '⭐ ' + t('promptPlaza.favorited') : '☆ ' + t('promptPlaza.favorite') }} {{ detailModal.data.favoritesCount || 0 }}
            </n-button>
            <n-button size="tiny" type="success" @click="handleDownload(detailModal.data)">
              ⬇️ {{ t('promptPlaza.download') }} {{ detailModal.data.downloadsCount || 0 }}
            </n-button>
            <n-button size="tiny" quaternary @click="handleCopyContent(detailModal.data.content)">
              📋 {{ t('promptPlaza.copy') }}
            </n-button>
          </n-space>
        </n-space>

        <div style="display: flex; gap: 16px">
          <div style="flex: 4; min-width: 0">
            <n-space vertical :size="8">
              <n-space :size="4" v-if="detailModal.data.tags">
                <n-tag v-for="tag in detailModal.data.tags.split(',').filter(t=>t)" :key="tag" size="small" round>{{ tag.trim() }}</n-tag>
              </n-space>
              <div style="max-height: 500px; overflow-y: auto">
                <MdPreview
                  :model-value="detailModal.data.content"
                  :theme="editorTheme"
                  style="text-align: left"
                />
              </div>
            </n-space>
          </div>
          <div style="flex: 1; min-width: 0">
            <n-space vertical :size="8" style="width: 100%">
              <n-text strong>{{ t('promptPlaza.commentsTitle', { count: detailModal.data.commentsCount || 0 }) }}</n-text>
              <n-input
                v-model:value="detailModal.newComment"
                type="textarea"
                :placeholder="detailModal.replyTo ? t('promptPlaza.replyTo', { username: detailModal.replyTo.user?.nickname || detailModal.replyTo.user?.username }) : t('promptPlaza.commentPlaceholder')"
                :rows="2"
              />
              <n-space justify="space-between" style="width: 100%">
                <n-text v-if="detailModal.replyTo" depth="3" style="font-size: 12px">
                  {{ t('promptPlaza.replyTo', { username: detailModal.replyTo.user?.nickname || detailModal.replyTo.user?.username }) }}
                  <n-button text size="tiny" type="error" @click="detailModal.replyTo = null">{{ t('promptPlaza.cancel') }}</n-button>
                </n-text>
                <span v-else />
                <n-button size="small" type="primary" @click="submitComment">{{ t('promptPlaza.submitComment') }}</n-button>
              </n-space>
              <n-spin :show="detailModal.commentLoading">
                <div style="max-height: 380px; overflow-y: auto; width: 100%">
                  <n-space vertical :size="12" style="width: 100%">
                    <n-card v-for="comment in detailModal.comments" :key="comment.id" size="small" embedded>
                      <template #header>
                        <n-space align="center" :size="8">
                          <n-text strong style="font-size: 13px">{{ comment.user?.nickname || comment.user?.username }}</n-text>
                          <n-text depth="3" style="font-size: 12px">{{ timeAgo(comment.createdAt) }}</n-text>
                        </n-space>
                      </template>
                      <n-text style="font-size: 13px; text-align: left; display: block">{{ comment.content }}</n-text>
                      <template #action>
                        <n-space :size="8">
                          <n-button text size="tiny" @click="detailModal.replyTo = comment">{{ t('promptPlaza.reply') }}</n-button>
                          <n-button
                            v-if="currentUser && comment.userId === currentUser.id"
                            text
                            size="tiny"
                            type="error"
                            @click="deleteComment(comment.id)"
                          >{{ t('promptPlaza.delete') }}</n-button>
                        </n-space>
                      </template>
                    </n-card>
                    <n-empty v-if="!detailModal.commentLoading && detailModal.comments.length === 0" :description="t('promptPlaza.noComments')" size="small" />
                  </n-space>
                </div>
              </n-spin>
            </n-space>
          </div>
        </div>
      </template>
    </n-modal>

    <n-modal v-model:show="loginModal.show" preset="card" style="width: 400px" :title="t('promptPlaza.account')">
      <n-tabs v-model:value="loginModal.tab" type="line">
        <n-tab-pane name="login" :tab="t('promptPlaza.login')">
          <n-space vertical :size="12">
            <n-input v-model:value="loginModal.username" :placeholder="t('promptPlaza.username')" />
            <n-input v-model:value="loginModal.password" type="password" :placeholder="t('promptPlaza.password')" show-password-on="click" />
            <n-button type="primary" block @click="handleLogin">{{ t('promptPlaza.login') }}</n-button>
          </n-space>
        </n-tab-pane>
        <n-tab-pane name="register" :tab="t('promptPlaza.register')">
          <n-space vertical :size="12">
            <n-input v-model:value="loginModal.username" :placeholder="t('promptPlaza.usernamePlaceholder')" />
            <n-input v-model:value="loginModal.password" type="password" :placeholder="t('promptPlaza.passwordPlaceholder')" show-password-on="click" />
            <n-input v-model:value="loginModal.nickname" :placeholder="t('promptPlaza.nicknamePlaceholder')" />
            <n-button type="primary" block @click="handleRegister">{{ t('promptPlaza.register') }}</n-button>
          </n-space>
        </n-tab-pane>
      </n-tabs>
    </n-modal>

    <n-modal v-model:show="createModal.show" preset="card" style="width: 1100px; max-width: 95vw" :title="t('promptPlaza.publishPrompt')">
      <n-space vertical :size="12">
        <n-input v-model:value="createModal.title" :placeholder="t('promptPlaza.title')" />
        <n-space :size="8">
          <n-input v-model:value="createModal.category" :placeholder="t('promptPlaza.categoryPlaceholder')" style="width: 240px" />
          <n-input v-model:value="createModal.tags" :placeholder="t('promptPlaza.tagsPlaceholder')" style="width: 240px" />
        </n-space>
        <n-input v-model:value="createModal.description" :placeholder="t('promptPlaza.descriptionPlaceholder')" type="textarea" :rows="2" />
        <MdEditor
          v-model="createModal.content"
          :theme="editorTheme"
          :placeholder="t('promptPlaza.contentPlaceholder')"
          style="height: 400px"
        />
        <n-space justify="end">
          <n-button @click="createModal.show = false">{{ t('promptPlaza.cancel') }}</n-button>
          <n-button type="primary" @click="handleCreate">{{ t('promptPlaza.publish') }}</n-button>
        </n-space>
      </n-space>
    </n-modal>

    <n-modal v-model:show="editModal.show" preset="card" style="width: 1100px; max-width: 95vw" :title="t('promptPlaza.edit')">
      <n-space vertical :size="12">
        <n-input v-model:value="editModal.title" :placeholder="t('promptPlaza.title')" />
        <n-space :size="8">
          <n-input v-model:value="editModal.category" :placeholder="t('promptPlaza.categoryPlaceholder')" style="width: 240px" />
          <n-input v-model:value="editModal.tags" :placeholder="t('promptPlaza.tagsPlaceholder')" style="width: 240px" />
        </n-space>
        <n-input v-model:value="editModal.description" :placeholder="t('promptPlaza.descriptionPlaceholder')" type="textarea" :rows="2" />
        <MdEditor
          v-model="editModal.content"
          :theme="editorTheme"
          :placeholder="t('promptPlaza.contentPlaceholder')"
          style="height: 400px"
        />
        <n-space align="center">
          <n-text>{{ t('promptPlaza.public') }}</n-text>
          <n-switch v-model:value="editModal.isPublic" />
        </n-space>
        <n-space justify="end">
          <n-button @click="editModal.show = false">{{ t('promptPlaza.cancel') }}</n-button>
          <n-button type="primary" :loading="editModal.loading" @click="handleEdit">{{ t('promptPlaza.save') }}</n-button>
        </n-space>
      </n-space>
    </n-modal>

    <n-modal v-model:show="rankingModal.show" preset="card" style="width: 700px; max-width: 95vw" :title="t('promptPlaza.rankingTitle')">
      <n-space vertical :size="12">
        <n-space :size="8">
          <n-text depth="3" style="font-size: 13px">{{ t('promptPlaza.type') }}</n-text>
          <n-button :type="rankingModal.type === 'hot' ? 'primary' : 'default'" size="small" @click="showRanking('hot', rankingModal.range)">🔥 {{ t('promptPlaza.comprehensiveHeat') }}</n-button>
          <n-button :type="rankingModal.type === 'likes' ? 'primary' : 'default'" size="small" @click="showRanking('likes', rankingModal.range)">❤️ {{ t('promptPlaza.hottest') }}</n-button>
          <n-button :type="rankingModal.type === 'downloads' ? 'primary' : 'default'" size="small" @click="showRanking('downloads', rankingModal.range)">⬇️ {{ t('promptPlaza.downloads') }}</n-button>
          <n-button :type="rankingModal.type === 'favorites' ? 'primary' : 'default'" size="small" @click="showRanking('favorites', rankingModal.range)">⭐ {{ t('promptPlaza.favorite') }}</n-button>
          <n-divider vertical />
          <n-text depth="3" style="font-size: 13px">{{ t('promptPlaza.time') }}</n-text>
          <n-button :type="rankingModal.range === 'all' ? 'primary' : 'default'" size="small" @click="showRanking(rankingModal.type, 'all')">{{ t('promptPlaza.all') }}</n-button>
          <n-button :type="rankingModal.range === 'daily' ? 'primary' : 'default'" size="small" @click="showRanking(rankingModal.type, 'daily')">{{ t('promptPlaza.today') }}</n-button>
          <n-button :type="rankingModal.range === 'weekly' ? 'primary' : 'default'" size="small" @click="showRanking(rankingModal.type, 'weekly')">{{ t('promptPlaza.thisWeek') }}</n-button>
          <n-button :type="rankingModal.range === 'monthly' ? 'primary' : 'default'" size="small" @click="showRanking(rankingModal.type, 'monthly')">{{ t('promptPlaza.thisMonth') }}</n-button>
        </n-space>

        <n-spin :show="rankingModal.loading">
          <n-list bordered>
            <n-list-item v-for="item in rankingModal.list" :key="item.id" style="cursor: pointer" @click="rankingModal.show = false; showDetail(item.id)">
              <n-space align="center" :size="12">
                <n-tag
                  :type="item.rank <= 3 ? 'error' : 'default'"
                  round
                  size="small"
                  style="min-width: 28px; text-align: center"
                >{{ item.rank }}</n-tag>
                <n-text strong>{{ item.title }}</n-text>
                <n-text depth="3" style="font-size: 12px">
                  {{ item.user?.nickname || item.user?.username || t('promptPlaza.anonymous') }}
                </n-text>
                <n-space :size="8" style="font-size: 12px">
                  <n-text depth="3">❤️ {{ item.likesCount || 0 }}</n-text>
                  <n-text depth="3">⬇️ {{ item.downloadsCount || 0 }}</n-text>
                  <n-text depth="3">⭐ {{ item.favoritesCount || 0 }}</n-text>
                  <n-text depth="3">💬 {{ item.commentsCount || 0 }}</n-text>
                  <n-text v-if="item.hotScore" type="warning" style="font-size: 12px">🔥 {{ item.hotScore }}</n-text>
                </n-space>
              </n-space>
            </n-list-item>
          </n-list>
          <n-empty v-if="!rankingModal.loading && rankingModal.list.length === 0" :description="t('promptPlaza.noRankingData')" />
        </n-spin>
      </n-space>
    </n-modal>
  </div>
</template>

<style scoped>
:deep(.md-editor-preview) {
  padding: 8px 12px;
}
:deep(.md-editor-preview-wrapper) {
  padding: 0;
}
:deep(.md-editor-preview p),
:deep(.md-editor-preview h1),
:deep(.md-editor-preview h2),
:deep(.md-editor-preview h3),
:deep(.md-editor-preview h4),
:deep(.md-editor-preview h5),
:deep(.md-editor-preview h6),
:deep(.md-editor-preview ul),
:deep(.md-editor-preview ol),
:deep(.md-editor-preview blockquote),
:deep(.md-editor-preview pre),
:deep(.md-editor-preview div),
:deep(.md-editor-content p),
:deep(.md-editor-content h1),
:deep(.md-editor-content h2),
:deep(.md-editor-content h3),
:deep(.md-editor-content h4),
:deep(.md-editor-content h5),
:deep(.md-editor-content h6),
:deep(.md-editor-content ul),
:deep(.md-editor-content ol),
:deep(.md-editor-content blockquote),
:deep(.md-editor-content pre),
:deep(.md-editor-content div) {
  text-align: left;
}
</style>
