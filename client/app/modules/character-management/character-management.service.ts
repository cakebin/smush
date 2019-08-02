import { Injectable, Inject  } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, tap, finalize, retryWhen, delay, take } from 'rxjs/operators';
import { ICharacterViewModel } from '../../app.view-models';

@Injectable()
export class CharacterManagementService {

    public characters: BehaviorSubject<ICharacterViewModel[]> = new BehaviorSubject<ICharacterViewModel[]>(null);

    constructor(
        private httpClient: HttpClient,
        @Inject('CharacterApiUrl') private apiUrl: string,
    ) {
        this.loadAllCharacters();
    }
    public loadAllCharacters(): void {
        this.httpClient.get<ICharacterViewModel[]>(`${this.apiUrl}/getall`).pipe(
            // Retry in case we're attempting to get data when the user is still being re-authed
            retryWhen(errors => errors.pipe(delay(1000), take(3)))
        ).subscribe(
            res => {
                this.characters.next(res);
                this.characters.pipe(
                    publish(),
                    refCount()
                );
            }
        );
    }
    public createCharacter(char: ICharacterViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/create`, char);
    }
    public updateCharacter(updatedChar: ICharacterViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/update`, updatedChar);
    }
    public deleteCharacter(charId: number): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/delete`, charId);
    }
}
