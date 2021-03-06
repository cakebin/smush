import { NgModule, ModuleWithProviders } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

// Non-angular dependencies
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

// Confirmation modal
import { ConfirmModalModule } from './components/confirm-modal/confirm-modal.module';
import { ConfirmModalService } from './components/confirm-modal/confirm-modal.service';

// Toast
import { ToastModule } from './components/toast/toast.module';
import { ToastService } from './components/toast/toast.service';
import { ToastComponent } from './components/toast/toast.component';

// Slide panel
import { SlidePanelModule } from './components/slide-panel/slide-panel.module';

// Various components
import { SortableTableHeaderComponent } from './components/sortable-table-header/sortable-table-header.component';
import { MaskedNumberInputComponent } from './components/masked-number-input/masked-number-input.component';
import { TypeaheadComponent } from './components/typeahead/typeahead.component';
import { TagInputComponent } from './components/tag-input/tag-input.component';

// Pipes
import { ToNumberPipe } from './components/string-to-number/string-to-number.pipe';


@NgModule({
  declarations: [
    SortableTableHeaderComponent,
    TypeaheadComponent,
    MaskedNumberInputComponent,
    TagInputComponent,
    ToNumberPipe,
  ],
  exports: [
    FormsModule,
    NgbModule,
    CommonModule,
    FontAwesomeModule,
    SortableTableHeaderComponent,
    TypeaheadComponent,
    MaskedNumberInputComponent,
    ToastComponent,
    SlidePanelModule,
    TagInputComponent,
    ConfirmModalModule,
    ToNumberPipe,
  ],
  imports: [
    CommonModule,
    BrowserAnimationsModule,
    FormsModule,
    FontAwesomeModule,
    NgbModule,
    ToastModule,
    ConfirmModalModule,
  ],
})
export class CommonUXModule {
  static forRoot(): ModuleWithProviders {
    return {
        ngModule: CommonUXModule,
        providers: [
          ToastService,
          ConfirmModalService,
        ]
    };
  }
}
