import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { CommonUXModule } from './modules/common-ux/common-ux.module';
import { MatchManagementModule } from './modules/match-management/match-management.module';
import { TopNavBar } from './components/top-nav-bar/top-nav-bar.component';

@NgModule({
  declarations: [
    AppComponent,
    TopNavBar,
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
