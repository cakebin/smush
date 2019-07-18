import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { TypeaheadComponent } from './components/typeahead/typeahead.component';


@NgModule({
  declarations: [
    TypeaheadComponent,
  ],
  exports: [
    TypeaheadComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    NgbModule,
  ],
})
export class CommonUXModule { }
