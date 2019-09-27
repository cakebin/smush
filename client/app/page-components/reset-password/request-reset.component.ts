import { Component } from '@angular/core';

@Component({
  selector: 'request-reset',
  templateUrl: './request-reset.component.html',
})
export class RequestResetComponent {
  public emailAddress: string = '';
  public isSaving: boolean = false;

  constructor() {
  }
}
