import { Component, OnInit, ViewChild, AfterViewInit } from '@angular/core';
import { IUserViewModel, UserViewModel } from '../../app.view-models';
import { MaskedNumberInputComponent } from '../../modules/common-ux/components/masked-number-input/masked-number-input.component';
import { TypeaheadComponent } from '../../modules/common-ux/components/typeahead/typeahead.component';
import { CommonUXService } from '../../modules/common-ux/common-ux.service';
import { UserManagementService } from '../../modules/user-management/user-management.service';
import { faQuestionCircle } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'profile-edit',
  templateUrl: './profile-edit.component.html',
})
export class ProfileEditComponent implements OnInit, AfterViewInit {
  @ViewChild('defaultCharacterNameInput', { static: false }) private defaultCharacterNameInput: TypeaheadComponent;
  @ViewChild('defaultCharacterGspInput', { static: false }) private defaultCharacterGspInput: MaskedNumberInputComponent;

  public user: IUserViewModel = new UserViewModel();
  public defaultCharacterGspString: string = '';

  public showFooterWarnings = false;
  public warnings: string[] = [];
  public isSaving = false;

  public faQuestionCircle = faQuestionCircle;

  constructor(
    private commonUXService: CommonUXService,
    private userService: UserManagementService,
    ) {
  }

  ngOnInit() {
    // Subscribe to the user data (could change from other components on the page)
    this.userService.cachedUser.subscribe({
      next: res => {
        if (res) {
          this.user = res;
          this.defaultCharacterGspString = this.user.defaultCharacterGsp.toString();
          this._formatValues();
      }
      },
      error: err => {
        this.commonUXService.showDangerToast('Unable to get user data.');
        console.error(err);
      }
    });
  }

  ngAfterViewInit() {
     this._formatValues();
  }

  private _formatValues(): void {
    // OKAY. So Angular has some sometimes-problematic design issues that require nasty hacks to get around.
    setTimeout(() => {
      if (this.defaultCharacterNameInput && this.user.defaultCharacterName) {
        this.defaultCharacterNameInput.setDefaultValue(this.user.defaultCharacterName);
      }
      if (this.defaultCharacterGspInput && this.user.defaultCharacterGsp) {
        this.defaultCharacterGspInput.setValue(this.user.defaultCharacterGsp);
      }
    }, 0);
  }
}
