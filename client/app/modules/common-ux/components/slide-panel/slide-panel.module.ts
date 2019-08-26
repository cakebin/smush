import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { SlidePanelComponent } from './slide-panel.component';
import { SlidePanelService } from './slide-panel.service';

// Janky thing written from scratch, not very generic for now
@NgModule({
  declarations: [
    SlidePanelComponent,
  ],
  exports: [
    SlidePanelComponent,
  ],
  imports: [
    CommonModule,
    NgbModule,
    BrowserAnimationsModule,
  ],
  providers: [
    SlidePanelService,
  ]
})
export class SlidePanelModule { }
