import { Component, Input, Output, EventEmitter } from '@angular/core';


@Component({
  selector: 'common-ux-masked-number-input',
  templateUrl: './masked-number-input.component.html'
})
export class MaskedNumberInputComponent {
  private _numberValue: string = '';
  @Input() set numberValue(value: string) {
    this._numberValue = this._formatNumber(value);
  }
  get numberValue(): string {
    return this._numberValue;
  }
  @Output() numberValueChange: EventEmitter<string> = new EventEmitter<string>();

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

  public formatNumberAndEmit(event: KeyboardEvent): void {
    // Skip for arrow keys
    if (['Left', 'Right', 'ArrowLeft', 'ArrowRight'].indexOf(event.key) !== -1) {
      return;
    }
    this.numberValueChange.emit(this.numberValue);
  }

  private _formatNumber(val: string): string {
    let newValue: string = val;
    newValue = val.replace(/\D/g, '');
    if (newValue !== '') {
      return parseInt(newValue, 10).toLocaleString();
    } else {
      return val;
    }
  }
}
