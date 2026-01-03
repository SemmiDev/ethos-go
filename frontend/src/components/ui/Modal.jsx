import { Fragment } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import { X, AlertTriangle } from 'lucide-react';
import { Button } from './Button';

export function Modal({ isOpen, onClose, title, children, size = 'md' }) {
    const sizes = {
        sm: 'max-w-sm',
        md: 'max-w-md',
        lg: 'max-w-2xl',
        xl: 'max-w-4xl',
        full: 'max-w-full',
    };

    const sizeClass = sizes[size] || sizes.md;

    return (
        <Transition appear show={isOpen} as={Fragment}>
            <Dialog as="div" className="relative z-50" onClose={onClose}>
                <Transition.Child
                    as={Fragment}
                    enter="ease-out duration-200"
                    enterFrom="opacity-0"
                    enterTo="opacity-100"
                    leave="ease-in duration-150"
                    leaveFrom="opacity-100"
                    leaveTo="opacity-0"
                >
                    <div className="fixed inset-0 bg-neutral/60 backdrop-blur-sm" />
                </Transition.Child>

                <div className="fixed inset-0 overflow-y-auto">
                    <div className="flex min-h-full items-center justify-center p-4">
                        <Transition.Child
                            as={Fragment}
                            enter="ease-out duration-200"
                            enterFrom="opacity-0 scale-95"
                            enterTo="opacity-100 scale-100"
                            leave="ease-in duration-150"
                            leaveFrom="opacity-100 scale-100"
                            leaveTo="opacity-0 scale-95"
                        >
                            <Dialog.Panel className={`w-full ${sizeClass} bg-base-100 rounded-lg shadow-xl border border-base-300 overflow-hidden`}>
                                {/* Header */}
                                <div className="flex items-center justify-between px-6 py-4 border-b border-base-200">
                                    {title && (
                                        <Dialog.Title as="h3" className="text-lg font-semibold text-base-content">
                                            {title}
                                        </Dialog.Title>
                                    )}
                                    <button
                                        className="p-1.5 rounded-md text-base-content/50 hover:text-base-content hover:bg-base-200 transition-colors"
                                        onClick={onClose}
                                    >
                                        <X size={18} />
                                    </button>
                                </div>

                                {/* Body */}
                                <div className="px-6 py-5 text-base-content/80">{children}</div>
                            </Dialog.Panel>
                        </Transition.Child>
                    </div>
                </div>
            </Dialog>
        </Transition>
    );
}

export function ConfirmModal({
    isOpen,
    onClose,
    onConfirm,
    title = 'Confirm Action',
    message = 'Are you sure you want to proceed?',
    confirmText = 'Confirm',
    cancelText = 'Cancel',
    variant = 'danger',
    isLoading = false,
}) {
    return (
        <Modal isOpen={isOpen} onClose={onClose} size="sm">
            <div className="text-center py-2">
                <div className="flex justify-center mb-5">
                    <div className={`p-4 rounded-full ${variant === 'danger' ? 'bg-error/10' : 'bg-primary/10'}`}>
                        <AlertTriangle size={28} className={variant === 'danger' ? 'text-error' : 'text-primary'} />
                    </div>
                </div>
                <h3 className="text-lg font-semibold text-base-content mb-2">{title}</h3>
                <p className="text-sm text-base-content/60 mb-6 max-w-xs mx-auto">{message}</p>
                <div className="flex justify-center gap-3">
                    <Button variant="secondary" onClick={onClose} disabled={isLoading}>
                        {cancelText}
                    </Button>
                    <Button variant={variant} onClick={onConfirm} isLoading={isLoading}>
                        {confirmText}
                    </Button>
                </div>
            </div>
        </Modal>
    );
}
