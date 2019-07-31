import { Injectable, Inject  } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, tap, finalize, map } from 'rxjs/operators';
import { IUserViewModel, LogInViewModel, IAuthServerResponse } from '../../app.view-models';

@Injectable()
export class UserManagementService {
    // This is a Timer but we can't easily import that type
    private _checkSessionInterval: any;

    public cachedUser: BehaviorSubject<IUserViewModel> = new BehaviorSubject<IUserViewModel>(null);

    constructor(
        private httpClient: HttpClient,
        private router: Router,
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
                // console.log('updateUser: Done updating user. Server returned:', res);
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
        const expiration = localStorage.getItem('smush_refresh_expire');
        const expiresAt = JSON.parse(expiration);
        if (new Date() <= new Date(expiresAt)) {
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
        this._runSessionCheck();

        this._checkSessionInterval = setInterval(() => {
            // Then check again in a minute
           this._runSessionCheck();
        }, 60000);
    }
    private _runSessionCheck() {
        const dateNow = new Date();
        const refreshExpiration = localStorage.getItem('smush_refresh_expire');
        const accessExpiration: string = localStorage.getItem('smush_access_expire');

        if (!refreshExpiration || !accessExpiration) {
            // This user never got and saved a token. Don't do anything (wait for them to log in)
            return;
        }

        // Only run the following if we have an expiration date in localstorage
        const dateNowMs = dateNow.getTime();
        const refreshExpireMs: number = new Date(JSON.parse(refreshExpiration)).getTime();
        const accessExpireMs: number = new Date(JSON.parse(accessExpiration)).getTime();

        console.log(
            'Checking session:\n\n',
            (new Date(dateNowMs).toLocaleDateString() + ' ' + new Date(dateNowMs).toLocaleTimeString()),
            ' => Date now\n',
            (new Date(accessExpireMs).toLocaleDateString() + ' ' + new Date(accessExpireMs).toLocaleTimeString()),
            ' => Access expire date\n',
            (new Date(refreshExpireMs).toLocaleDateString() + ' ' + new Date(refreshExpireMs).toLocaleTimeString()),
            ' => Refresh expire date',
        );

        if (dateNowMs > refreshExpireMs) {
            // It is after the refresh expiry date. Log user out and don't refresh their token.
            alert('Your session has expired. Please log in again.');
            this.logOut();
        } else {
            // We are still within the refresh range, so check the access expiration and see if we
            // need to refresh it (within 2 min of expiration) or get a new one (if it's gone).
            const accessExpired = dateNowMs > accessExpireMs;
            const accessAboutToExpire = (dateNowMs < accessExpireMs) && (accessExpireMs - dateNowMs < 120000);
            if (accessExpired || accessAboutToExpire) {
                if (accessExpired) {
                    console.warn('Access token has EXPIRED! Refresh token is still valid. Getting a new access token!');
                }
                if (accessAboutToExpire) {
                    console.warn('Access token is about to expire. Updating existing access token!');
                }
                // I was a bit concerned about creating a subscription every three minutes,
                // but it turns out HttpClient destroys subscriptions on completion of the request so memory leaks are not an issue.
                // https://stackoverflow.com/questions/35042929/is-it-necessary-to-unsubscribe-from-observables-created-by-http-methods
                this.httpClient.post(`${this.authApiUrl}/refresh`, this.cachedUser.value).subscribe(
                    (res: IAuthServerResponse) => {
                        console.warn('Api has refreshed token and responded with an updated expire time.\n\n',
                        new Date(accessExpireMs).toLocaleDateString() + ' ' + new Date(accessExpireMs).toLocaleTimeString(),
                        ' => Old access expiration\n',
                        new Date(res.accessExpiration).toLocaleDateString() + ' ' + new Date(res.accessExpiration).toLocaleTimeString(),
                        ' => New access expiration'
                        );

                        // Set the new updated access expiration date
                        localStorage.setItem('smush_access_expire', JSON.stringify(new Date(res.accessExpiration)));
                    });
                }
        }
    }
}
