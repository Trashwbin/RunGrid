import type {AppItem} from '../../types';

type AppTileProps = {
  item: AppItem;
  onLaunch?: (id: string) => void;
};

export function AppTile({item, onLaunch}: AppTileProps) {
  return (
    <button
      type="button"
      className="app-tile"
      title={item.name}
      onClick={() => onLaunch?.(item.id)}
    >
      <div className={`app-icon app-icon--${item.accent}`} aria-hidden="true">
        {item.iconUrl ? (
          <img src={item.iconUrl} alt="" />
        ) : (
          <span className="app-glyph">{item.glyph}</span>
        )}
      </div>
      <span className="app-name">{item.name}</span>
    </button>
  );
}
