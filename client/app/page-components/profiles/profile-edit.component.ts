import { Component, OnInit } from '@angular/core';
import { CommonUXService } from '../../modules/common-ux/common-ux.service';
import { UserManagementService } from '../../modules/user-management/user-management.service';
import { MatchManagementService } from 'client/app/modules/match-management/match-management.service';
import { IUserViewModel } from '../../app.view-models';
import { faQuestionCircle } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'profile-edit',
  templateUrl: './profile-edit.component.html',
})
export class ProfileEditComponent implements OnInit {
  public characters = this.matchService.characters;
  public user: IUserViewModel = {} as IUserViewModel;

  public set defaultCharacterGspString(value: string) {
    this._defaultCharacterGspString = value;
    this.user.defaultCharacterGsp = parseInt(value.replace(/\D/g, ''), 10);
  }
  public get defaultCharacterGspString(): string {
    return this._defaultCharacterGspString;
  }
  private _defaultCharacterGspString: string = '';

  public showFooterWarnings = false;
  public warnings: string[] = [];
  public isSaving = false;
  public faQuestionCircle = faQuestionCircle;

  constructor(
    private commonUxService: CommonUXService,
    private userService: UserManagementService,
    private matchService: MatchManagementService,
    ) {
  }

  ngOnInit() {
    // Subscribe to the user data (could change from other components on the page)
    this.userService.cachedUser.subscribe({
      next: res => {
        if (res) {
          this.user = res;
          this.defaultCharacterGspString = this.user.defaultCharacterGsp.toString();
        }
      },
      error: err => {
        this.commonUxService.showDangerToast('Unable to get user data.');
        console.error(err);
      }
    });
  }

  public updateUser(): void {
    this.userService.updateUser(this.user).subscribe(
      res => {
        this.commonUxService.showSuccessToast('User information updated!');
      },
      error => {
        this.commonUxService.showDangerToast('Error updating user information.');
        console.error(error);
      });
  }
}
