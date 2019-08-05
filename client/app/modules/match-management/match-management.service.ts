import { Injectable, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, tap, map, finalize, retryWhen, delay, take } from 'rxjs/operators';
import { IMatchViewModel, IServerResponse } from '../../app.view-models';

@Injectable()
export class MatchManagementService {

    public cachedMatches: BehaviorSubject<IMatchViewModel[]> = new BehaviorSubject<IMatchViewModel[]>(null);

    constructor(private httpClient: HttpClient, @Inject('MatchApiUrl') private apiUrl: string) {
    }

    public loadAllMatches(): void {
        this.httpClient.get<IServerResponse>(`${this.apiUrl}/getall`).pipe(
            // Retry that probably won't be necessary... it's here in case
            // https://stackoverflow.com/questions/44979131/rxjs-retry-with-delay-function
            retryWhen(errors => errors.pipe(delay(1000), take(3)))
        ).subscribe(
            (res: IServerResponse) => {
                if (res && res.data) {
                    this.cachedMatches.next(res.data.matches);
                    this.cachedMatches.pipe(
                        publish(),
                        refCount()
                    );
                }
            }
        );
    }
    public createMatch(match: IMatchViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/create`, match).pipe(
            tap((res: IServerResponse) => {
                if (res && res.data && res.data.match) {
                    const allMatches: IMatchViewModel[] = this.cachedMatches.value;
                    allMatches.push(res.data.match);
                    this.cachedMatches.next(allMatches);
                    this.cachedMatches.pipe(
                        publish(),
                        refCount()
                    );
                }
            }));
    }
    public updateMatch(updatedMatch: IMatchViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/update`, updatedMatch);
    }
    public deleteMatch(matchId: number): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/delete`, matchId);
    }
}
