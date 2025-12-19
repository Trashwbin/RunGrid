import {useEffect, useMemo, useState} from 'react';
import {
  PickIconSource,
  PickTargetFolder,
  PickTargetPath,
  PreviewIconFromSource,
} from '../../../wailsjs/go/main/App';
import {Icon} from '../ui/Icon';

export type EditDraft = {
  id: string;
  name: string;
  path: string;
  originalPath: string;
  iconUrl?: string;
  glyph: string;
  iconSource?: string;
  favorite: boolean;
  hidden: boolean;
};

type EditItemFormProps = {
  initialDraft: EditDraft;
  onChange: (next: EditDraft) => void;
};

export function EditItemForm({initialDraft, onChange}: EditItemFormProps) {
  const [draft, setDraft] = useState<EditDraft>(initialDraft);

  useEffect(() => {
    onChange(draft);
  }, [draft, onChange]);

  const iconLabel = useMemo(() => {
    if (!draft.iconSource) {
      return '跟随目标图标';
    }
    const normalized = draft.iconSource.replace(/\\/g, '/');
    const parts = normalized.split('/');
    return parts[parts.length - 1] || draft.iconSource;
  }, [draft.iconSource]);

  const handlePickIcon = async () => {
    try {
      const selected = await PickIconSource();
      const trimmed = selected.trim();
      if (!trimmed) {
        return;
      }
      const previewPath = await PreviewIconFromSource(trimmed);
      const previewUrl = toIconURL(previewPath);
      setDraft((prev) => ({
        ...prev,
        iconSource: trimmed,
        iconUrl: previewUrl ?? prev.iconUrl,
      }));
    } catch {
    }
  };

  const handlePickTargetPath = async () => {
    try {
      const selected = await PickTargetPath();
      const trimmed = selected.trim();
      if (!trimmed) {
        return;
      }
      setDraft((prev) => ({...prev, path: trimmed}));
    } catch {
    }
  };

  const handlePickTargetFolder = async () => {
    try {
      const selected = await PickTargetFolder();
      const trimmed = selected.trim();
      if (!trimmed) {
        return;
      }
      setDraft((prev) => ({...prev, path: trimmed}));
    } catch {
    }
  };

  return (
    <div className="edit-form">
      <label className="edit-field">
        <span className="edit-label">名称</span>
        <input
          type="text"
          className="edit-input"
          value={draft.name}
          onChange={(event) =>
            setDraft((prev) => ({...prev, name: event.target.value}))
          }
          placeholder="输入名称"
        />
      </label>
      <label className="edit-field">
        <span className="edit-label">目标路径</span>
        <div className="edit-input-row">
          <input
            type="text"
            className="edit-input"
            value={draft.path}
            onChange={(event) =>
              setDraft((prev) => ({...prev, path: event.target.value}))
            }
            placeholder="输入 exe/lnk/url 等路径"
          />
          <div className="edit-input-actions">
            <button
              type="button"
              className="edit-input-button"
              onClick={handlePickTargetPath}
            >
              选择文件
            </button>
            <button
              type="button"
              className="edit-input-button"
              onClick={handlePickTargetFolder}
            >
              选择文件夹
            </button>
          </div>
        </div>
      </label>
      <div className="edit-field">
        <span className="edit-label">图标</span>
        <div className="edit-icon-row">
          <div className="edit-icon-preview">
            {draft.iconUrl ? (
              <img src={draft.iconUrl} alt="" />
            ) : (
              <span className="edit-icon-glyph">{draft.glyph}</span>
            )}
          </div>
          <div className="edit-icon-meta">
            <span className="edit-icon-hint">{iconLabel}</span>
            <button type="button" className="edit-icon-button" onClick={handlePickIcon}>
              <Icon name="plus" size={14} />
              选择图标
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

function toIconURL(path: string | undefined) {
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
