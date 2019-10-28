import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { UserManagementService } from 'client/app/modules/user-management/user-management.service';

@Component({
  selector: 'reset',
  templateUrl: './reset.component.html',
})
export class PasswordResetComponent implements OnInit {
  public token: string = '';
  public newPassword: string = '';
  public newPasswordConfirm: string = '';

  public showWarnings: boolean = false;
  public requestSent: boolean = false;
  public resetSuccessful: boolean = false;
  public resetFailed: boolean = false;
  public tokenExpired: boolean = false;

  constructor(private route: ActivatedRoute, private userService: UserManagementService) {
  }

  ngOnInit() {
    this.route.queryParams.subscribe(params => {
        this.token = params.t;
        if (!params.e || params.e < new Date().getTime()) {
          this.tokenExpired = true;
        }
    });
  }

  public sendReset() {
    if (!this.newPassword || (this.newPassword !== this.newPasswordConfirm)) {
      return;
    }
    this.requestSent = true;
    this.userService.resetPassword(this.token, this.newPassword).subscribe(
      (res: boolean) => {
        if (res) {
          this.resetSuccessful = true;
          setTimeout(() => {
            window.location.href = '/home';
          }, 3000);
        } else {
          this.resetFailed = true;
        }
      },
      error => {
        this.resetFailed = true;
      },
      () => {
        this.newPassword = '';
        this.newPasswordConfirm = '';
      }
    );
  }
}
