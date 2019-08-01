import { HttpClientModule } from '@angular/common/http';
import { DecimalPipe } from '@angular/common';
import { NgModule, ModuleWithProviders } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';

import { CommonUXModule } from '../common-ux/common-ux.module';
import { CommonUxService } from '../common-ux/common-ux.service';

import { UserManagementService } from './user-management.service';


@NgModule({
  declarations: [],
  imports: [
    FormsModule,
    BrowserModule,
    HttpClientModule,
    CommonUXModule.forRoot(),
  ],
  exports: []
})
export class UserManagementModule {
  static forRoot(): ModuleWithProviders {
    return {
        ngModule: UserManagementModule,
        providers: [
          CommonUxService,
          UserManagementService,
          DecimalPipe,
          {
            provide: 'UserApiUrl',
            useValue: '/api/user'
          },
          {
            provide: 'AuthApiUrl',
            useValue: '/api/auth'
          },
        ]
    };
  }
}
