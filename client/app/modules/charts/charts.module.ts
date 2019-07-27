import { NgModule, ModuleWithProviders } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { CommonModule } from '@angular/common';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { NgxChartsModule } from '@swimlane/ngx-charts';
import { ChartTestComponent } from './bar-chart.component';

@NgModule({
  declarations: [
    ChartTestComponent,
  ],
  exports: [
    ChartTestComponent,
  ],
  imports: [
    CommonModule,
    NgbModule,
    BrowserAnimationsModule,
    NgxChartsModule,
  ],
})
export class ChartsModule {
  static forRoot(): ModuleWithProviders {
    return {
        ngModule: ChartsModule,
        providers: [
        ]
    };
  }
}
