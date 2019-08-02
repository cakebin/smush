import { Injectable, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, ReplaySubject } from 'rxjs';
import { publish, refCount, tap, finalize, retryWhen, delay, take } from 'rxjs/operators';
import { IMatchViewModel } from '../../app.view-models';

@Injectable()
export class MatchManagementService {

    public cachedMatches: ReplaySubject<IMatchViewModel[]> = new ReplaySubject<IMatchViewModel[]>();

    constructor(private httpClient: HttpClient, @Inject('MatchApiUrl') private apiUrl: string) {
        this.loadAllMatches();
    }

    public loadAllMatches(): void {
        this.httpClient.get<IMatchViewModel[]>(`${this.apiUrl}/getall`).pipe(
            // Retry in case we're attempting to get matches when the user is still being re-authed
            // https://stackoverflow.com/questions/44979131/rxjs-retry-with-delay-function
            retryWhen(errors => errors.pipe(delay(1000), take(3)))
        ).subscribe(
            res => {
                this.cachedMatches.next(res);
                this.cachedMatches.pipe(
                    publish(),
                    refCount()
                );
            }
        );
    }
    public createMatch(match: IMatchViewModel): Observable<{}> {
        const apiCreateMatch: Observable<{}> = this.httpClient.post(`${this.apiUrl}/create`, match);

        // SAKI: I might not keep this way of subscribing cause it's literally just a basic subscribe
        // in the form of stupid fancy pipe operators. It's so fucking extra.
        // I'm only trying this because I don't like the thought of a subscription to another subscription. TBD...
        return apiCreateMatch.pipe(
                    tap(res => {
                            console.log('createMatch: Done creating match. Server returned:', res);
                        }
                    ),
                    finalize(() => this.loadAllMatches())
                );
    }
    public getMatch(matchId: number): Observable<IMatchViewModel> {
        return this.httpClient.get<IMatchViewModel>(`${this.apiUrl}/get?id=${matchId}`);
    }
    public updateMatch(updatedMatch: IMatchViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/update`, updatedMatch);
    }
    public deleteMatch(matchId: number): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/delete`, matchId);
    }
}
