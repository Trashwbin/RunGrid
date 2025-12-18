import type {AppItem, Category, GroupTab, MenuItem} from '../types';

export const categories: Category[] = [
  {id: 'apps', label: '应用程序', icon: 'apps'},
  {id: 'system', label: '系统应用', icon: 'system'},
  {id: 'docs', label: '书籍文档', icon: 'docs'},
  {id: 'folders', label: '文件夹', icon: 'folder'},
];

export const groupTabs: GroupTab[] = [
  {id: 'all', label: '全部'},
  {id: 'dev', label: 'dev'},
];

export const menuItems: MenuItem[] = [
  {id: 'scan', label: '扫描快捷方式'},
  {id: 'dedupe', label: '清理重复项'},
  {id: 'clear', label: '清空项目'},
  {id: 'hide-desktop', label: '隐藏桌面图标'},
  {id: 'show-desktop', label: '显示桌面图标'},
  {id: 'help', label: '帮助文档'},
];

export const appItems: AppItem[] = [
  {
    id: 'apifox',
    name: 'Apifox',
    categoryId: 'apps',
    groupId: 'dev',
    accent: 'pink',
    glyph: 'A',
  },
  {
    id: 'docker',
    name: 'Docker Desktop',
    categoryId: 'apps',
    groupId: 'dev',
    accent: 'blue',
    glyph: 'D',
  },
  {
    id: 'postman',
    name: 'Postman',
    categoryId: 'apps',
    groupId: 'dev',
    accent: 'orange',
    glyph: 'P',
  },
  {
    id: 'vscode-insiders',
    name: 'Visual Studio Code - Insiders',
    categoryId: 'apps',
    groupId: 'dev',
    accent: 'teal',
    glyph: '<>',
  },
  {
    id: 'vscode',
    name: 'Visual Studio Code',
    categoryId: 'apps',
    groupId: 'dev',
    accent: 'azure',
    glyph: '<>',
  },
];
