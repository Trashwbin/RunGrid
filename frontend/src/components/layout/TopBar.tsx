import {useCallback, useRef, useState} from 'react';
import type {MenuItem} from '../../types';
import {useOutsideClick} from '../../hooks/useOutsideClick';
import {DropdownMenu} from '../ui/DropdownMenu';
import {Icon} from '../ui/Icon';
import {IconButton} from '../ui/IconButton';
import {WindowHide} from '../../../wailsjs/runtime/runtime';

type TopBarProps = {
  title: string;
  menuItems: MenuItem[];
  onMenuSelect?: (id: string) => void;
  onHide?: () => void;
};

export function TopBar({title, menuItems, onMenuSelect, onHide}: TopBarProps) {
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useOutsideClick(menuRef, () => setMenuOpen(false), menuOpen);

  const handleMenuToggle = useCallback(() => {
    setMenuOpen((prev) => !prev);
  }, []);

  const handleMenuSelect = useCallback(
    (id: string) => {
      onMenuSelect?.(id);
      setMenuOpen(false);
    },
    [onMenuSelect]
  );

  const handleClose = useCallback(() => {
    if (onHide) {
      onHide();
      return;
    }
    WindowHide();
  }, [onHide]);

  return (
    <header className="top-bar">
      <div className="brand">
        <span className="brand-mark" aria-hidden="true" />
        <span className="brand-name">{title}</span>
      </div>
      <div className="top-actions" ref={menuRef}>
        <IconButton
          label="菜单"
          icon={<Icon name="menu" />}
          isActive={menuOpen}
          onClick={handleMenuToggle}
        />
        <IconButton
          label="隐藏到托盘"
          icon={<Icon name="close" />}
          onClick={handleClose}
        />
        <DropdownMenu
          open={menuOpen}
          items={menuItems}
          onSelect={handleMenuSelect}
        />
      </div>
    </header>
  );
}
