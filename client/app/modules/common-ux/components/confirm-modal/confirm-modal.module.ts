import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { ConfirmModalDirective } from './confirm-modal.directive';
import { ConfirmModalComponent } from './confirm-modal.component';
import { ConfirmModalContainerComponent } from './confirm-modal-container.component';
import { ConfirmState, ConfirmModalService } from './confirm-modal.service';

// Adapted from:
// https://gist.github.com/jnizet/15c7a0ab4188c9ce6c79ca9840c71c4e

@NgModule({
  declarations: [
    ConfirmModalDirective,
    ConfirmModalComponent,
    ConfirmModalContainerComponent,
  ],
  entryComponents: [
    ConfirmModalComponent,
  ],
  exports: [
    // For cleaner markup, I am only exporting the container that holds the modal structure
    ConfirmModalContainerComponent,
  ],
  imports: [
    CommonModule,
    NgbModule,
  ],
  providers: [
    ConfirmModalService,
    ConfirmState,
  ]
})
export class ConfirmModalModule { }
