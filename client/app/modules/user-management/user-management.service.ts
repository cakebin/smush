import { Injectable, Inject  } from '@angular/core';
import { Router } from '@angular/router';
import { CommonUxService } from '../../modules/common-ux/common-ux.service';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, tap, map } from 'rxjs/operators';
import { IUserViewModel, LogInViewModel, IServerResponse, IUserCharacterViewModel, IChartUserViewModel } from '../../app.view-models';

@Injectable()
export class UserManagementService {
    // This is a Timer but we can't easily import that type
    private _checkSessionInterval: any;
    // The main cachedUser object that all pages are subscribed to
    public cachedUser: BehaviorSubject<IUserViewModel> = new BehaviorSubject<IUserViewModel>(null);

    constructor(
        private httpClient: HttpClient,
        private router: Router,
        private commonUxService: CommonUxService,
        @Inject('UserApiUrl') private apiUrl: string,
        @Inject('AuthApiUrl') private authApiUrl: string,
        @Inject('UserCharacterApiUrl') private userCharacterApiUrl: string,
    ) {
        // Start interval for login check (runs once a minute)
        this._startIntervalSessionCheck();
    }


    /*-----------------------
               Auth
    ------------------------*/

    public logIn(logInModel: LogInViewModel): Observable<IServerResponse> {
        return this.httpClient.post(`${this.authApiUrl}/login`, logInModel)
        .pipe(
            tap((res: IServerResponse) => {
                if (res.success && res.data) {
                    const user: IUserViewModel = res.data.user;
                    user.userCharacters = res.data.userCharacters;
                    localStorage.setItem('smush_user', JSON.stringify(user));
                    localStorage.setItem('smush_access_expire', JSON.stringify(new Date(res.data.accessExpiration)));
                    localStorage.setItem('smush_refresh_expire', JSON.stringify(new Date(res.data.refreshExpiration)));
                    // Don't save authenticated state
                    user.isAuthenticated = true;
                    this._updateCachedUser(user);
                }
            })
        );
    }
    public logOut(): void {
        this._clearLocalStorage();
        this.router.navigate(['/home']);

        if (this.cachedUser.value) {
            this.httpClient.post(`${this.authApiUrl}/logout`, this.cachedUser.value)
            .subscribe(res => {
                this.cachedUser.next(null);
                this.cachedUser.pipe(publish(), refCount());
            });
        }
    }


    /*-----------------------
               User
    ------------------------*/

    public getAllUsers(): Observable<IChartUserViewModel[]> {
        return this.httpClient.get(`${this.apiUrl}/getall`).pipe(
            map((res: IServerResponse) => {
                if (res.success && res.data) {
                    return res.data.users as IChartUserViewModel[];
                }
            })
        );
    }
    public createUser(user: IUserViewModel): Observable<{}> {
        return this.httpClient.post(`${this.authApiUrl}/register`, user);
    }
    public updateUser(updatedUser: IUserViewModel): Observable<{}> {
        updatedUser = this._prepareUserForApi(updatedUser);
        return this.httpClient.post(`${this.apiUrl}/update_profile`, updatedUser).pipe(
            tap((res: IServerResponse) => {
                if (res && res.data && res.data.user) {
                    res.data.user.userCharacters = res.data.userCharacters;
                    this._updateCachedUser(res.data.user, true);
                }
            })
        );
    }
    public deleteUser(userId: number): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/delete`, userId);
    }


    /*-----------------------
         User characters
    ------------------------*/

    public createUserCharacter(userCharacter: IUserCharacterViewModel): Observable<{}> {
        userCharacter.userId = this.cachedUser.value.userId;
        userCharacter = this._prepareUserCharacterForApi(userCharacter);
        return this.httpClient.post(`${this.userCharacterApiUrl}/create`, userCharacter).pipe(
            tap((res: IServerResponse) => {
                if (res && res.data && res.data.user) {
                    res.data.user.userCharacters = res.data.userCharacters;
                    this._updateCachedUser(res.data.user, true);
                }
            })
        );
    }
    public updateUserCharacter(userCharacter: IUserCharacterViewModel): Observable<IUserViewModel> {
        userCharacter.userId = this.cachedUser.value.userId;
        userCharacter = this._prepareUserCharacterForApi(userCharacter);
        return this.httpClient.post(`${this.userCharacterApiUrl}/update`, userCharacter).pipe(
            map((res: IServerResponse) => {
                if (res && res.data && res.data.user) {
                    res.data.user.userCharacters = res.data.userCharacters;
                    this._updateCachedUser(res.data.user, true);
                    return res.data.user;
                }
            })
        );
    }
    public setDefaultUserCharacter(userCharacter: IUserCharacterViewModel): void {
        userCharacter.userId = this.cachedUser.value.userId;
        this.httpClient.post(`${this.apiUrl}/update_default_user_character`, userCharacter).pipe(
            tap((res: IServerResponse) => {
                if (res && res.data && res.data.user) {
                    res.data.user.userCharacters = res.data.userCharacters;
                    this._updateCachedUser(res.data.user, true);
                }
            })
        ).subscribe();
    }
    public unsetDefaultUserCharacter(userCharacter: IUserCharacterViewModel): void {
        userCharacter.userId = this.cachedUser.value.userId;
        // If userCharacterId is null, the API will set user's defaultUserCharacterId to null
        userCharacter.userCharacterId = null;
        this.httpClient.post(`${this.apiUrl}/update_default_user_character`, userCharacter).pipe(
            tap((res: IServerResponse) => {
                if (res && res.data && res.data.user) {
                    res.data.user.userCharacters = res.data.userCharacters;
                    this._updateCachedUser(res.data.user, true);
                }
            })
        ).subscribe();
    }
    public deleteUserCharacter(userCharacter: IUserCharacterViewModel): void {
        userCharacter = this._prepareUserCharacterForApi(userCharacter);
        this.httpClient.post(`${this.userCharacterApiUrl}/delete`, userCharacter).pipe(
            tap((res: IServerResponse) => {
                if (res && res.data && res.data.user) {
                    res.data.user.userCharacters = res.data.userCharacters;
                    this._updateCachedUser(res.data.user, true);
                }
            })
        ).subscribe();
    }


    /*-----------------------
         Private helpers
    ------------------------*/

    private _prepareUserForApi(user: IUserViewModel): IUserViewModel {
        // Do all type conversions & other misc translations here before sending to API
        if (user.defaultUserCharacterGsp) {
            user.defaultUserCharacterGsp = parseInt(user.defaultUserCharacterGsp.toString().replace(/\D/g, ''), 10);
        }
        return user;
    }
    private _prepareUserCharacterForApi(userCharacter: IUserCharacterViewModel): IUserCharacterViewModel {
        // Do all type conversions & other misc translations here before sending to API
        if (userCharacter.characterGsp) {
            userCharacter.characterGsp = parseInt(userCharacter.characterGsp.toString().replace(/\D/g, ''), 10);
        }
        return userCharacter;
    }
    private _updateCachedUser(user: IUserViewModel, updateLocalStorage: boolean = false): void {
        if (updateLocalStorage) {
            localStorage.setItem('smush_user', JSON.stringify(user));
        }

        this.cachedUser.next(user);
        this.cachedUser.pipe(
            publish(),
            refCount()
        );
    }
    private _loadUserFromStorage() {
        const savedUserJson = localStorage.getItem('smush_user');
        const savedUser = JSON.parse(savedUserJson);
        if (savedUser) {
            this._updateCachedUser(savedUser);
        }
    }
    private _clearLocalStorage(): void {
        localStorage.removeItem('smush_user');
        localStorage.removeItem('smush_refresh_expire');
        localStorage.removeItem('smush_access_expire');
    }
    private _startIntervalSessionCheck() {
        if (this._checkSessionInterval) {
            clearInterval(this._checkSessionInterval);
        }
        // Run this first to make sure the user gets a cookie if they no longer have one
        this._runSessionCheck(true);
        this._checkSessionInterval = setInterval(() => {
            // Then check again every minute
           this._runSessionCheck();
        }, 60000);
    }
    private _runSessionCheck(isInitialCheck: boolean = false) {
        const dateNow = new Date();
        const refreshExpiration: string = localStorage.getItem('smush_refresh_expire');
        const accessExpiration: string = localStorage.getItem('smush_access_expire');

        if (!refreshExpiration || !accessExpiration) {
            // This user never got and saved a token. Don't try to refresh
            this._clearLocalStorage();
            return;
        }

        // Only run the following if we have an expiration date in localstorage
        const dateNowMs = dateNow.getTime();
        const refreshExpireMs: number = new Date(JSON.parse(refreshExpiration)).getTime();
        const accessExpireMs: number = new Date(JSON.parse(accessExpiration)).getTime();

        if (isInitialCheck) {
            // If this is on pageload, load the user because they still have both a refresh and access token token
            this._loadUserFromStorage();
        }

        if (dateNowMs >= refreshExpireMs) {
            this.logOut();
            // It is after the refresh expiry date. Log user out and don't refresh their token.
            if (!isInitialCheck) {
                this.commonUxService.openConfirmModal(
                    'You\'ve been logged out because your session has expired. Log in again to continue tracking matches :)',
                    'Session Expired',
                    true,
                    'Okey'
                );
            }
        } else {
            // We are still within the refresh range, so check the access expiration and see if we
            // need to refresh it (within 2 min of expiration) or get a new one (if it's gone).
            const accessExpired = dateNowMs > accessExpireMs;
            const accessAboutToExpire = (dateNowMs < accessExpireMs) && (accessExpireMs - dateNowMs < 120000);

            if (accessExpired || accessAboutToExpire) {
                this.httpClient.post(`${this.authApiUrl}/refresh`, this.cachedUser.value).subscribe(
                    (res: IServerResponse) => {
                        if (res && res.success && res.data) {
                            // Set the new updated access expiration date
                            localStorage.setItem('smush_access_expire', JSON.stringify(new Date(res.data.accessExpiration)));
                            if (isInitialCheck) {
                                // If this is on pageload, change isAuthenticated to true so we know to load the static data etc
                                this._setUserAuthenticated();
                            }
                        }
                    }
                );
            } else {
                if (isInitialCheck) {
                    // We didn't need to re-authenticate. The user still has a good access token.
                    // If this is on pageload, change isAuthenticated to true so we know to load the static data etc
                    this._setUserAuthenticated();
                }
            }
        }
    }
    private _setUserAuthenticated(): void {
        const authedUser: IUserViewModel = this.cachedUser.value;
        authedUser.isAuthenticated = true;
        this._updateCachedUser(authedUser);
    }
}
