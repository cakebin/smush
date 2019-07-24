export class MatchViewModel implements IMatchViewModel {
    constructor(
        public id: number = null,
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
export class UserViewModel implements IUserViewModel {
    constructor(
        public emailAddress: string = '',
        public userName: string = '',
        public defaultCharacterName: string = '',
        public defaultCharacterGsp: number = null,
    ) {
    }
}
export interface IMatchViewModel {
    id: number;
    userName: string;
    opponentCharacterName: string;
    opponentCharacterGsp: number;
    userCharacterName: string;
    userCharacterGsp: number;
    opponentTeabag: boolean;
    opponentCamp: boolean;
    opponentAwesome: boolean;
    userWin: boolean;
    created: Date;
}
export interface IUserViewModel {
    emailAddress: string;
    userName: string;
    defaultCharacterName: string;
    defaultCharacterGsp: number;
}