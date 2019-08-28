import { Injectable } from '@angular/core';
import { UserManagementService } from './modules/user-management/user-management.service';
import { IUserViewModel } from './app.view-models';
import { Router, ActivatedRouteSnapshot } from '@angular/router';

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
    public canActivate(route: ActivatedRouteSnapshot): boolean {
        // If user is not logged in, we automatically reroute to the home page
        if (!this._isAuthenticated()) {
            this.router.navigate(['/home']);
            return false;
        }

        // If a role is required for this route, check it
        const expectedRole = route.data.expectedRole;
        if (expectedRole) {
            if (expectedRole === 'administrator') {
                const hasRoleAdmin = this._user.userRoles.reduce((acc, cur) => {
                    return acc || cur.roleId == 1;
                  }, false);

                if (!this._isAuthenticated() || !hasRoleAdmin) {
                    this.router.navigate(['/home']);
                }
            }
        }

        // If the user is logged in and no roles are required, they can proceed
        return true;
    }

    private _isAuthenticated(): boolean {
        return (this._user != null);
    }
}
