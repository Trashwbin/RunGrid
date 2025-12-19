import {useEffect} from 'react';
import {useToastStore, type ToastState, type ToastTone} from '../../store/overlays';

const defaultDurations: Record<ToastTone, number> = {
  success: 2400,
  info: 2600,
  warning: 3200,
  error: 4200,
};

export function ToastHost() {
  const toasts = useToastStore((state) => state.toasts);

  if (toasts.length === 0) {
    return null;
  }

  return (
    <div className="toast-host" role="region" aria-live="polite">
      {toasts.map((toast) => (
        <ToastItem key={toast.id} toast={toast} />
      ))}
    </div>
  );
}

type ToastItemProps = {
  toast: ToastState;
};

function ToastItem({toast}: ToastItemProps) {
  const dismiss = useToastStore((state) => state.dismiss);
  const duration = toast.duration ?? defaultDurations[toast.type];

  useEffect(() => {
    if (duration <= 0) {
      return;
    }
    const timer = window.setTimeout(() => dismiss(toast.id), duration);
    return () => window.clearTimeout(timer);
  }, [duration, dismiss, toast.id]);

  return (
    <div
      className={`toast toast--${toast.type}`}
      role="status"
      onClick={() => dismiss(toast.id)}
    >
      <span className={`toast-dot toast-dot--${toast.type}`} aria-hidden="true" />
      <div className="toast-body">
        <div className="toast-title">{toast.title}</div>
        {toast.message ? (
          <div className="toast-message">{toast.message}</div>
        ) : null}
      </div>
    </div>
  );
}
