import { Component, Input, OnInit } from '@angular/core';
import { IMatchViewModel, ICharacterViewModel, IUserViewModel, IUserCharacterViewModel } from '../../../app.view-models';
import { MatchManagementService } from '../match-management.service';
import { CommonUxService } from '../../common-ux/common-ux.service';

@Component({
  selector: 'match-card',
  templateUrl: './match-card.component.html',
  styleUrls: ['./match-card.component.css']
})
export class MatchCardComponent implements OnInit{
  @Input() characters: ICharacterViewModel[] = [];
  @Input() match: IMatchViewModel = {} as IMatchViewModel;
  @Input() set user(user: IUserViewModel) {
    // Set user
    this._user = user;
    if (!user) {
      return;
    }
    // Calculate ownership of match
    this.isUserOwned = this.match.userId === this.user.userId;
  }
  get user(): IUserViewModel {
    return this._user;
  }
  private _user: IUserViewModel = {} as IUserViewModel;

  // Calculated display vars
  public userCharacterImagePath: string = '';
  public isUserOwned: boolean = false;

  // Form vars
  public editedMatch: IMatchViewModel = {} as IMatchViewModel;
  public boolOptions: any[] = [
    { name: 'Yes', value: true },
    { name: 'No', value: false },
  ];

  constructor(
    private matchService: MatchManagementService,
    private commonUxService: CommonUxService,
  ) {
  }

  ngOnInit() {
    if (this.match.altCostume) {
      this.userCharacterImagePath = '/static/assets/alt/' + this.match.userCharacterImage.replace('.png', '') +
      '_' + this.match.altCostume + '.png';
    } else {
      this.userCharacterImagePath = '/static/assets/full/' + this.match.userCharacterImage;
    }
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
