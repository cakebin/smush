import { Directive, TemplateRef } from '@angular/core';
import { ConfirmState } from './confirm-modal.service';

@Directive({
  selector: '[confirmModalTemplate]'
})
export class ConfirmModalDirective {
  constructor(confirmTemplate: TemplateRef<any>, modalState: ConfirmState) {
    modalState.template = confirmTemplate;
  }
}
