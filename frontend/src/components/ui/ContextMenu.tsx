import {useCallback, useEffect, useLayoutEffect, useRef, useState} from 'react';
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
  const [position, setPosition] = useState({x, y});

  useOutsideClick(menuRef, onClose, open);

  const updatePosition = useCallback(() => {
    const menu = menuRef.current;
    if (!menu) {
      return;
    }
    const rect = menu.getBoundingClientRect();
    const padding = 8;
    let nextX = x;
    let nextY = y;

    if (x + rect.width + padding > window.innerWidth) {
      nextX = x - rect.width;
    }
    if (y + rect.height + padding > window.innerHeight) {
      nextY = y - rect.height;
    }

    const maxX = window.innerWidth - rect.width - padding;
    const maxY = window.innerHeight - rect.height - padding;
    nextX = Math.max(padding, Math.min(nextX, maxX));
    nextY = Math.max(padding, Math.min(nextY, maxY));
    setPosition({x: nextX, y: nextY});
  }, [x, y]);

  useLayoutEffect(() => {
    if (!open) {
      return;
    }
    setPosition({x, y});
    const frame = window.requestAnimationFrame(updatePosition);
    return () => window.cancelAnimationFrame(frame);
  }, [open, x, y, updatePosition, items.length]);

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

  useEffect(() => {
    if (!open) {
      return;
    }
    const handleResize = () => updatePosition();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [open, updatePosition]);

  if (!open) {
    return null;
  }

  return (
    <div
      ref={menuRef}
      className="context-menu"
      style={{left: `${position.x}px`, top: `${position.y}px`}}
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
