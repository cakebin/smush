import { Component, OnInit, HostListener} from '@angular/core';
import { UserManagementService } from '../../modules/user-management/user-management.service';
import { CommonUXService } from '../../modules/common-ux/common-ux.service';
import { IUserViewModel, ILogInViewModel, IServerResponse } from 'client/app/app.view-models';

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
    public justCheckedCreateUserForm: boolean = false;

    constructor(private commonUxService: CommonUXService, private userService: UserManagementService) {
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
      this.userService.logIn();
    }
    public logOut(): void {
      this.userService.logOut();
    }
    public createUser(): void {
      if (!this.validateNewUser()) {
        this.commonUxService.showWarningToast('Please address highlighted errors.');
        return;
      }

      this.userService.createUser(this.newUser).subscribe(
        (res: IServerResponse) => {
            if (res.success) {
              this.commonUxService.showSuccessToast('Congratulations! Your account has been created.');
              this.newUser = {} as IUserViewModel;
              this.showLoginForm = true;
              this.justCheckedCreateUserForm = false;
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

    private validateNewUser(): boolean {
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
      this.justCheckedCreateUserForm = true;
      return isValid;
    }
}
