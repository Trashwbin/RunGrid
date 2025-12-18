import {type RefObject, useEffect} from 'react';

export function useOutsideClick<T extends HTMLElement>(
  ref: RefObject<T>,
  handler: () => void,
  active = true
) {
  useEffect(() => {
    if (!active) {
      return;
    }

    const onPointerDown = (event: MouseEvent) => {
      if (!ref.current || ref.current.contains(event.target as Node)) {
        return;
      }
      handler();
    };

    document.addEventListener('mousedown', onPointerDown);

    return () => {
      document.removeEventListener('mousedown', onPointerDown);
    };
  }, [active, handler, ref]);
}
