import { Injectable, Inject,  } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, delay } from 'rxjs/operators';
import { UserViewModel, IUserViewModel } from '../../app.view-models';

@Injectable()
export class UserManagementService {
    private _fakeUser: UserViewModel = new UserViewModel('joebin@gmail.com', 'Jerulfe', 'Joker', 5200000);
    public cachedUser: BehaviorSubject<IUserViewModel> = new BehaviorSubject<IUserViewModel>(null);

    constructor(private httpClient: HttpClient, @Inject('ApiUrl') private apiUrl: string) {
    }

    private _loadUser(): void {
        // this.httpClient.get<IUserViewModel>(`${this.apiUrl}/get/${userId}`)

        const fakeUserObservable: Observable<UserViewModel> = new Observable<UserViewModel>((observer) => {
            observer.next(this._fakeUser);
            observer.complete();
        }).pipe(delay(500));

        fakeUserObservable.subscribe(
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
        this._loadUser();
    }
    public logOut(): void {
        // Set cached user to nothing! Then publish the new NOTHINGNESS!
        this.cachedUser.next(null);
        this.cachedUser.pipe(publish(), refCount());
    }
    public getUser(): BehaviorSubject<IUserViewModel>{
        if(!this.cachedUser.value){
            this._loadUser();
        }
        return this.cachedUser;
    }
    public createUser(user: IUserViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/create`, user);
    }
    public updateUser(updatedUser: IUserViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/update`, updatedUser);
    }
    public deleteUser(userId: number): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/delete`, userId);
    }
}
