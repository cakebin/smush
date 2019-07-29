import { Component, TemplateRef } from '@angular/core';
import { ToastService } from './toast.service';


@Component({
  selector: 'common-ux-toasts',
  templateUrl: './toast.component.html',
  host: {'[class.ngb-toasts]': 'true'}
})
export class ToastComponent {
  constructor(public toastService: ToastService) {}
  isTemplate(toast: any) { return toast.textOrTpl instanceof TemplateRef; }
}
