import { Component, OnInit } from '@angular/core';
import { ConfirmOptions, ConfirmState } from './confirm-modal.service';


@Component({
  selector: 'confirm-modal-content',
  templateUrl: './confirm-modal.component.html',
})
export class ConfirmModalComponent implements OnInit {
  public options: ConfirmOptions;

  constructor(private modalState: ConfirmState) {
    this.options = this.modalState.options;
  }

  ngOnInit() {
    this.options = this.modalState.options;
  }

  yes() {
    this.modalState.modal.close('confirmed');
  }

  no() {
    this.modalState.modal.dismiss('rejected');
  }
}
