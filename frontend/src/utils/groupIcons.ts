import type {IconName} from '../types';

export type GroupIconValue = IconName | '';

export type GroupIconOption = {
  value: GroupIconValue;
  label: string;
  icon: IconName;
};

export const GROUP_ICON_OPTIONS: GroupIconOption[] = [
  {value: '', label: '无', icon: 'ban'},
  {value: 'apps', label: '应用', icon: 'apps'},
  {value: 'system', label: '系统', icon: 'system'},
  {value: 'docs', label: '文档', icon: 'docs'},
  {value: 'folder', label: '文件夹', icon: 'folder'},
  {value: 'link', label: '链接', icon: 'link'},
  {value: 'heart', label: '收藏', icon: 'heart'},
  {value: 'star', label: '星标', icon: 'star'},
  {value: 'home', label: '主页', icon: 'home'},
  {value: 'download', label: '下载', icon: 'download'},
  {value: 'upload', label: '上传', icon: 'upload'},
  {value: 'cloud', label: '云', icon: 'cloud'},
  {value: 'shield', label: '安全', icon: 'shield'},
  {value: 'code', label: '开发', icon: 'code'},
  {value: 'terminal', label: '终端', icon: 'terminal'},
  {value: 'palette', label: '设计', icon: 'palette'},
  {value: 'calendar', label: '日历', icon: 'calendar'},
  {value: 'camera', label: '相机', icon: 'camera'},
  {value: 'mail', label: '邮件', icon: 'mail'},
  {value: 'phone', label: '电话', icon: 'phone'},
  {value: 'cart', label: '购物', icon: 'cart'},
  {value: 'map', label: '地图', icon: 'map'},
  {value: 'chat', label: '聊天', icon: 'chat'},
  {value: 'bookmark', label: '书签', icon: 'bookmark'},
  {value: 'bolt', label: '效率', icon: 'bolt'},
  {value: 'key', label: '钥匙', icon: 'key'},
  {value: 'lock', label: '锁', icon: 'lock'},
  {value: 'chip', label: '硬件', icon: 'chip'},
  {value: 'rocket', label: '启动', icon: 'rocket'},
  {value: 'image', label: '图片', icon: 'image'},
  {value: 'puzzle', label: '组件', icon: 'puzzle'},
  {value: 'globe', label: '网络', icon: 'globe'},
  {value: 'database', label: '数据库', icon: 'database'},
  {value: 'gamepad', label: '游戏', icon: 'gamepad'},
  {value: 'video', label: '视频', icon: 'video'},
  {value: 'briefcase', label: '工作', icon: 'briefcase'},
  {value: 'tools', label: '工具', icon: 'tools'},
  {value: 'file', label: '文件', icon: 'file'},
  {value: 'music', label: '音乐', icon: 'music'},
  {value: 'user', label: '用户', icon: 'user'},
];

const groupIconSet = new Set<IconName>(
  GROUP_ICON_OPTIONS.map((option) => option.value).filter(
    (value): value is IconName => value !== ''
  )
);

export function toGroupIconName(value?: string): IconName | undefined {
  if (!value) {
    return undefined;
  }
  return groupIconSet.has(value as IconName) ? (value as IconName) : undefined;
}
