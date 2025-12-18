import type {AppItem} from '../../types';
import {AppTile} from './AppTile';

type AppGridProps = {
  items: AppItem[];
  isLoading?: boolean;
  error?: string | null;
  onAddItem?: () => void;
  onLaunch?: (id: string) => void;
};

export function AppGrid({
  items,
  isLoading = false,
  error = null,
  onAddItem,
  onLaunch,
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
      <div className="empty-state">
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
    <div className="app-grid">
      {items.map((item) => (
        <AppTile key={item.id} item={item} onLaunch={onLaunch} />
      ))}
    </div>
  );
}
