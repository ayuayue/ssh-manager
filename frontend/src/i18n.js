const messages = {
  zh: {
    nav: {
      config: '配置管理',
      keys: '密钥管理',
      hosts: '主机管理',
      logs: '操作日志',
    },
    config: {
      title: 'SSH 配置',
      reload: '重载',
      validate: '校验',
      saving: '保存中...',
      save: '保存',
      parsedHosts: '解析的主机',
      placeholder: 'SSH 配置内容...',
    },
    keys: {
      title: 'SSH 密钥',
      generate: '生成密钥',
      generateKey: '生成 SSH 密钥',
      noKeys: '未找到 SSH 密钥',
      generateFirst: '生成你的第一个密钥',
      keyType: '密钥类型',
      keySize: '密钥长度',
      curveSize: '曲线大小',
      keyName: '密钥名称',
      email: '邮箱 / 备注',
      passphrase: '密码短语',
      passphraseHint: '留空表示不设置密码',
      cancel: '取消',
      generating: '生成中...',
      generateBtn: '生成',
      viewPubKey: '查看公钥',
      viewPrivKey: '查看私钥',
      deleteKey: '删除密钥',
      deleteConfirm: '确定要删除密钥',
      deleteWarning: '此操作不可撤销！',
      copied: '已复制到剪贴板！',
      copyFailed: '复制失败',
      publicKey: '公钥',
      privateKey: '私钥',
      close: '关闭',
      copy: '复制',
      optional: '可选',
    },
    hosts: {
      title: '主机管理',
      import: '从配置导入',
      importing: '导入中...',
      addHost: '添加主机',
      editHost: '编辑主机',
      favorites: '收藏',
      history: '连接历史',
      searchPlaceholder: '搜索主机、用户、IP...',
      noFavorites: '还没有收藏的主机',
      importFromConfig: '从配置导入',
      addManually: '手动添加',
      noHistory: '没有连接历史',
      alias: '别名',
      hostname: '主机名 / IP',
      port: '端口',
      user: '用户',
      tags: '标签',
      group: '分组',
      tagsHint: '逗号分隔',
      test: '测试',
      edit: '编辑',
      del: '删除',
      deleteConfirm: '确定删除此主机？',
      justNow: '刚刚',
      minutesAgo: '分钟前',
      hoursAgo: '小时前',
      required: '别名和主机名必填',
    },
    theme: {
      light: '亮色模式',
      dark: '暗色模式',
    },
    sidebar: {
      expand: '展开侧栏',
      collapse: '收起侧栏',
    },
    lang: {
      zh: '中文',
      en: 'English',
    },
    logs: {
      title: '操作日志',
      allActions: '全部操作',
      reload: '刷新',
      loading: '加载日志...',
      empty: '暂无操作日志',
      configSave: '配置保存',
      configImport: '配置导入',
      keyGenerate: '密钥生成',
      keyDelete: '密钥删除',
      sshConnect: 'SSH 连接',
      justNow: '刚刚',
      minutesAgo: '分钟前',
      hoursAgo: '小时前',
    },
  },
  en: {
    nav: {
      config: 'Config',
      keys: 'Keys',
      hosts: 'Hosts',
      logs: 'Audit Logs',
    },
    config: {
      title: 'SSH Config',
      reload: 'Reload',
      validate: 'Validate',
      saving: 'Saving...',
      save: 'Save',
      parsedHosts: 'Parsed Hosts',
      placeholder: 'SSH config content...',
    },
    keys: {
      title: 'SSH Keys',
      generate: 'Generate Key',
      generateKey: 'Generate SSH Key',
      noKeys: 'No SSH keys found',
      generateFirst: 'Generate your first key',
      keyType: 'Key Type',
      keySize: 'Key Size',
      curveSize: 'Curve Size',
      keyName: 'Key Name',
      email: 'Email / Comment',
      passphrase: 'Passphrase',
      passphraseHint: 'Leave empty for no passphrase',
      cancel: 'Cancel',
      generating: 'Generating...',
      generateBtn: 'Generate',
      viewPubKey: 'View PubKey',
      viewPrivKey: 'View PrivKey',
      deleteKey: 'Delete Key',
      deleteConfirm: 'Are you sure you want to delete',
      deleteWarning: 'This cannot be undone!',
      copied: 'Copied to clipboard!',
      copyFailed: 'Failed to copy',
      publicKey: 'Public Key',
      privateKey: 'Private Key',
      close: 'Close',
      copy: 'Copy',
      optional: 'optional',
    },
    hosts: {
      title: 'Host Management',
      import: 'Import from Config',
      importing: 'Importing...',
      addHost: 'Add Host',
      editHost: 'Edit Host',
      favorites: 'Favorites',
      history: 'History',
      searchPlaceholder: 'Search hosts, users, IPs...',
      noFavorites: 'No favorite hosts yet',
      importFromConfig: 'Import from Config',
      addManually: 'Add Manually',
      noHistory: 'No connection history',
      alias: 'Alias',
      hostname: 'Hostname / IP',
      port: 'Port',
      user: 'User',
      tags: 'Tags',
      group: 'Group',
      tagsHint: 'comma separated',
      test: 'Test',
      edit: 'Edit',
      del: 'Del',
      deleteConfirm: 'Delete this host?',
      justNow: 'Just now',
      minutesAgo: 'm ago',
      hoursAgo: 'h ago',
      required: 'Alias and Hostname required',
    },
    theme: {
      light: 'Light Mode',
      dark: 'Dark Mode',
    },
    sidebar: {
      expand: 'Expand sidebar',
      collapse: 'Collapse sidebar',
    },
    lang: {
      zh: '中文',
      en: 'English',
    },
    logs: {
      title: 'Audit Logs',
      allActions: 'All Actions',
      reload: 'Refresh',
      loading: 'Loading logs...',
      empty: 'No audit logs yet',
      configSave: 'Config Saved',
      configImport: 'Config Imported',
      keyGenerate: 'Key Generated',
      keyDelete: 'Key Deleted',
      sshConnect: 'SSH Connected',
      justNow: 'Just now',
      minutesAgo: 'm ago',
      hoursAgo: 'h ago',
    },
  },
}

let currentLang = 'zh'

export function setLang(lang) {
  currentLang = lang
  localStorage.setItem('ssh-manager-lang', lang)
}

export function getLang() {
  if (currentLang === 'zh') return 'zh'
  const saved = localStorage.getItem('ssh-manager-lang')
  if (saved && messages[saved]) {
    currentLang = saved
  }
  return currentLang
}

export function t(key) {
  const keys = key.split('.')
  let obj = messages[getLang()]
  for (const k of keys) {
    if (obj && obj[k] !== undefined) {
      obj = obj[k]
    } else {
      return key
    }
  }
  return obj
}
