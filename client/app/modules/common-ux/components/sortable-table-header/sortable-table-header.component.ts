import {Component, Input, Output, EventEmitter, ViewChild } from '@angular/core';
import { faSortUp, faSortDown } from '@fortawesome/free-solid-svg-icons';
import { ISortEvent } from '../../common-ux.view-models';
import { SortableTableHeaderDirective } from '../../directives/sortable-table-header.directive';


@Component({
  selector: 'common-ux-sortable-table-header',
  template: `
    <span [sortable]="propertyName" (sort)="onSort($event)">
        {{displayName}}
        <fa-layers [fixedWidth]="true">
            <fa-icon sortUp [icon]="faSortUp"></fa-icon>
            <fa-icon sortDown [icon]="faSortDown"></fa-icon>
        </fa-layers>
    </span>
  `
})
export class SortableTableHeaderComponent {
    @Input() propertyName: string = "";
    @Input() displayName: string = "";
    @Output() sort = new EventEmitter<ISortEvent>();
    
    @ViewChild(SortableTableHeaderDirective, { static: false }) header: SortableTableHeaderDirective;
    
    // Icons need to be individually imported
    public faSortUp = faSortUp;
    public faSortDown = faSortDown;
    
    // Pass the event through to the table component (which will do the actual sorting)
    public onSort(event:ISortEvent){
        this.sort.emit(event);
    }
    public clearDirection(){
        this.header.clearDirection();
    }
}