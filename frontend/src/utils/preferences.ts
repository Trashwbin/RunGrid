export type Preferences = {
  focusSearchOnShow: boolean;
};

const STORAGE_KEY = 'rungrid.preferences';

export const DEFAULT_PREFERENCES: Preferences = {
  focusSearchOnShow: true,
};

export function loadPreferences(): Preferences {
  const next = {...DEFAULT_PREFERENCES};
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      return next;
    }
    const parsed = JSON.parse(raw);
    if (parsed && typeof parsed === 'object') {
      if (typeof parsed.focusSearchOnShow === 'boolean') {
        next.focusSearchOnShow = parsed.focusSearchOnShow;
      }
    }
  } catch {
    return next;
  }
  return next;
}

export function savePreferences(value: Preferences) {
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(value));
}
