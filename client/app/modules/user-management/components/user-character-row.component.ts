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

  public enterEditMode() {
    Object.assign(this.editedUserCharacter, this.userCharacter);
    this.isEditMode = true;
  }
  public leaveEditMode() {
    this.isEditMode = false;
    this.editedUserCharacter = {} as IUserCharacterViewModel;
  }
  public saveUserCharacter() {
    // Make api call
    console.log('making api call to save default char');
    Object.assign(this.userCharacter, this.editedUserCharacter);
    this.isEditMode = false;
  }
  public deleteUserCharacter() {
    // Make api call
    console.log('making api call to DELETE default char');
  }
  public setDefaultUserCharacter() {
    // Make api call
    console.log('making api call to set default char');
    this.isDefaultCharacter = true;
  }
  public unsetDefaultUserCharacter() {
    // Make api call
    console.log('making api call to UNSET default char');
    this.isDefaultCharacter = false;
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
}
