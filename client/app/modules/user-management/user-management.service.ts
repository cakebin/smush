import { Injectable, Inject  } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, tap, finalize } from 'rxjs/operators';
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
            console.log('Logged out. Server returned:', res);
             // Set cached user to nothing! Then publish the new NOTHINGNESS!
            this.cachedUser.next(null);
            this.cachedUser.pipe(publish(), refCount());
            // Send the user back to the home page
            this.router.navigate(['/home']);
        });
    }
    public createUser(user: IUserViewModel): Observable<{}> {
        console.log('Creating user. User model:', user);
        return this.httpClient.post(`${this.authApiUrl}/register`, user).pipe(
            tap(res => {
                console.log('createUser: Done creating user. Server returned:', res);
            }));
    }
    public updateUser(updatedUser: IUserViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/update`, updatedUser).pipe(
            tap(res => {
                console.log('updateUser: Done updating user. Server returned:', res);
            }
        ),
        finalize(() => this._loadUser(updatedUser))
        );
    }
    public deleteUser(userId: number): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/delete`, userId);
    }
}
