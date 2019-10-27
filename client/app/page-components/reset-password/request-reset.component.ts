import { Component } from '@angular/core';
import { UserManagementService } from 'client/app/modules/user-management/user-management.service';

@Component({
  selector: 'request-reset',
  templateUrl: './request-reset.component.html',
})
export class RequestResetComponent {
  public emailAddress: string = '';
  public isSending: boolean = false;
  public requestSent: boolean = false;

  constructor(private userService: UserManagementService) {
  }

  public sendRequest() {
    if (!this.emailAddress) {
      return;
    }
    this.isSending = true;
    this.userService.requestResetPassword(this.emailAddress).subscribe(
      res => {
        this.requestSent = true;
        this.isSending = false;
      },
      () => {
        this.isSending = false;
      }
    );
  }
}
