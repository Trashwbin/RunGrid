import type {Category, MenuItem} from '../types';

export const categories: Category[] = [
  {id: 'apps', label: '应用程序', icon: 'apps'},
  {id: 'system', label: '系统应用', icon: 'system'},
  {id: 'urls', label: '网址链接', icon: 'link'},
  {id: 'docs', label: '书籍文档', icon: 'docs'},
  {id: 'folders', label: '文件夹', icon: 'folder'},
];

export const menuItems: MenuItem[] = [
  {id: 'settings', label: '设置'},
  {id: 'scan', label: '扫描快捷方式'},
  {id: 'import-rules', label: '导入分组规则'},
  {id: 'sync-icons', label: '刷新图标缓存'},
  {id: 'clear', label: '清空项目'},
];
