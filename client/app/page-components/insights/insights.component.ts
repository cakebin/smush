import { Component, OnInit } from '@angular/core';
import { NgbDate } from '@ng-bootstrap/ng-bootstrap';
import { SingleSeries, DataItem } from '@swimlane/ngx-charts';
import { faCircleNotch, faCalendarAlt } from '@fortawesome/free-solid-svg-icons';

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
  public xAxisTickFormatting: (val: string) => string;
  public yAxisTickFormatting: (val: string) => string;

  private matches: IMatchViewModel[] = [];
  public startDate: NgbDate;
  public endDate: NgbDate;
  public chartUserId: number;

  public noFilteredDataToDisplay: boolean = false;
  public isLoading: boolean = false;
  public faCircleNotch = faCircleNotch;
  public faCalendarAlt = faCalendarAlt;

  constructor(
    private matchService: MatchManagementService,
    private commonUxService: CommonUXService,
    ) {
  }

  ngOnInit() {
    // It doesn't take long enough for this to load to warrant a spinner.
    // Will address later if this is actually an issue.
    // this.isLoading = true;
    this.matchService.cachedMatches.subscribe(res => {
      this.matches = res;
      this.publishCharacterUsageChartData();
      this.isLoading = false;
    },
    err => {
        this.commonUxService.showDangerToast('Unable to get data.');
        console.error(err);
        this.isLoading = false;
    });
  }
  public onDateSelect(event: any): void {
    this.publishCharacterUsageChartData();
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
      filteredData = filteredData.filter(match => {
        const matchCreateDate: Date = new Date(match.created);
        const matchNgbCreateDate: NgbDate = new NgbDate(
          matchCreateDate.getFullYear(),
          matchCreateDate.getMonth() + 1,
          matchCreateDate.getDate());
        return (matchNgbCreateDate.after(this.startDate) || matchNgbCreateDate.equals(this.startDate));
      });
    }
    if (this.endDate) {
      filteredData = filteredData.filter(match => {
        const matchCreateDate: Date = new Date(match.created);
        const matchNgbCreateDate: NgbDate = new NgbDate(
          matchCreateDate.getFullYear(),
          matchCreateDate.getMonth() + 1,
          matchCreateDate.getDate());
        return (matchNgbCreateDate.before(this.endDate) || matchNgbCreateDate.equals(this.endDate));
      });
    }
    if (this.chartUserId) {
      filteredData = filteredData.filter(match => match.userId === this.chartUserId);
    }
    if (!filteredData.length) {
      this.noFilteredDataToDisplay = true;
      return;
    } else {
      this.noFilteredDataToDisplay = false;
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

    // Sort by opponent character name alphabetically
    series = series.sort((a: DataItem, b: DataItem) => {
      if (a.name > b.name) {
        return 1;
      } else if (a.name === b.name) {
        return 0;
      } else if (a.name < b.name) {
        return -1;
      }
    });

    return series;
  }
}
