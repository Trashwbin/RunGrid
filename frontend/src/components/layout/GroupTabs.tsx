import {useCallback, useRef, useState} from 'react';
import type {GroupTab} from '../../types';
import {Icon} from '../ui/Icon';

type GroupTabsProps = {
  tabs: GroupTab[];
  activeId: string;
  onSelect: (id: string) => void;
  onAdd?: () => void;
  onOpenMenu?: (id: string, x: number, y: number) => void;
  onReorder?: (sourceId: string, targetId: string) => void;
};

export function GroupTabs({
  tabs,
  activeId,
  onSelect,
  onAdd,
  onOpenMenu,
  onReorder,
}: GroupTabsProps) {
  const [draggingId, setDraggingId] = useState<string | null>(null);
  const listRef = useRef<HTMLDivElement>(null);

  const handleWheel = useCallback((event: React.WheelEvent<HTMLDivElement>) => {
    const container = listRef.current;
    if (!container) {
      return;
    }
    const maxScroll = container.scrollWidth - container.clientWidth;
    if (maxScroll <= 0) {
      return;
    }
    const delta = Math.abs(event.deltaX) > Math.abs(event.deltaY)
      ? event.deltaX
      : event.deltaY;
    if (delta === 0) {
      return;
    }
    const next = Math.min(maxScroll, Math.max(0, container.scrollLeft + delta));
    if (next === container.scrollLeft) {
      return;
    }
    container.scrollLeft = next;
    event.preventDefault();
  }, []);

  return (
    <div className="group-tabs">
      <div className="group-tabs__list" ref={listRef} onWheel={handleWheel}>
        {tabs.map((tab) => (
          <button
            key={tab.id}
            type="button"
            className={`group-tab${activeId === tab.id ? ' is-active' : ''}${
              draggingId === tab.id ? ' is-dragging' : ''
            }`}
            onClick={() => onSelect(tab.id)}
            onContextMenu={(event) => {
              if (!onOpenMenu || tab.id === 'all') {
                return;
              }
              event.preventDefault();
              onOpenMenu(tab.id, event.clientX, event.clientY);
            }}
            draggable={Boolean(onReorder) && tab.id !== 'all'}
            onDragStart={(event) => {
              if (!onReorder || tab.id === 'all') {
                return;
              }
              event.dataTransfer.setData('text/plain', tab.id);
              event.dataTransfer.effectAllowed = 'move';
              setDraggingId(tab.id);
            }}
            onDragEnd={() => {
              setDraggingId(null);
            }}
            onDragOver={(event) => {
              if (!onReorder || tab.id === 'all') {
                return;
              }
              event.preventDefault();
              event.dataTransfer.dropEffect = 'move';
            }}
            onDrop={(event) => {
              if (!onReorder || tab.id === 'all') {
                return;
              }
              event.preventDefault();
              const sourceId = event.dataTransfer.getData('text/plain');
              if (!sourceId || sourceId === tab.id) {
                return;
              }
              onReorder(sourceId, tab.id);
            }}
          >
            {tab.icon ? <Icon name={tab.icon} size={14} /> : null}
            <span>{tab.label}</span>
          </button>
        ))}
      </div>
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
