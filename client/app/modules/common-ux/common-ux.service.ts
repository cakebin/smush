import { Injectable } from '@angular/core';
import { ToastService } from './components/toast/toast.service';


@Injectable()
export class CommonUXService {

    constructor(private toastService: ToastService) {
    }

    public showStandardToast(message: string): void {
        this.toastService.show(message);
    }
    public showSuccessToast(message: string): void {
        this.toastService.show(message, { classname: 'bg-success text-light', delay: 5000 });
    }
    public showWarningToast(message: string): void {
        this.toastService.show(message, { classname: 'bg-warning text-light', delay: 5000 });
    }
    public showDangerToast(message: string): void {
        this.toastService.show(message, { classname: 'bg-danger text-light', delay: 10000 });
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
    public debounce(functionToRun, delay: number, runImmediately: boolean = false): () => void {
        // Adapted from https://davidwalsh.name/javascript-debounce-function

        let timeout: NodeJS.Timer;
        return () => {
            const context = this;
            const args = arguments;

            const later = () => {
                timeout = null;
                if (!runImmediately) {
                    functionToRun.apply(context, args);
                }
            };
            const callNow: boolean = runImmediately && !timeout;
            clearTimeout(timeout);
            timeout = setTimeout(later, delay);
            if (callNow) {
                functionToRun.apply(context, args);
            }
        };
    }
}
