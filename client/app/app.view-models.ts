import { NgStyle } from '@angular/common';

export interface IMatchViewModel {
    matchId: number;
    userId: number;
    userName: string; // Viewmodel only, intended to be read-only for match display
    opponentCharacterId: number;
    opponentCharacterName: string;
    opponentCharacterGsp: number;
    userCharacterId: number;
    userCharacterName: string;
    userCharacterGsp: number;
    opponentTeabag: boolean;
    opponentCamp: boolean;
    opponentAwesome: boolean;
    userWin: boolean;
    created: Date; // Set on the server. Read-only
}
export class MatchViewModel implements IMatchViewModel {
    constructor(
        public matchId: number = null,
        public userId: number = null,
        public userName: string = '',
        public opponentCharacterId: number = null,
        public opponentCharacterName: string = '',
        public opponentCharacterGsp: number = null,
        public userCharacterId: number = null,
        public userCharacterName: string = '',
        public userCharacterGsp: number = null,
        public opponentTeabag: boolean = null,
        public opponentCamp: boolean = null,
        public opponentAwesome: boolean = null,
        public userWin: boolean = null,
        public created: Date = null,
    ) {
    }
}
export interface IUserViewModel {
    userId: number;
    emailAddress: string;
    password: string;
    passwordConfirm: string;
    userName: string;
    defaultCharacterId: number;
    defaultCharacterName: string;
    defaultCharacterGsp: number;
    defaultCharacterImageUrl: string; // Join on characters table for this
}
export class UserViewModel implements IUserViewModel {
    constructor(
        public userId: number = null,
        public emailAddress: string = '',
        public password: string = '',
        public passwordConfirm: string = '',
        public userName: string = '',
        public defaultCharacterId: number = null,
        public defaultCharacterName: string = '',
        public defaultCharacterGsp: number = null,
        public defaultCharacterImageUrl: string = '',
    ) {
    }
}
export interface ICharacterViewModel {
    characterId: number;
    characterName: string;
    characterStockImg: string;
    characterImg: string;
    characterArchetype: string;
}
export interface ITypeAheadViewModel {
    text: string;
    value: any;
}
export interface ILogInViewModel {
    emailAddress: string;
    password: string;
}
export class LogInViewModel implements ILogInViewModel {
    constructor(
        public emailAddress: string = '',
        public password: string = '',
    ) {
    }
}
export interface IServerResponse {
    success: boolean;
    error: any;
    data: any;
}

