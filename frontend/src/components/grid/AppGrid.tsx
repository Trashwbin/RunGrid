import type {AppItem} from '../../types';
import {AppTile} from './AppTile';

type AppGridProps = {
  items: AppItem[];
  isLoading?: boolean;
  error?: string | null;
};

export function AppGrid({items, isLoading = false, error = null}: AppGridProps) {
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
      </div>
    );
  }

  return (
    <div className="app-grid">
      {items.map((item) => (
        <AppTile key={item.id} item={item} />
      ))}
    </div>
  );
}
