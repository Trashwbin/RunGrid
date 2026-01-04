export type PanelPositionMode = 'center' | 'last' | 'cursor';
export type LaunchMode = 'single' | 'double';

export type Preferences = {
  focusSearchOnShow: boolean;
  panelPositionMode: PanelPositionMode;
  launchMode: LaunchMode;
  lastWindowPosition?: {x: number; y: number};
};

const STORAGE_KEY = 'rungrid.preferences';

export const DEFAULT_PREFERENCES: Preferences = {
  focusSearchOnShow: true,
  panelPositionMode: 'center',
  launchMode: 'single',
};

const PANEL_POSITION_MODES: PanelPositionMode[] = ['center', 'last', 'cursor'];
const LAUNCH_MODES: LaunchMode[] = ['single', 'double'];

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
      if (
        typeof parsed.panelPositionMode === 'string' &&
        PANEL_POSITION_MODES.includes(parsed.panelPositionMode as PanelPositionMode)
      ) {
        next.panelPositionMode = parsed.panelPositionMode as PanelPositionMode;
      }
      if (
        typeof parsed.launchMode === 'string' &&
        LAUNCH_MODES.includes(parsed.launchMode as LaunchMode)
      ) {
        next.launchMode = parsed.launchMode as LaunchMode;
      }
      if (
        parsed.lastWindowPosition &&
        typeof parsed.lastWindowPosition === 'object' &&
        typeof parsed.lastWindowPosition.x === 'number' &&
        typeof parsed.lastWindowPosition.y === 'number'
      ) {
        next.lastWindowPosition = {
          x: parsed.lastWindowPosition.x,
          y: parsed.lastWindowPosition.y,
        };
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
