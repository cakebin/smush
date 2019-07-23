import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { CommonUXModule } from './modules/common-ux/common-ux.module';
import { MatchManagementModule } from './modules/match-management/match-management.module';
import { TopNavBar } from './components/top-nav-bar/top-nav-bar.component';
import { MatchesComponent } from './components/matches/matches.component';
import { ProfileViewComponent } from './components/profiles/profile-view.component';
import { ProfileEditComponent } from './components/profiles/profile-edit.component';
import { InsightsComponent } from './components/insights/insights.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';

@NgModule({
  declarations: [
    AppComponent,
    PageNotFoundComponent,
    TopNavBar,
    MatchesComponent,
    ProfileViewComponent,
    ProfileEditComponent,
    InsightsComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    MatchManagementModule.forRoot(),
    CommonUXModule.forRoot(),
  ],
  bootstrap: [
    AppComponent,
  ],
})
export class AppModule {}
