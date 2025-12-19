import type {ReactNode} from 'react';
import {create} from 'zustand';

export type ModalKind = 'confirm' | 'form' | 'progress' | 'error' | 'custom';
export type ModalSize = 'sm' | 'md' | 'lg';
export type ModalTone = 'default' | 'danger';

export type ModalPayload = {
  id?: string;
  kind: ModalKind;
  title: string;
  description?: string;
  content?: ReactNode;
  size?: ModalSize;
  primaryLabel?: string;
  secondaryLabel?: string;
  tone?: ModalTone;
  progress?: number;
  path?: string;
  details?: string;
  closable?: boolean;
  backdropClose?: boolean;
  autoClose?: boolean;
  onConfirm?: () => void | Promise<void>;
  onCancel?: () => void | Promise<void>;
};

export type ModalState = Omit<ModalPayload, 'id'> & {id: string};

type ModalStore = {
  modals: ModalState[];
  openModal: (payload: ModalPayload) => string;
  updateModal: (id: string, patch: Partial<ModalPayload>) => void;
  closeModal: (id: string) => void;
  closeTop: () => void;
  clear: () => void;
};

const createId = () => {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
    return crypto.randomUUID();
  }
  return `modal_${Date.now()}_${Math.random().toString(16).slice(2)}`;
};

export const useModalStore = create<ModalStore>((set) => ({
  modals: [],
  openModal: (payload) => {
    const id = payload.id ?? createId();
    const next = {autoClose: true, backdropClose: true, ...payload, id};
    set((state) => ({modals: [...state.modals, next]}));
    return id;
  },
  updateModal: (id, patch) =>
    set((state) => ({
      modals: state.modals.map((modal) =>
        modal.id === id ? {...modal, ...patch} : modal
      ),
    })),
  closeModal: (id) =>
    set((state) => ({modals: state.modals.filter((modal) => modal.id !== id)})),
  closeTop: () =>
    set((state) => ({modals: state.modals.slice(0, -1)})),
  clear: () => set({modals: []}),
}));

export type ToastTone = 'success' | 'info' | 'warning' | 'error';

export type ToastPayload = {
  id?: string;
  type?: ToastTone;
  title: string;
  message?: string;
  duration?: number;
};

export type ToastState = {
  id: string;
  type: ToastTone;
  title: string;
  message?: string;
  duration?: number;
};

type ToastStore = {
  toasts: ToastState[];
  notify: (payload: ToastPayload) => string;
  dismiss: (id: string) => void;
  clear: () => void;
};

export const useToastStore = create<ToastStore>((set) => ({
  toasts: [],
  notify: (payload) => {
    const id = payload.id ?? createId();
    const next = {
      type: payload.type ?? 'info',
      ...payload,
      id,
    };
    set((state) => ({toasts: [...state.toasts, next]}));
    return id;
  },
  dismiss: (id) =>
    set((state) => ({toasts: state.toasts.filter((toast) => toast.id !== id)})),
  clear: () => set({toasts: []}),
}));
