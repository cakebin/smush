import { Component, OnInit, ViewChildren, QueryList } from '@angular/core';
import { faCheck, faTimes, faTrash, faPencilAlt } from '@fortawesome/free-solid-svg-icons';

import { IMatchViewModel } from '../../app.view-models';
import { CommonUXService } from '../common-ux/common-ux.service';
import { ISortEvent, SortEvent, SortDirection, HeaderViewModel } from '../common-ux/common-ux.view-models';
import { SortableTableHeaderComponent } from '../common-ux/components/sortable-table-header/sortable-table-header.component';
import { MatchManagementService } from './match-management.service';

@Component({
  selector: 'match-table-view',
  templateUrl: './match-table-view.component.html',
})
export class MatchTableViewComponent implements OnInit {
  public headerLabels:HeaderViewModel[] = [
    new HeaderViewModel('id', '#'),
    new HeaderViewModel('userCharacterName', 'Char'),
    new HeaderViewModel('userCharacterGsp', 'GSP'),
    new HeaderViewModel('opponentCharacterName', 'Opponent Char'),
    new HeaderViewModel('opponentCharacterGsp', 'Opponent GSP'),
    new HeaderViewModel('userWin', 'Win/Loss'),
    new HeaderViewModel('opponentAwesome', 'Chum'),
    new HeaderViewModel('opponentCamp', 'Camp'),
    new HeaderViewModel('opponentTeabag', 'TBag'),
    new HeaderViewModel('created', 'Created'),
  ];
  @ViewChildren(SortableTableHeaderComponent) headerComponents: QueryList<SortableTableHeaderComponent>;
  
  public matches: IMatchViewModel[] = [];

  public sortedMatches: IMatchViewModel[];
  public sortColumnName: string = '';
  public sortColumnDirection: SortDirection = '';
  public isLoading: boolean = false;

  public faCheck = faCheck;
  public faTimes = faTimes;
  public faTrash = faTrash;
  public faPencilAlt = faPencilAlt;

  constructor(
    private commonUXService:CommonUXService,
    private matchManagementService: MatchManagementService,
    ){
  }

  ngOnInit() {
    this.isLoading = true;
    this.matchManagementService.cachedMatches.subscribe({
      next: res => {
        this.isLoading = true;
        //console.warn("Got some new matches from the service:", res);
        this.sortedMatches = this.matches = res;
        this.initialSort();
      },
      error: err => {
        this.commonUXService.showDangerToast('Unable to get matches.');
        console.error(err);
      },
      complete: () => {
        this.isLoading = false;
      }
    });
  }

  public onSort({column, direction}: ISortEvent) {
    // Resetting all headers. This needs to be done in a parent, no way around it
    this.headerComponents.forEach(header => {
      if(header.propertyName !== column){
        header.clearDirection();
      }
    });

    // Sorting items
    if (direction === '') {
      this.sortColumnName = '';
      this.sortColumnDirection = '';
      this.sortedMatches = this.matches;
    } else {
      this.sortColumnName = column;
      this.sortColumnDirection = direction;     
      this.sortedMatches = [...this.matches].sort((a, b) => {
        const res = this.commonUXService.compare(a[column], b[column]);
        return direction === 'asc' ? res : -res;
      });
    }
  }

  public editMatch(match: IMatchViewModel): void {
    console.log('Editing match!', match);
  }
  public deleteMatch(matchId: number): void {
    console.log('DELETING match!', matchId);
  } 
  private initialSort(): void { 
    this.onSort(new SortEvent('id', 'desc'));
  }
}