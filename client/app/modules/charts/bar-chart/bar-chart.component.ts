import { Component, Input, OnInit } from '@angular/core';
import { SingleSeries } from '@swimlane/ngx-charts';

const SCHEME =  {
  name: 'ocean',
  selectable: false,
  group: 'Ordinal',
  domain: [
      '#1D68FB', '#33C0FC', '#4AFFFE', '#AFFFFF', '#FFFC63', '#FDBD2D', '#FC8A25', '#FA4F1E', '#FA141B', '#BA38D1'
  ]
};
const MINHEIGHT: number = 400;

@Component({
  selector: 'chart-bar-horizontal',
  templateUrl: './bar-chart.component.html',
  styleUrls: ['../charts.css']
})
export class BarChartComponent {
  private _data: SingleSeries = [];
  @Input() set data(value: SingleSeries) {
    this._data = value;
    const charHeight: number = 12;
    this.calcHeight = charHeight * this.data.length;
    if (this.calcHeight < MINHEIGHT) {
      this.calcHeight = MINHEIGHT;
    }
  }
  get data(): SingleSeries {
    return this._data;
  }

  @Input() xAxisLabel: string = '';
  @Input() yAxisLabel: string = '';
  @Input() xAxisTickFormatting;
  @Input() yAxisTickFormatting;
  // Weird-ass template handling: https://github.com/swimlane/ngx-charts/issues/736
  @Input() dataUnit: string = '';

  public colorScheme = SCHEME;
  public calcHeight: number = MINHEIGHT;

  constructor() {}
}
