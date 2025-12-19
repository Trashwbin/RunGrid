import {useEffect, useRef} from 'react';
import {useOutsideClick} from '../../hooks/useOutsideClick';

type ContextMenuItem = {
  id: string;
  label: string;
  disabled?: boolean;
  tone?: 'danger';
};

type ContextMenuProps = {
  open: boolean;
  x: number;
  y: number;
  items: ContextMenuItem[];
  onSelect: (id: string) => void;
  onClose: () => void;
};

export function ContextMenu({
  open,
  x,
  y,
  items,
  onSelect,
  onClose,
}: ContextMenuProps) {
  const menuRef = useRef<HTMLDivElement>(null);

  useOutsideClick(menuRef, onClose, open);

  useEffect(() => {
    if (!open) {
      return;
    }

    const handleKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClose();
      }
    };

    window.addEventListener('keydown', handleKey);
    return () => {
      window.removeEventListener('keydown', handleKey);
    };
  }, [open, onClose]);

  if (!open) {
    return null;
  }

  return (
    <div
      ref={menuRef}
      className="context-menu"
      style={{left: `${x}px`, top: `${y}px`}}
      role="menu"
    >
      {items.map((item) => (
        <button
          key={item.id}
          type="button"
          className={`context-menu-item${item.tone === 'danger' ? ' context-menu-item--danger' : ''}`}
          role="menuitem"
          disabled={item.disabled}
          onClick={() => {
            if (item.disabled) {
              return;
            }
            onSelect(item.id);
            onClose();
          }}
        >
          {item.label}
        </button>
      ))}
    </div>
  );
}
