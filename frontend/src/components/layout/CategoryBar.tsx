import type {Category} from '../../types';
import {Icon} from '../ui/Icon';

type CategoryBarProps = {
  categories: Category[];
  activeId: string;
  onSelect: (id: string) => void;
};

export function CategoryBar({categories, activeId, onSelect}: CategoryBarProps) {
  return (
    <nav className="category-bar" aria-label="主分类">
      {categories.map((category) => (
        <button
          key={category.id}
          type="button"
          className={`category-item${
            activeId === category.id ? ' is-active' : ''
          }`}
          onClick={() => onSelect(category.id)}
        >
          <Icon name={category.icon} size={16} />
          <span>{category.label}</span>
        </button>
      ))}
    </nav>
  );
}
