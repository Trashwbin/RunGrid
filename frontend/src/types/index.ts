export type IconName =
  | 'menu'
  | 'settings'
  | 'close'
  | 'apps'
  | 'system'
  | 'docs'
  | 'folder'
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

export type AppItem = {
  id: string;
  name: string;
  categoryId: string;
  groupId: string;
  accent: Accent;
  glyph: string;
  iconPath?: string;
};
