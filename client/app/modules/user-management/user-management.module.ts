import { HttpClientModule } from '@angular/common/http';
import { NgModule, ModuleWithProviders } from '@angular/core';

import { CommonUXModule } from '../common-ux/common-ux.module';
import { CommonUxService } from '../common-ux/common-ux.service';
import { UserManagementService } from './user-management.service';

import { UserCharacterRowComponent } from './components/user-character-row.component';
import { ProfileEditComponent } from './components/profile-edit.component';


@NgModule({
  declarations: [
    UserCharacterRowComponent,
    ProfileEditComponent,
  ],
  exports: [
    ProfileEditComponent,
  ],
  imports: [
    HttpClientModule,
    CommonUXModule,
  ]
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
          {
            provide: 'UserCharacterApiUrl',
            useValue: '/api/user/character'
          },
        ]
    };
  }
}
