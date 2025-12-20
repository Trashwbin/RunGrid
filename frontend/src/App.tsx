import {useCallback, useEffect, useMemo, useRef, useState} from 'react';
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
  ListScanRoots,
  OpenItemLocation,
  RefreshItemIcon,
  ScanShortcuts,
  SetFavorite,
  SyncIcons,
  UpdateItem,
  UpdateItemIconFromSource,
} from '../wailsjs/go/main/App';
import type {domain} from '../wailsjs/go/models';
import {EventsOn} from '../wailsjs/runtime/runtime';
import {AppGrid} from './components/grid/AppGrid';
import {CategoryBar} from './components/layout/CategoryBar';
import {GroupTabs} from './components/layout/GroupTabs';
import {SearchBar} from './components/layout/SearchBar';
import {TopBar} from './components/layout/TopBar';
import {ScanRootsEditor} from './components/scan/ScanRootsEditor';
import {EditItemForm, type EditDraft} from './components/item/EditItemForm';
import {ModalHost} from './components/overlay/ModalHost';
import {ToastHost} from './components/overlay/ToastHost';
import {ContextMenu} from './components/ui/ContextMenu';
import {ScrollArea} from './components/ui/ScrollArea';
import {SettingsModal} from './components/settings/SettingsModal';
import {useModalStore, useToastStore} from './store/overlays';
import {toAppItem, toGroupTab} from './utils/items';
import {loadHotkeys, saveHotkeys, type HotkeyConfig} from './utils/hotkeys';
import type {AppItem} from './types';

function App() {
  const [items, setItems] = useState<domain.Item[]>([]);
  const [groups, setGroups] = useState<domain.Group[]>([]);
  const [activeCategoryId, setActiveCategoryId] = useState(categories[0].id);
  const [activeGroupId, setActiveGroupId] = useState('all');
  const [query, setQuery] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [scanRoots, setScanRoots] = useState<string[]>([]);
  const [scanRootsReady, setScanRootsReady] = useState(false);
  const [iconVersion, setIconVersion] = useState(0);
  const scanRootsRef = useRef<string[]>([]);
  const editDraftRef = useRef<EditDraft | null>(null);
  const hotkeyDraftRef = useRef<HotkeyConfig | null>(null);
  const notify = useToastStore((state) => state.notify);
  const openModal = useModalStore((state) => state.openModal);
  const closeModal = useModalStore((state) => state.closeModal);
  const updateModal = useModalStore((state) => state.updateModal);
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

  const handleScanRootsChange = useCallback((next: string[]) => {
    setScanRoots(next);
    scanRootsRef.current = next;
  }, []);

  const bumpIconVersion = useCallback(() => {
    setIconVersion((prev) => prev + 1);
  }, []);

  useEffect(() => {
    const cached = window.localStorage.getItem('rungrid.scanRoots');
    if (cached) {
      try {
        const parsed = JSON.parse(cached);
        if (Array.isArray(parsed)) {
        handleScanRootsChange(normalizeRoots(parsed));
        setScanRootsReady(true);
        return;
      }
    } catch {
        window.localStorage.removeItem('rungrid.scanRoots');
      }
    }

    ListScanRoots()
      .then((roots) => {
        handleScanRootsChange(normalizeRoots(roots));
      })
      .catch((err) => {
        showError(
          err instanceof Error ? err.message : '无法读取扫描路径',
          '读取失败'
        );
      })
      .finally(() => setScanRootsReady(true));
  }, [handleScanRootsChange, showError]);

  useEffect(() => {
    if (!scanRootsReady) {
      return;
    }
    window.localStorage.setItem('rungrid.scanRoots', JSON.stringify(scanRoots));
  }, [scanRoots, scanRootsReady]);

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

  useEffect(() => {
    const off = EventsOn('icons:updated', () => {
      loadItems();
      bumpIconVersion();
    });
    return () => {
      off();
    };
  }, [loadItems, bumpIconVersion]);

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

  const startScan = useCallback(
    async (roots: string[]) => {
      const normalizedRoots = normalizeRoots(roots);
      if (normalizedRoots.length === 0) {
        notify({
          type: 'warning',
          title: '请先添加扫描目录',
          message: '至少保留一个有效路径。',
        });
        return;
      }

      const modalId = openModal({
        kind: 'progress',
        title: '正在扫描',
        description: `正在扫描 ${normalizedRoots.length} 个目录…`,
        closable: false,
        backdropClose: false,
      });
      const off = EventsOn('scan:progress', (payload: ScanProgressPayload) => {
        if (!payload) {
          return;
        }
        const patch: {
          path?: string;
          progress?: number;
          description?: string;
        } = {};
        if (payload.path) {
          patch.path = payload.path;
        }
        if (typeof payload.percent === 'number' && payload.percent >= 0) {
          patch.progress = payload.percent;
        }
        if (payload.rootTotal && payload.rootIndex) {
          patch.description = `正在扫描第 ${payload.rootIndex}/${payload.rootTotal} 个目录`;
        }
        if (Object.keys(patch).length > 0) {
          updateModal(modalId, patch);
        }
      });

      setIsLoading(true);
      setError(null);
      try {
        await ScanShortcuts(normalizedRoots);
        await loadItems();
        notify({type: 'success', title: '扫描完成', message: '应用列表已更新'});
      } catch (err) {
        showError(err instanceof Error ? err.message : '扫描失败', '扫描失败');
      } finally {
        setIsLoading(false);
        off();
        closeModal(modalId);
      }
    },
    [closeModal, loadItems, notify, openModal, showError, updateModal]
  );

  const handleMenuSelect = useCallback(
    async (id: string) => {
      if (id === 'settings') {
        const initial = loadHotkeys();
        hotkeyDraftRef.current = initial;
        const modalId = openModal({
          kind: 'form',
          title: '设置',
          description: '管理快捷键与偏好设置。',
          size: 'lg',
          primaryLabel: '保存',
          secondaryLabel: '关闭',
          autoClose: false,
          content: (
            <SettingsModal
              initialHotkeys={initial}
              onChange={(next) => {
                hotkeyDraftRef.current = next;
              }}
            />
          ),
          onConfirm: () => {
            if (hotkeyDraftRef.current) {
              saveHotkeys(hotkeyDraftRef.current);
            }
            notify({type: 'success', title: '设置已保存'});
            hotkeyDraftRef.current = null;
            closeModal(modalId);
          },
          onCancel: () => {
            hotkeyDraftRef.current = null;
            closeModal(modalId);
          },
        });
        return;
      }

      if (id === 'scan') {
        if (!scanRootsReady) {
          notify({
            type: 'info',
            title: '正在读取扫描目录',
            message: '请稍候再试。',
          });
          return;
        }
        openModal({
          kind: 'form',
          title: '扫描路径',
          description: '选择需要扫描的目录，可随时调整。',
          size: 'lg',
          primaryLabel: '开始扫描',
          secondaryLabel: '取消',
          content: (
            <ScanRootsEditor
              initialRoots={scanRoots}
              loading={!scanRootsReady}
              onChange={handleScanRootsChange}
            />
          ),
          onConfirm: () => {
            void startScan(scanRootsRef.current);
          },
        });
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
          bumpIconVersion();
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
    [
      closeModal,
      bumpIconVersion,
      loadItems,
      notify,
      openModal,
      scanRoots,
      scanRootsReady,
      showError,
      startScan,
    ]
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
    () => items.map((item, index) => toAppItem(item, index, iconVersion)),
    [items, iconVersion]
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

      if (actionId === 'edit') {
        const initialDraft: EditDraft = {
          id: current.id,
          name: current.name,
          path: current.path,
          originalPath: current.path,
          iconUrl: current.iconUrl,
          glyph: current.glyph,
          favorite: current.favorite,
          hidden: current.hidden,
        };
        editDraftRef.current = initialDraft;
        const modalId = openModal({
          kind: 'form',
          title: '编辑快捷方式',
          description: '更新名称、目标路径或替换图标。',
          size: 'lg',
          primaryLabel: '保存',
          secondaryLabel: '取消',
          autoClose: false,
          content: (
            <EditItemForm
              initialDraft={initialDraft}
              onChange={(next) => {
                editDraftRef.current = next;
              }}
            />
          ),
          onConfirm: async () => {
            const draft = editDraftRef.current;
            if (!draft) {
              return;
            }
            const name = draft.name.trim();
            const path = draft.path.trim();
            if (!name || !path) {
              notify({
                type: 'warning',
                title: '请补全信息',
                message: '名称和目标路径不能为空。',
              });
              return;
            }

            try {
              const pathChanged =
                draft.originalPath.trim().toLowerCase() !== path.toLowerCase();
              const updated = await UpdateItem({
                id: draft.id,
                name,
                path,
                type: pathChanged ? '' : current.type,
                icon_path: '',
                group_id: '',
                tags: current.tags,
                favorite: draft.favorite,
                hidden: draft.hidden,
              });
              let iconUpdated = false;

              if (draft.iconSource) {
                await UpdateItemIconFromSource(draft.id, draft.iconSource);
                iconUpdated = true;
              } else if (draft.originalPath !== path && updated.type !== 'url') {
                try {
                  await RefreshItemIcon(draft.id);
                  iconUpdated = true;
                } catch {
                }
              }

              await loadItems();
              if (iconUpdated) {
                bumpIconVersion();
              }
              notify({type: 'success', title: '已更新', message: name});
              closeModal(modalId);
              editDraftRef.current = null;
            } catch (err) {
              showError(err instanceof Error ? err.message : '更新失败');
            }
          },
          onCancel: () => {
            editDraftRef.current = null;
            closeModal(modalId);
          },
        });
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
          bumpIconVersion();
          notify({type: 'success', title: '图标已刷新', message: current.name});
        } catch (err) {
          showError(err instanceof Error ? err.message : '刷新图标失败');
        }
      }

      if (actionId === 'remove') {
        try {
          await DeleteItem(current.id);
          await loadItems();
          notify({type: 'success', title: '已移除', message: current.name});
        } catch (err) {
          showError(err instanceof Error ? err.message : '移除失败');
        }
      }
    },
    [
      menuState.item,
      loadItems,
      notify,
      openModal,
      showError,
      closeModal,
      bumpIconVersion,
    ]
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
        <ScrollArea className="grid-scroll" viewportClassName="grid-scroll__viewport">
          <AppGrid
            items={filteredItems}
            isLoading={isLoading}
            error={error}
            onAddItem={handleAddItem}
            onLaunch={handleLaunch}
            onOpenMenu={handleOpenMenu}
          />
        </ScrollArea>
      </div>
      <SearchBar value={query} onChange={setQuery} />
      <ContextMenu
        open={menuState.open}
        x={menuState.x}
        y={menuState.y}
        items={[
          {
            id: 'edit',
            label: '编辑',
          },
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

type ScanProgressPayload = {
  root?: string;
  path?: string;
  rootIndex?: number;
  rootTotal?: number;
  scanned?: number;
  percent?: number;
};

function normalizeRoots(roots: string[]) {
  const next: string[] = [];
  const seen = new Set<string>();
  for (const root of roots) {
    const trimmed = String(root ?? '').trim();
    if (!trimmed) {
      continue;
    }
    const key = trimmed.toLowerCase();
    if (seen.has(key)) {
      continue;
    }
    seen.add(key);
    next.push(trimmed);
  }
  return next;
}
