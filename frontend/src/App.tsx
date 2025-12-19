import {useCallback, useEffect, useMemo, useState} from 'react';
import './App.css';
import {categories, menuItems} from './data/mock';
import {
  ClearItems,
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
import {ModalHost} from './components/overlay/ModalHost';
import {ToastHost} from './components/overlay/ToastHost';
import {ContextMenu} from './components/ui/ContextMenu';
import {useModalStore, useToastStore} from './store/overlays';
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
  const notify = useToastStore((state) => state.notify);
  const openModal = useModalStore((state) => state.openModal);
  const closeModal = useModalStore((state) => state.closeModal);
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

  const showError = useCallback(
    (message: string, title = '操作失败') => {
      setError(message);
      notify({type: 'error', title, message});
    },
    [notify]
  );

  const loadGroups = useCallback(async () => {
    try {
      const data = await ListGroups();
      setGroups(data);
    } catch (err) {
      showError(err instanceof Error ? err.message : '无法加载分组', '加载失败');
    }
  }, [showError]);

  const loadItems = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await ListItems(activeGroupId, query);
      setItems(data);
    } catch (err) {
      showError(err instanceof Error ? err.message : '无法加载项目', '加载失败');
    } finally {
      setIsLoading(false);
    }
  }, [activeGroupId, query, showError]);

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
      notify({type: 'success', title: '分组已创建', message: name});
    } catch (err) {
      showError(err instanceof Error ? err.message : '无法创建分组');
    }
  }, [groups.length, loadGroups, notify, showError]);

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
      notify({type: 'success', title: '项目已添加', message: name});
    } catch (err) {
      showError(err instanceof Error ? err.message : '无法创建项目');
    }
  }, [activeCategoryId, activeGroupId, loadItems, notify, showError]);

  const handleMenuSelect = useCallback(
    async (id: string) => {
      if (id === 'scan') {
        const modalId = openModal({
          kind: 'progress',
          title: '正在扫描',
          description: '正在读取快捷方式与系统应用，请稍候…',
          closable: false,
          backdropClose: false,
        });
        setIsLoading(true);
        setError(null);
        try {
          await ScanShortcuts();
          await loadItems();
          notify({type: 'success', title: '扫描完成', message: '应用列表已更新'});
        } catch (err) {
          showError(err instanceof Error ? err.message : '扫描失败', '扫描失败');
        } finally {
          setIsLoading(false);
          closeModal(modalId);
        }
        return;
      }

      if (id === 'sync-icons') {
        const modalId = openModal({
          kind: 'progress',
          title: '正在刷新图标',
          description: '图标缓存更新中…',
          closable: false,
          backdropClose: false,
        });
        setIsLoading(true);
        setError(null);
        try {
          await SyncIcons();
          await loadItems();
          notify({type: 'success', title: '图标已刷新'});
        } catch (err) {
          showError(err instanceof Error ? err.message : '刷新图标失败', '刷新失败');
        } finally {
          setIsLoading(false);
          closeModal(modalId);
        }
        return;
      }

      if (id === 'clear') {
        openModal({
          kind: 'confirm',
          title: '确认清空所有项目？',
          description: '此操作不可恢复，仍要继续吗？',
          tone: 'danger',
          primaryLabel: '清空',
          secondaryLabel: '取消',
          onConfirm: async () => {
            setIsLoading(true);
            setError(null);
            try {
              await ClearItems();
              await loadItems();
              notify({type: 'success', title: '清空完成', message: '列表已重置'});
            } catch (err) {
              showError(err instanceof Error ? err.message : '清空失败', '清空失败');
            } finally {
              setIsLoading(false);
            }
          },
        });
      }
    },
    [closeModal, loadItems, notify, openModal, showError]
  );

  const handleLaunch = useCallback(
    async (id: string) => {
      try {
        await LaunchItem(id);
        await loadItems();
      } catch (err) {
        showError(err instanceof Error ? err.message : '启动失败', '启动失败');
      }
    },
    [loadItems, showError]
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
          notify({
            type: 'success',
            title: current.favorite ? '已取消收藏' : '已收藏',
            message: current.name,
          });
        } catch (err) {
          showError(err instanceof Error ? err.message : '更新收藏失败');
        }
      }

      if (actionId === 'open-location') {
        try {
          await OpenItemLocation(current.id);
        } catch (err) {
          showError(err instanceof Error ? err.message : '无法打开所在目录');
        }
      }

      if (actionId === 'refresh-icon') {
        try {
          await RefreshItemIcon(current.id);
          await loadItems();
          notify({type: 'success', title: '图标已刷新', message: current.name});
        } catch (err) {
          showError(err instanceof Error ? err.message : '刷新图标失败');
        }
      }

      if (actionId === 'remove') {
        openModal({
          kind: 'confirm',
          title: `移除 "${current.name}"？`,
          description: '移除后会从列表中消失，可通过重新扫描恢复。',
          tone: 'danger',
          primaryLabel: '移除',
          secondaryLabel: '取消',
          onConfirm: async () => {
            try {
              await DeleteItem(current.id);
              await loadItems();
              notify({type: 'success', title: '已移除', message: current.name});
            } catch (err) {
              showError(err instanceof Error ? err.message : '移除失败');
            }
          },
        });
      }
    },
    [menuState.item, loadItems, notify, openModal, showError]
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
        <div className="grid-scroll">
          <AppGrid
            items={filteredItems}
            isLoading={isLoading}
            error={error}
            onAddItem={handleAddItem}
            onLaunch={handleLaunch}
            onOpenMenu={handleOpenMenu}
          />
        </div>
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
      <ToastHost />
      <ModalHost />
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
    case 'urls':
      return 'url';
    case 'system':
      return 'system';
    case 'apps':
    default:
      return 'app';
  }
}

function isWebPath(path: string) {
  const trimmed = path.trim().toLowerCase();
  return trimmed.startsWith('http://') || trimmed.startsWith('https://');
}
