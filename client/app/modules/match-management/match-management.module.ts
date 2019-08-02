import { HttpClientModule } from '@angular/common/http';
import { DecimalPipe } from '@angular/common';
import { NgModule, ModuleWithProviders } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { CommonUXModule } from '../common-ux/common-ux.module';
import { CommonUxService } from '../common-ux/common-ux.service';
import { CharacterManagementModule } from '../character-management/character-management.module';
import { CharacterManagementService } from '../character-management/character-management.service';

import { MatchInputFormComponent } from './match-input-form.component';
import { MatchTableViewComponent } from './match-table-view.component';
import { MatchManagementService } from './match-management.service';



@NgModule({
  declarations: [
    MatchInputFormComponent,
    MatchTableViewComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    CharacterManagementModule.forRoot(),
    CommonUXModule.forRoot(),
  ],
  exports: [
    MatchInputFormComponent,
    MatchTableViewComponent,
  ]
})
export class MatchManagementModule {
  // Put providers here instead of above so they are only loaded once.
  static forRoot(): ModuleWithProviders {
    return {
        ngModule: MatchManagementModule,
        providers: [
          CommonUxService,
          MatchManagementService,
          CharacterManagementService,
          DecimalPipe,
          {
            provide: 'MatchApiUrl',
            useValue: '/api/match'
          },
        ]
    };
  }
}
