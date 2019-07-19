import { NgModule, ModuleWithProviders } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { CommonUXModule } from './modules/common-ux/common-ux.module';
import { MatchManagementModule } from './modules/match-management/match-management.module';


@NgModule({
  declarations: [
    AppComponent,
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
