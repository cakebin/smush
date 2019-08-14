import { Component, OnInit, Input, HostBinding } from '@angular/core';
import { CommonUxService } from '../../common-ux/common-ux.service';
import { UserManagementService } from '../user-management.service';
import { ICharacterViewModel, IUserCharacterViewModel } from '../../../app.view-models';

@Component({
  selector: 'tr[user-character-row]',
  templateUrl: './user-character-row.component.html'
})
export class UserCharacterRowComponent implements OnInit {
  @Input() userCharacter: IUserCharacterViewModel = {} as IUserCharacterViewModel;
  @Input() characters: ICharacterViewModel[] = [];
  @Input() set isDefaultCharacter(value: boolean) {
    this._isDefaultCharacter = value;
    if (value) {
      this.rowClass = 'bg-primary text-white';
    } else {
      this.rowClass = '';
    }
  }
  get isDefaultCharacter(): boolean {
    return this._isDefaultCharacter;
  }
  private _isDefaultCharacter: boolean = false;

  @HostBinding('class') rowClass: string = '';
  public isEditMode: boolean = false;
  public editedUserCharacter: IUserCharacterViewModel = {} as IUserCharacterViewModel;

  constructor(
    private commonUxService: CommonUxService,
    private userService: UserManagementService,
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
      res => {
        this.leaveEditMode();
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
