import { Component, Input, Output, EventEmitter } from '@angular/core';
import { NgbTypeaheadSelectItemEvent } from '@ng-bootstrap/ng-bootstrap';
import { Observable } from 'rxjs';
import { debounceTime, distinctUntilChanged, map } from 'rxjs/operators';


@Component({
  selector: 'common-ux-typeahead',
  templateUrl: './typeahead.component.html'
})
export class TypeaheadComponent {
  @Input() clearOnSelect: boolean = false;
  @Input() size: '' | 'sm' | 'lg' = '';
  @Input() textPropertyName: string = '';
  @Input() valuePropertyName: string = '';
  @Output() selectItem: EventEmitter<any> = new EventEmitter<any>();
  @Input() set items(itemValues: any[]) {
    this._items = itemValues;
    // Because the items are sometimes delayed, we need to reselect any defaults we might have
    // that had to be skipped when items weren't provided
    this.value = this._lastValue;
  }
  get items(): any[] {
    return this._items;
  }
  @Input() set value(value: any) {
    this._lastValue = value;

    // Either setting it to null
    if (this.items == null || value == null) {
      this.selectedItem = null;
      return;
    }
    // Or giving it a valid value
    this.selectedItem = this.items.find(i => {
      return i[this.valuePropertyName] === value;
    });
  }
  get value(): any {
    if (this.selectedItem) {
      return this.selectedItem[this.valuePropertyName];
    } else {
      return null;
    }
  }

  private _items: any[] = [];
  private _lastValue: any;
  public selectedItem: any;

  constructor() { }

  public itemFormatter = (result: any) => result[this.textPropertyName];
  public search = (text$: Observable<string>) =>
    text$.pipe(
      debounceTime(200),
      distinctUntilChanged(),
      map(term => {
        if (term.length < 1) {
          return [];
        } else {
          return this.items.filter(v => {
            return v[this.textPropertyName].toLowerCase().indexOf(term.toLowerCase()) > -1;
          }).slice(0, 10);
        }
      })
    )

  public onBlur() {
    // If the user has cleared the input and blurred out, we need to output a blank value manually
    // because the typeahead does not recognise this as an input "event" per se
    if (this.selectedItem == null) {
      this.selectItem.emit(null);
    }
  }
  public onSelect(eventObject: NgbTypeaheadSelectItemEvent, input: any): void {
    if (this.clearOnSelect) {
      // https://stackoverflow.com/questions/39783936/how-to-clear-the-typeahead-input-after-a-result-is-selected
      eventObject.preventDefault();
      input.value = '';
      this.selectedItem = null;
    }
    this.selectItem.emit(eventObject.item);
  }
}
