/**
 * DIY 文案入口（站点级）
 *
 * 用法：
 * 1. 修改本文件即可调整首页 Dock、博客/说说/留言板/关于我页面与其核心组件文案。
 * 2. 新页面需要可配置文案时，优先在本文件新增字段，再由页面读取，不要直接写死在组件模板里。
 */

export interface SiteNavItem {
  href: string
  label: string
  auth?: 'authed' | 'guest'
}

export interface HomeDockItem extends SiteNavItem {
  icon: string
  name: string
}

export const siteCopy = {
  brand: {
    logoAlt: 'nanamiku logo',
    text: 'NANAMIKU BLOG',
  },
  home: {
    metaTitle: 'nanamiku blog',
    metaDescription: 'nanamiku blog 开屏页',
    heroTitle: '薰逸の猫窝',
    heroSubtitle: '「月が綺麗ですね, 風も優しいですね」。',
    dockItems: [
      { name: '访客', href: '/about', icon: 'person' },
      { name: '博客', href: '/blog', icon: 'book' },
      { name: '说说', href: '/moments', icon: 'moments' },
      { name: '留言', href: '/guestbook', icon: 'message' },
      { name: '友链', href: '/friends', icon: 'link' },
      { name: '登录', href: '/login', icon: 'login', auth: 'guest' },
      { name: '后台', href: '/admin', icon: 'admin', auth: 'authed' },
    ] as HomeDockItem[],
  },
  blogTopNav: {
    ariaPrefix: '前往',
    primaryLinks: [
      { href: '/blog', label: '博客首页' },
      { href: '/moments', label: '说说' },
      { href: '/guestbook', label: '留言板' },
      { href: '/friends', label: '友情链接' },
      { href: '/about', label: '关于' },
    ] as SiteNavItem[],
    actionLinks: [
      { href: '/admin', label: '后台', auth: 'authed' },
      { href: '/login', label: '登录', auth: 'guest' },
    ] as SiteNavItem[],
  },
  blogIndex: {
    metaTitle: 'Miku Blog',
    metaDescription: '博客主页',
    heroBadge: 'CREATOR SPACE',
    heroTitle: 'NanaMiku Blog',
    heroDescription: '我是 Nana Miku，专注写 Astro、Vue Islands 与 Go 全栈实践。这里不是纯文章流，而是我持续构建的技术创作空间。',
    heroActions: [
      { href: '#latest-posts', label: '看最新文章' },
      { href: '#archive', label: '进入归档' },
      { href: '/about', label: '关于我' },
    ],
    quickStats: [
      { label: '技术栈', value: 'Astro + Vue + Go + PostgreSQL' },
      { label: '本月主线', value: '首页信息架构升级与阅读路径优化' },
      { label: '当前项目', value: 'Glass UI Lab · 液态玻璃组件实验集' },
    ],
    focusCard: {
      badge: '本月在做什么',
      title: '把博客从文章列表升级为创作者入口',
      description: '重点在信息分层、阅读引导与个人品牌辨识度，让访客在 10 秒内知道你是谁、写什么、从哪里开始读。',
      footnote: '文章来自后台管理面板发布',
    },
    scrollCue: {
      label: '向下阅读',
      ariaLabel: '向下阅读',
    },
    authorCard: {
      avatarAlt: '博主头像',
      name: 'Nana Miku',
      role: 'Front-end Developer / Writer',
      bio: '写前端、写系统、写日常。希望每篇都能让人少走一点弯路。',
      location: 'China',
      since: 'Since 2026',
      skills: ['Astro', 'Vue', 'Go', 'TypeScript', 'Tailwind'],
      aboutCta: '查看完整个人介绍',
    },
    socialLinks: [
      {
        label: 'GitHub',
        href: 'https://github.com',
        icon: 'M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 00-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0020 4.77 5.07 5.07 0 0019.91 1S18.73.65 16 2.48a13.38 13.38 0 00-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 005 4.77a5.44 5.44 0 00-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 009 18.13V22',
      },
      {
        label: 'Twitter',
        href: 'https://twitter.com',
        icon: 'M23 3a10.9 10.9 0 01-3.14 1.53 4.48 4.48 0 00-7.86 3v1A10.66 10.66 0 013 4s-4 9 5 13a11.64 11.64 0 01-7 2c9 5 20 0 20-11.5a4.5 4.5 0 00-.08-.83A7.72 7.72 0 0023 3z',
      },
      {
        label: 'Mail',
        href: 'mailto:hello@nanamiku.blog',
        icon: 'M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2zm16 2l-8 5-8-5',
      },
    ],
    nowTitle: 'Now / 近况',
    nowItems: ['重构博客首页信息架构', '整理 Astro + Vue Islands 模板', '记录前端交互动效实践'],
    playlistTitle: '正在听',
    playlist: ['宇多田光 - First Love', 'FKJ - Ylang Ylang', '椎名林檎 - 丸ノ内サディスティック'],
    backToTopAria: '返回顶部',
  },
  momentsPage: {
    metaTitle: '说说 | Miku Blog',
    metaDescription: 'Miku Blog 说说动态',
    heroBadge: 'MOMENTS',
    title: '说说',
    description: '碎片化的日常记录。随想、随拍、随心。',
    navAria: '页面导航',
    quickLinks: [
      { href: '/blog', label: '博客首页', ariaLabel: '前往博客首页' },
      { href: '/guestbook', label: '留言板', ariaLabel: '前往留言板' },
      { href: '/friends', label: '友情链接', ariaLabel: '前往友情链接' },
    ],
  },
  guestbookPage: {
    metaTitle: '留言板 | Miku Blog',
    metaDescription: 'Miku Blog 访客留言板',
    heroBadge: 'GUESTBOOK',
    title: '访客留言板',
    description: '这里是站点的公共交流区。欢迎留下你的想法、建议或一句路过打卡。',
    navAria: '页面导航',
    quickLinks: [
      { href: '/blog', label: '博客首页', ariaLabel: '前往博客首页' },
      { href: '/friends', label: '友情链接', ariaLabel: '前往友情链接页面' },
    ],
  },
  aboutPage: {
    metaTitle: '关于我 | Miku Blog',
    metaDescription: 'Nana Miku 的个人介绍与创作地图',
    heroBadge: 'ABOUT CREATOR',
    heroTitle: '我是 Nana Miku',
    heroDescription: '我把这个博客当作长期创作者空间来做。关注 Astro、Vue Islands 与 Go 全栈实践，持续写可复用、可落地、可维护的经验总结。',
    identityTags: ['Astro', 'Vue Islands', 'Go', 'Tailwind', 'Nano Stores'],
    heroActions: [
      { href: '/blog#latest-posts', label: '看精选文章' },
      { href: '#projects', label: '看主推项目' },
      { href: 'mailto:hello@nanamiku.blog?subject=来自博客的联系', label: '联系我' },
    ],
    stats: {
      postsLabel: '已发布文章',
      focusLabel: '写作主线',
      focusValue: '前端体验 + 全栈实战',
      statusLabel: '当前状态',
      statusValue: '持续迭代中',
    },
    profileCard: {
      avatarAlt: 'Nana Miku 头像',
      name: 'Nana Miku',
      role: 'Front-end Developer / Writer',
      nowTitle: 'NOW / 近况',
      quote: '让每篇文章都成为下一次开发可以直接复用的经验块。',
    },
    nowItems: [
      '重构博客首页信息架构与阅读路径',
      '打磨 LiquidGlassCard 的组件化用法',
      '沉淀 Astro + Vue Islands 项目模板',
    ],
    introCards: [
      {
        title: '我在做什么',
        description: '围绕个人博客做体验与结构升级，让页面既有设计感，也有明确的信息引导和内容路径。',
      },
      {
        title: '我为什么写作',
        description: '写作是我整理认知和沉淀方法的方式。希望读者看到后，能立刻拿去解决真实问题。',
      },
      {
        title: '我关心的问题',
        description: '信息层次是否清晰、组件能否复用、交互是否克制、系统是否能长期迭代。',
      },
    ],
    timeline: {
      title: '成长时间线',
      subtitle: '从写作到全栈，再到创作者空间化。',
      resultPrefix: '结果：',
    },
    milestones: [
      {
        year: '2022',
        title: '开始系统写作',
        summary: '从零散笔记转向结构化输出，建立个人技术写作节奏。',
        result: '累计产出 20+ 篇可复用实践文章',
      },
      {
        year: '2023',
        title: '聚焦前端工程化',
        summary: '主攻组件复用、样式规范与协作流程，减少重复开发。',
        result: '沉淀多套前端页面模板与规范',
      },
      {
        year: '2024',
        title: '全栈化推进',
        summary: '将前端体验与 Go 后端能力联动，关注性能与可维护性。',
        result: '形成 Astro + Go 的完整博客技术路线',
      },
      {
        year: '2026',
        title: '创作者空间升级',
        summary: '把博客从文章列表升级为有路径、有节奏的创作者入口。',
        result: '上线 About / Blog 的新信息架构',
      },
    ],
    capabilityGroups: [
      {
        title: '前端体验',
        desc: '关注视觉层次、交互节奏与页面信息密度。',
        stack: ['Astro', 'Vue 3', 'Tailwind CSS', '组件化设计'],
      },
      {
        title: '工程化',
        desc: '用规范和流程提升迭代效率与可维护性。',
        stack: ['TypeScript', '模块化结构', '内容驱动页面', '性能优化'],
      },
      {
        title: '后端协同',
        desc: '强调接口稳定性、数据组织和全栈落地能力。',
        stack: ['Go', 'Hertz', 'sqlc', 'PostgreSQL / Redis'],
      },
    ],
    projectsSection: {
      title: '主推项目',
      subtitle: '我最希望访客优先进入的 3 条路径。',
      cta: '查看详情',
    },
    featuredProjects: [
      {
        name: 'Miku Blog Engine',
        focus: 'Astro + Go 的轻量博客引擎',
        role: '负责信息架构、前后端联调、视觉规范落地',
        metric: '首屏信息密度与阅读路径明显增强',
        href: '/blog',
      },
      {
        name: 'Glass UI Lab',
        focus: '液态玻璃 UI 组件实验库',
        role: '负责组件 API 设计与视觉统一策略',
        metric: '形成可复用卡片与面板实现模式',
        href: '/blog',
      },
      {
        name: 'Nano Store Patterns',
        focus: '跨 Astro / Vue 的状态协作模式',
        role: '负责状态拆分方案与示例页面验证',
        metric: '降低跨页面交互耦合与重复状态逻辑',
        href: '/blog',
      },
    ],
    writingMapSection: {
      title: '写作地图',
      subtitle: '如果你是第一次来，可以这样阅读。',
      pathPrefix: 'PATH',
      cta: '进入',
    },
    writingMap: [
      {
        title: '入门先看',
        desc: '先从精选文章了解我最核心的技术观点与实现方式。',
        href: '/blog#latest-posts',
      },
      {
        title: '系列专栏',
        desc: '按专题连续阅读，适合系统学习同一方向。',
        href: '/blog#archive',
      },
      {
        title: '最近更新',
        desc: '快速查看我当前在推进的主线和近期输出。',
        href: '/blog#latest-posts',
      },
    ],
    monthlyGoalsTitle: '本月目标',
    monthlyGoals: ['完成 4 篇高质量文章', '完善项目展示页', '优化移动端阅读体验'],
    listeningTitle: '正在听',
    listeningNow: ['宇多田光 - First Love', 'FKJ - Ylang Ylang', '椎名林檎 - 丸ノ内サディスティック'],
    signature: {
      title: '固定签名',
      description: '技术不是为了炫技，而是为了让复杂问题更清晰、让协作更顺畅。',
      footer: '愿每次更新都比上次更有价值。',
    },
    contactSection: {
      title: '如果你也在做类似方向，欢迎交流',
      subtitle: '可以聊前端体验、工程化方案、博客搭建与内容结构。',
      emailButton: '发邮件',
      emailHref: 'mailto:hello@nanamiku.blog?subject=合作交流',
      githubButton: 'GitHub',
    },
    githubUsername: 'awkker',
    socialLinks: [
      { label: 'QQ', href: 'https://qm.qq.com/q/c7DY18rEju' },
      { label: 'Bilibili', href: 'https://space.bilibili.com/1969160969' },
      { label: '抖音', href: 'https://www.douyin.com/user/MS4wLjABAAAATzdjtBBrLLCn69TtPMeseuEUzztbNZzw-9f13adrfiM' },
      { label: '小红书', href: 'https://www.xiaohongshu.com/user/profile/6427cf87000000002901166e' },
    ],
    backToTopAria: '返回顶部',
  },
  components: {
    blogFeed: {
      retry: '重试',
      empty: '暂无已发布的文章',
      featuredBadge: '置顶阅读',
      readSuffix: ' 次阅读',
      likeSuffix: ' 点赞',
      publishedPrefix: '发布：',
      shortLikeSuffix: ' 赞',
      loadError: '加载文章失败，请检查后端服务是否运行',
    },
    latestMoments: {
      title: 'MOMENTS',
      viewAll: '查看全部 →',
      loadFailed: '加载失败',
      retry: '重试',
      empty: '暂无说说',
      imageAltPrefix: '图片 ',
    },
    momentsBoard: {
      countSuffix: ' 条说说',
      loadErrorFallback: '说说加载失败，请稍后再试。',
      submitError: '发布失败，请检查网络后重试。',
      emptyTitle: '还没有说说',
      emptyDescription: '发布第一条说说，记录此刻的想法。',
    },
    momentCard: {
      imageAltPrefix: '图片 ',
      nicknamePlaceholder: '你的昵称',
      commentPlaceholder: '写下评论...',
      send: '发送',
      previewAlt: '预览大图',
      shareTitleSuffix: ' 的说说',
      copiedToast: '链接已复制到剪贴板',
      commentSentToast: '评论已发送',
      commentFailedToast: '评论发送失败，请稍后重试',
    },
    guestbookBoard: {
      postTitle: '发布留言',
      nicknameLabel: '昵称',
      nicknamePlaceholder: '请输入昵称',
      nicknameAria: '昵称',
      websiteLabel: '网址（可选）',
      websitePlaceholder: 'https://example.com',
      websiteAria: '个人网址',
      messageLabel: '留言内容',
      messagePlaceholder: '想说点什么？欢迎留下你的想法。',
      messageAria: '留言内容',
      submitAria: '发送留言',
      submitLoading: '发送中...',
      submitIdle: '发布',
      sortTabs: [
        { label: '最热', value: 'hot' },
        { label: '最新', value: 'newest' },
        { label: '最早', value: 'oldest' },
      ],
      countSuffix: ' 条留言',
      loadErrorFallback: '留言读取失败，请稍后再试。',
      submitError: '留言提交失败，请检查网络后重试。',
      emptyTitle: '还没有留言',
      emptyDescription: '成为第一个在这里留下足迹的人。',
      validation: {
        nicknameRequired: '昵称不能为空',
        messageRequired: '留言内容不能为空',
        websiteInvalid: '网址格式不合法，请输入 http(s) 链接',
      },
      toasts: {
        submitSuccess: '留言已提交，等待审核',
        submitFailed: '留言发送失败，请稍后重试',
        replySuccess: '回复已提交，等待审核',
        replyFailed: '回复发送失败，请稍后重试',
      },
    },
    guestbookMessageCard: {
      upvoteAria: '赞同',
      downvoteAria: '反对',
      replyButton: '回复',
      replyCountSuffix: ' 条回复',
      nicknamePlaceholder: '你的昵称',
      contentPlaceholder: '写下你的回复...',
      cancel: '取消',
      submitReply: '回复',
    },
    aboutGithubProfile: {
      loadError: 'GitHub 数据暂时无法加载，请稍后刷新页面重试。',
      viewProfile: '查看 GitHub 主页 →',
      stats: {
        repos: 'Repos',
        stars: 'Stars',
        followers: 'Followers',
      },
      activityTitle: 'GitHub 活动概览',
      activitySubtitle: '近 12 个月提交趋势',
      monthSuffix: '月',
      techStackTitle: '技术栈',
      techStackSubtitle: '基于 GitHub 仓库语言自动分析',
      recentReposTitle: '最近活跃项目',
      recentReposSubtitle: '来自 GitHub 最近有更新的仓库',
    },
  },
}
