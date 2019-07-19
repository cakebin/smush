import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { CommonUXModule } from './modules/common-ux/common-ux.module';
import { CommonUXService } from './modules/common-ux/common-ux.service';

import { MatchInputFormComponent } from './components/match-input-form.component';

@NgModule({
  declarations: [
    AppComponent,
    MatchInputFormComponent,
  ],
  imports: [
    NgbModule,
    FormsModule,
    BrowserModule,
    AppRoutingModule,
    CommonUXModule.forRoot()
  ],
  providers: [
    CommonUXService,
  ],
  bootstrap: [
    AppComponent,
  ],
})
export class AppModule { }
