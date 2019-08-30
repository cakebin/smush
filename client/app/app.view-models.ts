export interface IMatchViewModel {
    matchId: number;
    userId: number;
    userName: string; // Viewmodel only, intended to be read-only for match display
    opponentCharacterId: number;
    opponentCharacterName: string;
    opponentCharacterGsp: number | string;
    opponentCharacterImage: string;
    userCharacterId: number;
    userCharacterName: string;
    userCharacterGsp: number | string;
    userCharacterImage: string;
    matchTags: IMatchTagViewModel[]; // Viewmodel only, from MatchTags join
    userWin: boolean;
    altCostume: number; // Viewmodel only, from UserCharacters join
    created: Date; // Set on the server. Read-only
    editMode: boolean; // Viewmodel only
    isNew: boolean; // Viewmodel only
}
export class MatchViewModel implements IMatchViewModel {
    constructor(
        public matchId: number = null,
        public userId: number = null,
        public userName: string = '',
        public opponentCharacterId: number = null,
        public opponentCharacterName: string = '',
        public opponentCharacterGsp: number = null,
        public opponentCharacterImage: string = '',
        public userCharacterId: number = null,
        public userCharacterName: string = '',
        public userCharacterGsp: number = null,
        public userCharacterImage: string = '',
        public matchTags: IMatchTagViewModel[] = [],
        public userWin: boolean = null,
        public altCostume: number = null,
        public created: Date = null,
        public editMode: boolean = false,
        public isNew: boolean = false,
    ) {
    }
}
export interface IUserViewModel {
    userId: number;
    emailAddress: string;
    password: string;
    passwordConfirm: string;
    userName: string;
    defaultUserCharacterId: number;
    defaultUserCharacterGsp: number;
    defaultCharacterId: number;
    defaultCharacterName: string;
    userCharacters: IUserCharacterViewModel[];
    isAuthenticated: boolean;
    userRoles: IUserRoleViewModel[];
}
export class UserViewModel implements IUserViewModel {
    constructor(
        public userId: number = null,
        public emailAddress: string = '',
        public password: string = '',
        public passwordConfirm: string = '',
        public userName: string = '',
        public defaultUserCharacterId: number = null,
        public defaultUserCharacterGsp: number = null,
        public defaultCharacterId: number = null,
        public defaultCharacterName: string = '',
        public userCharacters: IUserCharacterViewModel[] = [],
        public isAuthenticated: boolean = false,
        public userRoles: IUserRoleViewModel[] = [],
    ) {
    }
}
export interface IUserCharacterViewModel {
    userCharacterId: number;
    characterGsp: number|string;
    altCostume: number;
    characterId: number;
    characterName: string;
    userId: number;
    editMode: boolean;
}
export interface IUserRoleViewModel {
    userRoleId: number;
    userId: number;
    roleId: number;
    roleName: string;
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
export interface ITagViewModel {
    tagId: number;
    tagName: string;
    editMode: boolean;
}
export class TagViewModel implements ITagViewModel {
    constructor(
        public tagId: number = null,
        public tagName: string = '',
        public editMode: boolean = false,
    ) {
    }
}
export interface IMatchTagViewModel {
    matchTagId: number;
    matchId: number;
    tagId: number;
    tagName: string;
}
export interface IChartViewModel {
    chartId: number;
    chartName: string;
  }
export class ChartViewModel implements IChartViewModel {
    constructor(
      public chartId: number,
      public chartName: string,
      public chartType: 'bar' | 'line',
    ) {}
}
export interface IChartUserViewModel {
    userId: number;
    userName: string;
}
