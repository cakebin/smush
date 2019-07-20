import { Directive, EventEmitter, Input, Output } from '@angular/core';
import { ISortEvent, SortDirection } from '../common-ux.view-models';

const rotate: {[key: string]: SortDirection} = { 'asc': 'desc', 'desc': '', '': 'asc' };

@Directive({
  selector: 'th span[sortable]',
  host: {
    '[class.asc]': 'direction === "asc"',
    '[class.desc]': 'direction === "desc"',
    '(click)': 'rotate()'
  }
})
export class SortableTableHeaderDirective {

  @Input() sortable: string;
  @Output() sort = new EventEmitter<ISortEvent>();

  // This doesn't need to be an input. It's a state that only this directive keeps track of.
  // Trying to access it by setting it as an input through a grandparent isn't working, so we have a setter below.
  public direction: SortDirection = '';

  public rotate() {
    this.direction = rotate[this.direction];
    this.sort.emit({column: this.sortable, direction: this.direction});
  }
  public clearDirection() {
    this.direction = '';
  }
}