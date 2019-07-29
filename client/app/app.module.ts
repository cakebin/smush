import { NgModule } from '@angular/core';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';


import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { CommonUXModule } from './modules/common-ux/common-ux.module';
import { AuthGuardService } from './app-auth-guard.service';
import { UserManagementModule } from './modules/user-management/user-management.module';
import { MatchManagementModule } from './modules/match-management/match-management.module';

import { UserManagementService } from './modules/user-management/user-management.service';
import { TopNavBarComponent } from './page-components/top-nav-bar/top-nav-bar.component';
import { MatchesComponent } from './page-components/matches/matches.component';
import { ProfileEditComponent } from './page-components/profiles/profile-edit.component';
import { InsightsComponent } from './page-components/insights/insights.component';
import { PageNotFoundComponent } from './page-components/page-not-found/page-not-found.component';
import { HomeComponent } from './page-components/home/home.component';
import { ChartsModule } from './modules/charts/charts.module';
import { AuthInterceptor } from './auth.interceptor';


@NgModule({
  declarations: [
    AppComponent,
    PageNotFoundComponent,
    TopNavBarComponent,
    HomeComponent,
    MatchesComponent,
    ProfileEditComponent,
    InsightsComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    UserManagementModule.forRoot(),
    MatchManagementModule.forRoot(),
    CommonUXModule.forRoot(),
    ChartsModule,
  ],
  bootstrap: [
    AppComponent,
  ],
  providers: [
    UserManagementService,
    AuthGuardService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true,
    },
  ]
})
export class AppModule {}
