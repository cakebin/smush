import { Component, OnInit } from '@angular/core';

import { IUserViewModel } from '../../app.view-models';
import { CommonUXService } from '../common-ux/common-ux.service';
import { UserManagementService } from './user-management.service';

@Component({
  selector: 'user-view',
  templateUrl: './user-view.component.html',
})
export class UserViewComponent implements OnInit {
  public user: IUserViewModel;

  constructor(
    private commonUXService:CommonUXService,
    private userManagementService: UserManagementService,
    ){
  }

  ngOnInit() {
  
  }
}