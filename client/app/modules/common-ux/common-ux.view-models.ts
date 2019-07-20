// For types, interfaces, and classes only!
// Methods won't work in here. Put those in the commonux service instead.

export type SortDirection = 'asc' | 'desc' | '';
export interface ISortEvent {
    column: string;
    direction: SortDirection;
}
export interface IHeaderViewModel {
    propertyName: string;
    displayName: string;
}
export class HeaderViewModel implements IHeaderViewModel {
	constructor(
        public propertyName: string,
        public displayName: string,
    ) {
	}
}