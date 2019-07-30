import { Injectable, Inject  } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, tap, finalize, map } from 'rxjs/operators';
import { IUserViewModel, LogInViewModel, IAuthServerResponse } from '../../app.view-models';

@Injectable()
export class UserManagementService {
    public cachedUser: BehaviorSubject<IUserViewModel> = new BehaviorSubject<IUserViewModel>(null);

    constructor(
        private httpClient: HttpClient,
        private router: Router,
        @Inject('UserApiUrl') private apiUrl: string,
        @Inject('AuthApiUrl') private authApiUrl: string,
    ) {
        // Check for existing login
        this.checkLogIn();
    }

    private _loadUser(user: IUserViewModel): void {
        this.cachedUser.next(user);
        this.cachedUser.pipe(
            publish(),
            refCount()
        );
    }
    public logIn(logInModel: LogInViewModel): Observable<IAuthServerResponse> {
        return this.httpClient.post(`${this.authApiUrl}/login`, logInModel)
        .pipe(
            tap((res: IAuthServerResponse) => {
                if (res.success) {
                    localStorage.setItem('smush_user', JSON.stringify(res.user));
                    localStorage.setItem('smush_expire', JSON.stringify(new Date(res.expiration)));
                    this._loadUser(res.user);
                }
            })
        );
    }
    public checkLogIn() {
        const expiration = localStorage.getItem('smush_expire');
        const expiresAt = JSON.parse(expiration);
        if (new Date() > new Date(expiresAt)) {
            // alert('Your login session has expired.');
            // this.logOut();
        } else {
            const savedUserJson = localStorage.getItem('smush_user');
            const savedUser = JSON.parse(savedUserJson);
            if (savedUser) {
                this._loadUser(savedUser);
            }
        }
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
            localStorage.removeItem('smush_expire');
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
}
