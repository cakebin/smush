import { Component, OnInit, HostListener } from '@angular/core';
import { CommonUxService } from '../../modules/common-ux/common-ux.service';
import { UserManagementService } from '../../modules/user-management/user-management.service';
import { CharacterManagementService } from 'client/app/modules/character-management/character-management.service';
import { IUserViewModel, ICharacterViewModel } from '../../app.view-models';
import { faQuestionCircle } from '@fortawesome/free-solid-svg-icons';


class ISavedCharacter {
  id: number;
  name: string;
  gsp: number;
  isDefault: boolean;
  editMode: boolean;
}
class SavedCharacter implements ISavedCharacter {
  constructor(
    public id: number = null,
    public name: string = '',
    public gsp: number = null,
    public isDefault: boolean = false,
    public editMode: boolean = false,
  ) {}
}

@Component({
  selector: 'profile-edit',
  templateUrl: './profile-edit.component.html'
})
export class ProfileEditComponent implements OnInit {
  public savedCharactersTestData: ISavedCharacter[] = [];

  public editChar: ISavedCharacter = {} as ISavedCharacter;

  public characters: ICharacterViewModel[] = [];
  public user: IUserViewModel = {} as IUserViewModel;
  public editedUser: IUserViewModel = {} as IUserViewModel;

  public set defaultCharacterGspString(value: string) {
    this._defaultCharacterGspString = value;
    this.editedUser.defaultCharacterGsp = parseInt(value.replace(/\D/g, ''), 10);
  }
  public get defaultCharacterGspString(): string {
    return this._defaultCharacterGspString;
  }
  private _defaultCharacterGspString: string = '';

  public showFooterWarnings = false;
  public warnings: string[] = [];
  public isSaving = false;
  public faQuestionCircle = faQuestionCircle;
  public formChanged: boolean = false;

  constructor(
    private commonUxService: CommonUxService,
    private userService: UserManagementService,
    private characterService: CharacterManagementService,
    ) {
  }

  @HostListener('keyup', ['$event'])
  onKeyUp() {
    this.formChanged = this.getChangedStatus();
  }

  ngOnInit() {
    // Subscribe to the user data (could change from other components on the page)
    this.userService.cachedUser.subscribe({
      next: res => {
        if (res) {
          Object.assign(this.user, res);
          Object.assign(this.editedUser, res);

          if (this.editedUser.defaultCharacterGsp) {
            this.defaultCharacterGspString = this.editedUser.defaultCharacterGsp.toString();
          }
        }
      },
      error: err => {
        this.commonUxService.showDangerToast('Unable to get user data.');
        console.error(err);
      }
    });
    this.characterService.characters.subscribe(
      res => {
        if (res) {
          this.characters = res;

          // TEST DATA BEFORE WE HAVE A SAVED CHAR TABLE
          this.savedCharactersTestData = res.slice(0, 5).map(char => {
            return new SavedCharacter(char.characterId, char.characterName, 98000);
          });
          this.savedCharactersTestData[3].isDefault = true;
        }
      }
    );
  }

  public setDefaultSavedCharacter(defaultCharId: number) {
    this.savedCharactersTestData.forEach(char => {
      if (char.id === defaultCharId) {
        char.isDefault = true;
      } else {
        char.isDefault = false;
      }
    });
  }
  public setSavedCharacterEditMode(editCharId: number) {
    this.savedCharactersTestData.forEach(char => {
      if (char.id === editCharId) {
        char.editMode = true;
      } else {
        char.editMode = false;
      }
    });
  }
  public onSelectSavedCharacter(event: ICharacterViewModel, savedChar: ISavedCharacter) {
    savedChar.id = event.characterId;
    savedChar.name = event.characterName;
  }
  public onSelectDefaultCharacter(event: ICharacterViewModel): void {
      this.editedUser.defaultCharacterId = event.characterId;
      this.formChanged = this.getChangedStatus();
  }
  public updateUser(): void {
    this.userService.updateUser(this.editedUser).subscribe(
      res => {
        // Copy changes from edited user to the actual user object
        Object.assign(this.user, this.editedUser);
        this.formChanged = this.getChangedStatus();
        this.commonUxService.showSuccessToast('User information updated!');
      },
      error => {
        this.commonUxService.showDangerToast('Unable to update user information.');
        console.error(error);
      });
  }

  public getChangedStatus(): boolean {
    const keys: string[] = Object.keys(this.editedUser);
    let formChanged: boolean = false;
    keys.forEach(k => {
      if (!Object.is(this.user[k], this.editedUser[k])) {
        formChanged = true;
      }
    });
    return formChanged;
  }
}
