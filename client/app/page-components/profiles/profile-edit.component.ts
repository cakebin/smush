import { Component, OnInit, HostListener } from '@angular/core';
import { CommonUxService } from '../../modules/common-ux/common-ux.service';
import { UserManagementService } from '../../modules/user-management/user-management.service';
import { MatchManagementService } from 'client/app/modules/match-management/match-management.service';
import { IUserViewModel } from '../../app.view-models';
import { faQuestionCircle } from '@fortawesome/free-solid-svg-icons';


class ISavedCharacter {
  name: string;
  gsp: number;
}
class SavedCharacter implements ISavedCharacter {
  constructor(
    public name: string = '',
    public gsp: number = null
  ) {}
}

@Component({
  selector: 'profile-edit',
  templateUrl: './profile-edit.component.html',
  styleUrls: ['./profile-edit.component.css']
})
export class ProfileEditComponent implements OnInit {
  public savedCharactersTestData: SavedCharacter[] = [
    new SavedCharacter('Ike', 4300000),
    new SavedCharacter('Palutena', 4200000),
    new SavedCharacter('Lucina', 4500000),
    new SavedCharacter('Marth', 5100000),
    new SavedCharacter('Joker', 5200000),
  ];
  public editChar: SavedCharacter = {} as SavedCharacter;

  public characters = this.matchService.characters;
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
    private matchService: MatchManagementService,
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
  }
  public onSelectDefaultCharacter(event: string): void {
      this.editedUser.defaultCharacterName = event;
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
