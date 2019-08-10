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
        this.httpClient.get<IServerResponse>(`${this.apiUrl}/getall`).subscribe(
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
                    // Set "isNew" property for highlighting
                    res.data.match.isNew = true;

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
    public updateMatch(updatedMatch: IMatchViewModel): Observable<IMatchViewModel> {
        if (updatedMatch.userCharacterGsp) {
            updatedMatch.userCharacterGsp = parseInt(updatedMatch.userCharacterGsp.toString().replace(/\D/g, ''), 10);
        }
        if (updatedMatch.opponentCharacterGsp) {
            updatedMatch.opponentCharacterGsp = parseInt(updatedMatch.opponentCharacterGsp.toString().replace(/\D/g, ''), 10);
        }

        return this.httpClient.post(`${this.apiUrl}/update`, updatedMatch).pipe(
            map((res: IServerResponse) => {
                if (res && res.data && res.data.match) {
                    const serverMatch: IMatchViewModel = res.data.match;
                    const allMatches: IMatchViewModel[] = this.cachedMatches.value;
                    const index = allMatches.findIndex(m => m.matchId === serverMatch.matchId);

                    allMatches[index] = serverMatch;
                    this.cachedMatches.next(allMatches);
                    this.cachedMatches.pipe(
                        publish(),
                        refCount()
                    );
                    return serverMatch;
                } else {
                    return null;
                }
            })
        );
    }
    public deleteMatch(matchId: number): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/delete`, matchId);
    }
}
