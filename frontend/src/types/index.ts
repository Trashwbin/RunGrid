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
  | 'heart';

export type Category = {
  id: string;
  label: string;
  icon: IconName;
};

export type GroupTab = {
  id: string;
  label: string;
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
