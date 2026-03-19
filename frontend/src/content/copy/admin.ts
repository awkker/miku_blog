/**
 * DIY 文案入口（后台）
 *
 * 用法：
 * 1. 后台页面出现可见文案，优先放到这里再在组件中引用。
 * 2. 该文件是运营/设计同学改文案的首选入口，避免在多个 Vue 文件里到处找字符串。
 */

export const adminCopy = {
  layout: {
    defaultPageTitle: '仪表盘',
  },
  dashboard: {
    toolbar: {
      filter: 'Filter',
      prevWindowAria: '上一时间窗口',
      nextWindowAria: '下一时间窗口',
      pathFilter: 'Path Filter',
      pathPlaceholder: '输入路径关键字',
    },
    sections: {
      trafficOverviewTitle: 'Traffic Overview',
      trafficOverviewSubtitle: 'Visitors 和 Views 的时间分布',
      recentActivitiesTitle: '近期事项',
      recentActivitiesSubtitle: '管理员最近操作记录',
      pages: 'Pages',
      sources: 'Sources',
      environment: 'Environment',
      location: 'Location',
      geoDistribution: 'Geo Distribution',
      byCountry: 'By Country',
      traffic: 'Traffic',
    },
    ranges: {
      h24: 'Last 24 hours',
      d7: 'Last 7 days',
      d30: 'Last 30 days',
    },
    tabs: {
      path: 'Path',
      entry: 'Entry',
      exit: 'Exit',
      referrers: 'Referrers',
      channels: 'Channels',
      browsers: 'Browsers',
      os: 'OS',
      devices: 'Devices',
      countries: 'Countries',
      regions: 'Regions',
      cities: 'Cities',
    },
    table: {
      path: 'Path',
      source: 'Source',
      name: 'Name',
      region: 'Region',
      visitors: 'Visitors',
      percent: '%',
    },
    stats: {
      visitors: 'Visitors',
      visits: 'Visits',
      views: 'Views',
      bounceRate: 'Bounce rate',
      visitDuration: 'Visit duration',
    },
    actions: {
      approve: '通过了',
      reject: '拒绝了',
      delete: '删除了',
      create: '创建了',
      update: '更新了',
      publish: '发布了',
      unpublish: '下架了',
      schedule: '定时发布了',
    },
    targets: {
      comment: '评论',
      post: '文章',
      friend_link: '友链',
      guestbook: '留言',
      moment: '说说',
    },
    notices: {
      pendingPrefix: '评论留言待审核：',
      pendingSuffix: ' 条',
      pendingSubtitle: '建议优先处理审核队列，避免积压',
      draftPrefix: '草稿待发布：',
      draftSuffix: ' 篇',
      draftSubtitle: '可前往文章管理检查内容后发布',
    },
    mapFallback: {
      loading: '世界地图加载中...',
      empty: '暂无国家分布数据',
      failed: '无法加载地图底图，已在下方显示国家排行',
    },
    common: {
      loading: '加载中...',
      emptyActivities: '暂无操作记录',
      operatorPrefix: '操作人：',
      activityTargetPrefix: '一条',
      unknown: 'Unknown',
      windowFallback: 'Window',
      less: 'Less',
      more: 'More',
      justNow: '刚刚',
      minutesAgo: '分钟前',
      hoursAgo: '小时前',
      daysAgo: '天前',
      noData: '--',
    },
    degradedMessage: 'Analytics 数据暂不可用，当前显示占位数据。请先在 `backend` 目录执行 `go run cmd/migrate/main.go`。',
  },
}
