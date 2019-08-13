import { Component, OnInit, HostListener } from '@angular/core';
import { CommonUxService } from '../../common-ux/common-ux.service';
import { UserManagementService } from '../user-management.service';
import { CharacterManagementService } from 'client/app/modules/character-management/character-management.service';
import { IUserViewModel, IUserCharacterViewModel, ICharacterViewModel } from '../../../app.view-models';
import { faQuestionCircle } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'profile-edit',
  templateUrl: './profile-edit.component.html'
})
export class ProfileEditComponent implements OnInit {
  public newUserCharacter: IUserCharacterViewModel = {} as IUserCharacterViewModel;

  public characters: ICharacterViewModel[] = [];
  public user: IUserViewModel = {} as IUserViewModel;
  public editedUser: IUserViewModel = {} as IUserViewModel;

  public showFooterWarnings = false;
  public warnings: string[] = [];
  public isSaving = false;
  public formChanged: boolean = false;
  public faQuestionCircle = faQuestionCircle;

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
        }
      }
    );
  }

  public updateUser(): void {
    this.userService.updateUser(this.editedUser).subscribe(
      res => {
        // Copy changes from edited user to the actual user object
        Object.assign(this.user, this.editedUser);
        this.formChanged = this.getChangedStatus();
        this.commonUxService.showStandardToast('User information updated!');
      },
      error => {
        this.commonUxService.showDangerToast('Unable to update user information.');
        console.error(error);
      });
  }
  public addUserCharacter(): void {
    if (!this.newUserCharacter.characterId) {
      this.commonUxService.showWarningToast('Please select a character before adding a user character.');
      return;
    }
    if (this.user.userCharacters.find(c => c.characterId === this.newUserCharacter.characterId)) {
      this.commonUxService.showWarningToast('This user character already exists.');
      return;
    }

    console.log('api call to add user character (fake adding for now)');
    this.user.userCharacters.push(this.newUserCharacter);
    // Clear "add character" inputs
    this.newUserCharacter = {
      characterId: null,
      characterGsp: null
    } as IUserCharacterViewModel;
  }
  public onSelectNewUserCharacter(event: ICharacterViewModel) {
    if (event) {
      this.newUserCharacter.characterId = event.characterId;
      this.newUserCharacter.characterName = event.characterName;
    } else {
      this.newUserCharacter.characterId = null;
      this.newUserCharacter.characterName = '';
    }
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
