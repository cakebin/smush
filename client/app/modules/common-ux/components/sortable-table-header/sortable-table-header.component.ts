import { Component, Input, Output, EventEmitter, ViewChild, HostBinding, OnInit} from '@angular/core';
import { faSortUp, faSortDown } from '@fortawesome/free-solid-svg-icons';
import { ISortEvent, SortDirection } from '../../common-ux.view-models';

const rotate: {[key: string]: SortDirection} = { asc: 'desc', desc: '', '': 'asc' };

// Do not style the parent th here! It won't show up.
@Component({
  selector: '[common-ux-sortable-table-header]',
  template: `
    <span>
        {{displayName}}
        <fa-layers [fixedWidth]="true">
            <fa-icon sortUp [icon]="faSortUp"></fa-icon>
            <fa-icon sortDown [icon]="faSortDown"></fa-icon>
        </fa-layers>
    </span>
  `,
  host: {
    '[class.asc]': 'direction === "asc"',
    '[class.desc]': 'direction === "desc"',
    '(click)': 'rotate()'
  }
})
export class SortableTableHeaderComponent implements OnInit {
    @Input() propertyName: string = '';
    @Input() displayName: string = '';
    @Input() width: string = '';
    @Output() sort = new EventEmitter<ISortEvent>();

    @HostBinding('style.width') widthBinding: string = this.width;

    public direction: SortDirection = '';

    // Icons need to be individually imported
    public faSortUp = faSortUp;
    public faSortDown = faSortDown;

    ngOnInit() {
      this.widthBinding = this.width;
    }

    // Pass the event through to the table component (which will do the actual sorting)
    public onSort(event: ISortEvent) {
        this.sort.emit(event);
    }
    public clearDirection() {
        this.direction = '';
    }
    public rotate() {
        this.direction = rotate[this.direction];
        this.sort.emit({column: this.propertyName, direction: this.direction});
    }
}
