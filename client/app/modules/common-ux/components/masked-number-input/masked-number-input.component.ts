import { Component, Input, Output, EventEmitter, ViewChild, ElementRef, OnInit } from '@angular/core';


@Component({
  selector: 'common-ux-masked-number-input',
  templateUrl: './masked-number-input.component.html'
})
export class MaskedNumberInputComponent implements OnInit {
  @ViewChild('inputField', { static: false }) inputField: ElementRef;

  private _numberValue: string = '';
  @Input() size: '' | 'sm' | 'lg' = '';
  @Input() set numberValue(value: string) {
    this._numberValue = value;
  }
  get numberValue(): string {
    return this._numberValue;
  }
  @Output() numberValueChange: EventEmitter<string> = new EventEmitter<string>();

  ngOnInit() {
    // Format once on init to show a pretty default value
    this.numberValue = this._formatNumber(this.numberValue);
  }

  public checkKeyInput(event: KeyboardEvent): boolean {
    const isNumber = event.key.replace(/\D/g, '').length;
    const isNav = ['Backspace', 'Delete', 'Left', 'Right', 'ArrowLeft', 'ArrowRight'].indexOf(event.key) !== -1;
    // To-do: allow appcommands like cmd+[r, c, v...]
    if (isNav || isNumber) {
      return true;
    } else {
      return false;
    }
  }

  public formatAndEmitValue(): void {
    this.numberValue = this._formatNumber(this.numberValue);
    this.numberValueChange.emit(this.numberValue);
  }

  private _formatNumber(val: string | number): string {
    let newStringValue: string = '';

    if (typeof val === 'string') {
      newStringValue = val;
    } else if (typeof val === 'number') {
      newStringValue = val.toString();
    }

    newStringValue = newStringValue.replace(/\D/g, '');
    if (newStringValue !== '') {
      return parseInt(newStringValue, 10).toLocaleString();
    } else {
      return newStringValue;
    }
  }
}
