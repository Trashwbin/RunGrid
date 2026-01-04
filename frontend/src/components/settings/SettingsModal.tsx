import {useCallback, useEffect, useMemo, useState} from 'react';
import {ScrollArea} from '../ui/ScrollArea';
import {
  DEFAULT_HOTKEYS,
  HOTKEY_ACTIONS,
  type HotkeyConfig,
} from '../../utils/hotkeys';
import {
  DEFAULT_PREFERENCES,
  type LaunchMode,
  type PanelPositionMode,
  type Preferences,
} from '../../utils/preferences';
import './SettingsModal.css';

type SettingsModalProps = {
  initialHotkeys: HotkeyConfig;
  onChange: (next: HotkeyConfig) => void;
  initialPreferences: Preferences;
  onPreferencesChange: (next: Preferences) => void;
};

type TabId = 'general' | 'scan' | 'hotkeys' | 'storage';

const tabs: {id: TabId; label: string}[] = [
  {id: 'general', label: '通用设置'},
  {id: 'scan', label: '文件扫描'},
  {id: 'hotkeys', label: '快捷键'},
  {id: 'storage', label: '数据存储'},
];

const modifierKeys = new Set(['Shift', 'Control', 'Alt', 'Meta']);
const panelPositionOptions: Array<{value: PanelPositionMode; label: string}> = [
  {value: 'center', label: '居中显示'},
  {value: 'last', label: '上次位置'},
  {value: 'cursor', label: '跟随鼠标'},
];
const launchModeOptions: Array<{value: LaunchMode; label: string}> = [
  {value: 'single', label: '单击启动'},
  {value: 'double', label: '双击启动'},
];

export function SettingsModal({
  initialHotkeys,
  onChange,
  initialPreferences,
  onPreferencesChange,
}: SettingsModalProps) {
  const [activeTab, setActiveTab] = useState<TabId>('hotkeys');
  const [hotkeys, setHotkeys] = useState<HotkeyConfig>({
    ...DEFAULT_HOTKEYS,
    ...initialHotkeys,
  });
  const [preferences, setPreferences] = useState<Preferences>({
    ...DEFAULT_PREFERENCES,
    ...initialPreferences,
  });
  const [recordingId, setRecordingId] = useState<string | null>(null);

  useEffect(() => {
    setHotkeys({...DEFAULT_HOTKEYS, ...initialHotkeys});
  }, [initialHotkeys]);

  useEffect(() => {
    setPreferences({...DEFAULT_PREFERENCES, ...initialPreferences});
  }, [initialPreferences]);

  useEffect(() => {
    onChange(hotkeys);
  }, [hotkeys, onChange]);

  useEffect(() => {
    onPreferencesChange(preferences);
  }, [onPreferencesChange, preferences]);

  useEffect(() => {
    if (!recordingId) {
      return;
    }

    const handleKeydown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setRecordingId(null);
        return;
      }

      const next = normalizeHotkey(event);
      if (!next) {
        return;
      }
      event.preventDefault();
      setHotkeys((prev) => ({...prev, [recordingId]: next}));
      setRecordingId(null);
    };

    window.addEventListener('keydown', handleKeydown);
    return () => {
      window.removeEventListener('keydown', handleKeydown);
    };
  }, [recordingId]);

  useEffect(() => {
    if (!recordingId) {
      return;
    }

    const handlePointerMove = (event: PointerEvent) => {
      if (event.buttons === 0) {
        return;
      }
      setRecordingId(null);
    };

    window.addEventListener('pointermove', handlePointerMove);
    return () => {
      window.removeEventListener('pointermove', handlePointerMove);
    };
  }, [recordingId]);

  const handleReset = useCallback(() => {
    setHotkeys({...DEFAULT_HOTKEYS});
  }, []);

  const handleClear = useCallback((id: string) => {
    setHotkeys((prev) => ({...prev, [id]: ''}));
  }, []);

  const handleRecord = useCallback((id: string) => {
    setRecordingId((prev) => (prev === id ? null : id));
  }, []);

  const activeContent = useMemo(() => {
    if (activeTab === 'general') {
      return (
        <div className="settings-general">
          <div className="settings-section-header">
            <div>
              <h3 className="settings-section-title">通用设置</h3>
              <p className="settings-section-desc">
                控制唤出面板时的默认交互。
              </p>
            </div>
          </div>
          <div className="settings-option">
            <div className="settings-option-info">
              <div className="settings-option-title">唤出时聚焦搜索框</div>
              <div className="settings-option-desc">
                显示面板后自动定位到搜索输入框。
              </div>
            </div>
            <label className="settings-switch" aria-label="唤出时聚焦搜索框">
              <input
                type="checkbox"
                checked={preferences.focusSearchOnShow}
                onChange={() =>
                  setPreferences((prev) => ({
                    ...prev,
                    focusSearchOnShow: !prev.focusSearchOnShow,
                  }))
                }
              />
              <span className="settings-switch-track">
                <span className="settings-switch-thumb" />
              </span>
            </label>
          </div>
          <div className="settings-option">
            <div className="settings-option-info">
              <div className="settings-option-title">面板呼出位置</div>
              <div className="settings-option-desc">
                选择唤出面板时窗口的显示位置。
              </div>
            </div>
            <div
              className="settings-choice-group"
              role="radiogroup"
              aria-label="面板呼出位置"
            >
              {panelPositionOptions.map((option) => (
                <button
                  key={option.value}
                  type="button"
                  className={`settings-choice${preferences.panelPositionMode === option.value ? ' is-active' : ''}`}
                  aria-pressed={preferences.panelPositionMode === option.value}
                  onClick={() =>
                    setPreferences((prev) => ({
                      ...prev,
                      panelPositionMode: option.value,
                    }))
                  }
                >
                  {option.label}
                </button>
              ))}
            </div>
          </div>
          <div className="settings-option">
            <div className="settings-option-info">
              <div className="settings-option-title">启动方式</div>
              <div className="settings-option-desc">
                选择点击图标时的启动行为。
              </div>
            </div>
            <div
              className="settings-choice-group"
              role="radiogroup"
              aria-label="启动方式"
            >
              {launchModeOptions.map((option) => (
                <button
                  key={option.value}
                  type="button"
                  className={`settings-choice${preferences.launchMode === option.value ? ' is-active' : ''}`}
                  aria-pressed={preferences.launchMode === option.value}
                  onClick={() =>
                    setPreferences((prev) => ({
                      ...prev,
                      launchMode: option.value,
                    }))
                  }
                >
                  {option.label}
                </button>
              ))}
            </div>
          </div>
        </div>
      );
    }

    if (activeTab === 'hotkeys') {
      return (
        <div className="settings-hotkeys">
          <div className="settings-section-header">
            <div>
              <h3 className="settings-section-title">快捷键设置</h3>
              <p className="settings-section-desc">
                点击录入并按下组合键，保存后会记录到本地。
              </p>
            </div>
            <div className="settings-section-actions">
              <button
                type="button"
                className="settings-ghost-button"
                onClick={handleReset}
              >
                恢复默认
              </button>
            </div>
          </div>
          <ScrollArea
            className="settings-hotkey-scroll"
            viewportClassName="settings-hotkey-scroll__viewport"
            contentClassName="settings-hotkey-scroll__content"
          >
            <div className="hotkey-list">
              {HOTKEY_ACTIONS.map((action) => {
                const value = hotkeys[action.id];
                const isRecording = recordingId === action.id;
                return (
                  <div className="hotkey-row" key={action.id}>
                    <div className="hotkey-info">
                      <div className="hotkey-label">{action.label}</div>
                      <div className="hotkey-desc">{action.description}</div>
                    </div>
                    <div className="hotkey-control">
                      <button
                        type="button"
                        className={`hotkey-capture${isRecording ? ' is-recording' : ''}`}
                        onClick={() => handleRecord(action.id)}
                      >
                        {isRecording
                          ? '按下组合键...'
                          : value || '点击录入'}
                      </button>
                      <button
                        type="button"
                        className="hotkey-clear"
                        onClick={() => handleClear(action.id)}
                      >
                        清除
                      </button>
                    </div>
                  </div>
                );
              })}
            </div>
          </ScrollArea>
        </div>
      );
    }

    return (
      <div className="settings-placeholder">
        <div className="settings-placeholder-title">功能即将上线</div>
        <p className="settings-placeholder-desc">这一模块正在打磨中。</p>
      </div>
    );
  }, [
    activeTab,
    handleClear,
    handleRecord,
    handleReset,
    hotkeys,
    recordingId,
    preferences,
  ]);

  return (
    <div className="settings-modal">
      <div className="settings-tabs">
        {tabs.map((tab) => (
          <button
            type="button"
            key={tab.id}
            className={`settings-tab${activeTab === tab.id ? ' is-active' : ''}`}
            onClick={() => setActiveTab(tab.id)}
          >
            {tab.label}
          </button>
        ))}
      </div>
      <div className="settings-content">{activeContent}</div>
    </div>
  );
}

function normalizeHotkey(event: KeyboardEvent) {
  const parts: string[] = [];
  if (event.ctrlKey) {
    parts.push('Ctrl');
  }
  if (event.altKey) {
    parts.push('Alt');
  }
  if (event.shiftKey) {
    parts.push('Shift');
  }
  if (event.metaKey) {
    parts.push('Win');
  }

  const key = normalizeKey(event.key);
  if (!key || modifierKeys.has(key)) {
    return '';
  }
  parts.push(key);
  return parts.join('+');
}

function normalizeKey(key: string) {
  if (!key) {
    return '';
  }
  if (key === ' ') {
    return 'Space';
  }
  if (key === 'ArrowUp') {
    return 'Up';
  }
  if (key === 'ArrowDown') {
    return 'Down';
  }
  if (key === 'ArrowLeft') {
    return 'Left';
  }
  if (key === 'ArrowRight') {
    return 'Right';
  }
  if (key.length === 1) {
    return key.toUpperCase();
  }
  return key;
}
