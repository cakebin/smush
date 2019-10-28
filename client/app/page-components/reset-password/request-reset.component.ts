import { Component } from '@angular/core';
import { UserManagementService } from 'client/app/modules/user-management/user-management.service';

@Component({
  selector: 'request-reset',
  templateUrl: './request-reset.component.html',
})
export class RequestResetComponent {
  public emailAddress: string = '';
  public isSending: boolean = false;
  public requestSuccess: boolean = false;
  public requestFailed: boolean = false;

  constructor(private userService: UserManagementService) {
  }

  public sendRequest() {
    if (!this.emailAddress) {
      return;
    }
    this.isSending = true;
    this.userService.requestResetPassword(this.emailAddress).subscribe(
      (res: boolean) => {
        if (res) {
          this.requestSuccess = true;
        } else {
          this.requestFailed = true;
        }
      },
      error => {
        this.requestFailed = true;
      },
      () => {
        this.isSending = false;
      }
    );
  }
}
