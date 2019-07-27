export interface IMatchViewModel {
    matchId: number;
    userId: number;
    userName: string; // Viewmodel only, intended to be read-only for match display
    opponentCharacterName: string;
    opponentCharacterGsp: number;
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
        public opponentCharacterName: string = '',
        public opponentCharacterGsp: number = null,
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
    userName: string;
    defaultCharacterName: string;
    defaultCharacterGsp: number;
}
export class UserViewModel implements IUserViewModel {
    constructor(
        public userId: number = null,
        public emailAddress: string = '',
        public userName: string = '',
        public defaultCharacterName: string = '',
        public defaultCharacterGsp: number = null,
    ) {
    }
}

