import { HttpClientModule } from '@angular/common/http';
import { DecimalPipe } from '@angular/common';
import { NgModule, ModuleWithProviders } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';

import { CommonUXModule } from '../common-ux/common-ux.module';
import { CommonUXService } from '../common-ux/common-ux.service';

import { UserEditFormComponent } from './user-edit-form.component';
import { UserViewComponent } from './user-view.component';
import { UserManagementService } from './user-management.service';


@NgModule({
  declarations: [
    UserEditFormComponent,
    UserViewComponent,
  ],
  imports: [
    FormsModule,
    BrowserModule,
    HttpClientModule,
    CommonUXModule.forRoot(),
  ],
  exports: [
    UserEditFormComponent,
    UserViewComponent,
  ]
})
export class UserManagementModule {
  static forRoot(): ModuleWithProviders {
    return {
        ngModule: UserManagementModule,
        providers: [
          CommonUXService,
          UserManagementService,
          DecimalPipe,
          {
            provide: 'ApiUrl',
            useValue: '/api/user'
          },
        ]
    }
  }
}
