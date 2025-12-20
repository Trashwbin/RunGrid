import {useCallback, useEffect, useMemo, useState} from 'react';
import {ScrollArea} from '../ui/ScrollArea';
import {DEFAULT_HOTKEYS, type HotkeyConfig} from '../../utils/hotkeys';
import './SettingsModal.css';

type SettingsModalProps = {
  initialHotkeys: HotkeyConfig;
  onChange: (next: HotkeyConfig) => void;
};

type TabId = 'general' | 'scan' | 'hotkeys' | 'storage' | 'vip';

const tabs: {id: TabId; label: string}[] = [
  {id: 'general', label: '通用设置'},
  {id: 'scan', label: '文件扫描'},
  {id: 'hotkeys', label: '快捷键'},
  {id: 'storage', label: '数据存储'},
  {id: 'vip', label: '我的会员'},
];

const hotkeyActions = [
  {
    id: 'toggle-app',
    label: '唤出/隐藏主窗口',
    description: '快速显示或隐藏 RunGrid。',
  },
  {
    id: 'quick-search',
    label: '聚焦搜索框',
    description: '在任意界面聚焦搜索输入。',
  },
  {
    id: 'open-settings',
    label: '打开设置面板',
    description: '快速打开设置窗口。',
  },
];

const modifierKeys = new Set(['Shift', 'Control', 'Alt', 'Meta']);

export function SettingsModal({initialHotkeys, onChange}: SettingsModalProps) {
  const [activeTab, setActiveTab] = useState<TabId>('hotkeys');
  const [hotkeys, setHotkeys] = useState<HotkeyConfig>({
    ...DEFAULT_HOTKEYS,
    ...initialHotkeys,
  });
  const [recordingId, setRecordingId] = useState<string | null>(null);

  useEffect(() => {
    setHotkeys({...DEFAULT_HOTKEYS, ...initialHotkeys});
  }, [initialHotkeys]);

  useEffect(() => {
    onChange(hotkeys);
  }, [hotkeys, onChange]);

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
              {hotkeyActions.map((action) => {
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
          <div className="settings-footnote">
            当前仅保存热键配置，实际全局热键将在后续接入。
          </div>
        </div>
      );
    }

    return (
      <div className="settings-placeholder">
        <div className="settings-placeholder-title">功能即将上线</div>
        <p className="settings-placeholder-desc">这一模块正在打磨中。</p>
      </div>
    );
  }, [activeTab, handleClear, handleRecord, handleReset, hotkeys, recordingId]);

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
