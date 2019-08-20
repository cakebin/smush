import { Component, Input, HostBinding, OnInit, TemplateRef } from '@angular/core';
import { IMatchViewModel, ICharacterViewModel, ITagViewModel, IMatchTagViewModel } from '../../../app.view-models';
import { faCheck, faTimes, faTrash, faPencilAlt } from '@fortawesome/free-solid-svg-icons';
import { MatchManagementService } from '../match-management.service';
import { CommonUxService } from '../../common-ux/common-ux.service';
import { NgbPopover } from '@ng-bootstrap/ng-bootstrap';

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

  public editedMatch: IMatchViewModel = {} as IMatchViewModel;
  public editedMatchTags: ITagViewModel[] = []; // Will add to match on save
  public warnings: string[] = [];

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
      matchTags: [],
      userWin: originalMatch.userWin === undefined ? null : originalMatch.userWin,
      created: originalMatch.created
    } as IMatchViewModel;

    // Tags need to be copied over so we don't send a reference to the original tags
    Object.assign(this.editedMatchTags, originalMatch.matchTags);
  }
  public deleteMatch(match: IMatchViewModel): void {
    this.commonUxService.openConfirmModal(
      'Removing match against ' + match.opponentCharacterName + '.',
      'Delete match',
      false,
      'Nuke it').then(
      confirm => {
        this.matchService.deleteMatch(match);
      },
      reject => {
        // Do nothing
      }
    );
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

  public updateSelectedTags(event: ITagViewModel[]): void {
    if (event != null) {
      this.editedMatchTags = event;
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
  public validateMatch(): boolean {
    this.warnings = [];
    if (!this.editedMatch.opponentCharacterId) {
      this.warnings.push('Opponent character required.');
    }
    if (!this.editedMatch.userCharacterId && this.editedMatch.userCharacterGsp) {
      this.warnings.push('User GSP must be associated with a user character.');
    }
    if (this.warnings.length) {
      return false;
    } else {
      return true;
    }
  }
  public openWarningPopover(popover: NgbPopover) {
    if (!this.validateMatch()) {
      popover.popoverTitle = 'Invalid match';
      popover.ngbPopover = this.warnings.join(' ');
      popover.open();
    }
  }
}
