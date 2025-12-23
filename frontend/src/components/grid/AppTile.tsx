import type {AppItem} from '../../types';

type AppTileProps = {
  item: AppItem;
  selected?: boolean;
  selectionMode?: boolean;
  onLaunch?: (id: string) => void;
  onOpenMenu?: (item: AppItem, x: number, y: number) => void;
  onSelect?: (id: string, multi: boolean) => void;
};

export function AppTile({
  item,
  selected = false,
  selectionMode = false,
  onLaunch,
  onOpenMenu,
  onSelect,
}: AppTileProps) {
  return (
    <button
      type="button"
      className={`app-tile${selected ? ' is-selected' : ''}`}
      title={item.name}
      onClick={(event) => {
        const multi = event.ctrlKey || event.metaKey;
        if (selectionMode || multi) {
          event.preventDefault();
          onSelect?.(item.id, true);
          return;
        }
        onLaunch?.(item.id);
      }}
      onContextMenu={(event) => {
        event.preventDefault();
        event.stopPropagation();
        onOpenMenu?.(item, event.clientX, event.clientY);
      }}
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
