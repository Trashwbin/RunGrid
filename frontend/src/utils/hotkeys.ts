export type HotkeyConfig = Record<string, string>;

export type HotkeyAction = {
  id: string;
  label: string;
  description: string;
};

const STORAGE_KEY = 'rungrid.hotkeys';

export const DEFAULT_HOTKEYS: HotkeyConfig = {
  'toggle-app': 'Alt+Space',
};

export const HOTKEY_ACTIONS: HotkeyAction[] = [
  {
    id: 'toggle-app',
    label: '唤出/隐藏主窗口',
    description: '快速显示或隐藏 RunGrid。',
  },
];

export function loadHotkeys(): HotkeyConfig {
  const next = {...DEFAULT_HOTKEYS};
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      return next;
    }
    const parsed = JSON.parse(raw);
    if (parsed && typeof parsed === 'object') {
      for (const [key, value] of Object.entries(parsed)) {
        if (typeof value === 'string') {
          next[key] = value;
        }
      }
    }
  } catch {
    return next;
  }
  return next;
}

export function saveHotkeys(value: HotkeyConfig) {
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(value));
}

export function toHotkeyBindings(config: HotkeyConfig) {
  return HOTKEY_ACTIONS.map((action) => ({
    id: action.id,
    keys: config[action.id] ?? '',
  }));
}
