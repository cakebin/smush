import { Component, Input, EventEmitter, Output, HostBinding } from '@angular/core';
import { IMatchViewModel, ICharacterViewModel } from '../../../app.view-models';
import { faCheck, faTimes, faTrash, faPencilAlt } from '@fortawesome/free-solid-svg-icons';
import { MatchManagementService } from '../match-management.service';
import { CommonUxService } from '../../common-ux/common-ux.service';

@Component({
  selector: '[match-row]',
  templateUrl: './match-row.component.html'
})
export class MatchRowComponent {
  @Input() match: IMatchViewModel = {} as IMatchViewModel;
  @Input() isUserOwned: boolean = false;
  @Input() characters: ICharacterViewModel[] = [];
  @HostBinding('style.height') trHeight: string = '40px';

  public editedMatch: IMatchViewModel = {} as IMatchViewModel;

  public faCheck = faCheck;
  public faTimes = faTimes;
  public faTrash = faTrash;
  public faPencilAlt = faPencilAlt;
  public boolOptions: any[] = [
    { name: 'Yes', value: true },
    { name: 'No', value: false },
  ];

  constructor(
    private matchService: MatchManagementService,
    private commonUxService: CommonUxService,
  ) {
  }

  public editMatch(originalMatch: IMatchViewModel): void {
    originalMatch.editMode = true;

    // Properties don't exist on editedMatch if they aren't filled in,
    // so we need to make sure we have all relevant fields, null or not
    this.editedMatch = {
      matchId: originalMatch.matchId,
      userId: originalMatch.userId,
      userName: originalMatch.userName,
      userCharacterId: originalMatch.userCharacterId,
      userCharacterGsp: originalMatch.userCharacterGsp,
      opponentCharacterId: originalMatch.opponentCharacterId,
      opponentCharacterGsp: originalMatch.opponentCharacterGsp,
      opponentAwesome: originalMatch.opponentAwesome === undefined ? null : originalMatch.opponentAwesome,
      opponentTeabag: originalMatch.opponentTeabag === undefined ? null : originalMatch.opponentTeabag,
      opponentCamp: originalMatch.opponentCamp === undefined ? null : originalMatch.opponentCamp,
      userWin: originalMatch.userWin === undefined ? null : originalMatch.userWin,
      created: originalMatch.created
    } as IMatchViewModel;
  }
  public saveChanges(): void {
    if (!this.editedMatch.opponentCharacterId) {
      this.commonUxService.showWarningToast('Opponent character required.');
      return;
    }
    this.matchService.updateMatch(this.editedMatch).subscribe(
      (res: IMatchViewModel) => {
        if (res) {
          this.match = res;
          this.editedMatch = res;
          this.resetState();
        }
      });
  }
  public resetState(): void {
    this.editedMatch = {} as IMatchViewModel;
    this.match.editMode = false;
  }
  public onSelectEditOpponentCharacter(event: ICharacterViewModel): void {
    if (event) {
      this.editedMatch.opponentCharacterId = event.characterId;
    } else {
      this.editedMatch.opponentCharacterId = null;
    }
  }
  public onSelectEditUserCharacter(event: ICharacterViewModel): void {
    if (event) {
      this.editedMatch.userCharacterId = event.characterId;
    } else {
      this.editedMatch.userCharacterId = null;
    }
  }
}
