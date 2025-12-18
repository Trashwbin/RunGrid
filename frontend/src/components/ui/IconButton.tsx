import type {ReactNode} from 'react';

type IconButtonProps = {
  label: string;
  icon: ReactNode;
  onClick?: () => void;
  isActive?: boolean;
  className?: string;
};

export function IconButton({
  label,
  icon,
  onClick,
  isActive = false,
  className,
}: IconButtonProps) {
  return (
    <button
      type="button"
      className={`icon-button${isActive ? ' is-active' : ''}${
        className ? ` ${className}` : ''
      }`}
      aria-pressed={isActive}
      aria-label={label}
      title={label}
      onClick={onClick}
    >
      {icon}
    </button>
  );
}
