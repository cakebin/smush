import { HttpClientModule } from '@angular/common/http';
import { NgModule, ModuleWithProviders, ApplicationRef } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';


import { CommonUXModule } from '../common-ux/common-ux.module';
import { CommonUXService } from '../common-ux/common-ux.service';

import { MatchInputFormComponent } from './match-input-form.component';
import { MatchManagementService } from './match-management.service';


@NgModule({
  declarations: [
    MatchInputFormComponent,
  ],
  imports: [
    NgbModule,
    FormsModule,
    BrowserModule,
    HttpClientModule,
    CommonUXModule.forRoot(),
  ],
  exports: [
    MatchInputFormComponent,
  ]
})
export class MatchManagementModule {
  // Put providers here instead of above so they are only loaded once.
  static forRoot(): ModuleWithProviders {
    return {
        ngModule: MatchManagementModule,
        providers: [
          CommonUXService,
          MatchManagementService,
          {
            provide: 'ApiUrl',
            useValue: '/api/matches'
          },
        ]
    }
  }
}
