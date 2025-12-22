import type {GroupTab} from '../../types';
import {Icon} from '../ui/Icon';

type GroupTabsProps = {
  tabs: GroupTab[];
  activeId: string;
  onSelect: (id: string) => void;
  onAdd?: () => void;
};

export function GroupTabs({tabs, activeId, onSelect, onAdd}: GroupTabsProps) {
  return (
    <div className="group-tabs">
      {tabs.map((tab) => (
        <button
          key={tab.id}
          type="button"
          className={`group-tab${activeId === tab.id ? ' is-active' : ''}`}
          onClick={() => onSelect(tab.id)}
        >
          {tab.icon ? <Icon name={tab.icon} size={14} /> : null}
          <span>{tab.label}</span>
        </button>
      ))}
      <button
        type="button"
        className="group-tab add-tab"
        onClick={onAdd}
        aria-label="添加分组"
      >
        <Icon name="plus" size={14} />
      </button>
    </div>
  );
}
