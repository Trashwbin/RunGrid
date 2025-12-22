export type IconName =
  | 'menu'
  | 'settings'
  | 'close'
  | 'apps'
  | 'system'
  | 'docs'
  | 'folder'
  | 'link'
  | 'plus'
  | 'heart'
  | 'ban'
  | 'star'
  | 'home'
  | 'download'
  | 'upload'
  | 'cloud'
  | 'shield'
  | 'code'
  | 'terminal'
  | 'palette'
  | 'calendar'
  | 'camera'
  | 'mail'
  | 'phone'
  | 'cart'
  | 'map'
  | 'chat'
  | 'bookmark'
  | 'bolt'
  | 'key'
  | 'lock'
  | 'chip'
  | 'rocket'
  | 'image'
  | 'puzzle'
  | 'globe'
  | 'database'
  | 'gamepad'
  | 'video'
  | 'briefcase'
  | 'tools'
  | 'file'
  | 'music'
  | 'user';

export type Category = {
  id: string;
  label: string;
  icon: IconName;
};

export type GroupTab = {
  id: string;
  label: string;
  icon?: IconName;
};

export type MenuItem = {
  id: string;
  label: string;
};

export type Accent = 'pink' | 'blue' | 'orange' | 'teal' | 'azure' | 'indigo';

export type AppItemType = 'app' | 'system' | 'doc' | 'folder' | 'url';

export type AppItem = {
  id: string;
  name: string;
  categoryId: string;
  groupId: string;
  type: AppItemType;
  path: string;
  accent: Accent;
  glyph: string;
  iconUrl?: string;
  tags: string[];
  favorite: boolean;
  hidden: boolean;
};
