import { Component, OnInit } from '@angular/core';
import { NgbDate } from '@ng-bootstrap/ng-bootstrap';
import { SingleSeries, MultiSeries, DataItem, Series } from '@swimlane/ngx-charts';
import { faCalendarAlt } from '@fortawesome/free-solid-svg-icons';

import { MatchManagementService } from 'client/app/modules/match-management/match-management.service';
import { IMatchViewModel, IUserViewModel, IChartViewModel, ChartViewModel, IChartUserViewModel } from 'client/app/app.view-models';
import { CommonUxService } from 'client/app/modules/common-ux/common-ux.service';
import { UserManagementService } from 'client/app/modules/user-management/user-management.service';

const charts: IChartViewModel[] = [
  new ChartViewModel(1, 'Opponent character usage', 'bar'),
  new ChartViewModel(2, 'User GSP over time', 'line'),
];

@Component({
  selector: 'insights',
  templateUrl: './insights.component.html',
  styleUrls: ['./insights.component.css']
})
export class InsightsComponent implements OnInit {
  public user: IUserViewModel;
  public matches: IMatchViewModel[] = [];
  public users: IChartUserViewModel[] = [];

  public charts = charts;
  public selectedChart: IChartViewModel = this.charts[0];

  public startDate: NgbDate;
  public endDate: NgbDate;
  public selectedUser: IChartUserViewModel = null;
  public sortType: string = 'use';
  public sortOrder: string = 'desc';

  public chartData: SingleSeries = [];
  public multiSeriesChartData: MultiSeries = [];
  public dataUnit: string = '';
  public xAxisLabel: string = '';
  public yAxisLabel: string = '';
  public xAxisTickFormatting: (val: string) => string;
  public yAxisTickFormatting: (val: string) => string;

  public noFilteredDataToDisplay: boolean = false;
  public isInitialLoad: boolean = true;
  public faCalendarAlt = faCalendarAlt;
  public fillerPercents = [65, 100, 26, 70, 30, 27, 22, 15, 30, 60, 95];

  constructor(
    private matchService: MatchManagementService,
    private commonUxService: CommonUxService,
    private userService: UserManagementService,
    ) {
  }

  ngOnInit() {
    this.matchService.cachedMatches.subscribe(res => {
      if (res && res.length) {
        this.matches = res;
        this.updateAndRenderChart();
        this.isInitialLoad = false;
      }
    },
    err => {
        this.commonUxService.showDangerToast('Unable to get data.');
        console.error(err);
    });
    this.userService.cachedUser.subscribe(res => {
      if (res) {
        this.user = res;
        this.users = [
          { userId: res.userId, userName: res.userName } as IChartUserViewModel
        ];
      }
    });
    this.userService.getAllUsers().subscribe((res: IChartUserViewModel[]) => {
      if (res) {
        this.users = res;
      }
    });
  }


  /*---------------------
          Filters
  ----------------------*/

  public onDateSelect(event: any): void {
    this.updateAndRenderChart();
  }


  /*---------------------
      Chart publishers
  ----------------------*/

  public updateAndRenderChart() {
    this._clearData();
    const chart = this.selectedChart;

    switch (chart.chartId) {
      case 1:
        this.chartData = this._getCharacterUsageChartData();
        this.dataUnit = 'percent';
        this.xAxisLabel = 'Usage';
        this.yAxisLabel = 'Character';
        this.xAxisTickFormatting = (val: string) => val + '%';
        this.yAxisTickFormatting = (val: string) => val;

        if (this.sortType == null) {
          this.sortType = 'use';
        }
        if (this.sortOrder == null) {
          this.sortOrder = 'desc';
        }

        break;
      case 2:
        if (this.selectedUser == null) {
          this.selectedUser = this.users.find(u => u.userId === this.user.userId);
        }
        this.multiSeriesChartData = this._getUserGspChartData();
        this.xAxisLabel = 'Date';
        this.yAxisLabel = 'GSP';
        this.xAxisTickFormatting = (val: string) => this._getFormattedDate(val);
        this.yAxisTickFormatting = (val: string) => parseInt(val, 10).toLocaleString();
        this.sortType = null;
        this.sortOrder = null;
        break;
    }
  }


  /*------------------
      Data methods
  -------------------*/
  private _getFilteredData(): IMatchViewModel[] {
    let filteredData: IMatchViewModel[] = [];
    Object.assign(filteredData, this.matches);

    // Filter the data based on user-given constraints
    if (this.selectedUser) {
      filteredData = filteredData.filter(match => match.userId === this.selectedUser.userId);
    }
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

    if (!filteredData || !filteredData.length) {
      this.noFilteredDataToDisplay = true;
      return;
    } else {
      this.noFilteredDataToDisplay = false;
    }

    return filteredData;
  }
  private _getUserGspChartData(): MultiSeries {
    const multiSeries: MultiSeries = [];
    const filteredData = this._getFilteredData();

    if (!filteredData) {
      return null;
    }

    filteredData.map(match => {
      if (!match.userCharacterName || !match.userCharacterGsp) {
        return null;
      }
      // Push match data into relevant character series
      let characterSeries = multiSeries.find(s => s.name === match.userCharacterName);
      if (!characterSeries) {
        const newIndex = multiSeries.push({
          name: match.userCharacterName,
          series: []
        } as Series);
        characterSeries = multiSeries[newIndex - 1];
      }
      characterSeries.series.push({
        name: new Date(match.created),
        value: match.userCharacterGsp as number
      } as DataItem);
    });

    return multiSeries;
  }

  private _getCharacterUsageChartData(): SingleSeries {
    let series: SingleSeries = [];
    const filteredData = this._getFilteredData();

    if (!filteredData) {
      return null;
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

    // Sort as specified by user
    series = this._sortSeries(series);

    return series;
  }

  private _sortSeries(series: SingleSeries): SingleSeries {
    let sortMultiplier;
    if (this.sortOrder === 'asc') {
      sortMultiplier = 1;
    }

    if (this.sortOrder === 'desc') {
      sortMultiplier = -1;
    }

    if (this.sortType === 'alpha') {
      series = series.sort((a: DataItem, b: DataItem) => {
        let sortValue;
        if (a.name > b.name) {
          sortValue = 1;
        } else if (a.name === b.name) {
          sortValue = 0;
        } else if (a.name < b.name) {
          sortValue = -1;
        }

        return sortValue * sortMultiplier;
      });
    }

    if (this.sortType === 'use') {
      series = series.sort((a: DataItem, b: DataItem) => {
        let sortValue;
        if (a.value > b.value) {
          sortValue = 1;
        } else if (a.value === b.value) {
          sortValue = 0;
        } else if (a.value < b.value) {
          sortValue = -1;
        }

        return sortValue * sortMultiplier;
      });
    }

    return series;
  }
  private _getFormattedDate(dateString: string, showTime: boolean = false) {
    const date: Date  = new Date(dateString);

    let month: number | string = date.getMonth() + 1;
    let day: number | string = date.getDate();
    let hour: number | string = date.getHours();
    let min: number | string = date.getMinutes();
    let sec: number | string = date.getSeconds();

    month = (month < 10 ? '0' : '') + month;
    day = (day < 10 ? '0' : '') + day;
    hour = (hour < 10 ? '0' : '') + hour;
    min = (min < 10 ? '0' : '') + min;
    sec = (sec < 10 ? '0' : '') + sec;

    let str: string = date.getFullYear() + '/' + month + '/' + day + ' ' +  hour + ':' + min + ':' + sec;

    if (!showTime) {
      str = date.getFullYear() + '/' + month + '/' + day;
    }

    return str;
  }
  private _clearData(): void {
    this.chartData = null;
    this.multiSeriesChartData = null;
  }
}
