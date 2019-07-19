import { NgModule, ModuleWithProviders } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { TypeaheadComponent } from './components/typeahead/typeahead.component';
import { ToastModule } from './components/toast/toast.module';
import { ToastComponent } from './components/toast/toast.component';
import { ToastService } from './components/toast/toast.service';
import { NumberMaskDirective } from './directives/number-mask/number-mask.directive';


@NgModule({
  declarations: [
    TypeaheadComponent,
    NumberMaskDirective,
  ],
  exports: [
    TypeaheadComponent,
    ToastComponent,
    NumberMaskDirective,
  ],
  imports: [
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
