import type {domain} from '../../wailsjs/go/models';
import type {Accent, AppItem, GroupTab} from '../types';

const accentPalette: Accent[] = [
  'pink',
  'blue',
  'orange',
  'teal',
  'azure',
  'indigo',
];

export function toGroupTab(group: domain.Group): GroupTab {
  return {
    id: group.id,
    label: group.name,
  };
}

export function toAppItem(item: domain.Item, index: number): AppItem {
  return {
    id: item.id,
    name: item.name,
    categoryId: mapTypeToCategory(item.type),
    groupId: item.group_id || 'all',
    type: item.type as AppItem['type'],
    path: item.path,
    accent: accentPalette[index % accentPalette.length],
    glyph: makeGlyph(item.name),
    iconUrl: toIconURL(item.icon_path),
    favorite: item.favorite,
  };
}

function mapTypeToCategory(type: string): string {
  switch (type) {
    case 'folder':
      return 'folders';
    case 'doc':
      return 'docs';
    case 'url':
      return 'urls';
    case 'system':
      return 'system';
    case 'app':
    default:
      return 'apps';
  }
}

function makeGlyph(name: string): string {
  const trimmed = name.trim();
  if (!trimmed) {
    return '?';
  }

  const parts = trimmed.split(/\s+/);
  if (parts.length === 1) {
    return parts[0].slice(0, 2).toUpperCase();
  }

  return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
}

function toIconURL(path: string | undefined): string | undefined {
  if (!path) {
    return undefined;
  }

  const normalized = path.replace(/\\/g, '/');
  const fileName = normalized.split('/').pop();
  if (!fileName) {
    return undefined;
  }

  return `/icons/${encodeURIComponent(fileName)}`;
}
