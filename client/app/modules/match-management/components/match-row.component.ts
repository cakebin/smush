import { Component, Input, HostBinding, OnInit } from '@angular/core';
import { IMatchViewModel, ICharacterViewModel, ITagViewModel, IMatchTagViewModel } from '../../../app.view-models';
import { faCheck, faTimes, faTrash, faPencilAlt } from '@fortawesome/free-solid-svg-icons';
import { MatchManagementService } from '../match-management.service';
import { CommonUxService } from '../../common-ux/common-ux.service';

@Component({
  selector: '[match-row]',
  templateUrl: './match-row.component.html'
})
export class MatchRowComponent implements OnInit {
  @Input() match: IMatchViewModel = {} as IMatchViewModel;
  @Input() isUserOwned: boolean = false;
  @Input() characters: ICharacterViewModel[] = [];
  @Input() tags: ITagViewModel[] = [];
  @HostBinding('style.height') trHeight: string = '40px';

  public editedMatchTags: ITagViewModel[] = []; // Will add to match on save
  public newTag: ITagViewModel = null;

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

  ngOnInit() {
    this.match.editMode = false;
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
      matchTags: null,
      userWin: originalMatch.userWin === undefined ? null : originalMatch.userWin,
      created: originalMatch.created
    } as IMatchViewModel;

    Object.assign(this.editedMatchTags, originalMatch.matchTags);
  }
  public saveChanges(): void {
    if (!this.editedMatch.opponentCharacterId) {
      this.commonUxService.showWarningToast('Opponent character required.');
      return;
    }

    this.editedMatch.matchTags = this.editedMatchTags.map(t => {
      return {
        matchTagId: null,
        matchId: this.editedMatch.matchId,
        tagId: t.tagId,
        tagName: t.tagName
      } as IMatchTagViewModel;
    });

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
  public removeTag(tag: ITagViewModel): void {
    const tagIndex: number = this.editedMatchTags.findIndex(t => t.tagId === tag.tagId);
    this.editedMatchTags.splice(tagIndex, 1);
  }
  public onSelectTag(event: ITagViewModel): void {
    if (event != null) {
      if (!this.editedMatchTags.find(t => t.tagId === event.tagId)) {
        this.editedMatchTags.push(event);
      }
    }
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
