import { Component, OnInit, ViewChildren, QueryList } from '@angular/core';
import { faCheck, faTimes, faTrash, faPencilAlt } from '@fortawesome/free-solid-svg-icons';

import { IMatchViewModel, ICharacterViewModel, IUserViewModel } from '../../app.view-models';
import { CommonUxService } from '../common-ux/common-ux.service';
import { ISortEvent, SortEvent, SortDirection, HeaderViewModel } from '../common-ux/common-ux.view-models';
import { SortableTableHeaderComponent } from '../common-ux/components/sortable-table-header/sortable-table-header.component';
import { MatchManagementService } from './match-management.service';
import { CharacterManagementService } from '../character-management/character-management.service';
import { UserManagementService } from '../user-management/user-management.service';

@Component({
  selector: 'match-table-view',
  templateUrl: './match-table-view.component.html',
  styles: [`
    /deep/ td {
      white-space: nowrap;
    }
    .table-striped tbody tr.highlight {
      animation: highlight 1500ms ease-out;
    }
    @keyframes highlight {
      0% {
        background-color: #ffc107;
      }
      75% {
        background-color: #ffc107;
      }
      100 {
        background-color: initial;
      }
    }
  `]
})
export class MatchTableViewComponent implements OnInit {
  public headerLabels: HeaderViewModel[] = [
    new HeaderViewModel('userName', 'User', '100px'),
    new HeaderViewModel('userCharacterName', 'Char', '150px'),
    new HeaderViewModel('userCharacterGsp', 'GSP', '150px'),
    new HeaderViewModel('opponentCharacterName', 'Opponent Char', '150px'),
    new HeaderViewModel('opponentCharacterGsp', 'Opponent GSP', '150px'),
    new HeaderViewModel('userWin', 'Win', '80px'),
    new HeaderViewModel('opponentAwesome', 'Chum', '80px'),
    new HeaderViewModel('opponentCamp', 'Camp', '80px'),
    new HeaderViewModel('opponentTeabag', 'TBag', '80px'),
    new HeaderViewModel('created', 'Created', '120px'),
  ];
  @ViewChildren(SortableTableHeaderComponent) headerComponents: QueryList<SortableTableHeaderComponent>;

  public matches: IMatchViewModel[] = [];
  public user: IUserViewModel;
  public characters: ICharacterViewModel[] = [];

  public sortedMatches: IMatchViewModel[];
  public sortColumnName: string = '';
  public sortColumnDirection: SortDirection = '';
  public isInitialLoad: boolean = true;

  // Match editing
  public editedMatch: IMatchViewModel = {} as IMatchViewModel;

  public faCheck = faCheck;
  public faTimes = faTimes;
  public faTrash = faTrash;
  public faPencilAlt = faPencilAlt;

  constructor(
    private commonUXService: CommonUxService,
    private matchService: MatchManagementService,
    private characterService: CharacterManagementService,
    private userService: UserManagementService,
    ) {
  }

  ngOnInit() {
    this.matchService.cachedMatches.subscribe({
      next: res => {
        if (res) {
          this.sortedMatches = res;
          this.matches = res;

          if (this.isInitialLoad) {
            this._initialSort();
            this.isInitialLoad = false;
          } else {
            this._sort();
          }
        }
      },
      error: err => {
        this.commonUXService.showDangerToast('Unable to get matches.');
        console.error(err);
      }
    });
    this.characterService.characters.subscribe(
      (res: ICharacterViewModel[]) => {
        this.characters = res;
      });
    this.userService.cachedUser.subscribe(
      (res: IUserViewModel) => {
        this.user = res;
      });
  }

  public onSort({column, direction}: ISortEvent) {
    // Resetting all headers. This needs to be done in a parent, no way around it
    if (this.headerComponents) {
      this.headerComponents.forEach(header => {
        if (header.propertyName !== column) {
          header.clearDirection();
        }
      });
    }
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
  public onSelectEditUserCharacter(event: ICharacterViewModel): void {
    if (event) {
      this.editedMatch.userCharacterId = event.characterId;
    }
  }
  public onSelectEditOpponentCharacter(event: ICharacterViewModel): void {
    if (event) {
      this.editedMatch.opponentCharacterId = event.characterId;
    }
  }

  /*------------------------
       Private helpers
  -------------------------*/
  private _initialSort(): void {
    this.onSort(new SortEvent('created', 'desc'));
  }
  private _sort(): void {
    // For programmatic sorting on match load using existing column / direction
    this.onSort(new SortEvent(this.sortColumnName, this.sortColumnDirection));
  }
}
