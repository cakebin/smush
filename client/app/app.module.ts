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

// Pages
import { TopNavBarComponent } from './page-components/top-nav-bar/top-nav-bar.component';
import { HomeComponent } from './page-components/home/home.component';
import { MatchesComponent } from './page-components/matches/matches.component';
import { InsightsComponent } from './page-components/insights/insights.component';
import { ProfileEditComponent } from './page-components/profiles/profile-edit.component';
import { PageNotFoundComponent } from './page-components/page-not-found/page-not-found.component';


@NgModule({
  declarations: [
    AppComponent,
    PageNotFoundComponent,
    TopNavBarComponent,
    HomeComponent,
    MatchesComponent,
    InsightsComponent,
    ProfileEditComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    ChartsModule,
    CommonUXModule.forRoot(),
    UserManagementModule.forRoot(),
    MatchManagementModule.forRoot(),
  ],
  bootstrap: [
    AppComponent,
  ],
  providers: [
    AuthGuardService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true,
    },
  ]
})
export class AppModule {}
