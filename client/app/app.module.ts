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
import { TagManagementModule } from './modules/tag-management/tag-management.module';

// Pages
import { TopNavBarComponent } from './page-components/top-nav-bar/top-nav-bar.component';
import { HomeComponent } from './page-components/home/home.component';
import { MatchesComponent } from './page-components/matches/matches.component';
import { InsightsComponent } from './page-components/insights/insights.component';
import { AdminComponent } from './page-components/admin/admin.component';
import { AdminTagComponent } from './page-components/admin/admin-tag/admin-tag.component';
import { AdminCharacterComponent } from './page-components/admin/admin-character/admin-character.component';
import { ProfileComponent } from './page-components/profiles/profile.component';
import { PageNotFoundComponent } from './page-components/page-not-found/page-not-found.component';

// Services
import { TagManagementService } from 'client/app/modules/tag-management/tag-management.service';
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
    ProfileComponent,
    AdminTagComponent,
    AdminCharacterComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    ChartsModule,
    CommonUXModule.forRoot(),
    UserManagementModule.forRoot(),
    MatchManagementModule.forRoot(),
    CharacterManagementModule.forRoot(),
    TagManagementModule.forRoot(),
  ],
  bootstrap: [
    AppComponent,
  ],
  providers: [
    AuthGuardService,
    CommonUxService,
    CharacterManagementService,
    TagManagementService,
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
