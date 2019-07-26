import { Injectable } from '@angular/core';
import { UserManagementService } from './modules/user-management/user-management.service';
import { IUserViewModel } from './app.view-models';
import { Router } from '@angular/router';

@Injectable()
export class AuthGuardService {
    private _user: IUserViewModel;

    constructor(
        private userService: UserManagementService,
        private router: Router
    ) {
        this.userService.cachedUser.subscribe(
            res => {
                this._user = res;
            }
        );
    }

    // This is a method that needs to exist on an Auth Guard
    public canActivate(): boolean {
        // If the user is logged in, they can proceed
        if (this._isAuthenticated()) {
            return true;
        } else {
        // If the user is not logged in, send them to the home page
            this.router.navigate(['/home']);
        }
    }

    private _isAuthenticated(): boolean {
        return (this._user != null);
    }
}
