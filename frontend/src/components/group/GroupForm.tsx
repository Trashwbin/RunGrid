import {useEffect, useState} from 'react';
import {Icon} from '../ui/Icon';
import {
  GROUP_ICON_OPTIONS,
  type GroupIconValue,
} from '../../utils/groupIcons';
import './GroupForm.css';

export type GroupDraft = {
  name: string;
  icon: GroupIconValue;
};

type GroupFormProps = {
  initialDraft: GroupDraft;
  onChange: (next: GroupDraft) => void;
};

export function GroupForm({initialDraft, onChange}: GroupFormProps) {
  const [draft, setDraft] = useState<GroupDraft>(initialDraft);

  useEffect(() => {
    onChange(draft);
  }, [draft, onChange]);

  return (
    <div className="group-form">
      <label className="edit-field">
        <span className="edit-label">分组名称</span>
        <input
          type="text"
          className="edit-input"
          value={draft.name}
          onChange={(event) =>
            setDraft((prev) => ({...prev, name: event.target.value}))
          }
          placeholder="输入分组名称"
        />
      </label>
      <div className="group-icon-field">
        <span className="edit-label">图标</span>
        <div className="group-icon-grid" role="listbox" aria-label="选择分组图标">
          {GROUP_ICON_OPTIONS.map((option) => {
            const isSelected = draft.icon === option.value;
            return (
              <button
                key={option.value || 'none'}
                type="button"
                className={`group-icon-option${isSelected ? ' is-selected' : ''}`}
                onClick={() =>
                  setDraft((prev) => ({...prev, icon: option.value}))
                }
                aria-pressed={isSelected}
              >
                <Icon name={option.icon} size={18} />
              </button>
            );
          })}
        </div>
      </div>
    </div>
  );
}
