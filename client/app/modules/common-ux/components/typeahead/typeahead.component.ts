import { Component, Input, Output, EventEmitter } from '@angular/core';
import { NgbTypeaheadSelectItemEvent } from '@ng-bootstrap/ng-bootstrap';
import { Observable } from 'rxjs';
import { debounceTime, distinctUntilChanged, map } from 'rxjs/operators';


@Component({
  selector: 'common-ux-typeahead',
  templateUrl: './typeahead.component.html'
})
export class TypeaheadComponent {
  @Input() textPropertyName: string = '';
  @Input() valuePropertyName: string = '';
  @Input() items: any[] = [];
  @Output() selectItem: EventEmitter<any> = new EventEmitter<any>();
  @Input() set value(value: any) {
    // Either setting it to null, or giving it a valid value
    if (value == null) {
      this._selectedItem = null;
      this.inputDisplayValue = null;
    } else if (this.items.indexOf(value) !== -1) {
      this._selectedItem = this.items.find(i => i[this.valuePropertyName] === value);
      this.inputDisplayValue = this._selectedItem ? this._selectedItem[this.textPropertyName] : null;
    }
  }
  get value(): any {
    if (this._selectedItem) {
      return this._selectedItem[this.valuePropertyName];
    } else {
      return null;
    }
  }
  private _selectedItem: any;
  public inputDisplayValue: string = '';

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
    if (this.inputDisplayValue === '') {
      this.selectItem.emit('');
    } else if (this._selectedItem) {
      this.selectItem.emit(this._selectedItem);
    }
  }
  public onSelect(eventObject: NgbTypeaheadSelectItemEvent): void {
    this._selectedItem = eventObject.item;
    this.selectItem.emit(eventObject.item);
  }
  public clear(): void {
    this.inputDisplayValue = '';
    this._selectedItem = null;
  }


}
