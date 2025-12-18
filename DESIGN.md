### 项目代号
- 建议名称：**RunGrid** — 强调“快速启动 + 网格整理”，简短专业，贴合产品核心体验。

### 我们要做什么
- 构建一个轻量级的“图标整理 + 快捷启动器”，把桌面/开始菜单入口汇总到可配置网格中。
- 提供分组、搜索、快捷键唤起和托盘常驻，让启动动作变成“输入即达”。

### 为什么
- 桌面和开始菜单入口零散，常用程序难以快速定位，系统搜索对高频启动场景不够轻快。
- 需要一个本地、低占用、专注启动效率的工具，降低整理成本和认知负担。

### 怎么做
- 本地扫描桌面/开始菜单并建立索引，增量更新；去重与分组提升可管理性。
- 图标提取与缓存 + UI 虚拟化，保证大列表也能流畅加载与滚动。
- 快捷键直达搜索框，按收藏/最近使用/启动次数排序，缩短启动路径。
- 架构上采用 Go + Wails，系统能力在后端封装，前端专注交互与效率。

### 技术栈（Go + Wails）
- 后端：Go（Wails），纯 Go SQLite 驱动 `modernc.org/sqlite`，避免 CGO 依赖。
- 前端：Vue/React 任一，配虚拟化列表（如 react-window/virtually）以承载大量图标。
- 系统集成：go-ole 解析 .lnk；Windows API 抽图标（SHGetFileInfo/ExtractIconEx）；ShellExecute 启动 exe/lnk/UWP/URL。
- 快捷键与托盘：Wails 提供全局快捷键和托盘接口。

### 模块划分
- Scanner：并发扫描桌面、开始菜单、常见安装目录；解析 .lnk；识别 exe/UWP/url/文件夹。
- IconExtractor：抽取 ico → 转 png 缓存，缓存命名使用 path 的 hash；控制尺寸（如 128px）。
- Launcher：封装启动策略；路径校验（拒绝不存在/UNC 可疑路径）；URL 白名单协议。
- Deduper：路径规范化 + 文件信息比对；名称相似提示合并。
- SearchEngine：内存索引 name + 拼音首字母，前缀/模糊；排序（收藏 > 最近使用 > 启动次数）。
- Persistence：SQLite CRUD，索引字段 name/path/tags；启动计数/最近使用异步落库。
- Settings：主题、图标大小、网格密度、开机自启、监听桌面变更。
- Tray/Hotkey：托盘菜单、Alt+Space 等唤出窗口。

### 数据模型（示意）
- Item：ID, Name, Path, Type(app/url/folder/doc), IconPath, GroupID, Tags, Favorite, LaunchCount, LastUsedAt, Hidden
- Group：ID, Name, Order, Color
- Settings：Theme, IconSize, Density, AutoStart, WatchDesktop, Hotkey

### MVP 迭代
1) 分组展示 + 手动添加/删除 + 启动；搜索（含拼音首字母）。
2) 扫描桌面/开始菜单 + 去重 + 图标缓存。
3) 收藏/排序 + 托盘 + 全局快捷键。
4) 隐藏/显示桌面图标、开机自启、主题/密度设置。
5) 更新器/签名（如需）、云同步（可选）。

### 性能要点
- 并发扫描 + 限流 semaphore，避免磁盘打满。
- 图标提取异步化，结果落缓存；UI 虚拟化避免大批量渲染。
- 索引常驻内存；启动计数异步写入。
- 纯 Go 依赖减少启动体积和常驻内存。

### 安全与健壮性
- 启动前路径校验，URL 协议白名单。
- 授权/设置（若有）做签名或 HMAC；更新包哈希 + 签名验证，下载走 TLS + pinning。
- 去掉发布版调试符号，日志不输出敏感信息。
