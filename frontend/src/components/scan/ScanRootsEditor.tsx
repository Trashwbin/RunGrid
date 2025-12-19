import {useMemo, useState} from 'react';
import {PickScanRoot} from '../../../wailsjs/go/main/App';
import {Icon} from '../ui/Icon';
import {ScrollArea} from '../ui/ScrollArea';

type ScanRootsEditorProps = {
  initialRoots: string[];
  loading?: boolean;
  onChange: (next: string[]) => void;
};

export function ScanRootsEditor({
  initialRoots,
  loading = false,
  onChange,
}: ScanRootsEditorProps) {
  const [roots, setRoots] = useState<string[]>(() => initialRoots);

  const hasRoots = roots.length > 0;

  const normalizedRoots = useMemo(() => roots, [roots]);

  const handleRemove = (index: number) => {
    const next = normalizedRoots.filter((_, idx) => idx !== index);
    setRoots(next);
    onChange(next);
  };

  const handleAdd = async () => {
    try {
      const selected = await PickScanRoot();
      const trimmed = selected.trim();
      if (!trimmed) {
        return;
      }
      const exists = normalizedRoots.some(
        (root) => root.trim().toLowerCase() === trimmed.toLowerCase()
      );
      if (!exists) {
        const next = [...normalizedRoots, trimmed];
        setRoots(next);
        onChange(next);
      }
    } catch {
    }
  };

  if (loading) {
    return <p className="scan-roots-hint">正在读取扫描目录...</p>;
  }

  return (
    <div className="scan-roots">
      <p className="scan-roots-hint">
        默认会扫描桌面与开始菜单目录，你可以移除或追加路径。
      </p>
      <ScrollArea
        className="scan-root-list scroll-area--auto"
        viewportClassName="scan-root-list__viewport"
        contentClassName="scan-root-items"
      >
        {hasRoots ? (
          normalizedRoots.map((root, index) => (
            <div key={`${root}-${index}`} className="scan-root-item">
              <span className="scan-root-path" title={root}>
                {root}
              </span>
              <button
                type="button"
                className="scan-root-remove"
                onClick={() => handleRemove(index)}
                aria-label="移除扫描目录"
              >
                <Icon name="close" size={14} />
              </button>
            </div>
          ))
        ) : (
          <div className="scan-root-empty">暂无扫描目录</div>
        )}
      </ScrollArea>
      <div className="scan-root-add">
        <button type="button" onClick={handleAdd}>
          <Icon name="plus" size={16} />
          添加目录
        </button>
      </div>
    </div>
  );
}
