import {useCallback, useEffect, useMemo} from 'react';
import {useModalStore, type ModalState} from '../../store/overlays';
import {Icon} from '../ui/Icon';

const sizeClass: Record<NonNullable<ModalState['size']>, string> = {
  sm: 'modal-card--sm',
  md: 'modal-card--md',
  lg: 'modal-card--lg',
};

export function ModalHost() {
  const modals = useModalStore((state) => state.modals);
  const closeModal = useModalStore((state) => state.closeModal);

  useEffect(() => {
    if (modals.length === 0) {
      return;
    }
    const handleKey = (event: KeyboardEvent) => {
      if (event.key !== 'Escape') {
        return;
      }
      const top = modals[modals.length - 1];
      if (!top || top.closable === false) {
        return;
      }
      closeModal(top.id);
    };
    window.addEventListener('keydown', handleKey);
    return () => window.removeEventListener('keydown', handleKey);
  }, [modals, closeModal]);

  if (modals.length === 0) {
    return null;
  }

  return (
    <div className="modal-host">
      {modals.map((modal, index) => (
        <ModalView
          key={modal.id}
          modal={modal}
          zIndex={200 + index * 2}
          onClose={closeModal}
        />
      ))}
    </div>
  );
}

type ModalViewProps = {
  modal: ModalState;
  zIndex: number;
  onClose: (id: string) => void;
};

function ModalView({modal, zIndex, onClose}: ModalViewProps) {
  const size = modal.size ?? 'md';
  const primaryLabel =
    modal.primaryLabel ?? (modal.kind === 'error' ? '知道了' : '确定');
  const secondaryLabel =
    modal.secondaryLabel ?? (modal.kind === 'confirm' ? '取消' : undefined);
  const showActions =
    modal.kind !== 'progress' &&
    (modal.primaryLabel ||
      modal.secondaryLabel ||
      modal.kind === 'confirm' ||
      modal.kind === 'error');
  const canClose = modal.closable !== false;

  const handleBackdrop = useCallback(() => {
    if (!canClose || modal.backdropClose === false) {
      return;
    }
    onClose(modal.id);
  }, [canClose, modal.backdropClose, modal.id, onClose]);

  const handleConfirm = useCallback(async () => {
    if (modal.onConfirm) {
      await modal.onConfirm();
    }
    if (modal.autoClose !== false) {
      onClose(modal.id);
    }
  }, [modal, onClose]);

  const handleCancel = useCallback(async () => {
    if (modal.onCancel) {
      await modal.onCancel();
    }
    if (modal.autoClose !== false) {
      onClose(modal.id);
    }
  }, [modal, onClose]);

  const body = useMemo(() => renderModalBody(modal), [modal]);

  return (
    <div
      className="modal-backdrop"
      style={{zIndex}}
      onClick={handleBackdrop}
      role="presentation"
    >
      <div
        className={`modal-card ${sizeClass[size]} modal-card--${modal.kind}`}
        role="dialog"
        aria-modal="true"
        aria-label={modal.title}
        onClick={(event) => event.stopPropagation()}
      >
        <div className="modal-header">
          <div className="modal-header-text">
            <h2 className="modal-title">{modal.title}</h2>
            {modal.description &&
            modal.kind !== 'progress' &&
            modal.kind !== 'error' ? (
              <p className="modal-description">{modal.description}</p>
            ) : null}
          </div>
          {canClose ? (
            <button
              type="button"
              className="modal-close"
              onClick={() => onClose(modal.id)}
              aria-label="关闭"
            >
              <Icon name="close" size={16} />
            </button>
          ) : null}
        </div>
        <div className="modal-body">{body}</div>
        {showActions ? (
          <div className="modal-actions">
            {secondaryLabel ? (
              <button
                type="button"
                className="modal-button modal-button--secondary"
                onClick={handleCancel}
              >
                {secondaryLabel}
              </button>
            ) : null}
            <button
              type="button"
              className={`modal-button modal-button--primary${modal.tone === 'danger' ? ' modal-button--danger' : ''}`}
              onClick={handleConfirm}
            >
              {primaryLabel}
            </button>
          </div>
        ) : null}
      </div>
    </div>
  );
}

function renderModalBody(modal: ModalState) {
  if (modal.content) {
    return modal.content;
  }

  if (modal.kind === 'progress') {
    const progress =
      typeof modal.progress === 'number'
        ? Math.min(100, Math.max(0, modal.progress))
        : null;
    const isIndeterminate = progress === null;
    return (
      <div className="modal-progress">
        {modal.description ? (
          <p className="modal-description">{modal.description}</p>
        ) : null}
        <div
          className={`progress-track ${isIndeterminate ? 'is-indeterminate' : ''}`}
        >
          <div
            className="progress-fill"
            style={progress !== null ? {width: `${progress}%`} : undefined}
          />
        </div>
        <div className="progress-meta">
          <span className="progress-percent">
            {progress !== null ? `${progress}%` : '正在扫描...'}
          </span>
          {modal.path ? (
            <span className="progress-path" title={modal.path}>
              {modal.path}
            </span>
          ) : null}
        </div>
      </div>
    );
  }

  if (modal.kind === 'error') {
    return (
      <div className="modal-error">
        {modal.description ? (
          <p className="modal-description">{modal.description}</p>
        ) : null}
        {modal.details ? (
          <pre className="modal-code">{modal.details}</pre>
        ) : null}
      </div>
    );
  }

  return null;
}
