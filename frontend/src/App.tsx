import {useEffect, useMemo, useState} from 'react';
import './App.css';
import {categories, menuItems} from './data/mock';
import {ListGroups, ListItems} from '../wailsjs/go/main/App';
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

  useEffect(() => {
    let active = true;

    const loadGroups = async () => {
      try {
        const data = await ListGroups();
        if (active) {
          setGroups(data);
        }
      } catch (err) {
        if (active) {
          setError(err instanceof Error ? err.message : '无法加载分组');
        }
      }
    };

    loadGroups();

    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    let active = true;
    setIsLoading(true);
    setError(null);

    const loadItems = async () => {
      try {
        const data = await ListItems(activeGroupId, query);
        if (active) {
          setItems(data);
        }
      } catch (err) {
        if (active) {
          setError(err instanceof Error ? err.message : '无法加载项目');
        }
      } finally {
        if (active) {
          setIsLoading(false);
        }
      }
    };

    loadItems();

    return () => {
      active = false;
    };
  }, [activeGroupId, query]);

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
        />
        <AppGrid items={filteredItems} isLoading={isLoading} error={error} />
      </div>
      <SearchBar value={query} onChange={setQuery} />
    </div>
  );
}

export default App;
