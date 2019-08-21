import { Component, Input, OnInit } from '@angular/core';
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
  selector: 'chart-line',
  templateUrl: './line-chart.component.html',
  styleUrls: ['../charts.css']
})
export class LineChartComponent implements OnInit {
  @Input() data: SingleSeries = [];
  @Input() xAxisLabel: string = '';
  @Input() yAxisLabel: string = '';
  @Input() xAxisTickFormatting;
  @Input() yAxisTickFormatting;

  public colorScheme = oceanScheme;

  constructor() {}

  ngOnInit() {
  }
}
