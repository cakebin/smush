import { NgModule, ModuleWithProviders } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { CommonModule } from '@angular/common';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { NgxChartsModule } from '@swimlane/ngx-charts';
import { BarChartComponent } from './bar-chart.component';

@NgModule({
  declarations: [
    BarChartComponent,
  ],
  exports: [
    BarChartComponent,
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
