import { Injectable, Inject  } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, retryWhen, delay, take, tap } from 'rxjs/operators';
import { ICharacterViewModel, IServerResponse } from '../../app.view-models';

@Injectable()
export class CharacterManagementService {

    public cachedCharacters: BehaviorSubject<ICharacterViewModel[]> = new BehaviorSubject<ICharacterViewModel[]>(null);

    constructor(
        private httpClient: HttpClient,
        @Inject('CharacterApiUrl') private apiUrl: string,
    ) {
    }
    public loadAllCharacters(): void {
        this.httpClient.get<IServerResponse>(`${this.apiUrl}/getall`).subscribe(
            (res: IServerResponse) => {
                if (res && res.data && res.data.characters) {
                    this.cachedCharacters.next(res.data.characters);
                    this.cachedCharacters.pipe(
                        publish(),
                        refCount()
                    );
                }
            }
        );
    }
    public createCharacter(char: ICharacterViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/create`, char).pipe(
            tap((res: IServerResponse) => {
                if (res && res.data && res.data.character) {
                    const allCharacters: ICharacterViewModel[] = this.cachedCharacters.value;
                    allCharacters.push(res.data.character);
                    this.cachedCharacters.next(allCharacters);
                    this.cachedCharacters.pipe(
                        publish(),
                        refCount()
                    );
                }
            })
        );
    }
    public updateCharacter(updatedChar: ICharacterViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/update`, updatedChar).pipe(
            tap((res: IServerResponse) => {
                if (res && res.data && res.data.character) {
                    // Replace old character with updated character in a copy of cached characters
                    const updatedCharacterFromServer = res.data.character;
                    const allCharacters: ICharacterViewModel[] = this.cachedCharacters.value;
                    const charIndex: number = allCharacters.findIndex(
                        c => c.characterId === updatedCharacterFromServer.characterId);
                    Object.assign(allCharacters[charIndex], updatedCharacterFromServer);

                    // Overwrite cache with updated copy
                    this.cachedCharacters.next(allCharacters);
                    this.cachedCharacters.pipe(
                        publish(),
                        refCount()
                    );
                }
            })
        );
    }
}
