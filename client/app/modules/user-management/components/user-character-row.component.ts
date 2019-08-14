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
    console.log('making api call to update user char');
    this.userService.updateUserCharacter(this.editedUserCharacter).subscribe(
      res => {
        console.log('done updating userChar');
        this.leaveEditMode();
      }
    );
  }
  public deleteUserCharacter() {
    console.log('making api call to DELETE user char');
    this.userService.deleteUserCharacter(this.userCharacter);
  }
  public setDefaultUserCharacter() {
    console.log('making api call to set default char');
    this.userService.setDefaultUserCharacter(this.userCharacter);
  }
  public unsetDefaultUserCharacter() {
    console.log('making api call to UNSET default char');
    this.userService.unsetDefaultUserCharacter(this.userCharacter);
  }
}
