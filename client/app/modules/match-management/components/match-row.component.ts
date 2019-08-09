import { Component, Input, EventEmitter, Output } from '@angular/core';
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
  @Input() characters: ICharacterViewModel[] = [];

  public editedMatch: IMatchViewModel = {} as IMatchViewModel;

  public faCheck = faCheck;
  public faTimes = faTimes;
  public faTrash = faTrash;
  public faPencilAlt = faPencilAlt;

  constructor(
    private matchService: MatchManagementService,
    private commonUxService: CommonUxService,
  ) {
  }

  public editMatch(originalMatch: IMatchViewModel): void {
    originalMatch.editMode = true;
    Object.assign(this.editedMatch, originalMatch);
  }
  public saveChanges(): void {
    this.matchService.updateMatch(this.editedMatch).subscribe(
      (res: IMatchViewModel) => {
        if (res) {
          // Object.assign(this.match, res);
          this.resetState();
          this.commonUxService.showSuccessToast('Match updated.');
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
    }
  }
  public onSelectEditUserCharacter(event: ICharacterViewModel): void {
    if (event) {
      this.editedMatch.userCharacterId = event.characterId;
    }
  }
}
