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
    case 'ban':
      return (
        <svg {...commonProps} aria-hidden="true">
          <circle cx="12" cy="12" r="8" />
          <path d="M7 7l10 10" />
        </svg>
      );
    case 'star':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M12 4l2.6 5.2 5.7.8-4.1 4 1 5.6-5.2-2.7-5.2 2.7 1-5.6-4.1-4 5.7-.8z" />
        </svg>
      );
    case 'home':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M4 11l8-6 8 6v8a1 1 0 0 1-1 1h-5v-5h-6v5H5a1 1 0 0 1-1-1z" />
        </svg>
      );
    case 'download':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M12 4v10M8 10l4 4 4-4M4 20h16" />
        </svg>
      );
    case 'upload':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M12 20V10M8 14l4-4 4 4M4 4h16" />
        </svg>
      );
    case 'cloud':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M7 18h9a4 4 0 0 0 .5-8 5 5 0 0 0-9.5 1A3.5 3.5 0 0 0 7 18z" />
        </svg>
      );
    case 'shield':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M12 3l7 3v5c0 5-3 9-7 11-4-2-7-6-7-11V6z" />
        </svg>
      );
    case 'code':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M8 8l-4 4 4 4M16 8l4 4-4 4" />
        </svg>
      );
    case 'terminal':
      return (
        <svg {...commonProps} aria-hidden="true">
          <rect x="4" y="5" width="16" height="14" rx="2" />
          <path d="M7 10l2 2-2 2M11 14h4" />
        </svg>
      );
    case 'palette':
      return (
        <svg {...commonProps} aria-hidden="true">
          <circle cx="12" cy="12" r="8" />
          <circle cx="9" cy="10" r="1" />
          <circle cx="12" cy="8" r="1" />
          <circle cx="15" cy="10" r="1" />
        </svg>
      );
    case 'calendar':
      return (
        <svg {...commonProps} aria-hidden="true">
          <rect x="4" y="6" width="16" height="14" rx="2" />
          <path d="M8 3v4M16 3v4M4 10h16" />
        </svg>
      );
    case 'camera':
      return (
        <svg {...commonProps} aria-hidden="true">
          <rect x="4" y="7" width="16" height="12" rx="2" />
          <path d="M9 7l1-2h4l1 2" />
          <circle cx="12" cy="13" r="3" />
        </svg>
      );
    case 'mail':
      return (
        <svg {...commonProps} aria-hidden="true">
          <rect x="4" y="6" width="16" height="12" rx="2" />
          <path d="M4 7l8 6 8-6" />
        </svg>
      );
    case 'phone':
      return (
        <svg {...commonProps} aria-hidden="true">
          <rect x="7" y="4" width="10" height="16" rx="2" />
          <circle cx="12" cy="17" r="1" />
        </svg>
      );
    case 'cart':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M5 6h2l2 9h9l2-6H8" />
          <circle cx="10" cy="19" r="1" />
          <circle cx="17" cy="19" r="1" />
        </svg>
      );
    case 'map':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M4 6l6-2 6 2 4-2v14l-4 2-6-2-6 2z" />
          <path d="M10 4v14M16 6v14" />
        </svg>
      );
    case 'chat':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M4 6h16v9H8l-4 3z" />
        </svg>
      );
    case 'bookmark':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M6 4h12v16l-6-3-6 3z" />
        </svg>
      );
    case 'bolt':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M13 2L6 12h5l-1 10 7-10h-5z" />
        </svg>
      );
    case 'key':
      return (
        <svg {...commonProps} aria-hidden="true">
          <circle cx="8" cy="14" r="3" />
          <path d="M11 14h9M16 14v-2M19 14v-2" />
        </svg>
      );
    case 'lock':
      return (
        <svg {...commonProps} aria-hidden="true">
          <rect x="5" y="11" width="14" height="9" rx="2" />
          <path d="M8 11V8a4 4 0 0 1 8 0v3" />
        </svg>
      );
    case 'chip':
      return (
        <svg {...commonProps} aria-hidden="true">
          <rect x="7" y="7" width="10" height="10" rx="2" />
          <path d="M9 1v4M15 1v4M9 19v4M15 19v4M1 9h4M1 15h4M19 9h4M19 15h4" />
        </svg>
      );
    case 'rocket':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M12 3c3 2 5 5 5 9l-5 5-5-5c0-4 2-7 5-9z" />
          <circle cx="12" cy="11" r="1.5" />
          <path d="M9 17l-3 3M15 17l3 3" />
        </svg>
      );
    case 'image':
      return (
        <svg {...commonProps} aria-hidden="true">
          <rect x="4" y="5" width="16" height="14" rx="2" />
          <circle cx="9" cy="10" r="1.5" />
          <path d="M4 15l4-4 3 3 3-3 6 6" />
        </svg>
      );
    case 'puzzle':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M9 3h3a2 2 0 0 1 2 2v2h2a2 2 0 1 1 0 4h-2v2a2 2 0 0 1-2 2h-2v2a2 2 0 1 1-4 0v-2H7a2 2 0 0 1-2-2v-2H3a2 2 0 1 1 0-4h2V5a2 2 0 0 1 2-2h2z" />
        </svg>
      );
    case 'globe':
      return (
        <svg {...commonProps} aria-hidden="true">
          <circle cx="12" cy="12" r="8" />
          <path d="M4 12h16" />
          <path d="M12 4c2.5 2.5 2.5 11.5 0 16" />
          <path d="M12 4c-2.5 2.5-2.5 11.5 0 16" />
        </svg>
      );
    case 'database':
      return (
        <svg {...commonProps} aria-hidden="true">
          <ellipse cx="12" cy="6" rx="7" ry="3" />
          <path d="M5 6v8c0 1.7 3.1 3 7 3s7-1.3 7-3V6" />
          <path d="M5 10c0 1.7 3.1 3 7 3s7-1.3 7-3" />
        </svg>
      );
    case 'gamepad':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M6 10h12a4 4 0 0 1 3.7 5.4l-1.2 3a3 3 0 0 1-2.8 1.6H6.3a3 3 0 0 1-2.8-1.6l-1.2-3A4 4 0 0 1 6 10z" />
          <path d="M9 14H7" />
          <path d="M8 13v2" />
          <circle cx="16.5" cy="14" r="1" />
          <circle cx="18.5" cy="12.5" r="1" />
        </svg>
      );
    case 'video':
      return (
        <svg {...commonProps} aria-hidden="true">
          <rect x="4" y="6" width="16" height="12" rx="2" />
          <path d="M10 9l5 3-5 3z" />
        </svg>
      );
    case 'briefcase':
      return (
        <svg {...commonProps} aria-hidden="true">
          <rect x="4" y="8" width="16" height="10" rx="2" />
          <path d="M9 8V6h6v2" />
          <path d="M4 12h16" />
        </svg>
      );
    case 'tools':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M7 7l10 10" />
          <path d="M17 7l-4 4" />
          <path d="M7 17l4-4" />
        </svg>
      );
    case 'file':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M6 4h8l4 4v12H6z" />
          <path d="M14 4v4h4" />
        </svg>
      );
    case 'music':
      return (
        <svg {...commonProps} aria-hidden="true">
          <path d="M11 14V6l7-2v8" />
          <circle cx="9" cy="18" r="2" />
          <circle cx="16" cy="16" r="2" />
        </svg>
      );
    case 'user':
      return (
        <svg {...commonProps} aria-hidden="true">
          <circle cx="12" cy="9" r="3" />
          <path d="M5 20a7 7 0 0 1 14 0" />
        </svg>
      );
    default:
      return null;
  }
}
