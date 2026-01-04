import type {AppItem} from '../../types';
import {AppTile} from './AppTile';

type AppGridProps = {
  items: AppItem[];
  isLoading?: boolean;
  error?: string | null;
  onAddItem?: () => void;
  onLaunch?: (id: string) => void;
  onOpenMenu?: (item: AppItem, x: number, y: number) => void;
  onOpenGridMenu?: (x: number, y: number) => void;
  selectedIds?: string[];
  onSelectItem?: (id: string, multi: boolean) => void;
  selectionMode?: boolean;
  launchMode?: 'single' | 'double';
  focusedId?: string | null;
  onFocusItem?: (id: string) => void;
  onClearFocus?: () => void;
};

export function AppGrid({
  items,
  isLoading = false,
  error = null,
  onAddItem,
  onLaunch,
  onOpenMenu,
  onOpenGridMenu,
  selectedIds = [],
  onSelectItem,
  selectionMode = false,
  launchMode = 'single',
  focusedId = null,
  onFocusItem,
  onClearFocus,
}: AppGridProps) {
  if (error) {
    return (
      <div className="empty-state">
        <p>加载失败</p>
        <span>{error}</span>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="empty-state">
        <p>加载中</p>
        <span>正在同步应用列表...</span>
      </div>
    );
  }

  if (items.length === 0) {
    return (
      <div
        className="empty-state"
        onContextMenu={(event) => {
          if (!onOpenGridMenu) {
            return;
          }
          event.preventDefault();
          onOpenGridMenu(event.clientX, event.clientY);
        }}
      >
        <p>暂无项目</p>
        <span>试试更换分组或搜索</span>
        {onAddItem ? (
          <button type="button" className="empty-action" onClick={onAddItem}>
            新增项目
          </button>
        ) : null}
      </div>
    );
  }

  return (
    <div
      className="app-grid"
      onClick={(event) => {
        if (!onClearFocus) {
          return;
        }
        if (event.target !== event.currentTarget) {
          return;
        }
        onClearFocus();
      }}
      onContextMenu={(event) => {
        if (!onOpenGridMenu) {
          return;
        }
        if (event.target !== event.currentTarget) {
          return;
        }
        event.preventDefault();
        onOpenGridMenu(event.clientX, event.clientY);
      }}
    >
      {items.map((item) => (
        <AppTile
          key={item.id}
          item={item}
          selected={selectedIds.includes(item.id)}
          focused={focusedId === item.id}
          selectionMode={selectionMode}
          launchMode={launchMode}
          onLaunch={onLaunch}
          onOpenMenu={onOpenMenu}
          onSelect={onSelectItem}
          onFocus={onFocusItem}
        />
      ))}
    </div>
  );
}
