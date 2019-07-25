import { Component, OnInit, ViewChild } from '@angular/core';
import { DecimalPipe } from '@angular/common';
import { IUserViewModel, UserViewModel } from '../../app.view-models';
import { NumberMaskDirective } from '../../modules/common-ux/directives/number-mask.directive';
import { TypeaheadComponent } from '../../modules/common-ux/components/typeahead/typeahead.component';
import { CommonUXService } from '../../modules/common-ux/common-ux.service';
import { UserManagementService } from '../../modules/user-management/user-management.service';

@Component({
  selector: 'profile-edit',
  templateUrl: './profile-edit.component.html',
})
export class ProfileEditComponent implements OnInit {
  @ViewChild('defaultCharacterNameInput', { static: false }) private defaultCharacterNameInput: TypeaheadComponent;
  @ViewChild('defaultCharacterGspInput', { static: false }) private defaultCharacterGspInput: NumberMaskDirective;

  public user: IUserViewModel = new UserViewModel();

  public showFooterWarnings = false;
  public warnings: string[] = [];
  public isSaving = false;

  constructor(
    private commonUXService: CommonUXService,
    private userService: UserManagementService,
    ) {
  }

  ngOnInit() {
    this.userService.cachedUser.subscribe({
      next: res => {
        if (res) {
          this.user = res;
          if (this.user.defaultCharacterName) {
            // This doesn't actually work yet. Need to write out this method
            // this.opponentCharacterNameInput.select(this.match.userCharacterName);
          }
          if (this.user.defaultCharacterGsp) {
            this.defaultCharacterGspInput.setValue(res.defaultCharacterGsp);
          }
      }
      },
      error: err => {
        this.commonUXService.showDangerToast('Unable to get user data.');
        console.error(err);
      },
      complete: () => {
      }
    });
  }
}