export type HotkeyConfig = Record<string, string>;

const STORAGE_KEY = 'rungrid.hotkeys';

export const DEFAULT_HOTKEYS: HotkeyConfig = {
  'toggle-app': 'Alt+Space',
  'quick-search': 'Alt+F',
  'open-settings': 'Ctrl+,',
};

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
