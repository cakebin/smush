import { Component, OnInit, ViewChildren, QueryList } from '@angular/core';
import { of } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { faCheck, faTimes, faTrash, faPencilAlt } from '@fortawesome/free-solid-svg-icons';

import { IMatchViewModel } from '../../app.view-models';
import { CommonUXService } from '../common-ux/common-ux.service';
import { ISortEvent, HeaderViewModel } from '../common-ux/common-ux.view-models';
import { SortableTableHeaderComponent } from '../common-ux/components/sortable-table-header/sortable-table-header.component';
import { MatchManagementService } from './match-management.service';

@Component({
  selector: 'match-table-view',
  templateUrl: './match-table-view.component.html',
})
export class MatchTableViewComponent implements OnInit {
  public headerLabels:HeaderViewModel[] = [
    new HeaderViewModel('matchId', '#'),
    new HeaderViewModel('userCharacterName', 'User Character'),
    new HeaderViewModel('userCharacterGsp', 'User GSP'),
    new HeaderViewModel('opponentCharacterName', 'Opponent Character'),
    new HeaderViewModel('opponentCharacterGsp', 'Opponent GSP'),
    new HeaderViewModel('userWin', 'Win/Loss'),
    new HeaderViewModel('opponentAwesome', 'Chum'),
    new HeaderViewModel('opponentCamp', 'Camp'),
    new HeaderViewModel('opponentTeabag', 'TBag'),
  ];
  @ViewChildren(SortableTableHeaderComponent) headerComponents: QueryList<SortableTableHeaderComponent>;
  
  public matches: IMatchViewModel[] = [];

  public sortedMatches: IMatchViewModel[];
  public sortColumnName: string = '';
  public sortColumnDirection: string = '';
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
  
  // Watch all matches and update this component when they change
  private allMatches = this.matchManagementService.getAllMatches().pipe(
    map(res => {
      if (!res) {
        throw new Error('Matches expected!');
      }
      return res;
    }),
      catchError(err => of([]))
    );

  ngOnInit() {
    this.isLoading = true;
    this.subscribeToMatches(); 
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

  private subscribeToMatches(){
    this.allMatches.subscribe(
      result => {
        if(result){
          this.sortedMatches = this.matches = result;
        }
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
