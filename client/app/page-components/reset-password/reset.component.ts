import { Component } from '@angular/core';

@Component({
  selector: 'reset',
  templateUrl: './reset.component.html',
})
export class PasswordResetComponent {
  public newPassword: string = '';
  public newPasswordConfirm: string = '';
  public isSaving: boolean = false;
  public showWarnings: boolean = false;

  constructor() {
  }
}
