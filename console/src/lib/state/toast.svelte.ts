export type ToastType = 'success' | 'error' | 'info';

export interface Toast {
  id: string;
  type: ToastType;
  message: string;
  duration: number;
  createdAt: number;
}

let toasts = $state<Toast[]>([]);

function add(type: ToastType, message: string, duration = 4000) {
  // Deduplicate: don't show the same message if it's already visible
  const isDuplicate = toasts.some((t) => t.message === message && t.type === type);
  if (isDuplicate) return;

  const id = Math.random().toString(36).substring(2, 9);
  const toast: Toast = { id, type, message, duration, createdAt: Date.now() };
  toasts.push(toast);

  if (duration > 0) {
    setTimeout(() => remove(id), duration);
  }
}

function remove(id: string) {
  const index = toasts.findIndex((t) => t.id === id);
  if (index !== -1) toasts.splice(index, 1);
}

function clear() {
  toasts.splice(0, toasts.length);
}

export const toast = {
  get toasts() {
    return toasts;
  },
  success(message: string, duration = 4000) {
    add('success', message, duration);
  },
  error(message: string, duration = 6000) {
    add('error', message, duration);
  },
  info(message: string, duration = 4000) {
    add('info', message, duration);
  },
  remove,
  clear,
};
