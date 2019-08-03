import { Injectable, Inject  } from '@angular/core';
import { Router } from '@angular/router';
import { CommonUxService } from '../../modules/common-ux/common-ux.service';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, tap, finalize } from 'rxjs/operators';
import { IUserViewModel, LogInViewModel, IAuthServerResponse } from '../../app.view-models';

@Injectable()
export class UserManagementService {
    // This is a Timer but we can't easily import that type
    private _checkSessionInterval: any;

    public cachedUser: BehaviorSubject<IUserViewModel> = new BehaviorSubject<IUserViewModel>(null);

    constructor(
        private httpClient: HttpClient,
        private router: Router,
        private commonUxService: CommonUxService,
        @Inject('UserApiUrl') private apiUrl: string,
        @Inject('AuthApiUrl') private authApiUrl: string,
    ) {
        // Check and instantiate existing login
        this._onSiteLoad();
        // Start interval for login check (runs once a minute)
        this._startIntervalSessionCheck();
    }

    public logIn(logInModel: LogInViewModel): Observable<IAuthServerResponse> {
        return this.httpClient.post(`${this.authApiUrl}/login`, logInModel)
        .pipe(
            tap((res: IAuthServerResponse) => {
                if (res.success) {
                    localStorage.setItem('smush_user', JSON.stringify(res.user));
                    localStorage.setItem('smush_access_expire', JSON.stringify(new Date(res.accessExpiration)));
                    localStorage.setItem('smush_refresh_expire', JSON.stringify(new Date(res.refreshExpiration)));
                    this._loadUser(res.user);
                }
            })
        );
    }
    public logOut(): void {
        if (!this.cachedUser.value) {
            return;
        }
        this.httpClient.post(`${this.authApiUrl}/logout`, this.cachedUser.value)
        .subscribe(res => {
            this.cachedUser.next(null);
            this.cachedUser.pipe(publish(), refCount());
            localStorage.removeItem('smush_user');
            localStorage.removeItem('smush_refresh_expire');
            localStorage.removeItem('smush_access_expire');
            this.router.navigate(['/home']);
        });
    }
    public createUser(user: IUserViewModel): Observable<{}> {
        return this.httpClient.post(`${this.authApiUrl}/register`, user);
    }
    public updateUser(updatedUser: IUserViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/update`, updatedUser).pipe(
            tap(res => {
                localStorage.setItem('smush_user', JSON.stringify(updatedUser));
            }
        ),
        finalize(() => this._loadUser(updatedUser))
        );
    }
    public deleteUser(userId: number): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/delete`, userId);
    }





    // Private auth handlers
    private _loadUser(user: IUserViewModel): void {
        this.cachedUser.next(user);
        this.cachedUser.pipe(
            publish(),
            refCount()
        );
    }
    private _onSiteLoad() {
        const refreshExpiration = localStorage.getItem('smush_refresh_expire');
        const refreshExpiresAt = JSON.parse(refreshExpiration);
        if (new Date() <= new Date(refreshExpiresAt)) {
            const savedUserJson = localStorage.getItem('smush_user');
            const savedUser = JSON.parse(savedUserJson);
            if (savedUser) {
                this._loadUser(savedUser);
            }
        }
    }
    private _startIntervalSessionCheck() {
        if (this._checkSessionInterval) {
            clearInterval(this._checkSessionInterval);
        }

        // Run this first to make sure the user gets a cookie if they no longer have one
        this._runSessionCheck(true);

        this._checkSessionInterval = setInterval(() => {
            // Then check again in a minute
           this._runSessionCheck();
        }, 60000);
    }
    private _runSessionCheck(isInitialCheck: boolean = false) {
        const dateNow = new Date();
        const refreshExpiration: string = localStorage.getItem('smush_refresh_expire');
        const accessExpiration: string = localStorage.getItem('smush_access_expire');

        if (!refreshExpiration || !accessExpiration) {
            // This user never got and saved a token. Don't do anything (wait for them to log in)
            return;
        }

        // Only run the following if we have an expiration date in localstorage
        const dateNowMs = dateNow.getTime();
        const refreshExpireMs: number = new Date(JSON.parse(refreshExpiration)).getTime();
        const accessExpireMs: number = new Date(JSON.parse(accessExpiration)).getTime();

        if (dateNowMs >= refreshExpireMs) {
            // It is after the refresh expiry date. Log user out and don't refresh their token.
            if (!isInitialCheck) {
                this.commonUxService.openConfirmModal(
                    'You\'ve been logged out because your session has expired. Log in again to continue tracking matches :)',
                    'Session Expired',
                    'Okey'
                );
            }
            this.logOut();
        } else {
            // We are still within the refresh range, so check the access expiration and see if we
            // need to refresh it (within 2 min of expiration) or get a new one (if it's gone).
            const accessExpired = dateNowMs > accessExpireMs;
            const accessAboutToExpire = (dateNowMs < accessExpireMs) && (accessExpireMs - dateNowMs < 120000);
            if (accessExpired || accessAboutToExpire) {
                // I was a bit concerned about creating a subscription every three minutes,
                // but it turns out HttpClient destroys subscriptions on completion of the request so memory leaks are not an issue.
                // https://stackoverflow.com/questions/35042929/is-it-necessary-to-unsubscribe-from-observables-created-by-http-methods
                this.httpClient.post(`${this.authApiUrl}/refresh`, this.cachedUser.value).subscribe(
                    (res: IAuthServerResponse) => {
                        // Set the new updated access expiration date
                        localStorage.setItem('smush_access_expire', JSON.stringify(new Date(res.accessExpiration)));
                    });
                }
        }
    }
}
