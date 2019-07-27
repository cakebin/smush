import { Component, Input } from '@angular/core';
import { SingleSeries } from '@swimlane/ngx-charts';

@Component({
  selector: 'chart-bar-horizontal',
  template: `
    <ngx-charts-bar-horizontal class="chart-container"
      [results]="data"
      [xAxisLabel]="xAxisLabel"
      [yAxisLabel]="yAxisLabel"
      [view]="[700,500]"
      [xAxis]="true"
      [yAxis]="true"
      [showXAxisLabel]="true"
      [showYAxisLabel]="true"
      [xAxisTickFormatting]="xAxisTickFormatting"
      [yAxisTickFormatting]="yAxisTickFormatting"
    >
    </ngx-charts-bar-horizontal>
  `,
  styles: [
    `
    .chart-container {
      position: absolute;
      left: 50%;
      margin-left:-40px;
      transform: translate(-50%);
    }
    `
  ]
})
export class ChartTestComponent {
  @Input() data: SingleSeries = [];
  @Input() xAxisLabel: string = '';
  @Input() yAxisLabel: string = '';
  @Input() xAxisTickFormatting;
  @Input() yAxisTickFormatting;

  public colorScheme = {
      name: 'solar',
      selectable: true,
      group: 'Continuous',
      domain: [
          '#fff8e1', '#ffecb3', '#ffe082', '#ffd54f', '#ffca28', '#ffc107', '#ffb300', '#ffa000', '#ff8f00', '#ff6f00'
      ]
    };

    public testData: SingleSeries = [
    {
      name: 'Germany',
      value: 40632,
      extra: {
        code: 'de'
      }
    },
    {
      name: 'United States',
      value: 50000,
      extra: {
        code: 'us'
      }
    },
    {
      name: 'France',
      value: 36745,
      extra: {
        code: 'fr'
      }
    },
    {
      name: 'United Kingdom',
      value: 36240,
      extra: {
        code: 'uk'
      }
    },
    {
      name: 'Spain',
      value: 33000,
      extra: {
        code: 'es'
      }
    },
    {
      name: 'Italy',
      value: 35800,
      extra: {
        code: 'it'
      }
    }
  ];

  constructor() {}
}
