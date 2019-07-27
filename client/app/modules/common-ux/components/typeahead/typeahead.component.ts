import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { NgbTypeaheadSelectItemEvent } from '@ng-bootstrap/ng-bootstrap';
import { Observable } from 'rxjs';
import { debounceTime, distinctUntilChanged, map } from 'rxjs/operators';


@Component({
  selector: 'common-ux-typeahead',
  templateUrl: './typeahead.component.html'
})
export class TypeaheadComponent implements OnInit {
  @Input() items: string[] = [];
  @Output() selectItem: EventEmitter<string> = new EventEmitter<string>();
  @Input() set defaultItem(value: string) {
    // Either setting it to null, or giving it a valid value
    if (!value || this.items.indexOf(value) !== -1) {
      this._selectedValue = value;
      this.inputValue = value;
    }
  }
  get defaultItem(): string {
    return this._defaultItem;
  }
  private _defaultItem: string = '';
  private _selectedValue: string;
  public inputValue: string = '';

  constructor() { }

  ngOnInit() {
  }

  search = (text$: Observable<string>) =>
    text$.pipe(
      debounceTime(200),
      distinctUntilChanged(),
      map(term => {
        if (term.length < 1) {
          return [];
        } else {
          return this.items.filter(v => {
            return v.toLowerCase().indexOf(term.toLowerCase()) > -1;
          }).slice(0, 10);
        }
      })
    )

  public onBlur() {
    // If the user has cleared the input and blurred out, we need to output a blank value manually
    // because the typeahead does not recognise this as an input "event" per se
    if (this.inputValue === '') {
      this.selectItem.emit('');
    } else if (this._selectedValue) {
      this.selectItem.emit(this._selectedValue);
    }
  }
  public onSelect(eventObject: NgbTypeaheadSelectItemEvent): void {
    this._selectedValue = eventObject.item;
    this.selectItem.emit(eventObject.item);
  }
  public clear(): void {
    this.inputValue = '';
    this._selectedValue = '';
  }

}
