import { NgModule, ModuleWithProviders } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { TypeaheadComponent } from './components/typeahead/typeahead.component';
import { ToastModule } from './components/toast/toast.module';
import { ToastComponent } from './components/toast/toast.component';
import { ToastService } from './components/toast/toast.service';
import { NumberMaskDirective } from './directives/number-mask.directive';
import { SortableTableHeaderDirective } from './directives/sortable-table-header.directive';
import { SortableTableHeaderComponent } from './components/sortable-table-header/sortable-table-header.component';

@NgModule({
  declarations: [
    TypeaheadComponent,
    NumberMaskDirective,
    SortableTableHeaderDirective,
    SortableTableHeaderComponent,
  ],
  exports: [
    FormsModule,
    NgbModule,
    FontAwesomeModule,
    TypeaheadComponent,
    ToastComponent,
    NumberMaskDirective,
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
