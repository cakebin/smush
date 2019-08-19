import { Component, Input, OnInit } from '@angular/core';
import { IMatchViewModel, ICharacterViewModel, IUserViewModel, ITagViewModel } from '../../../app.view-models';
import { MatchCardEditComponent } from './match-card-edit.component';
import { MatchManagementService } from '../match-management.service';
import { CommonUxService } from '../../common-ux/common-ux.service';
import { reject } from 'q';

@Component({
  selector: 'match-card',
  templateUrl: './match-card.component.html',
  styleUrls: ['./match-card.component.css']
})
export class MatchCardComponent implements OnInit {
  @Input() tags: ITagViewModel[] = [];
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
      matchTags: [],
      userWin: originalMatch.userWin === undefined ? null : originalMatch.userWin,
      created: originalMatch.created,
    } as IMatchViewModel;

    // Tags need to be copied over so we don't send a reference to the original tags
    Object.assign(this.editedMatch.matchTags, originalMatch.matchTags);

    const modalRef = this.commonUxService.openModal(MatchCardEditComponent);
    modalRef.componentInstance.editedMatch = this.editedMatch;
    modalRef.componentInstance.tags = this.tags;
    modalRef.componentInstance.characters = this.characters;

    modalRef.result.then(
      result => {
        // User saved changes
        if (result) {
          this.match = result;
        }
      });
  }
  public deleteMatch(match: IMatchViewModel): void {
    this.commonUxService.openConfirmModal(
      'Removing match against ' + match.opponentCharacterName + '.',
      'Delete match',
      false,
      'Nuke it').then(
      confirm => {
        this.matchService.deleteMatch(match);
      }, dismiss => {
        // Cancelled. Do nothing
      });
  }
}
