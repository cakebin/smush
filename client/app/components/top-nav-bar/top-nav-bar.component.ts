import { Component, OnInit } from '@angular/core';
import { UserManagementService } from '../../modules/user-management/user-management.service';
import { CommonUXService } from '../../modules/common-ux/common-ux.service';
import { IUserViewModel } from 'client/app/app.view-models';

@Component({
  selector: 'top-nav-bar',
  templateUrl: './top-nav-bar.component.html',
  styleUrls: ['./top-nav-bar.component.css']
})
export class TopNavBarComponent implements OnInit {
    public user: IUserViewModel;

    constructor(private commonUXService: CommonUXService, private userService: UserManagementService){
    }

    ngOnInit() {
      // Subscribe to whatever user will end up logging in at some point
      this.userService.cachedUser.subscribe({
        next: res => {
          this.user = res;
        },
        error: err => {
          this.commonUXService.showDangerToast('Unable to get user data.');
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
}
