import { Component, OnInit, Input } from '@angular/core';
import { NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { CommonUxService } from '../../common-ux/common-ux.service';
import { IMatchViewModel, ITagViewModel, IMatchTagViewModel, ICharacterViewModel } from 'client/app/app.view-models';
import { MatchManagementService } from '../match-management.service';
import { faQuestionCircle } from '@fortawesome/free-solid-svg-icons';


@Component({
  selector: 'match-card-edit-component',
  templateUrl: './match-card-edit.component.html',
  // Styles don't seem to apply to dynamically generated components,
  // so I have to manually include it here
  styleUrls: ['../../common-ux/common-ux.css']
})
export class MatchCardEditComponent implements OnInit {
  @Input() tags: ITagViewModel[] = [];
  @Input() characters: ICharacterViewModel[] = [];
  @Input() editedMatch: IMatchViewModel = {} as IMatchViewModel;

  public editedMatchTags: ITagViewModel[] = []; // Will add to match on save
  public warnings: string[] = [];
  public isSaving: boolean = false;
  public faQuestionCircle = faQuestionCircle;

  constructor(
    private activeModalRef: NgbActiveModal,
    private commonUxService: CommonUxService,
    private matchService: MatchManagementService,
  ) {
  }

  ngOnInit() {
    this.editedMatchTags = this.editedMatch.matchTags;
  }


  /*----------------------
     Typeahead handlers
  ----------------------*/

  public onSelectOpponentCharacter(event: ICharacterViewModel): void {
    // Event properties aren't accessible in the template
    if (event == null) {
      this.editedMatch.opponentCharacterId = null;
    } else {
      this.editedMatch.opponentCharacterId = event.characterId;
    }
  }
  public onSelectUserCharacter(event: ICharacterViewModel): void {
    // Event properties aren't accessible in the template
    if (event == null) {
      this.editedMatch.userCharacterId = null;
    } else {
      this.editedMatch.userCharacterId = event.characterId;
    }
  }


  /*----------------------
        Save changes
  ----------------------*/

  public saveChanges(): void {
    if (!this.validateMatch()) {
      this.warnings.forEach(warningMessage => {
        this.commonUxService.showWarningToast(warningMessage);
      });
      return;
    }
    this.isSaving = true;
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
          this.activeModalRef.close(res);
        }
      }, error => {
        this.commonUxService.showDangerToast('The match could not be updated.');
        console.error(error);
      }, () => {
        this.isSaving = false;
      });
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
  public close() {
    this.activeModalRef.close();
  }
}
