import { Component, OnInit, Input, HostBinding } from '@angular/core';
import { CommonUxService } from '../../common-ux/common-ux.service';
import { UserManagementService } from '../user-management.service';
import { ICharacterViewModel, IUserCharacterViewModel, IUserViewModel } from '../../../app.view-models';
import { MatchManagementService } from '../../match-management/match-management.service';

@Component({
  selector: 'tr[user-character-row]',
  templateUrl: './user-character-row.component.html'
})
export class UserCharacterRowComponent implements OnInit {
  @Input() characters: ICharacterViewModel[] = [];
  @Input() userCharacter: IUserCharacterViewModel = {} as IUserCharacterViewModel;
  @Input() set user(value: IUserViewModel) {
    this._user = value;
    if (!value) {
      return;
    }
    if (value.defaultUserCharacterId === this.userCharacter.userCharacterId) {
      this.isDefaultCharacter = true;
      this.rowClass = 'bg-primary text-white';
    } else {
      this.isDefaultCharacter = false;
      this.rowClass = '';
    }
  }
  get user(): IUserViewModel {
    return this._user;
  }
  private _user: IUserViewModel = {} as IUserViewModel;


  @HostBinding('class') rowClass: string = '';
  public isDefaultCharacter: boolean = false;
  public isEditMode: boolean = false;
  public editedUserCharacter: IUserCharacterViewModel = {} as IUserCharacterViewModel;

  constructor(
    private commonUxService: CommonUxService,
    private userService: UserManagementService,
    private matchService: MatchManagementService,
  ) { }

  ngOnInit() {
  }

  // Template-related methods
  public enterEditMode() {
    Object.assign(this.editedUserCharacter, this.userCharacter);
    this.isEditMode = true;
  }
  public leaveEditMode() {
    this.isEditMode = false;
    this.editedUserCharacter = {} as IUserCharacterViewModel;
  }
  public onSelectCharacter(event: ICharacterViewModel) {
    if (event) {
      this.editedUserCharacter.characterId = event.characterId;
      this.editedUserCharacter.characterName = event.characterName;
    } else {
      this.editedUserCharacter.characterId = null;
      this.editedUserCharacter.characterName = '';
    }
  }


  // Api-related methods
  public updateUserCharacter() {
    this.userService.updateUserCharacter(this.editedUserCharacter).subscribe(
      (res: IUserViewModel) => {
        this.leaveEditMode();
        this.matchService.updateCachedMatchesWithAltCostume(res);
      }
    );
  }
  public deleteUserCharacter() {
    this.userService.deleteUserCharacter(this.userCharacter);
  }
  public setDefaultUserCharacter() {
    this.userService.setDefaultUserCharacter(this.userCharacter);
  }
  public unsetDefaultUserCharacter() {
    this.userService.unsetDefaultUserCharacter(this.userCharacter);
  }
}
