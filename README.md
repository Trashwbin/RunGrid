# RunGrid（启动网格）

轻量级图标整理 + 快捷启动器，聚合桌面/开始菜单入口，提供分组、搜索、托盘常驻与全局快捷键。

## 功能概览

- 扫描桌面与开始菜单，自动发现快捷方式与应用入口
- 分组管理、搜索、收藏与启动统计
- 图标提取与本地缓存，启动体验更轻快
- 托盘常驻 + 全局快捷键唤出
- 分组规则导入：基于 `target_name` 一键归类
- 启动方式可选：单击启动 / 双击启动
- 面板关闭时机可选：不自动关闭 / 启动后 / 失焦后 / 启动或失焦

## 目录结构

- `backend/` Go 后端（扫描、图标提取、启动、存储）
- `frontend/` React 前端（界面、交互、设置）
- `assets/` 应用图标等资源（`assets/icons/app.png`、`assets/icons/tray.ico`）
- `rules/` 分组规则示例（`rules/devrule.json`）
- `bland/` 品牌资产（设计源文件/导出图）
- `test/` 本地测试脚本（默认忽略）

## 开发与构建

依赖：Go、Node.js、Wails v2

开发模式：
```
wails dev
```

生产构建：
```
wails build
```

## 分组规则导入

通过菜单「导入分组规则」选择 JSON 文件，按 `target_name` 自动归类。

规则结构（简化）：
```json
{
  "version": "1.0",
  "groups": [
    {"id": "dev", "name": "开发", "category": "app", "order": 10, "color": "#2F80ED", "icon": "code"}
  ],
  "rules": [
    {"group_id": "dev", "match": {"target_name": ["code.exe", "postman.exe"]}}
  ]
}
```

说明：
- `category` 取值：`app` / `system` / `doc` / `folder` / `url`
- `target_name` 建议全部小写（扫描时会规范化为小写）
- 同一个 `target_name` 命中后只归入一个分组

## 数据存储

SQLite 数据库默认位于用户配置目录下的 `rungrid/rungrid.db`。

## 设计与决策

项目设计文档：`DESIGN.md`
