import {useMemo, useState} from 'react';
import './App.css';
import {appItems, categories, groupTabs, menuItems} from './data/mock';
import {AppGrid} from './components/grid/AppGrid';
import {CategoryBar} from './components/layout/CategoryBar';
import {GroupTabs} from './components/layout/GroupTabs';
import {SearchBar} from './components/layout/SearchBar';
import {TopBar} from './components/layout/TopBar';

function App() {
  const [activeCategoryId, setActiveCategoryId] = useState(categories[0].id);
  const [activeGroupId, setActiveGroupId] = useState(groupTabs[0].id);
  const [query, setQuery] = useState('');

  const filteredItems = useMemo(() => {
    return appItems.filter((item) => {
      const matchesCategory = item.categoryId === activeCategoryId;
      const matchesGroup =
        activeGroupId === 'all' || item.groupId === activeGroupId;
      const matchesQuery = query
        ? item.name.toLowerCase().includes(query.trim().toLowerCase())
        : true;

      return matchesCategory && matchesGroup && matchesQuery;
    });
  }, [activeCategoryId, activeGroupId, query]);

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
        <AppGrid items={filteredItems} />
      </div>
      <SearchBar value={query} onChange={setQuery} />
    </div>
  );
}

export default App;
