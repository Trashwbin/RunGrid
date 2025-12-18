import type {AppItem} from '../../types';
import {AppTile} from './AppTile';

type AppGridProps = {
  items: AppItem[];
};

export function AppGrid({items}: AppGridProps) {
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
