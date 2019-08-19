import { HttpClientModule } from '@angular/common/http';
import { DecimalPipe } from '@angular/common';
import { NgModule, ModuleWithProviders } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { CommonUXModule } from '../common-ux/common-ux.module';
import { CommonUxService } from '../common-ux/common-ux.service';

import { MatchInputFormComponent } from './match-input-form.component';
import { MatchViewComponent } from './match-view.component';
import { MatchRowComponent} from './components/match-row.component';
import { MatchCardComponent } from './components/match-card.component';
import { MatchCardEditComponent } from './components/match-card-edit.component';


import { MatchManagementService } from './match-management.service';
import { CharacterManagementModule } from '../character-management/character-management.module';
import { CharacterManagementService } from '../character-management/character-management.service';




@NgModule({
  declarations: [
    MatchInputFormComponent,
    MatchViewComponent,
    MatchRowComponent,
    MatchCardComponent,
    MatchCardEditComponent
  ],
  entryComponents: [
    MatchCardEditComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    CommonUXModule.forRoot(),
    CharacterManagementModule.forRoot(),
  ],
  exports: [
    MatchInputFormComponent,
    MatchViewComponent,
  ]
})
export class MatchManagementModule {
  // Put providers here instead of above so they are only loaded once.
  static forRoot(): ModuleWithProviders {
    return {
        ngModule: MatchManagementModule,
        providers: [
          CommonUxService,
          CharacterManagementService,
          MatchManagementService,
          DecimalPipe,
          {
            provide: 'MatchApiUrl',
            useValue: '/api/match'
          },
        ]
    };
  }
}
