import { Injectable, TemplateRef } from '@angular/core';
import { ToastService } from './components/toast/toast.service';
import { ConfirmModalService, ConfirmOptions } from './components/confirm-modal/confirm-modal.service';

@Injectable()
export class CommonUxService {

    constructor(
        private toastService: ToastService,
        private confirmModalService: ConfirmModalService,
    ) {
    }

    public showStandardToast(message: string | TemplateRef<any>, titleText: string = ''): void {
        this.toastService.show(message, { title: titleText });
    }
    public showSuccessToast(message: string | TemplateRef<any>, titleText: string = 'Success'): void {
        this.toastService.show(message, { classname: 'bg-success text-light', delay: 5000, title: titleText });
    }
    public showWarningToast(message: string | TemplateRef<any>, titleText: string = 'Warning'): void {
        this.toastService.show(message, { classname: 'bg-warning text-light', delay: 5000, title: titleText });
    }
    public showDangerToast(message: string | TemplateRef<any>, titleText: string = 'Error'): void {
        this.toastService.show(message, { classname: 'bg-danger text-light', delay: 8000, title: titleText });
    }
    public openConfirmModal(
        messageText: string,
        titleText: string = 'Confirmation',
        hideRejectBool: boolean = false,
        confirmLabelText: string = '',
        rejectLabelText: string = ''): Promise<any> {
        return this.confirmModalService.open({
            message: messageText,
            title: titleText,
            hideReject: hideRejectBool,
            confirmLabel: confirmLabelText,
            rejectLabel: rejectLabelText } as ConfirmOptions);
    }

    // Utility methods
    public compare(v1: number, v2: number): number {
        if (v1 < v2) {
            return -1;
        } else if (v1 > v2) {
            return 1;
        } else if (v1 === v2) {
            return 0;
        } else if (v1 == null && v2 != null) {
            return -1;
        } else if (v1 != null && v2 == null) {
            return 1;
        } else {
            return 0;
        }
    }
}
