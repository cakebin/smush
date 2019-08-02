import { HttpClientModule } from '@angular/common/http';
import { NgModule, ModuleWithProviders } from '@angular/core';

import { CommonUxService } from '../common-ux/common-ux.service';
import { UserManagementService } from './user-management.service';


@NgModule({
  declarations: [],
  imports: [
    HttpClientModule,
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
