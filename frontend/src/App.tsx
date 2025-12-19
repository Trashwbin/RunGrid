import {useCallback, useEffect, useMemo, useState} from 'react';
import './App.css';
import {categories, menuItems} from './data/mock';
import {
  CreateGroup,
  CreateItem,
  DeleteItem,
  LaunchItem,
  ListGroups,
  ListItems,
  OpenItemLocation,
  RefreshItemIcon,
  ScanShortcuts,
  SetFavorite,
  SyncIcons,
} from '../wailsjs/go/main/App';
import type {domain} from '../wailsjs/go/models';
import {AppGrid} from './components/grid/AppGrid';
import {CategoryBar} from './components/layout/CategoryBar';
import {GroupTabs} from './components/layout/GroupTabs';
import {SearchBar} from './components/layout/SearchBar';
import {TopBar} from './components/layout/TopBar';
import {ContextMenu} from './components/ui/ContextMenu';
import {toAppItem, toGroupTab} from './utils/items';
import type {AppItem} from './types';

function App() {
  const [items, setItems] = useState<domain.Item[]>([]);
  const [groups, setGroups] = useState<domain.Group[]>([]);
  const [activeCategoryId, setActiveCategoryId] = useState(categories[0].id);
  const [activeGroupId, setActiveGroupId] = useState('all');
  const [query, setQuery] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [menuState, setMenuState] = useState<{
    open: boolean;
    x: number;
    y: number;
    item: AppItem | null;
  }>({
    open: false,
    x: 0,
    y: 0,
    item: null,
  });

  const loadGroups = useCallback(async () => {
    try {
      const data = await ListGroups();
      setGroups(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : '无法加载分组');
    }
  }, []);

  const loadItems = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await ListItems(activeGroupId, query);
      setItems(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : '无法加载项目');
    } finally {
      setIsLoading(false);
    }
  }, [activeGroupId, query]);

  useEffect(() => {
    loadGroups();
  }, [loadGroups]);

  useEffect(() => {
    loadItems();
  }, [loadItems]);

  const handleAddGroup = useCallback(async () => {
    const name = window.prompt('分组名称');
    if (!name) {
      return;
    }

    try {
      const newGroup = await CreateGroup({
        name,
        order: groups.length,
        color: '#4f7dff',
      });
      setActiveGroupId(newGroup.id);
      await loadGroups();
    } catch (err) {
      setError(err instanceof Error ? err.message : '无法创建分组');
    }
  }, [groups.length, loadGroups]);

  const handleAddItem = useCallback(async () => {
    if (activeGroupId === 'all') {
      window.alert('请先选择一个分组');
      return;
    }

    const name = window.prompt('项目名称');
    if (!name) {
      return;
    }
    const path = window.prompt('目标路径');
    if (!path) {
      return;
    }

    try {
      await CreateItem({
        name,
        path,
        type: mapCategoryToType(activeCategoryId),
        icon_path: '',
        group_id: activeGroupId,
        tags: [],
        favorite: false,
        hidden: false,
      });
      await loadItems();
    } catch (err) {
      setError(err instanceof Error ? err.message : '无法创建项目');
    }
  }, [activeCategoryId, activeGroupId, loadItems]);

  const handleMenuSelect = useCallback(
    async (id: string) => {
      if (id === 'scan') {
        setIsLoading(true);
        setError(null);
        try {
          await ScanShortcuts();
          await loadItems();
        } catch (err) {
          setError(err instanceof Error ? err.message : '扫描失败');
        } finally {
          setIsLoading(false);
        }
        return;
      }

      if (id === 'sync-icons') {
        setIsLoading(true);
        setError(null);
        try {
          await SyncIcons();
          await loadItems();
        } catch (err) {
          setError(err instanceof Error ? err.message : '刷新图标失败');
        } finally {
          setIsLoading(false);
        }
      }
    },
    [loadItems]
  );

  const handleLaunch = useCallback(
    async (id: string) => {
      try {
        await LaunchItem(id);
        await loadItems();
      } catch (err) {
        setError(err instanceof Error ? err.message : '启动失败');
      }
    },
    [loadItems]
  );


  const groupTabs = useMemo(
    () => [{id: 'all', label: '全部'}, ...groups.map(toGroupTab)],
    [groups]
  );

  const appItems = useMemo(
    () => items.map((item, index) => toAppItem(item, index)),
    [items]
  );

  const filteredItems = useMemo(
    () => appItems.filter((item) => item.categoryId === activeCategoryId),
    [appItems, activeCategoryId]
  );

  const handleOpenMenu = useCallback((item: AppItem, x: number, y: number) => {
    setMenuState({open: true, x, y, item});
  }, []);

  const handleCloseMenu = useCallback(() => {
    setMenuState((prev) => ({...prev, open: false}));
  }, []);

  const handleMenuAction = useCallback(
    async (actionId: string) => {
      const current = menuState.item;
      if (!current) {
        return;
      }

      if (actionId === 'favorite') {
        try {
          await SetFavorite(current.id, !current.favorite);
          await loadItems();
        } catch (err) {
          setError(err instanceof Error ? err.message : '更新收藏失败');
        }
      }

      if (actionId === 'open-location') {
        try {
          await OpenItemLocation(current.id);
        } catch (err) {
          setError(err instanceof Error ? err.message : '无法打开所在目录');
        }
      }

      if (actionId === 'refresh-icon') {
        try {
          await RefreshItemIcon(current.id);
          await loadItems();
        } catch (err) {
          setError(err instanceof Error ? err.message : '刷新图标失败');
        }
      }

      if (actionId === 'remove') {
        const confirmed = window.confirm(`确定移除 "${current.name}" 吗？`);
        if (!confirmed) {
          return;
        }
        try {
          await DeleteItem(current.id);
          await loadItems();
        } catch (err) {
          setError(err instanceof Error ? err.message : '移除失败');
        }
      }
    },
    [menuState.item, loadItems]
  );

  const menuItem = menuState.item;
  const hasMenuPath = Boolean(menuItem?.path?.trim());
  const isMenuWeb = menuItem ? isWebPath(menuItem.path) : false;

  return (
    <div className="app-shell">
      <TopBar title="RunGrid" menuItems={menuItems} onMenuSelect={handleMenuSelect} />
      <div className="content">
        <CategoryBar
          categories={categories}
          activeId={activeCategoryId}
          onSelect={setActiveCategoryId}
        />
        <GroupTabs
          tabs={groupTabs}
          activeId={activeGroupId}
          onSelect={setActiveGroupId}
          onAdd={handleAddGroup}
        />
        <AppGrid
          items={filteredItems}
          isLoading={isLoading}
          error={error}
          onAddItem={handleAddItem}
          onLaunch={handleLaunch}
          onOpenMenu={handleOpenMenu}
        />
      </div>
      <SearchBar value={query} onChange={setQuery} />
      <ContextMenu
        open={menuState.open}
        x={menuState.x}
        y={menuState.y}
        items={[
          {
            id: 'favorite',
            label: menuItem?.favorite ? '取消收藏' : '收藏',
          },
          {
            id: 'open-location',
            label: '打开所在目录',
            disabled: !menuItem || !hasMenuPath || isMenuWeb,
          },
          {
            id: 'refresh-icon',
            label: '刷新图标',
            disabled:
              !menuItem || !hasMenuPath || menuItem.type === 'url' || isMenuWeb,
          },
          {
            id: 'remove',
            label: '移除',
            tone: 'danger',
          },
        ]}
        onSelect={handleMenuAction}
        onClose={handleCloseMenu}
      />
    </div>
  );
}

export default App;

function mapCategoryToType(categoryId: string): domain.ItemInput['type'] {
  switch (categoryId) {
    case 'folders':
      return 'folder';
    case 'docs':
      return 'doc';
    case 'system':
      return 'app';
    case 'apps':
    default:
      return 'app';
  }
}

function isWebPath(path: string) {
  const trimmed = path.trim().toLowerCase();
  return trimmed.startsWith('http://') || trimmed.startsWith('https://');
}
