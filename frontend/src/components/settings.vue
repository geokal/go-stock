<script setup>
import {h, onBeforeUnmount, onMounted, ref, computed} from "vue";
import {useI18n} from 'vue-i18n'
import {
  AddPrompt,
  DelPrompt,
  ExportConfig,
  GetConfig,
  GetPromptTemplates,
  SendDingDingMessageByType,
  UpdateConfig,
  CheckSponsorCode,
  FetchAiModels,
  FetchAiModelInfo
} from "../../wailsjs/go/main/App";
import {NTag, NTooltip, NIcon, useMessage} from "naive-ui";
import {data, models} from "../../wailsjs/go/models";
import {EventsEmit} from "../../wailsjs/runtime";
import {HelpCircleFilledIcon, HelpIcon} from "tdesign-icons-vue-next";

const message = useMessage()
const { t, locale } = useI18n()

const formRef = ref(null)
const formValue = ref({
  ID: 1,
  tushareToken: '',
  iwencaiApiKey: '',
  emApiKey: '',
  dingPush: {
    enable: false,
    dingRobot: ''
  },
  localPush: {
    enable: true,
  },
  updateBasicInfoOnStart: false,
  refreshInterval: 1,
  openAI: {
    enable: false,
    aiConfigs: [], // AI配置列表
    prompt: "",
    questionTemplate: "{{stockName}}分析和总结",
    crawlTimeOut: 30,
    kDays: 30,
    httpProxy:"",
    httpProxyEnabled:false,
  },
  enableDanmu: false,
  browserPath: '',
  enableNews: false,
  darkTheme: true,
  enableFund: false,
  enablePushNews: true,
  enableOnlyPushRedNews: true,
  sponsorCode: "",
  httpProxy:"",
  httpProxyEnabled:false,
  enableAgent: false,
  qgqpBId: '',
  updateChannel: 'release',
  promptPlazaApiBase: '',
  language: 'zh-CN',
})

// 添加一个新的AI配置到列表
function addAiConfig() {
  formValue.value.openAI.aiConfigs.push(new data.AIConfig({
    name: '',
    baseUrl: 'https://api.deepseek.com',
    apiKey: '',
    modelName: 'deepseek-reasoner',
    temperature: 0.1,
    maxTokens: 8192,
    timeOut: 6000,
    httpProxy:"",
    httpProxyEnabled:false,
  }));
}

// 从列表中移除一个AI配置
function removeAiConfig(index) {
  const originalCount = formValue.value.openAI.aiConfigs.length;
  // 使用filter创建新数组确保响应式更新
  formValue.value.openAI.aiConfigs = formValue.value.openAI.aiConfigs.filter((_, i) => i !== index);
}

const updateChannelOptions = [
  { label: 'Release', value: 'release' },
  { label: 'Pre-release', value: 'pre' },
  { label: 'Dev', value: 'dev' },
]

const languageOptions = computed(() => [
  { label: t('settings.chinese'), value: 'zh-CN' },
  { label: 'English', value: 'en' },
])

async function fetchAiModels(aiConfig) {
  if (!aiConfig.baseUrl || !aiConfig.apiKey) {
    message.warning(t('settings.fillApiKeyFirst'))
    return
  }
  if (aiConfig._loadingModels) {
    return
  }
  aiConfig._loadingModels = true
  try {
    const list = await FetchAiModels(aiConfig.baseUrl, aiConfig.apiKey)
    const options = (list || []).map(id => ({ label: id, value: id }))
    aiConfig._modelOptions = options
    if (!aiConfig.modelName && options.length > 0) {
      aiConfig.modelName = options[0].value
      onModelNameChange(aiConfig, aiConfig.modelName)
    }
    if (!options.length) {
      message.warning(t('settings.noModelsFromApi'))
    }
  } catch (e) {
    console.error('FetchAiModels error', e)
    message.error(t('settings.fetchModelsFailed'))
  } finally {
    aiConfig._loadingModels = false
  }
}


const promptTemplates = ref([])

const aiPlatformOptions = computed(() => [
  { label: t('settings.platform.deepSeek') + ' (https://api.deepseek.com)', value: 'https://api.deepseek.com' },
  { label: t('settings.platform.siliconFlow') + ' (https://api.siliconflow.cn/v1)', value: 'https://api.siliconflow.cn/v1' },
  { label: t('settings.platform.zhipuAI') + ' (https://open.bigmodel.cn/api/paas/v4)', value: 'https://open.bigmodel.cn/api/paas/v4' },
  { label: t('settings.platform.byteDance') + ' (https://ark.cn-beijing.volces.com/api/v3)', value: 'https://ark.cn-beijing.volces.com/api/v3' },
  { label: t('settings.platform.alibabaCloud') + ' (https://dashscope.aliyuncs.com/compatible-mode/v1)', value: 'https://dashscope.aliyuncs.com/compatible-mode/v1' },
  { label: t('settings.platform.moonshot') + ' (https://api.moonshot.cn/v1)', value: 'https://api.moonshot.cn/v1' },
  { label: t('settings.platform.tencentHunyuan') + ' (https://api.hunyuan.cloud.tencent.com/v1)', value: 'https://api.hunyuan.cloud.tencent.com/v1' },
  { label: t('settings.platform.xunfeiSpark') + ' (https://spark-api-open.xf-yun.com/v1)', value: 'https://spark-api-open.xf-yun.com/v1' },
  { label: t('settings.platform.lingyiwanwu') + ' (https://api.lingyiwanwu.com/v1)', value: 'https://api.lingyiwanwu.com/v1' },
  { label: 'MiniMax (https://api.minimax.chat/v1)', value: 'https://api.minimax.chat/v1' },
  { label: 'Baichuan (https://api.baichuan-ai.com/v1)', value: 'https://api.baichuan-ai.com/v1' },
  { label: 'Baidu Qianfan (https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop)', value: 'https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop' },
  { label: 'OpenAI (https://api.openai.com/v1)', value: 'https://api.openai.com/v1' },
  { label: 'Azure OpenAI (https://YOUR_RESOURCE.openai.azure.com)', value: 'https://YOUR_RESOURCE.openai.azure.com' },
  { label: 'OpenRouter (https://openrouter.ai/api/v1)', value: 'https://openrouter.ai/api/v1' },
  { label: t('settings.platform.ollama') + ' (http://localhost:11434/v1)', value: 'http://localhost:11434/v1' },
])

function getPlatformName(baseUrl) {
  if (!baseUrl) return ''
  const platform = aiPlatformOptions.find(opt => opt.value === baseUrl)
  if (platform) {
    const idx = platform.label.indexOf(' (')
    return idx > 0 ? platform.label.substring(0, idx) : platform.label
  }
  return ''
}

function onBaseUrlChange(aiConfig, newBaseUrl) {
  const platformName = getPlatformName(newBaseUrl)
  if (platformName && aiConfig.name && !aiConfig.name.startsWith(platformName)) {
    aiConfig.name = platformName + '-' + aiConfig.name
  } else if (platformName && !aiConfig.name) {
    aiConfig.name = platformName
  }
}

function onModelNameChange(aiConfig, newModelName) {
  if (!newModelName) return
  const platformName = getPlatformName(aiConfig.baseUrl)
  const baseName = platformName || 'AI'
  
  if (!aiConfig.name) {
    aiConfig.name = baseName + '-' + newModelName
  } else if (aiConfig.name === platformName) {
    aiConfig.name = platformName + '-' + newModelName
  } else {
    const parts = aiConfig.name.split('-')
    if (parts.length >= 2 && parts[0] === platformName) {
      parts[parts.length - 1] = newModelName
      aiConfig.name = parts.join('-')
    } else if (!aiConfig.name.endsWith(newModelName)) {
      aiConfig.name = aiConfig.name + '-' + newModelName
    }
  }

  fetchModelInfo(aiConfig, newModelName)
}

async function fetchModelInfo(aiConfig, modelName) {
  if (!modelName || !aiConfig.baseUrl) return
  try {
    const info = await FetchAiModelInfo(aiConfig.baseUrl, aiConfig.apiKey || '', modelName)
    if (info && info.maxTokens > 0) {
      aiConfig.maxTokens = info.maxTokens
      const sourceLabel = info.source === 'api' ? t('settings.apiSource') : t('settings.builtIn')
      message.success(t('settings.autoSetMaxtokens', { modelName: modelName, maxTokens: info.maxTokens, source: sourceLabel }))
    }
  } catch (e) {
    console.error('FetchAiModelInfo error', e)
  }
}

onMounted(() => {
  GetConfig().then(res => {
    formValue.value.ID = res.ID
    formValue.value.tushareToken = res.tushareToken
    formValue.value.iwencaiApiKey = res.iwencaiApiKey || ''
    formValue.value.emApiKey = res.emApiKey || ''
    formValue.value.dingPush = {
      enable: res.dingPushEnable,
      dingRobot: res.dingRobot
    }
    formValue.value.localPush = {
      enable: res.localPushEnable,
    }
    formValue.value.updateBasicInfoOnStart = res.updateBasicInfoOnStart
    formValue.value.refreshInterval = res.refreshInterval
    // 加载AI配置
    formValue.value.openAI = {
      enable: res.openAiEnable,
      aiConfigs: res.aiConfigs || [],
      prompt: res.prompt,
      questionTemplate: res.questionTemplate ? res.questionTemplate : '{{stockName}}分析和总结',
      crawlTimeOut: res.crawlTimeOut,
      kDays: res.kDays,
      httpProxy:"",
      httpProxyEnabled:false,
    }


    formValue.value.enableDanmu = res.enableDanmu
    formValue.value.browserPath = res.browserPath
    formValue.value.enableNews = res.enableNews
    formValue.value.darkTheme = res.darkTheme
    formValue.value.enableFund = res.enableFund
    formValue.value.enablePushNews = res.enablePushNews
    formValue.value.enableOnlyPushRedNews = res.enableOnlyPushRedNews
    formValue.value.sponsorCode = res.sponsorCode
    formValue.value.httpProxy=res.httpProxy;
    formValue.value.httpProxyEnabled=res.httpProxyEnabled;
    formValue.value.enableAgent = res.enableAgent;
    formValue.value.qgqpBId = res.qgqpBId;
    formValue.value.updateChannel = res.updateChannel || 'release';
    formValue.value.promptPlazaApiBase = res.promptPlazaApiBase || '';
    formValue.value.language = res.language || 'zh-CN';
    locale.value = formValue.value.language;
  })

  // GetPromptTemplates("", "").then(res => {
  //   promptTemplates.value = res
  // })
})
onBeforeUnmount(() => {
  message.destroyAll()
})

function saveConfig() {
  console.log('开始保存设置', formValue.value);
  // 构建配置时，包含aiConfigs列表
  let config = new data.SettingConfig({
    ID: formValue.value.ID,
    dingPushEnable: formValue.value.dingPush.enable,
    dingRobot: formValue.value.dingPush.dingRobot,
    localPushEnable: formValue.value.localPush.enable,
    updateBasicInfoOnStart: formValue.value.updateBasicInfoOnStart,
    refreshInterval: formValue.value.refreshInterval,
    openAiEnable: formValue.value.openAI.enable,
    aiConfigs: formValue.value.openAI.aiConfigs,
    // 序列化aiConfigs列表以传递给后端
    tushareToken: formValue.value.tushareToken,
    iwencaiApiKey: formValue.value.iwencaiApiKey,
    emApiKey: formValue.value.emApiKey,
    prompt: formValue.value.openAI.prompt,
    questionTemplate: formValue.value.openAI.questionTemplate,
    crawlTimeOut: formValue.value.openAI.crawlTimeOut,
    kDays: formValue.value.openAI.kDays,
    enableDanmu: formValue.value.enableDanmu,
    browserPath: formValue.value.browserPath,
    enableNews: formValue.value.enableNews,
    darkTheme: formValue.value.darkTheme,
    enableFund: formValue.value.enableFund,
    enablePushNews: formValue.value.enablePushNews,
    enableOnlyPushRedNews: formValue.value.enableOnlyPushRedNews,
    sponsorCode: formValue.value.sponsorCode,
    httpProxy:formValue.value.httpProxy,
    httpProxyEnabled:formValue.value.httpProxyEnabled,
    enableAgent: formValue.value.enableAgent,
    qgqpBId: formValue.value.qgqpBId,
    updateChannel: formValue.value.updateChannel,
    promptPlazaApiBase: formValue.value.promptPlazaApiBase,
    language: formValue.value.language,
  })

  if (config.sponsorCode) {
    CheckSponsorCode(config.sponsorCode).then(res => {
      if (res.code) {
        UpdateConfig(config).then(res => {
          message.success(res)
          locale.value = formValue.value.language
          EventsEmit("updateSettings", config);
        })
      } else {
        message.error(res.msg)
      }
    })
  } else {
    UpdateConfig(config).then(res => {
      message.success(res)
      locale.value = formValue.value.language
      EventsEmit("updateSettings", config);
    })
  }
}


function getHeight() {
  return document.documentElement.clientHeight
}

function sendTestNotice() {
  let markdown = "### go-stock test\n" + new Date()
  let msg = '{' +
      '     "msgtype": "markdown",' +
      '     "markdown": {' +
      '         "title":"go-stock' + new Date() + '",' +
      '         "text": "' + markdown + '"' +
      '     },' +
      '      "at": {' +
      '          "isAtAll": true' +
      '      }' +
      ' }'

  SendDingDingMessageByType(msg, "test-" + new Date().getTime(), 1).then(res => {
    message.info(res)
  })
}

function exportConfig() {
  ExportConfig().then(res => {
    message.info(res)
  })
}

function importConfig() {
  let input = document.createElement('input');
  input.type = 'file';
  input.accept = '.json';
  input.onchange = (e) => {
    let file = e.target.files[0];
    let reader = new FileReader();
    reader.onload = (e) => {
      let config = JSON.parse(e.target.result);
      formValue.value.ID = config.ID
      formValue.value.tushareToken = config.tushareToken
      formValue.value.iwencaiApiKey = config.iwencaiApiKey || ''
      formValue.value.emApiKey = config.emApiKey || ''
      formValue.value.dingPush = {
        enable: config.dingPushEnable,
        dingRobot: config.dingRobot
      }
      formValue.value.localPush = {
        enable: config.localPushEnable,
      }
      formValue.value.updateBasicInfoOnStart = config.updateBasicInfoOnStart
      formValue.value.refreshInterval = config.refreshInterval
      // 导入AI配置
      formValue.value.openAI = {
        enable: config.openAiEnable,
        aiConfigs: config.aiConfigs || [],
        prompt: config.prompt,
        questionTemplate: config.questionTemplate,
        crawlTimeOut: config.crawlTimeOut,
        kDays: config.kDays
      }
      formValue.value.enableDanmu = config.enableDanmu
      formValue.value.browserPath = config.browserPath
      formValue.value.enableNews = config.enableNews
      formValue.value.darkTheme = config.darkTheme
      formValue.value.enableFund = config.enableFund
      formValue.value.enablePushNews = config.enablePushNews
      formValue.value.enableOnlyPushRedNews = config.enableOnlyPushRedNews
      formValue.value.sponsorCode = config.sponsorCode
      formValue.value.httpProxy=config.httpProxy
      formValue.value.httpProxyEnabled=config.httpProxyEnabled
      formValue.value.enableAgent = config.enableAgent
      formValue.value.qgqpBId = config.qgqpBId
      formValue.value.updateChannel = config.updateChannel || 'release'
    };
    reader.readAsText(file);
  };
  input.click();
}


window.onerror = function (event, source, lineno, colno, error) {
  EventsEmit("frontendError", {
    page: "settings.vue",
    message: event,
    source: source,
    lineno: lineno,
    colno: colno,
    error: error ? error.stack : null
  });
  return true;
};

const showManagePromptsModal = ref(false)
const promptTypeOptions = computed(() => [
  {label: t('promptTemplateList.systemPromptType'), value: '模型系统Prompt'},
  {label: t('promptTemplateList.userPromptType'), value: '模型用户Prompt'},])
const formPromptRef = ref(null)
const formPrompt = ref({
  ID: 0,
  Name: '',
  Content: '',
  Type: '',
})

function managePrompts() {
  formPrompt.value.ID = 0
  showManagePromptsModal.value = true
}

function savePrompt() {
  AddPrompt(formPrompt.value).then(res => {
    message.success(res)
    GetPromptTemplates("", "").then(res => {
      promptTemplates.value = res
    })
    showManagePromptsModal.value = false
  })
}

function editPrompt(prompt) {
  formPrompt.value.ID = prompt.ID
  formPrompt.value.Name = prompt.name
  formPrompt.value.Content = prompt.content
  formPrompt.value.Type = prompt.type
  showManagePromptsModal.value = true
}

function deletePrompt(ID) {
  DelPrompt(ID).then(res => {
    message.success(res)
    GetPromptTemplates("", "").then(res => {
      promptTemplates.value = res
    })
  })
}
</script>

<template>
  <n-flex justify="left" style="text-align: left; --wails-draggable:no-drag">
    <n-form ref="formRef" :label-placement="'left'" :label-align="'left'">
      <n-space vertical size="large">
        <n-card :title="() => h(NTag, { type: 'primary', bordered: false }, () => t('settings.basicSettings'))" size="small">
          <n-grid :cols="24" :x-gap="24" style="text-align: left">
<!--            <n-form-item-gi :span="10" label="Tushare Token：" path="tushareToken">
              <n-input type="text" placeholder="Tushare api token" v-model:value="formValue.tushareToken" clearable/>
            </n-form-item-gi>-->
            <n-form-item-gi :span="4" :label="t('settings.updateBasicInfoOnStart')" path="updateBasicInfoOnStart">
              <n-switch v-model:value="formValue.updateBasicInfoOnStart"/>
            </n-form-item-gi>
            <n-form-item-gi :span="4" :label="t('settings.refreshInterval')" path="refreshInterval">
              <n-input-number v-model:value="formValue.refreshInterval" :placeholder="t('settings.enterRefreshInterval')">
                <template #suffix>{{ t('settings.seconds') }}</template>
              </n-input-number>
            </n-form-item-gi>
            <n-form-item-gi :span="6" :label="t('settings.darkTheme')" path="darkTheme">
              <n-switch v-model:value="formValue.darkTheme"/>
            </n-form-item-gi>
            <n-form-item-gi :span="6" :label="t('settings.language')" path="language">
              <n-select v-model:value="formValue.language" :options="languageOptions" :placeholder="t('settings.selectLanguage')" />
            </n-form-item-gi>
            <n-form-item-gi :span="8" :label="t('settings.updateChannel')" path="updateChannel">
              <n-select v-model:value="formValue.updateChannel" :options="updateChannelOptions" />
              <n-tooltip placement="top">
                <template #trigger>
                  <n-icon color="#0e7a0d" size="20">
                    <HelpCircleFilledIcon />
                  </n-icon>
                </template>
                <template #default>
                  <n-gradient-text :type="'warning'">
                  <div style="max-width: 400px;text-align: left">
                    {{ t('settings.updateChannelNote') }}<br>
                    <b>{{ t('settings.releaseStable') }}</b>：{{ t('settings.releaseStableDesc') }}<br>
                    <b>{{ t('settings.preRelease') }}</b>：{{ t('settings.preReleaseDesc') }}<br>
                    <b>{{ t('settings.devBuild') }}</b>：{{ t('settings.devBuildDesc') }}
                  </div>
                  </n-gradient-text>
                </template>
              </n-tooltip>
            </n-form-item-gi>
            <n-form-item-gi :span="10" :label="t('settings.browserPath')" path="browserPath">
              <n-input type="text" :placeholder="t('settings.enterBrowserPath')" v-model:value="formValue.browserPath" clearable/>
            </n-form-item-gi>
<!--            <n-form-item-gi :span="3" label="指数基金：" path="enableFund">
              <n-switch v-model:value="formValue.enableFund"/>
            </n-form-item-gi>
            <n-form-item-gi :span="3" label="AI智能体：" path="enableAgent">
              <n-switch v-model:value="formValue.enableAgent"/>
            </n-form-item-gi>-->
            <n-form-item-gi :span="11" :label="t('settings.eastmoneyUniqueId')" path="qgqpBId">
              <n-input type="text" :placeholder="t('settings.enterEastmoneyId')" v-model:value="formValue.qgqpBId" clearable/>
              <n-tooltip placement="top">
                <template #trigger>
                  <n-icon color="#0e7a0d" size="20">
                    <HelpCircleFilledIcon />
                  </n-icon>
                </template>
                <template #default>
                  <n-gradient-text :type="'warning'">
                  <div style="max-width: 400px;text-align: left">
                    {{ t('settings.eastmoneyIdHelp') }}
                  </div>
                  </n-gradient-text>
                </template>
              </n-tooltip>
            </n-form-item-gi>

            <n-form-item-gi :span="11" :label="t('settings.wencaiApiKey')" path="iwencaiApiKey">
              <n-input type="password" :placeholder="t('settings.enterWencaiKey')" v-model:value="formValue.iwencaiApiKey" clearable show-password-on="click"/>
              <n-tooltip placement="top">
                <template #trigger>
                  <n-icon color="#0e7a0d" size="20">
                    <HelpCircleFilledIcon />
                  </n-icon>
                </template>
                <template #default>
                  <n-gradient-text :type="'warning'">
                  <div style="max-width: 400px;text-align: left">
                    {{ t('settings.wencaiHelp') }}
                  </div>
                  </n-gradient-text>
                </template>
              </n-tooltip>
            </n-form-item-gi>

            <n-form-item-gi :span="11" :label="t('settings.eastmoneyAiKey')" path="emApiKey">
              <n-input type="password" :placeholder="t('settings.enterEastmoneyAiKey')" v-model:value="formValue.emApiKey" clearable show-password-on="click"/>
              <n-tooltip placement="top">
                <template #trigger>
                  <n-icon color="#0e7a0d" size="20">
                    <HelpCircleFilledIcon />
                  </n-icon>
                </template>
                <template #default>
                  <n-gradient-text :type="'warning'">
                  <div style="max-width: 400px;text-align: left">
                    {{ t('settings.eastmoneyAiHelp') }}
                  </div>
                  </n-gradient-text>
                </template>
              </n-tooltip>
            </n-form-item-gi>

            <n-form-item-gi :span="11" :label="t('settings.sponsorCode')" path="sponsorCode">
              <n-input-group>
                <n-input :show-count="true" :placeholder="t('settings.enterSponsorCode')" v-model:value="formValue.sponsorCode">
                </n-input>
                <n-button type="success" secondary strong
                          @click="CheckSponsorCode(formValue.sponsorCode).then((res) => {message.warning(res.msg)})">{{ t('settings.verify') }}
                </n-button>
                <n-popover trigger="hover" placement="top">
                  <template #trigger>
                    <n-icon color="#0e7a0d" size="20">
                      <HelpCircleFilledIcon />
                    </n-icon>
                  </template>
                  <n-gradient-text :type="'warning'">
                    <div style="max-width: 400px;text-align: left">
                      {{ t('settings.sponsorCodeHowTo') }}<br>
                      {{ t('settings.sponsorCodeHowToDesc') }}
                    </div>
                  </n-gradient-text>
                </n-popover>
              </n-input-group>
            </n-form-item-gi>

            <n-form-item-gi :span="11" :label="t('settings.promptPlazaUrl')" path="promptPlazaApiBase">
              <n-input type="text" :placeholder="t('settings.promptPlazaUrlDefault')" v-model:value="formValue.promptPlazaApiBase" clearable/>
              <n-tooltip placement="top">
                <template #trigger>
                  <n-icon color="#0e7a0d" size="20">
                    <HelpCircleFilledIcon />
                  </n-icon>
                </template>
                <template #default>
                  <n-gradient-text :type="'warning'">
                  <div style="max-width: 400px;text-align: left">
                    {{ t('settings.promptPlazaUrlHelp') }}
                  </div>
                  </n-gradient-text>
                </template>
              </n-tooltip>
            </n-form-item-gi>
          </n-grid>
        </n-card>

        <n-card :title="() => h(NTag, { type: 'primary', bordered: false }, () => t('settings.notificationSettings'))" size="small">
          <n-grid :cols="24" :x-gap="24" style="text-align: left">
            <n-form-item-gi :span="3" :label="t('settings.dingtalkPush')" path="dingPush.enable">
              <n-switch v-model:value="formValue.dingPush.enable"/>
            </n-form-item-gi>
            <n-form-item-gi :span="3" :label="t('settings.localPush')" path="localPush.enable">
              <n-switch v-model:value="formValue.localPush.enable"/>
            </n-form-item-gi>
            <n-form-item-gi :span="3" :label="t('settings.danmuFeature')" path="enableDanmu">
              <n-switch v-model:value="formValue.enableDanmu"/>
            </n-form-item-gi>
            <n-form-item-gi :span="3" :label="t('settings.showScrollNews')" path="enableNews">
              <n-switch v-model:value="formValue.enableNews"/>
            </n-form-item-gi>
            <n-form-item-gi :span="3" :label="t('settings.marketNewsAlert')" path="enablePushNews">
              <n-switch v-model:value="formValue.enablePushNews"/>
            </n-form-item-gi>
            <n-form-item-gi v-if="formValue.enablePushNews" :span="4" :label="t('settings.onlyAlertRedOrWatched')" path="enableOnlyPushRedNews">
              <n-switch v-model:value="formValue.enableOnlyPushRedNews"/>
            </n-form-item-gi>

            <n-form-item-gi :span="22" v-if="formValue.dingPush.enable" :label="t('settings.dingtalkRobotUrl')"
                            path="dingPush.dingRobot">
              <n-input :placeholder="t('settings.enterDingtalkUrl')" v-model:value="formValue.dingPush.dingRobot"/>
              <n-button type="primary" @click="sendTestNotice">{{ t('settings.sendTestNotification') }}</n-button>
            </n-form-item-gi>
          </n-grid>
        </n-card>

        <n-card :title="() => h(NTag, { type: 'primary', bordered: false }, () => t('settings.aiSettings'))" size="small">
          <n-grid :cols="24" :x-gap="24" style="text-align: left;">
            <n-form-item-gi :span="24" :label="t('settings.aiStockAnalysis')" path="openAI.enable">
              <n-switch v-model:value="formValue.openAI.enable"/>
            </n-form-item-gi>

            <n-form-item-gi :span="6" v-if="formValue.openAI.enable" :label="t('settings.crawlerTimeout')"
                            :title="t('settings.crawlerTimeout')" path="openAI.crawlTimeOut">
              <n-input-number min="30" step="1" v-model:value="formValue.openAI.crawlTimeOut"/>
            </n-form-item-gi>
            <n-form-item-gi :span="4" v-if="formValue.openAI.enable" :title="t('settings.moreDaysMoreTokens')"
                            :label="t('settings.dailyKlineDays')" path="openAI.kDays">
              <n-input-number min="30" step="1" max="60" v-model:value="formValue.openAI.kDays"/>
            </n-form-item-gi>
            <n-form-item-gi :span="2" :label="t('settings.httpProxy')" path="httpProxyEnabled">
              <n-switch v-model:value="formValue.httpProxyEnabled"/>
            </n-form-item-gi>
            <n-form-item-gi :span="10" v-if="formValue.httpProxyEnabled" :title="t('settings.httpProxyAddress')"
                            :label="t('settings.httpProxyAddress')" path="httpProxy">
              <n-input type="text" :placeholder="t('settings.enterHttpProxy')" v-model:value="formValue.httpProxy" clearable/>
            </n-form-item-gi>


            <n-gi :span="24" v-if="formValue.openAI.enable">
              <n-divider :title="() => t('settings.defaultPromptSettings')" title-placement="left"></n-divider>
            </n-gi>
            <n-form-item-gi :span="12" v-if="formValue.openAI.enable" :label="t('settings.defaultSystemPrompt')" path="openAI.prompt">
              <n-input v-model:value="formValue.openAI.prompt" type="textarea" :show-count="true"
                       :placeholder="t('settings.enterSystemPrompt')" :autosize="{ minRows: 4, maxRows: 8 }"/>
            </n-form-item-gi>
            <n-form-item-gi :span="12" v-if="formValue.openAI.enable" :label="t('settings.defaultStockAnalysisPrompt')"
                            path="openAI.questionTemplate">
              <n-input v-model:value="formValue.openAI.questionTemplate" type="textarea" :show-count="true"
                       :placeholder="t('settings.enterStockAnalysisPrompt')"
                       :autosize="{ minRows: 4, maxRows: 8 }"/>
            </n-form-item-gi>

            <n-gi :span="24" v-if="formValue.openAI.enable">
              <n-divider :title="() => t('settings.aiModelConfig')" title-placement="left"></n-divider>
            </n-gi>
            <n-gi :span="24" v-if="formValue.openAI.enable">
              <n-space vertical>
                <n-card v-for="(aiConfig, index) in formValue.openAI.aiConfigs" :key="index" :bordered="true"
                        size="small">
                  <template #header>
                    <n-flex justify="space-between" align="center">
                      <n-text depth="3">{{ t('settings.aiConfigIndex', { index: index + 1 }) }}</n-text>
                      <n-button type="error" size="tiny" ghost @click="removeAiConfig(index)">{{ t('settings.removeAiConfig') }}</n-button>
                    </n-flex>
                  </template>
                  <n-grid :cols="24" :x-gap="24">
                    <n-form-item-gi :span="24" hidden :label="t('settings.configId')" :path="`openAI.aiConfigs[${index}].ID`">
                      <n-input type="text" :placeholder="t('settings.enterConfigId')" v-model:value="aiConfig.ID" clearable/>
                    </n-form-item-gi>
                    <n-form-item-gi :span="12" :label="t('settings.configName')" :path="`openAI.aiConfigs[${index}].name`">
                      <n-input type="text" :placeholder="t('settings.enterConfigName')" v-model:value="aiConfig.name" clearable/>
                    </n-form-item-gi>
                    <n-form-item-gi :span="12" :label="t('settings.apiAddress')" :path="`openAI.aiConfigs[${index}].baseUrl`">
                      <n-select
                        v-model:value="aiConfig.baseUrl"
                        :options="aiPlatformOptions"
                        filterable
                        tag
                        clearable
                        :placeholder="t('settings.enterApiAddress')"
                        @update:value="(val) => onBaseUrlChange(aiConfig, val)"
                      />
                    </n-form-item-gi>
                    <n-form-item-gi :span="12" :label="t('settings.apiKey')" :path="`openAI.aiConfigs[${index}].apiKey`">
                      <n-input type="password" :placeholder="t('settings.enterApiKey')" v-model:value="aiConfig.apiKey" clearable
                               show-password-on="click"/>
                    </n-form-item-gi>
                    <n-form-item-gi :span="8" :label="t('settings.modelName')" :path="`openAI.aiConfigs[${index}].modelName`">
                      <n-select
                        v-model:value="aiConfig.modelName"
                        :options="aiConfig._modelOptions || []"
                        filterable
                        tag
                        :loading="aiConfig._loadingModels"
                        :placeholder="t('settings.getModelListOrEnter')"
                        @click="fetchAiModels(aiConfig)"
                        @update:value="(val) => onModelNameChange(aiConfig, val)"
                      />
                    </n-form-item-gi>
                    <n-form-item-gi :span="5" label="Temperature" :path="`openAI.aiConfigs[${index}].temperature`">
                      <n-input-number :placeholder="t('settings.temperature')" v-model:value="aiConfig.temperature" :step="0.1"/>
                    </n-form-item-gi>
                    <n-form-item-gi :span="5" label="MaxTokens" :path="`openAI.aiConfigs[${index}].maxTokens`">
                      <n-input-number :placeholder="t('settings.maxTokens')" v-model:value="aiConfig.maxTokens"/>
                    </n-form-item-gi>
                    <n-form-item-gi :span="5" label="Timeout" :path="`openAI.aiConfigs[${index}].timeOut`">
                      <n-input-number min="60" step="1" :placeholder="t('settings.enterTimeout')" v-model:value="aiConfig.timeOut"/>
                    </n-form-item-gi>
                    <n-form-item-gi :span="12" :label="t('settings.proxyEnabled')" :path="`openAI.aiConfigs[${index}].httpProxyEnabled`">
                      <n-switch v-model:value="aiConfig.httpProxyEnabled"/>
                    </n-form-item-gi>
                    <n-form-item-gi :span="12" v-if="aiConfig.httpProxyEnabled" :title="t('settings.httpProxyAddress')" :path="`openAI.aiConfigs[${index}].httpProxy`">
                      <n-input type="text" :placeholder="t('settings.enterHttpProxy')" v-model:value="aiConfig.httpProxy" clearable/>
                    </n-form-item-gi>
                  </n-grid>
                </n-card>
                <n-button type="primary" dashed @click="addAiConfig" style="width: 100%;">{{ t('settings.addAiConfig') }}</n-button>
              </n-space>
            </n-gi>

            <n-gi :span="24">
              <n-divider/>
            </n-gi>

            <n-gi :span="24">
              <n-space vertical>
                <n-space justify="center">
<!--                  <n-button type="warning" @click="managePrompts">管理提示词模板</n-button>-->
                  <n-button type="primary" strong @click="saveConfig">{{ t('settings.saveSettings') }}</n-button>
                  <n-button type="info" @click="exportConfig">{{ t('settings.exportSettings') }}</n-button>
                  <n-button type="error" @click="importConfig">{{ t('settings.importSettings') }}</n-button>
                </n-space>

                <n-flex justify="start" style="margin-top: 10px" v-if="promptTemplates.length > 0">
                  <n-tag :bordered="false" type="warning">{{ t('settings.promptTemplate') }}:</n-tag>
                  <n-tag size="medium" secondary v-for="prompt in promptTemplates" closable
                         @close="deletePrompt(prompt.ID)" @click="editPrompt(prompt)" :title="prompt.content"
                         :type="prompt.type === '模型系统Prompt' ? 'success' : 'info'" :bordered="false">{{
                      prompt.name
                    }}
                  </n-tag>
                </n-flex>
              </n-space>
            </n-gi>

          </n-grid>
        </n-card>
      </n-space>
    </n-form>
  </n-flex>

  <n-modal v-model:show="showManagePromptsModal" closable :mask-closable="false">
    <n-card style="width: 800px; height: 600px; text-align: left" :bordered="false"
            :title="(formPrompt.ID > 0 ? t('settings.modify') : t('settings.add')) + ' ' + t('settings.promptTemplate')" size="huge" role="dialog" aria-modal="true">
      <n-form ref="formPromptRef" :label-placement="'left'" :label-align="'left'">
        <n-form-item :label="t('settings.promptName')">
          <n-input v-model:value="formPrompt.Name" :placeholder="t('settings.enterPromptName')"/>
        </n-form-item>
        <n-form-item :label="t('settings.promptType')">
          <n-select v-model:value="formPrompt.Type" :options="promptTypeOptions" :placeholder="t('settings.selectPromptType')"/>
        </n-form-item>
        <n-form-item :label="t('settings.promptContent')">
          <n-input v-model:value="formPrompt.Content" type="textarea" :show-count="true" :placeholder="t('settings.enterPromptContent')"
                   :autosize="{ minRows: 12, maxRows: 12, }"/>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-flex justify="end">
          <n-button type="primary" @click="savePrompt">{{ t('settings.save') }}</n-button>
          <n-button type="warning" @click="showManagePromptsModal = false">{{ t('settings.cancel') }}</n-button>
        </n-flex>
      </template>
    </n-card>
  </n-modal>
</template>

<style scoped>
.cardHeaderClass {
  font-size: 16px;
  font-weight: bold;
  color: red;
}
</style>