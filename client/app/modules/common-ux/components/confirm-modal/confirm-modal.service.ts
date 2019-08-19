import { Injectable, TemplateRef } from '@angular/core';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';


export interface ConfirmOptions {
  title: string;
  message: string;
  template: TemplateRef<any>;
  confirmLabel: string;
  rejectLabel: string;
  hideReject: boolean;
  modalOptions: NgbModalOptions;
}

@Injectable()
export class ConfirmState {
  options: ConfirmOptions;
  modal: NgbModalRef;
  template: TemplateRef<any>;
}

@Injectable({ providedIn: 'root' })
export class ConfirmModalService {

  constructor(private ngbModalService: NgbModal, private modalState: ConfirmState) {}

  public open(options: ConfirmOptions): Promise<any> {
    this.modalState.options = options;
    this.modalState.modal = this.ngbModalService.open(this.modalState.template, this.modalState.options.modalOptions);
    return this.modalState.modal.result;
  }
}
