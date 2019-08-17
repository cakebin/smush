import { Component, OnInit, ViewChild } from '@angular/core';

import { MatchViewModel, IMatchViewModel, IUserViewModel, ICharacterViewModel, ITagViewModel } from '../../app.view-models';
import { CommonUxService } from '../common-ux/common-ux.service';
import { MatchManagementService } from './match-management.service';
import { UserManagementService } from '../user-management/user-management.service';
import { CharacterManagementService } from '../character-management/character-management.service';
import { TagManagementService } from '../tag-management/tag-management.service';
import { TypeaheadComponent } from '../common-ux/components/typeahead/typeahead.component';


@Component({
  selector: 'match-input-form',
  templateUrl: './match-input-form.component.html',
})
export class MatchInputFormComponent implements OnInit {
  @ViewChild('tagInput', { static: false }) tagInputComponent: TypeaheadComponent;

  public match: IMatchViewModel = new MatchViewModel();
  public characters: ICharacterViewModel[] = [];
  public tags: ITagViewModel[] = [];

  // Fake array, will add to actual match later
  public matchTags: ITagViewModel[] = [];
  public newTag: ITagViewModel = null;

  public showFooterWarnings: boolean = false;
  public warnings: string[] = [];
  public isSaving: boolean = false;

  private user: IUserViewModel;

  constructor(
    private commonUxService: CommonUxService,
    private userService: UserManagementService,
    private matchService: MatchManagementService,
    private characterService: CharacterManagementService,
    private tagService: TagManagementService,
    ) {
  }

  ngOnInit() {
    this.userService.cachedUser.subscribe(
      res => {
        if (res) {
          this.user = res;
          this.match.userId = this.user.userId;
          if (res.defaultCharacterId) {
            this.match.userCharacterId = res.defaultCharacterId;
          }
          if (res.defaultUserCharacterGsp) {
            this.match.userCharacterGsp = res.defaultUserCharacterGsp;
          }
        }
    });
    this.characterService.cachedCharacters.subscribe(
      res => {
        this.characters = res;
      }
    );
    this.tagService.cachedTags.subscribe(
      res => {
        this.tags = res;
      }
    );
  }

  public createEntry(): void {
    if (!this.validateMatch()) {
      this.warnings.forEach(warningMessage => {
        this.commonUxService.showWarningToast(warningMessage);
      });
      return;
    }
    this.isSaving = true;
    this.matchService.createMatch(this.match).subscribe((res: number) => {
      // On success, do nothing
    }, error => {
      this.commonUxService.showDangerToast('Unable to save match.');
      console.error(error);
    }, () => {
      this._resetMatch();
      // Set footer warnings to false so it won't show up until the next mouseenter
      this.showFooterWarnings = false;
      this.isSaving = false;
    });
  }
  public removeTag(tag: ITagViewModel): void {
    const tagIndex: number = this.matchTags.findIndex(t => t.tagId === tag.tagId);
    this.matchTags.splice(tagIndex, 1);
  }
  public onSelectTag(event: ITagViewModel): void {
    if (event != null) {
      if (!this.matchTags.find(t => t.tagId === event.tagId)) {
        this.matchTags.push(event);
      }
    }
  }
  public onSelectOpponentCharacter(event: ICharacterViewModel): void {
    // Event properties aren't accessible in the template
    if (event == null) {
      this.match.opponentCharacterId = null;
    } else {
      this.match.opponentCharacterId = event.characterId;
    }
  }
  public onSelectUserCharacter(event: ICharacterViewModel): void {
    // Event properties aren't accessible in the template
    if (event == null) {
      this.match.userCharacterId = null;
    } else {
      this.match.userCharacterId = event.characterId;
    }
  }
  public validateMatch(): boolean {
    this.warnings = [];
    if (!this.match.opponentCharacterId) {
      this.warnings.push('Opponent character required.');
    }
    if (!this.match.userCharacterId && this.match.userCharacterGsp) {
      this.warnings.push('User GSP must be associated with a user character.');
    }
    if (this.warnings.length) {
      return false;
    } else {
      return true;
    }
  }

  /*-----------------------
       Private helpers
  ------------------------*/
  private _resetMatch(): void {
    this.match = {
      matchId: null,
      userId: this.user.userId,
      userName: null,
      userCharacterId: this.match.userCharacterId,
      userCharacterName: this.match.userCharacterName,
      userCharacterGsp: this.match.userCharacterGsp,
      opponentCharacterId: null,
      opponentCharacterName: null,
      opponentCharacterGsp: null,
      opponentAwesome: null,
      opponentTeabag: null,
      opponentCamp: null,
      userWin: null
    } as IMatchViewModel;
  }
}
