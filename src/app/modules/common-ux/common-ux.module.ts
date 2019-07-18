import { NgModule, ModuleWithProviders } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { TypeaheadComponent } from './components/typeahead/typeahead.component';
import { ToastModule } from './components/toast/toast.module';
import { ToastComponent } from './components/toast/toast.component';
import { ToastService } from './components/toast/toast.service';


@NgModule({
  declarations: [
    TypeaheadComponent,
  ],
  exports: [
    TypeaheadComponent,
    ToastComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    NgbModule,
    ToastModule,
  ],
  providers: [
    ToastService,
  ]
})
export class CommonUXModule {
  static forRoot(): ModuleWithProviders {
    return {
        ngModule: CommonUXModule,
        providers: [
        ]
    };
  }
}
