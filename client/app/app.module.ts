import { NgModule } from '@angular/core';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';

// Root
import { AppComponent } from './app.component';

// Routing
import { AppRoutingModule } from './app-routing.module';
import { AuthGuardService } from './app-auth-guard.service';
import { AuthInterceptor } from './auth.interceptor';

// Modules
import { CommonUXModule } from './modules/common-ux/common-ux.module';
import { ChartsModule } from './modules/charts/charts.module';
import { UserManagementModule } from './modules/user-management/user-management.module';
import { MatchManagementModule } from './modules/match-management/match-management.module';
import { CharacterManagementModule } from './modules/character-management/character-management.module';

// Pages
import { TopNavBarComponent } from './page-components/top-nav-bar/top-nav-bar.component';
import { HomeComponent } from './page-components/home/home.component';
import { MatchesComponent } from './page-components/matches/matches.component';
import { InsightsComponent } from './page-components/insights/insights.component';
import { ProfileEditComponent } from './page-components/profiles/profile-edit.component';
import { AdminComponent } from './page-components/admin/admin.component';
import { PageNotFoundComponent } from './page-components/page-not-found/page-not-found.component';

// Services
import { CharacterManagementService } from 'client/app/modules/character-management/character-management.service';
import { MatchManagementService } from './modules/match-management/match-management.service';
import { UserManagementService } from './modules/user-management/user-management.service';
import { CommonUxService } from './modules/common-ux/common-ux.service';


@NgModule({
  declarations: [
    AppComponent,
    PageNotFoundComponent,
    TopNavBarComponent,
    HomeComponent,
    MatchesComponent,
    InsightsComponent,
    AdminComponent,
    ProfileEditComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    ChartsModule,
    CommonUXModule.forRoot(),
    UserManagementModule.forRoot(),
    MatchManagementModule.forRoot(),
    CharacterManagementModule.forRoot(),
  ],
  bootstrap: [
    AppComponent,
  ],
  providers: [
    AuthGuardService,
    CommonUxService,
    CharacterManagementService,
    MatchManagementService,
    UserManagementService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true,
    },
  ]
})
export class AppModule {}
