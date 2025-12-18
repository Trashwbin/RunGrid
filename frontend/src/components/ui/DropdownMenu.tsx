import type {MenuItem} from '../../types';

type DropdownMenuProps = {
  open: boolean;
  items: MenuItem[];
  onSelect?: (id: string) => void;
};

export function DropdownMenu({open, items, onSelect}: DropdownMenuProps) {
  if (!open) {
    return null;
  }

  return (
    <div className="dropdown-menu" role="menu">
      {items.map((item) => (
        <button
          key={item.id}
          type="button"
          className="dropdown-item"
          role="menuitem"
          onClick={() => onSelect?.(item.id)}
        >
          {item.label}
        </button>
      ))}
    </div>
  );
}
