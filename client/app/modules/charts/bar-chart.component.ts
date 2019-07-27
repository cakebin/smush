import { Component, Input } from '@angular/core';
import { SingleSeries } from '@swimlane/ngx-charts';

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

  public colorScheme = {
    name: 'fire',
      selectable: true,
      group: 'Ordinal',
      domain: [
          '#ff3d00', '#bf360c', '#ff8f00', '#ff6f00', '#ff5722', '#e65100', '#ffca28', '#ffab00'
      ]
    };

  constructor() {}
}
