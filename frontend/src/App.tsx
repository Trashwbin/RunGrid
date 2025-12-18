import {useCallback, useEffect, useMemo, useState} from 'react';
import './App.css';
import {categories, menuItems} from './data/mock';
import {CreateGroup, CreateItem, ListGroups, ListItems} from '../wailsjs/go/main/App';
import type {domain} from '../wailsjs/go/models';
import {AppGrid} from './components/grid/AppGrid';
import {CategoryBar} from './components/layout/CategoryBar';
import {GroupTabs} from './components/layout/GroupTabs';
import {SearchBar} from './components/layout/SearchBar';
import {TopBar} from './components/layout/TopBar';
import {toAppItem, toGroupTab} from './utils/items';

function App() {
  const [items, setItems] = useState<domain.Item[]>([]);
  const [groups, setGroups] = useState<domain.Group[]>([]);
  const [activeCategoryId, setActiveCategoryId] = useState(categories[0].id);
  const [activeGroupId, setActiveGroupId] = useState('all');
  const [query, setQuery] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

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

  return (
    <div className="app-shell">
      <TopBar title="RunGrid" menuItems={menuItems} />
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
        />
      </div>
      <SearchBar value={query} onChange={setQuery} />
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
