import type {IconName} from '../../types';

type IconProps = {
  name: IconName;
  size?: number;
  className?: string;
};

export function Icon({name, size = 18, className}: IconProps) {
  const commonProps = {
    width: size,
    height: size,
    viewBox: '0 0 24 24',
    fill: 'none',
    stroke: 'currentColor',
    strokeWidth: 1.6,
    strokeLinecap: 'round' as const,
    strokeLinejoin: 'round' as const,
    className,
  };

  switch (name) {
    case 'menu':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M4 6h16M4 12h16M4 18h12" />
        </svg>
      );
    case 'settings':
      return (
        <svg {...commonProps} aria-hidden="true">
          <circle cx="12" cy="12" r="7" />
          <path d="M12 9v6M9 12h6" />
        </svg>
      );
    case 'close':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M6 6l12 12M18 6l-12 12" />
        </svg>
      );
    case 'apps':
      return (
        <svg {...commonProps} aria-hidden="true">
          <rect x="4" y="4" width="6" height="6" rx="1" />
          <rect x="14" y="4" width="6" height="6" rx="1" />
          <rect x="4" y="14" width="6" height="6" rx="1" />
          <rect x="14" y="14" width="6" height="6" rx="1" />
        </svg>
      );
    case 'system':
      return (
        <svg {...commonProps} aria-hidden="true">
          <rect x="4" y="5" width="16" height="14" rx="2" />
          <path d="M8 9h8M8 13h5" />
        </svg>
      );
    case 'docs':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M6 5h9l3 3v11H6z" />
          <path d="M15 5v3h3" />
        </svg>
      );
    case 'folder':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M3.5 7h6l2 2h9v8.5a2 2 0 0 1-2 2h-13a2 2 0 0 1-2-2z" />
        </svg>
      );
    case 'link':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M10 13a4 4 0 0 1 0-6l1.5-1.5a4 4 0 0 1 6 6L16 12" />
          <path d="M14 11a4 4 0 0 1 0 6l-1.5 1.5a4 4 0 0 1-6-6L8 12" />
        </svg>
      );
    case 'plus':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M12 5v14M5 12h14" />
        </svg>
      );
    case 'heart':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M12 19s-7-4.5-7-9a4 4 0 0 1 7-2 4 4 0 0 1 7 2c0 4.5-7 9-7 9z" />
        </svg>
      );
    default:
      return null;
  }
}
