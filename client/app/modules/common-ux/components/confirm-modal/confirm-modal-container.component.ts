import { Component } from '@angular/core';


@Component({
  selector: 'common-ux-confirmation-modal',
  template: `
    <div *confirmModalTemplate>
      <confirm-modal-content></confirm-modal-content>
    </div>
  `,
})
export class ConfirmModalContainerComponent {
}
