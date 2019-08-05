import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { faBars } from '@fortawesome/free-solid-svg-icons';

import { UserManagementService } from '../../modules/user-management/user-management.service';
import { CommonUxService } from '../../modules/common-ux/common-ux.service';
import { IUserViewModel, ILogInViewModel, IServerResponse, LogInViewModel } from 'client/app/app.view-models';

@Component({
  selector: 'top-nav-bar',
  templateUrl: './top-nav-bar.component.html',
  styleUrls: ['./top-nav-bar.component.css']
})
export class TopNavBarComponent implements OnInit {
    public user: IUserViewModel;
    public newUser: IUserViewModel = {} as IUserViewModel;
    public logInModel: ILogInViewModel = {} as ILogInViewModel;

    public showLoginForm: boolean = true;
    public showRegistrationFormWarnings: boolean = false;
    public invalidEmailPassword: boolean = false;
    public paneVisible: boolean = false;

    public faBars = faBars;

    constructor(
      private commonUxService: CommonUxService,
      private userService: UserManagementService,
      private router: Router,
    ) {
    }

    ngOnInit() {
      // Subscribe to whatever user will end up logging in at some point
      this.userService.cachedUser.subscribe({
        next: res => {
          this.user = res;
        },
        error: err => {
          this.commonUxService.showDangerToast('Unable to get user data.');
          console.error(err);
        },
        complete: () => {
        }
      });
    }

    public logIn(): void {
      this.userService.logIn(this.logInModel).subscribe((res: IServerResponse) => {
        if (res.success) {
          this.resetPane();
          this.router.navigate(['/matches']);
        } else {
          this.invalidEmailPassword = true;
        }
      }, error => {
        this.invalidEmailPassword = true;
        console.error(error);
      });
    }
    public logOut(): void {
      this.userService.logOut();
      this.paneVisible = false;
    }
    public createUser(): void {
      if (!this._validateNewUser()) {
        this.commonUxService.showWarningToast('Please address highlighted errors.');
        return;
      }

      this.userService.createUser(this.newUser).subscribe(
        (res: IServerResponse) => {
            if (res.success) {
              this.commonUxService.showSuccessToast('Congratulations! Your account has been created.');
              this.resetPane();
            } else {
              this.commonUxService.showDangerToast('Unable to create account.');
              console.error(res.error);
            }
          },
          error => {
            this.commonUxService.showDangerToast('Unable to create account.');
            console.error(error);
          }
        );
    }
    public togglePanelState(stateToSet: boolean = null): void {
      if (stateToSet !== null) {
        this.paneVisible = stateToSet;
      } else {
        this.paneVisible = !this.paneVisible;
      }

      // When the user closes the login panel,
      // change the form back to the login (not register) form and clear all warnings
      this.resetPane(false);
    }
    public resetPane(closePane: boolean = true): void {
      this.logInModel = {} as LogInViewModel;
      this.newUser = {} as IUserViewModel;

      this.showLoginForm = true;
      this.showRegistrationFormWarnings = false;
      this.invalidEmailPassword = false;

      if (closePane) {
        this.paneVisible = false;
      }
    }
    private _validateNewUser(): boolean {
      let isValid: boolean = true;

      if (!this.newUser.emailAddress) {
        isValid = false;
      }
      if (!this.newUser.userName) {
        isValid = false;
      }
      if (!this.newUser.password) {
        isValid = false;
      }
      if (!this.newUser.passwordConfirm) {
        isValid = false;
      }
      if (this.newUser.password !== this.newUser.passwordConfirm) {
        isValid = false;
      }
      this.showRegistrationFormWarnings = true;
      return isValid;
    }
}
