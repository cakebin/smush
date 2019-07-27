import { Component, OnInit } from '@angular/core';
import { SingleSeries, DataItem } from '@swimlane/ngx-charts';
import { MatchManagementService } from 'client/app/modules/match-management/match-management.service';
import { IMatchViewModel } from 'client/app/app.view-models';
import { CommonUXService } from 'client/app/modules/common-ux/common-ux.service';

@Component({
  selector: 'insights',
  templateUrl: './insights.component.html',
  styles: [`
    .insights-chart-container {
      margin-top:50px;
    }
  `]
})
export class InsightsComponent implements OnInit {
  public chartData: SingleSeries = [];
  public dataUnit: string = '';
  public xAxisLabel: string = '';
  public yAxisLabel: string = '';
  public xAxisTickFormatting;
  public yAxisTickFormatting;

  private matches: IMatchViewModel[] = [];
  public startDate: Date;
  public endDate: Date;
  public chartUserId: number;

  public isLoading: boolean = false;

  constructor(
    private matchService: MatchManagementService,
    private commonUxService: CommonUXService,
    ) {
  }

  ngOnInit() {
    this.isLoading = true;
    this.matchService.cachedMatches.subscribe({
      next: res => {
        this.matches = res;
        this.publishCharacterUsageChartData();
      },
      error: err => {
        this.commonUxService.showDangerToast('Unable to get data.');
        console.error(err);
      },
      complete: () => {
        this.isLoading = false;
      }
    });
  }

  public publishCharacterUsageChartData() {
    this.chartData = this._getCharacterUsageChartData();
    this.dataUnit = 'percent';
    this.xAxisLabel = 'Usage';
    this.yAxisLabel = 'Character';
    this.xAxisTickFormatting = (val: string) => val + '%';
  }

  private _getCharacterUsageChartData(): SingleSeries {
    let filteredData: IMatchViewModel[] = this.matches;
    let series: SingleSeries = [];

    // Filter the data based on user-given constraints
    if (this.startDate) {
      filteredData = filteredData.filter(match => match.created >= this.startDate);
    }
    if (this.endDate) {
      filteredData = filteredData.filter(match => match.created <= this.endDate);
    }
    if (this.chartUserId) {
      filteredData = filteredData.filter(match => match.userId === this.chartUserId);
    }
    // Group by and transform into DataItem objects simultaneously (cringe)
    series = filteredData.reduce((dataItemArray: SingleSeries, match: IMatchViewModel) => {
      const opponentCharacterName = match.opponentCharacterName;
      const storageIndex = dataItemArray.findIndex((characterGroup: DataItem) => characterGroup.name === opponentCharacterName);
      if (storageIndex === - 1) { // We don't have this character group yet, so make it
        dataItemArray.push({ name: opponentCharacterName, value: 1 } as DataItem);
      } else {
        // This specific type can be multiple types, so we need to
        // assert that it's a number before we can add a number to it
        (dataItemArray[storageIndex].value as number)++;
      }
      return dataItemArray;
    }, []);

    // Go through the integer counts and turn into percents for more usable visualisation
    series = series.map((dataItem: DataItem) => { return {
      name: dataItem.name, value: ((dataItem.value as number) / filteredData.length) * 100 } as DataItem;
    });

    return series;
  }
}
