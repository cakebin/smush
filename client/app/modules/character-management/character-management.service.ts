import { Injectable, Inject  } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, tap, finalize, retryWhen, delay, take } from 'rxjs/operators';
import { ICharacterViewModel } from '../../app.view-models';

@Injectable()
export class CharacterManagementService {
    private _characters = ['Bayonetta', 'Bowser', 'Bowser Jr.', 'Captain Falcon', 'Chrom', 'Cloud', 'Corrin',
    'Daisy', 'Dark Samus', 'Diddy Kong', 'Donkey Kong', 'Dr. Mario', 'Duck Hunt', 'Falco', 'Fox', 'Ganondorf',
    'Greninja', 'Ice Climbers', 'Ike', 'Incineroar', 'Inkling', 'Jigglypuff', 'Joker', 'Ken', 'King DeDeDe',
    'King K. Rool', 'Kirby', 'Link', 'Little Mac', 'Lucario', 'Lucas', 'Lucina', 'Luigi', 'Mario', 'Marth',
    'Mega Man', 'Meta Knight', 'Mewtwo', 'Mii Brawler', 'Mii Gunner', 'Mii Sword Fighter', 'Mr. Game & Watch',
    'Ness', 'Olimar', 'Pac-Man', 'Palutena', 'Peach', 'Pichu', 'Pikachu', 'Pit', 'Pokemon Trainer', 'Richter',
    'Ridley', 'Rob', 'Robin', 'Rosalina and Luma', 'Roy', 'Ryu', 'Samus', 'Sheik', 'Shulk', 'Simon', 'Snake',
    'Sonic', 'Toon Link', 'Villager', 'Wario', 'Wolf', 'Yoshi', 'Young Link', 'Wii-Fit Trainer', 'Zelda', 'Zero-Suit Samus'];

    public characters: BehaviorSubject<ICharacterViewModel[]> = new BehaviorSubject<ICharacterViewModel[]>(null);

    constructor(
        private httpClient: HttpClient,
        @Inject('CharacterApiUrl') private apiUrl: string,
    ) {
        this.loadAllCharacters();
    }
    public loadAllCharacters(): void {
        const characterViewModels: ICharacterViewModel[] = this._characters.map((charName: string, index) => {
            return {
                characterId: index,
                characterName: charName,
                characterStockImg: '',
                characterImg: '',
                characterArchetype: ''
            } as ICharacterViewModel;
        });

        this.characters.next(characterViewModels);
        this.characters.pipe(
            publish(),
            refCount()
        );

        /*
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
        */
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
