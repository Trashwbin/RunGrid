import {useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState} from 'react';
import type {ReactNode} from 'react';
import './ScrollArea.css';

type ScrollAreaProps = {
  className?: string;
  viewportClassName?: string;
  contentClassName?: string;
  children: ReactNode;
  minThumbHeight?: number;
};

type Metrics = {
  canScroll: boolean;
  thumbHeight: number;
  thumbTop: number;
};

const defaultMetrics: Metrics = {canScroll: false, thumbHeight: 0, thumbTop: 0};

export function ScrollArea({
  className,
  viewportClassName,
  contentClassName,
  children,
  minThumbHeight = 26,
}: ScrollAreaProps) {
  const viewportRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const [metrics, setMetrics] = useState<Metrics>(defaultMetrics);
  const [dragging, setDragging] = useState(false);
  const dragStartRef = useRef(0);
  const scrollStartRef = useRef(0);
  const metricsRef = useRef(metrics);

  const updateMetrics = useCallback(() => {
    const viewport = viewportRef.current;
    if (!viewport) {
      return;
    }
    const {scrollHeight, clientHeight, scrollTop} = viewport;
    if (scrollHeight <= clientHeight + 1) {
      setMetrics(defaultMetrics);
      return;
    }
    const trackHeight = clientHeight;
    const thumbHeight = Math.max(
      minThumbHeight,
      Math.round((clientHeight / scrollHeight) * trackHeight)
    );
    const scrollable = scrollHeight - clientHeight;
    const maxThumbTop = trackHeight - thumbHeight;
    const thumbTop =
      scrollable > 0 ? Math.round((scrollTop / scrollable) * maxThumbTop) : 0;
    setMetrics({canScroll: true, thumbHeight, thumbTop});
  }, [minThumbHeight]);

  useLayoutEffect(() => {
    updateMetrics();
  }, [updateMetrics, children]);

  useEffect(() => {
    metricsRef.current = metrics;
  }, [metrics]);

  useEffect(() => {
    const viewport = viewportRef.current;
    if (!viewport) {
      return;
    }
    const handleScroll = () => updateMetrics();
    viewport.addEventListener('scroll', handleScroll, {passive: true});
    const observer = new ResizeObserver(() => updateMetrics());
    observer.observe(viewport);
    if (contentRef.current) {
      observer.observe(contentRef.current);
    }
    return () => {
      viewport.removeEventListener('scroll', handleScroll);
      observer.disconnect();
    };
  }, [updateMetrics]);

  useEffect(() => {
    if (!dragging) {
      return;
    }

    const handleMove = (event: MouseEvent) => {
      const viewport = viewportRef.current;
      if (!viewport) {
        return;
      }
      const {scrollHeight, clientHeight} = viewport;
      const scrollable = scrollHeight - clientHeight;
      const trackHeight = clientHeight;
      const thumbHeight = metricsRef.current.thumbHeight;
      const maxThumbTop = trackHeight - thumbHeight;
      if (maxThumbTop <= 0) {
        return;
      }
      const delta = event.clientY - dragStartRef.current;
      const scrollDelta = (delta / maxThumbTop) * scrollable;
      viewport.scrollTop = scrollStartRef.current + scrollDelta;
    };

    const handleUp = () => {
      setDragging(false);
    };

    window.addEventListener('mousemove', handleMove);
    window.addEventListener('mouseup', handleUp);
    return () => {
      window.removeEventListener('mousemove', handleMove);
      window.removeEventListener('mouseup', handleUp);
    };
  }, [dragging]);

  const handleThumbMouseDown = useCallback(
    (event: React.MouseEvent<HTMLDivElement>) => {
      event.preventDefault();
      const viewport = viewportRef.current;
      if (!viewport) {
        return;
      }
      dragStartRef.current = event.clientY;
      scrollStartRef.current = viewport.scrollTop;
      setDragging(true);
    },
    []
  );

  const handleTrackMouseDown = useCallback(
    (event: React.MouseEvent<HTMLDivElement>) => {
      const viewport = viewportRef.current;
      if (!viewport) {
        return;
      }
      const target = event.target as HTMLElement;
      if (target.closest('.scroll-area__thumb')) {
        return;
      }
      const rect = event.currentTarget.getBoundingClientRect();
      const clickY = event.clientY - rect.top;
      const {scrollHeight, clientHeight} = viewport;
      const scrollable = scrollHeight - clientHeight;
      const thumbHeight = metricsRef.current.thumbHeight;
      const maxThumbTop = rect.height - thumbHeight;
      if (maxThumbTop <= 0) {
        return;
      }
      const nextThumbTop = Math.max(0, Math.min(clickY - thumbHeight / 2, maxThumbTop));
      viewport.scrollTop = (nextThumbTop / maxThumbTop) * scrollable;
    },
    []
  );

  const wrapperClass = useMemo(() => {
    const base = 'scroll-area';
    return className ? `${base} ${className}` : base;
  }, [className]);

  const viewportClass = useMemo(() => {
    const base = 'scroll-area__viewport';
    return viewportClassName ? `${base} ${viewportClassName}` : base;
  }, [viewportClassName]);

  const contentClass = useMemo(() => {
    const base = 'scroll-area__content';
    return contentClassName ? `${base} ${contentClassName}` : base;
  }, [contentClassName]);

  return (
    <div className={wrapperClass}>
      <div ref={viewportRef} className={viewportClass}>
        <div ref={contentRef} className={contentClass}>
          {children}
        </div>
      </div>
      {metrics.canScroll ? (
        <div
          className={`scroll-area__bar${dragging ? ' is-dragging' : ''}`}
          onMouseDown={handleTrackMouseDown}
        >
          <div
            className="scroll-area__thumb"
            style={{height: `${metrics.thumbHeight}px`, transform: `translateY(${metrics.thumbTop}px)`}}
            onMouseDown={handleThumbMouseDown}
          />
        </div>
      ) : null}
    </div>
  );
}
