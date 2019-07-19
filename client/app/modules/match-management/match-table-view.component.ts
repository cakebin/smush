import { Component, OnInit } from '@angular/core';
import { IMatchViewModel } from '../../app.view-models';
import { CommonUXService } from '../common-ux/common-ux.service';
import { MatchManagementService } from './match-management.service';

@Component({
  selector: 'match-table-view',
  templateUrl: './match-table-view.component.html',
})
export class MatchTableViewComponent implements OnInit {
  public matches: IMatchViewModel[] = [];
  public isLoading: boolean = false;

  constructor(
    private commonUXService:CommonUXService,
    private matchManagementService: MatchManagementService,
    ){
  }
  
  ngOnInit() {
    this.isLoading = true;

    this.matchManagementService.getAllMatches().subscribe(
      result => {
        if(result) this.matches = result;
      },
      error => {
        this.commonUXService.showDangerToast("Unable to get matches.");
      },
      () => {
        this.isLoading = false;
      }
    );
  }
}
