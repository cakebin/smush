import { Component, Input } from '@angular/core';
import { SingleSeries } from '@swimlane/ngx-charts';

const oceanScheme =  {
  name: 'ocean',
  selectable: false,
  group: 'Ordinal',
  domain: [
      '#1D68FB', '#33C0FC', '#4AFFFE', '#AFFFFF', '#FFFC63', '#FDBD2D', '#FC8A25', '#FA4F1E', '#FA141B', '#BA38D1'
  ]
};

@Component({
  selector: 'chart-bar-horizontal',
  templateUrl: './bar-chart.component.html',
  styleUrls: ['./bar-chart.component.css']
})
export class BarChartComponent {
  @Input() data: SingleSeries = [];
  @Input() xAxisLabel: string = '';
  @Input() yAxisLabel: string = '';
  @Input() xAxisTickFormatting;
  @Input() yAxisTickFormatting;
  // Weird-ass template handling: https://github.com/swimlane/ngx-charts/issues/736
  @Input() dataUnit: string = '';

  public colorScheme = oceanScheme;

  constructor() {}
}
