import { Component, OnInit } from '@angular/core';

import { MatchViewModel,
  IMatchViewModel,
  IUserViewModel,
  ICharacterViewModel,
  ITagViewModel,
  IMatchTagViewModel,
  IUserCharacterViewModel } from '../../app.view-models';
import { CommonUxService } from '../common-ux/common-ux.service';
import { MatchManagementService } from './match-management.service';
import { UserManagementService } from '../user-management/user-management.service';
import { CharacterManagementService } from '../character-management/character-management.service';
import { TagManagementService } from '../tag-management/tag-management.service';


@Component({
  selector: 'match-input-form',
  templateUrl: './match-input-form.component.html',
})
export class MatchInputFormComponent implements OnInit {
  public match: IMatchViewModel = new MatchViewModel();
  public characters: ICharacterViewModel[] = [];
  public tags: ITagViewModel[] = [];
  public userCharacter: IUserCharacterViewModel = {} as IUserCharacterViewModel;

  public matchTags: ITagViewModel[] = []; // Will add to match on save

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
          this.userCharacter.userId = this.userCharacter.userId;
          if (res.defaultCharacterId) {
            this.match.userCharacterId = res.defaultCharacterId;
            this.userCharacter.characterId = res.defaultCharacterId;
          }
          if (res.defaultUserCharacterGsp) {
            this.match.userCharacterGsp = res.defaultUserCharacterGsp;
            this.userCharacter.characterGsp = res.defaultUserCharacterGsp
          }
          if (res.defaultUserCharacterId) {
            this.userCharacter.userCharacterId = res.defaultUserCharacterId;
            const defaultUserChar: IUserCharacterViewModel = res.userCharacters.find((userChar) => {
              return userChar.userCharacterId === res.defaultUserCharacterId; 
            });
            this.userCharacter.altCostume = defaultUserChar.altCostume;
          }
          if (res.defaultCharacterName) {
            this.userCharacter.characterName = res.defaultCharacterName;
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

    // Format match tags
    this.match.matchTags = this.matchTags.map(t => {
      return {
        matchTagId: null,
        matchId: null,
        tagId: t.tagId,
        tagName: t.tagName
      } as IMatchTagViewModel;
    });

    this.matchService.createMatch(this.match).subscribe(
      res => {
        this.userCharacter.characterGsp = res.userCharacterGsp;
        this.userService.updateUserCharacter(this.userCharacter).subscribe();
      },
      error => {
        this.commonUxService.showDangerToast('Unable to save match.');
        console.error(error);
      },
      () => {
        this._resetMatch();
        // Set footer warnings to false so it won't show up until the next mouseenter
        this.showFooterWarnings = false;
        this.isSaving = false;
      }
     );
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
    this.matchTags = [];
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
      matchTags: [],
      userWin: null
    } as IMatchViewModel;
  }
}
