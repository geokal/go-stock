<script setup>
import { useI18n } from 'vue-i18n'
import { MdPreview } from 'md-editor-v3';
import 'md-editor-v3/lib/preview.css';
import {h, computed, nextTick, onBeforeUnmount, onMounted, ref} from 'vue';
import {CheckUpdate, GetConfig, GetVersionInfo,GetSponsorInfo,GetUserManual,OpenURL} from "../../wailsjs/go/main/App";
import {EventsOff, EventsOn,Environment} from "../../wailsjs/runtime";
import {NAvatar, NButton, NTree, useNotification,NText} from "naive-ui";
import { addMonths, format ,parse} from 'date-fns';
import { zhCN } from 'date-fns/locale';
const { t } = useI18n()
const updateLog = ref('');
const versionInfo = ref('');
const icon = ref('https://raw.githubusercontent.com/ArvinLovegood/go-stock/master/build/appicon.png');
const alipay =ref('https://github.com/ArvinLovegood/go-stock/raw/master/build/screenshot/alipay.jpg')
const wxpay =ref('https://github.com/ArvinLovegood/go-stock/raw/master/build/screenshot/wxpay.jpg')
const wxgzh =ref('https://github.com/ArvinLovegood/go-stock/raw/dev/build/screenshot/%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E7%99%BD%E8%89%B2%E7%89%88.png')
const notify = useNotification()
const vipLevel=ref("");
const vipStartTime=ref("");
const vipEndTime=ref("");
const expired=ref(false)
const showManual = ref(false)
const manualContent = ref('')
const manualId = 'manual-preview'
const darkTheme = ref(false)
const theme = computed(() => darkTheme.value ? 'dark' : 'light')
const manualScrollRef = ref(null)
const catalogList = ref([])

const buildCatalogTree = (headings) => {
  if (!headings.length) return []
  const roots = []
  const stack = []
  for (const h of headings) {
    const node = { key: h.text, label: h.text, level: h.level, children: [] }
    while (stack.length && stack[stack.length - 1].level >= h.level) {
      stack.pop()
    }
    if (stack.length) {
      stack[stack.length - 1].children.push(node)
    } else {
      roots.push(node)
    }
    stack.push(node)
  }
  const clean = (nodes) => {
    for (const n of nodes) {
      if (n.children.length === 0) delete n.children
      else clean(n.children)
    }
  }
  clean(roots)
  return roots
}

const catalogTree = computed(() => buildCatalogTree(catalogList.value))

const onTreeSelect = (keys) => {
  if (keys.length) scrollToHeading(keys[0])
}

const slugifyHeading = (text) => {
  return text
    .trim()
    .replace(/[^\w\u4e00-\u9fff]/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
}

const extractCatalog = () => {
  if (!manualScrollRef.value) return
  const headings = manualScrollRef.value.querySelectorAll('h1, h2, h3, h4, h5, h6')
  catalogList.value = Array.from(headings).map(h => ({
    text: h.textContent?.trim() || '',
    level: parseInt(h.tagName.slice(1))
  }))
}

const scrollToHeading = (headingText) => {
  if (!manualScrollRef.value) return
  const container = manualScrollRef.value
  const headings = container.querySelectorAll('h1, h2, h3, h4, h5, h6')
  for (const h of headings) {
    const text = h.textContent?.trim()
    if (text === headingText) {
      const containerRect = container.getBoundingClientRect()
      const headingRect = h.getBoundingClientRect()
      container.scrollTop += headingRect.top - containerRect.top - 10
      return
    }
  }
}

const openManual = () => {
  if (!manualContent.value) {
    GetUserManual().then(res => {
      manualContent.value = res
      showManual.value = true
      nextTick(() => { setTimeout(extractCatalog, 500) })
    })
  } else {
    showManual.value = true
    nextTick(() => { setTimeout(extractCatalog, 300) })
  }
}

onMounted(() => {
  document.title = t('about.aboutSoftware');
  GetConfig().then(res => {
    darkTheme.value = res.darkTheme
  })
  GetVersionInfo().then((res) => {
    updateLog.value = res.content;
    versionInfo.value = res.version;
    icon.value = res.icon;
    alipay.value=res.alipay;
    wxpay.value=res.wxpay;
    wxgzh.value=res.wxgzh;

    GetSponsorInfo().then((res) => {
      vipLevel.value = res.vipLevel;
      vipStartTime.value = res.vipStartTime;
      vipEndTime.value = res.vipEndTime;
      if (res.vipLevel) {
        if (res.vipEndTime < format(new Date(), 'yyyy-MM-dd HH:mm:ss')) {
          notify.warning({content: t('about.vipExpired')})
          expired.value = true;
        }
      }
    })

  });



})
onBeforeUnmount(() => {
  notify.destroyAll()
  EventsOff("updateVersion")
})

EventsOn("updateVersion",async (msg) => {
  const githubTimeStr = msg.published_at;
  const utcDate = new Date(githubTimeStr);
  const date = new Date(utcDate.getTime());
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  const formattedDate = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;

  notify.info({
    avatar: () =>
        h(NAvatar, {
          size: 'small',
          round: false,
          src: icon.value
        }),
    title: t('about.newVersionFound') + ': ' + msg.tag_name,
    content: () => {
      return h('div', {
        style: {
          'text-align': 'left',
          'font-size': '14px',
        }
      }, { default: () => msg.commit?.message })
    },
    duration: 5000,
    meta: t('about.publishTime') + ":"+formattedDate,
    action: () => {
      return h(NButton, {
        type: 'primary',
        size: 'small',
        onClick: () => {
          Environment().then(env => {
            switch (env.platform) {
              case 'windows':
                window.open(msg.html_url)
                break
              default :
                OpenURL(msg.html_url)
                break
            }
          })
        }
      }, { default: () => 'View' })
    }
  })
})

</script>

<template>
      <n-space vertical size="large"  style="--wails-draggable:no-drag">
        <n-card size="large">
          <n-divider title-placement="center">{{ t('about.aboutSoftware') }}</n-divider>
          <n-space vertical >
            <n-image width="100" :src="icon" />
            <h1>
              <n-badge v-if="!vipLevel"  :value="versionInfo" :offset="[80,10]"  type="success">
                <n-gradient-text type="info" :size="50" >go-stock</n-gradient-text>
              </n-badge>
              <n-badge v-if="vipLevel"  :value="versionInfo" :offset="[70,10]"  type="success">
                <n-gradient-text :type="expired?'error':'warning'" :size="50" >go-stock</n-gradient-text><n-tag :bordered="false" size="small" type="warning">VIP{{vipLevel}}</n-tag>
              </n-badge>
            </h1>
            <n-gradient-text  :type="expired?'error':'warning'" v-if="vipLevel" >{{ t('about.vipExpirationTime') }}：{{vipEndTime}}</n-gradient-text>
            <n-flex justify="center">
              <n-button size="tiny" @click="CheckUpdate(1)"  type="info" tertiary >{{ t('about.checkUpdate') }}</n-button>
              <n-button size="tiny" @click="openManual" type="success" tertiary >{{ t('about.viewUserManual') }}</n-button>
            </n-flex>
            <div style="justify-self: center;text-align: left" >
              <p>{{ t('about.description') }}</p>
              <p>{{ t('about.supportedStocks') }}</p>
              <p>{{ t('about.supportedPlatforms') }}</p>
              <p>
                <i style="color: crimson">{{ t('about.disclaimer') }}</i>
              </p>
              <p>
                {{ t('about.welcomeStar') }}：<a href="https://github.com/ArvinLovegood/go-stock" target="_blank">go-stock</a><n-divider vertical />
                <a href="https://github.com/ArvinLovegood/go-stock" target="_blank">GitHub</a><n-divider vertical />
                <a href="https://github.com/ArvinLovegood/go-stock/issues" target="_blank">{{ t('about.issues') }}</a><n-divider vertical />
                <a href="https://github.com/ArvinLovegood/go-stock/releases" target="_blank">{{ t('about.releases') }}</a><n-divider vertical />
              </p>
              <p v-if="updateLog">{{ t('about.updateLog') }}：{{updateLog}}</p>
              <p>{{ t('about.community') }}：<a href="https://go-stock.sparkmemory.top/" target="_blank">https://go-stock.sparkmemory.top/</a></p>
              <p>{{ t('about.qqGroup') }}：<a href="http://qm.qq.com/cgi-bin/qm/qr?_wv=1027&k=0YQ8qD3exahsD4YLNhzQTWe5ssstWC89&authKey=usOMMRFtIQDC%2FYcatHYapcxQbJ7PwXPHK9OypTXWzNjAq%2FRVvQu9bj2lRgb%2BSZ3p&noverify=0&group_code=491605333" target="_blank">491605333</a></p>
            </div>
          </n-space>
          <n-divider title-placement="center">{{ t('about.supportOpensource') }}</n-divider>
          <n-flex justify="center">
            <n-table  size="small" style="width: 820px">
              <n-thead>
                <n-tr>
                  <n-th>{{ t('about.sponsorPlan') }}</n-th>
                  <n-th>{{ t('about.sponsorLevel') }}</n-th>
                  <n-th>{{ t('about.benefits') }}</n-th>
                </n-tr>
              </n-thead>
              <n-tbody>
                <n-tr>
                  <n-td>0 RMB/month</n-td><n-td>vip0</n-td><n-td>All features, auto-update (from GitHub), resolve GitHub network issues yourself.</n-td>
                </n-tr>
                <n-tr>
                  <n-td>18.8 RMB/month<br>120 RMB/year</n-td><n-td>vip1</n-td><n-td>All features, auto-update (from CDN), fast updates. AI config guidance, prompt reference, etc.</n-td>
                </n-tr>
                <n-tr>
                  <n-td>28.8 RMB/month<br>240 RMB/year</n-td><n-td>vip2</n-td><n-td>vip1 all features, sync recent 24h market info on startup (including foreign media briefs), go-stock AI assistant (contact author WeChat/QQ for details)</n-td>
                </n-tr>
                <n-tr>
                  <n-td>X RMB/month</n-td><n-td>vipX</n-td><n-td>More plans based on go-stock open source project development... (GitHub README ad promotion)</n-td>
                </n-tr>
              </n-tbody>
            </n-table>
          </n-flex>
          <n-divider title-placement="center">{{ t('about.aboutAuthor') }}</n-divider>
          <n-space vertical>
            <n-avatar width="100" src="https://avatars.githubusercontent.com/u/7401917?v=4" />
            <h2><a href="https://github.com/ArvinLovegood" target="_blank">@ArvinLovegood</a></h2>
            <p>{{ t('about.followGithub') }}</p>
            <n-image width="300" :src="wxgzh" />
            <p>{{ t('about.buyCoffee') }}</p>
            <n-flex justify="center">
              <n-image width="200" :src="alipay" />
              <n-image width="200" :src="wxpay" />
            </n-flex>
          </n-space>
          <n-divider title-placement="center">{{ t('about.thanks') }}</n-divider>
          <div style="justify-self: center;text-align: left" >
            <p>
              {{ t('about.thanksDonors') }}：
              <n-gradient-text size="small" type="warning">*晨</n-gradient-text><n-divider vertical />
            </p>
            <p>
              {{ t('about.thanksDevelopers') }}：
              <a href="https://github.com/GiCo001" target="_blank">@Gico</a><n-divider vertical />
              <a href="https://github.com/CodeNoobLH" target="_blank">浓睡不消残酒</a><n-divider vertical />
              <a href="https://github.com/gnim2600" target="_blank">@gnim2600</a><n-divider vertical />
              <a href="https://github.com/XXXiaohuayanGGG" target="_blank">@XXXiaohuayanGGG</a><n-divider vertical />
              <a href="https://github.com/2lovecode" target="_blank">@2lovecode</a><n-divider vertical />
              <a href="https://github.com/JerryLookupU" target="_blank">@JerryLookupU</a><n-divider vertical />
            </p>
            <p>
              {{ t('about.thanksOpenSource') }}：
              <a href="https://github.com/wailsapp/wails" target="_blank">Wails</a><n-divider vertical />
              <a href="https://github.com/vuejs" target="_blank">Vue</a><n-divider vertical />
              <a href="https://github.com/tusen-ai/naive-ui" target="_blank">NaiveUI</a><n-divider vertical />
            </p>
          </div>
          <n-divider title-placement="center">{{ t('about.copyrightStatement') }}</n-divider>
          <div style="justify-self: center;text-align: left" >
            <p style="color: #FAA04A">{{ t('about.contactTip') }}</p>
            <p>
              {{ t('about.commercialAuthorization') }}
            </p>
            <n-divider/>
            <p>
              {{ t('about.opensourceBuildDesc') }}
            </p>
            <p>
              {{ t('about.techSupportDesc') }}
            </p>
            <p style="color: #FAA04A">{{ t('about.wechatNote') }}</p>
            <n-table id="support">
              <n-thead>
                <n-tr>
                  <n-th>{{ t('about.technicalSupport') }}</n-th><n-th>Sponsor (RMB)</n-th>
                </n-tr>
              </n-thead>
              <n-tbody>
                <n-tr>
                  <n-td>
                    QQ: 506808970, WeChat: ArvinLovegood
                  </n-td>
                  <n-td>
                    100/次
                  </n-td>
                </n-tr>
                <n-tr>
                  <n-td>
                    Long-term tech support (unlimited, new features priority, etc.)
                  </n-td>
                  <n-td>
                    5000
                  </n-td>
                </n-tr>
              </n-tbody>
            </n-table>
          </div>

        </n-card>

        <n-modal
          v-model:show="showManual"
          preset="card"
          :title="t('about.userManual')"
          style="width: 90vw; max-height: 90vh"
          :bordered="false"
          :segmented="{ content: true, footer: true }"
        >
          <div style="display: flex; max-height: 75vh;">
            <div v-if="catalogList.length" class="manual-catalog" style="width: 240px; min-width: 240px; border-right: 1px solid var(--n-border-color); padding: 8px 4px; overflow-y: auto;">
              <div style="font-weight: bold; margin-bottom: 8px; padding: 0 8px;">{{ t('about.catalog') }}</div>
              <n-tree
                :data="catalogTree"
                :block-line="true"
                :block-node="true"
                :selectable="true"
                :cancelable="false"
                default-expand-all
                key-field="key"
                label-field="label"
                children-field="children"
                @update:selected-keys="onTreeSelect"
              />
            </div>
            <div ref="manualScrollRef" style="flex: 1; overflow-y: auto; padding: 0 16px;">
              <MdPreview style="text-align: left;" :id="manualId" v-model="manualContent" :theme="theme" :preview-theme="'github'" :md-heading-id="slugifyHeading" @onHtmlChanged="extractCatalog" />
            </div>
          </div>
        </n-modal>
      </n-space>
</template>

<style scoped>
h1, h2 {
  margin: 0;
  padding: 6px 0;
}

p {
  margin: 2px 0;
}

ul {
  list-style-type: disc;
  padding-left: 20px;
}

a {
  color: #18a058;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

.manual-catalog > div:hover {
  color: #18a058;
}

.manual-catalog :deep(.n-tree-node-content) {
  text-align: left;
  justify-content: flex-start;
}

.manual-catalog :deep(.n-tree-node) {
  text-align: left;
}
</style>