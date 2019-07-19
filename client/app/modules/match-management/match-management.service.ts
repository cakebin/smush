import { Injectable, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { IMatchViewModel } from '../../app.view-models';


@Injectable()
export class MatchManagementService {
    constructor(private httpClient: HttpClient, @Inject('ApiUrl') private apiUrl: string) {
    }

    public getAllMatches(): Observable<IMatchViewModel[]> {
        return this.httpClient.get<IMatchViewModel[]>(`${this.apiUrl}/getall`);
    }
    public createMatch(match: IMatchViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/create`, match);
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
