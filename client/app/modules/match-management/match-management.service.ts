import { Injectable, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, tap, map } from 'rxjs/operators';
import { IMatchViewModel, IServerResponse, IUserViewModel, IUserCharacterViewModel } from '../../app.view-models';

@Injectable()
export class MatchManagementService {

    public cachedMatches: BehaviorSubject<IMatchViewModel[]> = new BehaviorSubject<IMatchViewModel[]>(null);

    constructor(private httpClient: HttpClient, @Inject('MatchApiUrl') private apiUrl: string) {
    }

    public loadAllMatches(): void {
        this.httpClient.get<IServerResponse>(`${this.apiUrl}/getall`).subscribe(
            (res: IServerResponse) => {
                if (res && res.data) {
                    this._updateCachedMatches(res.data.matches);
                }
            }
        );
    }
    public createMatch(match: IMatchViewModel): Observable<IMatchViewModel> {
        match = this._prepareMatchForApi(match);
        return this.httpClient.post(`${this.apiUrl}/create`, match).pipe(
            map((res: IServerResponse) => {
                if (res && res.data && res.data.match) {
                    const serverMatch: IMatchViewModel = res.data.match;
                    // Set "isNew" property for highlighting
                    serverMatch.isNew = true;
                    let allMatches: IMatchViewModel[] = this.cachedMatches.value;
                    allMatches.push(serverMatch);
                    this._updateCachedMatches(allMatches);

                    // Remove "isNew" property after 3 seconds so it won't re-highlight.
                    // The subscribers have no way of knowing if a match is new or old
                    // if they're not the component creating matches, so this service
                    // has to keep track of that info for them
                    setTimeout(() => {
                        allMatches = this.cachedMatches.value;
                        const latestMatch: IMatchViewModel = allMatches.find(m => m.matchId === res.data.match.matchId);
                        latestMatch.isNew = false;
                        this._updateCachedMatches(allMatches);
                    }, 3000);

                    return serverMatch;
                }
            }));
    }
    public updateMatch(updatedMatch: IMatchViewModel): Observable<IMatchViewModel> {
        updatedMatch = this._prepareMatchForApi(updatedMatch);
        return this.httpClient.post(`${this.apiUrl}/update`, updatedMatch).pipe(
            map((res: IServerResponse) => {
                if (res && res.data && res.data.match) {
                    const serverMatch: IMatchViewModel = res.data.match;
                    const allMatches: IMatchViewModel[] = this.cachedMatches.value;
                    const index = allMatches.findIndex(m => m.matchId === serverMatch.matchId);
                    allMatches[index] = serverMatch;
                    this._updateCachedMatches(allMatches);
                    return serverMatch;
                } else {
                    return null;
                }
            })
        );
    }
    public deleteMatch(match: IMatchViewModel): void {
        this.httpClient.post(`${this.apiUrl}/delete`, match).pipe(
            tap((res: IServerResponse) => {
                if (res && res.success) {
                    const allMatches: IMatchViewModel[] = this.cachedMatches.value;
                    const index = allMatches.findIndex(m => m.matchId === match.matchId);
                    allMatches.splice(index, 1);
                    this._updateCachedMatches(allMatches);
                }
            })
        ).subscribe();
    }
    public updateCachedMatchesWithUserName(updatedUser: IUserViewModel): void {
        const allMatches: IMatchViewModel[] = this.cachedMatches.value;
        allMatches.filter(m => m.userId === updatedUser.userId).forEach(m => {
            m.userName = updatedUser.userName;
        });
        this._updateCachedMatches(allMatches);
    }
    public updateCachedMatchesWithAltCostume(updatedUser: IUserViewModel): void {
        const allMatches: IMatchViewModel[] = this.cachedMatches.value;
        allMatches.filter(m => m.userId === updatedUser.userId).forEach(m => {
            const matchingUserChar: IUserCharacterViewModel = updatedUser.userCharacters.find(c => c.characterId === m.userCharacterId);
            if (matchingUserChar) {
                m.altCostume = matchingUserChar.altCostume;
            }
        });
        this._updateCachedMatches(allMatches);
    }


    /*-----------------------
         Private helpers
    ------------------------*/
    private _prepareMatchForApi(match: IMatchViewModel): IMatchViewModel {
        // Do all type conversions & other misc translations here before sending to API
        if (match.userCharacterGsp) {
            match.userCharacterGsp = parseInt(match.userCharacterGsp.toString().replace(/\D/g, ''), 10);
        }
        if (match.opponentCharacterGsp) {
            match.opponentCharacterGsp = parseInt(match.opponentCharacterGsp.toString().replace(/\D/g, ''), 10);
        }
        if (match.userCharacterGsp === '') {
            match.userCharacterGsp = null;
        }
        if (match.opponentCharacterGsp === '') {
            match.opponentCharacterGsp = null;
        }
        return match;
    }
    private _updateCachedMatches(updatedMatches: IMatchViewModel[]): void {
        this.cachedMatches.next(updatedMatches);
        this.cachedMatches.pipe(
            publish(),
            refCount()
        );
    }
}
