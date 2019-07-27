import { NgModule, ModuleWithProviders } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { TypeaheadComponent } from './components/typeahead/typeahead.component';
import { ToastModule } from './components/toast/toast.module';
import { ToastComponent } from './components/toast/toast.component';
import { MaskedNumberInputComponent } from './components/masked-number-input/masked-number-input.component';
import { ToastService } from './components/toast/toast.service';
import { SortableTableHeaderDirective } from './directives/sortable-table-header.directive';
import { SortableTableHeaderComponent } from './components/sortable-table-header/sortable-table-header.component';

@NgModule({
  declarations: [
    TypeaheadComponent,
    MaskedNumberInputComponent,
    SortableTableHeaderDirective,
    SortableTableHeaderComponent,
  ],
  exports: [
    FormsModule,
    NgbModule,
    FontAwesomeModule,
    TypeaheadComponent,
    MaskedNumberInputComponent,
    ToastComponent,
    SortableTableHeaderDirective,
    SortableTableHeaderComponent,
  ],
  imports: [
    FontAwesomeModule,
    CommonModule,
    FormsModule,
    NgbModule,
    ToastModule,
  ],
})
export class CommonUXModule {
  static forRoot(): ModuleWithProviders {
    return {
        ngModule: CommonUXModule,
        providers: [
          ToastService,
        ]
    };
  }
}
