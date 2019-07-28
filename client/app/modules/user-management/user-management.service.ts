import { Injectable, Inject  } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, tap, finalize } from 'rxjs/operators';
import { IUserViewModel } from '../../app.view-models';

@Injectable()
export class UserManagementService {
    public cachedUser: BehaviorSubject<IUserViewModel> = new BehaviorSubject<IUserViewModel>(null);

    constructor(
        private httpClient: HttpClient,
        private router: Router,
        @Inject('UserApiUrl') private apiUrl: string,
    ) {
    }

    private _loadUser(userId: number): void {
        this.httpClient.get<IUserViewModel>(`${this.apiUrl}/get/${userId}`).subscribe(
            res => {
                this.cachedUser.next(res);
                this.cachedUser.pipe(
                    publish(),
                    refCount()
                );
            }
        );
    }
    public logIn(): void {
        // Placeholder for actually logging in
        this._loadUser(1);
    }
    public logOut(): void {
        // Set cached user to nothing! Then publish the new NOTHINGNESS!
        this.cachedUser.next(null);
        this.cachedUser.pipe(publish(), refCount());
        // Send the user back to the home page
        this.router.navigate(['/home']);
    }
    public getUser(): BehaviorSubject<IUserViewModel> {
        if (!this.cachedUser.value) {
            this._loadUser(1);
        }
        return this.cachedUser;
    }
    public createUser(user: IUserViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/create`, user);
    }
    public updateUser(updatedUser: IUserViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/update`, updatedUser).pipe(
            tap(res => {
                console.log('updateUser: Done updating user. Server returned:', res);
            }
        ),
        finalize(() => this._loadUser(updatedUser.userId))
        );
    }
    public deleteUser(userId: number): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/delete`, userId);
    }
}
