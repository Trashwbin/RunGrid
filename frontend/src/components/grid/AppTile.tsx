import type {AppItem} from '../../types';

type AppTileProps = {
  item: AppItem;
};

export function AppTile({item}: AppTileProps) {
  return (
    <button type="button" className="app-tile" title={item.name}>
      <div className={`app-icon app-icon--${item.accent}`} aria-hidden="true">
        {item.iconPath ? (
          <img src={item.iconPath} alt="" />
        ) : (
          <span className="app-glyph">{item.glyph}</span>
        )}
      </div>
      <span className="app-name">{item.name}</span>
    </button>
  );
}
