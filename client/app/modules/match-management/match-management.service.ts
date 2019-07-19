import { Injectable, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { IMatchViewModel, MatchViewModel } from '../../app.view-models';

@Injectable()
export class MatchManagementService {
    private dummyMatches:IMatchViewModel[] = [
        new MatchViewModel(1, "Ice Climbers", 5010032, "Young Link", 4500000, false, false, true, true),
        new MatchViewModel(2, "Captain Falcon", 100190, "Lucina", 98250, false, false, null, false),
        new MatchViewModel(3, "Pikachu", 4510000, "Young Link", 4500000, true, true, null, null),
        new MatchViewModel(4, "Bowser Jr.", 4700200, "Young Link", 4600000, false, false, null, true),
        new MatchViewModel(5, "Pichu", 101000, "Lucina", 98000, false, false, null, null),
        new MatchViewModel(6, "Bowser", 95280, "Lucina", 98500, false, false, null, false),
        new MatchViewModel(7, "Jigglypuff", 99500, "Lucina", 95620, true, true, null, false),
        new MatchViewModel(8, "Falco", 5001500, "Joker", 5125000, false, false, true, null),
        new MatchViewModel(9, "Captain Falcon", 5571000, "Joker", 5125000, false, false, null, true),
        new MatchViewModel(10, "Captain Falcon", 5101400, "Joker", 5001500, true, false, null, false),
        new MatchViewModel(11, "Mario", 10000, "Joker", 5001500, false, true, null, true),
        new MatchViewModel(12, "Zero-Suit Samus", 5001500, "Joker", 5002000, true, true, null, false),
        new MatchViewModel(13, "Ganondorf", 5102000, "Joker", 5002000, true, null, false, false),
        new MatchViewModel(14, "Ganondorf", 4802000, "Marth", 4901500, false, null, true, true),
        new MatchViewModel(15, "Mario", 4982000, "Joker", 5001500, false, false, null, false),
        new MatchViewModel(16, "Mii Brawler", 2500023, "Ike", 2004012, true, true, null, true),
        new MatchViewModel(17, "Falco", 5120000, "Joker", 5002500, false, false, null, true),
        new MatchViewModel(18, "Megaman", 58250, "Wolf", 56000, true, true, null, false),
    ]

    constructor(private httpClient: HttpClient, @Inject('ApiUrl') private apiUrl: string) {
    }

    public getAllMatches(): Observable<IMatchViewModel[]> {
        return of(this.dummyMatches);
        //return this.httpClient.get<IMatchViewModel[]>(`${this.apiUrl}/getall`);
    }
    public createMatch(match: IMatchViewModel): Observable<{}> {
        let lastId:number = 0;
        this.dummyMatches.forEach(m => {if(m.matchId > lastId) lastId = m.matchId});
        match.matchId = lastId + 1;

        this.dummyMatches.push(match);
        return of();
        //return this.httpClient.post(`${this.apiUrl}/create`, match);
    }
    public getMatch(matchId: number): Observable<IMatchViewModel> {
        const result:MatchViewModel = this.dummyMatches.find(m => m.matchId == matchId);
        return of(result);
        //return this.httpClient.get<IMatchViewModel>(`${this.apiUrl}/get?id=${matchId}`);
    }
    public updateMatch(updatedMatch: IMatchViewModel): Observable<{}> {
        const oldMatchId:number = this.dummyMatches.findIndex(m => m.matchId == updatedMatch.matchId);
        this.dummyMatches[oldMatchId] = updatedMatch;

        return of();
        //return this.httpClient.post(`${this.apiUrl}/update`, updatedMatch);
    }
    public deleteMatch(matchId: number): Observable<{}> {
        const matchIndex:number = this.dummyMatches.findIndex(m => m.matchId == matchId);
        this.dummyMatches.splice(matchIndex, 1);

        return of();
        //return this.httpClient.post(`${this.apiUrl}/delete`, matchId);
    }
}
